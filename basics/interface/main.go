// main.go —— 接口专题（Interface）学习程序入口
//
// 接口是 Go 的"骨架"：标准库、Gin 框架、分层架构全靠它解耦。
// 本文件夹从定义到实战讲透，共 5 个演示文件 + 本入口文件：
//
//	01_define.go      —— 定义与隐式实现：契约、多态、嵌套、空接口 any
//	02_use.go         —— 三种使用场景：函数参数 / 结构体字段 / 工厂返回
//	03_switch_impl.go —— 换实现 + mock：repository 接口（分层架构核心）
//	04_stdlib.go      —— 标准库接口：Stringer / error / io.Writer
//	05_assertion.go   —— any + 类型断言 + type switch：从接口里掏类型
//
// 运行方法（二选一）：
//
//	cd basics/interface && go run .
//	go run ./basics/interface    （在项目根目录执行）
//
// 学习方法建议：
//  1. 先整篇跑一遍，看输出和注释对应着理解
//  2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//  3. 配合 basics/01_interface.go 的速览版一起看，效果更好
package main

import "fmt"

func main() {
	// demoDefine()     // 01 定义与隐式实现
	// demoUse() // 02 三种使用场景
	// demoSwitchImpl() // 03 换实现 + mock
	// demoStdlib()     // 04 标准库接口
	demoAssertion() // 05 断言 + any + type switch

	fmt.Println("\n🎉 接口专题学完 —— 之后看 Gin，会发现到处都是接口！")
}
