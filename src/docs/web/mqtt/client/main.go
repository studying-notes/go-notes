/*
 * @Date: 2021.12.13 9:57
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2021.12.13 9:57
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://localhost:1883"
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("client")
	opts.SetKeepAlive(2 * time.Second)
	// Subscribe 的 callback 为 nil 时默认调用
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	if token := client.Subscribe("sample", 0, nil); token.Wait() && token.Error() != nil {
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
		token := client.Publish("sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(6 * time.Second)

	// 取消订阅
	if token := client.Unsubscribe("sample"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Disconnect(250)

	time.Sleep(1 * time.Second)
}
