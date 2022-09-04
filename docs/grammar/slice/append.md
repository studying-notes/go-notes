---
date: 2020-11-14T22:21:59+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 切片 append 方法的陷阱"  # 文章标题
url:  "posts/go/docs/grammar/slice/append"  # 设置网页永久链接
tags: [ "Go", "append" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 函数原型

```go
func append(slice []Type, elems ...Type) []Type
```

## 会改变切片的地址

`append` 的本质是向切片中追加数据，而随着切片中元素逐渐增加，当切片底层的数组将满时，切片会发生扩容，扩容会导致产生一个新的切片（拥有容量更大的底层数组）。

`append` 每个追加元素，都有可能触发切片扩容，也即有可能返回一个新的切片，这也是 `append` 函数声明中返回值为切片的原因。实际使用中应该总是接收该返回值。

不管初始切片长度为多少，不接收 `append` 返回都是有极大风险的。

## 可以追加 nil 值

向切片中追加一个 `nil` 值是完全不会报错的，如下代码所示：

```
slice := append(slice, nil)
```

经过追加后，slice 的长度递增 1。

实际上 `nil` 是一个预定义的值，即空值，所以完全有理由向切片中追加。

```go

```
