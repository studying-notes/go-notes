package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		cookie, _ := c.Cookie("key")
		c.SetCookie("key", "value", 3600,
			"/", "localhost", false, true)
		fmt.Println(cookie)
	})
	_ = r.Run()
}
