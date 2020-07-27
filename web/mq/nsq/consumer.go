package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MyHandler struct {
	Title string
}

func (handler *MyHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("%s recv from %v, msg:%v\n",
		handler.Title, msg.NSQDAddress, string(msg.Body))
	return nil
}

func InitConsumer(topic string, channel string, address string) error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		panic(err)
	}
	consumer := &MyHandler{
		Title: "first",
	}
	c.AddHandler(consumer)
	return c.ConnectToNSQLookupd(address)
}

func main() {
	err := InitConsumer("topic_demo", "first", "127.0.0.1:4161")
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal)        // 定义一个信号的通道
	signal.Notify(c, syscall.SIGINT) // 转发键盘中断信号
	<-c                              // 阻塞
}
