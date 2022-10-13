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

通道（channel）是 Go 语言中提供协程间通信的独特方式。将通道作为 Go 语言中的一等公民，是 Go 语言遵循 CSP 并发编程模式的结果，这种模型最重要的思想是通过通道来传递消息。同时，通道借助 Go 语言调度器的设计，可以高效实现通道的堵塞 / 唤醒，进一步实现通道多路复用的机制。

## CSP 并发编程

[Go 语言的 GPM 调度器](../gpm.md)

## 通道基本使用方式

[通道基本使用方式](operation.md)

## 通道底层原理

[通道底层原理](underlying_principle.md)

```go

```
