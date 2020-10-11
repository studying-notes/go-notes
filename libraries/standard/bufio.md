---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "通过 bufio 标准库获取用户输入"  # 文章标题
url:  "posts/go/libraries/standard/bufio"  # 设置网页链接，默认使用文件名
tags: [ "go", "bufio", "io" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

有时候我们想完整获取输入的内容，而输入的内容可能包含空格，这种情况下可以使用 bufio 包来实现。

```go
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter: ")
	s, _ := reader.ReadString('\n')
	fmt.Printf("%#v\n", strings.TrimSpace(s))
}
```

## Fscan 系列

这几个函数功能分别类似于 fmt.Scan、fmt.Scanf、fmt.Scanln 三个函数，只不过它们不是从标准输入中读取数据而是从 io.Reader 中读取数据。

```go
func Fscan(r io.Reader, a ...interface{}) (n int, err error)
func Fscanln(r io.Reader, a ...interface{}) (n int, err error)
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
```

## Sscan系列

这几个函数功能分别类似于fmt.Sscan、fmt.Sscanf、fmt.Sscanln三个函数，只不过它们不是从标准输入中读取数据而是从指定字符串中读取数据。

```go
func Sscan(str string, a ...interface{}) (n int, err error)
func Sscanln(str string, a ...interface{}) (n int, err error)
func Sscanf(str string, format string, a ...interface{}) (n int, err error)
```
