// main.go —— goroutine + channel（并发）学习程序入口
//
// 本文件夹是一个完整的 Go 学习项目，共 4 个演示文件 + 本入口文件：
//
//   01_basics.go    —— go 启动协程、主函数退出 = 协程全没（最经典的坑）
//   02_channel.go   —— channel 收发、无缓冲 = 同步阻塞（本章核心）
//   03_buffered.go  —— 有缓冲 channel + 生产者-消费者 + 关闭 channel
//   04_waitgroup.go —— WaitGroup 正规等待 + 数据竞争预告
//
// 运行方法（二选一）：
//   cd goroutine && go run .
//   go run ./goroutine    （在项目根目录执行）
//
// 学习方法建议（和 slices / map / pointers 一样）：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 本章和前几章最大的不同：并发是"不确定的"
//      —— 同一段代码多跑几遍，输出顺序可能不一样，这是正常的，不是 bug
package main

import "fmt"

func main() {
	demoBasics()     // 01 启动 goroutine
	demoChannel()    // 02 channel 收发
	demoBuffered()   // 03 有缓冲 + 生产者-消费者
	demoWaitGroup()  // 04 WaitGroup 等待

	fmt.Println("\n🎉 全部演示完成！建议打开每个文件，动手改一改再跑一遍")
}
