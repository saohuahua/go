// main.go —— sync + select + context（Day7 并发控制）学习程序入口
//
// 本文件夹是 goroutine 的续篇，共 4 个演示文件 + 本入口文件：
//
//   01_mutex.go    —— Mutex 修数据竞争 + RWMutex（本章最核心）
//   02_select.go   —— select 多路复用 channel + 超时（核心）
//   03_context.go  —— Context 取消 / 超时 / 传值（核心，最抽象，慢点读）
//   04_once_map.go —— 次重点：sync.Once 只跑一次 + sync.Map 并发安全 map
//
// 运行方法（二选一）：
//   cd sync_context && go run .
//   go run ./sync_context    （在项目根目录执行）
//
// 学习方法建议（和 goroutine 一样，但本章更依赖"自己动手"）：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 本章是前端知识覆盖最少的一章，每节都按
//      "问题 → 工具 → 语法 → 实验"来写 —— 看不懂就回去重看"问题"部分
package main

import "fmt"

func main() {
	demoMutex()    // 01 锁：修数据竞争
	demoSelect()   // 02 select：多路复用
	demoContext()  // 03 Context：取消 / 超时 / 传值
	demoOnceMap()  // 04 Once + sync.Map

	fmt.Println("\n🎉 全部演示完成！Day7 学完，Go 并发三件套（锁 / select / Context）就齐了")
}
