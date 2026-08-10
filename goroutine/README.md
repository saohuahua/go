# Go 并发（Goroutine + Channel）学习

> 面向「前端（TS/JS）转 Go 后端」的学习代码，注释尽量直白，建议边跑边改。
> ⚠️ 本章是 Go 与前端差距最大的地方：Promise 是单线程上的"排队"，goroutine 是真的"另开跑道"。
> ⚠️ 并发是"不确定的"——同一段代码多跑几遍，输出顺序都可能不一样，这是正常的。

## 怎么运行

```bash
cd goroutine
go run .          # 一次跑完所有演示
```

## 注释标记（Better Comments 插件）

| 标记 | 颜色 | 含义 |
| --- | --- | --- |
| `//!` | 红色 | 必踩的坑 / 必记的重点（面试会问） |
| `//?` | 蓝色 | 思考题：先想答案，再对注释里的答案 |
| `//*` | 绿色 | 关键铺垫 / 前端对照，秒懂的类比 |

> 没装插件的话：VSCode 扩展搜 "Better Comments" 安装即可

推荐学习顺序（也是建议的阅读顺序）：

| 文件 | 内容 | 对应学习路线 |
| --- | --- | --- |
| [01_basics.go](01_basics.go) | go 启动协程、主函数退出 = 协程全没（最坑） | Day6 基础 |
| [02_channel.go](02_channel.go) | channel 收发、无缓冲 = 同步阻塞（核心） | Day6 核心 |
| [03_buffered.go](03_buffered.go) | 有缓冲 channel、生产者-消费者、close | Day6 实战 |
| [04_waitgroup.go](04_waitgroup.go) | WaitGroup 正规等待、数据竞争预告 | Day6 收尾 |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| `go f()` | 另开一条跑道跑 f，立刻返回、不等结果 | Promise 的异步，但是真并行 |
| 主函数退出 | 所有协程一起被杀，不会等你 | JS 事件循环会等任务排完 |
| `make(chan int)` | 造一根协程间通信的水管 | postMessage / 任务队列 |
| `ch <- v` / `v := <-ch` | 发送 / 接收（方向看箭头） | 发消息 / 收消息 |
| 无缓冲 channel | 收发必须"同时在场"，同步阻塞 | 无对应（JS 全异步） |
| 有缓冲 channel | 塞满才等，像储物柜 | 有限长度的队列 |
| `len/cap(ch)` | 当前排队数 / 总容量 | 同 slice 的 len/cap |
| `close(ch)` | 发送方声明"不会再有了" | —— |
| `range ch` | 取到 channel 关闭自动结束 | —— |
| `v, ok := <-ch` | ok=false = 已关闭取空 | 同 map 的 comma-ok |
| WaitGroup | 计数器等一批协程全部跑完 | `Promise.all` |
| `defer wg.Done()` | 干完必通知，防漏计数 | —— |

## 面试速背（八股，复习时只看这一节）

- **go 启动协程**：`go f()` 立刻返回；主函数一退出，所有协程被杀（不等待）
- **goroutine 不是线程**：轻量（栈可增长、KB 级起），一个进程能开几万个
- **channel 是引用类型**：传参就是同一个管道（呼应 map/02）
- **无缓冲 channel = 同步**：收发必须成对，单独发/收会死锁
- **有缓冲 channel**：`make(chan T, n)`，塞满才阻塞；len=排队数 cap=容量
- **close + range**：只有发送方该 close；`range ch` 收到关闭自动退出
- **channel 的 comma-ok**：`v, ok := <-ch`，ok=false = 已关闭（呼应 map/01）
- **写已关闭的 channel = panic**；读已关闭的 channel = 零值，安全
- **WaitGroup**：Add 在开协程前、Done 用 defer、Wait 阻塞到归零
- **格言**：不要通过共享内存来通信，要通过通信来共享内存
- **GMP 模型**：G=协程 M=线程 P=调度器。现在知道"协程被调度到线程上跑"即可，深挖在路线文档第四阶段面试清单

## 最容易踩的 3 个坑

1. **主函数退出 = 协程全没**：`go f()` 后不等待就 return → f 根本没跑
2. **死锁**：无缓冲 channel 只发不收 / 只收不发；WaitGroup 忘 Done 或 Add 位置写错
3. **数据竞争**：多个协程同时读写同一变量，结果不确定（用 `go run -race` 能检测）

## 下一章预告（Day7）

`select`（多路复用多个 channel）、`sync.Mutex / RWMutex`（锁，解决数据竞争）、`Context`（超时 / 取消传播）。学到那再说。
