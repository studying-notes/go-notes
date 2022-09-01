---
date: 2022-09-01T09:00:33+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "strings.Title 被废弃"  # 文章标题
url:  "posts/go/quickstart/feature/strings/title"  # 设置网页永久链接
tags: [ "Go", "title" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 功能

strings.Title 会将每个单词的首字母变成大写字母。

strings 中还有一个函数：ToTitle，它的作用和 ToUpper 类似，所有字符全部变成大写，而不只是首字母。不过 ToTitle 和 ToUpper 的区别特别微小，是 Unicode 相关的规定。

## 原因

strings.Title 的规则是使用单词边界，不能正确处理 Unicode 标点。

```go
fmt.Println(strings.Title("here comes o'brian"))
```

期望输出：`Here Comes O'brian`，但 strings.Title 的结果是：`Here Comes O'Brian`。

## 继任者

在 strings.Title 中提到，可以使用 `golang.org/x/text/cases` 代替 strings.Title，具体来说就是 cases.Title。

```bash
go get golang.org/x/text/cases 
```

该包提供了通用和特定于语言的 case map，其中有一个 Title 函数，签名如下：

```go
func Title(t language.Tag, opts ...Option) Caser
```

第一个参数是 language.Tag 类型，表示 BCP 47 种语言标记。它用于指定特定语言或区域设置的实例。所有语言标记值都保证格式良好。

第二个参数是不定参数，类型是 Option，这是一个函数类型：

```go
type Option func(o options) options
```

它被用来修改 Caser 的行为，cases 包可以找到相关 Option 的实例。

cases.Title 的返回类型是 Caser，这是一个结构体，这里我们只关心它的 String 方法，它接收一个字符串，并返回一个经过 Caser 处理过后的字符串。

所以，针对上文 strings.Title 的场景，可以改为 cases.Title 实现。

```go
package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	fmt.Println(strings.Title("Golang is awesome!"))

	caser := cases.Title(language.English)
	fmt.Println(caser.String("here comes o'brian"))
}
```

```go

```
