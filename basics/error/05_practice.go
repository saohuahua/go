// 05_practice.go —— 后端实战规范：不吞错误 / 边界打日志 / panic·recover
//
// 到这一节，你已经会用错误了。但"怎么写才规范"才是面试官看代码的第一眼。
// 三个最重要的习惯 + 一个"几乎不用、但面试会问"的 panic/recover。
package main

import (
	"errors"
	"fmt"
)

// loadConfig 模拟读取配置：这里故意让它失败，演示下面的处理方式
func loadConfig() error {
	return errors.New("配置文件缺失")
}

func demoPractice() {
	fmt.Println("========== 05 后端实战规范：不吞错误 / 边界打日志 / panic·recover ==========")

	// ---------- ① 不吞错误：错误必须"处理或上报"，二选一 ----------
	// 反模式（千万别学）：错误拿到手就当没看见，bug 直接蒸发
	//   if err := loadConfig(); err != nil {
	//       // 什么都不做 —— 配置没加载，后面全跑偏，你还找不到原因
	//   }
	// 正解二选一：
	//   a. 向上返回：return fmt.Errorf("加载配置失败: %w", err)
	//   b. 记录日志：log.Printf("加载配置失败: %v", err)
	if err := loadConfig(); err != nil {
		fmt.Printf("① 正确处理：把错误上报/记录 —— %v\n", err)
	}

	// ---------- ② 日志只在上层打，下层只负责 return+包裹 ----------
	// 常见错误：每一层都 log.Printf，同一个错误打 3 遍，日志爆炸
	// 规范：下层（repo/service）用 %w 带上下文传上去；上层（handler/入口）统一打日志
	fmt.Println("② 日志只在上层打一次 —— 下层用 %w 带上下文（03），上层统一记录")

	// ---------- ③ panic/recover：几乎用不到，但面试会问（提一嘴）----------
	// panic = 程序炸了（不可恢复），只留给两种场景：
	//   1. 不可能发生的错（switch 没写 default 且前面全 return）
	//   2. 程序初始化失败（配置/数据库连不上，起不来直接崩）
	// 业务错误一律 return error，绝对不要 panic —— panic 会崩掉整个进程！
	fmt.Println("③ panic 只留给「不该发生的事」；业务错误一律 error")
	fmt.Println("③ Gin 的 Recovery 中间件会自动接住 panic，防止一个请求 500 崩掉整个服务")

	// ---------- ④ recover：在 defer 里接住 panic（呼应 defer 专题）----------
	// 匿名函数立即执行（IIFE，前端也有）：defer+recover 必须在同一个函数内才有效
	func() {
		defer func() { // defer：这个函数退出前执行（03_defer 专题讲过）
			if r := recover(); r != nil { // recover 接住 panic，程序不崩
				fmt.Printf("④ defer 里 recover 接住了 panic：%v —— 程序没崩 ✅\n", r)
			}
		}()
		fmt.Println("④ 正常执行…")
		panic("这里本不该 panic 的") // panic 之后函数立即中断，后面任何代码都不会执行
	}()
	fmt.Println("④ 回到主流程继续跑（recover 接住后，外层不受影响）")

	// ---------- ⑤ errors.Join（Go 1.20+，提一嘴）：把多个错误合并成一个 ----------
	// 场景：表单校验想一次性返回所有问题，而不是"遇到一个就 return"
	all := errors.Join(
		errors.New("用户名太短"),
		errors.New("密码格式不对"),
	)
	fmt.Println("⑤ errors.Join 合并多个错误：", all)

	fmt.Println("\n🔑 规范三连：错误不吞掉 / 日志只在上层 / 业务错不 panic；recover 交给 Gin 的 Recovery")
}
