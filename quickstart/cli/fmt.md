---
date: 2022-07-19T10:44:53+08:00
author: "Rustle Karl"

title: "go fmt 命令详解"
url:  "posts/go/quickstart/cli/fmt"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

```shell
go help fmt
```

格式化 .go 源文件。

## 用法

```shell
go fmt [-n] [-x] [packages]
```

实际上调用 gofmt 程序，相当于：

```shell
gofmt -l -w
```

- `-n` 打印将执行的命令
- `-x` 在执行时打印命令

如果需要指定其他参数，可以直接执行 gofmt 命令。

```shell

```

```shell

```
