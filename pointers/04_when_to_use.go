// 04_when_to_use.go —— 总结：什么时候该用指针？+ new 初始化
//
// 这节是八股复习的浓缩版：指针的"使用场景"一次说清，
// 面试答"什么时候用指针"就背这一节。
package main

import "fmt"

// demoSummary 演示指针的使用场景、new、以及面试速背
func demoSummary() {
	fmt.Println("========== 04 什么时候用指针 ==========")


	// ---------- 1. 场景一：表示"可能没有"的值（可空） ----------
	// int 本身不能是 null，但 *int 可以（nil）
	var ptr *int // nil = "没有"
	if ptr == nil {
		fmt.Println("1. nil 指针 = 明确的「没有」：判断 ptr == nil 就知道")
	}
	//* 前端对照：JS 里 string/number 也可以是 null/undefined；
	//* Go 里只有指针能"空" → 需要可空时就用 *T（服务端很常见）


	// ---------- 2. 场景二：new —— 创建指针的简写 ----------
	// new(int) = 创建一个 int 并返回它的地址，值自动是零值 0
	n := new(int)
	fmt.Println("2. new(int)：*n =", *n, "（自动初始化为零值 0）")
	// 记忆：make 只用于 slice/map/channel；new 用于其他所有类型
	// 但日常更常用 &User{...} 这种字面量写法，new 偶尔见


	// ---------- 3. 场景三：直接声明结构体并拿指针 ----------
	u := &User{Name: "张三"} // 一步到位：创建结构体 + 取地址
	fmt.Println("3. &User{...} 直接得到指针：u.Name =", u.Name)


	// ---------- 4. 面试速背：指针的四个使用场景 ----------
	fmt.Println("\n4. 面试答「什么时候用指针」背这四条：")
	fmt.Println("   ① 修改调用方的变量 / 结构体字段（传指针、指针接收者）")
	fmt.Println("   ② 结构体很大时避免整块拷贝（省内存、省时间）")
	fmt.Println("   ③ 表示可空：nil 判断（模拟 JS 的 null）")
	fmt.Println("   ④ 方法接收者要改字段 → 指针接收者")
	fmt.Println("   —— 只是读值、小变量，直接用值就行，别滥用指针")
}
