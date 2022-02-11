---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "NSQ 简介"  # 文章标题
url:  "posts/go/libraries/tripartite/mq/nsq"  # 设置网页永久链接
tags: [ "go", "mq", "nsq" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## NSQ 介绍

[NSQ](https://nsq.io/) 是 Go 语言编写的一个开源的实时分布式内存消息队列，其性能十分优异。 NSQ 的优势有以下优势：

1. NSQ 提倡分布式和分散的拓扑，没有单点故障，支持容错和高可用性，并提供可靠的消息交付保证；
2. NSQ 支持横向扩展，没有任何集中式代理；
3. NSQ 易于配置和部署，并且内置了管理界面。

## NSQ的应用场景

通常来说，消息队列都适用以下场景。

### 异步处理

参照下图利用消息队列把业务流程中的非关键流程异步化，从而显著降低业务请求的响应时间。

![](imgs/nsq1.png)

### 应用解耦

通过使用消息队列将不同的业务逻辑解耦，降低系统间的耦合，提高系统的健壮性。后续有其他业务要使用订单数据可直接订阅消息队列，提高系统的灵活性。

![](imgs/nsq2.png)

### 流量削峰

类似秒杀等场景下，某一时间可能会产生大量的请求，使用消息队列能够为后端处理请求提供一定的缓冲区，保证后端服务的稳定性。

![](imgs/nsq3.png)

## 下载安装

```
https://nsq.io/deployment/installing.html
```

## NSQ组件

### nsqlookupd

`nsqlookupd` 是维护所有 `nsqd` 状态、提供服务发现的守护进程。它能为消费者查找特定 `topic` 下的 `nsqd` 提供了运行时的自动发现服务。 它不维持持久状态，也不需要与任何其他 `nsqlookupd` 实例协调以满足查询。因此根据你系统的冗余要求尽可能多地部署 `nsqlookupd` 节点。它们消耗的资源很少，可以与其他服务共存。

```shell
# 启动
nsqlookupd
```

### nsqd

nsqd是一个守护进程，它接收、排队并向客户端发送消息。

启动 `nsqd`，指定 `-broadcast-address=127.0.0.1` 来配置广播地址。

```shell
nsqd --lookupd-tcp-address=127.0.0.1:4160
```

如果部署了多个 `nsqlookupd` 节点，可以指定多个 `-lookupd-tcp-address`。

### nsqadmin

一个实时监控集群状态、执行各种管理任务的 Web 管理平台。 启动 `nsqadmin`，指定 `nsqlookupd` 地址:

```shell
nsqadmin --lookupd-http-address=127.0.0.1:4161
```

## 官方示例

1. in one shell, start `nsqlookupd`:

   ```
   $ nsqlookupd
   ```

2. in another shell, start `nsqd`:

   ```
   $ nsqd --lookupd-tcp-address=127.0.0.1:4160
   ```

3. in another shell, start `nsqadmin`:

   ```
   $ nsqadmin --lookupd-http-address=127.0.0.1:4161
   ```

4. publish an initial message (creates the topic in the cluster, too):

   ```
   $ curl -d 'hello world 1' 'http://127.0.0.1:4151/pub?topic=test'
   ```

5. finally, in another shell, start `nsq_to_file`:

   ```
   $ nsq_to_file --topic=test --output-dir=/tmp --lookupd-http-address=127.0.0.1:4161
   ```

6. publish more messages to `nsqd`:

   ```
   $ curl -d 'hello world 2' 'http://127.0.0.1:4151/pub?topic=test'
   $ curl -d 'hello world 3' 'http://127.0.0.1:4151/pub?topic=test'
   ```

7. to verify things worked as expected, in a web browser open `http://127.0.0.1:4171/` to view the `nsqadmin` UI and see statistics. Also, check the contents of the log files (`test.*.log`) written to `/tmp`.

## NSQ 架构

### NSQ 工作模式

![](imgs/nsq4.png)


### Topic 和 Channel

每个 nsqd 实例旨在一次处理多个数据流。这些数据流称为 `“topics”`，一个 `topic` 具有 1 个或多个 `“channels”`。每个 `channel` 都会收到 `topic` 所有消息的副本，实际上下游的服务是通过对应的 `channel` 来消费 `topic` 消息。

`topic` 和 `channel` 不是预先配置的。`topic` 在首次使用时创建，方法是将其发布到指定 `topic`，或者订阅指定 `topic` 上的` channel`。`channel` 是通过订阅指定的 `channel` 在第一次使用时创建的。

`topic` 和 `channel` 都相互独立地缓冲数据，防止缓慢的消费者导致其他 `chennel` 的积压，同样适用于 `topic` 级别。

`channel` 可以并且通常会连接多个客户端。假设所有连接的客户端都处于准备接收消息的状态，则每条消息将被传递到随机客户端。例如：

![](imgs/nsq5.gif)

总之，消息是从 `topic -> channel`（每个 `channel` 接收该 `topic` 的所有消息的副本）多播的，但是从 `channel -> consumers` 均匀分布（每个消费者接收该 `channel` 的一部分消息）。

### NSQ 接收和发送消息流程

![](imgs/nsq6.png)

## NSQ 特性

- 消息默认不持久化，可以配置成持久化模式。NSQ 采用的方式是内存+硬盘的模式，当内存到达一定程度时就会将数据持久化到硬盘。
  - 如果将 `--mem-queue-size` 设置为 0，所有的消息将会存储到磁盘。
  - 服务器重启时也会将当时在内存中的消息持久化。
- 每条消息至少传递一次。
- 消息不保证有序。

## Go 操作 NSQ

```
go get -u github.com/nsqio/go-nsq
```

### 生产者

```go
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
```

### 消费者

```go
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
```
