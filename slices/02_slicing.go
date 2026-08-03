// 02_slicing.go —— 切片截取：Go 里最经典的一个"坑"，必须搞懂
//
// ⚠️ 重点对比（前端同学最容易踩坑的地方）：
//   JS：arr.slice(1, 3)  → 返回一个【新数组】，和原数组完全独立
//   Go：base[1:3]        → 返回一个【视图】，与原切片【共享同一块内存】！
//
// 也就是说：改了视图里的元素，原切片也会变，因为它们是同一片内存。
package main

import "fmt"

// demoSlicing 演示切片截取后共享底层数组
func demoSlicing() {
	fmt.Println("========== 02 切片截取：共享底层数组 ==========")

	base := []int{1, 2, 3, 4, 5}
	view := base[1:3] // 取下标 1、2 → [2, 3]

	fmt.Println("base =", base) // [1 2 3 4 5]
	fmt.Println("view =", view) // [2 3]

	// 用 %p 打印"第一个元素的地址"—— 两个地址一模一样，证明是同一块内存
	fmt.Printf("base 首元素地址: %p\n", &base[0])
	fmt.Printf("view 首元素地址: %p\n", &view[0])

	// 修改 view，base 跟着变！
	view[0] = 999
	fmt.Println("\n修改 view[0] = 999 之后：")
	fmt.Println("view =", view) // [999 3]
	fmt.Println("base =", base) // [1 999 3 4 5] ← base 也被改了！

	fmt.Println("\n🔑 记忆法：截取只是「开了一个窗口」，没有复制数据")
	fmt.Println("   如果你想要 JS 那种独立的副本，请用 copy（见 04 节）")

	// ---------- 进阶：三索引截取，限制容量 ----------
	// 语法：base[low:high:max]，第三个参数 max 用来限制结果的容量
	// 好处：容量被卡死，之后 append 一定会扩容，不会覆盖原数组后面的数据
	nums := []int{1, 2, 3, 4, 5}
	normal := nums[1:3]   // 不限制容量 → cap 自动 = 底层剩余空间 4
	safe := nums[1:3:3]   // 限制容量为 max-low = 2
	fmt.Println("\nnums      =", nums)
	fmt.Println("normal    =", normal, " cap =", cap(normal))
	fmt.Println("safe      =", safe, " cap =", cap(safe))
	fmt.Println("safe 的 cap 只有 2，之后 append 必定扩容 → 不会污染 nums")
}
