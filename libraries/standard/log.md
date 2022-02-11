---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "log - 日志"  # 文章标题
url:  "posts/go/libraries/standard/log"  # 设置网页链接，默认使用文件名
tags: [ "go", "log", "logger" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

Log 标准库定义了 Logger 类型，该类型提供了一些格式化输出的方法，也提供了一个预定义的“标准” Logger，可以通过调用函数 Print 系列（Print|Printf|Println）、Fatal 系列（Fatal|Fatalf|Fatalln）、和 Panic 系列（Panic|Panicf|Panicln）来使用，比自行创建一个 Logger 对象更容易使用。

```go
func main() {
	log.Println("A Println() log.")
	// 中断程序
	log.Fatalln("A Fatalln() log.")
	// 中断程序
	log.Panicln("A Panicln() log.")
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, s)
	panic(s)
}
```

- Logger 会打印每条日志信息的日期时间，默认输出到系统的标准错误；
- `Fatal` 系列函数会在写入日志信息后调用 `os.Exit(1)`；
- `Panic` 系列函数会在写入日志信息后 `panic`。

## 配置 Logger

### 标准 Logger 的配置

Log 标准库中的 Flags 函数会返回标准 Logger 的输出配置，而 `SetFlags` 函数用来设置标准 Logger 的输出配置。

```go
func Flags() int  // 以数值编号
func SetFlags(flag int)
```

### flag 选项

Log 标准库提供了如下的 flag 选项，它们是一系列定义好的常量。输出时，每项之间以一个冒号分隔。

```go
const (
	Ldate         = 1 << iota     // 日期: 2009/01/23
	Ltime                         // 时间: 01:23:23
	Lmicroseconds                 //微秒: 01:23:23.123123
	Llongfile                     // 文件全路径+行号: /a/b/c/d.go:23
	Lshortfile                    // 文件名+行号: d.go:23 覆盖 Llongfile
	LUTC                          // UTC 时间
	LstdFlags     = Ldate | Ltime // 标准 logger 的初始值
)
```

```go
func main() {
    log.SetFlags(log.Ldate)  // 默认
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("A log.")
}
```

### 配置日志前缀

```go
func Prefix() string
func SetPrefix(prefix string)
```

```go
func main() {
    log.SetPrefix("[go-notes]")
	log.Println("A log.")
}
```

### 配置日志输出位置

```go
func SetOutput(w io.Writer)
```

`SetOutput` 函数用来设置标准 logger 的输出目的地，默认是标准错误输出。

把日志输出到同目录下的 `go.log` 文件中：

```go
func init() {
	logFile, _ := os.OpenFile("go.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("A log.")
}
```

### 创建新 Logger 对象

Log 标准库中还提供了一个创建新 Logger 对象的构造函数。

```go
func New(out io.Writer, prefix string, flag int) *Logger
```

```go
func main() {
	logger := log.New(os.Stdout, "<New>", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println("Customized Logger.")
}
```
