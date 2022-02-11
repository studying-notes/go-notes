---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "fmt - 获取用户输入"  # 文章标题
url:  "posts/go/libraries/standard/fmt"  # 设置网页链接，默认使用文件名
tags: [ "go", "fmt", "io" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

Go 语言 fmt 包下有 `fmt.Scan`、`fmt.Scanf`、`fmt.Scanln` 三个函数，可以在程序运行过程中从标准输入获取用户的输入。

## fmt.Scan

```go
func Scan(a ...interface{}) (n int, err error)
```

- Scan 从标准输入扫描文本，读取由空白符分隔的值保存到传递给本函数的参数中，换行符视为空白符。
- 返回成功扫描的数据个数和遇到的任何错误。如果读取的数据个数比提供的参数少，会返回一个错误报告原因。

```go
func main() {
	var (
		name    string
		age     int
		married bool
	)
	_, _ = fmt.Scan(&name, &age, &married)
	fmt.Printf("Result name: %s age: %d married: %t \n", name, age, married)
}
```

输入：

```
karl 18 false
```

输出：

```
Result name: karl age: 18 married: false
```

## fmt.Scanf

```go
func Scanf(format string, a ...interface{}) (n int, err error)
```

- Scanf 从标准输入扫描文本，根据 format 参数指定的格式去读取由空白符分隔的值保存到传递给本函数的参数中。
- 返回成功扫描的数据个数和遇到的任何错误。如果读取的数据个数比提供的参数少，会返回一个错误报告原因。

```go
func main() {
	var (
		name    string
		age     int
		married bool
	)
	_, _ = fmt.Scanf("1:%s 2:%d 3:%t", &name, &age, &married)
	fmt.Printf("Result name: %s age: %d married: %t \n", name, age, married)
}
```

> `fmt.Scanf` 不同于 `fmt.Scan` 简单的以空格作为输入数据的分隔符，`fmt.Scanf` 为输入数据指定了具体的输入内容格式，只有按照格式输入数据才会被扫描并存入对应变量。

必须按指定格式输入：

```
1:karl 2:18 3:false
```

输出：

```
Result name: karl age: 18 married: false
```


## fmt.Scanln

```go
func Scanln(a ...interface{}) (n int, err error)
```

- Scanln 类似 Scan，它在遇到换行时停止扫描，即使参数不够，而 Scan 输入的参数不够的话换行不会停止扫描，而是继续等待用户输入。Scanln 最后一个数据后面必须有换行或者到达结束位置。
- 返回成功扫描的数据个数和遇到的任何错误。如果读取的数据个数比提供的参数少，会返回一个错误报告原因。
