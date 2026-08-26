// 04_validate.go —— validator 参数校验 + 统一错误响应
//
// TS 类型只在编译期管住自己，管不住运行时进来的网络请求
// 后端永远不能相信前端传参——校验是后端自己的事
// Gin 内置 go-playground/validator：规则写在 struct tag 上，ShouldBind 时自动跑
//
// 前端对照：运行时版的 zod / Element Plus 的 rules，只是规则写进 tag 里
package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


// apiResp 统一响应结构：{code, msg, data} 是行业惯例
// 前端对照：axios 响应拦截器按 code 统一拦截报错，消费的就是这个结构
type apiResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"` // omitempty：data 为空时整个字段不输出
}

// ok / fail 是全项目唯二的响应出口，结构才不会百花齐放
// 真实项目里它们一般抽到 pkg/response 包里
func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, apiResp{Code: 0, Msg: "ok", Data: data})
}

func fail(c *gin.Context, status int, msg string) {
	c.JSON(status, apiResp{Code: status, Msg: msg})
}


// registerInput 常用规则一览（min/max 在字符串和数字上含义不同，见下方注释）
type registerInput struct {
	// string 上：min/max 管长度
	Username string `json:"username" binding:"required,min=3,max=20"`
	// omitempty：没传就跳过后面的校验；email：格式校验
	Email string `json:"email" binding:"omitempty,email"`
	// oneof：枚举，空格分隔
	Role string `json:"role" binding:"required,oneof=admin editor viewer"`
	// string 上 min=8 是长度不是数值，别看错
	Password string `json:"password" binding:"required,min=8"`
	// int 上：gt/gte/lt/lte 管数值
	Age int `json:"age" binding:"required,gt=0"`
}


// createOrderInput 嵌套校验：dive = 钻进切片校验每个元素
// 没有 dive 只校验「切片非空」，元素内容不管
type createOrderInput struct {
	Items []orderItem `json:"items" binding:"required,dive"`
}

type orderItem struct {
	SKU string `json:"sku" binding:"required"`
	Qty int    `json:"qty" binding:"required,gt=0"`
}

func demoValidate() {
	fmt.Println("\n========== 04 参数校验 + 统一响应 ==========")

	router := gin.New()

	router.POST("/register", func(c *gin.Context) {
		var input registerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			bindFail(c, err)
			return
		}
		ok(c, gin.H{"username": input.Username})
	})

	router.POST("/orders", func(c *gin.Context) {
		var input createOrderInput
		if err := c.ShouldBindJSON(&input); err != nil {
			bindFail(c, err)
			return
		}
		ok(c, gin.H{"items": len(input.Items)})
	})

	//! required 认为「零值 = 没传」：age 传 0 会因 required 挂掉（0 是 int 的零值）
	//! 允许 0 的字段别用 required；要区分「没传」和「传了 0」就用指针 *int

	//? 校验为什么要放在 handler 最前面，而不是等业务代码报错？
	//? 校验失败是 400「客户端的错」；混进业务层就变成 500「服务端的错」，排查方向全反

	showGinRequest(router, http.MethodPost, "/register",
		`{"username":"saohua","role":"viewer","password":"12345678","age":25}`)
	showGinRequest(router, http.MethodPost, "/register",
		`{"username":"ab","role":"viewer","password":"12345678","age":25}`)
	showGinRequest(router, http.MethodPost, "/register",
		`{"username":"saohua","role":"root","password":"12345678","age":25}`)
	showGinRequest(router, http.MethodPost, "/register",
		`{"username":"saohua","role":"viewer","password":"12345678","age":0}`)
	showGinRequest(router, http.MethodPost, "/register",
		`{"username":"saohua","email":"not-an-email","role":"viewer","password":"12345678","age":25}`)
	showGinRequest(router, http.MethodPost, "/orders", `{"items":[{"sku":"A-1","qty":2}]}`)
	showGinRequest(router, http.MethodPost, "/orders", `{"items":[{"sku":"A-1","qty":0}]}`)

	fmt.Println("🔑 校验规则写在 tag 里，绑定即校验；错误出口全项目只有一个")
}


// bindFail 把绑定/校验错误翻译成人话
// validator 报的是英文结构化错误，直接丢给前端等于没报
func bindFail(c *gin.Context, err error) {
	// errors.As：判断错误链里是不是「校验错误」；JSON 格式错误走不进这个分支
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, translate(fe))
		}
		fail(c, http.StatusBadRequest, strings.Join(msgs, "；"))
		return
	}
	fail(c, http.StatusBadRequest, "参数格式错误："+err.Error())
}


// translate 单条规则翻成中文（真实项目做成「规则 → 文案」映射表，这里演示思路）
func translate(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " 不能为空"
	case "min", "max":
		// 同一个 tag 两副面孔：string 管长度，数字管大小——用 Kind 区分
		if fe.Kind() == reflect.String {
			return fe.Field() + " 长度不能" + map[string]string{"min": "小于", "max": "大于"}[fe.Tag()] + " " + fe.Param()
		}
		return fe.Field() + " 不能" + map[string]string{"min": "小于", "max": "大于"}[fe.Tag()] + " " + fe.Param()
	case "gt":
		return fe.Field() + " 必须大于 " + fe.Param()
	case "oneof":
		return fe.Field() + " 只能是 " + strings.Join(strings.Fields(fe.Param()), " / ")
	case "email":
		return fe.Field() + " 格式不正确"
	default:
		return fe.Field() + " 校验未通过：" + fe.Tag()
	}
}
