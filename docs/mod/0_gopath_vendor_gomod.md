---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "从 gopath 到 gomod 历程"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/mod/gopath_vendor_gomod"  # 设置网页永久链接
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

## GOPATH

### GOROOT 是什么

通常我们说安装 Go 语言，实际上安装的是 Go 编译器和 Go 标准库，二者位于同一个安装包中。

假如你在Windows上使用Installer安装的话，它们将会被默认安装到 `c:\go`目录下，该目录即 GOROOT 目录，里面保存了开发 GO 程序所需要的所有组件，比如编译器、标准库和文档等等。

同时安装程序还会自动帮你设置 GOROOT 环境变量，如下图所示：

[![DF0RfJ.png](https://s3.ax1x.com/2020/11/15/DF0RfJ.png)](https://imgchr.com/i/DF0RfJ)

另外，安装程序还会把 `c:\go\bin` 目录添加到系统的 `PATH` 环境变量中，如下图所示：

[![DF0fp9.png](https://s3.ax1x.com/2020/11/15/DF0fp9.png)](https://imgchr.com/i/DF0fp9)

该目录主要是 GO 语言开发包中提供的二进程可执行程序。

所以，GOROOT 实际上是指示 GO 语言安装目录的环境变量，属于 GO 语言顶级目录。

### GOPATH 是什么

安装完 Go 语言，接下来就要写自己的 Hello World 项目了。实际上 Go 语言项目是由一个或多个 package 组成的，这些 package 按照来源分为以下几种：

- 标准库
- 第三方库
- 项目私有库

其中标准库的 package 全部位于 GOROOT 环境变量指示的目录中，而第三方库和项目私有库都位于 GOPATH 环境变量所指示的目录中。

实际上，安装 GO 语言时，安装程序会设置一个默认的 GOPATH 环境变量，如下所示：

[![DF0bkD.png](https://s3.ax1x.com/2020/11/15/DF0bkD.png)](https://imgchr.com/i/DF0bkD)

与 GOROOT 不同的是，GOPATH 环境变量位于用户域，因为每个用户都可以创建自己的工作空间而互不干扰。
用户的项目需要位于 `GOPATH` 下的 `src/` 目录中。

所以 GOPATH 指示用户工作空间目录的环境变量，它属于用户域范畴的。

### 依赖查找

当某个 package 需要引用其他包时，编译器就会依次从 `GOROOT/src/` 和 `GOPATH/src/` 中去查找，如果某个包从 GOROOT 下找到的话，就不再到 GOPATH 目录下查找，所以如果项目中开发的包名与标准库相同的话，将会被自动忽略。

### GOPATH的缺点

GOPATH 的优点是足够简单，但它不能很好的满足实际项目的工程需求。

比如，你有两个项目 A 和 B，他们都引用某个第三方库 T，但这两个项目使用了不同的 T 版本，即：

- 项目A 使用T v1.0
- 项目B 使用T v2.0

由于编译器依赖查找固定从 GOPATH/src 下查找 `GOPATH/src/T`，所以，无法在同一个 GOPATH 目录下保存第三方库 T 的两个版本。所以项目 A、B 无法共享同一个 GOPATH，需要各自维护一个，这给广大软件工程师带来了极大的困扰。

针对 GOPATH 的缺点，GO 语言社区提供了 Vendor 机制，从此依赖管理进入第二个阶段：将项目的依赖包私有化。


vendor 机制提供了一个机制让项目的依赖隔离而不互相干扰。

自 Go 1.6 版本起，vendor 机制正式启用，它允许把项目的依赖全部放到一个位于本项目的 vendor 目录中，这个 vendor 目录可以简单理解成私有的 GOPATH 目录。即编译时，优先从 vendor 中寻找依赖包，如果 vendor 中找不到再到 GOPATH 中寻找。

## vendor

### vendor 目录位置

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

### 搜索顺序

上面的例子中，在编译main.go时，编译器搜索依赖包顺序为：

1. 从 `github.com/constabulary/example-gsftp/cmd/gsftp/` 下寻找 vendor 目录，没有找到，继续从上层查找；
2. 从 `github.com/constabulary/example-gsftp/cmd/` 下寻找 vendor 目录，没有找到，继续从上层查找；
3. 从 `github.com/constabulary/example-gsftp/` 下寻找 vendor 目录，从 vendor 目录中查找依赖包，结束；

如果 `github.com/constabulary/example-gsftp/` 下的 vendor 目录中没有依赖包，则返回到 GOPATH 目录继续查找，这就是前面介绍的 GOPATH 机制了。

从上面的搜索顺序可以看出，实际上 vendor 目录可以存在于项目的任意目录的。但非常不推荐这么做，因为如果 vendor 目录过于分散，很可能会出现同一个依赖包，在项目的多个 vendor 中出现多次，这样依赖包会多次编译进二进制文件，从而造成二进制大小急剧变大。同时，也很可能出现一个项目中使用同一个依赖包的多个版本的情况，这种情况往往应该避免。

### vendor 存在的问题

vendor 很好的解决了多项目间的隔离问题，但是位于 vendor 中的依赖包无法指定版本，某个依赖包，在把它放入 vendor 的那刻起，它就固定在当时版本，项目的使用者很难识别出你所使用的依赖版本。

比起这个，更严重的问题是上面提到的二进制急剧扩大问题，比如你依赖某个开源包 A 和 B，但 A 中也有一个 vendor 目录，其中也放了 B，那么你的项目中将会出现两个开源包 B。再进一步，如果这两个开源包 B 版本不一致呢？如果二者不兼容，那后果将是灾难性的。

但是，不得不说，vendor 能够解决绝大部分项目中的问题，如果你项目在使用 vendor，也绝对没有问题。一直到 Go 1.11 版本，官方社区推出了 Modules 机制，从此 Go 的版本管理走进第三个时代。

## gomod

在 Go v1.11 版本中，Module 特性被首次引入，这标志着 Go 的依赖管理开始进入第三个阶段。

Go Module 相比 GOPATH 和 vendor 而言功能强大得多，它基本上完全解决了 GOPATH 和 vendor 时代遗留的问题。

我们知道，GOPATH 时代最大的困扰是无法让多个项目共享同一个 pakage 的不同版本，在 vendor 时代，通过把每个项目依赖的 package 放到 vendor 中可以解决这个困扰，但是使用 vendor 的问题是无法很好的管理依赖的 package，比如升级 package。

虽然 Go Module 能够解决 GOPATH 和 vendor 时代遗留的问题，但需要注意的是 Go Module 不是 GOPATH 和 vendor 的演进，理解这个对于接下来正确理解 Go Module 非常重要。

Go Module 更像是一种全新的依赖管理方案，它涉及一系列的特性，但究其核心，它主要解决两个重要的问题：

- 准确的记录项目依赖；
- 可重复的构建；

准确的记录项目依赖，是指你的项目依赖哪些 package、以及 package 的版本可以非常精确。比如你的项目依赖 `github.com/prometheus/client_golang`，且必须是 `v1.0.0` 版本，那么你可以通过 Go Module 指定（具体指定方法后面会介绍），任何人在任何环境下编译你的项目，都必须要使用 `github.com/prometheus/client_golang` 的 `v1.0.0` 版本。

可重复的构建是指，项目无论在谁的环境中（同平台）构建，其产物都是相同的。回想一下 GOPATH 时代，虽然大家拥有同一个项目的代码，但由于各自的 GOPATH 中 `github.com/prometheus/client_golang` 版本不一样，虽然项目可以构建，但构建出的可执行文件很可能是不同的。可重复构建至关重要，避免出现 “ 我这运行没问题，肯定是你环境问题 ” 等类似问题出现。

一旦项目的依赖被准确记录了，就很容易做到重复构建。
