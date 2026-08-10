// 04_waitgroup.go —— WaitGroup：正规地"等所有协程跑完"
//
// 前端对照（先记这句）：
//
//	WaitGroup ≈ Promise.all —— 等一批任务全部完成，再继续往下走
package main

import (
	"fmt"
	"sync"
)

// demoWaitGroup 演示用 WaitGroup 等所有 goroutine 完成（替代 01 的 time.Sleep 笨办法）
func demoWaitGroup() {
	fmt.Println("========== 04 WaitGroup：正规等待 ==========")

	// 01 节我们用 time.Sleep 等，那是笨办法：睡多久全凭猜
	// 正规武器：sync.WaitGroup —— 一个"计数器"
	//   Add(1)：登记一个要等的任务（计数器 +1）
	//   Done()：任务干完了（计数器 -1）
	//   Wait()：阻塞到计数器归 0，也就是所有人都干完了
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // 先登记，再开协程 —— 顺序不能反
		go func(n int) {
			defer wg.Done() // 干完一定通知；defer 保证中途出错也会 -1，不漏
			fmt.Println("干活了", n)
		}(i)
	}

	wg.Wait() // 5 个全部 Done 之前，这行后面的代码不会执行
	fmt.Println("5 个协程全部跑完 ✅")

	// 融会贯通：wg.Done() 用 defer 包 —— 呼应你之前学过的"收尾必做"思想
	// defer = "函数结束前最后执行"。手写 wg.Done() 在 return 前一行很容易漏，
	// 漏了 = 计数器永不归 0 → Wait() 永远死等 → 死锁

	//! 常见 bug ①：忘写 Done() → 计数器永不归 0 → Wait() 死等
	//! 常见 bug ②：Add 写在 goroutine 内部 → 协程还没 Add，Wait 就提前放行了
	//! 规则一句话：Add 在开协程前调，Done 在协程里靠 defer 调

	//? 思考：把 go func(n int){...}(i) 的传参改成闭包捕获 i，会怎样？
	//? 答：Go 1.22+ 已修复循环变量问题（每轮一个独立 i），结果一样；
	//?     但"传值=快照 / 捕获=引用"的语义（pointers/02）仍是面试常考点

	// ---------- 预告 Day7：数据竞争 ----------
	// 想象 5 个协程同时做 count++，结果可能不是 5！
	// 多个协程同时读写同一内存 = 数据竞争，结果不确定
	// 解法预告：sync.Mutex 锁，或干脆用 channel 通信。Go 的格言（面试常考）：
	//   "Don't communicate by sharing memory; share memory by communicating."
	//   —— 不要靠共享内存来通信，要靠通信来共享内存
}
