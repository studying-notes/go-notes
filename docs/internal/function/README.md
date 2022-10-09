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

函数是程序中为了执行特定任务而存在的一系列执行代码。函数接受输入并返回输出，执行程序的过程可以看作一系列函数的调用过程。Go 语言中最重要的函数为 main 函数，其是程序执行用户代码的入口，在每个程序中都需要存在。

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

接口定义了方法的集合，这些方法依托于运行时的接口对象，因此**接口对应的方法是在运行时动态绑定的**。Go 语言通过**隐式接口机制**实现了鸭子面向对象模型。

## 函数基本使用方式

使用函数具有减少冗余、隐藏信息、提高代码清晰度等优点。在 Go 语言中，函数是一等公民（first-class），这意味着可以将它看作变量，并且它可以作为参数传递、返回及赋值。

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

## 返回值命名

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

其中 defer 语句延迟执行了一个匿名函数，因为这个匿名函数**捕获了外部函数的局部变量** v，这种函数我们一般称为闭包。

## 函数闭包与陷阱

Go语言同样支持匿名函数和闭包。闭包（Closure）是在函数作为一类公民的编程语言中实现词法绑定的一种技术，闭包包含了函数的入口地址和其关联的环境。闭包和普通函数最大的区别在于，闭包函数中可以引用闭包外的变量。

**当匿名函数引用了外部作用域中的变量时就成了闭包函数**，闭包函数是函数式编程语言的核心。

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
        i := i // 定义一个循环体内局部变量i
        defer func(){ println(i) } ()
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

第二种方式是将迭代变量通过闭包函数的参数传入，defer 语句会马上对调用参数求值。两种方式都是可以工作的。不过一般来说，在 for 循环内部执行 defer 语句并不是一个好的习惯，此处仅为示例，不建议使用。

## 参数传值与传引用

Go 语言中，如果以切片为参数调用函数，有时候会给人一种参数采用了传引用的方式的假象：因为在被调用函数内部可以修改传入的切片的元素。其实，**任何可以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数**。

函数参数传值的规范更准确说是**只针对数据结构中固定的部分传值**，例如字符串或切片对应结构体中的指针和字符串长度结构体传值，但是并**不包含指针间接指向的内容**。将切片类型的参数替换为类似 `reflect.SliceHeader` 结构体就能很好理解切片传值的含义了：

```go
func twice(x []int) {
    for i := range x {
        x[i] *= 2
    }
}

type IntSliceHeader struct {
    Data []int
    Len  int
    Cap  int
}

func twice(x IntSliceHeader) {
    for i := 0; i < x.Len; i++ {
        x.Data[i] *= 2
    }
}
```

因为切片中的底层数组部分通过隐式指针传递（指针本身依然是传值的，但是指针指向的却是同一份的数据），所以被调用函数可以通过指针修改调用参数切片中的数据。除数据之外，切片结构还包含了切片长度和切片容量信息，这两个信息也是传值的。如果被调用函数中修改了 Len 或 Cap 信息，就无法反映到调用参数的切片中，这时候我们一般会通过返回修改后的切片来更新之前的切片。这也是内置的 `append ()` 必须要返回一个切片的原因。

## 函数栈

栈在不同的场景具有不同的含义，有时候指一种先入后出的数据结构，有时候指操作系统组织内存的形式。

在大多数现代计算机系统中，**每个线程都有一个被称为栈的内存区域**，其遵循一种先入先出（FIFO）的形式，增长方向为从高地址到低地址。

当函数执行时，函数的参数、返回地址、局部变量会被压入栈中，当函数退出时，这些数据会被回收。当函数还没有退出就调用另一个函数时，形成了一条函数调用链。

例如，函数 A 调用了函数 B，被调函数 B 至少需要存储调用方函数 A 提供的返回地址的位置，以便在函数 B 执行完毕后，能够立即返回函数 A 之前的位置继续执行。

每个函数在执行过程中都使用一块栈内存来保存**返回地址、局部变量、函数参数**等，我们将这一块区域称为函数的栈帧（stack frame）。

当发生函数调用时，因为调用函数没有执行完毕，其栈内存中保存的数据还有用，所以被调用函数不能覆盖调用函数的栈帧，只能把被调用函数的栈帧压栈，等被调用函数执行完毕后再让栈帧出栈。这样，栈的大小就会随着函数调用层级的增加而扩大，随函数的返回而缩小，也就是说，函数的调用层级越深，消耗的栈空间越大。

因为数据是以先进先出的方式添加和删除的，所以基于栈的内存分配非常简单，并且通常比基于堆的动态内存分配快得多。另外，当函数退出时，栈上的内存会自动高效地回收，这是垃圾回收最初的形式。

维护和管理函数的栈帧非常重要，对于高级编程语言来说，栈帧通常是隐藏的。例如，Go 语言借助编译器，在开发中不用关心局部变量在栈中的布局与释放。

许多计算机指令集在硬件级别提供了用于管理栈的特殊指令，例如，80x86 指令集提供的 SP 用于管理栈，以 A 函数调用 B 函数为例，普遍的函数栈结构如图 9-1 所示。

![](../../../assets/images/docs/internal/function/README/图9-1%20普遍的函数调用栈结构.png)

## 栈帧结构

```go
package main

func mul(a, b int) int {
	return a * b
}

func main() {
	mul(3, 4)
}
```

由于高级语言为开发者隐藏了函数调用的细节，所以分析栈结构需要用一些特殊的手段，可以通过调试器或者打印汇编代码的方式进行分析。上例中的程序虽然简单，但通过编译器优化后，会识别出并不需要调用该函数，导致查看汇编代码时不能得到想要的结果。因此在调试时需要禁止编译器的优化及函数内联。

```
go tool compile -S -N -l main.go
```

> 不同系统上的汇编结果不一样。

```
main.main STEXT size=54 args=0x0 locals=0x18 funcid=0x0 align=0x0
	0x0000 00000 (main.go:7)	TEXT	main.main(SB), ABIInternal, $24-0
	0x0000 00000 (main.go:7)	CMPQ	SP, 16(R14)
	0x0004 00004 (main.go:7)	PCDATA	$0, $-2
	0x0004 00004 (main.go:7)	JLS	47
	0x0006 00006 (main.go:7)	PCDATA	$0, $-1
	0x0006 00006 (main.go:7)	SUBQ	$24, SP
	0x000a 00010 (main.go:7)	MOVQ	BP, 16(SP)
	0x000f 00015 (main.go:7)	LEAQ	16(SP), BP
	0x0014 00020 (main.go:7)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0014 00020 (main.go:7)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0014 00020 (main.go:8)	MOVL	$3, AX
	0x0019 00025 (main.go:8)	MOVL	$4, BX
	0x001e 00030 (main.go:8)	PCDATA	$1, $0
	0x001e 00030 (main.go:8)	NOP
	0x0020 00032 (main.go:8)	CALL	main.mul(SB)
	0x0025 00037 (main.go:9)	MOVQ	16(SP), BP
	0x002a 00042 (main.go:9)	ADDQ	$24, SP
	0x002e 00046 (main.go:9)	RET
	0x002f 00047 (main.go:9)	NOP
	0x002f 00047 (main.go:7)	PCDATA	$1, $-1
	0x002f 00047 (main.go:7)	PCDATA	$0, $-2
	0x002f 00047 (main.go:7)	CALL	runtime.morestack_noctxt(SB)
	0x0034 00052 (main.go:7)	PCDATA	$0, $-1
	0x0034 00052 (main.go:7)	JMP	0
```

第 2 行代码中声明的 $32-0 与函数的栈帧有关，其中， $32 代表当前栈帧会被分配的字节数，后面的 0 代表函数参数与返回值的字节数。由于 main 函数中没有参数也没有返回值，因此为 0。第 3 行代码 SUBQ $32，SP 将当前的 SP 寄存器减去 32 字节，这意味着当前的函数栈增加了 32 字节，图 9-2 描述了该例中 main 函数的栈帧结构。

![](../../../assets/images/docs/internal/function/README/图9-2%20main函数的栈帧结构.png)

第 4 行的 `MOVQ BP, 24(SP)` 用于将当前 BP 寄存器的值存储到栈帧的顶部。并在第 5 行通过 `LEAQ 24(SP)`，将当前 BP 寄存器的值指向栈基地址。第 6 行和第 7 行将需要调用的函数 mul 的参数存储到栈顶。

调用 mul 函数前，BP 寄存器与函数参数的位置如图 9-3 所示，从图中可以看出，函数调用时参数的压栈操作是从右到左进行的。`mul(3, 4)` 中的第 2 个参数 4 首先压栈，随后第 1 个参数 3 压栈。

![](../../../assets/images/docs/internal/function/README/图9-3%20调用mul函数前BP寄存器与函数参数的位置.png)

## 函数调用链结构与特性

当 main 函数调用 mul 函数时，进入一个新的栈帧，mul 函数的汇编代码如下。

```
main.mul STEXT nosplit size=60 args=0x10 locals=0x10 funcid=0x0 align=0x0
	0x0000 00000 (main.go:3)	TEXT	main.mul(SB), NOSPLIT|ABIInternal, $16-16
	0x0000 00000 (main.go:3)	SUBQ	$16, SP
	0x0004 00004 (main.go:3)	MOVQ	BP, 8(SP)
	0x0009 00009 (main.go:3)	LEAQ	8(SP), BP
	0x000e 00014 (main.go:3)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x000e 00014 (main.go:3)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x000e 00014 (main.go:3)	FUNCDATA	$5, main.mul.arginfo1(SB)
	0x000e 00014 (main.go:3)	MOVQ	AX, main.a+24(SP)
	0x0013 00019 (main.go:3)	MOVQ	BX, main.b+32(SP)
	0x0018 00024 (main.go:3)	MOVQ	$0, main.~r0(SP)
	0x0020 00032 (main.go:4)	MOVQ	main.b+32(SP), AX
	0x0025 00037 (main.go:4)	MOVQ	main.a+24(SP), CX
	0x002a 00042 (main.go:4)	IMULQ	CX, AX
	0x002e 00046 (main.go:4)	MOVQ	AX, main.~r0(SP)
	0x0032 00050 (main.go:4)	MOVQ	8(SP), BP
	0x0037 00055 (main.go:4)	ADDQ	$16, SP
	0x003b 00059 (main.go:4)	RET
```

第 1 行汇编代码中的 $0-24 表明当前函数栈帧中不会被分配任何字节，但是参数和返回值的大小为 24 字节，这 24 字节存储在调用者 main 函数的栈帧中。

在上一小节的汇编代码中，`CALL "".mul(SB)` 用于执行对于 mul 函数的调用，该指令有一个隐含的操作是将 SP 寄存器减 8，并存储其返回地址。该指令是 mul 函数返回后 main 函数执行的下一条指令。当 mul 函数返回时，会执行 RET 指令。该指令暗含着获取存储在栈帧顶部的返回地址，并跳转到该处执行的操作。所以在如图 9-4 所示 mul 函数还未返回时的栈帧结构中，a、b 两个参数分别对应 8(SP)、16(SP) 所在的位置，并且最后将返回值存储在 24(SP) 处。

![](../../../assets/images/docs/internal/function/README/图9-4%20调用mul函数还未返回时的栈帧结构.png)

当 main 函数返回时，MOVQ 24(SP)，BP 指令会还原 main 函数的调用者的 BP 地址，并且 `ADDQ $32`，SP 将为 SP 寄存器加上 32 字节，意味着收缩栈。最后 RET 指令会跳转到调用者函数处继续执行。当然，由于是 main 函数，所以执行程序退出操作。

Go 语言函数的参数和返回值存储在栈中，然而许多主流的编程语言会将参数和返回值存储在寄存器中。存储在栈中的好处在于所有平台都可以使用相同的约定，从而容易开发出可移植、跨平台的代码，同时这种方式简化了协程、延迟调用和反射调用的实现。寄存器的值不能跨函数调用、存活，这简化了垃圾回收期间的栈扫描和对栈扩容的处理。

将参数和返回值存储在栈中的约定存在一些性能问题。尽管现代高性能 CPU 在很大程度上优化了栈访问，但是访问寄存器中的数据仍比访问栈中的数据快 40%。此外，这种调用约定引起了额外的内存通信，降低了效率。

## 堆栈信息

下例通过简单的函数调用过程为 main 函数调用了 trace 函数。

```go
package main

func trace(array []int, a, b int) int {
	panic("not implemented")
	return 0
}

func main() {
	trace([]int{1, 2, 3, 4, 5}, 1, 3)
}
```

运行该程序时，会输出如下堆栈信息。注意，此时使用了-gcflags="-l" 禁止函数的内联优化，否则内联函数中不会打印函数的参数，在运行时会输出当前协程所在的堆栈。

```
go run -gcflags="-l" main.go
```

```
panic: not implemented

goroutine 1 [running]:
main.trace({0x0?, 0x4047f8?, 0xc000042770?}, 0x404739?, 0x60?)
        main.go:4 +0x27
main.main()
        main.go:9 +0x6f
exit status 2
```

其中，输出的第 1 行 `panic: not implemented` 会给出程序终止运行的原因。

`goroutine 1 [running]` 代表当前协程的 ID 及状态，触发堆栈信息的协程将会放在最上方。

接下来是当前协程函数调用链的具体信息。main.trace 为当前协程正在运行的函数，函数后面可以看到传递的具体参数。

trace 函数看起来有 3 个参数，但是由于切片在运行时的结构如下：

```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

所以在函数传递的过程中其实完成了一次该结构的复制。

第 1 个参数 0xc000046760 代表切片的地址，第 2 个参数 0x3 代表切片的长度为 3，第 3 个参数 0x3 代表切片的容量为 3。接下来是参数 a、b。

main.go:4+0x27 代表当前函数所在的文件位置及行号。其中，+0x27 代表当前函数中下一个要执行的指令的位置，其是距离当前函数起始位置的偏移量。

在堆栈信息中，还会依次列出调用 trace 函数的函数调用链。例如，trace 函数的调用者为 main 函数，同样会打印出 main 函数调用 trace 函数的文件所在的位置和行号。+0x6f 也是 main 函数的 PC 偏移量，对应着 trace 函数返回后，main 函数将执行的下一条指令。

可以看出，堆栈信息是一种非常有用的排查问题的方法，同时，可以通过函数参数信息得知函数调用时传递的参数，帮助读者学习 Go 语言内部类型的结构，以及值传递和指针传递的区别。值得一提的是，Go 语言可以通过配置 GOTRACEBACK 环境变量在程序异常终止时生成 core dump 文件，生成的文件可以由 dlv 或者 gdb 等高级的调试工具进行分析调试。

## 栈扩容与栈转移原理

Go 语言在线程的基础上实现了用户态更加轻量级的协程，线程的栈大小一般是在创建时指定的，为了避免出现栈溢出（stack overflow）的错误，默认的栈大小会相对较大（例如 2MB）。

而在 Go 语言中，每个协程都有一个栈，并且在 Go 1.4 之后，每个栈的大小在初始化的时候都是 2KB。程序中经常有成千上万的协程存在，可以预料到，Go 语言中的栈是可以扩容的。

Go 语言中最大的协程栈在 64 位操作系统中为 1GB，在 32 位系统中为 250MB。在 Go 语言中，栈的大小不用开发者手动调整，都是在运行时实现的。栈的管理有两个重要的问题：触发扩容的时机及栈调整的方式。

在函数序言阶段判断是否需要对栈进行扩容，由编译器自动插入判断指令，如果满足适当的条件则对栈进行扩容。

执行 main 函数的开始阶段，首先从线程局部存储中获取代表当前协程的结构体 g。

```go
// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic    *_panic // innermost panic - offset known to liblink
	_defer    *_defer // innermost defer
	m         *m      // current m; offset known to arm liblink
	sched     gobuf
	syscallsp uintptr // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc uintptr // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp  uintptr // expected sp at top of stack, to check in traceback
	// param is a generic pointer parameter field used to pass
	// values in particular contexts where other storage for the
	// parameter would be difficult to find. It is currently used
	// in three ways:
	// 1. When a channel operation wakes up a blocked goroutine, it sets param to
	//    point to the sudog of the completed blocking operation.
	// 2. By gcAssistAlloc1 to signal back to its caller that the goroutine completed
	//    the GC cycle. It is unsafe to do so in any other way, because the goroutine's
	//    stack may have moved in the meantime.
	// 3. By debugCallWrap to pass parameters to a new goroutine because allocating a
	//    closure in the runtime is forbidden.
	param        unsafe.Pointer
	atomicstatus atomic.Uint32
	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid         uint64
	schedlink    guintptr
	waitsince    int64      // approx time when the g become blocked
	waitreason   waitReason // if status==Gwaiting

	preempt       bool // preemption signal, duplicates stackguard0 = stackpreempt
	preemptStop   bool // transition to _Gpreempted on preemption; otherwise, just deschedule
	preemptShrink bool // shrink stack at synchronous safe point

	// asyncSafePoint is set if g is stopped at an asynchronous
	// safe point. This means there are frames on the stack
	// without precise pointer information.
	asyncSafePoint bool

	paniconfault bool // panic (instead of crash) on unexpected fault address
	gcscandone   bool // g has scanned stack; protected by _Gscan bit in status
	throwsplit   bool // must not split stack
	// activeStackChans indicates that there are unlocked channels
	// pointing into this goroutine's stack. If true, stack
	// copying needs to acquire channel locks to protect these
	// areas of the stack.
	activeStackChans bool
	// parkingOnChan indicates that the goroutine is about to
	// park on a chansend or chanrecv. Used to signal an unsafe point
	// for stack shrinking.
	parkingOnChan atomic.Bool

	raceignore     int8     // ignore race detection events
	sysblocktraced bool     // StartTrace has emitted EvGoInSyscall about this goroutine
	tracking       bool     // whether we're tracking this G for sched latency statistics
	trackingSeq    uint8    // used to decide whether to track this G
	trackingStamp  int64    // timestamp of when the G last started being tracked
	runnableTime   int64    // the amount of time spent runnable, cleared when running, only used when tracking
	sysexitticks   int64    // cputicks when syscall has returned (for tracing)
	traceseq       uint64   // trace event sequencer
	tracelastp     puintptr // last P emitted an event for this goroutine
	lockedm        muintptr
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
	gopc           uintptr         // pc of go statement that created this goroutine
	ancestors      *[]ancestorInfo // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
	startpc        uintptr         // pc of goroutine function
	racectx        uintptr
	waiting        *sudog         // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
	cgoCtxt        []uintptr      // cgo traceback context
	labels         unsafe.Pointer // profiler labels
	timer          *timer         // cached timer for time.Sleep
	selectDone     atomic.Uint32  // are we participating in a select and did someone win the race?

	// goroutineProfiled indicates the status of this goroutine's stack for the
	// current in-progress goroutine profile
	goroutineProfiled goroutineProfileStateHolder

	// Per-G GC state

	// gcAssistBytes is this G's GC assist credit in terms of
	// bytes allocated. If this is positive, then the G has credit
	// to allocate gcAssistBytes bytes without assisting. If this
	// is negative, then the G must correct this by performing
	// scan work. We track this in bytes to make it fast to update
	// and check for debt in the malloc hot path. The assist ratio
	// determines how this corresponds to scan work debt.
	gcAssistBytes int64
}
```

可以看到结构体 g 的第 1 个成员 stack 占 16 字节（lo 和 hi 各占 8 字节），所以结构体 g 变量的起始位置偏移 16 就对应到 stackguard0 字段。main 函数的第 2 条指令 `CMPQ SP, 16(CX)` 会比较栈顶寄存器 rsp 的值是否比 stackguard0 的值小，如果 rsp 的值更小，则说明当前 g 的栈要用完了，有溢出风险，需要调用 `morestack_noctxt` 函数来扩栈。

stackguard0 会在初始化时将 `stack.lo+_StackGuard`，`_StackGuard` 设置为 896 字节，stack.lo 为当前栈的栈顶。如果出现图 9-5 中栈寄存器 SP 小于 stackguard0 的情况，则表明当前栈空间不够，stackguard0 除了用于栈的扩容，还用于协程抢占。

![](../../../assets/images/docs/internal/function/README/图9-5%20栈寄存器SP小于stackguard0.png)

`src/runtime/stack.go`

可以看到，在函数序言阶段如果判断出需要扩容，则会跳转调用运行时 morestack_noctxt 函数，函数调用链为 `morestack_noctxt()->morestack()->newstack()`，核心代码位于 newstack 函数中。Newstack 函数不仅会处理扩容，还会处理协程的抢占。

```go
	// Allocate a bigger segment and move the stack.
	oldsize := gp.stack.hi - gp.stack.lo
	newsize := oldsize * 2

	// Make sure we grow at least as much as needed to fit the new frame.
	// (This is just an optimization - the caller of morestack will
	// recheck the bounds on return.)
	if f := findfunc(gp.sched.pc); f.valid() {
		max := uintptr(funcMaxSPDelta(f))
		needed := max + _StackGuard
		used := gp.stack.hi - gp.sched.sp
		for newsize-used < needed {
			newsize *= 2
		}
	}

	if stackguard0 == stackForceMove {
		// Forced stack movement used for debugging.
		// Don't double the stack (or we may quickly run out
		// if this is done repeatedly).
		newsize = oldsize
	}

	if newsize > maxstacksize || newsize > maxstackceiling {
		if maxstacksize < maxstackceiling {
			print("runtime: goroutine stack exceeds ", maxstacksize, "-byte limit\n")
		} else {
			print("runtime: goroutine stack exceeds ", maxstackceiling, "-byte limit\n")
		}
		print("runtime: sp=", hex(sp), " stack=[", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
		throw("stack overflow")
	}

	// The goroutine must be executing in order to call newstack,
	// so it must be Grunning (or Gscanrunning).
	casgstatus(gp, _Grunning, _Gcopystack)

	// The concurrent GC will not scan the stack while we are doing the copy since
	// the gp is in a Gcopystack status.
	copystack(gp, newsize)
	if stackDebug >= 1 {
		print("stack grow done\n")
	}
	casgstatus(gp, _Gcopystack, _Grunning)
	gogo(&gp.sched)
```

如上所示，newstack 函数首先通过栈底地址与栈顶地址相减计算旧栈的大小，并计算新栈的大小。新栈为旧栈的两倍大小。在 64 位操作系统中，如果栈大小超过了 1GB 则直接报错为 stack overflow。

栈扩容的重要一步是将旧栈的内容转移到新栈中。栈扩容首先将协程的状态设置为 _Gcopystack，以便在垃圾回收状态下不会扫描该协程栈带来错误。栈复制并不像直接复制内存那样简单，如果栈中包含了引用栈中其他地址的指针，那么该指针需要对应到新栈中的地址，copystack 函数会分配一个新栈的内存。为了应对频繁的栈调整，对获取栈的内存进行了许多优化，特别是对小栈。在 Linux 操作系统下，会对 2KB/4KB/8KB/16KB 的小栈进行专门的优化，即在全局及每个逻辑处理器(P)中预先分配这些小栈的缓存池，从而避免频繁地申请堆内存。

![](../../../assets/images/docs/internal/function/README/图9-6%20栈的全局与本地缓存池结构.png)

栈的全局与本地缓存池结构如图 9-6 所示，每个逻辑处理器中的缓存池都来自全局缓存池(stackpool)。mcache 有时可能不存在（例如在调整 P 的大小后），这时需要直接从全局缓存池获取栈缓存。对于大栈，其大小不确定，虽然也有一个全局的缓存池，但不会预先放入多个栈，当栈被销毁时，如果被销毁的栈为大栈则放入全局缓存池中。当全局缓存池中找不到对应大小的栈时，会从堆区分配。

在分配到新栈后，如果有指针指向旧栈，那么需要将其调整到新栈中。在调整时有一个额外的步骤是调整 sudog，由于通道在阻塞的情况下存储的元素可能指向了栈上的指针，因此需要调整。接着需要将旧栈的大小复制到新栈中，这涉及借助 memmove 函数进行内存复制。

内存复制完成后，需要调整当前栈的 SP 寄存器和新的 stackguard0，并记录新的栈顶与栈底。扩容最关键的一步是在新栈中调整指针。因为新栈中的指针可能指向旧栈，旧栈一旦释放就会出现严重的问题。图 9-7 描述了栈扩容的过程，copystack 函数会遍历新栈上所有的栈帧信息，并遍历其中所有可能有指针的位置。一旦发现指针指向旧栈，就会调整当前的指针使其指向新栈。

![](../../../assets/images/docs/internal/function/README/图9-7%20栈扩容的过程.png)

栈的转移如图 9-8 所示，调整后，栈指针将指向新栈中的地址。

![](../../../assets/images/docs/internal/function/README/图9-8%20栈的转移.png)

## 栈调试

一种特别的方式是在源码级别进行调试，Go 语言在源码级别提供了栈相关的多种级别的调试、用户调试栈的扩容及栈的分配等。但是这些静态常量并没有暴露给用户，如果要使用这些变量，则需要直接修改 Go 的源码并重新进行编译。

`src/runtime/stack.go`

```go
const (
	// stackDebug == 0: no logging
	//            == 1: logging of per-stack operations
	//            == 2: logging of per-frame operations
	//            == 3: logging of per-word updates
	//            == 4: logging of per-word reads
	stackDebug       = 0
	stackFromSystem  = 0 // allocate stacks from system memory instead of the heap
	stackFaultOnFree = 0 // old stacks are mapped noaccess to detect use after free
	stackPoisonCopy  = 0 // fill stack that should not be accessed with garbage, to detect bad dereferences during copy
	stackNoCache     = 0 // disable per-P small stack caches

	// check the BP links during traceback.
	debugCheckBP = false
)
```

```go

```
