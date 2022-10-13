---
date: 2022-09-09T09:17:43+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言底层原理剖析"  # 文章标题
url:  "posts/go/docs/internal/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 基础知识

1. [了解 CPU 寄存器](base/register.md)
2. [了解内存](base/memory.md)
3. [了解汇编指令](base/assembly.md)
4. [了解 Plan9 汇编](base/assembly_plan9.md)

## Go 语言底层原理剖析

1. [深入 Go 语言编译器](compiler/README.md)
2. [浮点数设计原理](float/README.md)
3. [类型推断](type_inference/README.md)
4. [常量与隐式类型转换](constant/README.md)
5. [字符串本质与实现](string/README.md)
6. [数组](array/README.md)
7. [切片使用方法与底层原理](slice/README.md)
8. [哈希表](map/README.md)
9. [函数与栈](function/README.md)
10. [defer延迟调用](defer/README.md)
11. [异常与异常捕获](panic/README.md)
12. [接口与程序设计模式](interface/README.md)
13. [反射](reflect/README.md)
14. [协程初探](goroutine/README.md)
15. [深入协程设计与调度原理](goroutine/design/README.md)
16. [通道与协程间通信](goroutine/channel/README.md)
17. 并发控制
    1.  [context 处理协程退出](goroutine/context/README.md)
    2. [数据争用检查](goroutine/race/README.md)
    3. [锁](goroutine/lock/README.md)
18. [内存分配管理](memory/README.md)
    - [Go 内存分配](memory/alloc.md)
    - [Go 内存逃逸分析](memory/escape.md)
    - [Go 垃圾回收](memory/gc.md)
19. [垃圾回收初探](gc/README.md)
20. [深入垃圾回收全流程](gc/underlying_principle.md)
21. [特征分析与事件追踪](debug/README.md)
