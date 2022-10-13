---
date: 2022-10-12T16:07:30+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "函数栈"  # 文章标题
url:  "posts/go/docs/internal/function/stack"  # 设置网页永久链接
tags: [ "Go", "stack" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 栈

栈在算法中指一种先入后出的数据结构，在操作系统中指组织内存的形式。

**每个系统线程都有一个被称为栈的内存区域**，其遵循一种先入先出(FIFO)的形式，增长方向为从高地址到低地址。

## 进程的内存布局

严格说来这里讲的是进程在虚拟地址空间中的布局。

操作系统把磁盘上的可执行文件加载到内存运行之前，会做很多工作，其中很重要的一件事情就是把可执行文件中的代码，数据放在内存中合适的位置，并分配和初始化程序运行过程中所必须的堆栈，所有准备工作完成后操作系统才会调度程序起来运行。

程序运行时在内存中的布局图：

![](../../../assets/images/docs/internal/goroutine/scheduler/function_call_stack/141394582f737218.png)

进程在内存中的布局主要分为 4 个区域：代码区，数据区，堆和栈。

- **代码区**，包括能被 CPU 执行的机器代码（指令）和**只读数据**比如**字符串常量**，程序一旦加载完成代码区的大小就不会再变化了。

- **数据区**，包括程序的全局变量和静态变量（c 语言有静态变量，go 没有），与代码区一样，程序加载完毕后数据区的大小也不会发生改变。

- **堆**，程序**运行时动态分配的内存都位于堆中**，这部分内存由内存分配器负责管理。**该区域的大小会随着程序的运行而变化**，即当我们向堆请求分配内存但分配器发现堆中的内存不足时，它会向操作系统内核申请向高地址方向扩展堆的大小，而当我们释放内存把它归还给堆时如果内存分配器发现剩余空闲内存太多则又会向操作系统请求向低地址方向收缩堆的大小。从这个内存申请和释放流程可以看出，我们从堆上分配的内存用完之后必须归还给堆，否则内存分配器可能会反复向操作系统申请扩展堆的大小从而导致堆内存越用越多，最后出现内存不足，这就是所谓的内存泄漏。

## 函数调用栈

函数调用栈简称栈，在程序运行过程中，不管是函数的执行还是函数调用，栈都起着非常重要的作用，它主要被用来：

- 保存函数的**局部变量**；
- 向被调用函数**传递参数**；
- 返回函数的**返回值**；
- 保存函数的**返回地址**。

返回地址是指从被调用函数返回后调用者应该继续执行的指令地址。

当函数执行时，函数的参数、返回地址、局部变量会被压入栈中，**当函数退出时，这些数据会被回收**。当函数还没有退出就调用另一个函数时，形成了一条函数调用链。

例如，函数 A 调用了函数 B，被调函数 B 至少需要存储调用方函数 A 提供的返回地址的位置，以便在函数 B 执行完毕后，能够立即返回函数 A 之前的位置继续执行。

每个函数在执行过程中都使用一块栈内存来保存**返回地址、局部变量、函数参数**等，我们将这一块区域称为函数的**栈帧**(stack frame)。

当发生函数调用时，因为调用函数没有执行完毕，其栈内存中保存的数据还有用，所以被调用函数不能覆盖调用函数的栈帧，只能把被调用函数的栈帧压栈，等被调用函数执行完毕后再让栈帧出栈。这样，**栈的大小就会随着函数调用层级的增加而扩大，随函数的返回而缩小**，也就是说，**函数的调用层级越深，消耗的栈空间越大**。

因为数据是以先进先出的方式添加和删除的，所以**基于栈的内存分配相对简单**，并且通常比基于堆的动态内存分配快得多。另外，当函数退出时，**栈上的内存会自动高效地回收**，这是垃圾回收最初的形式。

维护和管理函数的栈帧非常重要，对于高级编程语言来说，栈帧通常是隐藏的。例如，Go 语言借助编译器，在开发中不用关心局部变量在栈中的布局与释放。

许多计算机指令集在硬件级别提供了用于管理栈的特殊指令，例如，x86 指令集提供的 SP 用于管理栈，以 A 函数调用 B 函数为例，普遍的函数栈结构如图 9-1 所示。

![](../../../assets/images/docs/internal/function/stack/图9-1%20普遍的函数调用栈结构.png)

## 栈之间的关系

下面用两个图例来说明一下函数调用栈以及 SP/BP 与栈之间的关系。假设现在有如下函数调用链且正在执行函数 C()：

```
A()->B()->C()
```

则函数 ABC 的栈帧以及 rsp/rbp 的状态大致如下图所示：

![image](../../../assets/images/docs/internal/goroutine/scheduler/function_call_stack/5321d64484f5f889.png)

对于上图，有几点需要说明一下：

- 调用函数时，**参数和返回值都是存放在调用者的栈帧之中**，而不是在被调函数之中；
- 目前正在执行 C 函数，且函数调用链为 A()->B()->C()，所以以栈帧为单位来看的话，C 函数的栈帧目前位于栈顶；
- CPU 硬件寄存器 rsp 指向整个栈的栈顶，当然它也指向 C 函数的栈帧的栈顶，而 rbp 寄存器指向的是 C 函数栈帧的起始位置；
- 虽然图中 ABC 三个函数的栈帧看起来都差不多大，但事实上在真实的程序中，每个函数的栈帧大小可能都不同，因为不同的函数局部变量的个数以及所占内存的大小都不尽相同；
- 有些编译器比如 **gcc 会把参数和返回值放在寄存器中而不是栈中**，**go 语言中函数的参数和返回值都是放在栈上的**；

随着程序的运行，如果 C、B 两个函数都执行完成并返回到了 A 函数继续执行，则栈状态如下图：

![](../../../assets/images/docs/internal/goroutine/scheduler/function_call_stack/c98550de539d27e2.png)

因为 C、B 两个函数都已经执行完成并返回到了 A 函数之中，所以 C、B 两个函数的栈帧就已经被 POP 出栈了，也就是说它们所消耗的栈内存被自动回收了。因为现在正在执行 A 函数，所以寄存器 rbp 和 rsp 指向的是 A 函数的栈中的相应位置。如果 A 函数又继续调用了 D 函数的话，则栈又变成下面这个样子：

![](../../../assets/images/docs/internal/goroutine/scheduler/function_call_stack/11e8be87f54160dc.png)

可以看到，现在 D 函数的栈帧其实使用的是之前调用 B、C 两个函数所使用的栈内存，这没有问题，因为 B 和 C 函数已经执行完了，现在 D 函数重用了这块内存，这也是为什么**在 C 语言中绝对不要返回函数局部变量的地址，因为同一个地址的栈内存会被重用**，这就会造成意外的 bug，而 go 语言中没有这个限制，因为 **go 语言的编译器比较智能，当它发现程序返回了某个局部变量的地址，编译器会把这个变量放到堆上去，而不会放在栈上**。同样，这里我们还是需要注意 rbp 和 rsp 这两个寄存器现在指向了 D 函数的栈帧。从上面的分析我们可以看出，**寄存器 rbp 和 rsp 始终指向正在执行的函数的栈帧**。

最后，我们再来看一个递归函数的例子，假设有如下 go 语言代码片段：

```go
func f(n int) {
   if n <= 0 { //递归结束条件 n <= 0
       return
  }
   ......
   f(n - 1) //递归调用f函数自己
   ......
}
```

函数 f 是一个递归函数，f 函数会一直递归的调用自己直到参数 n 小于等于 0 为止，如果我们在其它某个函数里调用了 f(10)，而且现在正在执行 f(8) 的话，则其栈状态如下图所示：

![](../../../assets/images/docs/internal/goroutine/scheduler/function_call_stack/09ec60814ecdab98.png)

从上图可以看出，**即使是同一个函数，每次调用都会产生一个不同的栈帧**，因此对于递归函数，**每递归一次都会消耗一定的栈内存**，如果递归层数太多就有导致栈溢出的风险，这也是为什么我们**在实际的开发过程中应该尽量避免使用递归函数**的原因之一，另外一个原因是**递归函数执行效率比较低**，因为它要反复调用函数，而**调用函数有较大的性能开销**。

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

> 不同系统上的汇编结果可能不一样。

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

Go 语言函数的**参数和返回值存储在栈中**，然而许多主流的编程语言会将参数和返回值存储在寄存器中。存储在栈中的好处在于所有平台都可以使用相同的约定，从而容易开发出可移植、跨平台的代码，同时这种方式简化了协程、延迟调用和反射调用的实现。寄存器的值不能跨函数调用、存活，这简化了垃圾回收期间的栈扫描和对栈扩容的处理。

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

可以看出，堆栈信息是一种非常有用的排查问题的方法，同时，可以通过函数参数信息得知函数调用时传递的参数，学习 Go 语言内部类型的结构，以及值传递和指针传递的区别。

Go 语言可以通过配置 GOTRACEBACK 环境变量在程序异常终止时生成 core dump 文件，生成的文件可以由 dlv 或者 gdb 等高级的调试工具进行分析调试。

## 栈扩容

Go 语言在线程的基础上实现了用户态更加轻量级的协程，线程的栈大小一般是在创建时指定的，为了避免出现栈溢出(stack overflow)的错误，默认的栈大小会相对较大。

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
    ...
}
```

可以看到结构体 g 的第 1 个成员 stack 占 16 字节（lo 和 hi 各占 8 字节），所以结构体 g 变量的起始位置偏移 16 就对应到 stackguard0 字段。main 函数的第 2 条指令 `CMPQ SP, 16(CX)` 会比较栈顶寄存器 `rsp` 的值是否比 `stackguard0` 的值小，如果 `rsp` 的值更小，则说明当前 g 的栈要用完了，有溢出风险，需要调用 `morestack_noctxt` 函数来扩栈。

`stackguard0` 会在初始化时将 `stack.lo+_StackGuard`，`_StackGuard` 设置为 896 字节，`stack.lo` 为当前栈的栈顶。如果出现图 9-5 中栈寄存器 SP 小于 `stackguard0` 的情况，则表明当前栈空间不够，`stackguard0` 除了用于栈的扩容，还用于协程抢占。

![](../../../assets/images/docs/internal/function/stack/图9-5%20栈寄存器SP小于stackguard0.png)

在**函数序言阶段**如果判断出需要扩容，则会跳转调用运行时 `morestack_noctxt` 函数，函数调用链为 `morestack_noctxt()->morestack()->newstack()`，核心代码位于 `src/runtime/stack.go:newstack` 函数中，该函数不仅会处理扩容，还会处理协程的抢占。

```go
// Called from runtime·morestack when more stack is needed.
// Allocate larger stack and relocate to new stack.
// Stack growth is multiplicative, for constant amortized cost.
//
// g->atomicstatus will be Grunning or Gscanrunning upon entry.
// If the scheduler is trying to stop this g, then it will set preemptStop.
//
// This must be nowritebarrierrec because it can be called as part of
// stack growth from other nowritebarrierrec functions, but the
// compiler doesn't check this.
//
//go:nowritebarrierrec
func newstack() {
	...
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
	...
}
```

如上所示，newstack 函数首先通过栈底地址与栈顶地址相减计算旧栈的大小，并计算新栈的大小。新栈为旧栈的两倍大小。在 64 位操作系统中，如果栈大小超过了 1GB 则直接报错为 stack overflow。

## 栈转移

栈扩容的重要一步是将旧栈的内容转移到新栈中。栈扩容首先将协程的状态设置为 _Gcopystack，以便在垃圾回收状态下不会扫描该协程栈带来错误。栈复制并不像直接复制内存那样简单，如果栈中包含了引用栈中其他地址的指针，那么该指针需要对应到新栈中的地址，copystack 函数会分配一个新栈的内存。

为了应对频繁的栈调整，对获取栈的内存进行了许多优化，特别是对小栈。在 Linux 操作系统下，会对 2KB/4KB/8KB/16KB 的小栈进行专门的优化，即在全局及每个逻辑处理器(P)中预先分配这些小栈的缓存池，从而避免频繁地申请堆内存。

![](../../../assets/images/docs/internal/function/stack/图9-6%20栈的全局与本地缓存池结构.png)

栈的全局与本地缓存池结构如图 9-6 所示，每个逻辑处理器中的缓存池都来自全局缓存池(stackpool)。mcache 有时可能不存在（例如在调整 P 的大小后），这时需要直接从全局缓存池获取栈缓存。对于大栈，其大小不确定，虽然也有一个全局的缓存池，但不会预先放入多个栈，当栈被销毁时，如果被销毁的栈为大栈则放入全局缓存池中。**当全局缓存池中找不到对应大小的栈时，会从堆区分配。**

在分配到新栈后，如果有指针指向旧栈，那么需要将其调整到新栈中。在调整时有一个额外的步骤是调整 sudog，由于通道在阻塞的情况下存储的元素可能指向了栈上的指针，因此需要调整。接着需要将旧栈的大小复制到新栈中，这涉及借助 memmove 函数进行内存复制。

内存复制完成后，需要调整当前栈的 SP 寄存器和新的 stackguard0，并记录新的栈顶与栈底。扩容最关键的一步是在新栈中调整指针。因为新栈中的指针可能指向旧栈，旧栈一旦释放就会出现严重的问题。图 9-7 描述了栈扩容的过程，copystack 函数会遍历新栈上所有的栈帧信息，并遍历其中所有可能有指针的位置。一旦发现指针指向旧栈，就会调整当前的指针使其指向新栈。

![](../../../assets/images/docs/internal/function/stack/图9-7%20栈扩容的过程.png)

栈的转移如图 9-8 所示，调整后，栈指针将指向新栈中的地址。

![](../../../assets/images/docs/internal/function/stack/图9-8%20栈的转移.png)

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
