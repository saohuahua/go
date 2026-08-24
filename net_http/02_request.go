// 02_request.go —— Request：读取 method、URL、query 和 header
//
// HTTP 请求信息都在 *http.Request 中
// Gin 的 c.Query 和 c.GetHeader 本质上是对这些字段的封装
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func demoRequest() {
	fmt.Println("\n========== 02 Request：读取请求 ==========")

	// URL 的 ? 后是 query string，同名 key 可以有多个值
	// 所以 Query() 返回 url.Values = map[string][]string（值是切片，因为 key 可重复）
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/articles?page=2&tag=go&tag=http",
		nil, // 第三个参数是 body；GET 没有 body，传 nil
	)
	req.Header.Set("X-Request-ID", "req-123")

	fmt.Println("① method =", req.Method)
	fmt.Println("② path =", req.URL.Path)
	fmt.Println("③ page =", req.URL.Query().Get("page"))
	fmt.Println("④ tags =", req.URL.Query()["tag"])
	fmt.Println("⑤ X-Request-ID =", req.Header.Get("X-Request-ID"))

	// Get 取不到时返回空字符串，但没法区分「没传」和「传了空值」
	// 用 map 的 comma-ok：query["keyword"] 直接当下标取，返回 (值, 是否存在)
	query := req.URL.Query()
	_, hasKeyword := query["keyword"]
	fmt.Println("⑥ keyword 是否传入 =", hasKeyword)

	//! 服务端必须按 method 分流，不能只依赖前端约定
	switch req.Method {
	case http.MethodGet:
		fmt.Println("⑦ GET：读取资源")
	case http.MethodPost:
		fmt.Println("⑦ POST：创建资源")
	default:
		fmt.Println("⑦ 当前 demo 未处理的方法")
	}

	//* Gin 对照：req.URL.Query().Get("page") → c.Query("page")
	//* Gin 对照：req.Header.Get("X-Request-ID") → c.GetHeader("X-Request-ID")
	fmt.Println("🔑 Request 是输入，先确认 method，再读取 URL、query、header 或 body")
}
