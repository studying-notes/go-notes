---
date: 2020-11-15T16:41:20+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 内存逃逸分析"  # 文章标题
url:  "posts/go/docs/internal/memory/escape"  # 设置网页永久链接
tags: [ "Go", "escape" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [前言](#前言)
- [逃逸策略](#逃逸策略)
- [逃逸场景](#逃逸场景)
	- [指针逃逸](#指针逃逸)
	- [栈空间不足逃逸](#栈空间不足逃逸)
	- [动态类型逃逸](#动态类型逃逸)
	- [闭包引用对象逃逸](#闭包引用对象逃逸)
- [逃逸总结](#逃逸总结)
- [函数传递指针真的比传值效率高吗？](#函数传递指针真的比传值效率高吗)

## 前言

逃逸分析（Escape analysis）是指由编译器决定内存分配的位置，不需要程序员指定。

函数中申请一个新的对象

- 如果分配在栈中，则**函数执行结束可自动将内存回收**，生命周期仅在函数执行期间；
- 如果分配在堆中，则**函数执行结束可交给 GC（垃圾回收）处理**；

有了逃逸分析，**返回函数局部变量**将变得可能，除此之外，逃逸分析还跟闭包息息相关，了解哪些场景下对象会逃逸至关重要。

## 逃逸策略

每当函数中申请新的对象，编译器会**根据该对象是否被函数外部引用来决定是否逃逸**：

1. 如果函数外部**没有引用，则优先放到栈中**；
2. 如果函数外部**存在引用，则必定放到堆中**；

对于函数外部**没有引用的对象，也有可能放到堆**中，比如**内存过大超过栈的存储能力**。

## 逃逸场景

### 指针逃逸

我们知道 Go 可以返回局部变量指针，这其实是一个典型的变量逃逸案例，示例代码如下：

```go
package main

type Student struct {
	Name string
	Age  int
}

func StudentRegister(name string, age int) *Student {
	s := new(Student)
	s.Name = name
	s.Age = age
	return s
}

func main() {
	StudentRegister("Jim", 18)
}
```

函数 StudentRegister() 内部 s 为局部变量，其值通过函数返回值返回，s 本身为一指针，其指向的内存地址不会是栈而是堆，这就是典型的逃逸案例。

通过编译参数 -gcflag=-m 可以查看编译过程中的逃逸分析：

```bash
go build -gcflags=-m
```

```
.\main.go:8:6: can inline StudentRegister
.\main.go:15:6: can inline main
.\main.go:16:17: inlining call to StudentRegister
.\main.go:8:22: leaking param: name
.\main.go:9:10: new(Student) escapes to heap
.\main.go:16:17: new(Student) does not escape
```

可见在 StudentRegister() 函数中，也即代码第 9 行显示 "escapes to heap"，代表该行内存分配发生了逃逸现象。

如果返回的不是指针则不会发生逃逸。

### 栈空间不足逃逸

看下面的代码，是否会产生逃逸呢？

```go
package main

func Slice() {
	s := make([]int, 1000, 1000)

	for index, _ := range s {
		s[index] = index
	}
}

func main() {
	Slice()
}
```

上面代码 Slice() 函数中分配了一个 1000 个长度的切片，是否逃逸取决于栈空间是否足够大。

```bash
go build -gcflags=-m
```

```
.\main.go:3:6: can inline Slice
.\main.go:11:6: can inline main
.\main.go:12:7: inlining call to Slice
.\main.go:4:11: make([]int, 1000, 1000) does not escape
.\main.go:12:7: make([]int, 1000, 1000) does not escape
```

我们发现此处并没有发生逃逸。那么把切片长度扩大 10 倍即 10000 会如何呢?

```
.\main.go:3:6: can inline Slice
.\main.go:11:6: can inline main
.\main.go:12:7: inlining call to Slice
.\main.go:4:11: make([]int, 10000, 10000) escapes to heap
.\main.go:12:7: make([]int, 10000, 10000) escapes to heap
```

我们发现当切片长度扩大到 10000 时就会逃逸。

实际上当栈空间不足以存放当前对象时或无法判断当前切片长度时会将对象分配到堆中。

### 动态类型逃逸

很多函数参数为 interface 类型，比如

```go
fmt.Println(a...interface{})
```

编译期间很难确定其参数的具体类型，也会产生逃逸。

如下代码所示：

```go
package main

import "fmt"

func main() {
    s := "Escape"
    fmt.Println(s)
}
```

上述代码 s 变量只是一个 string 类型变量，调用 fmt.Println() 时会产生逃逸：

```shell
go build -gcflags=-m
```

```
.\main.go:7: s escapes to heap
.\main.go:7: main ... argument does not escape
```

### 闭包引用对象逃逸

```go
func Fibonacci() func() int {
    a, b := 0, 1
    return func() int {
        a, b = b, a+b
        return a
    }
}
```

该函数返回一个闭包，闭包引用了函数的局部变量 a 和 b，使用时通过该函数获取该闭包，然后每次执行闭包都会依次输出 Fibonacci 数列。

完整的示例程序如下所示：

```go
package main

import "fmt"

func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	f := Fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Printf("Fibonacci: %d\n", f())
	}
}
```

上述代码通过 Fibonacci() 获取一个闭包，每次执行闭包就会打印一个 Fibonacci 数值。

Fibonacci() 函数中原本属于局部变量的 a 和 b 由于闭包的引用，不得不将二者放到堆上，以致产生逃逸：

```
.\main.go:5:6: can inline Fibonacci
.\main.go:7:9: can inline Fibonacci.func1
.\main.go:14:16: inlining call to Fibonacci
.\main.go:7:9: can inline main.func1
.\main.go:17:34: inlining call to main.func1
.\main.go:17:13: inlining call to fmt.Printf
.\main.go:6:2: moved to heap: a
.\main.go:6:5: moved to heap: b
.\main.go:7:9: func literal escapes to heap
.\main.go:14:16: func literal does not escape
.\main.go:17:13: ... argument does not escape
.\main.go:17:34: ~R0 escapes to heap
```

## 逃逸总结

- 栈上分配内存比在堆中分配内存有更高的效率

- 栈上分配的内存不需要 GC 处理

- 堆上分配的内存使用完毕会交给 GC 处理

- 逃逸分析目的是决定内分配地址是栈还是堆

- 逃逸分析在编译阶段完成


## 函数传递指针真的比传值效率高吗？

我们知道传递指针可以减少底层值的拷贝，可以提高效率，但是如果**拷贝的数据量小**，由于**指针传递会产生逃逸**，可能会使用堆，也可能会增加 GC 的负担，所以**传递指针不一定是高效的**。

```go

```
