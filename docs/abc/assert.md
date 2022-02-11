---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 断言与类型转换"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/abc/assert"  # 设置网页永久链接
tags: [ "go", "assert" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 直接断言

在确信类型是正确的情况下可以直接断言：

```go
t := i.(T)
```

这个表达式可以断言一个接口对象 `i` 里不是 `nil`，并且接口对象 `i` 存储的值的类型是 `T`，如果断言成功，就会返回值给 `t`，如果断言失败，就会触发 `panic`。

## 判别断言

通过判别防止触发 `panic`。

```go
t, ok:= i.(T)
```

断言成功就会返回其类型给 `t`，并且此时 `ok` 的值为 `true`，表示断言成功；断言失败，不会触发 `panic`，而是将 `ok` 的值设为 `false` ，表示断言失败，此时 `t` 为 `T` 的零值。

## 类型区分

通过 Type Switch 断言比一个一个进行类型断言更简单、直接、高效。

```go
func assertType(i interface{}) {
    switch x := i.(type) {
    case int:
        fmt.Println(x, "is int")
    case string:
        fmt.Println(x, "is string")
    case nil:
        fmt.Println(x, "is nil")
    default:
        fmt.Println(x, "not type matched")
    }
}
```
