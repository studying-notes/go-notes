package main

import "github.com/gin-gonic/gin"

func pingRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	r := pingRouter()
	r.Run()
}
