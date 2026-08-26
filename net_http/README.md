# net/http：HTTP Server 预热

> 面向「前端（TS/JS）转 Go 后端」的学习代码
> 本章不追求手写 Web 框架，只建立 Gin 必需的底层心智模型

## 怎么运行

```bash
cd net_http
go run .
```

## 注释标记

| 标记 | 含义 |
| --- | --- |
| `//!` | 必踩的坑 / 必记重点 |
| `//?` | 思考题 |
| `//*` | 关键铺垫 / 前端对照 |

> Better Comments 标记只出现在函数体内

## 学习顺序

| 文件 | 内容 | Gin 对照 |
| --- | --- | --- |
| [01_server.go](01_server.go) | `http.Handler`、`HandlerFunc`、Server | `gin.Engine`、`gin.HandlerFunc` |
| [02_request.go](02_request.go) | method、path、query、header | `c.Query`、`c.GetHeader` |
| [03_response_json.go](03_response_json.go) | 状态码、JSON、struct tag | `c.JSON`、`c.ShouldBindJSON` |
| [04_router_middleware.go](04_router_middleware.go) | 路由分发、中间件、请求 Context | `router.GET`、`router.Use` |
| [05_todo_api.go](05_todo_api.go) | 内存 TODO API、读写锁、错误响应 | Gin TODO API |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| `http.Handler` | 处理请求的契约，只有 `ServeHTTP` | 路由回调的统一签名 |
| `HandlerFunc` | 让普通函数也满足 `Handler` | 回调函数 |
| `*http.Request` | 请求的 method、URL、header、body 都在这里 | Request / Axios 配置 |
| `http.ResponseWriter` | 设置 header、status，写回响应 body | 服务端版 `res` |
| `json` tag | 控制 Go 字段序列化后的 JSON 名 | TS DTO 字段映射 |
| 路由 | method + path 分发到 Handler | Vue Router 的匹配思维 |
| 中间件 | 包住 Handler，统一处理日志、鉴权等 | Axios 拦截器 / 路由守卫 |
| `r.Context()` | 请求取消和请求元信息沿调用链传递 | `AbortController` |
| `RWMutex` | 保护并发读写的内存 TODO 数据 | JS 单线程通常不需要 |

## 中间件专题：洋葱模型全解析

> 对应 [04_router_middleware.go](04_router_middleware.go)，把「组装 → 请求穿 → 取回数据」整条链路说透

### 1. 中间件是什么

```go
type middleware func(http.Handler) http.Handler
```

吃进一个 `http.Handler`，返回一个「包好一层的新 Handler」。它不改原 handler，只返回新的——所以调用时必须接住返回值。

前端对照：装饰器 / Axios 拦截器 / Express 的 `app.use`

### 2. 组装：applyMiddleware 为什么倒着包

```go
handler := applyMiddleware(mux, requestIDMiddleware, loggerMiddleware)
// 目标结构 = requestID( logger( mux ) )，requestID 最外层、最先进入
```

参数第 1 个 mux 是最里层，中间件按「想让它多外层」排序。要套成 `requestID( logger( mux ) )`，只能从最里层开始一个个往外套——先套 logger，再套 requestID，这就是 `for i := len(middlewares)-1; i >= 0; i--` **倒着循环**的原因。

```
i=1: handler = logger(mux)                → logger 包住 mux
i=0: handler = requestID( logger(mux) )   → requestID 包住整体
```

`handler = middlewares[i](handler)` 这个赋值必须写：中间结果不存回去，包好的那层就丢了（白包）。

### 3. 请求完整路径（核心）

`handler.ServeHTTP(rec, req)` 触发后（真实服务器最后也是这一句），请求从最外层穿到最里层再穿出来：

```
handler.ServeHTTP(rec, req)
│
├─ requestIDMiddleware  外层：ctx 挂 id → r.WithContext 换新请求 → next
│     └─ loggerMiddleware  中层：记开始时间 → next
│           └─ mux 路由器：匹配 "GET /hello" → 调业务 handler
│                 └─ 业务 handler：从 ctx 取回 id → 写响应
│           └─ logger 退出：打印耗时
└─ requestID 退出：无收尾代码
```

三个关键点：
- **w 一路原样传**：`http.ResponseWriter` 从最外层传到最里层，中间件一般不动它
- **r 换过一次**：`r.WithContext(ctx)` 返回一份拷贝（ctx 换新的、其余照旧），只有下游能读到新 ctx；原 r 的 ctx 不变
- **洋葱进出顺序**：先进先出、后进后出，每层的「进」和「出」都夹着 `next.ServeHTTP(w, r)`

### 4. 请求元信息靠 Context 传，不能存全局

```go
ctx := context.WithValue(r.Context(), requestIDKey{}, "req-demo-001")  // ① 挂上去
next.ServeHTTP(w, r.WithContext(ctx))                                 // ② 换新请求传下去
requestID, _ := r.Context().Value(requestIDKey{}).(string)            // ③ 最里层取回来
```

- `WithValue` **不改原 ctx**，返回派生新 ctx，① ② 两步都必须用返回值
- `.(string)` 是类型断言：`Value()` 返回 `interface{}`（空接口），断言回字符串；正式代码应写成 `id, ok := ...` 再 `if !ok` 兜底
- 为什么不存全局变量/struct：HTTP 请求并发执行，共享一个全局会串数据（A 请求写进去、B 请求读到）；ctx 跟着这次请求走，天然隔离

### 5. 路由 vs 中间件：分工

| | 路由（mux） | 中间件 |
| --- | --- | --- |
| 管什么 | 按 method + path 分发到业务 handler | 所有请求的通用逻辑（日志、鉴权、挂 id） |
| 会不会拒请求 | 会：路径存在但方法不对 → 自动 **405** + `Allow` 头；路径不存在 → 404 | 不会，只包不拦（拦不拦由它自己决定） |
| 只对某些方法生效？ | 是，`"GET /hello"` 只匹配 GET | 否，对所有方法都穿一遍 |
| Gin 对应 | `router.GET(...)` | `router.Use(...)` |

## 面试速背

- **`http.Handler`**：只要实现 `ServeHTTP(ResponseWriter, *Request)` 就能处理 HTTP 请求
- **`http.HandlerFunc`**：函数适配器，让普通函数也实现 `Handler`
- **请求读取**：`r.Method` 取方法，`r.URL.Path` 取路径，`r.URL.Query().Get` 取 query，`r.Header.Get` 取请求头
- **响应顺序**：先 `Header().Set`，再 `WriteHeader`，最后写 body
- **默认状态码**：不显式 `WriteHeader` 时，第一次 `Write` 默认返回 200
- **请求 body**：Server 会在请求结束后自动关闭 `r.Body`，Handler 不需要 `defer r.Body.Close()`
- **JSON 输入错误**：`Decoder.Decode` 失败返回 400，业务错误不能 panic
- **中间件结构**：`func(next http.Handler) http.Handler`，一层包一层形成洋葱模型
- **请求 Context**：从 `r.Context()` 派生并用 `r.WithContext(ctx)` 传给下游，不能存进 struct
- **内存共享数据**：HTTP 请求并发执行，map 和计数器要用 `Mutex` 或 `RWMutex` 保护
- **标准库路由**：Go 1.22+ 的 `ServeMux` 支持 `"GET /path"` 这种 method + path 匹配

## 最容易踩的坑

1. **先写 body 再改状态码或 header**：响应可能已自动以 200 写出，后续设置不生效
2. **JSON 解码错误还继续执行**：输入无效就立即返回 400，避免零值进入业务逻辑
3. **把 Context 存进 struct**：Context 只活在一条请求调用链中，作为函数第一个参数向下传
4. **普通 map 扛并发 HTTP 请求**：会产生 data race，严重时 `concurrent map writes` 直接崩溃
5. **手写路由时不区分 404 和 405**：路径不存在是 404，路径存在但 method 不支持是 405

## 怎么自测

```bash
go test ./net_http
```

测试覆盖 TODO API 的成功路径、JSON 和参数错误、404 与 405、并发创建

> 当前 Windows 386 环境不支持 `go test -race`
> 代码已使用 `RWMutex` 保护共享 map 和自增 ID，切换到支持 race detector 的环境后再运行 `go test -race ./net_http`

## 下一步

进入 `gin/` 学习目录

1. Gin 路由、query/path/JSON 参数绑定、JSON 响应
2. CORS、Logger、Recovery 中间件
3. validator 参数校验和统一错误处理
4. JWT 认证
5. Swagger 和分层项目骨架

> 把 [05_todo_api.go](05_todo_api.go) 用 Gin 重写一次
> 重点观察 Gin 如何替你完成路由分发、参数绑定和 JSON 响应
