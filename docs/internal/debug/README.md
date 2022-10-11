---
date: 2022-10-11T11:36:25+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "特征分析与事件追踪"  # 文章标题
url:  "posts/go/docs/internal/debug/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

Go 语言原生支持对于程序运行时重要指标或特征的分析，这种支持体现在：

- 根据程序编译参数与运行参数的不同，可以进行不同类型和水平的调试。
- 官方提供了许多调试库，可用于程序的调试。
- 运行时可以保留重要的特征指标和状态，有许多工具可以分析甚至可视化程序运行的状态和过程。

Go 语言中的 pprof 指对于指标或特征的分析（Profiling），通过分析不仅可以查找到程序中的错误（内存泄漏、race 冲突、协程泄漏），也能对程序进行优化（例如 CPU 利用率不足）。

由于 Go 语言运行时的指标不对外暴露，因此有标准库 net/http/pprof 和 runtime/pprof 用于与外界交互。

其中 net/http/pprof 提供了一种通过 http 访问的便利方式，用于用户调试和获取样本特征数据。对特征文件进行分析要依赖谷歌推出的分析工具 pprof，该工具在 Go 安装时即存在。

## pprof 的使用方式

[pprof 的使用方式](pprof/README.md)

## trace 事件追踪

[trace 事件追踪](trace/README.md)

```go

```
