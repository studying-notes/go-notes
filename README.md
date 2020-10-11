---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 学习笔记目录"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/readme"  # 设置网页永久链接
tags: [ "go", "toc" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

> 纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。

## 预准备

{{<card src="posts/go/doc/install">}}

## 语言基础

{{<card src="posts/go/abc/array">}}
{{<card src="posts/go/abc/string">}}
{{<card src="posts/go/abc/slice">}}
{{<card src="posts/go/abc/func">}}
{{<card src="posts/go/abc/method">}}
{{<card src="posts/go/abc/interface">}}
{{<card src="posts/go/abc/goroutine">}}
{{<card src="posts/go/abc/concurrent">}}
{{<card src="posts/go/abc/error">}}
{{<card src="posts/go/abc/attention">}}

### 类型转换

{{<card src="posts/go/abc/assert">}}
{{<card src="posts/go/libraries/standard/strconv">}}

## 算法与数据结构

## CGO 编程

{{<card src="posts/go/cgo/quickstart">}}
{{<card src="posts/go/cgo/intro">}}
{{<card src="posts/go/cgo/dll">}}
{{<card src="posts/go/cgo/func">}}
{{<card src="posts/go/cgo/link">}}
{{<card src="posts/go/cgo/type">}}
{{<card src="posts/go/cgo/internal">}}

## 命令行

{{<card src="posts/go/cmd/compile">}}

### 构建命令行程序

{{<card src="posts/go/libraries/tripartite/cobra">}}
{{<card src="posts/go/libraries/standard/flag">}}

## 输入输出

### 标准输入输出

{{<card src="posts/go/libraries/standard/bufio">}}
{{<card src="posts/go/libraries/standard/fmt">}}

### 模板引擎

{{<card src="posts/go/libraries/standard/template">}}

### 数据库

{{<card src="posts/go/libraries/tripartite/gorm">}}
{{<card src="posts/go/libraries/tripartite/sqlx">}}
{{<card src="posts/go/io/sqlite">}}
{{<card src="posts/go/libraries/tripartite/sqlcipher">}}
{{<card src="posts/go/io/mysql">}}
{{<card src="posts/go/io/redis">}}
{{<card src="posts/go/io/mongo">}}

### 文件读写

{{<card src="posts/go/libraries/tripartite/fsnotify">}}
{{<card src="posts/go/io/excel">}}

### 序列化

{{<card src="posts/go/libraries/standard/json">}}

### 日志

{{<card src="posts/go/libraries/standard/log">}}
{{<card src="posts/go/libraries/tripartite/logrus">}}
{{<card src="posts/go/libraries/tripartite/zap">}}

### 程序配置

{{<card src="posts/go/libraries/tripartite/viper">}}

## 并发编程

{{<card src="posts/go/libraries/standard/sync/pool">}}
{{<card src="posts/go/libraries/standard/context">}}

## 系统

### 执行命令

{{<card src="posts/go/libraries/standard/exec">}}

### 硬件监测

{{<card src="posts/go/libraries/tripartite/gopsutil">}}

### 时间

{{<card src="posts/go/libraries/standard/time">}}

## 网络编程

{{<card src="posts/go/web/http/cookie">}}
{{<card src="posts/go/web/http/httpclient">}}
{{<card src="posts/go/web/grpc">}}
{{<card src="posts/go/web/mqtt/intro">}}

### 消息队列

{{<card src="posts/go/web/mq/intro">}}
{{<card src="posts/go/web/mq/kafka">}}
{{<card src="posts/go/web/mq/nsq">}}
{{<card src="posts/go/web/mq/rabbitmq">}}

## 测试与性能

{{<card src="posts/go/doc/benchmark">}}
{{<card src="posts/go/doc/unittest">}}
