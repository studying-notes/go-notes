---
date: 2020-11-15T19:58:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 示例测试"  # 文章标题
url:  "posts/go/libraries/standard/testing/example"  # 设置网页永久链接
tags: [ "Go", "example" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

创建两个文件，其中 `example.go` 为源代码文件，`example_test.go` 为测试文件。

## 源代码文件

源代码文件 `example.go` 中包含 `SayHello()`、`SayGoodbye()` 和 `PrintNames()` 三个方法，如下所示：

```go
package gotest

import "fmt"

// SayHello 打印一行字符串
func SayHello() {
    fmt.Println("Hello World")
}

// SayGoodbye 打印两行字符串
func SayGoodbye() {
    fmt.Println("Hello,")
    fmt.Println("goodbye")
}

// PrintNames 打印学生姓名
func PrintNames() {
    students := make(map[int]string, 4)
    students[1] = "Jim"
    students[2] = "Bob"
    students[3] = "Tom"
    students[4] = "Sue"
    for _, value := range students {
        fmt.Println(value)
    }
}
```

这几个方法打印内容略有不同，分别代表一种典型的场景：

* SayHello()：只有一行打印输出
* SayGoodbye()：有两行打印输出
* PrintNames()：有多行打印输出，且由于 Map 数据结构的原因，多行打印次序是随机的。

## 测试文件

测试文件 `example_test.go` 中包含 3 个测试方法，于源代码文件中的 3 个方法一一对应，测试文件如下所示：

```go
package gotest

// 检测单行输出
func ExampleSayHello() {
    SayHello()
    // OutPut: Hello World
}

// 检测多行输出
func ExampleSayGoodbye() {
    SayGoodbye()
    // OutPut:
    // Hello,
    // goodbye
}

// 检测乱序输出
func ExamplePrintNames() {
    PrintNames()
    // Unordered output:
    // Jim
    // Bob
    // Tom
    // Sue
}
```

例子测试函数命名规则为 "ExampleXxx"，其中 "Xxx" 为自定义的标识，通常为待测函数名称。

这三个测试函数分别代表三种场景：

* ExampleSayHello()：待测试函数只有一行输出，使用 "// OutPut : " 检测。
* ExampleSayGoodbye()：待测试函数有多行输出，使用 "// OutPut : " 检测，其中期望值也是多行。
* ExamplePrintNames()：待测试函数有多行输出，但输出次序不确定，使用 "// Unordered output : " 检测。

注：字符串比较时会忽略前后的空白字符。

## 执行测试

命令行下，使用 `go test` 或 `go test example_test.go` 命令即可启动测试，如下所示：

```
go test example_test.go
```

```
ok      command-line-arguments  0.331s
```

## 小结

1. 例子测试函数名需要以 "Example" 开头；
2. 最终测试函数必须输出字符串；
3. 检测单行输出格式为 “// Output : < 期望字符串 >” ；
4. 检测多行输出格式为 “// Output : \ < 期望字符串 > \ < 期望字符串 >”，每个期望字符串占一行；
5. 检测无序输出格式为 "// Unordered output : \ < 期望字符串 > \ < 期望字符串 >"，每个期望字符串占一行；
6. 测试字符串时会自动忽略字符串前后的空白字符；
7. 如果测试函数中没有 “Output” 标识，则该测试函数不会被执行；
8. 执行测试可以使用 `go test`，此时该目录下的其他测试文件也会一并执行；
9. 执行测试可以使用 `go test <xxx_test.go>`，此时仅执行特定文件中的测试函数。

```go

```
