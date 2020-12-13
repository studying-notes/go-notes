---
date: 2020-11-15T13:33:47+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go select 语句理解"  # 文章标题
url:  "posts/go/abc/select"  # 设置网页永久链接
tags: [ "go", "select" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 前言

select 是 go 在语言层面提供的多路 IO 复用的机制，其可以检测多个 channel 是否 ready (即是否可读或可写)，使用起来非常方便。

## 热身环节

### 题目一

下面的程序输出是什么？

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    chan1 := make(chan int)
    chan2 := make(chan int)

    go func() {
        chan1 <- 1
        time.Sleep(5 * time.Second)
    }()

    go func() {
        chan2 <- 1
        time.Sleep(5 * time.Second)
    }()

    select {
    case <-chan1:
        fmt.Println("chan1 ready.")
    case <-chan2:
        fmt.Println("chan2 ready.")
    default:
        fmt.Println("default")
    }

    fmt.Println("main exit.")
}
```

程序中声明两个 channel，分别为 chan1 和 chan2，依次启动两个协程，分别向两个 channel 中写入一个数据就进入睡眠。select 语句两个 case 分别检测 chan1 和 chan2 是否可读，如果都不可读则执行 default 语句。

参考答案：  

select 中各个 case 执行顺序是随机的，如果某个 case 中的 channel 已经 ready，则执行相应的语句并退出 select 流程，如果所有 case 中的 channel 都未 ready，则执行 default 中的语句然后退出 select 流程。另外，由于启动的协程和 select 语句并不能保证执行顺序，所以也有可能 select 执行时协程还未向 channel 中写入数据，所以 select 直接执行 default 语句并退出。所以，以下三种输出都有可能：

可能的输出一：

```
chan1 ready.
main exit.
```

可能的输出二：

```
chan2 ready.
main exit.
```

可能的输出三：

```
default
main exit.
```

### 题目二

下面的程序执行到 select 时会发生什么？

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    chan1 := make(chan int)
    chan2 := make(chan int)

    writeFlag := false
    go func() {
        for {
            if writeFlag {
                chan1 <- 1
            }
            time.Sleep(time.Second)
        }
    }()

    go func() {
        for {
            if writeFlag {
                chan2 <- 1
            }
            time.Sleep(time.Second)
        }
    }()

    select {
    case <-chan1:
        fmt.Println("chan1 ready.")
    case <-chan2:
        fmt.Println("chan2 ready.")
    }

    fmt.Println("main exit.")
}
```

程序中声明两个 channel，分别为 chan1 和 chan2，依次启动两个协程，协程会判断一个 bool 类型的变量 writeFlag 来决定是否要向 channel 中写入数据，由于 writeFlag 永远为 false，所以实际上协程什么也没做。select 语句两个 case 分别检测 chan1 和 chan2 是否可读，这个 select 语句不包含 default 语句。

参考答案：

select 会按照随机的顺序检测各 case 语句中 channel 是否 ready，如果某个 case 中的 channel 已经 ready 则执行相应的 case 语句然后退出 select 流程，如果所有的 channel 都未 ready 且没有 default 的话，则会阻塞等待各个 channel。所以上述程序会一直阻塞。

### 题目三

下面程序有什么问题？

```go
package main

import (
    "fmt"
)

func main() {
    chan1 := make(chan int)
    chan2 := make(chan int)

    go func() {
        close(chan1)
    }()

    go func() {
        close(chan2)
    }()

    select {
    case <-chan1:
        fmt.Println("chan1 ready.")
    case <-chan2:
        fmt.Println("chan2 ready.")
    }

    fmt.Println("main exit.")
}
```

程序中声明两个 channel，分别为 chan1 和 chan2，依次启动两个协程，协程分别关闭两个 channel。select 语句两个 case 分别检测 chan1 和 chan2 是否可读，这个 select 语句不包含 default 语句。

参考答案：

select 会按照随机的顺序检测各 case 语句中 channel 是否 ready，考虑到已关闭的 channel 也是可读的，所以上述程序中 select 不会阻塞，具体执行哪个 case 语句具是随机的。

### 题目四

下面程序会发生什么？

```go
package main

func main() {
    select {
    }
}
```

上面程序中只有一个空的 select 语句。

参考答案：

对于空的 select 语句，程序会被阻塞，准确的说是当前协程被阻塞，同时 go 自带死锁检测机制，当发现当前协程再也没有机会被唤醒时，则会 panic。所以上述程序会 panic。


## 实现原理

go 实现 select 时，定义了一个数据结构表示每个 case 语句 ( 含 defaut，default 实际上是一种特殊的 case)，select 执行过程可以类比成一个函数，函数输入 case 数组，输出选中的 case，然后程序流程转到选中的 case 块。

### case 数据结构

源码包 `src/runtime/select.go:scase` 定义了表示 case 语句的数据结构：

```go
type scase struct {
	c           *hchan         // chan
	kind        uint16
	elem        unsafe.Pointer // data element
}
```

scase.c 为当前 case 语句所操作的 channel 指针，这也说明了一个 case 语句只能操作一个 channel。scase.kind 表示该 case 的类型，分为读 channel、写 channel 和 default，三种类型分别由常量定义：

- caseRecv：case 语句中尝试读取 scase.c 中的数据；
- caseSend：case 语句中尝试向 scase.c 中写入数据；
- caseDefault：default 语句

scase.elem 表示缓冲区地址，根据 scase.kind 不同，有不同的用途：

- scase.kind == caseRecv：scase.elem 表示读出 channel 的数据存放地址；
- scase.kind == caseSend：scase.elem 表示将要写入 channel 的数据存放地址；

### select 实现逻辑

源码包 `src/runtime/select.go:selectgo()` 定义了 select 选择 case 的函数：

```go
func selectgo(cas0 *scase, order0 *uint16, ncases int) (int, bool)
```

函数参数：
- cas0 为 scase 数组的首地址，selectgo() 就是从这些 scase 中找出一个返回。
- order0 为一个两倍 cas0 数组长度的 buffer，保存 scase 随机序列 pollorder 和 scase 中 channel 地址序列 lockorder
	 - pollorder：每次 selectgo 执行都会把 scase 序列打乱，以达到随机检测 case 的目的。
	 - lockorder：所有 case 语句中 channel 序列，以达到去重防止对 channel 加锁时重复加锁的目的。
- ncases 表示 scase 数组的长度

函数返回值：

1. int： 选中 case 的编号，这个 case 编号跟代码一致
2. bool: 是否成功从 channle 中读取了数据，如果选中的 case 是从 channel 中读数据，则该返回值表示是否读取成功。

selectgo 实现伪代码如下：

```go
func selectgo(cas0 *scase, order0 *uint16, ncases int) (int, bool) {
    //1. 锁定scase语句中所有的channel
    //2. 按照随机顺序检测scase中的channel是否ready
    //  # 如果case可读，则读取channel中数据，解锁所有的channel，然后返回(case index, true)
    //  # 如果case可写，则将数据写入channel，解锁所有的channel，然后返回(case index, false)
    //  # 所有case都未ready，则解锁所有的channel，然后返回（default index, false）
    //3. 所有case都未ready，且没有default语句
    //  # 将当前协程加入到所有channel的等待队列
    //  # 当将协程转入阻塞，等待被唤醒
    //4. 唤醒后返回channel对应的case index
    //  # 如果是读操作，解锁所有的channel，然后返回(case index, true)
    //  # 如果是写操作，解锁所有的channel，然后返回(case index, false)
}
```

特别说明：对于读 channel 的 case 来说，如 `case elem, ok := <-chan1 : `, 如果 channel 有可能被其他协程关闭的情况下，一定要检测读取是否成功，因为 close 的 channel 也有可能返回，此时 ok == false。

## 小结

- select 语句中除 default 外，每个 case 操作一个 channel，要么读要么写
- select 语句中除 default 外，各 case 执行顺序是随机的
- select 语句中如果没有 default 语句，则会阻塞等待任一 case
- select 语句中读操作要判断是否成功读取，关闭的 channel 也可以读取
