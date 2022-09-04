---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 并发同步互斥锁"  # 文章标题
url:  "posts/go/libraries/standard/sync/mutex"  # 设置网页永久链接
tags: [ "Go", "mutex" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

一般情况下，原子操作都是通过“互斥”访问来保证的，通常由特殊的 CPU 指令提供保护。当然，如果仅仅是想模拟粗粒度的原子操作，可以借助于 `sync.Mutex` 来实现：

```go
import (
    "sync"
)

var total struct {
    sync.Mutex
    value int
}

func worker(wg *sync.WaitGroup) {
    defer wg.Done()

    for i := 0; i <= 10000; i++ {
        total.Lock()
        total.value += i
        total.Unlock()
    }
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go worker(&wg)
    go worker(&wg)
    wg.Wait()

    fmt.Println(total.value)
}
```

在 `worker` 的循环中，为了保证 `total.value += i` 的原子性，我们通过 `sync.Mutex` 加锁和解锁来保证该语句在同一时刻只被一个线程访问。对多线程模型的程序而言，**进出临界区前后进行加锁和解锁都是必需的**。如果没有锁的保护，`total` 的最终值将由于多线程之间的竞争而可能不正确。

```go

```
