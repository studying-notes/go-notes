---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "限流器/限速器"  # 文章标题
url:  "posts/gin/project/limiter"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 简介

限流会导致用户在短时间内（这个时间段是毫秒级的）系统不可用，一般我们衡量系统处理能力的指标是每秒的 QP 或者 TPS，假设系统每秒的流量阈值是 1000，理论上一秒内有第 1001 个请求进来时，那么这个请求就会被限流。

## 计数器法

### 普通计数

```go
package limiter

import (
	"sync/atomic"
	"time"
)

// CountLimiter 普通计数限流器，缺点在于不够均匀，瞬时流量可能超过限制
type CountLimiter struct {
	counter  int64         // 计数器
	max      int64         // 最大数量
	interval time.Duration // 间隔时间
	last     time.Time     // 上一次时间
}

func NewCountLimiter(max int64, interval time.Duration) *CountLimiter {
	return &CountLimiter{
		counter:  0,
		max:      max,
		interval: interval,
		last:     time.Now(), // 可能不必要
	}
}

func (c *CountLimiter) Allow() bool {
	current := time.Now()

	// 超过时间计数清零
	if current.After(c.last.Add(c.interval)) {
		atomic.StoreInt64(&c.counter, 1)
		c.last = current
		return true
	}

	// 取出一个
	atomic.AddInt64(&c.counter, 1)

	// 判断是否超过限流个数
	if c.counter <= c.max {
		return true
	}
	return false
}
```

### 滑动窗口

```URL
https://wangbjun.site/2020/coding/golang/limiter.html
```


### 漏桶算法

我们把水比作是请求，漏桶比作是系统处理能力极限，水先进入到漏桶里，漏桶里的水按一定速率流出，当流出的速率小于流入的速率时，由于漏桶容量有限，后续进入的水直接溢出（拒绝请求），以此实现限流。

[![BxyqVH.jpg](https://s3.ax1x.com/2020/11/12/BxyqVH.jpg)](https://imgchr.com/i/BxyqVH)

```go

```

```go

```

### 令牌桶算法

系统会维护一个令牌（token）桶，以一个恒定的速度往桶里放入令牌（token），这时如果有请求进来想要被处理，则需要先从桶里获取一个令牌（token），当桶里没有令牌（token）可取时，则该请求将被拒绝服务。令牌桶算法通过控制桶的容量、发放令牌的速率，来达到对请求的限制。

[![Bx6ZR0.jpg](https://s3.ax1x.com/2020/11/12/Bx6ZR0.jpg)](https://imgchr.com/i/Bx6ZR0)

`golang.org/x/time/rate`

Limter 限制时间的发生频率，采用令牌池的算法实现。这个池子一开始容量为 b，装满 b 个令牌，然后每秒往里面填充 r 个令牌。
由于令牌池中最多有 b 个令牌，所以一次最多只能允许 b 个事件发生

Limter 提供三中主要的函数 Allow, Reserve, and Wait. 大部分时候使用 Wait。

首先创建一个 rate.Limiter, 其有两个参数，第一个参数为每秒发生多少次事件，第二个参数是其缓存最大可存多少个事件。

```go
// 每秒产生 200 * cpu 个数个令牌，最多存储 200 * cpu 个数个令牌
limiter = rate.NewLimiter(rate.Limit(200*NumCPU()), 200*runtime.NumCPU())
```

rate.Limiter 提供了三类方法用来限速

- Wait / WaitN 当没有可用或足够的事件时，将阻塞等待 推荐实际程序中使用这个方法
- Allow / AllowN 当没有可用或足够的事件时，返回 false 
- Reserve / ReserveN 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件

**基本用法**

```go
package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	limiter := rate.NewLimiter(
		rate.Limit(1), // 每秒产生的数量
		10,            // 桶容量大小
	)
	for i := 0; i < 100; i++ {
		for limiter.Allow() {
			fmt.Println(i)
		}
		time.Sleep(time.Second)
	}
}
```

详解文章：

```URL
https://www.cyhone.com/articles/usage-of-golang-rate/
```

