// 01_basics.go —— 上手：初始化 + nil map 大坑 + 增删改查 + comma-ok
//
// 先打个比喻，方便你从 JS 思维迁移过来：
//
//	map = 一本「字典」：给你一个 key，立刻翻到对应的 value
//	JS 里你用的对象 {} / new Map() 就是 map，Go 里叫 map
//
// 这一节学完你就能"用" map 了，核心就一句话：
//
//	用之前先 make，读之前先确认 key 在不在。
package main

import "fmt"

// demoBasics 演示 map 初始化、nil map 大坑、增删改查、comma-ok
func demoBasics() {
	fmt.Println("========== 01 上手：初始化 + 增删改查 ==========")


	// ---------- 1. 三种初始化方式 ----------
	// 方式一：字面量（带数据初始化，最直观）
	m1 := map[string]int{"a": 1, "b": 2}
	fmt.Println("字面量      m1 =", m1)


	// 方式二：make（先建空 map 再填，最常用）
	m2 := make(map[string]int)
	m2["c"] = 3
	fmt.Println("make        m2 =", m2)


	// 方式三：make 时给容量提示（性能优化：预估数据量，少扩容几次）
	m3 := make(map[string]int, 100)
	m3["d"] = 4 // 容量只是"预估"，想放多少还是能放多少
	fmt.Println("make+容量   m3 =", m3, " len =", len(m3))


	// ---------- 2. nil map 大坑（❌ 必踩） ----------
	var m4 map[string]int // 只声明、不初始化 → m4 是 nil
	fmt.Println("\nm4 == nil ?", m4 == nil) // true

	fmt.Println("读 nil map 不会崩，返回零值：m4[\"a\"] =", m4["a"]) // 0

	//! 写 nil map 直接 panic 崩溃！先把下面这行解开注释试一次：
	//! m4["a"] = 1  // ❌ panic: assignment to entry in nil map

	fmt.Println("🔑 记忆法：nil map 只许读、不许写；想写必须先 make 或字面量")
	fmt.Println("   前端对照：JS 里 const obj = {} 天生就能 obj.x = 1，Go 不行！")


	// ---------- 3. 增 / 改 / 删 / 长度 ----------
	scores := make(map[string]int)
	scores["语文"] = 90    // 没有 → 新增
	scores["语文"] = 95    // 已有 → 覆盖（改）
	scores["英语"] = 60    // 
	delete(scores, "语文") // 删除（删不存在的 key 也不报错）
	// delete(scores, "数学") // 删除（删不存在的 key 也不报错）
	fmt.Println("\n增改删之后 scores =", scores, " len =", len(scores))


	// ---------- 4. 读：comma-ok 判断 key 是否存在 ----------
	// JS：obj[key] 读不到 → undefined，一眼看出"没有"
	// Go：读不到 → 返回【零值】（0 / ""），看起来像"有数据"！
	v, ok := scores["数学"] // 数学不存在
	fmt.Println("\nv =", v, " ok =", ok) // 0 false ← ok 才是"是否存在"的答案

	//? 思考：如果只写 v := scores["数学"]，v 是 0 —— 你怎么知道是"考了 0 分"还是"没这门课"？
	//? 答：靠 ok。所以判断存在，永远用 v, ok := m[k] 写法
}
