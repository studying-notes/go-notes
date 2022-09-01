---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 模块深入讲解 8"  # 文章标题
url:  "posts/go/docs/module/8_module_pseudo_version"  # 设置网页永久链接
tags: [ "go" ]
categories: [ "Go 学习笔记" ]

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

在 `go.mod` 中通常使用语义化版本来标记依赖，比如 `v1.2.3`、`v0.1.5` 等。因为 `go.mod` 文件通常是 `go` 命令自动生成并修改的，所以实际上是 `go` 命令习惯使用语义化版本。

诸如 `v1.2.3` 和 `v0.1.5` 这样的语义化版本，实际是某个 commit ID 的标记，真正的版本还是 commit ID。比如 `github.com/renhongcai/gomodule` 项目的 `v1.5.0` 对应的真实版本为 `20e9757b072283e5f57be41405fe7aaf867db220`。

由于语义化版本比 `commit ID` 更直观（方便交流与比较版本大小），所以一般情况下使用语义化版本。

## 什么是伪版本

在实际项目中，有时不得不直接使用一个 `commit ID`，比如某项目发布了 `v1.5.0` 版本，但随即又修复了一个 bug（引入一个新的 commit ID），而且没有发布新的版本。此时，如果我们希望使用最新的版本，就需要直接引用最新的 `commit ID`，而不是之前的语义化版本 `v1.5.0`。

使用 `commit ID` 的版本在 Go 语言中称为 `pseudo-version`，可译为 " 伪版本 "。

伪版本的版本号通常会使用 `vx.y.z-yyyymmddhhmmss-abcdefabcdef` 格式，其中 `vx.y.z` 看上去像是一个真实的语义化版本，但通常并不存在该版本，所以称为伪版本。另外 `abcdefabcdef` 表示某个 commit ID 的前 12 位，而 `yyyymmddhhmmss` 则表示该 commit 的提交时间，方便做版本比较。

使用伪版本的 `go.mod` 举例如下：

```
require (
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
)
```

## 伪版本风格

伪版本格式都为 `vx.y.z-yyyymmddhhmmss-abcdefabcdef`，但 `vx.y.z` 部分在不同情况下略有区别，有时可能是 `vx.y.z-pre.0` 或者 `vx.y.z -0`，甚至 `vx.y.z-dev.2.0` 等。

`vx.y.z` 的具体格式取决于所引用 `commit ID` 之前的版本号，如果所引用 `commit ID` 之前的最新的 tag 版本为 `v1.5.0`，那么伪版本号则在其基础上增加一个标记，即 `v1.5.1 -0`，看上去像是下一个版本一样。

实际使用中 `go` 命令会帮我们自动生成伪版本，不需要手动计算，所以此处我们仅做基本说明。

## 如何获取伪版本

我们使用具体的例子还演示如何使用伪版本。在仓库 `github.com/renhongcai/gomodule` 中存在 `v1.5.0` tag 版本，在 `v1.5.0` 之后又提交了一个 commit，并没有发布新的版本。

为了方便描述，我们把 `v1.5.0` 对应的 commit 称为 `commit-A`，而其随后的 commit 称为 `commit-B`。

如果我们要使用 commit-A，即 `v1.5.0`，可使用 `go get github.com/renhongcai/gomodule@v1.5.0` 命令：

```bash
go get github.com/renhongcai/gomodule@v1.5.0
```

```
go: finding github.com/renhongcai/gomodule v1.5.0
go: downloading github.com/renhongcai/gomodule v1.5.0
go: extracting github.com/renhongcai/gomodule v1.5.0
go: finding github.com/renhongcai/indirect v1.0.1
```

此时，如果存在 `go.mod` 文件，`github.com/renhongcai/gomodule` 体现在 `go.mod` 文件的版本为 `v1.5.0`。

如果我们要使用 commit-B，可使用 `go get github.com/renhongcai/gomodule@6eb27062747a458a27fb05fceff6e3175e5eca95` 命令（可以使用完整的 commit id，也可以只使用前 12 位）：

```
go get github.com/renhongcai/gomodule@6eb27062747a458a27fb05fceff6e3175e5eca95
```

```
go: finding github.com 6eb27062747a458a27fb05fceff6e3175e5eca95
go: finding github.com/renhongcai/gomodule 6eb27062747a458a27fb05fceff6e3175e5eca95
go: finding github.com/renhongcai 6eb27062747a458a27fb05fceff6e3175e5eca95
go: downloading github.com/renhongcai/gomodule v1.5.1-0.20200203082525-6eb27062747a
go: extracting github.com/renhongcai/gomodule v1.5.1-0.20200203082525-6eb27062747a
go: finding github.com/renhongcai/indirect v1.0.2
```

此时，可以看到生成的伪版本号为 `v1.5.1 -0.20200203082525 -6eb27062747a`，当前最新版本为 `v1.5.0`，`go` 命令生成伪版本号时自动增加了版本。此时，如果存在 `go.mod` 文件的话，`github.com/renhongcai/gomodule` 体现在 `go.mod` 文件中的版本则为该伪版本号。
