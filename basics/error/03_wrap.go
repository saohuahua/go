// 03_wrap.go —— 错误链：%w 包裹 + errors.Is / errors.As（面试高频，重点）
//
// 场景：错误从底层一路传上来，每一层都想"加一句说明"（我在哪层、参数是啥）。
//
//	用 fmt.Errorf + %w 包一层：既能带上上下文，又能把底层错误留在"链"里。
//	之后 errors.Is / errors.As 就能穿透这串链条，判断"最里面是不是这个错"。
//
// 前端对照：
//
//	JS：只能 instanceof 检查最外层错误，内层要靠 err.cause 手动挖
//	Go：errors.Is / errors.As 自动穿透包裹链 —— Go 1.13 之后的核心能力，面试必问
package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ---------- 模拟三层业务：repo → service → handler ----------
// 哨兵错误定义在 02_make.go，这里直接用 —— 哨兵的意义就是全局共用、统一判断

// userRepo 最底层：只返回哨兵错误，不解释
func userRepo(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	if id == 999 {
		return ErrNotFound
	}
	return nil
}

// userService 中间层：用 %w 包一层，加上"我在哪层、参数是啥"
func userService(id int) error {
	if err := userRepo(id); err != nil {
		return fmt.Errorf("userService: 查询用户失败（id=%d）: %w", id, err)
	}
	return nil
}

// userHandler 最外层：再包一层，handler 是最接近用户的层
func userHandler(id int) error {
	if err := userService(id); err != nil {
		return fmt.Errorf("userHandler: 获取用户失败: %w", err)
	}
	return nil
}

// wrapWithV / wrapWithW：对比 %v 和 %w 的区别（见 demo ③）
func wrapWithV() error {
	return fmt.Errorf("用 %%v 包：%v", ErrNotFound) // %v 只是把文本拼进去，底层错误丢了
}

func wrapWithW() error {
	return fmt.Errorf("用 %%w 包：%w", ErrNotFound) // %w 把底层错误留在错误链里
}

func demoWrap() {
	fmt.Println("========== 03 错误链：%w 包裹 + errors.Is / errors.As ==========")

	// ---------- ① 链式错误：每一层都加一句上下文，一眼看出错在哪一层 ----------
	err := userHandler(999)
	fmt.Println("① 完整错误信息（层层带上下文）：", err)

	// ---------- ② errors.Is：穿透包裹链，判断"最里面是不是这个错" ----------
	// Is 顺着 %w 的链一路往下找，找到 == 哨兵的错误就返回 true
	if errors.Is(err, ErrNotFound) {
		fmt.Println("② errors.Is(err, ErrNotFound) = true ✅ 隔着 2 层也认得出")
	}
	// 判断错了目标 → false（认错对象，什么都不打印）
	if errors.Is(err, ErrInvalidID) {
		fmt.Println("② 认错目标（不会打印）")
	}

	// ---------- ③ %v vs %w：一个丢错误，一个留错误 ----------
	// %v：把文本拼进新错误，底层错误彻底没了，Is 穿透不了
	// （Printf 里的 %%v 是转义，打印出字面量 %v；真正的 %v 负责格式化后面的参数）
	fmt.Printf("③ %%v 包的：%v → Is 能认吗？%v\n", wrapWithV(), errors.Is(wrapWithV(), ErrNotFound))
	// %w：底层错误留在链里，Is 能一路穿透找到
	fmt.Printf("③ %%w 包的：%v → Is 能认吗？%v\n", wrapWithW(), errors.Is(wrapWithW(), ErrNotFound))

	// ---------- ④ errors.As：找"类型"而不是找"值" ----------
	// 标准库例子：strconv.ParseInt 解析失败返回 *strconv.NumError，As 能把它掏出来
	_, parseErr := strconv.ParseInt("abc", 10, 64)
	var numErr *strconv.NumError
	// 注意：第二个参数传 &numErr（指针的指针）—— 它是个"口袋"，As 找到类型后把值装进去
	if errors.As(parseErr, &numErr) {
		fmt.Printf("④ 掏出 *strconv.NumError：函数=%s 输入=%q\n", numErr.Func, numErr.Num)
	}
	// 区别一句话：Is 判断"是不是某个错误值"，As 判断"是不是某种错误类型"（能取数据）
	// 自定义错误类型 + As 取结构化数据，04 专门演示
	fmt.Println("④ errors.As 找类型、errors.Is 找值 —— 两个都能穿透包裹链")

	fmt.Println("\n🔑 只要路径上用了 %w，Is/As 就能穿透整条链 —— 这是 Go 1.13 的核心")
}
