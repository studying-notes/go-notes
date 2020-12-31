---
date: 2020-11-15T20:40:59+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "理解 Go Module 机制"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/doc/mod/gomod"  # 设置网页永久链接
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

事实上，Go Module 是一个非常复杂的特性，一下子全盘托出其特性，往往会让人产生疑惑，所以接下来的章节，我们希望逐个介绍其特性，
并且，尽可以附以实例，希望大家也跟我一样手动实践一下，以加深认识。
