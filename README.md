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
weight: 10 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

> 纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。

## 目录结构

- libraires 库
  - libraries/standard 标准库
  - libraries/tripartite 第三方库
- src 测试源码
- abc 深入语言实现
- cmds 基本命令
- quickstart 快速开始
- docs 学习笔记文档
- learn-src 知名库源码学习

## 预准备

{{<card src="posts/go/doc/install">}}
{{<card src="posts/go/doc/uninstall">}}

## 语言基础

### 数据结构

{{<card src="posts/go/abc/array">}}
{{<card src="posts/go/abc/string">}}
{{<card src="posts/go/abc/slice">}}
{{<card src="posts/go/abc/map">}}
{{<card src="posts/go/abc/func">}}
{{<card src="posts/go/abc/struct">}}
{{<card src="posts/go/abc/method">}}
{{<card src="posts/go/abc/interface">}}
{{<card src="posts/go/abc/goroutine2">}}
{{<card src="posts/go/abc/channel">}}
{{<card src="posts/go/abc/reflect">}}
{{<card src="posts/go/abc/append">}}
{{<card src="posts/go/abc/iota">}}
{{<card src="posts/go/abc/attention">}}

### 控制结构

{{<card src="posts/go/abc/defer">}}
{{<card src="posts/go/abc/recover">}}
{{<card src="posts/go/abc/error">}}
{{<card src="posts/go/abc/select">}}
{{<card src="posts/go/abc/range">}}
{{<card src="posts/go/abc/range2">}}
{{<card src="posts/go/abc/mutex">}}
{{<card src="posts/go/abc/rwmutex">}}

### 内存管理

{{<card src="posts/go/abc/memory/alloc">}}
{{<card src="posts/go/abc/memory/gc">}}
{{<card src="posts/go/abc/memory/escape">}}

## 并发控制

我们考虑这么一种场景，协程 A 执行过程中需要创建子协程 A1、A2、A3...An，协程 A 创建完子协程后就等待子协程退出。针对这种场景，Go 提供了三种解决方案：

- Channel : 使用 channel 控制子协程
- WaitGroup : 使用信号量机制控制子协程
- Context : 使用上下文控制子协程

三种方案各有优劣，比如 Channel 优点是实现简单，清晰易懂，WaitGroup 优点是子协程个数动态可调整，Context 优点是对子协程派生出来的孙子协程的控制。

{{<card src="posts/go/abc/concurrent/goroutine">}}
{{<card src="posts/go/abc/concurrent/concurrent">}}
{{<card src="posts/go/abc/concurrent/channel">}}
{{<card src="posts/go/abc/concurrent/waitgroup">}}
{{<card src="posts/go/abc/concurrent/context">}}
{{<card src="posts/go/libraries/standard/sync/pool">}}
{{<card src="posts/go/libraries/standard/context">}}

### 类型转换

{{<card src="posts/go/abc/assert">}}
{{<card src="posts/go/libraries/standard/strconv">}}

### 语法糖

语法糖（Syntactic Sugar），Go 中最常用的语法糖莫过于赋值符 `:=`，其次，表示函数变参的 `...`。

## 测试与性能

* 单元测试 - 指对软件中的最小可测试单元进行检查和验证，比如对一个函数的测试。
* 性能测试 - 也称基准测试，可以测试一段程序的性能，可以得到时间消耗、内存使用情况的报告。
* 示例测试 - 示例测试，广泛应用于 Go 源码和各种开源框架中，用于展示某个包或某个方法的用法。

{{<card src="posts/go/doc/test/benchmark_test">}}
{{<card src="posts/go/doc/test/unit_test">}}
{{<card src="posts/go/doc/test/example_test">}}
{{<card src="posts/go/doc/test/sub_test">}}
{{<card src="posts/go/doc/test/main_test">}}
{{<card src="posts/go/libraries/standard/pprof">}}

### 深入测试标准库

{{<card src="posts/go/abc/test/common">}}
{{<card src="posts/go/abc/test/tb_interface">}}
{{<card src="posts/go/abc/test/unit">}}
{{<card src="posts/go/abc/test/benchmark">}}
{{<card src="posts/go/abc/test/example">}}
{{<card src="posts/go/abc/test/main">}}
{{<card src="posts/go/abc/test/go_test">}}
{{<card src="posts/go/abc/test/go_test_params">}}
{{<card src="posts/go/abc/test/go_test_benchstat">}}

## 依赖管理

Go 语言依赖管理经历了三个阶段：

- GOPATH；
- vendor；
- Go Module；

{{<card src="posts/go/doc/mod/gopath">}}
{{<card src="posts/go/doc/mod/vendor">}}
{{<card src="posts/go/doc/mod/gomod">}}

### Go Module



## 算法与数据结构

{{<card src="posts/go/libraries/standard/rand">}}

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
{{<card src="posts/go/cmd/build">}}

### 构建命令行程序

{{<card src="posts/go/libraries/tripartite/cobra">}}
{{<card src="posts/go/libraries/standard/flag">}}

## 输入输出

### 标准输入输出

{{<card src="posts/go/libraries/standard/bufio">}}
{{<card src="posts/go/libraries/standard/fmt">}}
{{<card src="posts/go/libraries/standard/ioutil">}}

### 模板引擎

{{<card src="posts/go/libraries/standard/template">}}
{{<card src="posts/go/libraries/standard/regexp">}}

### 数据库

{{<card src="posts/go/libraries/tripartite/gorm">}}
{{<card src="posts/go/libraries/tripartite/sqlx">}}
{{<card src="posts/go/io/sqlite">}}
{{<card src="posts/go/libraries/tripartite/sqlcipher">}}
{{<card src="posts/go/io/mysql">}}
{{<card src="posts/go/io/redis">}}
{{<card src="posts/go/io/mongo">}}

### 图片处理

{{<card src="posts/go/io/image/draw">}}
{{<card src="posts/go/io/image/merge">}}

### 文件读写

{{<card src="posts/go/libraries/tripartite/fsnotify">}}
{{<card src="posts/go/io/excel">}}

## 网络服务

### HTTP 客户端

{{<card src="posts/go/web/http/httpclient">}}
{{<card src="posts/go/web/http/gout">}}

### HTTP 服务端

{{<card src="posts/go/web/http/cookie">}}
{{<card src="posts/go/web/grpc">}}
{{<card src="posts/go/web/mqtt/intro">}}

### 序列化

{{<card src="posts/go/libraries/standard/json">}}
{{<card src="posts/go/libraries/tripartite/gjson">}}

### 日志

{{<card src="posts/go/libraries/standard/log">}}
{{<card src="posts/go/libraries/tripartite/logrus">}}
{{<card src="posts/go/libraries/tripartite/zap">}}

### 程序配置

{{<card src="posts/go/libraries/tripartite/viper">}}

## 系统

### 定时任务

{{<card src="posts/go/libraries/tripartite/cron">}}

### 执行命令

{{<card src="posts/go/libraries/standard/exec">}}

### 服务

{{<card src="posts/go/libraries/standard/service">}}

### 硬件监测

{{<card src="posts/go/libraries/tripartite/gopsutil">}}

### 时间

{{<card src="posts/go/libraries/standard/time">}}
{{<card src="posts/go/doc/zone">}}

### 消息队列

{{<card src="posts/go/web/mq/intro">}}
{{<card src="posts/go/web/mq/kafka">}}
{{<card src="posts/go/web/mq/nsq">}}
{{<card src="posts/go/web/mq/rabbitmq">}}
