package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode("debug")

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {

		c.JSON(http.StatusOK,
			gin.H{
				"girls": "cute",
			},
		)
	})

	_ = r.Run(":18888")
}
