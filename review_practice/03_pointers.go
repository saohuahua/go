// 03_pointers.go —— 指针：修改 struct 本体，以及切片元素的正确取址方式
package main

import "fmt"

func demoPointers() {
	fmt.Println("\n========== 03 指针：原地修改 ==========")
	tasks := sampleTasks()

	check("完成 ID=1", markDone(tasks, 1), true)
	check("任务确实被原地修改", tasks[0].Done, true)

	task := Task{Title: "整理笔记"}
	task.addTag("go")
	check("指针接收者添加 tag", task.Tags, "[go]")
}

// TODO 03.1：找到 id 对应任务，将 Done 改为 true；找到返回 true，找不到返回 false。
// //! range 的 task 是元素副本。需要修改切片元素时，用 for i := range tasks，再取 &tasks[i]。
func markDone(tasks []Task, id int) bool {
	return false
}

// TODO 03.2：为 Task 实现 addTag 方法。
// 要求：使用指针接收者，让调用者的 Task 被修改。
func (task *Task) addTag(tag string) {
	// TODO
}
