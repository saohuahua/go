// 02_function.go —— 指针传参：为什么"函数里改了，外面变"
//
// 先复习你已经学过的（回想 slices/05 和 map/02）：
//   map           → 本身就是引用，传参天然能改外面
//   切片          → 内部藏着指针，改元素外面看得见（append 扩容除外）
//   int 等基本类型 → 传值（拷贝），函数里改的是副本，外面【不会变】
//
// 那如果我就是想改外面的基本类型变量呢？两个办法：
//   JS 的做法：return 接住返回值
//   Go 的新能力：传指针，函数里 *p = x 直接改原变量
//
//* 前端对照（关键差异）：
//*   JS：function f(x) { x = 100 } 改不了外面的 x（基本类型传值）
//*   Go：function f(p *int) { *p = 100 } 能改！这是 JS 没有的能力
package main

import "fmt"

// demoFunction 演示指针传参修改外部变量
func demoFunction() {
	fmt.Println("========== 02 指针传参 ==========")

	score := 60


	// ---------- 1. 传值：函数内改的是拷贝，外面不变 ----------
	passValue(score)
	fmt.Println("传值修改后   score =", score) // 60 ← 没变


	//! 2. 传指针：函数内 *p 改的是原变量（p 拿着 score 的门牌号）
	passPointer(&score)
	fmt.Println("传指针修改后 score =", score) // 95 ← 变了！

	fmt.Println("\n🔑 一句话：想让函数「改外面」，就传指针；不想被改，就传值（安全）")


	// ---------- 3. 补充：为什么 Go 不直接"默认引用"像 JS 一样？ ----------
	// Go 的选择：所有传参都是值拷贝，但【指针本身也是值】——
	// 传指针 = 把门牌号抄一份递过去，你和函数拿的是同一张门牌号 → 能改同一个房间
	// 好处：一眼能看出来"这个函数会不会改我的数据"（看签名里有没有 *）
	//* JS 对照：JS 是隐式的，你不知道传对象进去会不会被改；Go 看参数类型就知道
}

// passValue 传值：形参 s 是 score 的拷贝，改 100 只是改副本
func passValue(s int) {
	s = 100
}

// passPointer 传指针：形参 p 指向 score 本体，*p = 95 直接改它
func passPointer(p *int) {
	*p = 95
}
