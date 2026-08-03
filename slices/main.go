// main.go —— 切片（Slice）学习程序入口
//
// 本文件夹是一个完整的 Go 学习项目，共 5 个演示文件 + 本入口文件：
//
//   01_basics.go      —— 数组 vs 切片、三种创建方式、len/cap
//   02_slicing.go     —— 切片截取、共享底层数组（最经典的大坑）
//   03_append.go      —— append 追加 + 扩容机制（面试必问！）
//   04_copy_ops.go    —— copy 深拷贝、删除元素、清空切片
//   05_function.go    —— 切片传参的真相（值传递但不完全）
//
// 运行方法（二选一）：
//   cd slices && go run .
//   go run ./slices    （在项目根目录执行）
//
// 学习方法建议：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 每节都配了和 JS/TS 的对照，帮你把前端经验迁移过来
package main

import "fmt"

func main() {
	// 依次演示，建议逐个放开调用，边看输出边理解
	demoBasics()   // 01 数组 vs 切片、创建方式、len/cap
	demoSlicing()  // 02 截取与共享底层数组
	demoAppend()   // 03 append 与扩容机制
	demoCopy()     // 04 copy / 删除 / 清空
	demoFunction() // 05 传参陷阱

	fmt.Println("\n🎉 全部演示完成！建议打开每个文件，动手改一改再跑一遍")
}
