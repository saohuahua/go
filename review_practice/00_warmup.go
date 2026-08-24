// 00_warmup.go —— 热身：先把之后每一章都会用到的 Task 准备好
package main

import "fmt"

type Task struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Done     bool     `json:"done"`
	Tags     []string `json:"tags"`
	Priority int      `json:"priority"`
}

func sampleTasks() []Task {
	return []Task{
		{ID: 1, Title: "复习切片", Tags: []string{"go", "基础"}, Priority: 2},
		{ID: 2, Title: "写 HTTP API", Done: true, Tags: []string{"go", "http"}, Priority: 1},
		{ID: 3, Title: "学习 Context", Tags: []string{"go", "并发"}, Priority: 3},
	}
}

func demoWarmup() {
	fmt.Println("\n========== 00 热身：struct + for + 条件 ==========")
	tasks := sampleTasks()

	check("任务总数", taskCount(tasks), 3)
	check("已完成数", doneCount(tasks), 1)
	check("最高优先级标题", highestPriorityTitle(tasks), "学习 Context")
}

// TODO 00.1：返回任务数量。
// 限制：不要写死 3；任务为空时应该返回 0。
func taskCount(tasks []Task) int {
	return 0
}

// TODO 00.2：遍历 tasks，数出 Done 为 true 的任务。
// 前端对照：等价于 tasks.filter(task => task.done).length。
func doneCount(tasks []Task) int {
	return 0
}

// TODO 00.3：返回 Priority 最大的 Title。
// 限制：空切片返回空字符串；不要排序，单次遍历完成。
func highestPriorityTitle(tasks []Task) string {
	return ""
}
