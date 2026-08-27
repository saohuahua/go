// 05_todo_api.go —— 内存 TODO API：把 Handler、JSON、路由和锁串起来
//
// 接口
//
// GET    /health
// GET    /todos
// GET    /todos?id=1
// POST   /todos
// DELETE /todos?id=1
//
// 数据只存在内存中，程序重启后会丢失
// 现在重点是 HTTP 流程，持久化会在 MySQL 和 GORM 阶段处理
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
)

type todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// todoStore 管理共享的 TODO 数据
// HTTP Server 会并发调用 Handler，因此 map 和 nextID 必须加锁
// 前端对照：JS 单线程没有这个问题；Go 里多个 goroutine 同时读写 map 会 panic
// RWMutex = 读锁(R) + 写锁(W)：多个读可并发，写是排他的
type todoStore struct {
	mu     sync.RWMutex // 保护 todos 和 nextID 的读写锁
	nextID int          // 自增 ID 计数器，create 时 +1
	todos  map[int]todo // 真正的数据
}

// newTodoStore 新建空 store，ID 从 1 开始
func newTodoStore() *todoStore {
	return &todoStore{
		nextID: 1,
		todos:  make(map[int]todo),
	}
}

func (s *todoStore) create(input todo) todo {
	// 写操作：整个函数持写锁，保证 nextID++ 和写 map 是一个不可分割的整体
	s.mu.Lock()
	defer s.mu.Unlock()

	input.ID = s.nextID
	s.nextID++
	s.todos[input.ID] = input
	return input
}

func (s *todoStore) list() []todo {
	// 读操作：用 RLock，多个读请求可同时进行、互不阻塞
	s.mu.RLock()
	defer s.mu.RUnlock()

	// map 遍历顺序不固定，按 ID 从小到大填入切片，保证接口输出稳定（前端好断言）
	result := make([]todo, 0, len(s.todos))
	for id := 1; id < s.nextID; id++ {
		if item, ok := s.todos[id]; ok {
			result = append(result, item)
		}
	}
	return result
}

// get 按 id 查一个 TODO
// 查不到时返回 todo 零值 + false——Go 惯例（≈ JS 的 find 找不到返回 undefined）
func (s *todoStore) get(id int) (todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.todos[id]
	return item, ok
}

// delete 删除 TODO，返回是否真的删掉了（不存在时 false）
func (s *todoStore) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	return true
}

// todoAPI 是最小路由器：一个 Handler 里按 path 分流，相当于手写路由表
// 真实项目会由 Gin 根据 method 和 path 帮你分发到不同 Handler，省掉这段 switch
func todoAPI(store *todoStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/health":
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		case r.URL.Path == "/todos":
			handleTodos(w, r, store)
		default:
			writeAPIError(w, http.StatusNotFound, "路由不存在")
		}
	})
}

// handleTodos 按 method 分流到增删查 Handler
// 不支持的 method 回 405 + Allow 头（告诉客户端支持哪些方法）
func handleTodos(w http.ResponseWriter, r *http.Request, store *todoStore) {
	switch r.Method {
	case http.MethodGet:
		handleGetTodos(w, r, store)
	case http.MethodPost:
		handleCreateTodo(w, r, store)
	case http.MethodDelete:
		handleDeleteTodo(w, r, store)
	default:
		w.Header().Set("Allow", "GET, POST, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// handleGetTodos：不带 ?id= 返回全部；带 ?id= 查单个，参数非法或不存在回 400/404
func handleGetTodos(w http.ResponseWriter, r *http.Request, store *todoStore) {
	idText := r.URL.Query().Get("id")
	if idText == "" {
		writeJSON(w, http.StatusOK, store.list())
		return
	}

	id, ok := parseTodoID(w, idText)
	if !ok {
		return
	}
	item, ok := store.get(id)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "TODO 不存在")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// handleCreateTodo：解析 body → 校验 text 非空 → 入库，成功回 201 + 带 id 的完整对象
func handleCreateTodo(w http.ResponseWriter, r *http.Request, store *todoStore) {
	var input todo
	decoder := json.NewDecoder(r.Body)
	//! 解析失败 / text 为空都是客户端输入错误 → 4xx；Server 不因坏输入而 500 或 panic
	if err := decoder.Decode(&input); err != nil {
		writeAPIError(w, http.StatusBadRequest, "JSON 格式错误")
		return
	}
	if strings.TrimSpace(input.Text) == "" {
		writeAPIError(w, http.StatusBadRequest, "text 不能为空")
		return
	}

	created := store.create(todo{Text: strings.TrimSpace(input.Text), Done: input.Done})
	writeJSON(w, http.StatusCreated, created)
}

// handleDeleteTodo：删成功回 204（无 body）；查无此 id 回 404
func handleDeleteTodo(w http.ResponseWriter, r *http.Request, store *todoStore) {
	id, ok := parseTodoID(w, r.URL.Query().Get("id"))
	if !ok {
		return
	}
	if !store.delete(id) {
		writeAPIError(w, http.StatusNotFound, "TODO 不存在")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseTodoID 把 query 里的字符串 id 转成 int 并校验为正整数
// 失败时已写好 400 响应，返回 false，调用方直接 return
func parseTodoID(w http.ResponseWriter, text string) (int, bool) {
	id, err := strconv.Atoi(text)
	if err != nil || id <= 0 {
		writeAPIError(w, http.StatusBadRequest, "id 必须是正整数")
		return 0, false
	}
	return id, true
}

// writeAPIError 统一错误格式：{ "error": "..." }
func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func demoTodoAPI() {
	fmt.Println("\n========== 05 TODO API：最小 HTTP 闭环 ==========")

	api := todoAPI(newTodoStore())
	showTodoRequest(api, http.MethodGet, "/health", "")
	showTodoRequest(api, http.MethodPost, "/todos", `{"text":"学 net/http"}`)
	showTodoRequest(api, http.MethodPost, "/todos", `{"text":"吃 net/http"}`)
	showTodoRequest(api, http.MethodGet, "/todos", "")
	showTodoRequest(api, http.MethodGet, "/todos?id=1", "")
	showTodoRequest(api, http.MethodDelete, "/todos?id=1", "")
	showTodoRequest(api, http.MethodGet, "/todos?id=1", "")
	showTodoRequest(api, http.MethodPost, "/todos", `{"text":"   "}`)

	//! 同一 store 会被多个请求并发访问，因此写操作用 Lock，读操作用 RLock
	//! 没有锁时 map 并发读写会 data race，甚至直接 panic
	//* Gin 对照：手写 method + path 分流 → router.GET、router.POST、router.DELETE
	fmt.Println("🔑 一个 API 的闭环：路由 → 解析输入 → 校验 → 业务 → 状态码 + JSON")
}

// showTodoRequest 用 httptest 发一个假请求并打印「方法 路径 → 状态码 body」
// 不真实联网，纯内存模拟，方便在 demo 里演示整个 API 闭环
func showTodoRequest(handler http.Handler, method, path, body string) {
	req := httptest.NewRequest(method, "http://example.com"+path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	fmt.Printf("%s %-11s → %d %s\n", method, path, rec.Code, strings.TrimSpace(rec.Body.String()))
}
