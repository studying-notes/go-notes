---
date: 2022-09-01T09:07:26+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "字符串分割 Cut"  # 文章标题
url:  "posts/go/quickstart/feature/strings/cut"  # 设置网页永久链接
tags: [ "Go", "cut" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 功能

在此之前，strings.Index 系列函数常用于字符串的切割。

Cut 函数的签名如下：

```go
func Cut(s, sep string) (before, after string, found bool)
```

将字符串 s 在第一个 sep 处切割为两部分，分别存在 before 和 after 中。如果 s 中没有 sep，返回 s,"",false。

## 实现

```go
func Cut(s, sep string) (before, after string, found bool) {
	if i := Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
```

## 示例

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	addr := "192.168.1.1:8080"
	ip, port, ok := strings.Cut(addr, ":")

	if ok {
		fmt.Printf("ip: %s, port: %s\n", ip, port)
	}
}
```

```go

```
