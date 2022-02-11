---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 4"  # 文章标题
url:  "posts/go/docs/mod/4_module_exclude"  # 设置网页永久链接
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

`go.mod` 文件中的 `exclude` 指令用于排除某个包的特定版本，其与 `replace` 类似，也仅在当前 module 为 `main module` 时有效，其他项目引用当前项目时，`exclude` 指令会被忽略。

`exclude` 指令在实际的项目中很少被使用，因为很少会显式地排除某个包的某个版本，除非我们知道某个版本有严重 bug。
比如指令 `exclude github.com/google/uuid v1.1.0`，表示不使用 v1.1.0 版本。

下面我们还是使用`github.com/renhongcai/gomodule`来举例说明。

## 排除指定版本

在 `github.com/renhongcai/gomodule` 的 v1.3.0 版本中，我们的 `go.mod` 文件如下：

```
module github.com/renhongcai/gomodule

go 1.13

require (
  github.com/google/uuid v1.0.0
  golang.org/x/text v0.3.2
)

replace golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
```

`github.com/google/uuid v1.0.0` 说明我们期望使用 uuid 包的 `v1.0.0` 版本。

假如，当前 uuid 仅有 `v1.0.0`、`v1.1.0` 和 `v1.1.1` 三个版本可用，而且我们假定 `v1.1.0` 版本有严重 bug。
此时可以使用 `exclude` 指令将 uuid 的 `v1.1.0` 版本排除在外，即在 `go.mod` 文件添加如下内容：

```
exclude github.com/google/uuid v1.1.0
```

虽然我们暂时没有使用 uuid 的 `v1.1.0` 版本，但如果将来引用了其他包，正好其他包引用了 uuid 的 `v1.1.0` 版本的话，此时添加的 `exclude` 指令就会跳过 `v1.1.0` 版本。

下面我们创建 `github.com/renhongcai/exclude` 包来验证该问题。

#### 创建依赖包

为了进一步说明 `exclude` 用法，我们创建了一个仓库`github.com/renhongcai/exclude`，并在其中创建了一个 module `github.com/renhongcai/exclude`，其中 `go.mod` 文件（v1.0.0 版本）如下：

```
module github.com/renhongcai/exclude

go 1.13

require github.com/google/uuid v1.1.0

```

可以看出其依赖 `github.com/google/uuid` 的 `v1.1.0` 版本。创建 `github.com/renhongcai/exclude` 的目的是供 `github.com/renhongcai/gomodule` 使用的。

#### 使用依赖包

由于 `github.com/renhongcai/exclude` 也引用了 uuid 包且引用了更新版本的 uuid，那么在 `github.com/renhongcai/gomodule`引用 `github.com/renhongcai/exclude` 时，会被动的提升 uuid 的版本。

在没有添加 `exclude` 之前，编译时 `github.com/renhongcai/gomodule` 依赖的 uuid 版本会提升到 `v1.1.0`，与 `github.com/renhongcai/exclude` 保持一致，相应的 `go.mod` 也会被自动修改，如下所示：

```
module github.com/renhongcai/gomodule

go 1.13

require (
	github.com/google/uuid v1.1.0
	github.com/renhongcai/exclude v1.0.0
	golang.org/x/text v0.3.2
)

replace golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
```

但如果添加了 `exclude github.com/google/uuid v1.1.0` 指令后，编译时 `github.com/renhongcai/gomodule` 依赖的 uuid 版本会自动跳过 `v1.1.0`，即选择 `v1.1.1` 版本，相应的 `go.mod` 文件如下所示：

```
module github.com/renhongcai/gomodule

go 1.13

require (
	github.com/google/uuid v1.1.1
	github.com/renhongcai/exclude v1.0.0
	golang.org/x/text v0.3.2
)

replace golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2

exclude github.com/google/uuid v1.1.0
```

在本例中，在选择版本时，跳过 uuid `v1.1.0` 版本后还有 `v1.1.1` 版本可用，Go 命令行工具可以自动选择 `v1.1.1` 版本，但如果没有更新的版本时将会报错而无法编译。
