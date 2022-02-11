package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
)

func main() {

	gin.SetMode(gin.ReleaseMode)       // 关闭 debug 日志
	gin.DefaultWriter = ioutil.Discard // 丢弃全部日志

	r := gin.Default()
	// 接受任意类型、任意路径的请求，p 保存路径信息
	r.Any("/*p", func(c *gin.Context) {
		_, _ = io.Copy(os.Stdout, c.Request.Body) // 打印请求 body
		println()
		c.JSON(200, gin.H{
			"status":     "200",
			"error_code": "0",
			"error_msg":  "",
		})
	})
	_ = r.Run(":8080")
}
