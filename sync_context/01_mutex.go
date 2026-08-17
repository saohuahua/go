// 01_mutex.go —— Mutex 锁：解决"数据竞争"（Day6 04 节埋的坑，本章最核心）
//
// 这一章的前端对照几乎为零：JS 是单线程，变量读写天然安全，
// 根本不会遇到这个坑。所以本章从头讲 —— 问题先来，解法后到。
//
// 全新写法（之前没见过，正文第一次出现处都标了"新写法"）：
//   sync.Mutex / .Lock() / .Unlock()、defer mu.Unlock()、sync.RWMutex / RLock / RUnlock
package main

import (
	"fmt"
	"sync"
)

// demoMutex 演示：数据竞争长什么样 → 用 Mutex 修复 → RWMutex 读多写少
func demoMutex() {
	fmt.Println("========== 01 Mutex：锁 ==========")

	// ---------- 问题：先看清"数据竞争" ----------
	// count++ 看着是一行，其实是三步：读 count → 加 1 → 写回 count
	// 两个协程同时做，可能"读到同一个旧值"，写回去互相覆盖 → 结果丢了
	count := 0
	var wg sync.WaitGroup // WaitGroup 是 Day6 04 学过的，这里只用它"等全部跑完"
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++ // 读-改-写三步，多个协程同时做 → 互相覆盖
		}()
	}
	wg.Wait()
	fmt.Println("① 不加锁：5 个协程各 +1，结果 =", count, "（常常不是 5！）")
	// 多跑几遍看：可能是 5，也可能是 3、4 —— 每次都不一样，这就是"不确定"

	//! 数据竞争是 Go 最阴险的 bug：编译不报错、运行才偶发、还难复现
	//! 实操检测：go run -race . （竞态检测器）—— 输出里出现 race 字样 = 有竞争
	//! 回去跑 Day6 的 04 例子，用 -race 大概率也会报警

	// ---------- 解法：互斥锁 sync.Mutex ----------
	// 锁 = 一个"通行证"：拿到的人独占一段代码，其他人只能门口排队
	locked := 0
	var mu sync.Mutex // 零值即可直接用，不用初始化
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()    //! 新写法：抢锁。抢不到就阻塞排队等
			locked++     // 被锁保护的"临界区"：同一时刻只有一个协程能进来
			mu.Unlock()  //! 新写法：还锁。后面排队的人才能进来
		}()
	}
	wg.Wait()
	fmt.Println("② 加锁后：结果稳定 =", locked, "✅ 每次都一样")

	// 生产环境更稳的写法：Lock 之后紧跟 defer mu.Unlock()
	//   好处：函数中途 return 或 panic，defer 也会把锁还掉，绝不忘解锁
	//! 新写法：defer mu.Unlock() —— 把"解锁"登记到函数结束时自动执行
	//   呼应：和 wg.Done() 用 defer 是同一个道理 —— "收尾必做"交给 defer

	// ---------- 升级：RWMutex —— 读多写少的场景 ----------
	// 普通 Mutex：读和写互相排斥，读的人多了也互相排队 → 浪费
	// RWMutex：多个读者可同时进（读读不互斥），只有写者独占
	fmt.Println("\n③ RWMutex 读多写少：")
	var rw sync.RWMutex
	value := 0
	for i := 0; i < 3; i++ { // 3 个读者同时读
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			rw.RLock()     //! 新写法：读锁（共享）—— 读者之间不互斥
			fmt.Println("  读者", n, "读到 value =", value)
			rw.RUnlock()   //! 新写法：还读锁。RLock 配 RUnlock，Lock 配 Unlock
		}(i)
	}
	wg.Wait()
	fmt.Println("  注意看输出：三个读者的打印是穿插的（同时进来的证据）")

	//! 三个必记的坑：
	//! 1. Lock 后忘 Unlock → 其他协程全堵在门口 → 死锁。用 defer Unlock 兜底
	//! 2. 锁别放进 struct 再整体拷贝 → 复制的是"锁的副本"，等于没锁
	//! 3. 锁粒度过大（把不用保护的代码也圈进来）→ 并发变串行，性能白给

	//? 思考：为什么 RWMutex 敢让"读读不互斥"？
	//? 答：读不改数据，多个人一起读不会打架；只要有人在写，
	//?     别人读写都可能读到写了一半的数据，所以才要互斥

	fmt.Println("\n🔑 Mutex = 给共享内存装“门禁”：同一时刻只放一个人进去")
}
