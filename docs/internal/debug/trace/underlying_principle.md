---
date: 2022-10-11T14:45:09+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "trace 底层原理"  # 文章标题
url:  "posts/go/docs/internal/debug/trace/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

即便使用 net/http/pprof 包，底层仍然会调用 runtime/trace 功能。在 trace 的初始阶段需要首先 STW，然后获取协程的快照、状态、栈帧信息，接着设置 trace.enable = true 开启 GC，最后重新启动所有协程，如图 21-17 所示。

trace 提供了强大的内省功能，这种功能不是没有代价的，Go 语言在运行时源码中每个重要的事件处都加入了判断 trace.enabled 是否开启的条件，并编译到了程序中，当 trace 开启后，会触发 traceEvent 写入事件。

这些关键的事件包括协程的生命周期、协程堵塞、网络I/O、系统调用、垃圾回收等，根据事件的不同，可能保存和此事件相关的不同数量的参数及栈追踪数据。每个逻辑处理器P都有一个缓存（p.tracebuf），用于存储已经被序列化为字节的事件（Event），如图21-18所示。

![](../../../../assets/images/docs/internal/debug/trace/underlying_principle/图21-17%20tracec处理流程.png)

![](../../../../assets/images/docs/internal/debug/trace/underlying_principle/图21-18%20trace事件存储到逻辑处理器P的缓存中.png)

版本、时间戳、栈 ID、协程 ID 等整数信息使用 LEB128 编码，用于有效压缩数字的长度。字符串使用 UFT-8 编码。

每个逻辑处理器 P 的缓存都是有限度的，当超过了缓存限度后，逻辑处理器 P 中的 tracebuf 会转移到全局链表中，如图 21-19 所示。

同时，trace 工具会新开一个协程专门用于读取全局 trace 上的信息，此时全局的事件对象已经是序列化之后的字节数组，直接添加到文件中即可。另外，访问全局 trace 缓存需要加锁，当没有可以访问的对象时，读取协程会陷入休眠状态。

![](../../../../assets/images/docs/internal/debug/trace/underlying_principle/图21-19%20逻辑处理器P中的缓存溢出到全局链表中.png)

当指定的时间到期后，需要结束 trace 任务，程序会再次陷入 STW 状态，刷新逻辑处理器 P 上的 tracebuf 缓存，设置 trace.enabled = false，从而完成整个 trace 收集周期。

当完成收集工作并存储到文件后，go tool trace 完成对 trace 文件的解析并开启 http 服务供浏览器访问，在 Go 源码中可以看到具体的解析过程。trace 的 web 界面来自 trace-viewer 项目， trace-viewer 可以从多种事件格式中生成可视化效果，go tool trace 中使用了基于 JSON 的事件格式。

```go

```
