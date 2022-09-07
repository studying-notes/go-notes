package main

import (
	"encoding/hex"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.SetTrustedProxies([]string{})

	r.GET("/file_attachment/ascii", func(c *gin.Context) {
		filename := make([]byte, 10)
		if _, err := rand.Read(filename); err == nil {
			c.FileAttachment("main.go", hex.EncodeToString(filename)+".go")
		} else {
			c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		}
	})

	r.GET("/file_attachment/utf-8", func(c *gin.Context) {
		c.FileAttachment("main.go", "中文名.go")
	})

	r.GET("/file_attachment/exe", func(c *gin.Context) {
		c.FileAttachment("Setup.exe", "Setup.exe")
	})

	r.Run(":26535")
}
