// 03_buffered.go —— 有缓冲 channel：把"必须碰头"改成"先放队列里"
//
// 前端对照（先记这句）：
//
//	无缓冲 channel ≈ 两人必须当面交接
//	有缓冲 channel ≈ 家门口放了个储物柜：塞进去就走，对方空了再取
package main

import "fmt"

// demoBuffered 演示有缓冲 channel + 生产者-消费者 + 关闭 channel
func demoBuffered() {
	fmt.Println("========== 03 有缓冲 channel + 生产者-消费者 ==========")

	// make(chan int, 3)：第三个参数 = 缓冲容量（能塞 3 个不用等人收）
	buffered := make(chan int, 3)
	buffered <- 1
	buffered <- 2
	buffered <- 3
	fmt.Println("连续塞了 3 个：len =", len(buffered), " cap =", cap(buffered))

	// 融会贯通：channel 也有 len / cap！和 slice（slices/01）语义一模一样
	//   slice   ：len = 元素个数，cap = 底层数组容量
	//   channel ：len = 现在排队几个，cap = 总共能排几个

	//! 塞满再塞会怎样？→ 阻塞。在 main 里直接 = 死锁（解开下面一行试一次）：
	//! buffered <- 4 // panic: deadlock

	// ---------- 经典模式：生产者-消费者 ----------
	// 生产者：不断往管道里放任务；消费者：不断从管道里取任务
	// 两者只认识这个 channel，互相不认识 —— 这就是"解耦"
	jobs := make(chan int, 5)

	// 生产者：一条协程，放 5 个任务，放完 close 通知"没了"
	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
			fmt.Println("生产了任务", i)
		}
		close(jobs) // 关闭 = 告诉消费者：不会再有了
	}()

	// 消费者：main 里用 range 一直取，取到 channel 关闭自动结束
	for job := range jobs {
		fmt.Println("消费了任务", job)
	}
	fmt.Println("channel 已关闭，range 自动退出")

	// 注意看输出：生产和消费是交替穿插的，顺序不固定
	// → 和 01 的"10 个协程顺序随机"是同一个现象，现在你应该不惊讶了

	//* 前端对照：这就像一个"任务队列"—— 生产者往里 push，消费者往外 shift
	//          两边不用互相等，比 02 无缓冲的"必须碰头"松耦合得多

	// ---------- 读已关闭的 channel：comma-ok ----------
	closed := make(chan int)
	close(closed)
	zero, ok := <-closed
	fmt.Println("读已关闭的 channel：v =", zero, " ok =", ok) // 0 false

	// 融会贯通：v, ok := <-ch 的 ok 和 map（map/01）的 comma-ok 一模一样！
	//   map     ：ok = key 是否存在
	//   channel ：ok = 还有没有值（false = 已关闭且取空了）

	//! 反过来绝对不行：往已关闭的 channel 里发送 → panic: send on closed channel
	//! closed <- 1 // 解开会崩

	//? 思考：为什么"读已关闭"安全、"写已关闭"崩溃？
	//? 答：读取可以拿零值兜底（呼应 map 读不到返回零值）；
	//?     写入是"我断言这个管道还在用"，关了还写 = 逻辑错误，Go 直接爆给你看

	fmt.Println("\n🔑 有缓冲 = 解耦（塞完就走）；close 由发送方负责，收方只管读")
}
