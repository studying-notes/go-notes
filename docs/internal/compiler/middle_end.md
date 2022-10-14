---
date: 2022-10-01T10:47:26+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "中间端"  # 文章标题
url:  "posts/go/docs/internal/compiler/middle-end"  # 设置网页永久链接
tags: [ "Go", "middle-end" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 源码分布

这一阶段的相关源码分布在：

* `cmd/compile/internal/deadcode` (dead code elimination)
* `cmd/compile/internal/inline` (function call inlining)
* `cmd/compile/internal/devirtualize` (devirtualization of known interface method calls)
* `cmd/compile/internal/escape` (escape analysis)

中间端对 **IR 表示**执行了几个优化过程：死代码消除、函数调用内联、去虚拟化和逃逸分析。

## 死代码消除

```go
	// Eliminate some obviously dead code.
	// Must happen after typechecking.
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			deadcode.Func(n.(*ir.Func))
		}
	}
```

## 函数调用内联

```go
	// Inlining
	base.Timer.Start("fe", "inlining")
	if base.Flag.LowerL != 0 {
		inline.InlinePackage()
	}
	noder.MakeWrappers(typecheck.Target) // must happen after inlining
```

```go
// InlinePackage finds functions that can be inlined and clones them before walk expands them.
func InlinePackage() {
	ir.VisitFuncsBottomUp(typecheck.Target.Decls, func(list []*ir.Func, recursive bool) {
		numfns := numNonClosures(list)
		for _, n := range list {
			if !recursive || numfns > 1 {
				// We allow inlining if there is no
				// recursion, or the recursion cycle is
				// across more than one function.
				CanInline(n)
			} else {
				if base.Flag.LowerM > 1 {
					fmt.Printf("%v: cannot inline %v: recursive\n", ir.Line(n), n.Nname)
				}
			}
			InlineCalls(n)
		}
	})
}
```

函数内联指将较小的函数直接组合进调用者的函数。这是现代编译器优化的一种核心技术。函数内联的优势在于，可以减少函数调用带来的开销。

对于 Go 语言来说，函数调用的成本在于参数与返回值栈复制、较小的栈寄存器开销以及函数序言部分的检查栈扩容（Go 语言中的栈是可以动态扩容的）。

同时，函数内联是其他一些编译器优化的基础。我们可以通过一段简单的程序衡量函数内联带来的效率提升，如下所示，使用 bench 对 max 函数调用进行测试。当我们在函数的注释前方加上 `//go:noinline` 时，代表当前函数是禁止进行函数内联优化的。取消该注释后，max 函数将会对其进行内联优化。

```go
package main

import "testing"

//go:noinline
func MaxNoinline(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxInline(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var result int

func BenchmarkMaxNoinline(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		result = MaxNoinline(-1, i)
	}
	result = r
}

func BenchmarkMaxInline(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = MaxInline(-1, i)
	}
	result = r
}
```

通过下面的 bench 对比结果可以看出，在内联后，max 函数的执行时间显著少于非内联函数调用花费的时间，这里的消耗主要来自函数调用增加的执行指令。

```
go test -bench=".*"
```

```go
BenchmarkMaxNoinline-12         787131966                1.580 ns/op
BenchmarkMaxInline-12           1000000000               0.3520 ns/op
```

Go 语言编译器会计算函数内联花费的成本，只有执行相对简单的函数时才会内联。

函数内联的核心逻辑位于 inline/inl.go 中。

当函数内部有 for、range、go、select 等语句时，该函数不会被内联，当函数执行过于复杂（例如太多的语句或者函数为递归函数）时，也不会执行内联。

另外，如果函数前的注释中有 `go:noinline` 标识，则该函数不会执行内联。如果希望程序中所有的函数都不执行内联操作，那么可以添加编译器选项“-l”。

```
go build -gcflags="-l" main.go
```

```
go tool compile -l main.go
```

在调试时，可以获取当前函数是否可以内联，以及不可以内联的原因。

```go
package main

func small() string {
	return "small"
}

func fib(index int) int {
	if index < 2 {
		return index
	}
	return fib(index-1) + fib(index-2)
}

func main() {
	small()
	fib(10)
}
```

在上面的代码中，当在编译时加入 -m=2 标志时，可以打印出函数的内联调试信息。可以看出，small 函数可以被内联，而 fib（斐波那契）函数为递归函数，不能被内联。

```
go tool compile -m=2 main.go | grep inline
```

```
main.go:3:6: can inline small with cost 2 as: func() string { return "small" }
main.go:7:6: cannot inline fib: recursive
main.go:14:6: can inline main with cost 64 as: func() { small(); fib(10) }
```

当函数可以被内联时，该函数将被纳入调用函数。

如下所示，a:=b+f(1），其中，f 函数可以被内联。

```go
func f(n int) int {
	return n + 1
}

func main() {
	b := 1
	a := b + f(1)
}
```

函数参数与返回值在编译器内联阶段都将转换为声明语句，并通过 goto 语义跳转到调用者函数语句中，上述代码的转换形式如下，在后续编译器阶段还将对该内联结构做进一步优化。

```
n := 1
~r1 := n + 1
goto end
end:
  a := b + ~r1
```

## 去虚拟化

```go
	// Devirtualize.
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			devirtualize.Func(n.(*ir.Func))
		}
	}
	ir.CurFunc = nil
```

## 逃逸分析

[逃逸分析](escape.md)

## 变量捕获

变量捕获主要是针对闭包场景而言的，由于闭包函数中可能引用闭包外的变量，因此变量捕获需要明确在闭包中通过值引用或地址引用的方式来捕获变量。

下面的例子中有一个闭包函数，在闭包内引入了闭包外的 a、b 变量，由于变量 a 在闭包之后进行了其他赋值操作，因此在闭包中，a、b 变量的引用方式会有所不同。在闭包中，必须采取地址引用的方式对变量 a 进行操作，而对变量 b 的引用将通过直接值传递的方式进行。

```go
package main

import "fmt"

func main() {
	a := 1
	b := 2

	go func() {
		fmt.Println(a, b)
	}()

	a = 3
}
```

在 Go 语言编译的过程中，可以通过如下方式查看当前程序闭包变量捕获的情况。从输出中可以看出，a 采取 ref 引用传递的方式，而 b 采取了值传递的方式。assign=true 代表变量 a 在闭包完成后又进行了赋值操作。

```
go tool compile -m=2 main.go | grep capturing
```

```
main.go:6:2: main capturing by ref: a (addr=false assign=true width=8)
main.go:7:2: main capturing by value: b (addr=false assign=false width=8)
```

## 闭包重写

在前面的阶段，编译器完成了闭包变量的捕获用于决定是通过指针引用还是值引用的方式传递外部变量。在完成逃逸分析后，下一个优化的阶段为闭包重写。

闭包重写分为闭包定义后被立即调用和闭包定义后不被立即调用两种情况。在闭包被立即调用的情况下，闭包只能被调用一次，这时可以将闭包转换为普通函数的调用形式。

```go
func do() {
	a := 1
	func() {
		fmt.Println(a)
		a = 2
	}()
}
```

上面的闭包最终会被转换为类似正常函数调用的形式，如下所示，由于变量 a 为引用传递，因此构造的新的函数参数应该为 int 指针类型。如果变量是值引用的，那么构造的新的函数参数应该为 int 类型。

```go
func do() {
	a := 1
	func1(&a)
}

func func1(a *int) {
	fmt.Println(*a)
	*a = 2
}
```

如果闭包定义后不被立即调用，而是后续调用，那么同一个闭包可能被调用多次，这时需要创建闭包对象。

如果变量是按值引用的，并且该变量占用的存储空间小于 2×sizeof(int)，那么通过在函数体内创建局部变量的形式来产生该变量。如果变量通过指针或值引用，但是占用存储空间较大，那么捕获的变量(var)转换成指针类型的“&var”。这两种方式都需要在函数序言阶段将变量初始化为捕获变量的值。
