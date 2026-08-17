// 04_once_map.go —— 次重点：sync.Once（只跑一次）+ sync.Map（并发安全的 map）
//
// 这一节的目标是"知道有它们、会用、说清适用场景"，不用深挖原理。
// 前端同样没有可对照的知识（JS 单线程没有"并发"这回事，更不需要并发 map），
// 所以按"问题 → 工具"两段式来读。
package main

import (
	"fmt"
	"sync"
)

// demoOnceMap 演示 sync.Once 只执行一次 + sync.Map 并发安全
func demoOnceMap() {
	fmt.Println("========== 04 Once + sync.Map ==========")

	// ---------- 问题①：初始化只该成功一次 ----------
	// 经典场景：全局单例的"懒加载"。多个协程同时进来抢着初始化，
	// 如果各自初始化一遍，就会开出重复的资源（连接、实例）。
	//! 新写法：sync.Once / once.Do(f) —— 保证 f 在整个程序里只跑一次
	var once sync.Once
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			once.Do(func() { // 多个协程抢着调，但只有第一个真正执行
				fmt.Println("  第", n, "个协程进来，但初始化只执行这一次")
			})
		}(i)
	}
	wg.Wait()
	// 底层实现 = 互斥锁 + 一个标记位（面试问到说这句就够了）

	// ---------- 问题②：普通 map 并发写会直接崩溃 ----------
	// 多个协程同时写同一个 map → 运行时 fatal error，连报错的机会都不给
	// 这是 map/02 没提的"并发陷阱"，今天补上。解封下面代码试一次：
	fmt.Println("\n① 普通 map 并发写会崩溃，解封下面代码试一次：")
	// m := make(map[string]int)
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		m["key"]++ // fatal error: concurrent map writes
	// 	}()
	// }
	// wg.Wait()

	// ---------- 解法：sync.Map —— 并发安全的 map ----------
	//! 新写法：sync.Map / .Store() / .Load() / .LoadOrStore() / .Range()
	// 适用场景：读多写少、key 一旦写入基本不变（配置表、缓存）
	// 注意：它不是万能 map —— 写入频繁时性能反而打不过"Mutex + 普通 map"
	var sm sync.Map
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sm.Store("key", n) // 并发写，不崩
		}(i)
	}
	wg.Wait()
	v, _ := sm.Load("key") // Load 双返回值 v, ok —— 呼应 map 的 comma-ok
	fmt.Println("② sync.Map 并发写了 10 次，最后 key =", v)

	// 其它常用方法：
	//   LoadOrStore(k, v)  有就返回旧的，没有才存新的 —— 缓存"只打一次"利器
	//   Delete(k)          删除
	//   Range(fn)          遍历（顺序不保证，呼应 map/03 的随机顺序）

	//? 思考：什么时候用 sync.Map，什么时候用 Mutex + 普通 map？
	//? 答：读多写少、key 稳定（缓存类）→ sync.Map 省心；写多 / 逻辑复杂
	//?     → Mutex + 普通 map 更灵活也更快

	fmt.Println("\n🔑 Once = 只成功一次；sync.Map = 并发安全的 map，但不是万能钥匙")
}
