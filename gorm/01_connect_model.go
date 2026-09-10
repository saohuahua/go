// 01_connect_model.go GORM 连接 MySQL 模型定义 建表
//
// ORM 全称对象关系映射 让 Go 结构体和数据库表一一对应
// 你写 Go 代码 GORM 负责翻译成 SQL 执行 再把结果装回结构体
// 大部分场景不用手写 SQL 但你得懂点 SQL 才能看懂 GORM 在干什么
//
// 本文件是全文件夹的地基
// connectDB        连接加建表 后面 5 个文件直接复用这个函数
// User 等 5 个模型  后面所有 demo 共用这一套模型
//
// 运行前提 在本目录先执行 docker compose up -d 把 MySQL 跑起来
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DSN 数据库连接字符串
// 格式 用户名 密码 @tcp(地址 端口)/库名?参数
// 生产环境不要硬编码 放配置文件或环境变量 这里为了学习直观
const dsn = "root:root@tcp(127.0.0.1:3306)/learn_gorm?charset=utf8mb4&parseTime=True&loc=Local"

// User 用户表 演示最常用的模型写法
//
// 内嵌 gorm.Model 等于白拿 4 个字段
// ID        uint      自增主键
// CreatedAt time.Time 创建时间 插入时自动填
// UpdatedAt time.Time 更新时间 每次更新自动刷新
// DeletedAt time.Time 删除时间 软删除的载体 03 专门讲
//
// 常用 gorm tag 速查
// size 50      字符串最大长度
// not null     非空
// uniqueIndex  唯一索引 重复插入直接报错
// index        普通索引 加速查询
// default 18   默认值
// type text    指定列类型 长文本用
// -            忽略该字段 不建列
type User struct {
	gorm.Model
	Name     string `gorm:"size:50;not null"`
	Email    string `gorm:"size:100;uniqueIndex"`
	Password string `gorm:"size:100"`         // 存密文 明文密码永远不进库 05 演示
	Age      int    `gorm:"default:18;index"` // 常查的字段建索引

	Articles []Article // HasMany 一对多 一个用户多篇文章 不需要 tag
	Tags     []Tag     `gorm:"many2many:user_tags"` // Many2Many 多对多 自动建中间表 user_tags
}

// Article 文章表
// UserID 是 BelongsTo 的外键 默认命名规则 拥有者的模型名加 ID
type Article struct {
	gorm.Model
	Title   string `gorm:"size:100;not null"`
	Content string `gorm:"type:text"` // 长文本 不限长度
	UserID  uint   // BelongsTo 外键 文章属于某个用户
	User    User   // BelongsTo 关联字段 声明属于谁 04 专门讲
	Tags    []Tag  `gorm:"many2many:article_tags"` // 文章和标签也是多对多
}

// Tag 标签表
type Tag struct {
	gorm.Model
	Name string `gorm:"size:50;not null"`
}

// Account 银行账户 演示事务专用
// 金额用整数存分 不用浮点 浮点算钱会丢精度 这是常识也是面试题
type Account struct {
	gorm.Model
	Owner   string `gorm:"size:50;not null"`
	Balance int    `gorm:"not null;default:0"` // 单位分
}

// LogRow 演示自定义表名
// GORM 默认表名是结构体名的蛇形加复数 User 变 users LogRow 变 log_rows
// 不想要默认名 实现 TableName 方法想叫什么叫什么
type LogRow struct {
	ID        uint
	Message   string `gorm:"size:200"`
	CreatedAt time.Time
}

// TableName 自定义表名
func (LogRow) TableName() string {
	return "sys_logs"
}

// connectDB 连接 MySQL 并建表 全文件夹复用
// AutoMigrate 放在这里 任何一个 demo 单独跑都不怕缺表 建过就跳过 幂等
func connectDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 物理外键约束会在删数据 迁移时碍事 真实项目更常用逻辑外键
		// 外键关系只在代码层面保证 建表时不生成物理约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 学习期把默认日志调静默 输出干净
		// 想看真实 SQL 就在那条查询前加 .Debug() 03 和 04 有示例
		// 真实项目一般用 logger.Warn 保留慢查询和错误日志
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("连接 MySQL 失败 %v\n请先在本目录执行 docker compose up -d", err)
	}

	// AutoMigrate 按模型定义建表 表已存在就对比补齐缺的列
	//! AutoMigrate 只加列 不删列 不改列类型 学习开发阶段用着爽
	//! 生产环境改表要走 SQL 迁移脚本 每次变更有记录 可回滚
	if err := db.AutoMigrate(&User{}, &Article{}, &Tag{}, &Account{}, &LogRow{}); err != nil {
		log.Fatalf("建表失败 %v", err)
	}
	return db
}

func demoConnectModel() {
	fmt.Println("\n========== 01 连接与模型 ==========")
	db := connectDB()

	// 建表成果检查 HasTable 查一张表存不存在
	fmt.Println("users 建好 =", db.Migrator().HasTable(&User{}))
	fmt.Println("articles 建好 =", db.Migrator().HasTable(&Article{}))
	fmt.Println("user_tags 中间表建好 =", db.Migrator().HasTable("user_tags"))
	fmt.Println("sys_logs 自定义表名建好 =", db.Migrator().HasTable("sys_logs"))

	// gorm.Model 四字段的实际效果 插一条看自动填充
	resetTables(db, "users")
	u := User{Name: "张三", Email: "zhangsan@t.com", Age: 25}
	db.Create(&u)
	fmt.Printf("ID 自增 %d  CreatedAt 自动填 %s  UpdatedAt 自动填 %s\n",
		u.ID, u.CreatedAt.Format("15:04:05"), u.UpdatedAt.Format("15:04:05"))

	//? Email 上有 uniqueIndex 再插一个 zhangsan@t.com 会怎样
	//? Create 返回 error MySQL 报 Duplicate entry 业务里就是提示邮箱已注册

	//! DSN 里 parseTime 不写 True 表里的时间装不进 time.Time 直接报错
	//! 新手遇到的第一个 GORM 报错基本是它 拿不准就对照文件顶部的 dsn 抄

	fmt.Println("🔑 内嵌 gorm.Model 主键和两个时间字段全自动维护 DeletedAt 留到 03 讲")
}

// resetTables 清空演示表并把自增 ID 归零 让每次运行的输出完全一致
// TRUNCATE 是 MySQL 的清表命令 比 DELETE 快 还会重置自增计数器
// GORM 没有对应 API 所以用 Exec 执行原生 SQL 这里表名写死没有注入风险
func resetTables(db *gorm.DB, tables ...string) {
	for _, t := range tables {
		db.Exec("TRUNCATE TABLE " + t)
	}
}
