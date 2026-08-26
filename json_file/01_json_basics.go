// 01_json_basics.go —— JSON 编解码基础：Marshal / Unmarshal / struct tag
//
// Go 和 JSON 之间的桥梁就是 struct tag —— 一份 tag 同时管「序列化」和「反序列化」
//
// 前端对照
// JSON.stringify ↔ json.Marshal
// JSON.parse     ↔ json.Unmarshal
// TS 的 DTO + 字段重映射 ↔ Go 的 struct + json tag
//
// 呼应 net_http/03：writeJSON 和 Decoder.Decode 的底层就是本章这套
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// user 演示用结构体：json tag 决定 JSON 里的字段名
// tag 写在反引号里（不能换成普通双引号字符串），格式：`json:"小写名"`
type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"` // 对外暴露的字段
	Password string `json:"-"`     // 敏感字段：- 表示序列化时直接忽略
	Age      int    // 没有 tag → JSON 字段名就用 Go 字段名 "Age"
}

func demoJSONBasics() {
	fmt.Println("\n========== 01 JSON 基础：Marshal / Unmarshal / tag ==========")

	// ---------- ① 序列化：struct → JSON ----------
	u := user{ID: 1, Name: "小明", Email: "xm@test.com", Password: "123456", Age: 18}

	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println("① Marshal 失败：", err)
		return
	}
	fmt.Printf("① Marshal：Password 被 tag \"-\" 干掉了 → %s\n", data)

	// MarshalIndent：带缩进的好看版本（写配置文件、调试时用）
	pretty, _ := json.MarshalIndent(u, "", "  ") // 这种简单 struct 不可能 Marshal 失败
	fmt.Printf("① MarshalIndent（好看版）：\n%s\n", pretty)

	// ---------- ② 反序列化：JSON → struct ----------
	// Unmarshal 第二个参数必须传指针：Go 值传递，传值只改到一份拷贝
	// 前端对照：JSON.parse 返回一个新对象；Go 是「把数据填进你给的结构体」
	var u2 user
	if err := json.Unmarshal([]byte(`{"id":7,"name":"小红","Age":20}`), &u2); err != nil {
		fmt.Println("② Unmarshal 失败：", err)
		return
	}

	//? 字段匹配是【大小写不敏感】的：JSON 的 "Age" 能填进没 tag 的 Age 字段
	fmt.Printf("② Unmarshal：ID=%d, Name=%s, Age=%d\n", u2.ID, u2.Name, u2.Age)

	// ---------- ③ 缺字段=零值，多字段=忽略 ----------
	//? JSON 里没 email → u3.Email 保持零值 ""，不报错
	//? JSON 里多出的 "extra" → 默认直接丢弃（Gin 的参数绑定也是这个行为）
	var u3 user
	_ = json.Unmarshal([]byte(`{"id":1,"name":"小刚","extra":"多余字段"}`), &u3)
	fmt.Printf("③ 缺字段=零值，多字段=忽略 → Email=%q, Age=%d\n", u3.Email, u3.Age)

	// ---------- ④ 类型不匹配 → 返回错误 ----------
	// JSON 的字符串 "abc" 塞不进 int 字段
	var u4 user
	err = json.Unmarshal([]byte(`{"id":"abc"}`), &u4)
	fmt.Printf("④ 类型不匹配 → err = %v\n", err)

	// ---------- ⑤ HTML 转义坑：Marshal 默认转义 < > & ----------
	// 防 XSS 的默认行为，但 URL 里的 & 会变成 \u0026，返回给前端时要注意
	url := "https://a.com?x=1&y=2"
	escaped, _ := json.Marshal(map[string]string{"url": url})
	fmt.Printf("⑤ 默认转义：%s\n", escaped)

	//! 想关掉转义要用 Encoder + SetEscapeHTML(false)，Marshal 本身没有开关
	//! Gin 的 c.JSON 默认就是「不转义」的（内部走的 Encoder）
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(map[string]string{"url": url})
	fmt.Printf("⑤ 关闭转义：%q（注意 Encode 自带结尾换行）\n", buf.String())

	//* Gin 对照：Marshal / Unmarshal 就是 c.JSON / ShouldBindJSON 的底层
	fmt.Println("🔑 一份 json tag 同时管进和出：JSON 长什么样，tag 说了算")
}
