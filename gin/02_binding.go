// 02_binding.go —— 参数绑定：query、路径参数、JSON body 三种来源
//
// net/http 阶段：query 手动 r.URL.Query().Get、body 手动 json.NewDecoder
// Gin 阶段：struct 挂 tag + ShouldBind 系列方法全搞定
//
// 绑定方法按「数据来源」选，一个 struct 可以同时挂多种 tag：
//
//	来源            方法               tag      前端对照
//	?query 参数     ShouldBindQuery    form     req.query
//	路径参数 :id    ShouldBindUri      uri      route.params
//	JSON body      ShouldBindJSON     json     req.body
//
// json tag 一身二职：请求体绑定的字段名 + 响应序列化的字段名
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


// createUserInput 演示一个字段挂多种 tag（binding 是校验规则，04 章细讲）
type createUserInput struct {
	Name string `json:"name" form:"name" uri:"name" binding:"required"`
	Age  int    `json:"age" form:"age" uri:"age" binding:"gte=0,lte=150"`
}

// listUserQuery 绑定 ?page=1&size=20 这种 query
// form 管「进来」，json 管「出去」——返回给前端的字段名也统一小驼峰
type listUserQuery struct {
	Page int    `form:"page" json:"page"`
	Size int    `form:"size" json:"size"`
	KW   string `form:"kw" json:"kw"`
}

// userIDURI 绑定路径里的 :id，顺手把「正整数」的校验也做了
type userIDURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func demoBinding() {
	fmt.Println("\n========== 02 参数绑定 ==========")

	router := gin.New()


	// ---- 取一两个零散参数：便捷方法最省事 ----
	router.GET("/search", func(c *gin.Context) {
		// Query：参数不存在时返回空串，不报错
		keyword := c.Query("kw")
		// DefaultQuery：没传时用默认值，分页场景常用
		page := c.DefaultQuery("page", "1")
		// GetQuery 返回 (值, 是否存在)：能区分「没传」和「传了空串」
		sort, exists := c.GetQuery("sort")
		c.JSON(http.StatusOK, gin.H{
			"kw": keyword, "page": page,
			"sort": sort, "sortExists": exists,
		})
	})


	// ---- 字段多时：struct 绑定 query ----
	router.GET("/users", func(c *gin.Context) {
		var q listUserQuery
		// 类型对不上（page=abc）会返回 error，而不是静默给 0——手写版做不到这点
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, q)
	})


	// ---- 路径参数：Uri 绑定 + 校验 ----
	router.GET("/users/:id", func(c *gin.Context) {
		var u userIDURI
		if err := c.ShouldBindUri(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID})
	})


	// ---- JSON body：用得最多的一个 ----
	router.POST("/users", func(c *gin.Context) {
		var input createUserInput
		// ShouldBindJSON = 读 body + 反序列化 + 触发 binding 校验
		// 前端对照：Express 里 body 已被中间件解析，这里连解析带校验一起做了
		if err := c.ShouldBindJSON(&input); err != nil {
			// ShouldBind 系列只返回 error、不写响应（Bind 才会自动写 400）
			// 错误怎么回给前端由自己掌控——生产代码统一错误处理就靠它
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// JSON 里没出现的字段保持零值，不报错
		c.JSON(http.StatusCreated, input)
	})


	//! body 是流，读一次就耗尽：同一请求里第二次 ShouldBindJSON 必然失败
	//! 表现为「代码一模一样、第二次绑定却报 EOF」的经典灵异现象
	//! 确实要多次绑定时用 c.ShouldBindBodyWith（它会缓存住 body）
	router.POST("/bind-twice", func(c *gin.Context) {
		var a, b createUserInput
		err1 := c.ShouldBindJSON(&a)
		err2 := c.ShouldBindJSON(&b)
		fmt.Printf("   第一次绑定 err=%v；第二次绑定 err=%v\n", err1, err2)
		c.JSON(http.StatusOK, gin.H{"a": a, "b": b})
	})

	//? age 传 0 和不传 age，绑定后 Age 都是 0，怎么区分？
	//? 字段改指针 *int：nil = 没传，指向 0 = 真传了 0（PATCH 局部更新就靠这招）

	showGinRequest(router, http.MethodGet, "/search?kw=go&page=2", "")
	showGinRequest(router, http.MethodGet, "/search?kw=go&sort=", "")
	showGinRequest(router, http.MethodGet, "/search", "")
	showGinRequest(router, http.MethodGet, "/users?page=2&size=10", "")
	showGinRequest(router, http.MethodGet, "/users?page=abc", "")
	showGinRequest(router, http.MethodGet, "/users/42", "")
	showGinRequest(router, http.MethodGet, "/users/-1", "")
	showGinRequest(router, http.MethodPost, "/users", `{"name":"saohua","age":25}`)
	showGinRequest(router, http.MethodPost, "/users", `{"age":25}`)
	showGinRequest(router, http.MethodPost, "/bind-twice", `{"name":"x"}`)

	fmt.Println("🔑 三种来源三个方法：ShouldBindQuery / ShouldBindUri / ShouldBindJSON")
}
