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


	// ---------- 0. 快速回顾：之前学的 slice/map 传参 ----------
	//* 前端对照：JS "基本类型传值、对象传引用" —— Go 类似但更显式

	n := 42
	changeInt(n)
	fmt.Println("int 传值 →", n) // 42，没变（和 JS number 一样）

	s := []int{1, 2, 3}
	changeSlice(s)
	fmt.Println("slice 改元素 →", s) // [99 2 3]，变了！（header 里藏着指针）

	m := map[string]int{"a": 1}
	changeMap(m)
	fmt.Println("map 改 key →", m) // map[a:1 b:2]，变了！（本身就是引用）

	//! 🔑 所以：int 改不了外面，slice/map 能改 —— 那 int 想改怎么办？→ 传指针（往下看）


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
	//
	//* JS 对照：
	//*   JS 里 function f(obj) {} —— 你传个对象进去，看签名根本不知道它会不会改你的 obj
	//*   必须去【读函数体】才能发现 obj.name = "xxx" 这种偷偷修改
	//*
	//*   Go 里 func f(u User)  → 没有 * → 值拷贝 → 它改不了外面，放心传
	//*         func f(u *User) → 有  * → 传指针 → 它能改外面，调用时要注意
	//*   不用读函数体，光看签名里有没有 * 就知道了 —— 这就是"显式 > 隐式"
}

// ---------- 回顾用的辅助函数 ----------

func changeInt(x int)          { x = 0 }        // 拷贝，改不了外面
func changeSlice(s []int)      { s[0] = 99 }    // header 里的指针指向同一底层数组
func changeMap(m map[string]int) { m["b"] = 2 } // map 本身就是引用

// ---------- 指针传参的辅助函数 ----------

// passValue 传值：形参 s 是 score 的拷贝，改 100 只是改副本
func passValue(s int) {
	s = 100
}

// passPointer 传指针：形参 p 指向 score 本体，*p = 95 直接改它
func passPointer(p *int) {
	*p = 95
}
