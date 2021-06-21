package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// 只发送消息，不接收消息
func main() {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("publisher")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 定时每隔几秒发送一次消息
	ticker := time.NewTicker(6 * time.Second)

	// 一共有两种类型进行测试
	for {
		person := Person{Name: "Wetness", Age: 18, Parents: []string{"Gomez", "Morita"}}
		if text, err := json.Marshal(Message{Type: "1", Content: person}); err != nil {
			panic(err)
		} else {
			token := c.Publish(topic, 0, false, text)
			token.Wait()
		}

		time.Sleep(3 * time.Second)

		thing := Thing{Length: 10, Width: 20, Height: 30}
		if text, err := json.Marshal(Message{Type: "2", Content: thing}); err != nil {
			panic(err)
		} else {
			token := c.Publish(topic, 0, false, text)
			token.Wait()
		}

		<-ticker.C
	}
}
