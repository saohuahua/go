// 04_stdlib.go —— 标准库的接口：Go 生态怎么用接口（看懂标准库，重点）
//
// Go 标准库本身就是接口的"活教材"，四个写代码必碰的：
//
//	fmt.Stringer  —— 你定义 String() string，fmt.Println 就能友好打印你的类型
//	error         —— 你定义 Error() string，你的类型就是错误（呼应 basics/02）
//	io.Writer     —— 你定义 Write(p []byte) (int, error)，任何地方都能写你
//	http.Handler  —— 你定义 ServeHTTP(w, r)，你就能被注册成路由（Gin 前瞻）
//
// 关键认识：这些接口都是"使用者"（标准库）定义的，你的类型根本不知道。
// 你只是恰好有这个方法，就自动能用了 —— 这就是隐式实现的力量。
package main

import (
	"bytes"
	"fmt"
	"strings"
)

// ---------- ① fmt.Stringer：实现 String()，打印就好看 ----------
// Money 金额类型：默认打印 "3.5" 很丑；实现 String() 后打印 "$3.50"
type Money float64

func (m Money) String() string {
	return fmt.Sprintf("$%.2f", float64(m))
}

// ---------- ② error 接口：实现 Error() string 就是错误 ----------
// 呼应 basics/02 的自定义错误：这里是同一个机制
type MyError struct{ Code int }

func (e *MyError) Error() string {
	return fmt.Sprintf("错误码 %d", e.Code)
}

// ---------- ③ io.Writer：自定义一个"会计数"的写入器 ----------
// bytes.Buffer 已经实现了 io.Writer，我们再写一个自己的：
type CountingWriter struct {
	count int // 记下总共写进来多少字节
}

// Write 方法签名必须和 io.Writer 契约完全一致：Write(p []byte) (int, error)
func (w *CountingWriter) Write(p []byte) (int, error) {
	w.count += len(p)
	return len(p), nil
}

func demoStdlib() {
	fmt.Println("========== 04 标准库接口：Stringer / error / io.Writer ==========")

	// ---------- ① Stringer ----------
	// fmt.Println 内部发现 Money 实现了 Stringer，就用你的 String() 打印
	fmt.Println("① 自己类型 + String()：", Money(3.5), Money(19.99))

	// ---------- ② error 接口 ----------
	err := &MyError{Code: 404}
	fmt.Println("② 自定义错误：", err)

	// ---------- ③ io.Writer：任何"能 Write 的东西"都能被写入 ----------
	// bytes.Buffer 满足 io.Writer → 往"内存缓冲区"写
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "写入内存缓冲")
	fmt.Println("③ bytes.Buffer 写进去：", strings.TrimSpace(buf.String()))

	// 自己的 CountingWriter 也满足 io.Writer → 同样的写法，换成"计数"
	cw := &CountingWriter{}
	fmt.Fprintln(cw, "hello") // fmt.Fprintln 只要求一个"能 Write 的东西"
	fmt.Println("③ CountingWriter：写进了", cw.count, "字节")

	// ---------- ④ http.Handler（Gin 前瞻，提一嘴）----------
	// 标准库：谁有 ServeHTTP(w, r)，谁就能被注册成 HTTP 处理器
	// Gin 的 HandlerFunc 就是适配这个接口 —— Week3 学 Gin 时天天见
	fmt.Println("④ http.Handler：Gin 的路由处理器都是它（学到 Week3 再展开）")

	fmt.Println("\n🔑 标准库接口 = 你的类型 + 一个方法 = 白嫖一个生态")
}
