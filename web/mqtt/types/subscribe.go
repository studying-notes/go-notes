package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// 收到消息时的处理函数
var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	b := msg.Payload()
	//log.Printf("----receive----\n%s\n", b)

	model := make(map[string]interface{})
	_ = json.Unmarshal(b, &model)
	t := model["type"].(string)
	switch t {
	case "1":
		// 这种方法解析出来都是空接口
		//p := Message{Content: Person{}}
		//_ = json.Unmarshal(b, &p)
		//fmt.Printf("%#v\n\n", p.Content)

		p2 := struct {
			Type    string `json:"type"`
			Content Person `json:"content"`
		}{}
		_ = json.Unmarshal(b, &p2)
		fmt.Printf("%#v\n\n", p2.Content)

	case "2":
		//q := Message{Content: Thing{}}
		//_ = json.Unmarshal(b, &q)
		//fmt.Printf("%#v\n\n", q.Content)

		q2 := struct {
			Type    string `json:"type"`
			Content Thing `json:"content"`
		}{}
		_ = json.Unmarshal(b, &q2)
		fmt.Printf("%#v\n\n", q2.Content)
	}
}

// 只接收消息，不发送消息
// go run models.go publish.go
func main() {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("subscriber")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	if token := c.Subscribe(topic, 0, msgHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
	}
}
