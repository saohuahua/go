// 03_middleware.go —— Gin 中间件：三种挂法、Next/Abort、Set/Get 传值
//
// net_http/04 手写过 middleware(next http.Handler) http.Handler 再倒序组装
// Gin 把中间件和 handler 排进同一个数组，c.Next() 负责往下走
// 洋葱模型没变，只是不用自己组装了
//
// 三种挂法：
//
//	router.Use(m)            全局：每个请求都穿
//	group.Use(m)             组级：只穿这组路由
//	router.GET("/x", m, h)   单路由：夹在 path 和 handler 之间
//
// 前端对照：Axios 拦截器 / Express 的 app.use / Vue 路由守卫
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func demoMiddleware() {
	fmt.Println("\n========== 03 中间件 ==========")

	router := gin.New()

	// 全局中间件：先 Use 的在最外层、先进后出（和 net_http 倒序组装的结果一致）
	router.Use(timerMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	// 单路由中间件：只保护这一条路由
	// 前端对照：某个路由单独加的 beforeEach 守卫
	router.GET("/secret", gateMiddleware("123"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "机密数据"})
	})

	// 组级中间件：整组路由统一鉴权，组内 handler 不用重复写
	api := router.Group("/api")
	api.Use(fakeAuthMiddleware())
	{
		api.GET("/me", func(c *gin.Context) {
			// Get 返回 any 要类型断言；GetInt/GetString 是内置快捷方法
			// 对照 net_http/04：context.WithValue + 自定义 key + 断言，三步并成一步
			c.JSON(http.StatusOK, gin.H{"userID": c.GetInt("userID")})
		})
	}

	showGinRequest(router, http.MethodGet, "/ping", "")
	showGinRequest(router, http.MethodGet, "/secret", "")
	showGinRequest(router, http.MethodGet, "/secret?password=123", "")
	showGinRequest(router, http.MethodGet, "/api/me", "")
	showGinRequest(router, http.MethodGet, "/api/me?token=abc", "")

	fmt.Println("🔑 Use 的顺序 = 洋葱的层级；Abort 拦截后必须 return")
}


// timerMiddleware 计时中间件：演示 Next 前后是「进」「出」两个世界
func timerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// c.Next() 放行：执行后面的中间件和 handler，全部跑完才回到这一行
		// Next 之前 = 请求还没处理；Next 之后 = 响应已经写出去了（改状态码无效）
		c.Next()
		// 这里的耗时是「整条链」的耗时，不是 handler 一家
		fmt.Printf("   [timer] %s %s 耗时 %s\n",
			c.Request.Method, c.Request.URL.Path, time.Since(start).Round(time.Microsecond))
	}
}


// gateMiddleware 拦截演示：Abort + return 缺一不可
func gateMiddleware(password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("password") != password {
			// AbortWithStatusJSON = 写 401 响应 + 打「后续 handler 别再执行」的标记
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "密码不对"})
			//! Abort 只阻止「后续中间件和 handler」，不会跳出当前函数
			//! 这行 return 不写，本函数里后面的代码照样执行——中间件一长，漏 return 就是事故
			return
		}
		// 校验通过可以不写 c.Next()：中间件返回后，引擎自动执行链上的下一个
		// 但只有显式调用 c.Next() 才拥有「响应之后」的世界（洋葱的出口半段）
	}
}


// fakeAuthMiddleware 演示 Set/Get：中间件鉴定身份，handler 直接取用
func fakeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("token") != "abc" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token 无效"})
			return
		}
		// Set 塞进「本次请求」的上下文，链路上任何 handler 都能 Get 到
		// 并发请求各自有各自的 c，不会串数据（对照 net_http/04 用 ctx 传 requestID）
		c.Set("userID", 42)
	}
}
