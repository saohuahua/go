// 01_interface.go —— 接口：Go 的"能力契约"（前端最难迁移的一节）
//
// 前端对照一句话（先记这个）：
//
//	TS：interface + 类写 implements X，显式声明"我实现了你"
//	Go：只要方法签名对上，自动就算实现（结构类型 / duck typing）
//	也就是说：TS 要自己写 implements 才认，Go 不用写，形状对上就认
//
// 为什么重要：Gin 和标准库全靠接口解耦 —— 调用方只依赖"能力"，不依赖具体类型。
//
//	能看懂接口，才能看懂 http.Handler、Gin 中间件签名、写 mock 测试。
//
// 本节一条主线：定义契约 → 两个实现 → 多态 → 断言掏出来 → any 装万物 → type switch 分拣
package main

import "fmt"

// ---------- 1. 定义契约（接口）：只有方法签名，没有实现 ----------

// Speaker 接口：谁有 Speak() string 方法，谁就是 Speaker
type Speaker interface {
	Speak() string
}

// Runner 接口：谁有 Run() string 方法，谁就是 Runner
type Runner interface {
	Run() string
}

// Animal 接口：接口也能嵌套组合（把几个契约拼成一个更大的契约）
type Animal interface {
	Speaker
	Runner
}

// ---------- 2. 两个"实现者"：注意！没写任何 implements 关键字 ----------

// Dog 有 Speak + Run 两个方法 → 自动实现 Speaker / Runner / Animal
type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + "：汪汪！" }
func (d Dog) Run() string   { return d.Name + " 撒腿就跑" }

// Cat 也一样，方法签名对上就自动实现
type Cat struct{ Name string }

func (c Cat) Speak() string { return c.Name + "：喵呜~" }
func (c Cat) Run() string   { return c.Name + " 轻巧地走" }

type Duck struct{ Name string }

func (d Duck) Speak() string { return d.Name + ":嘎嘎嘎" }

func demoInterface() {
	fmt.Println("========== 01 接口：能力契约 + 隐式实现 ==========")

	// ---------- ① 多态：同一个接口变量，装不同类型 ----------
	var s Speaker       // 声明一个接口变量（此刻是 nil）
	s = Dog{Name: "旺财"} // 装进一只狗
	fmt.Println("① 接口装 Dog：", s.Speak())

	s = Cat{Name: "咪咪"} // 同一变量换装一只猫
	fmt.Println("① 接口装 Cat：", s.Speak())

	// ---------- ② 面向接口编程：函数只认契约，不认具体类型 ----------
	// 函数接收 Speaker，狗能传、猫能传 —— 以后加个 Duck 也不用改这个函数
	printSpeak(Dog{Name: "大黄"})
	printSpeak(Cat{Name: "小白"})
	printSpeak(Duck{Name: "唐老鸭"})

	// ---------- ③ 类型断言：把具体类型从接口里"掏"出来 ----------
	// 场景：接口方法不够用，想用具体类型专属的字段/方法
	var a Animal = Dog{Name: "阿黄"} // Animal 是更大的契约，Dog 能满足
	if dog, ok := a.(Dog); ok {    // comma-ok 断言：成功 ok=true，dog 是 Dog 类型
		fmt.Println("③ 断言成功：", dog.Name, "是只狗，能看它的专属字段")
	}
	if cat, ok := a.(Cat); ok {
		fmt.Println("③ 不会走到：", cat.Name)
	} else {
		fmt.Println("③ 断言失败走 else：a 里面装的不是猫")
	}

	// ---------- ④ 空接口 any：能装任何东西 ----------
	// Gin 的 c.JSON(200, gin.H{"msg": "ok"}) 里 gin.H 本质是 map[string]interface{}
	var anything any
	anything = 42
	anything = "换个字符串"
	anything = Dog{Name: "小黑"}
	fmt.Println("④ any 随便装：", anything)

	// ---------- ⑤ type switch：一步到位分拣 any 里的类型 ----------
	// 前端对照：TS 里 typeof + 收窄（narrowing），Go 用 type switch 干同一件事
	classify(42)
	classify("hello")
	classify(3.14)

	//! 面试高频坑①：装了"nil 指针"的接口 ≠ nil 接口
	//!   接口底层是两个字段：(类型, 值)。p 是 nil，但装进接口后类型是 *Dog，
	//!   所以 s2 == nil 是 false —— 想判"里面是不是空"得先断言，别直接 s2 == nil
	var p *Dog         // nil 指针
	var s2 Speaker = p // 装进接口：类型是 *Dog，值是 nil
	fmt.Println("⑥ 接口装 nil 指针：s2 == nil 吗？", s2 == nil, "（false！类型是 *Dog 不是空接口）")

	//! 面试高频坑②：值接收者 vs 指针接收者实现接口
	//!   值接收者 → 值和指针都算实现（Go 自动 &d 调方法）
	//!   指针接收者 → 只有指针算实现，值不算（编译直接报错）
	var s3 Speaker = &Dog{Name: "指针狗"} // 能编译：值接收者方法，指针也能调
	fmt.Println("⑦ 指针也满足接口：", s3.Speak())

	fmt.Println("\n🔑 一句话：接口 = 能力契约，隐式实现；用 any 装万物，用断言掏出来")
}

// printSpeak 只认 Speaker 契约，不关心具体是谁（新增类型都不用改这个函数）
func printSpeak(s Speaker) {
	fmt.Println("②", s.Speak())
}

// classify 用 type switch 分拣 any 里的类型
func classify(x any) {
	switch v := x.(type) { // 语法：switch 变量.(type)，每个 case 写一个类型
	case int:
		fmt.Println("    type switch：", v, "是 int")
	case string:
		fmt.Println("    type switch：", v, "是 string")
	case float64:
		fmt.Println("    type switch：", v, "是 float64")
	default:
		fmt.Println("    type switch：不认识这个类型")
	}
}
