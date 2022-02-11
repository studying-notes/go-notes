---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 2"  # 文章标题
url:  "posts/go/docs/mod/2_module_quickstart"  # 设置网页永久链接
tags: [ "go", "gomod" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## Go module 到底是做什么的？

我们在前面的章节已介绍过，但还是想强调一下，Go module 实际上只是精准的记录项目的依赖情况，包括每个依赖的精确版本号，仅此而矣。

那么，为什么需要记录这些依赖情况，或者记录这些依赖有什么好处呢？

试想一下，在编译某个项目时，第三方包的版本往往是可以替换的，如果不能精确的控制所使用的第三方包的版本，最终构建出的可执行文件从本质上是不同的，这会给问题诊断带来极大的困扰。

接下来，我们从一个Hello World项目开始，逐步介绍如何初始化module、如何记录依赖的版本信息。

项目托管在GitHub `https://github.com/renhongcai/gomodule` 中，并使用版本号区别使用 go module 的阶段。

- v1.0.0 未引用任何第三方包，也未使用 go module
- v1.1.0 未引用任何第三方包，已开始使用 go module，但没有任何外部依赖
- v1.2.0 引用了第三方包，并更新了项目依赖

需要注意的是，下面的例子统一使用 go 1.13 版本，如果你使用 go 1.11 或者 go 1.12，运行效果可能略有不同。本文最后部分我们尽量尝试记录一些版本间的差异，以供参考。

## Hello World

在 v1.0.0 版本时，项目只包含一个 main.go 文件，只是一个简单的字符串打印：

```golang
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

此时，项目还没有引用任何第三方包，也未使用 go module。

## 初始化 module

一个项目若要使用 Go module，那么其本身需要先成为一个 module，也即需要一个 module 名字。

在 Go module 机制下，项目的 module 名字以及其依赖信息记录在一个名为 `go.mod` 的文件中，该文件可以手动创建，也可以使用 `go mod init` 命令自动生成。推荐自动生成的方法，如下所示：

```
[root@ecs-d8b6 gomodule]# go mod init github.com/renhongcai/gomodule
go: creating new go.mod: module github.com/renhongcai/gomodule
```

完整的 `go mod init` 命令格式为 `go mod init [module]`：

其中 `[module]` 为 module 名字，如果不填，`go mod init` 会尝试从版本控制系统或 import 的注释中猜测一个。这里推荐指定明确的 module 名字，因为猜测有时需要一些额外的条件，比如 Go 1.13 版本，只有项目位于 GOPATH 中才可以正确运行，而 Go 1.11 版本则没有此要求。

上面的命令会自动创建一个 `go.mod` 文件，其中包括 module 名字，以及我们所使用的 Go 版本：

```
[root@ecs-d8b6 gomodule]# cat go.mod
module github.com/renhongcai/gomodule

go 1.13
```

`go.mod` 文件中的版本号 `go 1.13` 是在 Go 1.12 引入的，意思是开发此项目的 Go 语言版本，并不是编译该项目所限制的 Go 语言版本，但是如果项目中使用了 Go 1.13 的新特性，而你使用 Go 1.11 编译的话，编译失败时，编译器会提示你版本不匹配。

由于我们的项目还没有使用任何第三方包，所以 `go.mod` 中并没有记录依赖包的任何信息。我们把自动生成的 `go.mod` 提交，然后我们尝试引用一个第三方包。

## 管理依赖

现在我们准备引用一个第三方包 `github.com/google/uuid` 来生成一个 UUID，这样就会产生一个依赖，代码如下：

```go
package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New().String()
	fmt.Println("UUID: ", id)
}
```

在开始编译以前，我们先使用 `go get` 来分析依赖情况，并会自动下载依赖：

```
[root@ecs-d8b6 gomodule]# go get
go: finding github.com/google/uuid v1.1.1
go: downloading github.com/google/uuid v1.1.1
go: extracting github.com/google/uuid v1.1.1
```

从输出内容来看，`go get` 帮我们定位到可以使用 `github.com/google/uuid` 的 v1.1.1 版本，并下载再解压它们。

注意：`go get` 总是获取依赖的最新版本，如果 `github.com/google/uuid` 发布了新的版本，输出的版本信息会相应的变化。关于 Go Module 机制中版本选择我们将在后续的章节详细介绍。

`go get` 命令会自动修改 `go.mod`文件：

```
[root@ecs-d8b6 gomodule]# cat go.mod
module github.com/renhongcai/gomodule

go 1.13

require github.com/google/uuid v1.1.1
```

可以看到，现在 `go.mod` 中增加了 `require github.com/google/uuid v1.1.1` 内容，表示当前项目依赖 `github.com/google/uuid` 的 `v1.1.1` 版本，这就是我们所说的 `go.mod` 记录的依赖信息。

由于这是当前项目第一次引用外部依赖，`go get` 命令还会生成一个 `go.sum` 文件，记录依赖包的 hash 值：

```
[root@ecs-d8b6 gomodule]# cat go.sum
github.com/google/uuid v1.1.1 h1:Gkbcsh/GbpXz7lPftLA3P6TYMwjCLYm83jiFQZF/3gY=
github.com/google/uuid v1.1.1/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

该文件通过记录每个依赖包的 hash 值，来确保依赖包没有被篡改。关于此部分内容我们在此暂不展开介绍，留待后面的章节详细介绍。

经 `go get` 修改的 `go.mod` 和创建的 `go.sum` 都需要提交到代码库，这样别人获取到项目代码，编译时就会使用项目所要求的依赖版本。

至此，项目已经有一个依赖包，并且可以编译执行了，每次运行都会生成一个独一无二的 UUID：

```
[root@ecs-d8b6 gomodule]# go run main.go
UUID:  20047f5a-1a2a-4c00-bfcd-66af6c67bdfb
```

注：如果你没有使用 `go get` 在执行之前下载依赖，而是直接使用 `go build main.go` 运行项目的话，依赖包也会被自动下载。但是在 `v1.13.4` 中有个 bug，即此时生成的 `go.mod` 中显示的依赖信息则会是`require github.com/google/uuid v1.1.1 // indirect`，注意行末的 `indirect` 表示间接依赖，这明显是错误的，因为我们直接 `import` 的。

## 版本间差异

由于 Go module 在 Go v1.11 初次引入，历经 Go v1.12、v1.13 的发展，其实现细节上已有了一些变化，按照之前的规划 Go module 将会在 v1.14 定型，推荐尽可能使用最新版本，否则可能会产生一些困扰。

比如，在v1.11 中使用 `go mod init` 初始化项目时，不填写 `module` 名称是没有问题，但在 v1.13 中，如果项目不在 GOPATH 目录中，则必须填写 `module` 名称。
