// 01_define.go —— 接口的定义与隐式实现（语法基础）
//
// 接口 = 一组方法签名的集合 = 一份"能力契约"。
// 谁的方法签名满足这份契约，谁就自动实现了它 —— 不用写任何 implements 关键字。
//
// 前端对照一句话（先记这个）：
//
//	TS：interface Speaker {...}，类要写 implements Speaker 才认
//	Go：方法签名对上就自动算实现（结构类型 / duck typing）
//	相当于 TS 的 structural typing 在 Go 里是默认行为
package main

import "fmt"

// ---------- 1. 定义接口 ----------

// Speaker 接口：契约 = "必须有一个 Speak() string 方法"
type Speaker interface {
	Speak() string
}

// Runner 接口：契约 = "必须有一个 Run() string 方法"
type Runner interface {
	Run() string
}

// Animal 接口：接口也能嵌套组合（把多个契约拼成一个）
// 想实现 Animal，必须先满足 Speaker + Runner 两个契约
type Animal interface {
	Speaker
	Runner
}

// ---------- 2. 实现者：注意！没写任何 implements ----------

// Dog 有 Speak + Run 两个方法 → 自动实现 Speaker / Runner / Animal
type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + "：汪汪！" }
func (d Dog) Run() string   { return d.Name + " 撒腿就跑" }

// Cat 也一样，方法签名对上就自动实现
type Cat struct{ Name string }

func (c Cat) Speak() string { return c.Name + "：喵呜~" }
func (c Cat) Run() string   { return c.Name + " 轻巧地走" }

// Duck 只实现了 Speak，没实现 Run
// → 它满足 Speaker，但不满足 Runner / Animal（这就是"部分满足契约"）
type Duck struct{ Name string }

func (d Duck) Speak() string { return d.Name + "：嘎嘎！" }

func demoDefine() {
	fmt.Println("========== 01 接口：定义 + 隐式实现 ==========")

	// ---------- ① 多态：同一个接口变量，装不同类型的实现 ----------
	var s Speaker       // 声明一个接口变量（此刻是 nil）
	s = Dog{Name: "旺财"} // 装狗
	fmt.Println("① 接口装 Dog：", s.Speak())
	s = Cat{Name: "咪咪"}    // 同一变量换装猫
	fmt.Println("① 接口装 Cat：", s.Speak())
	s = Duck{Name: "唐老鸭"} // 再换鸭子 —— 契约只要求会 Speak，鸭子会
	fmt.Println("① 接口装 Duck：", s.Speak())

	// ---------- ② 部分满足契约：Duck 只有 Speak，装不进 Animal ----------
	// var a Animal = Duck{Name: "x"} // ❌ 编译报错：Duck 没有 Run()，不满足 Animal
	var a Animal = Dog{Name: "阿黄"} // ✅ Dog 两个方法都有，满足
	fmt.Println("② Dog 满足 Animal：", a.Speak(), "|", a.Run())

	// ---------- ③ 接口变量能调的方法 = 接口声明的方法 ----------
	// s 的静态类型是 Speaker，所以只能调 Speak() —— 想用 Run() 得先断言（见 05）
	// s.Run() // ❌ 编译报错：Speaker 契约里没有 Run
	// 这就是"只认契约"：接口给你什么能力，你就能用什么

	// ---------- ④ 空接口 any：能装任何东西 ----------
	// 所有类型都满足空接口（契约 = 0 个方法）
	// 前端对照：TS 的 unknown（能装万物，但要收窄才能用）
	var anything any
	anything = 42
	anything = Dog{Name: "小黑"}
	fmt.Println("④ any 随便装：", anything)

	fmt.Println("\n🔑 接口 = 契约：隐式实现 + 多态；any = 没有方法的接口，装万物")
}
