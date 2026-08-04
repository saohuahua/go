// main.go —— 指针（Pointer）学习程序入口
//
// 本文件夹是一个 Go 学习项目，共 4 个演示文件 + 本入口文件：
//
//   01_basics.go       —— 认识指针：& 取地址、* 解引用、nil 指针
//   02_function.go     —— 指针传参：为什么函数里改了，外面变
//   03_struct.go       —— 结构体指针：值接收者 vs 指针接收者（必踩坑）
//   04_when_to_use.go  —— 总结：什么时候用指针 + new + 面试速背
//
// 运行方法（二选一）：
//   cd pointers && go run .
//   go run ./pointers    （在项目根目录执行）
//
// 学习方法建议（和 slices / map 一样）：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 每节都配了和 JS/TS 的对照 —— 指针是前端没有的概念，多跑几遍
package main

import "fmt"

func main() {
	demoBasics()   // 01 认识 & 和 *
	demoFunction() // 02 指针传参
	demoStruct()   // 03 结构体指针
	demoSummary()  // 04 总结 + 面试速背

	fmt.Println("\n🎉 全部演示完成！建议打开每个文件，动手改一改再跑一遍")
}
