---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "RabbitMQ 简介"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/web/mq/rabbitmq"  # 设置网页永久链接
tags: [ "go", "mq", "rabbitmq" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 预准备

1. 下载安装 Erlang 开发环境

```
https://www.erlang.org/downloads
```

```powershell
> erl -v
Eshell V11.0  (abort with ^G)
1>
```

2. 下载安装 RabbitMQ

```
https://github.com/rabbitmq/rabbitmq-server/releases
```

默认注册为服务自动启动。

3. 下载安装 AMQP 协议实现库

```
go get github.com/streadway/amqp
```

## 安装 RabbitMQ 管理页面

1. 开启 Web 管理插件

```powershell
> rabbitmq-plugins enable rabbitmq_management
Enabling plugins on node rabbit@WHITE:
rabbitmq_management
The following plugins have been configured:
  rabbitmq_management
  rabbitmq_management_agent
  rabbitmq_web_dispatch
Applying plugin configuration to rabbit@WHITE...
The following plugins have been enabled:
  rabbitmq_management
  rabbitmq_management_agent
  rabbitmq_web_dispatch
```

2. 访问 `http://localhost:15672`

默认用户 `guest` 登录，密码也为 `guest`，即可进入管理界面。

## RabbitMQ 组件

RabbitMQ 中进行消息控制的组件可以分为以下几部分：

- EXCHANGE：路由部件，控制消息的转发路径；
- QUEUE：消息队列，可以有多个消费者从队列中读取消息；
- CONSUMER：消息的消费者。

可以单独用 `queue` 进行消息传递，也可以利用 `exchange` 与 `queue` 构建多种消息模式，主要包括 `fanout`、`direct` 和 `topic` 方式。

![](imgs/rabbitmq.png)
