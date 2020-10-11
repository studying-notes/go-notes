---
date: 2020-09-19T13:39:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 内存池 sync.Pool"  # 文章标题
url:  "posts/go/libraries/standard/sync/pool"  # 设置网页链接，默认使用文件名
tags: [ "go", "exec" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
toc: false  # 是否自动生成目录
---

`sync.Pool` 可以作为保存临时取还对象的一个“池子”， Pool 里装的对象可以被无通知地被回收。

对于很多需要重复分配、回收内存的地方，`sync.Pool` 是一个很好的选择。频繁地分配、回收内存会给 GC 带来一定的负担，而 `sync.Pool` 可以将暂时不用的对象缓存起来，待下次需要的时候直接使用，不用再次经过内存分配，复用对象的内存，减轻 GC 的压力，提升系统的性能。

`sync.Pool` 是协程安全的，设置好对象的 `New` 函数，用于在 `Pool` 里没有缓存的对象时，创建一个。之后，在程序的任何地方、任何时候仅通过 `Get()`、`Put()` 方法就可以取、还对象了。

## 适用场景

当多个 goroutine 都需要创建同一个对象的时候，如果 goroutine 数过多，导致对象的创建数目剧增进而导致 GC 压力增大。形成 “并发大—占用内存大—GC 缓慢—处理并发能力降低—并发更大” 这样的恶性循环。

## 简单示例

```go
package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func InitPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Person")
			return new(Person)
		},
	}
}

func main() {
    InitPool()
    // 返回的是空接口类型
    // 有必要进行类型转换
	p := pool.Get().(*Person)
	fmt.Println("first from pool, get:", p)

	p.Name = "first"
	fmt.Println("set p.Name =", p.Name)

	pool.Put(p)

	fmt.Println("pool has a Person, get:", pool.Get().(*Person))
	fmt.Println("pool has no Person, get:", pool.Get().(*Person))
}
```

首先，需要初始化 `Pool`，即设置好 `New` 函数。当调用 Get 方法时，如果池子里缓存了对象，就直接返回缓存的对象。如果没有，则调用 `New` 函数创建一个新的对象。

`Get` 方法取出来的对象和上次 `Put` 进去的对象实际上是同一个，`Pool` 没有做任何“清空”的处理。但是在实际的并发使用场景中，无法保证这种顺序，最好的做法是在 `Put` 前，将对象清空。
