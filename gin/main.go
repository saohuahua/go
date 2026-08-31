// main.go —— Gin：路由、绑定、中间件、校验、JWT、TODO API
//
// 本文件夹共 6 个演示文件 + 本入口文件
//
// 01_router.go     —— 路由注册、路径参数、路由分组
// 02_binding.go    —— query / 路径 / JSON 三种参数绑定
// 03_middleware.go —— 中间件三种挂法、Next/Abort、Set/Get
// 04_validate.go   —— validator 校验、统一错误响应
// 05_jwt.go        —— JWT 签发与校验、鉴权中间件
// 06_todo_api.go   —— 用 Gin 重写内存 TODO API
//
// 运行方法
//
// cd gin && go run .
// go run ./gin    （在项目根目录执行）
//
// 学习方法
//
// 1 先完整跑一遍看输出，再按编号对照注释精读
// 2 06 和 net_http/05 对照着读，差异就是 Gin 的价值
// 3 每个 //! 的坑都亲手解开试一次，再注释回去
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func gintest() {
	e := gin.Default()
	e.GET("/findUser/:uid/:name", findUser)
	// e.GET("/downloadFile/*filePath",UserPage)

	log.Fatalln(e.Run(":8080"))
}

func findUser(c *gin.Context) {
	uid := c.Param("uid")
	username := c.Param("username")
	c.String(http.StatusOK, "username is %s\n userid is %s", username, uid)
}

func main() {
	// TestMode：关掉 Gin 的 debug 启动打印，demo 输出干净
	// 真实项目设 gin.ReleaseMode：不打 debug 日志，性能也更好
	// gin.SetMode(gin.TestMode)

	// demoRouter()     // 01 路由
	// demoBinding()    // 02 参数绑定
	// demoMiddleware() // 03 中间件
	// demoValidate()   // 04 校验 + 统一响应
	// demoJWT()        // 05 JWT
	// demoTodoAPI()    // 06 TODO API

	fmt.Println("\n🎉 Gin 速通完成：下一步 GORM + MySQL，或直接搭分层项目骨架")
}
