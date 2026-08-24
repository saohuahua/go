// 02_maps.go —— map：把线性列表加工成便于查询和统计的数据
package main

import "fmt"

func demoMaps() {
	fmt.Println("\n========== 02 map：分组 / 计数 ==========")
	tasks := sampleTasks()

	byTag := groupTitlesByTag(tasks)
	check("go 标签任务", byTag["go"], "[复习切片 写 HTTP API 学习 Context]")
	check("http 标签任务", byTag["http"], "[写 HTTP API]")
	check("标签数量", countTags(tasks), "map[go:3 http:1 基础:1 并发:1]")
}

// TODO 02.1：按标签分组，key 是 tag，value 是拥有该 tag 的任务标题。
// //! map value 是 slice 时，必须写 groups[tag] = append(groups[tag], title)，不能漏掉赋值。
func groupTitlesByTag(tasks []Task) map[string][]string {
	return nil
}

// TODO 02.2：统计每个标签出现次数。
// 提示：map 取不存在的 key 得到 value 类型的零值；int 的零值是 0。
func countTags(tasks []Task) map[string]int {
	return nil
}
