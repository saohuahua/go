// 03_struct.go —— 结构体指针 + 方法接收者（前端最容易懵的地方）
//
// 铺垫：先讲清"值拷贝"对结构体意味着什么
//   结构体 = 你自己定义的、打包好的一堆变量（类似 TS 的 interface + class）
//   把结构体传给函数/方法时，默认是【整个拷一份】！
//   改拷贝不会影响原结构体 —— 和 JS 的对象引用【完全相反】！
//
//! ⚠️ 前端最经典的坑：
//!   JS：const user = {name:'小明'};  user.name = '小红'   → 直接改成功
//!   Go：值接收者方法里改 user.Name                        → 【外面没变】！
//!   Go：指针接收者方法里改                                → 外面变
package main

import "fmt"

// User 演示用结构体（定义在包级别，多个函数都能用）
type User struct {
	Name string
	Age  int
}

// demoStruct 演示结构体传值 vs 传指针、值接收者 vs 指针接收者
func demoStruct() {
	fmt.Println("========== 03 结构体指针 ==========")

	user := User{Name: "小明", Age: 18}


	// ---------- 1. 函数传参：值拷贝 vs 指针 ----------
	changeAgeByValue(user) // 传值：函数里改的是拷贝
	fmt.Println("传值改年龄后    user.Age =", user.Age) // 18 ← 没变

	changeAgeByPointer(&user) // 传指针：函数里改的是本体
	fmt.Println("传指针改年龄后  user.Age =", user.Age) // 25 ← 变了！

	//! 为什么传值不行？结构体是一整块数据，传值 = 整块复制一份
	//! 函数里改的是副本，改完就扔了。想改原结构体 → 必须传指针 &user
	//* 前端对照：JS 对象传参就能改（隐式引用）；Go 结构体默认"值"，
	//*          要 JS 的效果必须显式 & —— 这是 Go 和 JS 最大的思维差异


	// ---------- 2. 方法接收者：值接收者 vs 指针接收者 ----------
	// 方法 = 挂在类型上的函数（类似 JS class 里的方法）
	user.SetNameValue("小红") // 值接收者：方法内改的是拷贝
	fmt.Println("\n值接收者 SetNameValue 后   user.Name =", user.Name) // 小明 ← 没变！

	user.SetNamePointer("小红") // 指针接收者：方法内改的是本体
	fmt.Println("指针接收者 SetNamePointer 后 user.Name =", user.Name) // 小红 ← 变了！

	//! 面试高频题：值接收者和指针接收者怎么选？
	//! 答：方法里要改字段 → 指针接收者；只是读字段 → 值接收者
	//? 思考：既然值接收者改不了外面，那它还有什么用？
	//? 答：① 只读字段的方法（方法不该改状态，更安全）
	//?     ② 小结构体拷贝开销可忽略；③ 想让方法操作"快照"时
	// 补充经验：结构体很大时一律用指针接收者（避免每次调用整块拷贝）
}

// changeAgeByValue 传值：形参 u 是 user 的完整拷贝
func changeAgeByValue(u User) {
	u.Age = 25
}

// changeAgeByPointer 传指针：形参 u 指向 user 本体
func changeAgeByPointer(u *User) {
	u.Age = 25 // 语法糖：u.Age 等价于 (*u).Age，Go 自动帮你解引用
}

// SetNameValue 值接收者：方法内改的是拷贝
func (u User) SetNameValue(name string) {
	u.Name = name
}

// SetNamePointer 指针接收者：方法内改的是本体
func (u *User) SetNamePointer(name string) {
	u.Name = name
}
