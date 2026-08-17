// 03_defer.go —— defer：函数"收尾动作"的保险（你一直在用，今天学透）
//
// 前端对照一句话（先记这个）：
//
//	JS：try { ... } finally { ... } —— 不管中间怎样，finally 保证执行
//	Go：defer 函数() —— 注册一个"离开当前函数时执行"的动作
//	两个区别要注意：Go 是栈（后注册先执行），且参数在 defer 那行就求值了
//
// 为什么重要：关文件、关连接、还锁、cancel()、wg.Done() 全靠它兜底。
//   不学透 defer，资源清理就是"漏一次埋一个 bug"，还特别难查。
//
// 你早就用过的 defer：wg.Done()（goroutine/04）、mu.Unlock()（sync_context/01）、
// cancel()（sync_context/03）—— 今天补上它背后的全部机制。
package main

import "fmt"

func demoDefer() {
	fmt.Println("========== 03 defer：延迟执行 + 收尾保险 ==========")

	// ---------- ① 基本用法：函数结束统一执行（类似 finally）----------
	basicDefer()

	// ---------- ② 栈序：多个 defer = 后注册的先执行（LIFO）----------
	stackDefer()

	// ---------- ③ 参数立即求值（面试必考！）----------
	argEval()

	// ---------- ④ defer + 命名返回值 ----------
	fmt.Println("④ 返回值被 defer 改成了：", namedReturn())

	// ---------- ⑤ defer + recover 接住 panic ----------
	recoverDemo()
}

func basicDefer() {
	defer fmt.Println("① [defer] 我最后执行")
	fmt.Println("① 正文第一行")
	fmt.Println("① 正文第二行")
}

func stackDefer() {
	defer fmt.Println("② 最先注册 → 最后执行")
	defer fmt.Println("② 中间注册")
	defer fmt.Println("② 最后注册 → 最先执行")
	fmt.Println("② 函数体跑完，开始倒着执行 defer 栈")
}

func argEval() {
	x := 10
	defer fmt.Println("③ defer 里看到的 x =", x)          // ① 立即把 x=10 拷贝好
	x = 100                                              // ② 再改，不影响已经拷好的值
	fmt.Println("③ 函数里把 x 改成：", x)
	// 想取"最新值"？用闭包包一层（呼应 04 闭包）：
	defer func() { fmt.Println("③ 闭包版 defer 看到最新 x =", x) }()
	//! 坑：defer 的参数【在 defer 那一行就求值】，不是函数结束时求值
	//! 面试必考题：defer fmt.Println(x) 之后 x 再变，打印的还是旧值
}

// namedReturn 演示 defer + 命名返回值：return 之后、函数返回之前，defer 还能改它
func namedReturn() (n int) {
	defer func() {
		n += 100 // return 先填好 1，这里再 +100
	}()
	return 1 // 返回值先定为 1，defer 改成 101
}

func recoverDemo() {
	defer func() {
		if r := recover(); r != nil { // recover 只能在 defer 里调用才有效
			fmt.Println("⑤ 接住 panic：", r, "—— 程序没崩，继续往下走")
		}
	}()
	fmt.Println("⑤ 开始一个可能 panic 的操作...")
	panic("某处出大问题了") // 正常会崩掉整个程序，被上面的 defer 接住
	// panic 后面的代码根本不会执行 —— 程序立刻跳去执行已注册的 defer
}
