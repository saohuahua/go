// 04_closure.go —— 函数是一等公民 + 闭包（Gin 中间件的灵魂）
//
// 前端对照一句话（先记这个）：
//
//	JS：const f = () => {} —— 函数就是值，能存、能传、能返回
//	Go：完全一样！func(参数) 返回类型 就是函数的"类型签名"
//
// 为什么重要：Gin 中间件 = 外层函数返回内层闭包，外层配置被闭包"记住"。
//   认证、CORS、日志、限流全是这个模式 —— 不学闭包，中间件只能死背模板。
//
// 本节一条主线：函数能存 → 函数能传 → 函数能返回 → 闭包记住外层 → 中间件模式
package main

import "fmt"

func demoClosure() {
	fmt.Println("========== 04 闭包：函数就是值 ==========")

	// ---------- ① 函数能赋值给变量 ----------
	// 变量 add 的类型是：func(int, int) int
	add := func(a, b int) int { return a + b }
	fmt.Println("① 函数当值用：add(1, 2) =", add(1, 2))

	// ---------- ② 函数能当参数传（回调）----------
	// 前端对照：JS 的 arr.map(x => x * 2) 就是传回调
	apply(3, 4, add)                    // 传函数变量
	apply(10, 2, func(a, b int) int {   // 直接传匿名函数（字面量）
		return a * b
	})

	// ---------- ③ 函数能当返回值：工厂函数，返回一个"记住配置"的函数 ----------
	double := makeMultiplier(2) // 返回一个"乘 2"的函数
	triple := makeMultiplier(3) // 返回一个"乘 3"的函数
	fmt.Println("③ 工厂函数：", 5, "×2 =", double(5), "，5 ×3 =", triple(5))

	// ---------- ④ 闭包捕获：内层函数"记住并修改"外层变量 ----------
	counter := makeCounter()
	fmt.Println("④ 闭包计数：", counter(), counter(), counter()) // 1 2 3

	// ---------- ⑤ 中间件模式（Gin 的雏形）----------
	// logFactory 外层记住前缀，返回的闭包在真正干活时用这个前缀
	infoLog := logFactory("[INFO]")
	infoLog("用户登录成功")
	infoLog("订单已创建")

	// ---------- ⑥ 循环里闭包捕获循环变量（经典坑）----------
	fmt.Println("\n⑥ 循环闭包捕获：")
	var funcs []func()
	for i := 0; i < 3; i++ {
		i := i // 拷贝当前轮的 i（Go 1.22+ 已自动修，老版本必须这样写）
		funcs = append(funcs, func() {
			fmt.Println("   闭包看到 i =", i)
		})
	}
	for _, f := range funcs {
		f() // 输出 0 1 2：每个闭包记住的是自己那一轮的 i 副本
	}
	//! 前端对照：JS 老经典坑 —— var 循环里 setTimeout 全打印同一个值
	//!   Go 1.22 前一样：不拷贝 i，所有闭包记住同一个循环变量，全打印 3
	//!   Go 1.22+ 已修复（每轮一个独立变量），但面试还爱问老版本行为

	fmt.Println("\n🔑 一句话：函数是值 → 闭包记住外层 → 中间件 = 工厂函数返回闭包")
}

// apply 接收一个函数类型的参数 op，在里面回调它
func apply(a, b int, op func(int, int) int) {
	fmt.Println("② 回调：", a, "和", b, "运算结果 =", op(a, b))
}

// makeMultiplier 工厂函数：返回闭包，闭包记住外层 factor
func makeMultiplier(factor int) func(int) int {
	return func(x int) int { return x * factor } // factor 被捕获，不会丢
}

// makeCounter 返回闭包：闭包记住并累加外层 count
func makeCounter() func() int {
	count := 0
	return func() int {
		count++ // 修改外层变量（捕获 = 引用，呼应 01 说的"引用"语义）
		return count
	}
}

// logFactory 中间件雏形：外层记住前缀，返回的闭包干活时用
func logFactory(prefix string) func(msg string) {
	return func(msg string) {
		fmt.Println("⑤", prefix, msg)
	}
}
