package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

const (
	topic  = "types"
	broker = "tcp://localhost:1883"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Parents []string `json:"parents"`
}

type Thing struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func init() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0) // 打印调试信息
	mqtt.ERROR = log.New(os.Stdout, "", 0) // 打印错误信息
}
