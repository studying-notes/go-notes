---
date: 2020-11-16T21:08:21+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go := 赋值符号"  # 文章标题
url:  "posts/go/docs/grammar/assign"  # 设置网页永久链接
tags: [ "Go", "assign" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 多变量赋值可能会重新声明

```go
func fun1() {
    i := 0
    i, j := 1, 2
    fmt.Printf("i = %d, j = %d\n", i, j)
}
```

程序输出如下：

```
i = 1, j = 2
```

前一个语句中已经声明了 i, 为什么还可以再次声明呢？

重新声明并没有什么问题，它并没有引入新的变量，只是把变量的值改变了，这是 Go 提供的一个语法糖。我们所说的重新声明不会引入问题要满足一个前提，变量声明要在同一个作用域中出现。如果出现在不同的作用域，那很可能就创建了新的同名变量。

```go
func fun2(i int) {
    i := 0
    fmt.Println(i)
}
```

在这里，形参已经声明了变量 i，使用`:=`再次声明是不允许的。

- 当 ` := ` 左侧存在新变量时（如 field2），那么已声明的变量（如 offset）则会被重新声明，不会有其他额外副作用。
- 当 ` := ` 左侧没有新变量是不允许的，编译会提示 `no new variable on left side of := `。

## 不能用于函数外部

简短变量场景只能用于函数中，使用 ` := ` 来声明和初始化全局变量是行不通的。

比如，像下面这样：

```go
package main

import "fmt"

rule := "Short variable declarations" // syntax error: non-declaration statement outside function body
```

这里的编译错误提示 `syntax error : non-declaration statement outside function body`，表示非声明语句不能出现在函数外部。可以理解成 ` := ` 实际上会拆分成两个语句，即声明和赋值。赋值语句不能出现在函数外部的。

## 变量作用域问题

```go
func fun3() {
    i, j := 0, 0

    if true {
        j, k := 1, 1
        fmt.Printf("j = %d, k = %d\n", j, k)
    }

    fmt.Printf("i = %d, j = %d\n", i, j)
}
```

程序输出如下：

```
j = 1, k = 1
i = 0, j = 0
```

这里要注意的是，block `if` 中声明的 j，与上面的 j 属于不同的作用域。

```go

```
