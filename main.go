// main.go —— Go 语言入门第一课
// Go 程序从 main 包的 main 函数开始执行
package main

import "fmt"

func main() {
	// fmt.Println 会在控制台打印一行文字
	fmt.Println("你好，Go！")

	// 1. 变量声明：var 变量名 类型
	var name string = "小明"
	var age int = 18

	// 2. 简短声明：类型自动推断（只能在函数内使用）
	score := 99.5

	// 3. 使用 fmt.Printf 做格式化输出
	fmt.Printf("我叫 %s，今年 %d 岁，考试得了 %.1f 分\n", name, age, score)

	// 4. 条件判断
	if score >= 90 {
		fmt.Println("成绩优秀！")
	} else if score >= 60 {
		fmt.Println("及格了")
	} else {
		fmt.Println("要加油哦")
	}

	// 5. 循环：Go 只有 for，没有 while / do-while
	for i := 1; i <= 3; i++ {
		fmt.Printf("第 %d 次循环\n", i)
	}

	// 6. 调用自己写的函数
	greet("Alice")
	result := add(10, 20)
	fmt.Println("10 + 20 =", result)
}

// greet 是一个自定义函数：func 函数名(参数 类型) 返回类型
func greet(person string) {
	fmt.Printf("Hello, %s!\n", person)
}

// add 接收两个整数，返回它们的和
func add(a int, b int) int {
	return a + b
}
