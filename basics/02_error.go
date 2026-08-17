// 02_error.go —— 错误处理：Go 没有 try-catch，错误是返回值
//
// 前端对照一句话（先记这个）：
//
//	JS：throw + try/catch —— 错误被"抛出"，靠外层捕获
//	Go：函数多返回一个 error —— 错误被"返回"，调用方自己判断
//	好处：不会"忘了 catch 就悄悄崩"；代价：每个函数都要显式处理
//
// 为什么重要：后端业务代码里 if err != nil 能占 30%~50%。
//   面试官看 Go 代码，第一眼就看你错误处理规不规范。
//
// 本节一条主线：错误长什么样 → 造错误 → 包错误（带上下文）→ 判断错误链
package main

import (
	"errors"
	"fmt"
)

// ---------- 0. 哨兵错误（sentinel error）----------
// 在包级别定义"预设错误"，方便到处引用、统一判断（比随手写字符串靠谱）
var ErrInvalidID = errors.New("无效的 id")
var ErrNotFound = errors.New("找不到该记录")

// ---------- 1. 造错误 ----------

// divide 模拟一个可能失败的函数：除 0 返回错误
// Go 惯例：成功返回 (正常值, nil)，失败返回 (零值, 错误)
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为 0")
	}
	return a / b, nil
}

// getUser 模拟一层业务逻辑
// 用 fmt.Errorf + %w 把底层错误"包一层"：既带上上下文，又保留原始错误
func getUser(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("getUser: 非法 id=%d（%w）", id, ErrInvalidID)
	}
	if id == 99 {
		return "", fmt.Errorf("getUser: id=99 用户已被删除（%w）", ErrNotFound)
	}
	return "小明", nil
}

// ---------- 2. 自定义错误类型 ----------
// 除了错误信息，还能带结构化数据（Gin 参数校验的错误就是这种类型）
// 只要实现 Error() string 方法，就自动满足 error 接口（呼应 01 的隐式实现！）
type ValidateError struct {
	Field string
	Msg   string
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("字段 %s 校验失败：%s", e.Field, e.Msg)
}

// validateName 用自定义错误类型
func validateName(name string) error {
	if len(name) < 2 {
		return &ValidateError{Field: "name", Msg: "长度至少 2"}
	}
	return nil
}

func demoError() {
	fmt.Println("========== 02 错误处理：返回值模式 ==========")

	// ---------- ① 标准姿势：先判错，再走正常逻辑（背下来）----------
	r, err := divide(10, 2)
	if err != nil {
		fmt.Println("① 出错：", err)
	} else {
		fmt.Println("① 10/2 =", r)
	}

	// 失败分支：返回 (零值, 错误)
	if _, err := divide(1, 0); err != nil {
		fmt.Println("② 除 0：", err)
	}

	// ---------- ③ 错误包裹 + errors.Is 判断错误链 ----------
	// getUser 用 %w 包了一层，errors.Is 能穿透包裹链找到底层错误
	// 场景：错误被层层包过，你想判断"最里面是不是这个错"
	_, err = getUser(-1)
	fmt.Println("③ getUser(-1) 的错误信息：", err)
	if errors.Is(err, ErrInvalidID) {
		fmt.Println("③ errors.Is 判断：这确实是'无效 id'错误 ✅")
	}

	_, err = getUser(99)
	fmt.Println("③ getUser(99) 的错误信息：", err)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("③ errors.Is 判断：这确实是'找不到'错误 ✅")
	}

	// ---------- ④ errors.As：取出自定义错误类型里的结构化数据 ----------
	err = validateName("x")
	fmt.Println("④ 校验错误：", err)
	var ve *ValidateError
	if errors.As(err, &ve) { // As 沿错误链找"是不是这种类型"，找到就把数据塞进 ve
		fmt.Println("④ errors.As 取出：字段 =", ve.Field, "，信息 =", ve.Msg)
	}

	fmt.Println("\n🔑 错误处理三件套：err != nil 判断 / %w 包裹 / errors.Is·As 判断")
}
