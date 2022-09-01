---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 模块深入讲解 2"  # 文章标题
url:  "posts/go/docs/module/2_module_quickstart"  # 设置网页永久链接
tags: [ "go" ]
categories: [ "Go 学习笔记" ]

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

Go module 的作用是精准记录项目的依赖情况，包括每个依赖的精确版本号。

在编译某个项目时，第三方包的版本往往是可以替换的，如果不能精确的控制所使用的第三方包的版本，最终构建出的可执行文件从本质上是不同的。

## Hello World

项目只包含一个 main.go 文件，进行一个简单的字符串打印：

```golang
package main

import "fmt"

func main() {
	fmt.Println("Golang is awesome!")
}
```

## 初始化 module

一个项目若要使用 Go module，那么其本身需要先成为一个 module，也即需要一个 module 名字。

在 Go module 机制下，项目的 module 名字以及其依赖信息记录在一个名为 `go.mod` 的文件中，该文件可以手动创建，也可以使用 `go mod init` 命令自动生成。

```
go mod init example
```

完整的 `go mod init` 命令格式为 `go mod init [module]`：

其中 `[module]` 为 module 名字。

上面的命令会自动创建一个 `go.mod` 文件，其中包括 module 名字，以及我们所使用的 Go 版本：

```
module example

go 1.18
```

`go.mod` 文件中的版本号 `go 1.18` 是在 Go 1.12 引入的，意思是开发此项目的 Go 语言版本，并不是编译该项目所限制的 Go 语言版本，但是如果项目中使用了 Go 1.13 的新特性，而你使用 Go 1.11 编译的话，编译失败时，编译器会提示你版本不匹配。

由于我们的项目还没有使用任何第三方包，所以 `go.mod` 中并没有记录依赖包的任何信息。

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

在开始编译以前，我们先使用 `go get` 来分析依赖情况，并会自动下载依赖。

`go get` 总是获取依赖的最新版本，如果 `github.com/google/uuid` 发布了新的版本，输出的版本信息会相应的变化。

`go get` 命令会自动修改 `go.mod`文件：

```
module example

go 1.18

require github.com/google/uuid v1.3.0
```

可以看到，现在 `go.mod` 中增加了 `require github.com/google/uuid v1.3.0` 内容，表示当前项目依赖 `github.com/google/uuid` 的 `v1.3.0` 版本，这就是我们所说的 `go.mod` 记录的依赖信息。

由于这是当前项目第一次引用外部依赖，`go get` 命令还会生成一个 `go.sum` 文件，记录依赖包的 hash 值：

```
github.com/google/uuid v1.3.0 h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=
github.com/google/uuid v1.3.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
```

该文件通过记录每个依赖包的 hash 值，来确保依赖包没有被篡改。

经 `go get` 修改的 `go.mod` 和创建的 `go.sum` 都需要提交到代码库，这样别人获取到项目代码，编译时就会使用项目所要求的依赖版本。

如果你没有使用 `go get` 在执行之前下载依赖，而是直接使用 `go build main.go` 运行项目的话，依赖包也会被自动下载。
