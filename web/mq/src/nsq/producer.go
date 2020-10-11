package main

import (
	"bufio"
	"fmt"
	"github.com/nsqio/go-nsq"
	"os"
	"strings"
)

var producer *nsq.Producer

// 初始化生产者
func initProducer(s string) (err error) {
	cfg := nsq.NewConfig()
	producer, err = nsq.NewProducer(s, cfg)
	return err
}

func main() {
	nsqAddr := "127.0.0.1:4150"
	err := initProducer(nsqAddr)
	if err != nil {
		panic(err)
	}
	// 从标准输入读取
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read failed:", err)
			continue
		}
		data = strings.TrimSpace(data)
		if strings.EqualFold(data, "Q") {
			break
		}
		err = producer.Publish("topic_demo", []byte(data))
		if err != nil {
			fmt.Println("publish failed:", err)
			continue
		}
	}
}
