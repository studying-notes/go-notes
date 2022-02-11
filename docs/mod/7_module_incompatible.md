---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 7"  # 文章标题
url:  "posts/go/docs/mod/7_module_incompatible"  # 设置网页永久链接
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

在前面的章节中，我们介绍了Go module的版本选择机制，其中介绍了一个Module的版本号需要遵循`v<major>.<minor>.<patch>`的格式，此外，如果`major`版本号大于1时，其版本号还需要体现在Module名字中。

比如Module `github.com/RainbowMango/m`，如果其版本号增长到`v2.x.x`时，其Module名字也需要相应的改变为：
`github.com/RainbowMango/m/v2`。即，如果`major`版本号大于1时，需要在Module名字中体现版本。

那么如果Module的`major`版本号虽然变成了`v2.x.x`，但Module名字仍保持原样会怎么样呢？ 其他项目是否还可以引用呢？其他项目引用时有没有风险呢？这就是今天要讨论的内容。

## 能否引用不兼容的包

我们还是以Module `github.com/RainbowMango/m` 为例，假如其当前版本为`v3.6.0`，因为其Module名字未遵循Golang所推荐的风格，即Module名中附带版本信息，我们称这个Module为不规范的Module。

不规范的Module还是可以引用的，但跟引用规范的Module略有差别。

如果我们在项目A中引用了该module，使用命令`go mod tidy`，go 命令会自动查找Module m的最新版本，即`v3.6.0`。
由于Module为不规范的Module，为了加以区分，go 命令会在`go.mod`中增加`+incompatible` 标识。

```
require (
	github.com/RainbowMango/m v3.6.0+incompatible
)
```

除了增加`+incompatible`（不兼容）标识外，在其使用上没有区别。

## 如何处理incompatible

`go.mod`文件中出现`+incompatible`，说明你引用了一个不规范的Module，正常情况下，只能说明这个Module版本未遵循版本化语义规范。但引用这个规范的Module还是有些困扰，可能还会有一定的风险。

比如，我们拿某开源Module `github.com/blang/semver`为例，编写本文时，该Module最新版本为`v3.6.0`，但其`go.mod`中记录的Module却是：
```
module github.com/blang/semver
```

Module `github.com/blang/semver` 在另一个著名的开源软件`Kubernetes`（github.com/kubernetes/kubernetes）中被引用，那么`Kubernetes`的`go.mod`文件则会标记这个Module为`+incompatible`：

```
require (
	...
	github.com/blang/semver v3.5.0+incompatible
	...
）
```

站在`Kubernetes`的角度，此处的困扰在于，如果将来 `github.com/blang/semver`发布了新版本`v4.0.0`，但不幸的是Module名字仍然为`github.com/blang/semver`。那么，升级这个Module的版本将会变得困难。因为`v3.6.0`到`v4.0.0`跨越了大版本，按照语义化版本规范来解释说明发生了不兼容的改变，即然不兼容，项目维护者有必须对升级持谨慎态度，甚至放弃升级。

站在`github.com/blang/semver`的角度，如果迟迟不能将自身变得"规范"，那么其他项目有可能放弃本Module，转而使用其他更规范的Module来替代，开源项目如果没有使用者，也就走到了尽头。
