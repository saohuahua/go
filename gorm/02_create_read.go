// 02_create_read.go 增与查
//
// 增 Create   查 First Take Last Find Where
// 日常业务的大头就是这两个动作 练熟
//
// 配套 MySQL 知识
// 增就是 INSERT 查就是 SELECT GORM 只是帮你拼这几句话
package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func demoCreateRead() {
	fmt.Println("\n========== 02 增与查 ==========")
	db := connectDB()

	// 每章开头清空演示表 自增 ID 也归零 每次运行输出完全一致
	resetTables(db, "users")

	// ---- 增 ----

	// 单条插入
	// 关键点 Create 返回后 u.ID 已被回填 不用再查一次数据库拿主键
	u := User{Name: "张三", Email: "zhangsan@t.com", Password: "123456", Age: 25}
	result := db.Create(&u)
	fmt.Println("单条插入 ID 回填 =", u.ID, " 影响行数 =", result.RowsAffected)

	// 批量插入 传切片
	users := []User{
		{Name: "李四", Email: "lisi@t.com", Age: 30},
		{Name: "王五", Email: "wangwu@t.com", Age: 18},
		{Name: "老六", Email: "laoliu@t.com", Age: 30},
		{Name: "赵七", Email: "zhaoqi@t.com", Age: 28},
	}
	db.Create(&users)
	for i := range users {
		fmt.Printf("批量插入 %s ID = %d\n", users[i].Name, users[i].ID)
	}

	fmt.Println()

	// ---- 查 ----

	// First 按主键升序取第一条
	var first User
	db.First(&first)
	fmt.Println("First 主键最小 =", first.Name)

	// Take 不排序取一条 Last 按主键降序取第一条
	// 面试爱问三者区别 区别就在排不排序 怎么排
	var take User
	db.Take(&take)
	var last User
	db.Last(&last)
	fmt.Println("Last 主键最大 =", last.Name)

	// First 第二个参数直接传主键值 最快的查法
	// 用刚插入的李四的 ID 别写死数字 自增 ID 会随运行环境变化
	var byID User
	db.First(&byID, users[0].ID)
	fmt.Println("按主键直查 ID", byID.ID, "=", byID.Name)

	fmt.Println()

	// ---- Where 三种写法 ----

	// 写法一 字符串加问号占位符 最常用最安全
	//? 为什么不直接拼接字符串 db.Where("name = " + name)
	//? 输入里带个单引号就能改写你的 SQL 这就是 SQL 注入
	//? 占位符让输入永远只当数据不当代码 从根上防注入 面试必问
	var u1 User
	db.Where("name = ?", "张三").First(&u1)
	fmt.Println("字符串条件查到 =", u1.Name)

	// 写法二 map 等值匹配 零值条件照样生效
	var byMap []User
	db.Where(map[string]any{"age": 30}).Find(&byMap)
	fmt.Println("map 条件 age 30 命中", len(byMap), "人")

	// 写法三 struct 看着优雅但藏着本章最大的坑
	//! struct 条件会忽略零值字段 年龄 0 空字符串 false 全部被跳过
	//! 本意查 0 岁用户 结果生成了没有任何条件的查询 查回全表 5 条
	var byStruct []User
	db.Where(User{Age: 0}).Find(&byStruct)
	fmt.Println("struct 条件 Age 0 实际查回", len(byStruct), "条 不是 0 条 坑")
	// 想让 struct 条件支持零值 把字段定义成指针 *int 或 sql.NullInt32
	// 学习阶段记结论 等值条件无脑用字符串占位符或 map

	fmt.Println()

	// ---- 更多查询姿势 ----

	// Find 查多条 查不到不报错 返回空切片
	var nobody []User
	db.Where("age > ?", 99).Find(&nobody)
	fmt.Println("Find 查不到 len =", len(nobody), " 不报错")

	// First 查不到会报 ErrRecordNotFound 这是标准判空姿势
	// GORM 默认会给它打红色日志 这里日志已调静默所以看不到 真实项目会看到
	var missing User
	err := db.Where("age > ?", 99).First(&missing).Error
	fmt.Println("First 查不到 errors.Is 判定 =", errors.Is(err, gorm.ErrRecordNotFound))

	// Like 模糊查询 百分号是通配符 代表任意长度的任意字符
	var likes []User
	db.Where("name LIKE ?", "李%").Find(&likes)
	fmt.Println("Like 李开头的 =", likes[0].Name)

	// IN 多值匹配
	var ins []User
	db.Where("name IN ?", []string{"张三", "李四"}).Find(&ins)
	fmt.Println("IN 两个名字命中", len(ins), "人")

	// 多个 Where 链式叠加 条件之间是 AND
	var young []User
	db.Where("age >= ?", 20).Where("age < ?", 30).Find(&young)
	fmt.Println("20 到 30 岁之间", len(young), "人")

	fmt.Println()

	// ---- 排序 分页 统计 ----

	// Order 排序 desc 降序 asc 升序 不写默认升序
	// 分页三件套 Offset 跳过几条 Limit 取几条
	// 换算公式 第 N 页就是 Offset((N-1) * 每页条数).Limit(每页条数)
	var page []User
	db.Order("age desc").Offset(1).Limit(2).Find(&page)
	for i := range page {
		fmt.Printf("排序分页取到 %s %d岁\n", page[i].Name, page[i].Age)
	}

	// Count 计数 参数是 int64 指针
	var total int64
	db.Model(&User{}).Where("age = ?", 30).Count(&total)
	fmt.Println("30 岁人数 =", total)

	// Select 只查需要的列 不用的列不传输 不占内存
	var onlyName []User
	db.Model(&User{}).Select("name", "age").Find(&onlyName)
	fmt.Println("Select 只取姓名年龄 首个 =", onlyName[0].Name, onlyName[0].Age)

	fmt.Println("🔑 Create 后主键自动回填 等值条件用占位符或 map struct 有零值坑")
}
