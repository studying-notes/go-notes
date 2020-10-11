package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

const (
	topic    = "sample"
	broker   = "tcp://localhost:1883"
	user     = "xlj"
	password = "xlj"
)

type Message struct {
	ClientID string `json:"client_id"`
	Type     string `json:"type"`
	Data     string `json:"data,omitempty"`
	Time     int64  `json:"time"`
}

type Client struct {
	client mqtt.Client
	opt    *mqtt.ClientOptions
	locker *sync.Mutex
	// 消息收到之后处理函数
	observer func(c *Client, msg *Message)
}

func NewClient(clientID string) *Client {
	opt := mqtt.NewClientOptions().
		AddBroker(broker).
		SetUsername(user).
		SetPassword(password).
		SetClientID(clientID).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(120 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetWriteTimeout(10 * time.Second).
		SetOnConnectHandler(func(c mqtt.Client) {
			// 连接被建立后的回调函数
			fmt.Println("MQTT is connected. Client ID:", clientID)
		}).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			// 连接被关闭后的回调函数
			fmt.Println("MQTT is disconnected. Client ID:", clientID)
			fmt.Println("Closed reason:", err.Error())
		})
	return &Client{
		client: mqtt.NewClient(opt),
		opt:    opt,
		locker: &sync.Mutex{},
	}
}

func (c *Client) GetClientID() string {
	return c.opt.ClientID
}

// 确保连接
func (c *Client) Connect() error {
	if !c.client.IsConnected() {
		c.locker.Lock()
		defer c.locker.Unlock()
		if !c.client.IsConnected() { // 考虑并发情况再判断一次
			if token := c.client.Connect(); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}
	return nil
}

// 消费消息
func (c *Client) Subscribe(observer func(c *Client, msg *Message), qos byte, topics ...string) error {
	if len(topics) == 0 {
		return errors.New("the topic is empty")
	}
	if observer == nil {
		return errors.New("the observer func is nil")
	}
	if c.observer != nil {
		return errors.New("an existing observer subscribed on this client, " +
			"you must unsubscribe it before you subscribe a new observer")
	}
	c.observer = observer
	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = qos
	}
	c.client.SubscribeMultiple(filters, c.messageHandler)
	return nil
}

// 发布消息
func (c *Client) Publish(topic string, qos byte, retained bool, data []byte) error {
	if err := c.Connect(); err != nil {
		return err
	}
	// retained: 是否保留信息
	token := c.client.Publish(topic, qos, retained, data)
	if err := token.Error(); err != nil {
		return err
	}

	if !token.WaitTimeout(time.Second * 10) {
		return errors.New("publish wait timeout")
	}

	return nil
}

func (c *Client) messageHandler(cli mqtt.Client, msg mqtt.Message) {
	if c.observer == nil {
		return
	}
	if message, err := MsgDecoder(msg.Payload()); err != nil {
		log.Println(err.Error())
		return
	} else {
		c.observer(c, message)
	}
}

func (c *Client) Unsubscribe(topics ...string) {
	c.observer = nil
	c.client.Unsubscribe(topics...)
}

// 解析 JSON 数据，与官方的区别在于数值类型解析
func MsgDecoder(payload []byte) (*Message, error) {
	msg := &Message{}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	// as a interface{} instead of as a float64
	decoder.UseNumber() // 方便之后类型转换
	if err := decoder.Decode(&msg); err != nil {
		return nil, err
	}
	return msg, nil
}
