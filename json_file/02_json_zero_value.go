// 02_json_zero_value.go —— 零值困境：omitempty / 指针 / any / float64 精度
//
// Go 的「致命」设定：int 零值 0、string 零值 ""、bool 零值 false
// → 「用户没填」和「用户填了 0/false/""」在 struct 里长得一模一样！
//
// 前端对照：TS 有可选字段 age?: number（undefined = 没填）
// Go 没有可选字段，只能用「指针」或「omitempty」来模拟
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// profile 零值困境的三件套
type profile struct {
	Name  string         `json:"name"`
	Age   int            `json:"age,omitempty"` // 零值(0)时整个字段从 JSON 里消失
	Score *int           `json:"score"`         // 指针：nil → null；指向 0 → 0
	Ext   map[string]any `json:"ext,omitempty"` // 未知结构兜底；nil map 也不输出
}

func demoJSONZeroValue() {
	fmt.Println("\n========== 02 JSON 零值困境：omitempty / 指针 / any ==========")

	// ---------- ① omitempty：零值字段直接消失 ----------
	p1 := profile{Name: "小明"}
	out, _ := json.Marshal(p1)
	fmt.Printf("① age=0 被删、score=nil 输出 null → %s\n", out)

	//? 为什么需要 omitempty？返回列表时不把一堆 0/false/"" 的空字段全吐给前端
	//? 前端对照：TS 的 age?: number —— undefined 时 JSON.stringify 也不带它

	// ---------- ② 零值困境：0 分是真分数！ ----------
	// 「没参加考试」(null) 和「考了 0 分」(0) 必须能区分 → 用指针
	zero := 0
	p2 := profile{Name: "小红", Score: &zero, Ext: map[string]any{"vip": true}}
	out2, _ := json.Marshal(p2)
	fmt.Printf("② 指针+any 兜底：%s\n", out2)

	//? 指针字段加 omitempty 也行：nil 时连字段都不输出（「有值才出现」）
	//? Ext 这类 map/slice 的零值是 nil，不加 omitempty 会输出 "ext":null

	// ---------- ③ 请求进来：nil = 没传，非 nil = 传了（哪怕传 0） ----------
	// PATCH 局部更新的地基：只更新用户传了的字段
	var patched profile
	_ = json.Unmarshal([]byte(`{"name":"小刚","score":0}`), &patched)
	fmt.Printf("③ 传了 score:0 → Score == nil ? %v（非 nil，用户显式传了 0）\n", patched.Score == nil)

	var noScore profile
	_ = json.Unmarshal([]byte(`{"name":"小刚"}`), &noScore)
	fmt.Printf("③ 没传 score  → Score == nil ? %v（nil，说明用户没动这个字段）\n", noScore.Score == nil)

	// ---------- ④ map[string]any：结构未知的 JSON 兜底 ----------
	var payload map[string]any
	_ = json.Unmarshal([]byte(`{"id":1,"tags":["go","后端"],"meta":{"vip":true}}`), &payload)
	fmt.Printf("④ any 兜底：%+v\n", payload)

	//! ⚠️ any 解码时所有数字一律变 float64！
	//! 大整数会丢精度：9007199254740993 会变成 9007199254740992
	//! 这和 JS 的 Number 精度丢失是同一个坑（都是 IEEE 754 双精度）
	v := payload["id"].(float64) // 断言成 float64 才能当数字用
	fmt.Printf("④ id 的真实类型是 %T，值 %v\n", v, v)

	// v2 := payload["id"].(int) // 解开这行试试：断言失败直接 panic（int ≠ float64）

	//? 怎么破大数丢精度？
	//? 方案1：结构已知就定义 struct，大 ID 字段直接用 string（后端返回雪花ID 的惯例）
	//? 方案2：Decoder.UseNumber() → 数字变成 json.Number（本质是字符串），要算时再转
	dec := json.NewDecoder(strings.NewReader(`{"big":9007199254740993}`))
	dec.UseNumber()
	var m map[string]any
	_ = dec.Decode(&m)
	fmt.Printf("④ UseNumber：big = %v（json.Number，没丢精度）\n", m["big"])

	fmt.Println("🔑 表达「没填」用指针；接「未知结构」用 map+any；大数字警惕 float64")
}
