package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserInfo struct {
	User string
	Pass string
}

func authHandler(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 2001,
			"msg":    "错误参数",
		})
		return
	}

	if user.User == "user" && user.Pass == "pass" {
		token, _ := GenToken(user.User)
		c.JSON(http.StatusOK, gin.H{
			"status": 2000,
			"msg":    "success",
			"data":   gin.H{"token": token},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 2002,
		"msg":    "鉴权失败",
	})
	return
}

// 认证中间件，中间件完全可以写成路由处理函数
func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": 2003,
				"msg":    "no token",
			})
			c.Abort() // 确保该请求的其余处理程序不被调用，但不停止当前处理程序
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"status": 2004,
				"msg":    "格式错误",
			})
			c.Abort()
			return
		}
		ch, err := ParseToke(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 2003,
				"msg":    "invalid token",
			})
			c.Abort() // 确保该请求的其余处理程序不被调用，但不停止当前处理程序
			return
		}

		// 将当前请求的 user 信息保存到请求的上下文上
		c.Set("user", ch.User)

		// 后续的处理函数可以用过 c.Get("user") 来获取当前请求的用户信息
		//c.Next()
	}
}

func main() {
	r := gin.Default()
	r.GET("/", AuthMiddleware(), func(c *gin.Context) {
		user := c.MustGet("user").(string)
		c.JSON(http.StatusOK, gin.H{
			"status": 2000,
			"msg":    "success",
			"data":   gin.H{"user": user},
		})
	})
	r.POST("/auth", authHandler)
	r.Run()
}
