// 03_context.go —— Context：给协程装一根"刹车线"（取消 / 超时 / 传值）
//
// 这是全书最抽象的一节，前端几乎没有可对照的知识，所以慢慢来：
// 先看清它解决的问题，再学它长什么样。
//
// 问题引入（快速入门必看）：
//   Day6 只教了怎么"开"协程，没教怎么"关"。
//   go f() 一旦跑起来，除非 f 自己 return，没人能叫停它。
//   真实项目里"开出去收不回来"的协程 = 协程泄漏，开多了程序就崩。
//   Context 就是那根刹车线：上游一拉（取消 / 超时），下游所有协程一起停。
//
// 前端对照仅一处：fetch 的 AbortController —— abort() 一下，请求链全取消。
// 其余全是全新写法，正文第一次出现处都标了"新写法"。
package main

import (
	"context"
	"fmt"
	"time"
)

//! 新写法：type ctxKey string —— 定义一个"自定义类型"，用来当 WithValue 的 key
// 为什么不用现成的 string / int？防止和别处传的 key 撞车、互相覆盖（面试爱问）
type ctxKey string

// demoContext 演示 Context 三大作用：取消、超时、传值
func demoContext() {
	fmt.Println("========== 03 Context：取消 / 超时 / 传值 ==========")

	// ---------- ① WithCancel：手动喊停 ----------
	//! 新写法：context.Background() —— 根节点，所有 Context 的起点
	//! 新写法：context.WithCancel(父) —— 返回 ctx 和 cancel() 函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 铁律：生成了就必须 cancel，要么手动调，要么 defer

	go worker(ctx, "工人A")
	time.Sleep(100 * time.Millisecond)
	cancel() // 喊停！ctx.Done() 被关闭，worker 里的 select 立刻收到
	fmt.Println("① 已调用 cancel()，worker 应该马上退出")
	time.Sleep(50 * time.Millisecond) // 留 50ms 让 worker 打印（正规等法后面讲）

	// ---------- ② WithTimeout：到点自动取消 ----------
	//! 新写法：context.WithTimeout(父, 时长) —— 到点自动取消，不用手动调
	// 最常用场景：给"可能很慢的操作"设上限，超时自动放弃
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2() // 就算超时已自动取消，defer 再调一次也无害（重复 cancel 安全）

	go slowTask(ctx2, "慢任务")
	fmt.Println("② 慢任务只有 500ms 的额度，超了就自动放弃")
	time.Sleep(700 * time.Millisecond) // 等慢任务把"超时退出"打出来

	// ---------- ③ WithValue：沿调用链传元信息 ----------
	// 场景：把 requestId / 用户ID 这类"请求元信息"一路传给下游函数
	//! 新写法：context.WithValue(父, key, value) / ctx.Value(key)
	ctx3 := context.WithValue(context.Background(), ctxKey("requestId"), "req-12345")
	child(ctx3)
}

// worker 用 select 监听 ctx.Done()：被取消就优雅退出，而不是被强行杀掉
//! 新写法：ctx.Done() —— 一个 channel，cancel 或超时的那一刻被关闭
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done(): // 刹车线被拉 → 这里立刻触发
			fmt.Println(" ", name, "收到取消信号，正在清理并退出")
			return
		default:
			time.Sleep(20 * time.Millisecond) // 模拟干活：干一小段再回来检查
		}
	}
}

// slowTask 演示超时：到点了走 ctx.Done() 分支，不再傻等
func slowTask(ctx context.Context, name string) {
	select {
	case <-ctx.Done(): // 500ms 后这里触发
		fmt.Println(" ", name, "超时了，停止工作")
	case <-time.After(3 * time.Second): // 假设这个任务本来要 3 秒
		fmt.Println(" ", name, "干完了") // 走不到：3s > 500ms 的超时额度
	}
}

// child 演示 WithValue 传值：把 ctx 里的元信息取出来
func child(ctx context.Context) {
	v := ctx.Value(ctxKey("requestId")) // 会一路往父级找 key
	fmt.Println("③ 下游函数取到 requestId =", v)
}

//! 三个必记的坑：
//! 1. 忘 cancel() → 底层定时器不释放 → 资源泄漏。铁律：生成即 defer cancel()
//! 2. WithValue 的 key 别用内置 string / int → 可能和别处撞车。用自定义类型（ctxKey）
//! 3. 别把 ctx 存进 struct 字段 → Go 惯例：ctx 永远是函数第一个参数，顺着调用链往下传
