---
date: 2022-10-14T11:29:54+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "内存逃逸分析"  # 文章标题
url:  "posts/go/docs/internal/compiler/escape"  # 设置网页永久链接
tags: [ "Go", "escape" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 处理入口

```go
// Escape analysis.
// Required for moving heap allocations onto stack,
// which in turn is required by the closure implementation,
// which stores the addresses of stack variables into the closure.
// If the closure does not escape, it needs to be on the stack
// or else the stack copier will not update it.
// Large values are also moved off stack in escape analysis;
// because large values may contain pointers, it must happen early.
base.Timer.Start("fe", "escapes")
escape.Funcs(typecheck.Target.Decls)
```

## 简介

逃逸分析（Escape analysis）是指**由编译器决定内存分配的位置（栈区还是堆区）**，不需要程序员指定。

Go 语言能够通过编译时的逃逸分析识别这种问题，自动将该变量放置到堆区，并借助 Go 运行时的垃圾回收机制自动释放内存。编译器会尽可能地将变量放置到栈中，因为栈中的对象随着函数调用结束会被自动销毁，减轻运行时分配和垃圾回收的负担。

## 分配原则

在 Go 语言中，开发者模糊了栈区与堆区的差别，不管是字符串、数组字面量，还是通过 new、make 标识符创建的对象，**都既可能被分配到栈中，也可能被分配到堆中**，分配时，遵循以下两个原则：

- 原则 1：**指向栈上对象的指针不能被存储到堆中**
- 原则 2：**指向栈上对象的指针不能超过该栈对象的生命周期**

如果指针存储在全局变量或者其它数据结构中，它们也可能发生逃逸，这种情况是当前程序中的指针逃逸。逃逸分析需要确定指针所有可以存储的地方，保证指针的生命周期只在当前进程或线程中。

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

很多函数参数为 `interface` 类型，比如

```go
fmt.Println(a...interface{})
```

**编译期间很难确定其参数的具体类型**，也会产生逃逸。

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

## 函数传递指针真的比传值效率高吗？

我们知道传递指针可以减少底层值的拷贝，可以提高效率，但是如果**拷贝的数据量小**，由于**指针传递会产生逃逸**，可能会使用堆，也可能会增加 GC 的负担，所以**传递指针不一定是高效的**。

## 实现原理

Go 语言通过对抽象语法树的**静态数据流分析**（static data-flow analysis）来实现逃逸分析，这种方式构建了**带权重的有向图**。

导致内存逃逸的情况比较多，通常来讲就是如果变量的作用域不会扩大并且其行为或者大小能够在编译的时候确定，一般情况下都是分配到栈上，否则就可能发生内存逃逸分配到堆上。

简单的逃逸现象举例如下：

```go
var z *int

func escape() {
    a := 1
    z = &a
}
```

在上例中，变量 z 为全局变量，是一个指针。在函数中，变量 z 引用了变量 a 的地址。如果变量 a 被分配到栈中，那么最终程序将违背原则 2，即变量 z 超过了变量 a 的生命周期，因此变量 a 最终将被分配到堆中。

可以通过在编译时加入 -m=2 标志打印出编译时的逃逸分析信息。

如下所示，表明变量 a 将被放置到堆中。

```
go tool compile -m=2 main.go
```

```
main.go:5:6: can inline escape with cost 9 as: func() { a := 1; z = &a }
main.go:10:6: can inline main with cost 0 as: func() {  }
main.go:6:2: a escapes to heap:
main.go:6:2:   flow: {heap} = &a:
main.go:6:2:     from &a (address-of) at main.go:7:6
main.go:6:2:     from z = &a (assign) at main.go:7:4
main.go:6:2: moved to heap: a
```

Go 语言在编译时构建了带权重的有向图，其中权重可以表明当前变量引用与解引用的数量。

下例为 p 引用 q 时的权重，当权重大于 0 时，代表存在 * 解引用操作。当权重为 -1 时，代表存在 & 引用操作。

```go
p = &q // -1
p = q // 0
p = *q // 1
p = **q // 2
p = **&**&q // 2
```

并不是权重为 -1 就一定要逃逸，例如在下例中，虽然 z 引用了变量 a 的地址，但是由于变量 z 并没有超过变量 a 的声明周期，因此变量 a 与变量 z 都不需要逃逸。

```go
func f() int {
	a := 1
	z := &a
	return *z
}
```

为了理解编译器带权重的有向图，再来看一个更加复杂的例子。在该案例中有多次的引用与解引用过程。

```go
package main

var o *int

func main() {
	l := new(int)
	*l = 42
	m := &l
	n := &m // &&l
	o = **n // l
}
```

最终编译器在逃逸分析中的数据流分析，会被解析成如图 1-7 所示的带权重的有向图。

![](../../../assets/images/docs/internal/compiler/escape/图1-7%20逃逸分析带权重的有向图.png)

其中，节点代表变量，边代表变量之间的赋值，箭头代表赋值的方向，边上的数字代表当前赋值的引用或解引用的个数。节点的权重=前一个节点的权重 + 箭头上的数字，例如节点 m 的权重为 2-1 = 1，而节点 l 的权重为 1-1 = 0。

遍历和计算有向权重图的目的是找到权重为 -1 的节点，例如图 1-7 中的 new(int) 节点，它的节点变量地址会被传递到根节点 o 中，这时还需要考虑逃逸分析的分配原则，o 节点为全局变量，不能被分配在栈中，因此，new(int) 节点创建的变量会被分配到堆中。

```
go tool compile -m=2 main.go
```

```
main.go:5:6: can inline main with cost 27 as: func() { l := new(int); *l = 42; m := &l; n := &m; o = *(*n) }
main.go:6:10: new(int) escapes to heap:
main.go:6:10:   flow: l = &{storage for new(int)}:
main.go:6:10:     from new(int) (spill) at main.go:6:10
main.go:6:10:     from l := new(int) (assign) at main.go:6:4
main.go:6:10:   flow: m = &l:
main.go:6:10:     from &l (address-of) at main.go:8:7
main.go:6:10:     from m := &l (assign) at main.go:8:4
main.go:6:10:   flow: n = &m:
main.go:6:10:     from &m (address-of) at main.go:9:7
main.go:6:10:     from n := &m (assign) at main.go:9:4
main.go:6:10:   flow: {heap} = **n:
main.go:6:10:     from *n (indirection) at main.go:10:7
main.go:6:10:     from *(*n) (indirection) at main.go:10:6
main.go:6:10:     from o = *(*n) (assign) at main.go:10:4
main.go:6:10: new(int) escapes to heap
```

实际的情况更加复杂，因为一个节点可能拥有多条边（例如结构体），而节点之间可能出现环。Go 语言采用 Bellman Ford 算法遍历查找有向图中权重小于 0 的节点，核心逻辑位于 escape/escape.go 中。

```go

```
