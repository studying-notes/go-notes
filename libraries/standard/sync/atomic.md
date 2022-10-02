---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 并发同步原子操作"  # 文章标题
url:  "posts/go/libraries/standard/sync/atomic"  # 设置网页永久链接
tags: [ "Go", "atomic" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

**用互斥锁来保护一个数值型的共享资源麻烦且效率低下**。标准库的 `sync/atomic` 包对原子操作提供了丰富的支持。

```go
import (
	"fmt"
	"sync"
	"sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()
	fmt.Println(total)
}
```

`atomic.AddUint64()` 函数调用保证了 `total` 的读取、更新和保存是一个原子操作，因此在多线程中访问也是安全的。

```go

```
