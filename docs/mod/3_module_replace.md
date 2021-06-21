---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 3"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/mod/3_module_replace"  # 设置网页永久链接
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

`go.mod` 文件中通过`指令`声明 module 信息，用于控制命令行工具进行版本选择。一共有四个指令可供使用：

- module：声明 module 名称；
- require：声明依赖以及其版本号；
- replace：替换 require 中声明的依赖，使用另外的依赖及其版本号；
- exclude：禁用指定的依赖；

其中 `module` 和 `require` 我们前面已介绍过，`module` 用于指定 module 的名字，如`module github.com/renhongcai/gomodule`，那么其他项目引用该 module 时其 import 路径需要指定`github.com/renhongcai/gomodule`。`require` 用于指定依赖，如`require github.com/google/uuid v1.1.1`，该指令相当于告诉 `go build` 使用`github.com/google/uuid`的 `v1.1.1` 版本进行编译。

## replace 工作机制

顾名思义，`replace` 指替换，它指示编译工具替换 `require` 指定中出现的包，比如，我们在 `require` 中指定的依赖如下：

```
module github.com/renhongcai/gomodule  
  
go 1.13  
  
require github.com/google/uuid v1.1.1
```

此时，我们可以使用 `go list -m all` 命令查看最终选定的版本：

```
[root@ecs-d8b6 gomodule]# go list -m all
github.com/renhongcai/gomodule
github.com/google/uuid v1.1.1
```

毫无意外，最终选定的 uuid 版本正是我们在 require 中指定的版本 `v1.1.1`。

如果我们想使用 uuid 的 v1.1.0 版本进行构建，可以修改 require 指定，还可以使用 replace 来指定。

需要说明的是，正常情况下不需要使用 replace 来修改版本，最直接的办法是修改 require 即可，虽然 replace 也能够做到，但这不是 replace 的一般使用场景。

下面我们先通过一个简单的例子来说明 replace 的功能，随即介绍几种常见的使用场景。

比如，我们修改 `go.mod`，添加replace指令：

```
[root@ecs-d8b6 gomodule]# cat go.mod 
module github.com/renhongcai/gomodule

go 1.13

require github.com/google/uuid v1.1.1

replace github.com/google/uuid v1.1.1 => github.com/google/uuid v1.1.0
```

`replace github.com/google/uuid v1.1.1 => github.com/google/uuid v1.1.0`指定表示替换 uuid v1.1.1 版本为 v1.1.0，此时再次使用 `go list -m all` 命令查看最终选定的版本：

```
[root@ecs-d8b6 gomodule]# go list -m all 
github.com/renhongcai/gomodule
github.com/google/uuid v1.1.1 => github.com/google/uuid v1.1.0
```

可以看到其最终选择的 uuid 版本为 v1.1.0。如果你本地没有 v1.1.0 版本，你或许还会看到一条`go: finding github.com/google/uuid v1.1.0`信息，它表示在下载 uuid v1.1.0 包，也从侧面证明最终选择的版本为 v1.1.0。

到此，我们可以看出 `replace` 的作用了，它用于替换 `require` 中出现的包，它正常工作还需要满足两个条件：

第一，`replace` 仅在当前 module 为 `main module` 时有效，比如我们当前在编译`github.com/renhongcai/gomodule`，此时就是 `main module`，如果其他项目引用了`github.com/renhongcai/gomodule`，那么其他项目编译时，`replace` 就会被自动忽略。

第二，`replace` 指定中 ` = >` 前面的包及其版本号必须出现在 `require` 中才有效，否则指令无效，也会被忽略。
比如，上面的例子中，我们指定 `replace github.com/google/uuid => github.com/google/uuid v1.1.0`，或者指定 `replace github.com/google/uuid v1.0.9 => github.com/google/uuid v1.1.0`，二者均都无效。

## replace 使用场景

前面的例子中，我们使用 `replace` 替换 `require` 中的依赖，在实际项目中 `replace` 在项目中经常被使用，其中不乏一些精彩的用法。

但不管应用在哪种场景，其本质都一样，都是替换 `require` 中的依赖。

#### 替换无法下载的包

由于中国大陆网络问题，有些包无法顺利下载，比如 `golang.org` 组织下的包，值得庆幸的是这些包在 GitHub 都有镜像，此时就可以使用 GitHub 上的包来替换。

比如，项目中使用了 `golang.org/x/text` 包：

```go
package main

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	id := uuid.New().String()
	fmt.Println("UUID: ", id)

	p := message.NewPrinter(language.BritishEnglish)
	p.Printf("Number format: %v.\n", 1500)

	p = message.NewPrinter(language.Greek)
	p.Printf("Number format: %v.\n", 1500)
}
```

上面的简单例子，使用两种语言 `language.BritishEnglish` 和 `language.Greek` 分别打印数字 `1500`，来查看不同语言对数字格式的处理，一个是 `1,500`，另一个是 `1.500`。此时就会分别引入 `"golang.org/x/text/language"` 和 `"golang.org/x/text/message"`。

执行 `go get` 或 `go build` 命令时会就再次分析依赖情况，并更新 `go.mod` 文件。网络正常情况下，`go.mod` 文件将会变成下面的内容：

```
module github.com/renhongcai/gomodule

go 1.13

require (
	github.com/google/uuid v1.1.1
	golang.org/x/text v0.3.2
)

replace github.com/google/uuid v1.1.1 => github.com/google/uuid v1.1.0
```

我们看到，依赖 `golang.org/x/text` 被添加到了 require 中。（多条 require 语句会自动使用 `()` 合并）。此外，我们没有刻意指定 `golang.org/x/text` 的版本号，Go 命令行工具根据默认的版本计算规则使用了 v0.3.2 版本，此处我们暂不关心具体的版本号。

没有合适的网络代理情况下，`golang.org/x/text` 很可能无法下载。那么此时，就可以使用 `replace` 来让我们的项目使用 GitHub 上相应的镜像包。我们可以添加一条新的 `replace` 条目，如下所示：

```
replace (
	github.com/google/uuid v1.1.1 => github.com/google/uuid v1.1.0
	golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
)
```

此时，项目编译时就会从 GitHub 下载包。我们源代码中 import 路径 `golang.org/x/text/xxx` 不需要改变。

也许有读者会问，是否可以将 import 路径由 `golang.org/x/text/xxx` 改成 `github.com/golang/text/xxx`，这样一来，就不需要使用 replace 来替换包了。

遗憾的是，不可以。因为 `github.com/golang/text` 只是镜像仓库，其 `go.mod` 文件中定义的module还是 `module golang.org/x/text`，这个 module 名字直接决定了你的 import 的路径。

#### 调试依赖包

有时我们需要调试依赖包，此时就可以使用 `replace` 来修改依赖，如下所示：

```
replace (
github.com/google/uuid v1.1.1 => ../uuid
golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
)
```

语句 `github.com/google/uuid v1.1.1 => ../uuid` 使用本地的 uuid 来替换依赖包，此时，我们可以任意地修改 `../uuid` 目录的内容来进行调试。

除了使用相对路径，还可以使用绝对路径，甚至还可以使用自已的 fork 仓库。

#### 使用 fork 仓库

有时在使用开源的依赖包时发现了 bug，在开源版本还未修改或者没有新的版本发布时，你可以使用 fork 仓库，在 fork 仓库中进行 bug fix。
你可以在 fork 仓库上发布新的版本，并相应的修改 `go.mod` 来使用 fork 仓库。

比如，我 fork 了开源包 `github.com/google/uuid`，fork 仓库地址为 `github.com/RainbowMango/uuid`，那我们就可以在 fork 仓库里修改 bug 并发布新的版本 `v1.1.2`，此时使用 fork 仓库的项目中 `go.mod` 中 replace 部分可以相应的做如下修改：
```
github.com/google/uuid v1.1.1 => github.com/RainbowMango/uuid v1.1.2
```

需要说明的是，使用 fork 仓库仅仅是临时的做法，一旦开源版本变得可用，需要尽快切换到开源版本。

#### 禁止被依赖

另一种使用 `replace` 的场景是你的 module 不希望被直接引用，比如开源软件[kubernetes](https://github.com/kubernetes/kubernetes)，在它的 `go.mod` 中 `require` 部分有大量的 `v0.0.0` 依赖，比如：

```
module k8s.io/kubernetes

require (
	...
	k8s.io/api v0.0.0
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0
	k8s.io/cli-runtime v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/cloud-provider v0.0.0
	...
)
```

由于上面的依赖都不存在 v0.0.0 版本，所以其他项目直接依赖`k8s.io/kubernetes`时会因无法找到版本而无法使用。
因为 Kubernetes 不希望作为 module 被直接使用，其他项目可以使用 kubernetes 其他子组件。

kubernetes 对外隐藏了依赖版本号，其真实的依赖通过 `replace` 指定：

```
replace (
	k8s.io/api => ./staging/src/k8s.io/api
	k8s.io/apiextensions-apiserver => ./staging/src/k8s.io/apiextensions-apiserver
	k8s.io/apimachinery => ./staging/src/k8s.io/apimachinery
	k8s.io/apiserver => ./staging/src/k8s.io/apiserver
	k8s.io/cli-runtime => ./staging/src/k8s.io/cli-runtime
	k8s.io/client-go => ./staging/src/k8s.io/client-go
	k8s.io/cloud-provider => ./staging/src/k8s.io/cloud-provider
)
```

前面我们说过，`replace` 指令在当前模块不是 `main module ` 时会被自动忽略的，Kubernetes 正是利用了这一特性来实现对外隐藏依赖版本号来实现禁止直接引用的目的。
