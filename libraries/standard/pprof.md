---
date: 2020-11-11T16:27:36+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "PProf 之性能剖析"  # 文章标题
url:  "posts/go/libraries/standard/pprof"  # 设置网页永久链接
tags: [ "go", "pprof"]  # 标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## PProf 简介

在 Go 语言中，PProf 是分析性能、分析数据的工具，PProf 用 profile.proto 读取分析样本的集合，并生成可视化报告，以帮助分析数据（支持文本和图形报告）。

profile.proto 是一个 Protobuf v3 的描述文件，它描述了一组 callstack 和 symbolization 信息，它的作用是统计分析一组采样的调用栈，配置文件格式是 stacktrace。

### 采样方式

- runtime/pprof：采集程序（非 Server ）指定区块的运行数据进行分析。
- net/http/pprof：基于 HTTP Server 运行，并且可以采集运行时的数据进行分析。
- go test：通过运行测试用例，指定所需标识进行采集。

### 使用模式

- Report Generation：报告生成。
- Interactive Terminal Use：交互式终端使用。
- Web Interface：Web 界面。

### 可以做什么

- CPU Profiling：CPU 分析。按照一定的频率采集所监听的应用程序 CPU （含寄存器）的使用情况，确定应用程序在主动消耗 CPU 周期时花费时间的位置。
- Memory Profiling：内存分析。在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏。
- Block Profiling：阻塞分析。记录 goroutine 阻塞等待同步（包括定时器通道）的位置，默认不开启，需要调用 runtime.SetBlockProfileRate 进行设置。
- Mutex Profiling：互斥锁分析。报告互斥锁的竞争情况，默认不开启，需要调用 runtime.SetMutexProfileFraction 进行设置。
- Goroutine Profiling：goroutine 分析，可以对当前应用程序正在运行的 goroutine 进行堆栈跟踪和分析。这项功能在实际排查中会经常用到，因为很多问题出现时的表象就是 goroutine 暴增，而这时候我们要做的事情之一就是查看应用程序中的 goroutine 正在做什么事情，因为什么阻塞了，然后再进行下一步。

## PProf 的使用

```go
package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
```

### 通过访问网址

```url
http://localhost:6060/debug/pprof/
```

- allocs：查看过去所有内存分配的样本
- block：查看导致阻塞同步的堆栈跟踪
- cmdline：当前程序命令行的完整调用路径
- goroutine：查看当前所有运行的 goroutines 堆栈跟踪
- heap：查看活动对象的内存分配情况
- mutex：查 看 导 致 互 斥 锁 的 竞 争 持 有 者 的 堆 栈 跟 踪
- profile：默认进行 30s 的 CPU Profiling，会得到一个分析用的 profile 文件
- threadcreate：查看创建新 OS 线程的堆栈跟踪

### 通过交互式终端

```shell
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

输入查询命令 `top10`，查看对应资源开销（例如，CPU 就是执行耗时/开销、Memory 就是内存占用大小）排名前十的函数。

```shell
top10
```

- flat：函数自身的运行耗时。
- flat%：函数自身占 CPU 运行总耗时的比例。
- sum%：函数自身累积使用占 CPU 运行总耗时比例。
- cum：函数自身及其调用函数的运行总耗时。
- cum%：函数自身及其调用函数占 CPU 运行总耗时的比例。
- Name：函数名。


#### Heap Profile

分析应用程序常驻内存的占用情况。

```shell
go tool pprof http://localhost:6060/debug/pprof/heap
```

#### CPU Profile

```shell
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

#### Blocking Profile

```shell
go tool pprof http://localhost:6060/debug/pprof/block
```

#### Mutexes

```shell
go tool pprof http://localhost:6060/debug/pprof/mutex
```

#### Trace

```shell
wget -O trace.out http://localhost:6060/debug/pprof/trace?seconds=5
```

### 可视化界面

安装 graphviz 工具。

```shell
https://www2.graphviz.org/Packages/stable/windows/10/cmake/Release/x64/graphviz-install-2.44.1-win64.exe
```

```shell
dot -verson
```

```shell
wget http://localhost:6060/debug/pprof/profile
```

网页浏览：

```shell
go tool pprof -http=:6001 profile
```

## 通过测试用例做剖析

```shell
go test -bench. -cpuprofile=cpu.profile
```


```go

```

```go

```

```go

```

