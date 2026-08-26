// 06_todo_api.go —— 用 Gin 重写 net_http/05 的内存 TODO API
//
// 同一个需求两个框架各写一遍，Gin 的价值立刻可见
// （前端对照：jQuery 手动操作 DOM vs Vue 声明式渲染）
//
//	Gin 替你干掉的样板代码：
//
//	net_http 手写版                  Gin 版
//	手写 switch 按 method+path 分流   router.GET/POST/DELETE 声明即路由
//	json.NewDecoder(r.Body).Decode    c.ShouldBindJSON（还带校验）
//	writeJSON 手设 header + 状态码    c.JSON / ok / fail 一步到位
//	404 和 405 手动分支               路由表自动处理
package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)


// ginTodo binding:"required" 只能拦「没传」，拦不住「全是空格」——语义校验还得业务代码做
type ginTodo struct {
	ID   int    `json:"id"`
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}


// ginTodoStore 内存存储：RWMutex 保护并发读写
// 逻辑和 net_http/05 完全一致，此处不展开（RWMutex 忘了回 net_http 复习）
type ginTodoStore struct {
	mu     sync.RWMutex
	nextID int
	todos  map[int]ginTodo
}

func newGinTodoStore() *ginTodoStore {
	return &ginTodoStore{nextID: 1, todos: make(map[int]ginTodo)}
}

func (s *ginTodoStore) create(text string, done bool) ginTodo {
	// 写操作持写锁：nextID++ 和写 map 必须是一个不可分割的整体
	s.mu.Lock()
	defer s.mu.Unlock()
	t := ginTodo{ID: s.nextID, Text: text, Done: done}
	s.todos[t.ID] = t
	s.nextID++
	return t
}

func (s *ginTodoStore) list() []ginTodo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// map 遍历顺序不固定，按 ID 从小到大填入切片，保证输出稳定
	result := make([]ginTodo, 0, len(s.todos))
	for id := 1; id < s.nextID; id++ {
		if t, ok := s.todos[id]; ok {
			result = append(result, t)
		}
	}
	return result
}

func (s *ginTodoStore) get(id int) (ginTodo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	return t, ok
}

func (s *ginTodoStore) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	return true
}

func demoTodoAPI() {
	fmt.Println("\n========== 06 TODO API（Gin 重写版） ==========")

	router := gin.New()
	store := newGinTodoStore()

	v1 := router.Group("/api/v1")
	{
		// 列表 + 单个查询合在一条路由：带 ?id=1 查单个，不带查全部
		v1.GET("/todos", func(c *gin.Context) {
			var q struct {
				ID int `form:"id" binding:"omitempty,gt=0"`
			}
			// 对比手写版的 parseTodoID：绑定 + 校验一行搞定
			if err := c.ShouldBindQuery(&q); err != nil {
				fail(c, http.StatusBadRequest, "id 必须是正整数")
				return
			}
			if q.ID > 0 {
				t, found := store.get(q.ID)
				if !found {
					fail(c, http.StatusNotFound, "TODO 不存在")
					return
				}
				ok(c, t)
				return
			}
			ok(c, store.list())
		})

		v1.POST("/todos", func(c *gin.Context) {
			var input ginTodo
			if err := c.ShouldBindJSON(&input); err != nil {
				bindFail(c, err)
				return
			}
			// tag 校验管格式（非空串），业务校验管语义（不能纯空格）
			input.Text = strings.TrimSpace(input.Text)
			if input.Text == "" {
				fail(c, http.StatusBadRequest, "text 不能是空白")
				return
			}
			// 201 Created：新建资源的惯例状态码
			c.JSON(http.StatusCreated, apiResp{Code: 0, Msg: "ok", Data: store.create(input.Text, input.Done)})
		})

		v1.DELETE("/todos/:id", func(c *gin.Context) {
			var uri struct {
				ID int `uri:"id" binding:"required,gt=0"`
			}
			if err := c.ShouldBindUri(&uri); err != nil {
				fail(c, http.StatusBadRequest, "id 必须是正整数")
				return
			}
			if !store.delete(uri.ID) {
				fail(c, http.StatusNotFound, "TODO 不存在")
				return
			}
			// 204 No Content：删除成功的惯例，不带 body
			c.Status(http.StatusNoContent)
		})
	}

	showGinRequest(router, http.MethodPost, "/api/v1/todos", `{"text":"学 Gin 路由"}`)
	showGinRequest(router, http.MethodPost, "/api/v1/todos", `{"text":"学参数绑定"}`)
	showGinRequest(router, http.MethodGet, "/api/v1/todos", "")
	showGinRequest(router, http.MethodGet, "/api/v1/todos?id=1", "")
	showGinRequest(router, http.MethodDelete, "/api/v1/todos/1", "")
	showGinRequest(router, http.MethodGet, "/api/v1/todos?id=1", "")
	showGinRequest(router, http.MethodPost, "/api/v1/todos", `{}`)
	showGinRequest(router, http.MethodPost, "/api/v1/todos", `{"text":"   "}`)
	showGinRequest(router, http.MethodGet, "/api/v1/todos?id=abc", "")

	//! 对照着看 net_http/05_todo_api.go：路由、绑定、校验、响应全部声明式
	//! 手写版约 150 行，Gin 版不到 100 行，而且每一行都在说「业务」
	fmt.Println("🔑 框架的价值 = 把样板代码收走，让代码只剩业务")
}
