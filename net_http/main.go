// main.go —— net/http：理解 HTTP Server 的底层，再进入 Gin
//
// 本文件夹共 5 个演示文件 + 本入口文件
//
// 01_server.go            —— HTTP Server、Handler 接口、HandlerFunc
// 02_request.go           —— method、URL、query、header
// 03_response_json.go     —— 状态码、响应头、JSON、struct tag
// 04_router_middleware.go —— 路由分发、中间件链、请求 Context
// 05_todo_api.go          —— 内存版 TODO API
//
// 运行方法
//
// cd net_http && go run .
// go run ./net_http    （在项目根目录执行）
//
// 学习方法
//
// 1 按编号阅读，前四节分别拆开 HTTP Server 的一个部件
// 2 05 是最小 API，改请求、状态码和数据观察结果
// 3 学 Gin 时将 05 重写一次，对比框架省掉的样板代码
package main

import "fmt"

func main() {
	// demoServer() // 01 HTTP Server + Handler
	// demoRequest() // 02 读取请求
	// demoResponseJSON() // 03 写 JSON 响应
	// demoRouterMiddleware() // 04 路由 + 中间件
	demoTodoAPI() // 05 内存 TODO API

	fmt.Println("\n🎉 net/http 预热完成：现在知道 Gin 帮你封装了什么，可以进 Gin 了")
}
