---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "MQTT 学习笔记"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/web/mqtt/intro"  # 设置网页永久链接
tags: [ "go", "mqtt", "intro" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 客户端

```
$ go get github.com/eclipse/paho.mqtt.golang
```

## 服务端

### Linux - Mosquitto

#### 安装

1. 引入仓库

```shell
$ apt-add-repository ppa:mosquitto-dev/mosquitto-ppa
$ apt-get update
```

2. 安装 Mosquitto 服务端

```shell
$ apt-get install mosquitto
# 开发包
$ apt-get install mosquitto-dev
```

3. 安装 Mosquitto 客户端

```shell
$ apt-get install mosquitto-clients
```

4. 启动服务端

```shell
$ /etc/init.d/mosquitto start
$ service mosquitto start

# 关闭
$ /etc/init.d/mosquitto stop
$ service mosquitto stop
```

这种方法看不到输出。为了显示控制信息，可以手动启动服务：

```shell
$ mosquitto -v
```

5. 查询 Mosquitto 是否正确运行

```shell
$ service mosquitto status
* mosquitto is running
```

#### 测试

1. 注册一个 topic 进行接收

```shell
$ mosquitto_sub -h localhost -t "mqtt" -v
```

2. 发布消息到刚注册的 topic

```shell
$ mosquitto_pub -h localhost -t "mqtt" -m "Hi"
```

3. 自定义设置

`/etc/mosquitto/mosquitto.conf`

```ini
# 关闭匿名
allow_anonymous false

# 设置用户名和密码
password_file /etc/mosquitto/pwfile

# 配置访问控制列表
acl_file /etc/mosquitto/acl

# 配置端口
port 8000
```

4. 添加用户

```shell
$ mosquitto_passwd -c /etc/mosquitto/pwfile xlj
```

5. 添加 Topic 和用户的关系

`/etc/mosquitto/acl`

```ini
# 用户 test 只能发布以 test 为前缀的主题
# 订阅以 mqtt 开头的主题
user test
topic write test/#
topic read mqtt/#

user pibigstar
topic write mqtt/#
topic write mqtt/#
```

6. 重启服务

```shell
# 这种方式可能存在问题
$ /etc/init.d/mosquitto restart
# 建议手动关闭启动
```

7. 监听消费

```shell
$ mosquitto_sub -h localhost -t "mqtt" -v -u xlj -P xlj
```

8. 发布消息

```shell
$ mosquitto_pub -h localhost -t "mqtt" -m "Hi" -u xlj -P xlj
$ mosquitto_pub -h localhost -t "sample" -m '{"client_id":"xlj","type":"text","data":"Every Linux or Unix command executed by the shell script or user has an exit status.","time":1596812596}' -u xlj -P xlj
```

### Windows - HMQ

```
$ git clone https://github.com/fhmq/hmq.git
$ cd hmq
$ go run main.go
```

#### 用法

```
Usage: hmq [options]

Broker Options:
    -w,  --worker <number>            Worker num to process message, perfer (client num)/10. (default 1024)
    -p,  --port <port>                Use port for clients (default: 1883)
         --host <host>                Network host to listen on. (default "0.0.0.0")
    -ws, --wsport <port>              Use port for websocket monitoring
    -wsp,--wspath <path>              Use path for websocket monitoring
    -c,  --config <file>              Configuration file

Logging Options:
    -d, --debug <bool>                Enable debugging output (default false)
    -D                                Debug enabled

Cluster Options:
    -r,  --router  <rurl>             Router who maintenance cluster info
    -cp, --clusterport <cluster-port> Cluster listen port for others

Common Options:
    -h, --help                        Show this message
```

#### 配置

```json
{
	"workerNum": 4096,
	"port": "1883",
	"host": "0.0.0.0",
	"cluster": {
		"host": "0.0.0.0",
		"port": "1993"
	},
	"router": "127.0.0.1:9888",
	"wsPort": "1888",
	"wsPath": "/ws",
	"wsTLS": true,
	"tlsPort": "8883",
	"tlsHost": "0.0.0.0",
	"tlsInfo": {
		"verify": true,
		"caFile": "tls/ca/cacert.pem",
		"certFile": "tls/server/cert.pem",
		"keyFile": "tls/server/key.pem"
	},
	"plugins": {
		"auth": "authhttp",
		"bridge": "kafka"
	}
}
```

```

```
```

```
```

```
```

```
```

```
```

```
```

```
