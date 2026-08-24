// 05_errors.go —— 错误：区分“找不到”和其他失败，并保留错误上下文
package main

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("任务不存在")

type ValidationError struct {
	Field  string
	Reason string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s：%s", err.Field, err.Reason)
}

func demoErrors() {
	fmt.Println("\n========== 05 错误：哨兵 / 包裹 / 自定义类型 ==========")
	repo := newMemoryTaskRepo(sampleTasks())

	_, err := repo.FindByID(99)
	check("FindByID 返回哨兵错误", errors.Is(err, ErrTaskNotFound), true)

	err = completeTask(repo, 99)
	check("%w 后仍能 errors.Is", errors.Is(err, ErrTaskNotFound), true)

	err = validateTask(Task{Title: ""})
	var validationErr *ValidationError
	check("errors.As 取结构化错误", errors.As(err, &validationErr) && validationErr.Field == "title", true)
}

// TODO 05.1：修改 04 的 memoryTaskRepo.FindByID：找不到时 return Task{}, ErrTaskNotFound。
// TODO 05.2：修改 04 的 completeTask：给 FindByID / Save 错误加上下文，并使用 %w。
// 这两题要回到 04_interfaces.go 完成；这里不重复定义，避免同名冲突。

// TODO 05.3：标题去掉首尾空格后为空时，返回 *ValidationError{Field: "title", Reason: "不能为空"}。
// 提示：使用 strings.TrimSpace；正常时返回 nil。
func validateTask(task Task) error {
	return nil
}
