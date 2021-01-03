---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "理解 GOPATH"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/doc/mod/gopath"  # 设置网页永久链接
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

## GOROOT 是什么

通常我们说安装 Go 语言，实际上安装的是 Go 编译器和 Go 标准库，二者位于同一个安装包中。

假如你在Windows上使用Installer安装的话，它们将会被默认安装到`c:\go`目录下，该目录即 GOROOT 目录，里面保存了开发 GO 程序所需要的所有组件，比如编译器、标准库和文档等等。

同时安装程序还会自动帮你设置 GOROOT 环境变量，如下图所示：

[![DF0RfJ.png](https://s3.ax1x.com/2020/11/15/DF0RfJ.png)](https://imgchr.com/i/DF0RfJ)

另外，安装程序还会把 `c:\go\bin` 目录添加到系统的 `PATH` 环境变量中，如下图所示：

[![DF0fp9.png](https://s3.ax1x.com/2020/11/15/DF0fp9.png)](https://imgchr.com/i/DF0fp9)

该目录主要是 GO 语言开发包中提供的二进程可执行程序。

所以，GOROOT 实际上是指示 GO 语言安装目录的环境变量，属于 GO 语言顶级目录。

## GOPATH 是什么

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

## 依赖查找

当某个 package 需要引用其他包时，编译器就会依次从 `GOROOT/src/` 和 `GOPATH/src/` 中去查找，如果某个包从 GOROOT 下找到的话，就不再到 GOPATH 目录下查找，所以如果项目中开发的包名与标准库相同的话，将会被自动忽略。

## GOPATH的缺点

GOPATH 的优点是足够简单，但它不能很好的满足实际项目的工程需求。

比如，你有两个项目 A 和 B，他们都引用某个第三方库 T，但这两个项目使用了不同的 T 版本，即：

- 项目A 使用T v1.0
- 项目B 使用T v2.0

由于编译器依赖查找固定从 GOPATH/src 下查找 `GOPATH/src/T`，所以，无法在同一个 GOPATH 目录下保存第三方库 T 的两个版本。所以项目 A、B 无法共享同一个 GOPATH，需要各自维护一个，这给广大软件工程师带来了极大的困扰。

针对 GOPATH 的缺点，GO 语言社区提供了 Vendor 机制，从此依赖管理进入第二个阶段：将项目的依赖包私有化。
