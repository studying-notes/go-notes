---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 模块深入讲解 1"  # 文章标题
url:  "posts/go/docs/module/1_module_basic"  # 设置网页永久链接
tags: [ "go" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

在开始学习 module 机制之前，我们有必要初步了解一下其涉及的基本概念，比如到底什么是 module 等。

## Module 的定义

首先，module 是个新鲜又熟悉的概念。新鲜是指在以往 GOPATH 和 vendor 时代都没有提及，它是个新的词汇。

为什么说熟悉呢？因为它不是新的事物，事实上我们经常接触，这次只是官方给了一个统一的称呼。

拿开源项目 `https://github.com/blang/semver` 举例，这个项目是一个语义化版本处理库，当你的项目需要时可以在你的项目中 import，比如：

```golang
import "github.com/blang/semver"
```

`https://github.com/blang/semver` 项目中可以包含一个或多个 package，不管有多少 package，这些 package 都随项目一起发布，即当我们说`github.com/blang/semver`某个版本时，说的是整个项目，而不是具体的 package。此时项目`https://github.com/blang/semver`就是一个 module。

官方给 module 的定义是：`A module is a collection of related Go packages that are versioned together as a single unit.`。一组 package 的集合，一起被标记版本，即是一个 module。

通常而言，一个仓库包含一个 module（虽然也可以包含多个，但不推荐），所以仓库、module 和 package 的关系如下：

- 一个仓库包含一个或多个 Go module ；
- 每个 Go module 包含一个或多个 Go package ；
- 每个 package 包含一个或多个 Go 源文件；

此外，一个 module 的版本号规则必须遵循语义化规范（https://semver.org/），版本号必须使用格式 `v(major).(minor).(patch)`，比如 `v0.1.0`、`v1.2.3` 或`v1.5.0-rc.1`。

## 语义化版本规范

语义化版本（Semantic Versioning）已成为事实上的标准，几乎知名的开源项目都遵循该规范，更详细的信息请前往 https://semver.org/ 查看，在此只提炼一些要点，以便于后续的阅读。

版本格式 `v(major).(minor).(patch)` 中 major 指的是大版本，minor 指小版本，patch 指补丁版本。

- major : 当发生不兼容的改动时才可以增加 major 版本；比如 `v2.x.y` 与 `v1.x.y` 是不兼容的；
- minor : 当有新增特性时才可以增加该版本，比如 `v1.17.0` 是在 `v1.16.0` 基础上加了新的特性，同时兼容 `v1.16.0` ；
- patch : 当有 bug 修复时才可以 增加该版本，比如 `v1.17.1` 修复了 `v1.17.0` 上的 bug，没有新特性增加；

语义化版本规范的好处是，用户通过版本号就能了解版本信息。

除了上面介绍的基础概念以外，还有描述依赖的 `go.mod` 和记录 module 的 checksum 的 `go.sum` 等内容。
