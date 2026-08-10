// 02_reference.go —— 本质：map 是引用类型（传参 / 嵌套 / 复制陷阱）
//
// 先回想 slices/05 的结论：切片传参是值传递，但共享底层数组。
// map 更简单粗暴：它【本身就是引用类型】，传参传的就是那个 map 本体。
//
//   - 好消息（前端同学应该秒懂）：
//     JS 对象就是引用传递 → function f(o){ o.x = 1 } 会改到外面
//     Go 的 map 行为【完全一样】，这块对你反而比切片简单！
package main

import "fmt"

// demoReference 演示 map 的引用语义、嵌套 map、复制陷阱
func demoReference() {
	fmt.Println("========== 02 本质：引用类型 ==========")

	// ---------- 1. 传参：函数里改，外面跟着变 ----------
	ages := map[string]int{"小明": 10}
	addAge(ages)                    // 函数内加了一个 key
	fmt.Println("传参后 ages =", ages) // map[小明:10 小红:5] ← 外面也变了

	// 对比切片（slices/05）：append 要"接住返回值"才生效；map 什么都不用，天然生效

	// ---------- 2. 嵌套 map：map 套 map ----------
	// 场景：用户 → 用户的一堆属性
	users := map[string]map[string]int{
		"张三": {"年龄": 18, "身高": 175},
	}
	fmt.Println("\n嵌套 map users =", users)

	//! 嵌套 map 的经典 panic：只 make 外层，不初始化内层
	//! n := make(map[string]map[string]int)
	//! n["a"]["x"] = 1  // ❌ panic: assignment to entry in nil map
	//! 原因：n["a"] 本身是 nil map。正确写法：先 n["a"] = make(map[string]int)

	fmt.Println("🔑 嵌套 map 要一层一层 make，内层是 nil 就写不进去")

	// ---------- 3. 复制陷阱：m2 := m 不是复制！ ----------
	original := map[string]int{"a": 1, "d": 4}
	m2 := original // 只是让 m2 和 original 指向同一个 map
	m2["b"] = 2
	fmt.Println("\n改 m2 后 original =", original) // map[a:1 b:2] ← 也被改了！

	// 想要真正独立的副本 → 手动循环复制（Go 没有内置的 copy map）
	realCopy := make(map[string]int, len(original))
	for k, v := range original {
		realCopy[k] = v
		fmt.Println(realCopy[k]) // 顺序是乱的，随机的，不唯一
	}
	realCopy["c"] = 3
	fmt.Println("改 realCopy 后 original =", original) // map[a:1 b:2] ← 不受影响

	// 前端对照：JS 有 {...obj} 浅拷贝 / structuredClone，Go 只能自己循环
}

// addAge 演示：map 是引用类型，函数内直接改原 map，无需返回
func addAge(m map[string]int) {
	m["小红"] = 5
}
