# Gin：Web 框架速通

> 面向「前端（TS/JS）转 Go 后端」的学习代码
> 前置：先完成 [net_http/](../net_http/)——Gin 帮你封装的正是那一章手写的东西

## 怎么运行

```bash
cd gin
go run .
```

依赖（gin、jwt、validator）已写入根目录 go.mod，无需手动安装。

## 注释标记

| 标记 | 含义 |
| --- | --- |
| `//!` | 必踩的坑 / 必记重点 |
| `//?` | 思考题 |
| `//*` | 关键铺垫 / 前端对照 |

> Better Comments 标记只出现在函数体内

## 学习顺序

| 文件 | 内容 | net_http 对照 |
| --- | --- | --- |
| [01_router.go](01_router.go) | 路由注册、`:id`/`*rest`、分组、NoRoute | 手写 switch 分流 |
| [02_binding.go](02_binding.go) | Query/Uri/JSON 三种绑定、tag 映射 | `Query().Get`、`json.Decoder` |
| [03_middleware.go](03_middleware.go) | 中间件三种挂法、Next/Abort、Set/Get | 手写洋葱链 |
| [04_validate.go](04_validate.go) | binding tag、错误翻译、统一响应 | 手写 if 校验 |
| [05_jwt.go](05_jwt.go) | JWT 签发/校验、鉴权中间件 | 新知识 |
| [06_todo_api.go](06_todo_api.go) | 用 Gin 重写内存 TODO API | net_http/05 全文 |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| `gin.Engine` | 路由表 + 中间件链，本身实现 `http.Handler` | Express 的 app |
| `gin.Context` | 一次请求的全部上下文，请求结束即失效 | (req, res) 合体 |
| `gin.H` | `map[string]any`，写小 JSON 的速记 | 对象字面量 |
| `:id` / `*rest` | 匹配一段 / 匹配剩余全部 | 动态路由 / 通配路由 |
| `Group` | 公共前缀 + 组级中间件 | 嵌套路由 / 模块化路由 |
| `ShouldBind` 系列 | 按来源绑定 + 自动校验 | req.query / req.body |
| `c.Next()` | 放行；前后分别是洋葱的进和出 | 拦截器里 `await next()` |
| `c.Abort()` | 打「别再往下走」的标记，不跳出函数 | 守卫的 `next(false)` |
| `c.Set`/`c.Get` | 请求内传值，并发请求互相隔离 | 挂在请求上的 meta |
| `binding` tag | 校验规则写进 struct tag，绑定即校验 | zod / el-form rules |
| JWT | 带签名的 JSON，验签即验身份 | localStorage + 拦截器背后的事 |

## Gin vs 手写 net/http（本章核心价值）

> 06 就是这张表的完整对照，同一 API 两个框架各写一遍

| 手写 net/http | Gin | 省掉了什么 |
| --- | --- | --- |
| 手写 switch 按 method+path 分流 | `router.GET/POST/DELETE` | 整个路由分发 |
| `r.URL.Query().Get` + 手动转型 | `ShouldBindQuery` + form tag | 手动 strconv |
| `json.NewDecoder(r.Body).Decode` | `ShouldBindJSON` | 解析 + 校验二合一 |
| `writeJSON` 手设 header/状态码/body | `c.JSON` | 三步并一步 |
| 手写 `middleware(next)` 再倒序组装 | `router.Use` | 洋葱链组装 |
| `context.WithValue` + `WithContext` + 断言 | `c.Set` / `c.GetInt` | 三步并一步 |
| 手写 404/405 分支 | 路由表自动处理 | 兜底逻辑 |

## 中间件执行顺序（对照 net_http/04 的洋葱）

```
router.Use(timer())          ← 先 Use 的在最外层
router.GET("/x", gate(), h)  ← 单路由中间件夹在中间

请求进来：
timer 进 ─→ gate 进 ─→ h 执行并写响应 ─→ gate 出 ─→ timer 出
（c.Next() 之前是「进」，之后是「出」——响应已在 Next 里写出）
```

- `Abort` 后必须 `return`：Abort 只拦「后续中间件和 handler」，不跳出当前函数
- `c.Next()` 之后不要再改状态码：响应已经发出去了

## 面试速背

- **gin.New vs gin.Default**：Default = New + Logger（请求日志）+ Recovery（panic 兜底 500）
- **Gin 跑在 net/http 上**：gin.Engine 实现了 `http.Handler`，底层就是标准库 Server
- **三种绑定**：`ShouldBindQuery`（form tag）/ `ShouldBindUri`（uri tag）/ `ShouldBindJSON`（json tag）
- **ShouldBind vs Bind**：Bind 失败自动写 400；ShouldBind 只返回 error，响应自己控制（生产推荐）
- **body 只能读一次**：同一请求第二次 `ShouldBindJSON` 必报 EOF；要多次绑定用 `ShouldBindBodyWith`
- **中间件语义**：`c.Next()` 前的代码在 handler 前执行、后的在 handler 后执行；`Abort` 阻断后续链
- **Context 生命周期**：`*gin.Context` 只活在本次请求内，传给 goroutine 要用 `c.Copy()`
- **required 零值陷阱**：int 传 0、string 传 "" 都会被 required 拒绝；要区分「没传」和「零值」用指针字段
- **min/max 双义**：string 上管长度，数字上管大小
- **校验错误翻译**：`errors.As(err, &validator.ValidationErrors)` 断言后逐条翻译成中文
- **JWT 结构**：header.payload.signature；payload 放 claims，签名防篡改；服务端不存会话 = 无状态
- **解析时校验 alg**：防止算法混淆攻击（攻击者改 header 的 alg 骗过不检查的解析器）
- **401 vs 403**：401 没登录/token 无效；403 登录了但没权限
- **双 token**：access（短期干活）+ refresh（长期换新）→ 无感刷新

## 最容易踩的坑

1. **同一请求绑定两次 body**：第二次必然 EOF，灵异现象第一名
2. **required 把 0 当成没传**：age 传 0 报「不能为空」；允许零值就别用 required
3. **Abort 后忘写 return**：后续中间件确实没跑，但本函数剩下的代码照跑
4. **c.Next() 之后改响应**：响应已写出，改状态码无效
5. **把 *gin.Context 存进 struct / 直接传给 goroutine**：请求结束 Context 失效；要传就 `c.Copy()` 或只传需要的值
6. **路由参数冲突**：`/users/:id` 和 `/users/:name` 同位置不同名，启动直接 panic
7. **JWT 密钥硬编码 / 不设过期时间**：泄露密钥 = 任何人可伪造；不设 exp = token 永久有效

## 下一步

1. **搭分层项目骨架**（Week3 Day6-7）：按 `cmd / internal(handler,service,repository) / pkg` 分层，把 06 的 TODO API 拆进分层结构
2. **Swagger**：装 swag CLI，写注释生成 API 文档（做项目时一起）
3. 之后进 Week4：**GORM + MySQL**（数据持久化）

> 学习方法：每个 `//!` 的坑都亲手解开注释跑一次，看一眼真实报错，再注释回去
