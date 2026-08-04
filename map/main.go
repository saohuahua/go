// main.go —— map（Map）学习程序入口
//
// 本文件夹是一个 Go 学习项目，共 3 个演示文件 + 本入口文件：
//
//   01_basics.go    —— 上手：初始化 + nil map 大坑 + 增删改查 + comma-ok
//   02_reference.go —— 本质：map 是引用类型（传参 / 嵌套 / 复制陷阱）
//   03_advanced.go  —— 进阶：无序遍历 + 排序 + 面试坑 + 底层铺垫
//
// 运行方法（二选一）：
//   cd map && go run .
//   go run ./map    （在项目根目录执行）
//
// 学习方法建议（和 slices 一样）：
//   1. 先整篇跑一遍，看输出和注释对应着理解
//   2. 把 main 里某个 demoXxx() 注释掉，只留一个，自己改代码做实验
//   3. 每节都配了和 JS/TS 的对照，把熟悉的对象/Map 思维迁移过来
package main

import "fmt"

func main() {
	demoBasics()    // 01 上手：初始化 + 增删改查
	demoReference() // 02 本质：引用类型
	demoAdvanced()  // 03 进阶：遍历 + 面试坑

	fmt.Println("\n🎉 全部演示完成！建议打开每个文件，动手改一改再跑一遍")
}
