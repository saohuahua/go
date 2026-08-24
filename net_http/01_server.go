// 01_server.go —— HTTP Server：Handler 是请求处理入口
//
// Gin 最终也运行在 net/http 之上
// 先理解 Handler，之后才能看懂 gin.HandlerFunc 和 gin.Engine 的职责
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// helloHandler 实现 ServeHTTP，因此满足 http.Handler
//
// http.Handler 只有一个方法
//
// ServeHTTP(http.ResponseWriter, *http.Request)
type helloHandler struct{}

func (helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "你好，%s", r.URL.Path)
}

func demoServer() {
	fmt.Println("========== 01 HTTP Server：Handler 接口 ==========")

	//* 前端对照：类似给 GET /hello 注册回调，只是请求和响应被显式传入
	var h http.Handler = helloHandler{}

	// httptest 在内存中模拟请求，不占端口且不阻塞 main
	req := httptest.NewRequest(http.MethodGet, "http://example.com/hello", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	fmt.Printf("① struct Handler → status=%d, body=%q\n", rec.Code, rec.Body.String())

	// HandlerFunc 是函数类型，它替普通函数实现了 ServeHTTP
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "函数也能当 Handler")
	})
	rec = httptest.NewRecorder()
	funcHandler.ServeHTTP(rec, req)
	fmt.Printf("② HandlerFunc → status=%d, body=%q\n", rec.Code, rec.Body.String())

	//! 真正启动服务会持续阻塞，通常写在 main 的最后
	//! err := http.ListenAndServe(":8080", h)
	//! if err != nil { log.Fatal(err) }

	fmt.Println("🔑 http.Handler 是 HTTP 处理契约，HandlerFunc 让普通函数也能满足它")
}
