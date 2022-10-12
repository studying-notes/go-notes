---
date: 2022-10-12T13:33:27+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "遍历函数"  # 文章标题
url:  "posts/go/docs/internal/compiler/walk"  # 设置网页永久链接
tags: [ "Go", "walk" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 源码分布

这一阶段的相关源码分布在：

* `cmd/compile/internal/walk` (order of evaluation, desugaring)

## 简介

IR 表示的最后一次传递是 “walk”，有两个目的：

1. 将复杂的语句分解为单独的、更简单的语句，引入临时变量并尊重评估顺序。

2. 将更高级别的 Go 结构简化为更原始的结构。例如，`switch` 语句变成了二分查找或跳转表，map 和 channel 上的操作被替换为运行时调用。

在该阶段会识别出声明但是并未被使用的变量，遍历函数中的声明和表达式，将某些代表操作的节点转换为运行时的具体函数执行。例如，获取 map 中的值会被转换为运行时 mapaccess2_fast64 函数。

```go
val, ok := m["key"]
// 转化为
autotmp_1, ok := runtime.mapaccess2_fast64(typeOf(m), m, "key")
val := *autotmp_1
```

字符串变量的拼接会被转换为调用运行时 concatstrings 函数。对于 new 操作，如果变量发生了逃逸，那么最终会调用运行时 newobject 函数将变量分配到堆区。for...range 语句会重写为更简单的 for 语句形式。

在执行 walk 函数遍历之前，编译器还需要对某些表达式和语句进行重新排序，例如将 `x/= y` 替换为 `x = x/y`。根据需要引入临时变量，以确保形式简单，例如 `x = m[k]` 或 `m[k] = x`，而 k 可以寻址。

```go

```
