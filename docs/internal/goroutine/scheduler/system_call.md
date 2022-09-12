---
date: 2022-09-11T15:23:03+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言调度器源码分析之系统调用"  # 文章标题
url:  "posts/go/docs/internal/goroutine/scheduler/system_call"  # 设置网页永久链接
tags: [ "Go", "system-call" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

**系统调用是指使用类似函数调用的方式调用操作系统提供的API。**

虽然从概念上来说系统调用和函数调用差不多，但本质上它们有很大的不同，**操作系统的代码位于内核地址空间**，而 CPU 在执行用户代码时特权等级很低，无权访问**需要最高优先级才能访问的内核地址空间的代码和数据**，所以不能通过简单的 call 指令直接调用操作系统提供的函数，而**需要使用特殊的指令进入操作系统内核完成指定的功能**。

另外，用户代码调用操作系统 API 也不是根据函数名直接调用，而是需要根据操作系统为每个 API 提供的一个整型编号来调用，AMD64 Linux 平台约定在进行系统调用时使用 rax 寄存器存放系统调用编号，同时约定使用 rdi, rsi, rdx, r10, r8 和 r9 来传递前 6 个系统调用参数。

打开文件，读写文件以及网络编程中的创建 socket 等等都使用了系统调用，我们没有感觉到系统调用的存在主要是因为我们使用的函数库或 package 把它们封装成了函数，我们只需要直接调用这些函数就可以了。比如有下面一段 go 代码：

```go
package main

import (
        "os"
)

func main() {
        fd, err := os.Open("./syscall.go")  // 将会使用系统调用打开文件
        fd.Close()  // 将会使用系统调用关闭文件
}
```

这里的 os.Open() 和 fd.Close() 函数最终都会通过系统调用进入操作系统内核完成相应的功能。

以 os.Open 为例，它最终会执行下面这段汇编代码来通过 openat 系统调用打开文件：

```
mov    0x10(%rsp),%rdi  #第1个参数
mov    0x18(%rsp),%rsi  #第2个参数
mov    0x20(%rsp),%rdx #第3个参数
mov    0x28(%rsp),%r10 #第4个参数
mov    0x30(%rsp),%r8  #第5个参数
mov    0x38(%rsp),%r9  #第6个参数
mov    0x8(%rsp),%rax  #系统调用编号 rax = 267，表示调用 openat 系统调用
syscall                           #系统调用指令，进入 Linux 内核
```

这里，代码首先把 6 个参数以及 openat 这个系统调用的编号 267 保存在了对应的寄存器中，然后使用 syscall 指令进入内核执行打开文件的功能。
