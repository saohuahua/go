// 05_transaction_hook.go Hook 生命周期钩子 事务
//
// Hook 挂在模型写入前后自动执行的函数 写一次全部生效
// 典型场景 插入前自动加密密码 查询后自动脱敏
//
// 事务把多条 SQL 捆成一捆 要么全成功要么全失败
// 典型场景 转账 一边扣钱失败 另一边的加钱必须跟着作废
//
// 配套 MySQL 知识
// 事务四大特性 ACID 原子性 一致性 隔离性 持久性 补 MySQL 时会细学
package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BeforeCreate 插入前触发
// 密码在这里加密 业务代码只管传明文 入库的永远是密文
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 学习造数场景允许空密码 真实项目这里应该直接拒绝
	if u.Password == "" {
		return nil
	}
	// 太短的密码直接拒绝 Hook 返回 error 这次 Create 就失败了
	if len(u.Password) < 6 {
		return errors.New("密码至少 6 位")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// AfterFind 查询后触发 查出来的数据先加工再交给业务
func (u *User) AfterFind(tx *gorm.DB) error {
	// 脱敏演示 密文也只给看个开头
	if len(u.Password) > 10 {
		u.Password = u.Password[:10] + "..."
	}
	return nil
}

func demoTransactionHook() {
	fmt.Println("\n========== 05 事务与 Hook ==========")
	db := connectDB()
	resetTables(db, "users", "accounts")

	// ---- Hook 效果 ----

	u := User{Name: "阿珍", Email: "azhen@t.com", Password: "123456", Age: 22}
	db.Create(&u)
	// Create 返回时 Password 已经变成密文 BeforeCreate 干的
	fmt.Println("写入前明文 123456 写入后 =", u.Password)

	var got User
	db.First(&got, u.ID)
	fmt.Println("读出来 AfterFind 自动脱敏 =", got.Password)

	// Hook 返回 error 写入直接失败
	short := User{Name: "短密码", Email: "short@t.com", Password: "123", Age: 20}
	if err := db.Create(&short).Error; err != nil {
		fmt.Println("短密码被 Hook 拦下 =", err)
	}

	fmt.Println()

	// ---- 事务 闭包式 最推荐 ----

	db.Create(&[]Account{
		{Owner: "阿珍", Balance: 10000},
		{Owner: "阿强", Balance: 5000},
	})

	//? 为什么不先查余额 在 Go 里减完再 UPDATE
	//? 两个请求同时读到 10000 各自减 100 都写回 就丢了 100 块 叫并发丢失更新
	//? gorm.Expr 生成 balance = balance - 100 让数据库一条 SQL 算完 天然避开

	// 转账成功版 闭包返回 nil 自动提交
	err := db.Transaction(func(tx *gorm.DB) error {
		// 注意事务里全用 tx 不用外面的 db 用了 db 就脱离事务了
		if err := tx.Model(&Account{}).Where("owner = ?", "阿珍").
			Update("balance", gorm.Expr("balance - ?", 100)).Error; err != nil {
			return err // 返回 error 整个事务自动回滚
		}
		if err := tx.Model(&Account{}).Where("owner = ?", "阿强").
			Update("balance", gorm.Expr("balance + ?", 100)).Error; err != nil {
			return err
		}
		return nil // 返回 nil 自动提交
	})
	fmt.Println("成功事务结果 =", err)
	printBalance(db, "阿珍")
	printBalance(db, "阿强")

	// 失败版 第二步故意报错 第一步的扣款自动回滚
	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&Account{}).Where("owner = ?", "阿珍").
			Update("balance", gorm.Expr("balance - ?", 100))
		// 模拟第二步失败 比如对方账户被冻结
		return errors.New("对方账户被冻结 转账失败")
	})
	fmt.Println("失败事务结果 =", err)
	printBalance(db, "阿珍") // 余额没变 第一步的扣款被回滚了
	printBalance(db, "阿强")

	fmt.Println()

	// ---- 事务 手动挡 ----

	// 想自己控制提交时机 用 Begin Commit Rollback
	tx := db.Begin()
	// defer 兜底回滚 事务没提交函数就退出时一律回滚
	// 已提交后再 Rollback 只返回一个错误 没有副作用
	defer tx.Rollback()
	if err := tx.Model(&Account{}).Where("owner = ?", "阿珍").
		Update("balance", gorm.Expr("balance - ?", 50)).Error; err != nil {
		fmt.Println("手动事务出错 回滚", err)
		return
	}
	tx.Commit()
	fmt.Println("手动事务提交完成")
	printBalance(db, "阿珍")

	// 嵌套事务用 SavePoint 记存档点 回滚只回到存档
	// tx.SavePoint("sp1")  tx.RollbackTo("sp1")  了解一下即可

	//? 高并发扣库存怎么防超卖
	//? 事务里给这行数据上锁 SELECT ... FOR UPDATE
	//? tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, id)
	//? 锁住后别人改不了 直到事务提交 面试高频

	fmt.Println("🔑 事务里只用 tx 闭包返回 error 即回滚 Hook 拦截等于写入失败")
}

// printBalance 打印某人余额 全章复用的小工具
func printBalance(db *gorm.DB, owner string) {
	var acc Account
	db.Where("owner = ?", owner).First(&acc)
	fmt.Printf("  %s 余额 = %d 分\n", owner, acc.Balance)
}
