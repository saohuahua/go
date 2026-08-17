// main.go —— 错误处理专题（Error）学习程序入口
//
// 错误是 Go 后端代码里出现频率最高的类型：if err != nil 能占业务代码的 30%~50%。
// 面试官看 Go 代码，第一眼就看错误处理规不规范 —— 这是"脸面"级别的主题。
//
// 本文件夹把错误处理讲透（basics/02_error.go 是速览版，这里是完整版）：
//
//   01_what.go     —— error 是什么：一个接口 + 返回值模式（Go 没有 try/catch）
//   02_make.go     —— 造错误三招：errors.New / fmt.Errorf / 哨兵错误
//   03_wrap.go     —— 错误链：%w 包裹 + errors.Is / errors.As（面试高频）
//   04_custom.go   —— 自定义错误类型：错误里带结构化数据（Gin 校验错误同款）
//   05_practice.go —— 后端实战规范：不吞错误 / 边界打日志 / panic·recover
//
// 主线就一条：错误长什么样 → 怎么造 → 怎么带上下文 → 怎么判断错误链 → 实战怎么写。
//
// 运行方法（二选一）：
//   cd basics/error && go run .
//   go run ./basics/error    （在项目根目录执行）
package main

import "fmt"

func main() {
	demoWhat()     // 01 error 是什么 + 返回值模式
	demoMake()     // 02 造错误：errors.New / fmt.Errorf / 哨兵错误
	demoWrap()     // 03 错误链：%w 包裹 + Is / As
	demoCustom()   // 04 自定义错误类型
	demoPractice() // 05 后端实战规范 + panic/recover

	fmt.Println("\n🎉 错误处理专题学完 —— 之后写 Gin handler，if err != nil 会用到手软！")
}
