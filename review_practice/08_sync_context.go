// 08_sync_context.go —— 给共享 map 加锁，并让慢任务可被取消
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type safeCounter struct {
	mu     sync.RWMutex
	counts map[string]int
}

func demoSyncContext() {
	fmt.Println("\n========== 08 Mutex / select / Context ==========")
	counter := safeCounter{counts: map[string]int{}}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.add("go")
		}()
	}
	wg.Wait()
	check("100 个协程累计", counter.get("go"), 100)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	check("Context 超时取消慢任务", waitForTask(ctx, 100*time.Millisecond), context.DeadlineExceeded)
}

// TODO 08.1：用写锁安全地让指定 tag +1。
func (counter *safeCounter) add(tag string) {
	// TODO
}

// TODO 08.2：用读锁安全地读取；不存在返回 0。
func (counter *safeCounter) get(tag string) int {
	return 0
}

// TODO 08.3：select 等待 workDuration；如果 ctx 先取消，返回 ctx.Err()。
// 提示：case <-time.After(workDuration) 时返回 nil。
func waitForTask(ctx context.Context, workDuration time.Duration) error {
	return nil
}
