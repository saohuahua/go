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
// 前端对照：tag ≈ Go 结构体字段映射成 TS DTO 字段名，序列化/反序列化共用同一份 tag
type article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Published bool   `json:"published"`
}

func demoResponseJSON() {
	fmt.Println("\n========== 03 Response：JSON 请求和响应 ==========")

	//? 为什么用 httptest？本 demo 没起真实服务器——NewRecorder 造一个假的 ResponseWriter，
	//? 不联网，写入全存内存，下面 rec.Code / rec.Header / rec.Body 再读出来看
	//* 前端对照：rec ≈ { status:0, headers:{}, body:'' } 的空对象，写完打开看
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusCreated, article{ID: 1, Title: "学 net/http", Published: true})
	// ↑ 把 status=201 + Content-Type + JSON body 全写进了 rec 的内存，下面三行读出来：
	fmt.Printf("① status=%d, Content-Type=%q, body=%s\n",
		rec.Code,
		rec.Header().Get("Content-Type"),
		strings.TrimSpace(rec.Body.String()),
	)

	// 从请求 body 解 JSON：NewDecoder(r.Body).Decode(&input) 把 body 填进 input
	// 前端对照：≈ JSON.parse(await req.text()) 再手动赋值；Go 用 tag 自动映射字段
	// strings.NewReader：字符串不能直接当 body（body 须是 io.Reader，能被流式读），包一层即可
	body := strings.NewReader(`{"title":"写 TODO API","published":false}`)
	req := httptest.NewRequest(http.MethodPost, "/articles", body) // 手搓一个「客户端发来的 POST 请求」，req.Body 即客户端 JSON
	var input article                                              // 全零结构体，准备接收
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		fmt.Println("② JSON 解析失败：", err)
		return
	}
	// Decode(&input) 必须传指针：Go 值传递，不传 & 只改到一份拷贝，外头拿不到结果
	// 请求没带 id 字段 → input.ID 保持零值 0；缺失字段不报错
	fmt.Printf("② 解码后：%+v\n", input) // %+v 带字段名打印，一眼看清每个字段

	//! Server 会自动关闭请求 body，不需要在 Handler 里 defer r.Body.Close()
	//! JSON 解码失败属于客户端输入错误，应返回 400，不能 panic

	//* Gin 对照：writeJSON(...) → c.JSON(...)
	//* Gin 对照：Decoder.Decode(&input) → c.ShouldBindJSON(&input)
	fmt.Println("🔑 响应先设 header，再写 status 和 body")
}

// writeJSON 统一写 JSON 响应
// 顺序不能反：先设 header → 再写状态码 → 最后写 body
// 因为 WriteHeader 一调用状态码就立刻发出；body 一写若还没状态码就默认 200
// 前端对照：等于 Express 里先 res.set(...) 再 res.status(...).json(...)
func writeJSON(w http.ResponseWriter, status int, value any) {
	// any 是 interface{} 的别名，表示「任意类型」，让本函数能序列化任何结构体
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// NewEncoder(w).Encode(v)：边序列化边写进 w，等价于 json.Marshal 后再 w.Write
	if err := json.NewEncoder(w).Encode(value); err != nil {
		// ResponseWriter 已写出 status 后不能可靠撤回，只能忽略或打日志
		return
	}
}
