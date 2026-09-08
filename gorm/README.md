# GORM：Go 的 ORM 速通

> 面向「前端转 Go 后端」的学习代码
> 前置：[gin/](../gin/) 已完成
> 数据库操作前端基本不碰 本章不强行前端对照 以理解概念为主
> ORM 思路一句话 和 Prisma / Sequelize 同款 用操作对象代替写 SQL
>
> 配套：MySQL 基础同步自补 增删改查 索引 JOIN 会用即可 复杂 SQL 不是本章重点

## 怎么运行

```bash
cd gorm
docker compose up -d   # 启动 MySQL 首次拉镜像稍等
go run .               # 跑全部 demo
```

- 依赖（gorm、mysql 驱动、bcrypt）已写入根目录 go.mod 无需手动安装
- 本机 3306 被占用：改 `docker-compose.yml` 的端口映射 再同步改 `01_connect_model.go` 的 dsn
- 想单独精读某章：到 `main.go` 注释掉其他调用
- 想彻底重置数据：`docker compose down -v` 再 `up -d`
- 每章开头自己清表造数据 单独跑任何一章都没问题

## 注释标记

| 标记 | 含义 |
| --- | --- |
| `//!` | 必踩的坑 / 必记重点 |
| `//?` | 思考题 |
| `//*` | 关键铺垫 |

> Better Comments 标记只出现在函数体内

## 学习顺序

| 文件 | 内容 | 关键点 |
| --- | --- | --- |
| [01_connect_model.go](01_connect_model.go) | 连接、模型定义、AutoMigrate、表名规则 | gorm.Model 四字段、parseTime |
| [02_create_read.go](02_create_read.go) | Create、First/Take/Last、Where、分页、Count | 主键回填、struct 零值坑 |
| [03_update_delete.go](03_update_delete.go) | Update/Updates/Save、软删除、原生 SQL、Scope | 软删除原理、Updates 零值 |
| [04_association.go](04_association.go) | 四种关联、Preload、Joins、Association | N+1 问题 面试必问 |
| [05_transaction_hook.go](05_transaction_hook.go) | 事务（闭包式、手动式）、Hook | Hook 返回 error 即拦截 |
| [06_performance.go](06_performance.go) | 连接池、批量、FindInBatches、预编译 | ConnMaxLifetime、批量插入 |

## 一句话总结

| 概念 | 一句话 |
| --- | --- |
| ORM | 操作 Go 结构体 GORM 负责翻译成 SQL |
| `gorm.Model` | ID + CreatedAt + UpdatedAt + DeletedAt 四件套 内嵌即得 |
| AutoMigrate | 按模型建表 只加列不删不改 学习用 生产走迁移脚本 |
| gorm tag | `size` 长度 `not null` 非空 `index` 索引 `uniqueIndex` 唯一 `type` 列类型 |
| `First`/`Take`/`Last` | 主键升序 / 不排序 / 主键降序 各取一条 |
| 主键回填 | `Create` 返回后 `u.ID` 已有值 不用再查 |
| `ErrRecordNotFound` | First 查不到报它 Find 查不到返回空切片不报错 |
| 软删除 | Delete 变 UPDATE deleted_at 查询自动加 `deleted_at IS NULL` |
| `Unscoped` | 绕过软删除 查已删的 或真删 |
| `Preload` | 关联数据一次查回 不用就是 N+1 |
| 事务 | 闭包式返回 error 自动回滚 里面只用 tx 不用 db |
| Hook | `BeforeCreate` 写前加工 `AfterFind` 读后加工 返回 error 即拦截 |
| `gorm.Expr` | `balance = balance - 100` 数据库端原子计算 防并发丢失更新 |
| 连接池 | MaxOpenConns 上限 MaxIdleConns 空闲 ConnMaxLifetime 寿命 必设 |

## 四种关联速记

```
BelongsTo  我属于对方   外键在自己身上    文章表存 user_id
HasOne     我有一个     外键在对方身上    证件表存 user_id
HasMany    我有很多     外键在对方身上    文章表存 user_id（和 HasOne 同一个外键 只是数量）
Many2Many  多对多       谁也不存谁       中间表 user_tags 存两边的 id
```

## 面试速背

- **软删除怎么实现**：模型内嵌 `gorm.Model`（有 `DeletedAt` 字段）Delete 时只写时间戳 查询自动拼 `WHERE deleted_at IS NULL` 真删用 `Unscoped`
- **First / Take / Last 区别**：First 按主键升序 Take 不排序 Last 主键降序
- **Where / Updates 传 struct 的坑**：零值字段（0、空串、false）被忽略 精确匹配和更新成零值必须用 map 或字符串占位符
- **N+1 问题**：查列表后循环查关联 = 1+N 条 SQL 解法 `Preload` 预加载 一条 IN 查回 条数与记录数无关
- **Preload vs Joins**：Preload 多发一条 IN 查询 适合多对多和大列表 Joins 单条 JOIN 适合 BelongsTo 一对一
- **事务两种写法**：闭包式 `db.Transaction(func(tx *gorm.DB) error)` 返回 error 自动回滚（推荐） 手动式 Begin/Commit/Rollback 加 `defer Rollback()` 兜底
- **事务里只用 tx**：用外层 db 的操作脱离事务 回滚救不了它
- **Hook 返回 error**：写入直接中断 事务场景整体回滚 典型用法 BeforeCreate 密码加密 AfterFind 脱敏
- **防 SQL 注入**：永远用 `?` 占位符 不拼字符串 原生 Raw/Exec 同样
- **PrepareStmt 防的是重复解析 不是注入** 防注入靠占位符 两码事
- **GORM v2 默认禁止无条件的全局 Update/Delete** 返回 `ErrMissingWhereClause` 防呆设计
- **连接池为什么必须设 ConnMaxLifetime**：MySQL 默认 8 小时断开空闲连接 池里拿到死连接报 invalid connection 且只在低峰期出现
- **金额用整数存分**：浮点算钱有精度丢失
- **并发扣余额**：`gorm.Expr("balance - ?")` 让数据库原子计算 避免「先读后写」的丢失更新 加行锁用 `Clauses(clause.Locking{Strength: "UPDATE"})`

## 最容易踩的坑

1. **DSN 忘写 parseTime=True**：时间字段扫描直接报错 新手第一个报错
2. **struct 条件零值被忽略**：`Where(User{Age: 0})` 变成无条件查全表
3. **Updates 想改成 0 传 struct**：零值被忽略 改 0 必须用 map
4. **关联字段查出来是空的**：GORM 默认不查关联 必须 `Preload`
5. **事务里用了外层 db**：那几条 SQL 不在事务里 回滚时救不了
6. **AutoMigrate 当生产迁移工具**：它只加列 不删不改 生产表结构变更走 SQL 迁移脚本
7. **不设 ConnMaxLifetime**：半夜低峰期拿到死连接报 invalid connection 白天测不出来
8. **循环单条插入**：N 条数据 N 次网络往返 用 `CreateInBatches`

## 下一步

1. **用 GORM 重写 TODO API**：落地分层架构 `handler → service → repository` 三层 存储从内存切片换成 MySQL（Week3 Day6-7 + GORM 实战合并完成）
2. **补 Swagger 文档**
3. 之后 Week4 后半段：**Redis**（缓存、Session、分布式锁）

> 学习方法：每个 `//!` 的坑都亲手试一次 看一眼真实报错 再改回来
