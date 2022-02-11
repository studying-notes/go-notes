---
date: 2020-09-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "time - 时间标准库"  # 文章标题
url:  "posts/go/libraries/standard/time"  # 设置网页链接，默认使用文件名
tags: [ "go", "log", "logger" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

- [定时器](#定时器)
- [循环定时器](#循环定时器)
- [时间格式化](#时间格式化)

## 定时器

Timer 定时器，到了设定的时间执行一次，通过 Reset 可实现 Ticker 相同功能。

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	timer := time.NewTimer(2 * time.Second)

	go func(t *time.Timer) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("Contains Timer", time.Now().Format("2006-01-02 15:04:05"))
			t.Reset(2 * time.Second) // 重置定时器，可实现类似 Ticker 的功能
		}
	}(timer)

	wg.Wait()
}
```

## 循环定时器

Ticker 每隔固定的时间执行一次。

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	ticker := time.NewTicker(2 * time.Second)

	go func(t *time.Ticker) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("Contains Ticker", time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker)

	wg.Wait()
}
```

## 时间格式化

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 系统默认格式打印当前时间
	t1 := time.Now()
	fmt.Println("t1:", t1)

	// 自定义格式
	t2 := t1.Format("2006-01-02 15:04:05")
	fmt.Println("t2:", t2)

	// 换个时间定义格式不行么？不行！
	t3 := t1.Format("2020-07-01 21:00:00")
	fmt.Println("t3:", t3)

	// 自定义解析时间字符串格式
	t4, _ := time.Parse("2006-01-02 15:04:05", "2018-10-01 14:51:00")
	fmt.Println("t4:", t4)
}
```

在Go语言中，强调必须显示参考时间的格式，因此每个布局字符串都是一个时间戳，而并非随便写的时间点。

对于 2006-01-02 15: 04: 05，可以将其记忆为2006年1月2日3点4分5秒。

```go

```

```go

```
