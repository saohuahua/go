// 09_http.go —— 最后一关：把之前的 Task、接口、error、Context 串成 HTTP API
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type requestIDKey struct{}

type taskAPI struct {
	repo TaskRepository
}

// httpMemoryRepo 是给 HTTP 练习用的并发安全实现。
// TODO 09.1：补全两个方法。读用 RLock，写用 Lock，并用 defer 解锁。
type httpMemoryRepo struct {
	mu    sync.RWMutex
	tasks map[int]Task
}

func newHTTPMemoryRepo(tasks []Task) *httpMemoryRepo {
	items := make(map[int]Task, len(tasks))
	for _, task := range tasks {
		items[task.ID] = task
	}
	return &httpMemoryRepo{tasks: items}
}

func (repo *httpMemoryRepo) FindByID(id int) (Task, error) {
	return Task{}, fmt.Errorf("TODO: HTTP repo find")
}

func (repo *httpMemoryRepo) Save(task Task) error {
	return fmt.Errorf("TODO: HTTP repo save")
}

func demoHTTP() {
	fmt.Println("\n========== 09 HTTP 综合：路由 / JSON / 中间件 ==========")
	api := taskAPI{repo: newHTTPMemoryRepo(sampleTasks())}
	handler := applyHTTPMiddleware(http.HandlerFunc(api.serveHTTP), requestID)

	// 先完成下面 TODO，再取消这段验证代码的注释。
	// req := httptest.NewRequest(http.MethodGet, "/tasks/2", nil)
	// rec := httptest.NewRecorder()
	// handler.ServeHTTP(rec, req)
	// check("GET /tasks/2", rec.Code, http.StatusOK)
	// check("响应含 request ID", rec.Header().Get("X-Request-ID"), "practice-001")
	//
	// wrongMethod := httptest.NewRequest(http.MethodPost, "/tasks/2", nil)
	// wrongMethodRec := httptest.NewRecorder()
	// handler.ServeHTTP(wrongMethodRec, wrongMethod)
	// check("POST /tasks/2", wrongMethodRec.Code, http.StatusMethodNotAllowed)
	_ = handler
}

// TODO 09.2：实现最小路由。
// - GET /tasks/{id}：调用 getTask
// - PATCH /tasks/{id}/done：调用 completeTaskHTTP
// - path 不匹配返回 404；path 匹配但 method 错误返回 405
// 提示：strings.TrimPrefix(r.URL.Path, "/tasks/")，再 Split。
func (api taskAPI) serveHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "TODO: route"})
}

// TODO 09.3：从 path 的 id 部分解析整数，查 repo 后写 JSON。
// - id 不是整数：400
// - errors.Is(err, ErrTaskNotFound)：404
// - 成功：200 + Task JSON
// Context 要作为 completeTaskWithContext 的第一个参数一路向下传。
func (api taskAPI) getTask(w http.ResponseWriter, r *http.Request, idText string) {
	writeJSON(w, http.StatusNotImplemented, nil)
}

// TODO 09.4：调用 completeTaskWithContext。
// - id 不是整数：400
// - 不存在：404
// - 成功：200 + 更新后的 Task JSON
func (api taskAPI) completeTaskHTTP(w http.ResponseWriter, r *http.Request, idText string) {
	writeJSON(w, http.StatusNotImplemented, nil)
}

// TODO 09.5：第 04 / 05 章 completeTask 的 Context 版。
// 要求：调用 repo 前先 select 检查 ctx.Done()；FindByID / Save 错误都用 %w 包裹。
func completeTaskWithContext(ctx context.Context, repo TaskRepository, id int) (Task, error) {
	return Task{}, fmt.Errorf("TODO: complete with context")
}

// TODO 09.6：中间件生成 request ID，放入派生 Context，并在响应头写 X-Request-ID。
// //* Gin 对照：c.Set("requestID", id) + c.Header("X-Request-ID", id)。
func requestID(next http.Handler) http.Handler {
	return next
}

type httpMiddleware func(http.Handler) http.Handler

func applyHTTPMiddleware(final http.Handler, middlewares ...httpMiddleware) http.Handler {
	handler := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

// 保留这些引用，方便你完成后在本文件直接扩展 POST JSON 绑定练习。
var _ = strconv.Atoi
var _ = strings.Split
