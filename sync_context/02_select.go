// 02_select.go —— select：从多个 channel 里挑一个"最先就绪"的
//
// 前端对照只有一处能迁移：select ≈ Promise.race（谁先落定用谁）。
// 其余语法是全新写法，正文第一次出现处都标了"新写法"。
//
// 先讲为什么要它（快速入门看这段）：
//   Day6 的 channel 是"一根水管一对一"：ch <- v / v := <-ch。
//   但真实场景常是一条协程同时盯着好几根水管（网络返回、用户操作、超时信号），
//   总不可能给每根水管都单独开一条协程去等。select 就是为"多路监听"而生的。
package main

import (
	"fmt"
	"time"
)

// demoSelect 演示 select 多路复用 + 超时 + default + for 持续监听
func demoSelect() {
	fmt.Println("========== 02 select：多路复用 ==========")

	//! 新写法：select 语句 —— 语法长得像 switch，但 case 全是 channel 操作
	//   select {
	//   case v := <-ch1:           // ch1 有数据 → 走这里
	//   case v := <-ch2:           // ch2 有数据 → 走这里
	//   case ch3 <- v:             // ch3 能塞进去 → 走这里（发送也算就绪）
	//   case <-time.After(1 * time.Second): // 全等太久 → 走这里（超时）
	//   }
	//! 规则：多个 case 同时就绪时，随机挑一个 —— 别依赖顺序，面试爱考

	// ---------- 场景①：两个 channel 都来消息，随机收一个 ----------
	fast := make(chan string, 1)
	slow := make(chan string, 1)
	fast <- "快的消息"
	slow <- "慢的消息"
	select {
	case msg := <-fast:
		fmt.Println("① 收到 fast：", msg)
	case msg := <-slow:
		fmt.Println("① 收到 slow：", msg)
	}
	// 两个都就绪 → 随机二选一。多跑几遍看，不一定每次都是同一个

	// ---------- 场景②：超时 —— 最实用的一个 ----------
	//! 新写法：time.After(d) —— 一根"过 d 之后才来水"的 channel
	// 用途：等一个"可能永远不来的结果"，等太久就放弃，不让协程死等
	never := make(chan int) // 一根永远不会有人发的 channel
	select {
	case v := <-never:
		fmt.Println("② 收到了", v)
	case <-time.After(500 * time.Millisecond): // 500ms 没等到就走这里
		fmt.Println("② 等了 500ms 没等到，超时放弃")
	}
	//! 把这条超时 case 去掉试试 → 所有 case 都阻塞 → 主协程直接死锁

	// ---------- 场景③：default —— 非阻塞检查 ----------
	// 所有 case 都没就绪 + 有 default → 直接走 default，一秒都不等
	// 适合"顺手看一眼有没有，没有就干别的"，而不是傻等
	empty := make(chan int, 1) // 一根空 channel
	select {
	case v := <-empty:
		fmt.Println("③ 有值", v)
	default:
		fmt.Println("③ empty 现在是空的 → 走 default，立刻返回")
	}

	// ---------- 场景④：for + select 持续监听 + 优雅退出 ----------
	// 经典模式：协程里死循环监听，收到停止信号才 break（03 的 Context 也干这个）
	stop := make(chan struct{}) //! 新写法：chan struct{} —— 只当"信号"，不传数据
	done := make(chan struct{}) // 协程退出的确认信号
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("④ 收到停止信号，退出循环")
				close(done) // 通知主函数：我真的退出了（比 sleep 猜时间正规）
				return
			case <-time.After(200 * time.Millisecond):
				fmt.Println("④ 每 200ms 例行检查一次…")
			}
		}
	}()
	time.Sleep(700 * time.Millisecond)
	close(stop) // 发停止信号（呼应 Day6 03：close 由发送方负责）
	<-done      // 阻塞等协程确认退出 —— 呼应 Day6 02 的 channel 同步

	//? 思考：select 和 switch 有什么区别？
	//? 答：switch 是"值匹配"（相等才走）；select 是"就绪选择"（谁先到走谁），
	//?     而且 case 只能是 channel 收发操作。长得像，本质完全不同

	fmt.Println("\n🔑 select = 多根水管的“分流器”：谁先来水处理谁，配超时防死等")
}
