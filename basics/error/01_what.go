// 01_what.go —— error 是什么：一个接口 + 返回值模式（Go 没有 try/catch）
//
// 前端对照一句话（先记这个）：
//
//	JS：throw + try/catch —— 错误被"抛出"，靠外层 catch 接住，不接就崩
//	Go：函数多返回一个 error —— 错误被"返回"，调用方自己判断
//	心态差别：JS 把错误当"意外"，Go 把错误当"返回值"（常态，必须显式处理）
//
// Go 的 error 其实就是一个接口（呼应 interface/04_stdlib.go）：
//
//	type error interface {
//	    Error() string // 返回一段错误描述文本
//	}
//
// 所以只要实现 Error() string，你的类型就自动是错误 —— 白嫖标准库的接口。
package main

import (
	"errors"
	"fmt"
)

// ---------- 返回错误的函数：多一个返回值 error ----------

// fetchData 模拟一个可能失败的调用，返回 (正常结果, error)。
// Go 惯例：成功 → (有值, nil)；失败 → (零值, 错误)。nil = 没出错。
func fetchData(ok bool) (string, error) {
	if !ok {
		// errors.New 是最简单的造错误，细节下一节 02 讲
		return "", errors.New("拉取数据失败")
	}
	return "数据内容", nil
}

// ping 模拟一个必然成功的调用：返回 nil 表示一切正常
func ping() error {
	return nil
}

func demoWhat() {
	fmt.Println("========== 01 error 是什么：接口 + 返回值模式 ==========")

	// ---------- ① 标准姿势：先判 err，再走正常逻辑（背下来，天天写）----------
	result, err := fetchData(true) // := 短声明：声明 + 赋值一步到位（≈ 前端 const + 赋值）
	if err != nil {
		fmt.Println("① 出错了：", err)
	} else {
		fmt.Println("① 成功拿到：", result)
	}

	// ---------- ② 失败时返回 (零值, error)；_ 表示"这个返回值我不要" ----------
	// 零值：字符串是 ""，数字是 0，指针/接口是 nil
	_, err = fetchData(false) // _ 空白标识符：丢掉第一个返回值，只留 err
	if err != nil {
		fmt.Println("② 失败分支 err 非 nil：", err)
	}

	// ---------- ③ error 可以比较：nil = 成功，非 nil = 失败 ----------
	if err := ping(); err == nil {
		fmt.Println("③ ping() 没出错（err == nil）✅")
	}

	// ---------- ④ error 是个接口：能装不同类型的错误 ----------
	// fmt.Println 遇到 error 会自动调它的 Error() 方法
	err = errors.New("一个普通错误")
	fmt.Println("④ 打印 error：", err)
	fmt.Println("④ err.Error() 拿到纯文本：", err.Error())

	fmt.Println("\n🔑 记住：Go 没有异常，错误是返回值；nil 是「没出错」，非 nil 是「出错」")
}
