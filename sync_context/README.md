# Go 并发控制（sync + select + Context）学习

> 面向「前端（TS/JS）转 Go 后端」的学习代码，注释尽量直白，建议边跑边改。
> ⚠️ 本章是 goroutine 的续篇：Day6 学"怎么开协程、怎么通信"，Day7 学"怎么控制协程"。
> ⚠️ 本章是**前端知识覆盖最少的一章**（JS 单线程，没有锁、并发 map 这些概念），
>    所以每节都按"问题 → 工具 → 语法 → 实验"来写 —— 看不懂就回去重看"问题"部分。

## 怎么运行

```bash
cd sync_context
go run .          # 一次跑完所有演示
```

## 注释标记（Better Comments 插件）

| 标记 | 颜色 | 含义 |
| --- | --- | --- |
| `//!` | 红色 | 必踩的坑 / 必记的重点（面试会问） |
| `//?` | 蓝色 | 思考题：先想答案，再对注释里的答案 |
| `//*` | 绿色 | 关键铺垫 / 类比 |

> 没装插件的话：VSCode 扩展搜 "Better Comments" 安装即可
> 本章凡是第一次出现的**全新语法**（之前没见过），正文都用红色 `//! 新写法：` 标出来，
> 读到就慢一点、停一下，那是本章真正要学的东西。

推荐学习顺序（也是建议的阅读顺序）：

| 文件 | 内容 | 对应学习路线 |
| --- | --- | --- |
| [01_mutex.go](01_mutex.go) | Mutex 修数据竞争、RWMutex 读多写少 | Day7 核心 |
| [02_select.go](02_select.go) | select 多路复用、超时、default | Day7 核心 |
| [03_context.go](03_context.go) | Context 取消/超时/传值、防协程泄漏 | Day7 核心（最抽象，慢读） |
| [04_once_map.go](04_once_map.go) | sync.Once 只跑一次、sync.Map 并发安全 map | Day7 次重点 |

## 一句话总结

> 本章前端能对照的只有两处：`select ≈ Promise.race`、`Context ≈ AbortController`。
> 其余靠 Go 内部的关联来记（都写进了"对照 / 关联"列）。

| 概念 | 一句话 | 对照 / 关联 |
| --- | --- | --- |
| 数据竞争 | 多协程同时读改同一变量，更新互相覆盖 | JS 单线程无此坑（反例） |
| `sync.Mutex` | 给共享内存装"门禁"，同时只放一个人进 | JS 单线程不需要 |
| `Lock` / `Unlock` | 抢锁 / 还锁；忘还 = 死锁 | —— |
| `defer mu.Unlock()` | 函数结束必还锁，防漏 | 呼应 `defer wg.Done()` |
| `sync.RWMutex` | 读读共享、写独占；读多写少更快 | —— |
| `select` | 从多个 channel 挑一个最先就绪的 | `Promise.race` |
| `time.After(d)` | 一根 d 后才来水的 channel，配 select 做超时 | `setTimeout` |
| `default` | 都没就绪立刻返回，不阻塞 | —— |
| `context.Background()` | Context 的根节点，所有 ctx 的起点 | —— |
| `WithCancel` | 手动 cancel() 喊停，下游全停 | `AbortController.abort()` |
| `WithTimeout` | 到点自动取消（超时） | —— |
| `WithValue` | 沿调用链传 requestId 等元信息 | —— |
| `ctx.Done()` | 被取消/超时时关闭的 channel，select 监听它 | —— |
| `sync.Once` | 保证某动作只执行一次（单例初始化） | 模块级单例 |
| `sync.Map` | 并发安全的 map，读多写少时用 | JS 单线程不需要 |

## 面试速背（八股，复习时只看这一节）

- **Mutex 防数据竞争**：`count++` 是读-改-写三步，多协程穿插会丢更新；`go run -race` 检测
- **锁三定律**：Lock 后必 Unlock（defer 兜底）；锁别拷进 struct；锁粒度过大 = 并发变串行
- **RWMutex**：读读不互斥、写独占；读多写少用它，写多时不如普通 Mutex
- **select 语义**：多个 case 就绪随机挑一个；全阻塞且无 default → 死锁
- **select 超时模式**：`case <-time.After(d)`，等不到就放弃，别让协程死等
- **Context 铁律**：生成即 `defer cancel()` 防定时器泄漏；key 用自定义类型；ctx 永远是函数第一个参数，不存进 struct
- **ctx.Done()**：取消或超时那一刻被关闭，协程用 select 监听它优雅退出
- **context 三种**：WithCancel（手动停）、WithTimeout/WithDeadline（到点停）、WithValue（传值）
- **sync.Once**：底层 = Mutex + 标记位；once.Do(f) 保证 f 只执行一次
- **sync.Map 适用**：读多写少、key 基本不变；写多时性能反而差
- **Go 并发格言**（面试必背）：不要通过共享内存来通信，要通过通信来共享内存

## 最容易踩的 3 个坑

1. **忘 Unlock / 忘 cancel()**：锁不还 → 死锁；Context 不 cancel → 定时器泄漏。都用 defer 兜底
2. **select 全阻塞无 default**：所有 case 都没就绪 → 主协程死锁
3. **普通 map 并发写**：直接 `fatal error: concurrent map writes` 崩溃，不是返回错误

## 下一章预告（第二阶段 Week3）

`net/http` 写 HTTP Server → **Gin 框架**（路由、中间件、参数绑定）。
Context 到那时就天天见了：每个 HTTP 请求自带一个 ctx，超时 / 取消自动帮你处理。
