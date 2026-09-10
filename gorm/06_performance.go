// 06_performance.go 连接池 批量操作 预编译 查询瘦身
//
// 数据量一上来 用法对不对性能差一个数量级
// 本章四件事 连接池调参 批量插入 分批读取 预编译
//
// 配套 MySQL 知识
// 索引和 EXPLAIN 是 MySQL 主场 补 MySQL 时重点学
// GORM 侧只要记住 常查的列建索引 模型里加 index tag
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func demoPerformance() {
	fmt.Println("\n========== 06 性能 ==========")
	db := connectDB()
	resetTables(db, "users")

	// ---- 连接池 ----

	// 建一条 MySQL 连接很贵 TCP 握手 账号认证 一套下来毫秒级
	// 连接池就是用完不关 还回去 下次直接拿 省掉重复建连
	// GORM 底层是标准库 database/sql 自带连接池 调三个参数就够
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(25)                 // 同时打开的连接上限 超出的请求排队等
	sqlDB.SetMaxIdleConns(10)                 // 空闲时保留的连接数 复用全靠它们
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 一条连接最长活多久 到点强制换新
	fmt.Println("连接池已设置 上限 25 空闲 10 寿命 5 分钟")

	//! ConnMaxLifetime 必须设 MySQL 默认 8 小时主动断开空闲连接
	//! 不设的话 池里拿着一条已被服务端关掉的连接去查 直接报错
	//! invalid connection 而且半夜低峰期才出现 白天怎么测都正常

	fmt.Println()

	// ---- 批量插入 ----

	// 反面 循环单条插入 2000 次就是 2000 条 SQL 2000 次网络往返
	start := time.Now()
	for i := 0; i < 2000; i++ {
		db.Create(&User{
			Name:  fmt.Sprintf("单条%04d", i),
			Email: fmt.Sprintf("single%04d@t.com", i),
		})
	}
	slowCost := time.Since(start)

	// 正面 CreateInBatches 每批一条 SQL
	start = time.Now()
	batch := make([]User, 0, 2000)
	for i := 0; i < 2000; i++ {
		batch = append(batch, User{
			Name:  fmt.Sprintf("批量%04d", i),
			Email: fmt.Sprintf("batch%04d@t.com", i),
		})
	}
	db.CreateInBatches(&batch, 500) // 每批 500 条 共 4 条 SQL
	fastCost := time.Since(start)

	fmt.Printf("循环单条 2000 条耗时 %v\n批量插入 2000 条耗时 %v\n", slowCost, fastCost)
	fmt.Println("本地都差这么多 想想线上网络更差的时候")

	fmt.Println()

	// ---- 分批读取 FindInBatches ----

	// 表里几十万行不能一次全拉进内存 分批读 一批一批处理
	var holders []User
	var total int64
	db.Model(&User{}).Where("name LIKE ?", "批量%").
		FindInBatches(&holders, 500, func(tx *gorm.DB, batchNo int) error {
			total += int64(len(holders))
			fmt.Printf("  第 %d 批 %d 条\n", batchNo, len(holders))
			return nil // 返回 error 提前终止
		})
	fmt.Println("分批读完共", total, "条")

	fmt.Println()

	// ---- 预编译 PrepareStmt ----

	// 同一条 SQL 反复执行 MySQL 每次都要重新解析生成执行计划 纯属浪费
	// PrepareStmt 开启后第一次解析 后续直接复用
	// 常见误区 预编译防的不是注入 防注入靠问号占位符 两码事
	session := db.Session(&gorm.Session{PrepareStmt: true})
	for i := 0; i < 3; i++ {
		var u User
		session.Where("name = ?", "单条0001").First(&u)
	}
	fmt.Println("预编译 Session 三连查 SQL 只解析一次")
	// 全局开启就是在 connectDB 的 gorm.Config 里加 PrepareStmt true

	// ---- 查询瘦身 ----

	// Select 只查需要的列 顺手养成习惯
	// 索引 where 和 order 常用的列记得建 模型里就是加 index tag
	// 一条查询走没走索引用 EXPLAIN 看 这是 MySQL 主场

	fmt.Println()
	fmt.Println("🔑 日常 CRUD 交给 GORM 复杂统计和多表 JOIN 手写 Raw 反而直接")
}
