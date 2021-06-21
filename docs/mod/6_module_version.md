---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 6"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/mod/6_module_version"  # 设置网页永久链接
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

在前面的章节中，我们使用过 `go get <pkg>` 来获取某个依赖，如果没有特别指定依赖的版本号，`go get` 会自动选择一个最优版本，并且如果本地有 `go.mod` 文件的话，还会自动更新 `go.mod` 文件。

事实上除了 `go get`，`go build` 和 `go mod tidy` 也会自动帮我们选择依赖的版本。这些命令选择依赖版本时都遵循一些规则，本节我们就开始介绍 Go module 涉及到的版本选择机制。

## 依赖包版本约定

关于如何管理依赖包的版本，Go 语言提供了一个规范，并且 Go 语言的演进过程中也一直遵循这个规范。

这个非强制性的规范主要围绕包的兼容性展开。对于如何处理依赖包的兼容性，根据是否支持 Go module 分别有不同的建议。

#### Go module 之前版本兼容性

在 Go v1.11（开始引入 Go module 的版本）之前，Go 语言建议依赖包需要保持向后兼容，这包括可导出的函数、变量、类型、常量等不可以随便删除。以函数为例，如果需要修改函数的入参，可以增加新的函数而不是直接修改原有的函数。

如果确实需要做一些打破兼容性的修改，建议创建新的包。

比如仓库 `github.com/RainbowMango/xxx` 中包含一个 package A，此时该仓库只有一个 package：

- `github.com/RainbowMango/xxx/A`

那么其他项目引用该依赖时的 import 路径为：

```
import "github.com/RainbowMango/xxx/A"
```

如果该依赖包需要引入一个不兼容的特性，可以在该仓库中增加一个新的 package A1，此时该仓库包含两个包：

- `github.com/RainbowMango/xxx/A`
- `github.com/RainbowMango/xxx/A1`

那么其他项目在升级依赖包版本后不需要修改原有的代码可以继续使用 package A，如果需要使用新的 package A1，只需要将 import 路径修改为 `import "github.com/RainbowMango/xxx/A1"` 并做相应的适配即可。

#### Go module 之后版本兼容性

从 Go v1.11 版本开始，随着 Go module 特性的引入，依赖包的兼容性要求有了进一步的延伸，Go module 开始关心依赖包版本管理系统（如 Git）中的版本号。尽管如此，兼容性要求的核心内容没有改变：

- 如果新 package 和旧的 package 拥有相同的 import 路径，那么新 package 必须兼容旧的 package ;
- 如果新的 package 不能兼容旧的 package，那么新的 package 需要更换 import 路径；

在前面的介绍中，我们知道 Go module 的 `go.mod` 中记录的 module 名字决定了 import 路径。例如，要引用 module `module github.com/renhongcai/indirect`中的内容时，其import路径需要为`import github.com/renhongcai/indirect`。

在Go module 时代，module 版本号要遵循语义化版本规范，即版本号格式为`v<major>.<minor>.<patch>`，如v1.2.3。当有不兼容的改变时，需要增加 `major` 版本号，如 v2.1.0。

Go module 规定，如果 `major` 版本号大于 `1`，则 `major` 版本号需要显式地标记在 module 名字中，如`module github.com/my/mod/v2`。这样做的好处是 Go module 会把`module github.com/my/mod/v2` 和 `module github.com/my/mod`视做两个 module，他们甚至可以被同时引用。

另外，如果 module 的版本为 `v0.x.x` 或 `v1.x.x` 则都不需要在 module 名字中体现版本号。

## 版本选择机制

Go 的多个命令行工具都有自动选择依赖版本的能力，如 `go build` 和 `go test`，当在源代码中增加了新的 import，这些命令将会自动选择一个最优的版本，并更新 `go.mod` 文件。

需要特别说明的是，如果 `go.mod` 文件中已标记了某个依赖包的版本号，则这些命令不会主动更新 `go.mod` 中的版本号。所谓自动更新版本号只在 `go.mod` 中缺失某些依赖或者依赖不匹配时才会发生。

#### 最新版本选择

当在源代码中新增加了一个 import，比如：

```
import "github.com/RainbowMango/M"
```

如果`go.mod`的require指令中并没有包含`github.com/RainbowMango/M`这个依赖，那么`go build` 或`go test`命令则会去`github.com/RainbowMango/M`仓库寻找最新的符合语义化版本规范的版本，比如v1.2.3，并在`go.mod`文件中增加一条require依赖：

```
require github.com/RainbowMango/M v1.2.3
```

这里，由于import路径里没有类似于`v2`或更高的版本号，所以版本选择时只会选择v1.x.x的版本，不会去选择v2.x.x或更高的版本。

#### 最小版本选择

有时记录在`go.mod`文件中的依赖包版本会随着引入其他依赖包而发生变化。

如下图所示：

![](images/gomodule_minimal_version.png)

Module A 依赖 Module M的v1.0.0版本，但之后 Module A 引入了 Module D，而Module D 依赖 Module M的v1.1.1版本，此时，由于依赖的传递，Module A也会选择v1.1.1版本。

需要注意的是，此时会自动选择最小可用的版本，而不是最新的tag版本。
