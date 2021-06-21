package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

// 简单处理，打印主题和信息
var handler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s\n\n", msg.Topic())
	fmt.Printf("Msg: %s\n", msg.Payload())
}

// 消息内容
type Content struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 消息结构
type Message struct {
	Type    int `json:"type"`
	Content `json:"content"`
}

func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0) // 打印调试信息
	mqtt.ERROR = log.New(os.Stdout, "", 0) // 打印错误信息

	broker := "tcp://localhost:1883"
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("client")
	opts.SetKeepAlive(2 * time.Second)
	// Subscribe 的 callback 为 nil 时默认调用
	opts.SetDefaultPublishHandler(handler)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	if token := c.Subscribe("sample", 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 发布消息
	for i := 0; i < 5; i++ {
		content := Content{Code: 200, Msg: fmt.Sprintf("This is msg #%d!", i)}
		text, err := json.Marshal(Message{Type: 0, Content: content})
		if err != nil {
			mqtt.ERROR.Println(err)
			continue
		}
		token := c.Publish("sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(6 * time.Second)

	// 取消订阅
	if token := c.Unsubscribe("sample"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
