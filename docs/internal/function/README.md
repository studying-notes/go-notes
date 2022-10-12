---
date: 2022-10-05T11:05:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "函数与栈"  # 文章标题
url:  "posts/go/docs/internal/function/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 简介

函数是程序中为了执行特定任务而存在的一系列执行代码。函数接受输入并返回输出，执行程序的过程可以看作一系列函数的调用过程。

Go 语言中最重要的函数为 main 函数，其是程序执行用户代码的入口，在每个程序中都需要存在。

Go 语言中的函数有具名和匿名之分：**具名函数一般对应于包级的函数**，是匿名函数的一种特例。

```go
// 具名函数
func Add(a, b int) int {
    return a+b
}

// 匿名函数
var Add = func(a, b int) int {
    return a+b
}
```

方法是绑定到一个具体类型的特殊函数，Go 语言中的方法是依托于类型的，必须**在编译时静态绑定**。

接口定义了方法的集合，这些方法依托于运行时的接口对象，因此**接口对应的方法是在运行时动态绑定的**。

## 函数用法

函数具有减少冗余、隐藏信息、提高代码清晰度等优点。在 Go 语言中，函数是一等公民(first-class)，这意味着可以将它看作变量，并且它可以作为参数传递、返回及赋值。

```go
// 函数作为返回值
func calc1(base int) (func(int) int, func(int) int) {
	add := func(i int) int {
		base += i
		return base
	}
	sub := func(i int) int {
		base -= i
		return base
	}
	return add, sub
}

// 函数作为参数
func calc2(base int, f func(int, int) int) int {
	return f(base, base)
}
```

Go 语言中的函数还具有多返回值的特点，多返回值最常用于返回 error 错误信息，从而被调用者捕获。

Go 语言中的函数**可以有多个参数和多个返回值**，参数和返回值都是**以传值的方式**和被调用者交换数据。在语法上，函数还支持可变数量的参数，**可变数量的参数必须是最后出现的参数**，可变数量的参数其实是一个**切片类型的参数**。

```go
// 多个参数和多个返回值
func Swap(a, b int) (int, int) {
    return b, a
}

// 可变数量的参数
// more 对应 []int 切片类型
func Sum(a int, more ...int) int {
    for _, v := range more {
        a += v
    }
    return a
}
```

当可变参数是一个空接口类型时，调用者是否解包可变参数会导致不同的结果：

```go
func Print(a ...interface{}) {
    fmt.Println(a...)
}

func main() {
    var a = []interface{}{123, "abc"}

    Print(a...) // 123 abc
    Print(a)    // [123 abc]
}
```

第一个 `Print` 调用时传入的参数是 `a...`，等价于直接调用 `Print(123, "abc")`。第二个 `Print` 调用传入的是未解包的 `a`，等价于直接调用 `Print([]interface{}{123, "abc"} )`。

## 命名返回值

不仅函数的参数可以有名字，也可以给函数的返回值命名：

```go
func Find(m map[int]int, key int) (value int, ok bool) {
    value, ok = m[key]
    return
}
```

如果返回值命名了，可以通过名字来修改返回值，也可以通过 defer 语句在 return 语句之后修改返回值：

```go
func Inc() (v int) {
    defer func(){ v++ } ()
    return 42
}
```

## 闭包与陷阱

Go 语言同样支持匿名函数和闭包。闭包（Closure）是在函数作为一类公民的编程语言中实现词法绑定的一种技术，闭包包含了函数的入口地址和其关联的环境。闭包和普通函数最大的区别在于，**闭包函数中可以引用闭包外的变量**。

闭包的这种以**引用方式访问外部变量**的行为可能会导致一些隐含的问题：

```go
func main() {
    for i := 0; i < 3; i++ {
        defer func(){ println(i) } ()
    }
}
// Output:
// 3
// 3
// 3
```

因为是闭包，在 for 迭代语句中，每个 defer 语句延迟执行的函数引用的都是同一个 i 迭代变量，在循环结束后这个变量的值为 3，因此最终输出的都是 3。

修复的思路是在每轮迭代中为每个 defer 语句的闭包函数生成独有的变量。可以用下面两种方式：

```go
func main() {
    for i := 0; i < 3; i++ {
        j := i // 定义一个循环体内局部变量j
        defer func(){ println(j) } ()
    }
}

func main() {
    for i := 0; i < 3; i++ {
        // 通过函数传入i
        // defer 语句会马上对调用参数求值
        defer func(i int){ println(i) } (i)
    }
}
```

第一种方法是在循环体内部再定义一个局部变量，这样每次迭代 defer 语句的闭包函数捕获的都是不同的变量，这些变量的值对应迭代时的值。

第二种方式是将迭代变量通过闭包函数的参数传入，defer 语句会马上对调用参数求值。两种方式都是可以工作的。

## 参数传值与传引用

Go 语言中，如果以切片为参数调用函数，有时候会给人一种参数采用了传引用的方式的假象：因为在被调用函数内部可以修改传入的切片的元素。其实，**任何可以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数**。

函数参数传值的规范更准确说是**只针对数据结构中固定的部分传值**，例如切片作为参数时，其结构体中的指针、容量和长度仍然是传值，但是并**不包含指针间接指向的内容**。将切片类型的参数替换为类似 `reflect.SliceHeader` 结构体就能很好理解切片传值的含义了：

```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

```go
package main

import "fmt"

func argument(s []int) {
	fmt.Printf("addr: %p, len: %d, cap: %d, val: %v\n", s, len(s), cap(s), s)
	s = append(s, 1, 2, 3)
	fmt.Printf("addr: %p, len: %d, cap: %d, val: %v\n", s, len(s), cap(s), s)
}

func main() {
	s := make([]int, 1, 4)
	fmt.Printf("addr: %p, len: %d, cap: %d, val: %v\n", s, len(s), cap(s), s)
	argument(s)
	fmt.Printf("addr: %p, len: %d, cap: %d, val: %v\n", s, len(s), cap(s), s)
}
```

```
addr: 0xc0000a4060, len: 1, cap: 4, val: [0]
addr: 0xc0000a4060, len: 1, cap: 4, val: [0]
addr: 0xc0000a4060, len: 4, cap: 4, val: [0 1 2 3]
addr: 0xc0000a4060, len: 1, cap: 4, val: [0]
```

因为切片中的底层数组部分通过隐式指针传递（指针本身依然是传值的，但是指针指向的却是同一份的数据），所以被调用函数可以通过指针修改调用参数切片中的数据。

除数据之外，切片结构还包含了切片长度和切片容量信息，这两个信息也是传值的。如果被调用函数中修改了 Len 或 Cap 信息，如上面的案例所示，就无法反映到调用参数的切片中，这时候我们一般会通过返回修改后的切片来更新之前的切片。这也是内置的 `append()` 必须要返回一个切片的原因。

## 函数栈

[函数栈](stack.md)

```go

```
