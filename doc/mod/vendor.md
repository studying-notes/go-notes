---
date: 2020-11-15T20:37:25+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "理解 vendor 机制"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/doc/mod/vendor"  # 设置网页永久链接
tags: [ "go", "gomod" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

vendor 机制提供了一个机制让项目的依赖隔离而不互相干扰。

自 Go 1.6 版本起，vendor 机制正式启用，它允许把项目的依赖全部放到一个位于本项目的 vendor 目录中，这个 vendor 目录可以简单理解成私有的 GOPATH 目录。即编译时，优先从 vendor 中寻找依赖包，如果 vendor 中找不到再到 GOPATH 中寻找。

## vendor 目录位置

一个项目可以有多个 vendor 目录，分别位于不同的目录级别，但建议每个项目只在根目录放置一个 vendor 目录。

假如你有一个`github.com/constabulary/example-gsftp`项目，项目目录结构如下：

```
$GOPATH
|	src/
|	|	github.com/constabulary/example-gsftp/
|	|	|	cmd/
|	|	|	|	gsftp/
|	|	|	|	|	main.go
```

其中 `main.go`中依赖如下几个包：

```go
import (
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)
```

在没有使用 vendor 目录时，若想编译这个项目，那么 GOPATH 目录结构应该是如下所示：

```
$GOPATH
|	src/
|	|	github.com/constabulary/example-gsftp/
|	|	golang.org/x/crypto/ssh
|	|	github.com/pkg/sftp
```

即，所有依赖的包，都位于 ` $ GOPATH/src` 下。

为了把所使用到的 `golang.org/x/crypto/ssh` 和 `github.com/pkg/sftp` 版本固化下来，那么可以使用 vendor 机制。

在项目 `github.com/constabulary/example-gsftp` 根目录下，创建一个 vendor 目录，并把 `golang.org/x/crypto/ssh` 和 `github.com/pkg/sftp` 存放到此处，让其成为项目的一部分。如下所示：

```
$GOPATH
|	src/
|	|	github.com/constabulary/example-gsftp/
|	|	|	cmd/
|	|	|	|	gsftp/
|	|	|	|	|	main.go
|	|	|	vendor/
|	|	|	|	github.com/pkg/sftp/
|	|	|	|	golang.org/x/crypto/ssh/
```

使用 vendor 的好处是在项目 `github.com/constabulary/example-gsftp` 发布时，把其所依赖的软件一并发布，编译时不会受到 GOPATH 目录的影响，即便 GOPATH 下也有一个同名但不同版本的依赖包。

## 搜索顺序

上面的例子中，在编译main.go时，编译器搜索依赖包顺序为：

1. 从 `github.com/constabulary/example-gsftp/cmd/gsftp/` 下寻找 vendor 目录，没有找到，继续从上层查找；
2. 从 `github.com/constabulary/example-gsftp/cmd/` 下寻找 vendor 目录，没有找到，继续从上层查找；
3. 从 `github.com/constabulary/example-gsftp/` 下寻找 vendor 目录，从 vendor 目录中查找依赖包，结束；

如果 `github.com/constabulary/example-gsftp/` 下的 vendor 目录中没有依赖包，则返回到 GOPATH 目录继续查找，这就是前面介绍的 GOPATH 机制了。

从上面的搜索顺序可以看出，实际上 vendor 目录可以存在于项目的任意目录的。但非常不推荐这么做，因为如果 vendor 目录过于分散，很可能会出现同一个依赖包，在项目的多个 vendor 中出现多次，这样依赖包会多次编译进二进制文件，从而造成二进制大小急剧变大。同时，也很可能出现一个项目中使用同一个依赖包的多个版本的情况，这种情况往往应该避免。

## vendor 存在的问题

vendor 很好的解决了多项目间的隔离问题，但是位于 vendor 中的依赖包无法指定版本，某个依赖包，在把它放入 vendor 的那刻起，它就固定在当时版本，项目的使用者很难识别出你所使用的依赖版本。

比起这个，更严重的问题是上面提到的二进制急剧扩大问题，比如你依赖某个开源包 A 和 B，但 A 中也有一个 vendor 目录，其中也放了 B，那么你的项目中将会出现两个开源包 B。再进一步，如果这两个开源包 B 版本不一致呢？如果二者不兼容，那后果将是灾难性的。

但是，不得不说，vendor 能够解决绝大部分项目中的问题，如果你项目在使用 vendor，也绝对没有问题。一直到 Go 1.11 版本，官方社区推出了 Modules 机制，从此 Go 的版本管理走进第三个时代。
