// 03_response_json.go —— Response：状态码、响应头、JSON 和 struct tag
//
// 本节只看状态码和 Content-Type
//
// 状态码说明结果
// Content-Type 告诉客户端如何解析 body
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// article 的 json tag 决定 JSON 字段名
// 没有 tag 时默认使用 Go 字段名，例如 Title 和 Published
type article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Published bool   `json:"published"`
}

func demoResponseJSON() {
	fmt.Println("\n========== 03 Response：JSON 请求和响应 ==========")

	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusCreated, article{ID: 1, Title: "学 net/http", Published: true})
	fmt.Printf("① status=%d, Content-Type=%q, body=%s\n",
		rec.Code,
		rec.Header().Get("Content-Type"),
		strings.TrimSpace(rec.Body.String()),
	)

	body := strings.NewReader(`{"title":"写 TODO API","published":false}`)
	req := httptest.NewRequest(http.MethodPost, "/articles", body)
	var input article
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		fmt.Println("② JSON 解析失败：", err)
		return
	}
	fmt.Printf("② 解码后：%+v\n", input)

	//! Server 会自动关闭请求 body，不需要在 Handler 里 defer r.Body.Close()
	//! JSON 解码失败属于客户端输入错误，应返回 400，不能 panic

	//* Gin 对照：writeJSON(...) → c.JSON(...)
	//* Gin 对照：Decoder.Decode(&input) → c.ShouldBindJSON(&input)
	fmt.Println("🔑 响应先设 header，再写 status 和 body")
}

// writeJSON 统一写 JSON 响应
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		// ResponseWriter 已写出 status 后不能可靠撤回
		return
	}
}
