// main.go GORM 入口 连接 增查改删 关联 事务 性能
//
// 本文件夹共 6 个演示文件 + 本入口文件
//
// 01_connect_model.go     连接 MySQL 模型定义 AutoMigrate 建表
// 02_create_read.go       增 查 Create First Where 分页
// 03_update_delete.go     改 删 软删除 原生 SQL Scope
// 04_association.go       四种关联 Preload 预加载 N+1 问题
// 05_transaction_hook.go  事务 生命周期 Hook
// 06_performance.go       连接池 批量 预编译 查询瘦身
//
// 运行方法
//
// cd gorm && go run .
// go run ./gorm    在项目根目录执行
//
// 学习方法
//
// 1 先 docker compose up -d 把 MySQL 跑起来 再 go run .
// 2 想精读某章就把其他的注释掉 只留一个跑
// 3 每个 //! 的坑都亲手试一次 再改回来
package main

import "fmt"

func main() {
	demoConnectModel()    // 01 连接与模型
	demoCreateRead()      // 02 增与查
	demoUpdateDelete()    // 03 改与删
	demoAssociation()     // 04 关联与预加载
	demoTransactionHook() // 05 事务与 Hook
	demoPerformance()     // 06 性能

	fmt.Println("\n🎉 GORM 速通完成 下一步用 GORM 重写 TODO API 落地分层架构")
}
