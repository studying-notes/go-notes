---
date: 2022-10-09T16:35:00+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "通道与协程间通信"  # 文章标题
url:  "posts/go/docs/internal/goroutine/channel/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

通道（channel）是 Go 语言中提供协程间通信的独特方式。通过通道交流的方式，Go 语言屏蔽了底层实现的诸多细节，使得并发编程更加简单快捷。将通道作为 Go 语言中的一等公民，是 Go 语言遵循 CSP 并发编程模式的结果，这种模型最重要的思想是通过通道来传递消息。同时，通道借助 Go 语言调度器的设计，可以高效实现通道的堵塞 / 唤醒，进一步实现通道多路复用的机制。

- [CSP 并发编程](#csp-并发编程)
- [通道基本使用方式](#通道基本使用方式)
- [通道底层原理](#通道底层原理)

## CSP 并发编程

在计算机科学中，CSP（Communicating Sequential Processes，通信顺序进程）是用于描述并发系统中交互模式的形式化语言，其通过通道传递消息。

Tony Hoare 在 1978 年第一次发表了关于 CSP 的想法。在过去，多线程或多进程程序通常采用共享内存进行交流，通过信号量等手段实现同步机制，但 Tony Hoare 通过同步交流（Synchronous Communication）的原理解决了交流与同步这两个问题。

在 CSP 语言中，通过命名的通道发送或接收值进行通信。在最初的设计中，通道是无缓冲的，因此发送操作会阻塞，直到被接收端接收后才能继续发送，从而提供了一种同步机制。

CSP 的思想后来影响了 Alef、Newsqueak、Limbo 等多种编程语言，并成为 Go 语言并发中的重要设计思想。

## 通道基本使用方式

[通道基本使用方式](operation.md)

## 通道底层原理

[通道底层原理](underlying_principle.md)

```go

```
