// 05_assertion.go —— any + 类型断言 + type switch：从接口里掏类型（重点）
//
// 前端对照一句话（先记这个）：
//
//	TS：unknown 类型想用，得先"收窄"（typeof / instanceof / as）
//	Go：any 类型想用，得先"断言"（v, ok := x.(T) / switch v := x.(type)）
//
// 场景：JSON 反序列化到 any、Gin 的 gin.H、第三方库返回 any ——
//   拿到的都是"不知道类型的东西"，想用必须掏出来。日常开发天天见。
package main

import "fmt"

func demoAssertion() {
	fmt.Println("========== 05 any + 类型断言 + type switch ==========")

	// ---------- ① 两种断言写法 ----------
	var x any = "hello"

	// 写法 A：单返回值 —— 类型不对直接 panic 崩溃！不推荐
	// s := x.(string)

	// 写法 B：comma-ok —— 安全，推荐（呼应 map / channel 的 comma-ok）
	if s, ok := x.(string); ok {
		fmt.Println("① 断言成功：x 是 string：", s)
	}
	if n, ok := x.(int); ok {
		fmt.Println("① 不会走到：", n)
	} else {
		fmt.Println("① 断言失败走 else：x 不是 int")
	}

	// ---------- ② type switch：多个类型一次分拣 ----------
	// 前端对照：TS 的 typeof + 收窄（narrowing）
	checkType(42)
	checkType("hi")
	checkType(3.14)
	checkType([]string{"a"})

	// ---------- ③ 实战：从 map[string]any 里取值 ----------
	// json.Unmarshal 到 any 得到的就是 map[string]any（模拟接口返回的未知数据）
	payload := map[string]any{
		"name":  "小明",
		"age":   float64(18), // 注意：JSON 数字反序列化后默认是 float64！
		"tags":  []any{"前端", "Go"},
		"isVIP": true,
		"nested": map[string]any{"city": "上海"},
	}

	// 取值写法：断言 + comma-ok（拿不到就用零值兜底）
	name, _ := payload["name"].(string)
	age, _ := payload["age"].(float64)
	fmt.Println("③ 从 map[string]any 取值：", name, "，", int(age), "岁")

	// 嵌套取值：先断言成 map，再掏下一层（JSON 多层结构就这么处理）
	if nested, ok := payload["nested"].(map[string]any); ok {
		city, _ := nested["city"].(string)
		fmt.Println("③ 嵌套取值：城市 =", city)
	}

	// ---------- ④ Gin 的 gin.H 真相 ----------
	// c.JSON(200, gin.H{"msg": "ok"}) —— gin.H 就是 map[string]any 的别名
	// 所以它才能"随便塞任何值"，最后序列化成 JSON 发给前端
	fmt.Println("④ gin.H = map[string]any 的别名，所以能随便塞值")

	//! 最经典的坑：JSON 里的数字反序列化后是 float64，不是 int！
	//!   v, ok := payload["age"].(int) → ok 是 false，拿不到
	//!   想拿 int：先断 float64 再转，或者用 json.Number 处理
	fmt.Println("\n🔑 any 装万物，断言掏出来；comma-ok 永远比单返回值安全")
}

// checkType 用 type switch 分拣类型
func checkType(v any) {
	switch t := v.(type) { // 语法：switch 变量.(type)，每个 case 写一个类型
	case int:
		fmt.Println("② int：", t)
	case string:
		fmt.Println("② string：", t)
	case float64:
		fmt.Println("② float64：", t)
	default:
		fmt.Println("② 其他类型：", t)
	}
}
