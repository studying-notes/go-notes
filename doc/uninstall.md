---
date: 2020-11-15T20:16:04+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 卸载指南"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/doc/uninstall"  # 设置网页永久链接
tags: [ "go", "config" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

当需要升级新的 Go 语言版本时，你需要先把旧版本删除。Go 语言版本升级过程实际上是 `删除旧版本`+`安装新版本`。

删除 Go 语言版本是安装新版本的逆过程，即把新版本安装时创建的目录、环境变量删除。

## 手动卸载旧版本

```shell
rm -rf /usr/bin/go
rm -rf /usr/local/go
apt purge golang-go
```

## 删除 Go 安装目录

通过 `go env` 命令查询安装目录，安装目录即 `GOROOT` 环境变量所指示的目录，如下所示：

```shell
# go env
GOPATH = "/root/go"
GOROOT = "/usr/local/go"
```

`go env` 命令会输出很多 Go 语言相关的环境变量，上面只保留了最关键的两个 `GOROOT` 和 `GOPATH`。

接下来使用 `rm` 命令删除 `GOROOT` 指向的目录即可，比如 ` # rm -rf /usr/local/go`。

## 删除残留的可执行文件

Go 程序在运行过程中会在 `GOPATH/bin` 目录下生成可执行文件，为了安全起见，也需要删除。

同样，使用 `rm` 命令删除即可，比如 ` # rm -rf /root/go/bin`。

**注：如果 GOPATH 包含多个目录，需要删除每个目录下的 bin 目录**。

## 删除环境变量

将环境变量 `GOPATH` 删除，该环境变量一般是前一次安装 Go 时人为设置的。

环境变量一般存在于以下几个文件中：
* /etc/profile
* /etc/bashrc
* ~ /.bash_profile
* ~ /.profile
* ~ /.bashrc

做完以上步骤，旧版本就被彻底的删除了。
