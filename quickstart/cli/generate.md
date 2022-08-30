---
date: 2022-07-19T10:10:18+08:00
author: "Rustle Karl"

title: "go generate 命令详解"
url:  "posts/go/quickstart/cli/generate"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

go generate 允许在 go 代码中来运行当前系统中已安装的程序，原则上可以运行任何程序，但是此命令设计的初衷是用来创建或者更新 go 源码文件。

必须手动执行 go generate，才会去解析执行 generate 指令，并且 go generate 命令不会执行 go 源代码。

单个运行只需要在对应目录（包）下执行：`go generate`；对应有多个目录，且又有多个 generate 需要执行的，在最上层目录下运行：`go generate ./...`。

一个简单的例子：

```go
package main

import "fmt"

//go:generate echo main

func main() {
	fmt.Println("Hello World")
}
```

在 linux 系统上运行 `go generate`，输出 main，就相当于你直接运行 `echo hello`。

运行 `go generate` 时，它将扫描与当前包相关的源代码文件，找出所有包含 `//go:generate` 的特殊注释，提取并执行该特殊注释后面的命令，命令为可执行程序，形同 shell 下面执行。

## 规则

- 必须在 .go 源码文件中。
- 每个源码文件可以包含多个 generate 特殊注释。
- 显示运行 go generate 命令时，才会执行特殊注释后面的命令。
- 如果前面的注释执行出错，则终止执行。
- // 与 go:generate 之间不能有空格。

## 可以使用的环境变量

```
$GOARCH：体系架构
$GOOS：操作系统
$GOFILE：当前处理中的文件名
$GOLINE：当前命令在文件中的行号
$GOPACKAGE：当前处理文件的包名
$DOLLAR：美元符号
```
