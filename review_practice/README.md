# Go 综合回顾练习：内存 TODO 服务

> 这不是再看一遍示例，而是从零补全一个小型「内存 TODO 服务」。
> 每一题都复用同一个 `Task`，前一章的产物会成为后一章的材料：数据处理 → 抽象 → 错误处理 → 并发控制 → HTTP API。

## 怎么练

```bash
# 在仓库根目录
# 只跑当前一章的检查（推荐按顺序）
go run ./review_practice -chapter=00
go run ./review_practice -chapter=01

# 看全部题目目前通过了多少项
go run ./review_practice

# 最后一章完成后，检查并发问题
go run -race ./review_practice -chapter=08
```

每个 `TODO` 都保留了一个**能编译但结果不正确**的占位实现；运行对应章节会显示 `❌`。完成后把占位实现替换掉，直到本章全是 `✅`。

> 建议：一次只打开一个编号文件，先根据「目标 / 限制 / 提示」自己写；卡住超过 10 分钟，再回看对应旧章节。不要一上来搜索答案。

## 学习顺序

| 章节 | 文件 | 练什么 | 回看旧目录 |
| --- | --- | --- | --- |
| 00 | [00_warmup.go](00_warmup.go) | struct、for、条件判断、零值 | pointers/03_struct.go |
| 01 | [01_slices.go](01_slices.go) | 遍历、过滤、删除、不改原切片 | slices/01、04、05 |
| 02 | [02_maps.go](02_maps.go) | 分组、计数、map value 是 slice | map/01、03 |
| 03 | [03_pointers.go](03_pointers.go) | `&tasks[i]`、指针接收者、原地修改 | pointers/02、03、04 |
| 04 | [04_interfaces.go](04_interfaces.go) | 小接口、隐式实现、依赖注入 | basics/interface/01、02、03 |
| 05 | [05_errors.go](05_errors.go) | 哨兵错误、`%w`、`Is`、自定义错误 | basics/error/02、03、04 |
| 06 | [06_defer_closure.go](06_defer_closure.go) | defer 栈序、recover、闭包工厂 | basics/03_defer.go、04_closure.go |
| 07 | [07_goroutine_channel.go](07_goroutine_channel.go) | goroutine、channel、关闭 channel、WaitGroup | goroutine/01～04 |
| 08 | [08_sync_context.go](08_sync_context.go) | Mutex、RWMutex、select、Context 取消 | sync_context/01～04 |
| 09 | [09_http.go](09_http.go) | JSON、路由、中间件、HTTP 状态码 | net_http/01～04 |

## 每章完成定义

- **00～08**：`go run ./review_practice -chapter=编号` 全部显示 `✅`。
- **08**：额外执行 `go run -race ./review_practice -chapter=08`，不能报 data race。
- **09**：补全 Handler 后执行 `go run ./review_practice -chapter=09`；再取消 `demoHTTP` 内最后几行注释，手动用浏览器 / curl 试请求。

## 最终项目会长成什么样

```text
HTTP 请求
  ↓
request ID / logger 中间件
  ↓
Task API Handler
  ↓
TaskRepository（接口，只依赖能力）
  ↓
MemoryTaskRepo（map + RWMutex）
```

这正是之后 Gin + Service + Repository 的缩小版：Gin 只是代替你处理路由、绑定 JSON 和中间件样板代码，接口、error、Context、并发安全这些核心思路不变。

## 练完后的速背

- 切片截取和 map 都不是 JS 的普通对象：切片可能共享底层数组，map 遍历无序且并发读写不安全。
- 需要修改 struct 本体时传指针；修改切片里的元素时拿 `&tasks[i]`，不能取 `&m[key]`。
- 接口保持小，业务依赖接口而不是具体实现；错误用 `%w` 保留链，调用方用 `errors.Is/As` 判断。
- `defer` 做收尾，后注册先执行；goroutine 必须有结果收集或退出机制，不能靠 `Sleep` 碰运气。
- 并发共享数据加锁；取消和超时沿 `context.Context` 传递；HTTP 的输入错误返回 400，资源不存在返回 404，创建成功返回 201。
