// 04_router_middleware.go —— 手写路由和中间件
//
// 路由解决“哪个 Handler 处理请求”
// 中间件解决“每个请求前后都要做什么”
// Gin 的核心价值就是把这两层封装得更方便
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// middleware 接收下游 Handler，返回包好后的新 Handler
type middleware func(http.Handler) http.Handler

type requestIDKey struct{}

func demoRouterMiddleware() {
	fmt.Println("\n========== 04 Router + Middleware：路由分发和洋葱链 ==========")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := r.Context().Value(requestIDKey{}).(string)
		fmt.Fprintf(w, "hello, requestID=%s", requestID)
	})

	// applyMiddleware 从后往前包
	// logger(requestID(mux)) 因此 logger 最先进入、最后退出
	handler := applyMiddleware(mux, requestIDMiddleware, loggerMiddleware)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/hello", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	fmt.Printf("① 命中 GET /hello → status=%d, body=%q\n", rec.Code, rec.Body.String())

	wrongMethod := httptest.NewRequest(http.MethodPost, "http://example.com/hello", nil)
	wrongMethodRec := httptest.NewRecorder()
	handler.ServeHTTP(wrongMethodRec, wrongMethod)
	fmt.Println("② POST /hello → status =", wrongMethodRec.Code)

	//* Gin 对照：mux.HandleFunc(...) → router.GET(...)
	//* Gin 对照：middleware → router.Use(...)
	fmt.Println("🔑 路由负责分发，中间件负责通用横切逻辑")
}

func applyMiddleware(final http.Handler, middlewares ...middleware) http.Handler {
	handler := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// requestIDMiddleware 将请求元信息放进派生 Context，再交给下游
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//* Context 沿调用链传递请求元信息，不要存进全局变量或 struct
		//* key 用自定义类型，避免和其他包的 string key 冲突
		ctx := context.WithValue(r.Context(), requestIDKey{}, "req-demo-001")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("   [logger] %s %s, cost=%s\n", r.Method, r.URL.Path, time.Since(started).Round(time.Microsecond))
	})
}
