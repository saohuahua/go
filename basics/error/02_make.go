// 02_make.go —— 造错误：errors.New / fmt.Errorf / 哨兵错误
//
// 前端对照一句话：
//
//	JS：new Error("...")  /  throw new Error(`余额不足，需要 ${need}`)
//	Go：errors.New("...") /  fmt.Errorf("余额不足，需要 %d", need)
//	差别：JS 造完要 throw 出去，Go 造完是 return 回去（01 说过，错误是返回值）
//
// 本节三个技能点：
//
//	① errors.New —— 最简单，一句话文本
//	② fmt.Errorf —— 能拼变量（带格式化的文本）
//	③ 哨兵错误  —— 把预设错误定义成包级变量，全局统一判断（重点）
package main

import (
	"errors"
	"fmt"
)

// ---------- 哨兵错误：包级变量，只创建一次，全局共用 ----------
// "哨兵" = 站岗的哨兵：提前在包入口站好预设错误，所有人都引用它、用它做判断。
// 好处：判断错误不靠"文本长啥样"，而靠"是不是同一个值"（配合 03 的 errors.Is）。
var ErrNotFound = errors.New("记录不存在")
var ErrInvalidID = errors.New("非法的 id")

func demoMake() {
	fmt.Println("========== 02 造错误：errors.New / fmt.Errorf / 哨兵错误 ==========")

	// ---------- ① errors.New：最简单，只有一个文本 ----------
	err := errors.New("余额不足")
	fmt.Println("① errors.New 造的：", err)

	// ---------- ② fmt.Errorf：要拼变量时用它 ----------
	// %d 是数字占位、%s 是字符串占位 —— 类似前端模板字符串 ${}，只是参数要排在后面
	need, have := 100, 30 // 多变量一起 := 声明（一次声明两个）
	fmt.Println("② fmt.Errorf 拼参数：", fmt.Errorf("余额不足：需要 %d，当前 %d", need, have))

	// ---------- ③ 哨兵错误：为什么要包级变量，而不是随手 New ----------
	//! 坑：每次 errors.New 都是"全新的错误值"，== 永远比不中！
	fmt.Println("③ 两次 errors.New 是同一个值吗：", errors.New("记录不存在") == errors.New("记录不存在"))
	// 解法：包级 var ErrNotFound 只创建一次，全局共用同一个实例，判断就能对上
	fmt.Println("③ 包级哨兵 ErrNotFound：", ErrNotFound)

	//! 坑：也别用"错误文本"判断（如 strings.Contains(err.Error(), "不存在")）
	//!     文本一改就失灵；正确做法：哨兵错误 + errors.Is（03 讲）
	fmt.Println("\n🔑 造错误两种：errors.New（纯文本）/ fmt.Errorf（拼参数）；判断用哨兵，别用文本")
}
