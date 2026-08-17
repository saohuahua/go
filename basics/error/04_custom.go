// 04_custom.go —— 自定义错误类型：错误里带结构化数据（Gin 校验错误同款）
//
// errors.New / fmt.Errorf 只能带"一句话文本"。有些场景你需要错误里携带数据：
//   错误码、HTTP 状态码、是哪个字段校验失败、给前端看的提示……
//   这时候自定义一个类型，让它实现 Error() string —— 它就是 error 了（隐式实现）。
//
// 前端对照：
//	JS：throw new ValidationError({ field, msg }) —— 自定义错误类 + instanceof 判断
//	Go：定义 struct + Error() 方法 = 自定义错误；errors.As 把类型从链里掏出来
package main

import (
	"errors"
	"fmt"
)

// AppError 业务错误：错误里能带"给上层用的数据"
// 字段贴近后端：HTTP 状态码（决定响应 4xx/5xx）+ 业务错误码 + 给前端的话术
type AppError struct {
	HTTPStatus int    // 这个错误该响应成哪个 HTTP 状态码
	Code       string // 业务错误码，前端可以拿去展示
	Message    string // 人类可读的提示
}

// Error() string：实现它，*AppError 就是 error（隐式实现，呼应接口专题）
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// fetchUser 底层返回 *AppError。
// 注意用 %w 包：这样上层才能用 errors.As 把它从链里掏出来。
func fetchUser(id int) error {
	if id < 0 {
		return fmt.Errorf("fetchUser: 参数校验失败（id=%d）: %w", id, &AppError{
			HTTPStatus: 400,
			Code:       "INVALID_PARAM",
			Message:    "用户 id 不能为负数",
		})
	}
	if id == 0 {
		return fmt.Errorf("fetchUser: 用户不存在: %w", &AppError{
			HTTPStatus: 404,
			Code:       "USER_NOT_FOUND",
			Message:    "查无此人",
		})
	}
	return nil
}

// handleGetUser 模拟 Gin 的 handler：把 error 翻译成 HTTP 响应。
// 不真起服务，只把"该返回什么状态码"打印出来示意。
func handleGetUser(id int) {
	err := fetchUser(id)
	if err == nil {
		fmt.Printf("✅ GET /user/%d → 200，返回用户数据\n", id)
		return
	}
	// 默认按 500；如果错误链里有 *AppError，就听它的状态码
	status := 500
	var appErr *AppError
	if errors.As(err, &appErr) { // &appErr = "口袋"的地址，As 找到类型后把值装进去
		status = appErr.HTTPStatus
		fmt.Printf("❌ GET /user/%d → %d，错误码 %s，提示「%s」\n", id, status, appErr.Code, appErr.Message)
		return
	}
	fmt.Printf("❌ GET /user/%d → %d 内部错误：%v\n", id, status, err)
}

func demoCustom() {
	fmt.Println("========== 04 自定义错误类型：错误里带结构化数据 ==========")

	// ---------- ① 直接打印：走的是 Error() 方法 ----------
	err := fetchUser(-1)
	fmt.Println("① fetchUser(-1) 的错误：", err)

	// ---------- ② errors.As 从链里掏出 AppError，拿到里面的数据 ----------
	// &AppError{...} 里的 & 是"取地址"：因为 Error() 挂在指针接收者上，*AppError 才实现 error
	var appErr *AppError
	if errors.As(err, &appErr) {
		fmt.Printf("② 掏出 AppError：HTTPStatus=%d，Code=%s，Message=%s\n",
			appErr.HTTPStatus, appErr.Code, appErr.Message)
	}

	// ---------- ③ 实战：handler 把错误翻译成 HTTP 响应 ----------
	handleGetUser(123) // 正常
	handleGetUser(-1)  // 参数错 → 400
	handleGetUser(0)   // 找不到 → 404

	fmt.Println("\n🔑 想带数据 → 自定义类型；想取数据 → errors.As；handler 拿到数据决定响应")
}
