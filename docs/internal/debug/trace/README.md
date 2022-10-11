---
date: 2022-10-11T14:00:32+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "trace 事件追踪"  # 文章标题
url:  "posts/go/docs/internal/debug/trace/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

trace 工具非常强大，提供了追踪到的运行时的完整事件和宏观视野。尽管如此，trace 仍然不是万能的，如果想查看协程内部函数占用 CPU 的时间、内存分配等详细信息，就需要结合 pprof 来实现。

## 用法与说明

在 pprof 的分析中，能够知道一段时间内的 CPU 占用、内存分配、协程堆栈信息。这些信息都是一段时间内数据的汇总，但是它们并没有提供整个周期内发生的事件，例如指定的 Goroutines 何时执行、执行了多长时间、什么时候陷入了堵塞、什么时候解除了堵塞、GC 如何影响单个 Goroutine 的执行、STW 中断花费的时间是否太长等。这就是在 Go1.5 之后推出的 trace 工具的强大之处，它提供了指定时间内程序发生的事件的完整信息，这些事件信息包括：

- 协程的创建、开始和结束。
- 协程的堵塞——系统调用、通道、锁。
- 网络 I/O 相关事件。
- 系统调用事件。
- 垃圾回收相关事件。

收集 trace 文件的方式和收集 pprof 特征文件的方式非常相似，有两种主要的方式，一种是在程序中调用 runtime/trace 包的接口：

```go
package main

import (
    "log"
    "os"
    "runtime/trace"
)

func main() {
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    err = trace.Start(f)
    if err != nil {
        log.Fatal(err)
    }
    defer trace.Stop()

    // ... do some work ...
}
```

另一种方式仍然是使用 http 服务器，net/http/pprof 库中集成了 trace 的接口，下例获取 20s 内的 trace 事件并存储到 trace.out 文件中。

```bash
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=20
```

当要对获取的文件进行分析时，需要使用 trace 工具。

```bash
go tool trace trace.out
```

执行后会默认自动打开浏览器。

其中最复杂、信息最丰富的是 View trace。点击后会出现交互式的可视化界面，用于显示整个执行周期内的完整事件。

## 分析场景

### 分析延迟问题

当程序中至关重要的协程长时间无法运行时，可能带来延迟问题。发生这种情况的原因有很多，例如系统调用被堵塞、通道 / 互斥锁上被堵塞、协程被运行时代码（例如 GC）堵塞，甚至可能是调度器没有按照预期的频率运行关键协程。

这些问题都可以通过 trace 查看。查看逻辑处理器的时间线，并查找到关键的协程长时间被阻塞的时间段，查看这段时间内发生的事件，有助于查找到延迟问题的根源。

### 诊断不良并行性

在程序中应该保持适当的协程数和并行率。如果一个预期会使用所有 CPU 的程序运行速度比预期要慢，那么可能是因为程序没有按照期望的并行。可以查找程序中的关键路径是否有并发，如果没有并发，则查看是否可以让这些关键路径并发从而提高效率。

## trace 底层原理

[trace 底层原理](underlying_principle.md)

```go

```
