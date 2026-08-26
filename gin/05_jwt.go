// 05_jwt.go —— JWT 签发与校验：登录发 token，中间件保护路由
//
// JWT = 三段字符串 header.payload.signature（用 . 连接）
//   payload 放用户信息（uid、过期时间），signature 用服务端密钥对前两段签名
//   改任何一个字节签名就对不上——这就是「防篡改」
// 服务端不存 token，收到后用密钥重算签名即可验证，这就是「无状态认证」
//
// 前端对照：登录后 token 存 localStorage/pinia，axios 请求拦截器统一加
//   Authorization: Bearer <token> 头——本章写的中间件就是这行代码背后服务端的事
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret HS256 是对称签名：签发和校验用同一把密钥
// ⚠️ 生产环境绝不能硬编码：从环境变量/配置读，泄露 = 任何人都能伪造 token
var jwtSecret = []byte("demo-secret-不安全-仅演示")


// userClaims 自定义 claims = 标准字段（嵌入 RegisteredClaims）+ 业务字段
type userClaims struct {
	UserID   int    `json:"uid"`
	Username string `json:"username"`
	// exp（过期）、iat（签发时间）、sub（主体）等标准字段都在这个嵌入结构里
	jwt.RegisteredClaims
}


// signToken 签发：claims → token 字符串
func signToken(userID int, username string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := userClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprintf("%d", userID),
			// 过期时间必设：不设的 token 永久有效，泄露后无法挽回
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	// HS256 = HMAC + SHA256，最常用的对称签名算法
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}


// parseToken 校验：token 字符串 → claims
func parseToken(tokenString string) (*userClaims, error) {
	var claims userClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		//? 为什么这里要检查签名方法，直接返回密钥不行吗？
		//? 防算法混淆攻击：攻击者把 header 里的 alg 改成别的算法，骗过不检查的解析器
		//? 只认 HS256 且密钥对得上，两道都过才算数
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不认识的签名算法 %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		// 签名不对 / 已过期 / 格式坏，全走这里（判断过期用 errors.Is(err, jwt.ErrTokenExpired)）
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token 无效")
	}
	return &claims, nil
}


// jwtAuthMiddleware 从 Authorization: Bearer <token> 校验身份
func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少 Bearer token"})
			return
		}
		claims, err := parseToken(strings.TrimPrefix(auth, "Bearer "))
		if err != nil {
			//! 401 = 没登录/token 无效；403 = 登录了但没权限，别混用
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token 无效：" + err.Error()})
			return
		}
		// 校验通过：身份塞进请求上下文，后面的 handler 直接取
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
	}
}

func demoJWT() {
	fmt.Println("\n========== 05 JWT ==========")

	//! jwtSecret 硬编码只是 demo；生产环境从配置读，正确写法：
	//! secret := []byte(os.Getenv("JWT_SECRET"))

	router := gin.New()

	// 登录：真实项目里先查库比对密码哈希，通过才签发；demo 直接发
	router.POST("/login", func(c *gin.Context) {
		token, err := signToken(42, "saohua", 15*time.Minute)
		if err != nil {
			fail(c, http.StatusInternalServerError, "签发失败")
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// 受保护路由：中间件挡在前面，handler 只管用身份
	router.GET("/me", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"userID": c.GetInt("userID"), "username": c.GetString("username")})
	})


	// 走一遍真实流程：登录拿 token → 带 token 访问
	loginBody := showJWTRequest(router, http.MethodPost, "/login", "", "")
	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal([]byte(loginBody), &loginResp)
	token := loginResp.Token

	// JWT 是三段式，用 . 切开：header（算法）. payload（claims）. signature（签名）
	parts := strings.Split(token, ".")
	fmt.Printf("   token 三段式：header=%s… payload=%s… signature=%s…\n",
		parts[0][:12], parts[1][:12], parts[2][:12])

	// 专门签发一个「一小时前就过期」的 token，演示过期场景
	expired, _ := signToken(42, "saohua", -time.Hour)

	//? access token 15 分钟就过期，用户正用着突然 401 怎么办？
	//? 双 token：access（短期，干活用）+ refresh（长期，只用来换新 access）
	//? access 过期拿 refresh 换新的 = 无感刷新；refresh 也过期才重新登录

	showJWTRequest(router, http.MethodGet, "/me", "", token)
	showJWTRequest(router, http.MethodGet, "/me", "", "")
	showJWTRequest(router, http.MethodGet, "/me", "", "伪造的token")
	showJWTRequest(router, http.MethodGet, "/me", "", expired)

	fmt.Println("🔑 JWT 无状态：服务端不存会话，验签即验身份；密钥就是一切")
}


// showJWTRequest 带鉴权头的内存请求（比 showGinRequest 多塞一个 Authorization）
// 返回响应 body，供调用方继续解析（demo 里用它取登录 token）
func showJWTRequest(router *gin.Engine, method, path, body, token string) string {
	req := httptest.NewRequest(method, "http://example.com"+path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// JWT token 很长，截断展示（按 rune 切，防止把中文切出乱码）
	text := strings.TrimSpace(rec.Body.String())
	if r := []rune(text); len(r) > 64 {
		text = string(r[:64]) + "…"
	}
	fmt.Printf("   %s %-8s → %d %s\n", method, path, rec.Code, text)
	return rec.Body.String()
}
