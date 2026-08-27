// 04_router_middleware.go —— 手写路由和中间件
//
// 路由解决“哪个 Handler 处理请求”
// 中间件解决“每个请求前后都要做什么”
// Gin 的核心价值就是把这两层封装得更方便
//
// 中间件的本质是「洋葱模型」：一层包一层
// 请求从最外层进 → 依次穿过各层 → 到达业务 Handler → 再反向依次退出
// 前端对照：middleware ≈ Axios 拦截器 / Express 的 app.use
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// middleware 接收下游 Handler，返回包好后的新 Handler
// 前端对照：这是一个「装饰器」——拿进来一个 handler，包一层再还回去
type middleware func(http.Handler) http.Handler

// requestIDKey 是空结构体，专门当 Context 的 key
// 用自定义类型而不是 string，是为了避免和其他包往 Context 里塞的 key 撞车
// 前端对照：类似用 Symbol 当唯一 key，而不是裸字符串
type requestIDKey struct{}

func demoRouterMiddleware() {
	fmt.Println("\n========== 04 Router + Middleware：路由分发和洋葱链 ==========")

	// ServeMux 是标准库自带的「路由器」，按 method + path 把请求分发给对应 Handler
	// Go 1.22 起支持 "GET /hello" 这种「方法 + 路径」写法；老版本只能写 "/hello"
	// 前端对照：mux ≈ Vue Router / Express 的 router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		// 从请求 Context 里取出中间件塞进去的 requestID（下面会讲它从哪来）
		requestID, _ := r.Context().Value(requestIDKey{}).(string)
		fmt.Fprintf(w, "hello, requestID=%s", requestID)
	})

	// applyMiddleware 从后往前包：先 logger 包住 mux，再用 requestID 包住整体
	// 最终结构 = requestID( logger( mux ) )，所以 requestID 最外层、最先进入
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

// middlewares 类型就是 []middleware, 切片
func applyMiddleware(final http.Handler, middlewares ...middleware) http.Handler {
	handler := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// requestIDMiddleware 将请求元信息放进派生 Context，再交给下游
// 前端对照：类似给每个请求挂一个 requestId，之后所有 handler 都能读到
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//* Context 沿调用链传递请求元信息，不要存进全局变量或 struct（会并发串数据）
		//* key 用自定义类型（空 struct），避免和其他包的 string key 冲突
		// WithValue 不改原 ctx，而是返回一个「派生」的新 ctx，所以必须用返回值
		ctx := context.WithValue(r.Context(), requestIDKey{}, "req-demo-001")
		// 把新 ctx 塞回请求再传给下游，下游用 r.Context() 就能读到
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 洋葱模型：next.ServeHTTP 之前是「进」，之后是「出」
		// 进入时记时间 → 交给下游处理 → 退出时打印耗时
		started := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("   [logger] %s %s, cost=%s\n", r.Method, r.URL.Path, time.Since(started).Round(time.Microsecond))
	})
}
