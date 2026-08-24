// main.go —— Go 综合回顾练习入口
//
// 本目录按从数据结构到 HTTP 的顺序，完成一个内存版 TODO 服务。
//
//   00_warmup.go            —— struct / 流程控制
//   01_slices.go            —— 切片操作
//   02_maps.go              —— map 分组和统计
//   03_pointers.go          —— 指针和原地修改
//   04_interfaces.go        —— 接口和依赖注入
//   05_errors.go            —— 错误链和自定义错误
//   06_defer_closure.go     —— defer / recover / 闭包
//   07_goroutine_channel.go —— goroutine / channel / WaitGroup
//   08_sync_context.go      —— 锁 / select / Context
//   09_http.go              —— net/http 综合实战
//
// 运行方法
//   go run ./review_practice
//   go run ./review_practice -chapter=01
//   go run -race ./review_practice -chapter=08
package main

import (
	"flag"
	"fmt"
)

var chapter = flag.String("chapter", "all", "运行章节：00 ~ 09，默认 all")

func main() {
	flag.Parse()

	demos := []struct {
		id   string
		run  func()
	}{
		{"00", demoWarmup},
		{"01", demoSlices},
		{"02", demoMaps},
		{"03", demoPointers},
		{"04", demoInterfaces},
		{"05", demoErrors},
		{"06", demoDeferClosure},
		{"07", demoGoroutineChannel},
		{"08", demoSyncContext},
		{"09", demoHTTP},
	}

	matched := false
	for _, demo := range demos {
		if *chapter == "all" || *chapter == demo.id {
			matched = true
			demo.run()
		}
	}
	if !matched {
		fmt.Printf("未知章节 %q；可选值：00 ~ 09 或 all\n", *chapter)
	}
}

func check(name string, got, want any) {
	if fmt.Sprint(got) == fmt.Sprint(want) {
		fmt.Printf("  ✅ %s：%v\n", name, got)
		return
	}
	fmt.Printf("  ❌ %s：得到 %v，期待 %v\n", name, got, want)
}
