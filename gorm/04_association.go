// 04_association.go 关联关系与预加载
//
// 四种关联一次看懂
// BelongsTo  我属于对方  外键在自己身上   文章表存 user_id 文章属于用户
// HasOne     我有一个    外键在对方身上   身份证表存 user_id 用户有一张证
// HasMany    我有很多    外键在对方身上   文章表存 user_id 用户有很多文章
// Many2Many  多对多      谁也不存谁      中间表 user_tags 存两边的 id
// HasOne 和 HasMany 是同一个外键的两种数量描述 一个是一条 一个是一堆
//
// 本章另一个主角 Preload 预加载 不用它就是面试必挂的 N+1 问题
package main

import (
	"fmt"
)

func demoAssociation() {
	fmt.Println("\n========== 04 关联与预加载 ==========")
	db := connectDB()

	// 清空四个表 顺序是先清中间表和子表 最后删 users
	db.Exec("DELETE FROM article_tags")
	db.Exec("DELETE FROM user_tags")
	db.Unscoped().Where("1 = 1").Delete(&Article{})
	db.Unscoped().Where("1 = 1").Delete(&Tag{})
	db.Unscoped().Where("1 = 1").Delete(&User{})

	// ---- 级联创建 一次 Create 连关联数据一起入库 ----

	// 用户带文章 文章带标签 一条 Create 全部入库
	// 外键自动回填 不用手工设置 article.user_id
	author := User{
		Name:  "作者甲",
		Email: "author@t.com",
		Age:   28,
		Articles: []Article{
			{Title: "Go 入门", Content: "从切片说起", Tags: []Tag{{Name: "Go"}, {Name: "教程"}}},
			{Title: "GORM 入门", Content: "从连接说起", Tags: []Tag{{Name: "Go"}}},
		},
	}
	db.Create(&author)
	for i := range author.Articles {
		fmt.Printf("级联创建 %s 的文章 %s 外键 UserID 自动回填 = %d\n",
			author.Name, author.Articles[i].Title, author.Articles[i].UserID)
	}

	// 再来一个用户 一篇文章 不带标签
	db.Create(&User{
		Name: "作者乙", Email: "author2@t.com", Age: 35,
		Articles: []Article{{Title: "跑步日记", Content: "今天五公里"}},
	})

	fmt.Println()

	// ---- 不预加载 关联字段是空的 ----

	var someone User
	db.First(&someone, "name = ?", "作者甲")
	fmt.Println("不 Preload 直接读关联 文章数 =", len(someone.Articles))
	// GORM 的关联字段默认不查 要显式 Preload 它才查
	// 这点不同于很多 ORM 的懒加载 懒加载是你访问时它自动发 SQL

	// ---- N+1 问题 面试必问 ----

	// 反面教材 先查所有用户 再循环查每个人的文章
	var users []User
	db.Find(&users) // 第 1 条 SQL
	for i := range users {
		var arts []Article
		db.Debug().Where("user_id = ?", users[i].ID).Find(&arts) // 每人再来 1 条
		fmt.Printf("  %s 有 %d 篇文章\n", users[i].Name, len(arts))
	}
	// Debug 打印能看到 2 个用户跑了 1 + 2 = 3 条 SQL
	// 用户涨到 1000 就是 1001 条 每条都是一次网络往返 接口就是这么被拖死的
	// 一条查列表 循环里再逐条查关联 这个模式就叫 N+1

	fmt.Println()

	// 正确姿势 Preload 预加载
	var users2 []User
	db.Debug().Preload("Articles").Find(&users2)
	// 只有 2 条 SQL 一条查 users 一条用 IN 把所有文章一次查回
	// SQL 条数和用户数无关 用户再多也是 2 条
	for i := range users2 {
		fmt.Printf("  %s 文章数 = %d\n", users2[i].Name, len(users2[i].Articles))
	}

	fmt.Println()

	// ---- Preload 的三种变体 ----

	// 条件预加载 第二个参数传条件 只预加载标题带 Go 的文章
	var goOnly []User
	db.Preload("Articles", "title LIKE ?", "%Go%").Find(&goOnly)
	for i := range goOnly {
		fmt.Printf("条件预加载 %s 只带 Go 文章 %d 篇\n", goOnly[i].Name, len(goOnly[i].Articles))
	}

	// 嵌套预加载 用户到文章再到文章的标签 用点号一路点下去
	var deep []User
	db.Preload("Articles.Tags").Find(&deep)
	for i := range deep {
		for j := range deep[i].Articles {
			fmt.Printf("嵌套预加载 %s 的 %s 标签数 = %d\n",
				deep[i].Name, deep[i].Articles[j].Title, len(deep[i].Articles[j].Tags))
		}
	}

	// Joins 预加载 适合 BelongsTo 这种一对一关系 一次 JOIN 全带回来
	var arts []Article
	db.Joins("User").Where("users.name = ?", "作者甲").Find(&arts)
	fmt.Println("Joins 预加载 第一篇属于", arts[0].User.Name)
	// 多对多和大列表别用 Joins 结果集会行数膨胀 老实用 Preload

	fmt.Println()

	// ---- Association 直接操作关联 ----

	// Count 数关联 Append 追加关联 中间表自动维护
	cnt := db.Model(&users2[0]).Association("Tags").Count()
	fmt.Println("作者甲当前用户标签数 =", cnt)
	db.Model(&users2[0]).Association("Tags").Append(&Tag{Name: "Vlog"})
	cnt = db.Model(&users2[0]).Association("Tags").Count()
	fmt.Println("Append 追加 Vlog 后 =", cnt)

	fmt.Println("🔑 关联字段默认不查 想要就 Preload 循环查关联就是 N+1")
}
