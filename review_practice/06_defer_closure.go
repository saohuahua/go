// 06_defer_closure.go —— defer 保证收尾；闭包把可配置行为“记”下来
package main

import (
	"fmt"
)

func demoDeferClosure() {
	fmt.Println("\n========== 06 defer + 闭包 ==========")
	check("defer 后进先出", deferOrder(), "[关闭数据库 关闭文件]")
	check("recover 接住 panic", safeRun(func() { panic("练习 panic") }), "练习 panic")

	withPrefix := newTitlePrefixer("[Go] ")
	check("闭包记住 prefix", withPrefix("复习"), "[Go] 复习")
}

// TODO 06.1：借助两个 defer，返回 ["关闭数据库", "关闭文件"]。
// 提示：defer 后注册先执行；命名返回值能在 defer 中被修改。
func deferOrder() (order []string) {
	return nil
}

// TODO 06.2：运行 fn；fn panic 时 recover 并返回 panic 内容的字符串；没 panic 返回 ""。
// //! recover 只能在 deferred 函数里生效。
func safeRun(fn func()) (message string) {
	return ""
}

// TODO 06.3：返回一个闭包。闭包接收 title，并加上创建时传入的 prefix。
// 前端对照：const newPrefixer = prefix => title => prefix + title。
func newTitlePrefixer(prefix string) func(string) string {
	return func(title string) string {
		return title
	}
}
