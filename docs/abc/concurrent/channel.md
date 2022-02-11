---
date: 2020-11-15T16:56:28+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 用 channel 控制子协程"  # 文章标题
url:  "posts/go/abc/concurrent/channel"  # 设置网页永久链接
tags: [ "go", "goroutine" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 前言

channel 一般用于协程之间的通信，channel 也可以用于并发控制。比如主协程启动 N 个子协程，主协程等待所有子协程退出后再继续后续流程，这种场景下 channel 也可轻易实现。

## 场景示例

下面程序展示一个使用 channel 控制子协程的例子：

```go
package main

import (
    "time"
    "fmt"
)

func Process(ch chan int) {
    //Do some work...
    time.Sleep(time.Second)

    ch <- 1 //管道中写入一个元素表示当前协程已结束
}

func main() {
    channels := make([]chan int, 10) //创建一个10个元素的切片，元素类型为channel

    for i:= 0; i < 10; i++ {
        channels[i] = make(chan int) //切片中放入一个channel
        go Process(channels[i])      //启动协程，传一个管道用于通信
    }

    for i, ch := range channels {  //遍历切片，等待子协程结束
        <-ch
        fmt.Println("Routine ", i, " quit!")
    }
}
```

上面程序通过创建 N 个 channel 来管理 N 个协程，每个协程都有一个 channel 用于跟父协程通信，父协程创建完所有协程后等待所有协程结束。

这个例子中，父协程仅仅是等待子协程结束，其实父协程也可以向管道中写入数据通知子协程结束，这时子协程需要定期地探测管道中是否有消息出现。

## 小结

使用 channel 来控制子协程的优点是实现简单，缺点是当需要大量创建协程时就需要有相同数量的 channel，而且对于子协程继续派生出来的协程不方便控制。

后面继续介绍的 WaitGroup、Context 看起来比 channel 优雅一些，在各种开源组件中使用频率比 channel 高得多。
