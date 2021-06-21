package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	clientID = "xlj"
	wg       sync.WaitGroup
)

func main() {
	client := NewClient(clientID)
	if err := client.Connect(); err != nil {
		panic(err.Error())
	}

	wg.Add(1)

	// 接收消息
	go func() {
		err := client.Subscribe(func(c *Client, msg *Message) {
			fmt.Printf("----receive----\n%+v\n", msg)
			wg.Done()
		}, 1, topic)
		if err != nil {
			panic(err.Error())
		}
	}()

	// 发送消息
	msg := &Message{
		ClientID: clientID,
		Type:     "text",
		Data: "Every Linux or Unix command executed by " +
			"the shell script or user has an exit status.",
		Time: time.Now().Unix(),
	}

	data, _ := json.Marshal(msg)
	if err := client.Publish(topic, 1, false, data); err != nil {
		panic(err.Error())
	}

	wg.Wait()
}
