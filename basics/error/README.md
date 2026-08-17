# Go 错误处理（Error）专题

> 面向「前端（TS/JS）转 Go 后端」的学习代码，注释尽量直白，建议边跑边改。
> ⚠️ 错误是 Go 后端代码里出现频率最高的类型：`if err != nil` 能占业务代码的 30%~50%。
>   面试官看代码第一眼就看错误处理规不规范 —— 这是"脸面"级主题。
>   （basics/02_error.go 是速览版，这里是完整版。）

## 怎么运行

```bash
cd basics/error
go run .          # 一次跑完所有演示
```

## 学习顺序（按阅读顺序学）

| 文件 | 内容 | 重点程度 |
| --- | --- | --- |
| [01_what.go](01_what.go) | error 接口 + 返回值模式（Go 没有 try/catch） | 🔥 基础 |
| [02_make.go](02_make.go) | 造错误：errors.New / fmt.Errorf / 哨兵错误 | 🔥 日常 |
| [03_wrap.go](03_wrap.go) | 错误链：%w 包裹 + errors.Is / errors.As | 🔥🔥 面试高频 |
| [04_custom.go](04_custom.go) | 自定义错误类型：错误带结构化数据 | 🔥 Gin 校验同款 |
| [05_practice.go](05_practice.go) | 后端规范：不吞错误 / 边界打日志 / panic·recover | 规范 |

## 一句话总结

| 概念 | 一句话 | 前端对照 |
| --- | --- | --- |
| error | 一个接口（`Error() string`），函数返回它 | throw + try/catch（Go 没有异常） |
| 返回值模式 | 成功 → (有值, nil)；失败 → (零值, 错误) | —— |
| errors.New | 造一个只有文本的错误 | `new Error('...')` |
| fmt.Errorf | 带格式化参数的错误 | 模板字符串拼接 |
| 哨兵错误 | 包级 `var ErrXxx` 预置错误，全局统一判断 | —— |
| %w 包裹 | 错误带上下文，又不丢底层错误 | —— |
| errors.Is | 穿透错误链判断"是不是这个哨兵" | instanceof（但能穿透包裹层） |
| errors.As | 穿透错误链取出"某种类型"的错误 | 类型转换 + instanceof |
| 自定义错误类型 | 实现 `Error() string` 就是错误，能带数据 | 自定义 Error 类 |
| panic/recover | 程序级崩溃 + 在 defer 里接住（几乎不用） | throw 全局 catch（Go 更少用） |

## 面试速背（八股，复习时只看这一节）

- **Go 没有异常**：错误是返回值，函数多返回一个 error；`err != nil` 判断
- **返回惯例**：成功 `(有值, nil)`，失败 `(零值, error)`；nil = 没出错
- **error 是接口**：`type error interface { Error() string }`，实现 `Error()` 就是错误
- **哨兵错误**：包级 `var ErrXxx = errors.New(...)`；别内联 New 后靠文本判断
- **%w vs %v**：`%w` 建立错误链（Is/As 能穿透），`%v` 只拼文本（底层错误丢失）
- **errors.Is**：`errors.Is(err, ErrXxx)` 穿透包裹链判断"是不是这个错误"（找值/哨兵）
- **errors.As**：`var e *XxxErr; errors.As(err, &e)` 穿透包裹链取"某种类型"（找类型/取数据）
- **不吞错误**：要么 return 带上下文，要么打日志，二选一别都不做
- **日志只在上层打**：下层 return+包裹，上层统一记录，别每层都打
- **panic 只用于不该发生的事**（启动失败 / 不可能的分支）；业务错误一律 error
- **recover 只能在 defer 里用**；Gin Recovery 中间件就是它（1.20+ 有 errors.Join 合并多个错误）

## 最容易踩的 4 个坑

1. **吞错误**：`_ = f()` 或拿到错误不打不返，bug 凭空消失（日志 / 返回二选一）
2. **用错误文本判断**：`strings.Contains(err.Error(), "不存在")` 文本一改就失灵，用哨兵 + Is
3. **%v 代替 %w**：包了但 Is/As 穿透不了，底层错误丢了
4. **业务错误 panic**：handler 里 panic 会崩掉整个进程，业务错误一律 return error

## 下一步

- **写 Gin 时回看**：handler 里 `if err != nil` 天天见；GORM 返回 `ErrRecordNotFound` 用 `errors.Is` 判断
- **学中间件时回看**：Recovery 中间件 = defer + recover；Logger = 上层统一打日志
- **第四阶段八股再补**：error 底层结构（errorString / wrapError 指针链）现在不用管
- errors.Join（1.20+）收集表单校验的多条错误，用到时再展开
