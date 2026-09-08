// 03_update_delete.go 改 删 软删除 原生SQL 条件复用
//
// 改 Update Updates Save   删 Delete
// 内嵌 gorm.Model 的模型删除时自动走软删除 这是本章重点
//
// 配套 MySQL 知识
// 改就是 UPDATE 删就是 DELETE
// 软删除其实是把 DELETE 换成一条 UPDATE 眼见为实见下方 Debug
package main

import (
	"fmt"

	"gorm.io/gorm"
)

// adultsScope 把常用查询条件封装成函数 到处复用
// 入参出参都是 *gorm.DB 这类函数叫 Scope
func adultsScope(db *gorm.DB) *gorm.DB {
	return db.Where("age >= ?", 18)
}

func demoUpdateDelete() {
	fmt.Println("\n========== 03 改与删 ==========")
	db := connectDB()
	db.Unscoped().Where("1 = 1").Delete(&User{})
	db.Create(&[]User{
		{Name: "张三", Email: "zhangsan@t.com", Age: 25},
		{Name: "李四", Email: "lisi@t.com", Age: 30},
		{Name: "王五", Email: "wangwu@t.com", Age: 17},
	})

	// ---- 改 ----

	// Update 改单列 必须带条件 返回值看影响行数
	res := db.Model(&User{}).Where("name = ?", "张三").Update("age", 26)
	fmt.Println("Update 影响行数 =", res.RowsAffected)

	// Updates 一次改多列
	//! Updates 传 struct 只更新非零字段 和 02 的 Where 同款零值坑
	var wang User
	db.First(&wang, "name = ?", "王五")
	db.Model(&wang).Updates(User{Age: 0}) // Age 0 是零值 被直接忽略
	db.First(&wang, wang.ID)
	fmt.Println("想把王五改成 0 岁 实际还是", wang.Age, "岁 struct 的 0 没写进去")

	// 想把字段更新成零值 必须用 map 精确指定
	db.Model(&wang).Updates(map[string]any{"age": 0, "name": "王五五"})
	db.First(&wang, wang.ID)
	fmt.Println("map 写 0 生效 现在叫", wang.Name, wang.Age, "岁")

	// Save 全字段保存 查出来 改完 整个存回去 零值也存
	var zhang User
	db.First(&zhang, "name = ?", "张三")
	zhang.Age = 27
	zhang.Email = "zhangsan2@t.com"
	db.Save(&zhang)
	fmt.Println("Save 全字段保存 张三 =", zhang.Age, "岁")

	fmt.Println()

	// ---- 删 ----

	// Delete 软删除 模型内嵌了 gorm.Model 就自动生效
	// Debug 会把真实 SQL 打到终端 眼见为实
	db.Debug().Where("name = ?", "李四").Delete(&User{})
	// 打印出来的不是 DELETE 而是
	// UPDATE users SET deleted_at='...' WHERE name='李四' AND deleted_at IS NULL
	// 数据没被删 只是把 deleted_at 打上时间戳
	// 之后所有查询自动拼上 WHERE deleted_at IS NULL 这就是软删除的全部原理

	var count int64
	db.Model(&User{}).Count(&count)
	fmt.Println("软删除后普通查询只剩", count, "人 数据其实还在表里")

	// Unscoped 绕过软删除 查得到已删的
	var ghost User
	db.Unscoped().Where("name = ?", "李四").First(&ghost)
	fmt.Println("Unscoped 查到已删除的李四 ID =", ghost.ID)

	// Unscoped 加 Delete 才是真删 物理删除
	db.Unscoped().Where("name = ?", "李四").Delete(&User{})
	var ghostCount int64
	db.Unscoped().Model(&User{}).Where("name = ?", "李四").Count(&ghostCount)
	fmt.Println("物理删除后连 Unscoped 都查不到了 =", ghostCount)

	fmt.Println()

	// ---- 原生 SQL ----

	// GORM 拼不出来的复杂查询 直接写 SQL
	// Raw 查询 用 Scan 把结果装进结构体
	var adults []User
	db.Raw("SELECT id, name, age FROM users WHERE age >= ?", 18).Scan(&adults)
	fmt.Println("Raw 查询成年人", len(adults), "人")

	// Exec 执行增删改 返回影响行数
	res = db.Exec("UPDATE users SET age = ? WHERE name = ?", 26, "王五五")
	fmt.Println("Exec 影响", res.RowsAffected, "行")

	//! 原生 SQL 里拼接用户输入等于把 SQL 注入大门敞开
	//! 哪怕是自己手写的 SQL 也永远用问号占位符

	// ---- Scope 条件复用 ----

	var scoped []User
	db.Scopes(adultsScope).Find(&scoped)
	fmt.Println("Scopes 复用成年人条件 命中", len(scoped), "人")

	fmt.Println("🔑 软删除本质是 UPDATE deleted_at 查询自动加 deleted_at IS NULL")
	fmt.Println("🔑 想更新成零值用 map 想全量保存用 Save")
}
