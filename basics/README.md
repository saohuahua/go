# Go 基础语法补课（接口 / error / defer / 闭包）

> 面向「前端（TS/JS）转 Go 后端」的学习代码，注释尽量直白，建议边跑边改。
> ⚠️ 本文件夹是**进 Gin 框架之前的前置补课**：
>   slices/map/pointers/goroutine/sync_context 已经把"核心用法"学完了，
>   但这 4 个基础语法是 Go 里最常用、却一直没系统学过的 —— 而 Gin 天天在用它们。

## 怎么运行

```bash
cd basics
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

| 文件 | 内容 | 对应 Gin / 项目 的用处 |
| --- | --- | --- |
| [01_interface.go](01_interface.go) | 接口：隐式实现、多态、any、类型断言 | http.Handler、分层架构解耦、mock 测试 |
| [02_error.go](02_error.go) | 错误处理：err!=nil、%w 包裹、errors.Is/As | 统一错误处理、GORM 的 ErrRecordNotFound |
| [03_defer.go](03_defer.go) | defer：栈序、参数立即求值、recover | 关连接、事务回滚、Recovery 中间件 |
| [04_closure.go](04_closure.go) | 闭包：函数当值、捕获、中间件模式 | Gin 中间件（认证 / CORS / 限流） |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| 接口 | 能力契约，方法签名对上就自动实现 | TS 要写 `implements`，Go 不用（结构类型） |
| 空接口 `any` | 能装任何类型的类型，取出来用要断言 | TS 的 `unknown`（不是 `any`） |
| 类型断言 | `v, ok := x.(T)` 从接口里掏具体类型 | TS 的 `as` / instanceof |
| type switch | `switch v := x.(type)` 一步分拣类型 | TS 的 typeof + 收窄 |
| error 返回值 | 函数多返回一个 error，调用方自己判断 | `try/catch`（Go 没有异常） |
| `%w` 包裹 | 错误带上上下文，又不丢底层错误 | —— |
| `errors.Is/As` | 穿透错误链判断"是不是这个错" | `instanceof`（但能穿透包裹层） |
| defer | 函数结束时的收尾动作，后注册先执行 | `finally`（栈序 + 参数立即求值不同） |
| recover | 在 defer 里接住 panic，程序不崩 | `try/catch`（只能配 defer 用） |
| 函数当值 | `func(参数) 返回类型` 就是函数签名 | `const f = () => {}` |
| 闭包 | 内层函数记住外层变量（引用语义） | JS 闭包，概念完全一样 |

## 面试速背（八股，复习时只看这一节）

- **接口是隐式实现**：方法签名对上就自动实现，不用写 implements（对比 TS 显式）
- **接口值 = (类型, 值) 对**：装了 nil 指针的接口 ≠ nil 接口，别直接 `s == nil`
- **值接收者实现接口，指针也算**；**指针接收者实现，值不算**（编译报错）
- **空接口 any**：能装任何类型；拿出来用要类型断言（`v, ok := x.(T)`）
- **错误处理三件套**：`if err != nil` 直接判断 / `fmt.Errorf`+`%w` 包裹带上下文 / `errors.Is`·`As` 判断错误链
- **defer 两个坑**：① 后注册先执行（栈）② 参数在 defer 那一行就求值（不是函数结束时）
- **defer + 命名返回值**：return 之后、返回之前还能改返回值
- **recover 只能在 defer 里用**，接住 panic 防程序崩溃（Gin Recovery 原理）
- **函数是一等公民**：能赋值、传参、返回；类型 = `func(参数) 返回`
- **闭包捕获 = 引用**（记住变量本身不是值）；循环变量是经典坑，Go 1.22+ 已修复

## 最容易踩的 4 个坑（每个主题一个）

1. **接口判空**：`var s Speaker = (*Dog)(nil)` 后 `s == nil` 是 false，判断"空不空"要先断言
2. **错误被吞掉**：`if _, err := f(); err != nil { return err }` 忘了 return / 打日志，bug 就丢了
3. **defer 参数立即求值**：`defer fmt.Println(x)` 打印的是 defer 那一刻的 x，不是函数结束时的
4. **循环里闭包捕获**：老版本不拷贝 i，闭包全记住同一个变量，全打印最后一个值

## 下一步

- 🔴 这 4 个是进 Gin 的门票，补完直接上 Week3（net/http → Gin）
- 🟡 struct tag + JSON、字符串/strconv 到 Gin 里**边用边学**，现在不用专门补
- 泛型（1.18+）面试会问一嘴，但写 Gin 用不到，放最后再补
