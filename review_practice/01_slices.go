// 01_slices.go —— 切片：过滤、删除，以及不意外修改调用者的数据
package main

import "fmt"

func demoSlices() {
	fmt.Println("\n========== 01 切片：筛选 / 删除 / copy ==========")
	tasks := sampleTasks()

	pending := pendingTasks(tasks)
	check("未完成任务", titles(pending), "[复习切片 学习 Context]")

	withoutFirst := removeTaskByID(tasks, 1)
	check("删除 ID=1", titles(withoutFirst), "[写 HTTP API 学习 Context]")
	check("原切片没有被改", titles(tasks), "[复习切片 写 HTTP API 学习 Context]")

	check("复制后修改互不影响", independentCopy(tasks), true)
}

// TODO 01.1：返回所有未完成任务，保持原先顺序。
// 限制：返回的新切片不能和 tasks 共用可写元素；后面改返回值的 Done，原 tasks 不应变化。
// 提示：Task 内还有 Tags 切片；本题先只要求 Task 这一层独立，Tags 深拷贝留给自己加练。
func pendingTasks(tasks []Task) []Task {
	return nil
}

// TODO 01.2：删除指定 ID，返回一个新切片。
// //! 不要写 return append(tasks[:i], tasks[i+1:]...)：它可能复用底层数组，改到原 tasks。
// 提示：make 一个 len 为 0、cap 合理的新切片，再 append 需要保留的元素。
func removeTaskByID(tasks []Task, id int) []Task {
	return tasks
}

// TODO 01.3：做一份独立副本，改副本第一个元素的 Title 后，验证原切片不受影响。
func independentCopy(tasks []Task) bool {
	return false
}

func titles(tasks []Task) []string {
	result := make([]string, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, task.Title)
	}
	return result
}
