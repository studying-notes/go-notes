---
date: 2020-11-15T16:56:45+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 用 WaitGroup 控制协程"  # 文章标题
url:  "posts/go/docs/internal/concurrent/waitgroup"  # 设置网页永久链接
tags: [ "Go", "waitgroup" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

WaitGroup 是 Golang 应用开发过程中经常使用的并发控制技术。

WaitGroup，可理解为 Wait-Goroutine-Group，即等待一组 goroutine 结束。比如某个 goroutine 需要等待其他几个 goroutine 全部完成，那么使用 WaitGroup 可以轻松实现。

## 示例

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		time.Sleep(1 * time.Second)

		fmt.Println("Goroutine 1 finished!")
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)

		fmt.Println("Goroutine 2 finished!")
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("All Goroutine finished!")
}
```

简单的说，上面程序中 wg 内部维护了一个计数器：

1. 启动 goroutine 前将计数器通过 Add(2) 将计数器设置为待启动的 goroutine 个数。

2. 启动 goroutine 后，使用 Wait() 方法阻塞自己，等待计数器变为 0。

3. 每个 goroutine 执行结束通过 Done() 方法将计数器减 1。

4. 计数器变为 0 后，阻塞的 goroutine 被唤醒。

## 信号量

信号量是 Unix 系统提供的一种保护共享资源的机制，用于防止多个线程同时访问某个资源。

可简单理解为信号量为一个数值：

- 当信号量 > 0 时，表示资源可用，获取信号量时系统自动将信号量减 1 ；

- 当信号量 == 0 时，表示资源暂不可用，获取信号量时，当前线程会进入睡眠，当信号量为正时被唤醒；

WaitGroup 实现中使用了信号量。

## 数据结构

源码包中 `src/sync/waitgroup.go:WaitGroup` 定义了其数据结构：

```go
type WaitGroup struct {
	state1 [3]uint32
}
```

state1 是个长度为 3 的数组，其中包含了 state 和一个信号量，而 state 实际上是两个计数器：

- counter：当前还未执行结束的 goroutine 计数器
- waiter count : 等待 goroutine-group 结束的 goroutine 数量，即有多少个等候者
- semaphore : 信号量

考虑到字节是否对齐，三者出现的位置不同，为简单起见，依照字节已对齐情况下，三者在内存中的位置如下所示：

![](https://dd-static.jd.com/ddimg/jfs/t1/27222/38/19271/9896/6312f359E45c28096/23475629b2c555ef.png)

WaitGroup 对外提供三个接口：

- Add(delta int) : 将 delta 值加到 counter 中
- Wait()：waiter 递增 1，并阻塞等待信号量 semaphore
- Done()：counter 递减 1，按照 waiter 数值释放相应次数信号量

下面分别介绍这三个函数的实现细节。

## Add(delta int)

Add() 做了两件事，一是把 delta 值累加到 counter 中，因为 delta 可以为负值，也就是说 counter 有可能变成 0 或负值，所以第二件事就是当 counter 值变为 0 时，根据 waiter 数值释放等量的信号量，把等待的 goroutine 全部唤醒，如果 counter 变为负值，则 panic。

Add() 伪代码如下：

```go
func (wg *WaitGroup) Add(delta int) {
    statep, semap := wg.state() //获取state和semaphore地址指针

    state := atomic.AddUint64(statep, uint64(delta)<<32) //把delta左移32位累加到state，即累加到counter中
    v := int32(state >> 32) //获取counter值
    w := uint32(state)      //获取waiter值

    if v < 0 {              //经过累加后counter值变为负值，panic
        panic("sync: negative WaitGroup counter")
    }

    //经过累加后，此时，counter >= 0
    //如果counter为正，说明不需要释放信号量，直接退出
    //如果waiter为零，说明没有等待者，也不需要释放信号量，直接退出
    if v > 0 || w == 0 {
        return
    }

    //此时，counter一定等于0，而waiter一定大于0（内部维护waiter，不会出现小于0的情况），
    //先把counter置为0，再释放waiter个数的信号量
    *statep = 0
    for ; w != 0; w-- {
        runtime_Semrelease(semap, false) //释放信号量，执行一次释放一个，唤醒一个等待者
    }
}
```

## Wait()

Wait() 方法也做了两件事，一是累加 waiter，二是阻塞等待信号量

```go
func (wg *WaitGroup) Wait() {
    statep, semap := wg.state() //获取state和semaphore地址指针
    for {
        state := atomic.LoadUint64(statep) //获取state值
        v := int32(state >> 32)            //获取counter值
        w := uint32(state)                 //获取waiter值
        if v == 0 {                        //如果counter值为0，说明所有goroutine都退出了，不需要待待，直接返回
            return
        }

        // 使用CAS（比较交换算法）累加waiter，累加可能会失败，失败后通过for loop下次重试
        if atomic.CompareAndSwapUint64(statep, state, state+1) {
            runtime_Semacquire(semap) //累加成功后，等待信号量唤醒自己
            return
        }
    }
}
```

这里用到了 CAS 算法保证有多个 goroutine 同时执行 Wait() 时也能正确累加 waiter。

## Done()

Done() 只做一件事，即把 counter 减 1，我们知道 Add() 可以接受负值，所以 Done 实际上只是调用了 Add(-1)。

源码如下：

```go
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}
```

Done() 的执行逻辑就转到了 Add()，实际上也正是最后一个完成的 goroutine 把等待者唤醒的。

```go

```
