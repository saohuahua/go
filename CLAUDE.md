# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 仓库是什么

前端转 Go 后端的**个人学习仓库**（目标：秋招 Go 后端 / 全栈）。学习总路线见 [doc/前端转全栈（Go方向）学习路线.md](doc/前端转全栈（Go方向）学习路线.md)。

## 常用命令

```bash
go run ./slices     # 运行切片学习项目（或 cd slices && go run .）
go run ./map        # 运行 map 学习项目
go run ./gin        # 运行 Gin 学习项目（外部依赖已写入根 go.mod）
go vet ./slices ./map ./gin   # 编译 + 静态检查
```

- 没有测试；每个学习文件夹是独立可运行的 `package main`，多个 main 包共存没问题
- 本仓库是"学习代码"风格，**不是 gofmt 严格风格**：文件头 doc comment 用三空格缩进（gofmt 会重排成 tab），body 注释夹带 emoji
- ⚠️ gofmt 会把**连续多空行折叠成 1 行**，把 doc comment 里的 `//*`/`//!`/`//?` 标记改坏（加空格）—— 保存时格式化会触发，别和它对着干

## 学习项目结构约定（每个主题一个文件夹）

- `01_xxx.go`、`02_xxx.go`… 编号演示文件，每个文件一个 `demoXxx()` 函数
- `main.go` 入口逐个调用 demo，文件头注释列出全部文件的作用和运行方法
- `README.md` 固定四节：怎么运行 → 学习顺序表 → **一句话总结**（含"前端对照"列）→ **面试速背**（八股复习锚点）→ 最容易踩的坑
- 新主题照此结构建；**拆分不要太细**：目标是前期快速上手、后期背八股时靠 README 的面试速背节回顾

## 注释风格（用户个人偏好，务必遵守）

- **抓核心、不说废话**，每个知识点尽量配 **JS/TS 前端对照**（用户是从前端转的，靠对照迁移思维）
- 必要处铺垫概念（如讲 map 先提一句哈希表），铺垫要有用，不堆背景
- Better Comments 插件标记，**只能用在函数体内**（doc comment 位置会被 gofmt 改坏）：
  - `//!` 红色：必踩的坑 / 必记重点（面试会问），写法上"注释给出正确做法 + 代码留一行 panic 让人解开试"
  - `//?` 蓝色：思考题（先想答案，再对注释里的答案）
  - `//*` 绿色：关键铺垫 / 前端对照的类比
- 段落之间喜欢多空 1-2 行留白（已知 gofmt 会折叠，能留则留）
- 关键结论用 🔑 / ⚠️ emoji 标记，和 fmt.Println 输出配合

## 当前进度（截至 2026-08）

- 已完成：`slices/`（切片）、`map/`（map）→ Day4 完成
- 已完成：`pointers/`（指针）→ Day5 完成（go mod 速背在 pointers/README 补充节）
- 已完成：`goroutine/`（Goroutine + Channel）→ Day6 完成
- 已完成：`sync_context/`（Mutex/select/Context 并发控制）→ Day7 完成
- 已完成：`basics/`（接口/error/defer/闭包 前置补课）→ 进 Gin 前补的 4 个基础语法
- 已完成：`basics/interface/`（接口专题：定义/使用场景/mock/标准库/断言）→ 接口完整版
- 已完成：`json_file/`（JSON 序列化/文件 IO/Reader-Writer）→ Week2 Day2
- 已完成：`net_http/`（HTTP Server 预热：Handler/请求响应/路由中间件/内存 TODO API）→ Week2 Day4-7
- 已完成：`review_practice/`（综合回顾练习：从零补全内存 TODO 服务）
- 已完成：`gin/`（Gin 框架：路由/参数绑定/中间件/校验/JWT/用 Gin 重写 TODO API）→ Week3 完成
- 已完成：`gorm/`（GORM + MySQL：连接/模型/CRUD/软删除/关联 Preload/事务/Hook/连接池性能）→ Week4 Day3-4（MySQL SQL 基础同步自补中，当前跑在本地 MySQL 8.0，gorm/ 也留了 docker-compose.yml 备用）
- 下一步：**用 GORM 重写 TODO 落地分层架构**（cmd/internal/pkg 分层 + Swagger）→ 之后 Week4 后半段 Redis
- 面试高频清单（切片扩容、Map 底层、GMP 等）在路线文档第四阶段，README 面试速背节对应着写
