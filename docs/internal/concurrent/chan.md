---
date: 2020-11-15T16:56:28+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 用 chan 控制协程"  # 文章标题
url:  "posts/go/docs/internal/concurrent/chan"  # 设置网页永久链接
tags: [ "Go", "chan" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

channel 一般用于协程之间的通信，channel 也可以用于并发控制。比如主协程启动 N 个子协程，主协程等待所有子协程退出后再继续后续流程，这种场景下 channel 也可轻易实现。

## 示例

```go
package main

import (
	"fmt"
	"time"
)

func Process(ch chan int) {
	time.Sleep(time.Second)

	ch <- 1
}

func main() {
	channels := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		channels[i] = make(chan int)
		go Process(channels[i])
	}

	for i, ch := range channels {
		<-ch
		fmt.Println("Routine ", i, " quit!")
	}
}
```

上面程序通过创建 N 个 channel 来管理 N 个协程，每个协程都有一个 channel 用于跟父协程通信，父协程创建完所有协程后等待所有协程结束。

这个例子中，父协程仅仅是等待子协程结束，其实父协程也可以向管道中写入数据通知子协程结束，这时子协程需要定期地探测管道中是否有消息出现。

## 小结

使用 channel 来控制子协程的优点是实现简单，缺点是当需要**大量创建协程时就需要有相同数量的 channel**，而且对于子协程继续派生出来的协程不方便控制。

WaitGroup、Context 看起来比 channel 优雅一些，在各种开源组件中使用频率比 channel 高得多。

```go

```
