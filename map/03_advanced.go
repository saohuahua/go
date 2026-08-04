// 03_advanced.go —— 进阶：无序遍历 + 排序 + 面试坑 + 底层铺垫
//
// 这节是八股复习的浓缩版，覆盖面试最爱问的 map 问题。
package main

import (
	"fmt"
	"sort"
)

// demoAdvanced 演示 map 遍历的无序性、排序、以及经典面试坑
func demoAdvanced() {
	fmt.Println("========== 03 进阶：遍历 + 面试坑 ==========")


	// ---------- 1. 遍历：顺序是【随机的】！ ----------
	city := map[string]int{"北京": 1, "上海": 2, "广州": 3}
	fmt.Println("遍历顺序每次都可能不同：")
	for k, v := range city {
		fmt.Printf("  %s → %d\n", k, v)
	}

	// 前端对照：JS 对象整数 key 升序 + 其他按插入序；JS Map 严格插入序
	//          Go map 完全无序！想要顺序 → 先取 key 再排序
	keys := make([]string, 0, len(city))
	for k := range city {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 排序后顺序固定
	fmt.Println("排序后再遍历：", keys)


	// ---------- 2. 面试坑：value 不能取地址 ----------
	ages := map[string]int{"张三": 18}

	//! p := &ages["张三"]  // ❌ 编译错误：cannot take the address of map index expression
	//! 原因：map 扩容（rehash）时元素会"搬家"，地址不稳定，Go 干脆不让取地址

	// 想改 value？直接 ages["张三"] = 19 就行，不需要地址
	ages["张三"] = 19
	fmt.Println("\n2. &m[k] 编译报错：map 元素地址不稳定，改值直接 ages[k] = v 即可（现在张三 =", ages["张三"], "）")


	// ---------- 3. 面试坑：什么能当 key ----------
	// key 必须"可比较"（== 能用）：string / int / bool / struct 都可以
	// 切片、map、函数【不能】当 key（编译报错）
	type Point struct {
		X, Y int
	}
	grid := map[Point]string{{1, 2}: "左上"}
	fmt.Println("3. struct 作 key：grid[{1,2}] =", grid[Point{1, 2}]) // 左上


	// ---------- 4. 面试坑：value 是切片，append 必须接住 ----------
	classes := map[string][]string{}
	classes["一班"] = append(classes["一班"], "小明")
	classes["一班"] = append(classes["一班"], "小红")
	fmt.Println("4. value 是切片：", classes) // map[一班:[小明 小红]]

	//! 必须写 classes["一班"] = append(...)，不能只写 append(classes["一班"], ...)
	//! 因为切片 append 扩容会换底层数组，必须把新切片接回去（slices/05）


	// ---------- 5. 底层铺垫（面试一句话版） ----------
	// map 底层 = 哈希表：数组 + 若干"桶"，每桶存约 8 个 key-value
	// 查找：算 key 的哈希 → 定位桶 → 桶内比对 → 平均 O(1)
	// 桶满 → 溢出到溢出桶 → 溢出过多 → 扩容 rehash
	// 面试答法：说出「哈希函数 → 桶 → 冲突 → 扩容」这条链路就够
	fmt.Println("\n5. 底层：哈希表 + 桶，平均 O(1)，冲突过多触发扩容 rehash")
}
