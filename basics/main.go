// main.go —— 基础语法补课（接口 / error / defer / 闭包）学习程序入口
//
// 本文件夹是进 Gin 框架之前的前置补课，共 4 个演示文件 + 本入口文件：
//
//   01_interface.go —— 接口：隐式实现、多态、any、类型断言（框架的契约层）
//   02_error.go     —— 错误处理：返回值模式、%w 包裹、errors.Is/As
//   03_defer.go     —— defer：栈序、参数立即求值、命名返回值、recover
//   04_closure.go   —— 闭包：函数当值、捕获、中间件模式（Gin 中间件的魂）
//
// 运行方法（二选一）：
//   cd basics && go run .
//   go run ./basics    （在项目根目录执行）
//
// 学习方法建议：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 04 闭包是中间件的基础，建议和后面 Gin 的中间件一起复习
package main

import "fmt"

func main() {
	demoInterface() // 01 接口：能力契约
	// demoError()     // 02 错误：返回值模式
	// demoDefer()     // 03 defer：收尾保险
	// demoClosure()   // 04 闭包：中间件的魂

	fmt.Println("\n🎉 4 个前置基础补完，可以进 Gin 了！（struct tag / JSON 到 Gin 里边用边学）")
}
