// 01_router.go —— Gin 路由：注册、路径参数、分组、404
//
// net/http 阶段用 switch 手写「method + path 分流」，Gin 把它变成一行声明式注册
// gin.Engine 实现了 http.Handler（有 ServeHTTP 方法），本体 = 路由表 + 中间件链
//
// 前端对照：router.GET("/users", h) ≈ Express 的 app.get("/users", h)
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func demoRouter() {
	fmt.Println("\n========== 01 Gin 路由 ==========")

	// gin.New() 返回零中间件的纯净 Engine，学习期输出干净
	// gin.Default() = gin.New() + Logger（打请求日志）+ Recovery（panic 兜底返回 500）
	// 真实项目直接用 Default：一个 handler panic 不至于打崩整个进程
	router := gin.New()

	// ---- 基本注册：method + path + handler ----

	// gin.HandlerFunc 签名固定：func(c *gin.Context)
	// c 是这次请求的全部上下文：读参数、写响应、传值都找它
	// 前端对照：c ≈ Express 的 (req, res) 合体
	router.GET("/health", func(c *gin.Context) {
		// c.JSON = marshal + 设 Content-Type + 写 body，net_http 里这三步是手写的
		// gin.H 就是 map[string]any，写小 JSON 不用先声明 struct
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ---- 路径参数：把 URL 中的一段变成变量 ----

	// :id 只匹配「一段」（两个斜杠之间的内容），Gin 里最常用的参数写法
	router.GET("/users/:id", func(c *gin.Context) {
		// Param 返回值永远是 string，要数字得自己 strconv.Atoi
		// 前端对照：Vue Router 的 /users/:id + route.params.id
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	})

	// *rest 匹配「剩余整段路径」（可以含斜杠），只能出现在路径末尾
	router.GET("/files/*rest", func(c *gin.Context) {
		// 取出来的值带开头的斜杠：/files/a/b/c → *rest = "/a/b/c"
		c.JSON(http.StatusOK, gin.H{"rest": c.Param("rest")})
	})

	// ---- 路由分组：公共前缀 + 组级中间件 ----

	// 按版本或模块分组；组上还能 .Use() 挂中间件（03 章演示）
	// 前端对照：Vue Router 的嵌套路由、按模块拆分的 router 文件
	api := router.Group("/api/v1")
	{
		// {} 纯粹是视觉分组，删掉照样能跑；写了组内路由一眼可辨
		api.GET("/todos", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"count": 2}) })
		api.POST("/todos", func(c *gin.Context) { c.Status(http.StatusCreated) })
	}

	// ---- 404 兜底：对应 net_http 手写版里的 default 分支 ----
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "路由不存在"})
	})

	//? GET /users/42/posts 会被 /users/:id 匹配吗？
	//? 不会。:id 只匹配一段、不跨斜杠；子路径要单独注册 /users/:id/posts

	//! 同一位置的参数名必须一致：/users/:id 和 /users/:name 冲突，注册时直接 panic
	//! router.GET("/users/:name", func(c *gin.Context) {})
	//! 解开上面一行再运行，亲眼看一次 panic，然后注释回去

	showGinRequest(router, http.MethodGet, "/health", "")
	showGinRequest(router, http.MethodGet, "/users/42", "")
	showGinRequest(router, http.MethodGet, "/users/42/posts", "")
	showGinRequest(router, http.MethodGet, "/files/a/b/c", "")
	showGinRequest(router, http.MethodGet, "/api/v1/todos", "")
	showGinRequest(router, http.MethodPost, "/api/v1/todos", "")
	showGinRequest(router, http.MethodGet, "/nope", "")

	fmt.Println("🔑 声明式路由：注册即生效；:id 匹配一段，*rest 匹配剩余全部")
}

// showGinRequest 给 router 发一条内存请求并打印结果，全文件夹的 demo 都复用它
// 不占端口、不阻塞 main（Engine 本身是 http.Handler，直接喂给 httptest 即可）
// 前端对照：等于 mock fetch，省得真的起服务器再 curl
func showGinRequest(router *gin.Engine, method, path, body string) {
	// body 为空时传 NewReader("") 也完全可用，不需要特判
	req := httptest.NewRequest(method, "http://example.com"+path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	fmt.Printf("   %s %-22s → %d %s\n", method, path, rec.Code, strings.TrimSpace(rec.Body.String()))
}
