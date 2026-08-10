// 02_slicing.go —— 切片截取：Go 里最经典的一个"坑"，必须搞懂
//
// ⚠️ 重点对比（前端同学最容易踩坑的地方）：
//
//	JS：arr.slice(1, 3)  → 返回一个【新数组】，和原数组完全独立
//	Go：base[1:3]        → 返回一个【视图】，与原切片【共享同一块内存】！
//
// 也就是说：改了视图里的元素，原切片也会变，因为它们是同一片内存。
package main

import "fmt"

// demoSlicing 演示切片截取后共享底层数组
func demoSlicing() {
	fmt.Println("========== 02 切片截取：共享底层数组 ==========")

	var base2 []int
	base2 = []int{1, 2, 3}

	fmt.Println(base2)

	base := []int{1, 2, 3, 4, 5}
	view := base[1:3] // 取下标 1、2 → [2, 3]

	fmt.Println("base =", base) // [1 2 3 4 5]
	fmt.Println("view =", view) // [2 3]

	// 用 %p 打印"第一个元素的地址"—— 两个地址一模一样，证明是同一块内存
	fmt.Printf("base 首元素地址: %p\n", &base[0])
	fmt.Printf("view 首元素地址: %p\n", &view[0])

	//! 修改 view，base 跟着变！
	view[0] = 999
	fmt.Println("\n修改 view[0] = 999 之后：")
	fmt.Println("view =", view) // [999 3]
	fmt.Println("base =", base) // [1 999 3 4 5] ← base 也被改了！

	fmt.Println("\n🔑 记忆法：截取只是「开了一个窗口」，没有复制数据")
	fmt.Println("   如果你想要 JS 那种独立的副本，请用 copy（见 04 节）")

	// ---------- 进阶：三索引截取，限制容量 ----------
	//* 语法：base[low:high:max]，第三个参数 max 用来限制结果的容量
	// 三句话记牢：len = high-low，cap = max-low，约束 0 <= low <= high <= max <= cap(base)
	nums := []int{1, 2, 3, 4, 5}
	normal := nums[1:3] // 不限制容量 → cap 自动 = 底层剩余空间 4
	safe := nums[1:3:3] // 限制容量为 max-low = 2
	fmt.Println("\nnums      =", nums)
	fmt.Println("normal    =", normal, " cap =", cap(normal))
	fmt.Println("safe      =", safe, " cap =", cap(safe))
	fmt.Println("safe 的 cap 只有 2，之后 append 必定扩容 → 不会污染 nums")

	//! 不写 max 的坑：append 没超 cap 就不扩容，直接写进底层数组 → 污染原数据
	nums2 := []int{1, 2, 3, 4, 5}
	view2 := nums2[1:3] // cap = 4，射程覆盖到底层剩余的全部空间
	view2 = append(view2, 999)
	fmt.Println("\n不限制容量 append 之后：")
	fmt.Println("view2 =", view2) // [2 3 999]
	fmt.Println("nums2 =", nums2) // [1 2 3 999 5] ← 被污染了！

	//* 三索引救法：cap 被卡死，append 一超 cap → 强制扩容、开新数组
	nums3 := []int{1, 2, 3, 4, 5}
	safe2 := nums3[1:3:3] // cap = 2，射程只有自己
	safe2 = append(safe2, 999)
	fmt.Println("\n限制容量 append 之后：")
	fmt.Println("safe2 =", safe2) // [2 3 999]
	fmt.Println("nums3 =", nums3) // [1 2 3 4 5] ← 毫发无损

	//? 思考：safe2[0] = 100 会不会影响 nums3[1]？
	//? 答：会！[low, high) 内的元素仍然共享同一块内存（三索引 ≠ 独立副本）
	//? 但注意：一旦 append 触发扩容，safe2 就指向新数组，之后改它不再影响 nums3

	//* 常用姿势：s[i:j:j]（第三参 = high）= "窗口封死，一 append 就翻脸扩容"
	// 用途：把 s[i:j:j] 传给可能 append 的函数，防止它污染你的数据
	//      原地删除惯用法：s = append(s[:i:i], s[i+1:]...)
	//! 注意：三索引不能释放底层数组内存；想要完全独立的副本请用 copy（见 04 节）
}
