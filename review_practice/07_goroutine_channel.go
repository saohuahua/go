// 07_goroutine_channel.go —— 并发任务的结果收集：不靠 Sleep 猜时间
package main

import (
	"fmt"
)

func demoGoroutineChannel() {
	fmt.Println("\n========== 07 goroutine + channel + WaitGroup ==========")
	check("并发统计的总优先级", sumPriorities(sampleTasks()), 6)
}

// TODO 07.1：每个 task 开一个 goroutine，把 task.Priority 发送到 results。
// TODO 07.2：用 sync.WaitGroup 等所有发送者结束，再 close(results)。
// TODO 07.3：range results 汇总并返回。
// //! 只有发送方全部结束才能关闭 channel；不要在多个 goroutine 里各自 close。
// //* WaitGroup ≈ Promise.all；channel 用来把每个异步结果交回主流程。
func sumPriorities(tasks []Task) int {
	return 0
}
