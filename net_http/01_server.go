// 01_server.go —— HTTP Server：Handler 是请求处理入口
//
// Gin 最终也运行在 net/http 之上
// 先理解 Handler，之后才能看懂 gin.HandlerFunc 和 gin.Engine 的职责
//
// 核心概念：http.Handler 是一个「只有一个方法的接口」
//
//	type Handler interface {
//		ServeHTTP(w http.ResponseWriter, r *http.Request)
//	}
//
// 谁实现 ServeHTTP，谁就能处理 HTTP 请求
// 前端对照：Handler ≈ 一个路由回调，只是请求/响应不是全局对象，而是显式传进来的两个参数
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
//
// Go 接口是「隐式实现」：不用写 implements，方法签名对得上就算实现
// 前端对照：类似 TS 的 structural typing（结构兼容即视为实现）
type helloHandler struct{}

// ServeHTTP 的两个参数就是一次 HTTP 交互的全部输入输出：
//   - w http.ResponseWriter：把响应写回客户端（状态码 / 响应头 / body）
//   - r *http.Request：读取请求（method / URL / header / body）
func (helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只写 body、不写状态码 → 默认返回 200（状态码一旦写出就无法再改）
	fmt.Fprintf(w, "你好，%s", r.URL.Path)
}

func demoServer() {
	fmt.Println("========== 01 HTTP Server：Handler 接口 ==========")

	//* 前端对照：类似给 GET /hello 注册回调，只是请求和响应被显式传入
	var h http.Handler = helloHandler{}

	// httptest 在内存里「假装」发一次请求、收一次响应，不占端口、不阻塞 main
	// 前端对照：等于 mock 一个 fetch，省得真的起服务器再 curl
	//   NewRequest(method, url, body) 构造请求；NewRecorder() 是能接住响应的「假响应器」
	//   然后手动 h.ServeHTTP(rec, req)：把 rec 当 w、req 当 r 传进去，等价于「处理一次请求」
	req := httptest.NewRequest(http.MethodGet, "http://example.com/hello", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	fmt.Printf("① struct Handler → status=%d, body=%q\n", rec.Code, rec.Body.String())

	// HandlerFunc 是函数类型，它替普通函数实现了 ServeHTTP
	// 于是「一个签名匹配的普通函数」也能直接当 Handler，不用单独声明 struct
	// 前端对照：把匿名函数 (req, res) => {} 直接注册成路由回调
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "函数也能当 Handler")
	})
	rec = httptest.NewRecorder()
	funcHandler.ServeHTTP(rec, req)
	fmt.Printf("② HandlerFunc → status=%d, body=%q\n", rec.Code, rec.Body.String())

	//! 真正启动服务会持续阻塞（监听端口、等请求进来），必须写在 main 最后
	//! 上面用 httptest 就是为了不阻塞，让 5 个 demo 依次跑完
	//! err := http.ListenAndServe(":8080", h)
	//! if err != nil { log.Fatal(err) }

	fmt.Println("🔑 http.Handler 是 HTTP 处理契约，HandlerFunc 让普通函数也能满足它")
}
