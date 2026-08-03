// 05_function.go —— 切片传参：为什么"改了外面变"、"append 了外面不变"？
//
// 这是 90% 初学者都会懵的问题，请务必亲手跑一遍并理解。
//
// 先记住结论：
//   ① Go 里所有参数都是【值传递】（拷贝一份传进去）
//   ② 但一个切片变量其实由三部分组成：{ 指向底层数组的指针, 长度, 容量 }
//   ③ 拷贝切片 = 拷贝这个"小结构体"，而结构体里的【指针指向的底层数组是同一个】！
//
// 由此推出一组现象：
//   改元素 → 改的是共享的底层数组 → 外面能看到
//   append 触发扩容 → 换了底层数组 → 外面的切片还指着旧数组 → 外面看不到
package main

import "fmt"

// demoFunction 演示切片传参的三种情况
func demoFunction() {
	fmt.Println("========== 05 切片传参的真相 ==========")

	nums := []int{1, 2, 3}

	modifyElement(nums) // 函数内改元素
	fmt.Println("情况1 函数内【改元素】  → nums =", nums) // [99 2 3] 变了

	appendInside(nums) // 函数内 append 但不返回
	fmt.Println("情况2 函数内【append】  → nums =", nums) // [99 2 3] 没变！

	nums = appendAndReturn(nums) // 函数内 append 并把返回值接住
	fmt.Println("情况3 函数内【append 且接住返回值】→ nums =", nums) // [99 2 3 4]

	fmt.Println("\n🔑 为什么 append 后外面不变？")
	fmt.Println("   因为外面 nums 的容量只有 3，函数内 append 第 4 个元素会触发扩容，")
	fmt.Println("   底层数组换成了新的：函数内的 s 指向新数组，外面的 nums 还指向旧数组。")
	fmt.Println("   想让外面生效，唯一办法就是接住返回值：nums = append(nums, ...)")
}

// 情况 1：修改元素 → 外面可见
// 原因：s 和外面 nums 的指针指向同一个底层数组，改这里 = 改那里
func modifyElement(s []int) {
	s[0] = 99
}

// 情况 2：append 且不返回 → 外面不可见
// 原因：外面 nums 的 cap = 3，这里加第 4 个触发扩容 → s 指向新数组
//      新数组只有函数内的 s 知道，外面的 nums 仍指向旧数组 → 看不到新元素
func appendInside(s []int) {
	s = append(s, 4)
}

// 情况 3：append 并返回 → 外面接住返回值才有效
func appendAndReturn(s []int) []int {
	return append(s, 4) // 把扩容后的"新切片"还出去
}
