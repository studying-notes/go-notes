package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.JSONFormatter{}
	f, _ := os.Create("gin.log")
	log.Out = f
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Out
	log.Level = logrus.InfoLevel
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn("A group of walrus emerges from the ocean.")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	_ = r.Run()
}
