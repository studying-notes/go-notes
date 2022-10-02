---
date: 2020-11-15T13:33:47+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 多路 IO 复用关键字 select"  # 文章标题
url:  "posts/go/docs/grammar/select"  # 设置网页永久链接
tags: [ "Go", "select" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

select 是 go 在语言层面提供的多路 IO 复用的机制，其可以检测多个 channel 是否 ready (即是否可读或可写)。

## 实现原理

go 实现 select 时，定义了一个数据结构表示每个 case 语句(含 defaut, default 实际上是一种特殊的 case)，select 执行过程可以类比成一个函数，函数输入 case 数组，输出选中的 case，然后程序流程转到选中的 case 块。

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
- caseDefault：default 语句。

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

1. int: 选中 case 的编号，这个 case 编号跟代码一致
2. bool: 是否成功从 channle 中读取了数据，如果选中的 case 是从 channel 中读数据，则该返回值表示是否读取成功。

selectgo 实现伪代码如下：

```go
func selectgo(cas0 *scase, order0 *uint16, ncases int) (int, bool) {
//1. 锁定 scase 语句中所有的 channel

 //2. 按照随机顺序检测 scase 中的 channel 是否 ready

 // # 如果 case 可读，则读取 channel 中数据，解锁所有的 channel，然后返回 (case index, true)

 // # 如果 case 可写，则将数据写入 channel，解锁所有的 channel，然后返回 (case index, false)

 // # 所有 case 都未 ready，则解锁所有的 channel，然后返回（default index, false）

 //3. 所有 case 都未 ready，且没有 default 语句

 // # 将当前协程加入到所有 channel 的等待队列

 // # 当将协程转入阻塞，等待被唤醒

 //4. 唤醒后返回 channel 对应的 case index

 // # 如果是读操作，解锁所有的 channel，然后返回 (case index, true)

 // # 如果是写操作，解锁所有的 channel，然后返回 (case index, false)
}
```

特别说明：对于读 channel 的 case 来说，如 `case elem, ok := <-chan1 : `, 如果 channel 有可能被其他协程关闭的情况下，一定要检测读取是否成功，因为 close 的 channel 也有可能返回，此时 ok == false。

- select 语句中除 default 外，每个 case 操作一个 channel，要么读要么写
- select 语句中除 default 外，各 case 执行顺序是随机的
- select 语句中如果没有 default 语句，则会阻塞等待任一 case
- select 语句中读操作要判断是否成功读取，关闭的 channel 也可以读取

## 示例

### 空语句

```go
package main

func main() {
    select {
    }
}
```

上面程序中只有一个空的 select 语句。

对于空的 select 语句，程序会被阻塞，准确的说是当前协程被阻塞，同时 go 自带死锁检测机制，当发现当前协程再也没有机会被唤醒时，则会 panic。所以上述程序会 panic。

```go

```
