// 04_interfaces.go —— 接口：业务代码只依赖“能做什么”，不依赖“具体怎么存”
package main

import "fmt"

// TaskRepository 保持小：这里只放当前 service 真正需要的两个能力。
type TaskRepository interface {
	FindByID(id int) (Task, error)
	Save(task Task) error
}

type memoryTaskRepo struct {
	tasks map[int]Task
}

func newMemoryTaskRepo(tasks []Task) *memoryTaskRepo {
	items := make(map[int]Task, len(tasks))
	for _, task := range tasks {
		items[task.ID] = task
	}
	return &memoryTaskRepo{tasks: items}
}

func demoInterfaces() {
	fmt.Println("\n========== 04 接口：依赖注入 ==========")
	repo := newMemoryTaskRepo(sampleTasks())

	check("通过接口完成任务", completeTask(repo, 3), nil)
	task, err := repo.FindByID(3)
	check("保存后的状态", err == nil && task.Done, true)
}

// TODO 04.1：让 *memoryTaskRepo 隐式实现 TaskRepository。
// FindByID：找不到时暂时返回任意非 nil error（第 05 章会统一规范）。
func (repo *memoryTaskRepo) FindByID(id int) (Task, error) {
	return Task{}, fmt.Errorf("TODO: find task")
}

func (repo *memoryTaskRepo) Save(task Task) error {
	return fmt.Errorf("TODO: save task")
}

// TODO 04.2：只通过 TaskRepository 完成任务，不要把参数写成 *memoryTaskRepo。
func completeTask(repo TaskRepository, id int) error {
	return fmt.Errorf("TODO: complete task")
}
