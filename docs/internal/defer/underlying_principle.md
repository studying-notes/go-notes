---
date: 2022-10-05T14:47:30+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "defer 底层原理"  # 文章标题
url:  "posts/go/docs/internal/defer/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 演进

Go 语言中 defer 的实现经历了复杂的演进过程，Go 1.13、Go 1.14 都经历了比较大的更新。

在 Go 1.13 之前，defer 是被**分配在堆区**的，尽管有全局的缓存池分配，仍然有比较大的性能问题，原因在于使用 defer 不仅涉及堆内存的分配，在一开始还需要存储 defer 函数中的参数，最后还需要将堆区数据转移到栈中执行，涉及内存的复制。因此，defer 比普通函数的直接调用要慢很多。

为了将调用 defer 函数的成本降到与调用普通函数相同。Go 1.13 在大部分情况下将 defer 语句放置在了栈中，避免在堆区分配、复制对象。但是其仍然和 Go 1.12 一样，需要将整个 defer 语句放置到一条链表中，从而能够在函数退出时，以 LIFO 的顺序执行。

将 defer 添加到链表中被认为是必不可少的，原因在于 defer 的数量可能是无限的，也可能是动态调用的，例如通过 for 或者 if 块包裹的 defer 语句，只有在运行时才能决定执行的个数。

在 Go 1.13 中包括两种策略，对于最多调用一个(at most once)语义的 defer 语句使用了**栈分配**的策略，而对于其他的方式，例如 for 循环体内部的 defer 语句，仍然采用了之前的堆分配策略。

在大部分情况下，程序中的 defer 涉及的都是比较简单的场景，这一改变也大幅度提高了 defer 的效率。defer 的操作时间从 Go 1.12 时的 50ns 降到 Go 1.13 时 35ns（直接调用大约花费 6ns）。Go 1.13 虽然进行了一定程度的优化，但仍然比直接调用慢了 5、6 倍左右。

Go 1.14 进一步对最多调用一次的 defer 语义进行了优化，通过编译时实现**内联优化**。因此，在 Go 1.14 之后，根据不同的场景，实际存在了 3 种实现 defer 的方式。

`src/cmd/compile/internal/ssagen/ssa.go`

```go
	case ir.ODEFER:
		n := n.(*ir.GoDeferStmt)
		if base.Debug.Defer > 0 {
			var defertype string
			if s.hasOpenDefers {
				defertype = "open-coded"
			} else if n.Esc() == ir.EscNever {
				defertype = "stack-allocated"
			} else {
				defertype = "heap-allocated"
			}
			base.WarnfAt(n.Pos(), "%s defer", defertype)
		}
		if s.hasOpenDefers {
			s.openDeferRecord(n.Call.(*ir.CallExpr))
		} else {
			d := callDefer
			if n.Esc() == ir.EscNever {
				d = callDeferStack
			}
			s.callResult(n.Call.(*ir.CallExpr), d)
		}
```

## 堆分配

在 Go 1.13 前，defer 全部使用在堆区分配的内存存储。

目前在大部分情况下，堆分配只会在循环结构中出现，例如在 for 循环结构中。

```go
func main() {
    for i := 0; i < 100; i++ {
        defer fmt.Println(i)
    }
}
```

在上面的循环 defer 中，当执行汇编代码时，会发现每一条 defer 语句都调用了运行时的 runtime.deferproc 函数。

- deferproc()： 在声明 defer 处调用，其将 defer 函数存入 goroutine 的链表中；
- deferreturn()：在 return 指令，准确的讲是在 ret 指令前调用，其将 defer 从 goroutine 链表中取出并执行。

可以简单这么理解，在编译阶段，声明 defer 处插入了函数 deferproc()，在函数 return 前插入了函数 deferreturn()。

```
go tool compile -S -N -l main.go
```

```
	0x00a4 00164 (main.go:7)	CALL	runtime.deferproc(SB)
```

在函数退出前，调用了运行时 runtime.deferreturn 函数。

```
	0x00d4 00212 (main.go:9)	CALL	runtime.deferreturn(SB)
	0x00d9 00217 (main.go:9)	MOVQ	56(SP), BP
	0x00de 00222 (main.go:9)	ADDQ	$64, SP
	0x00e2 00226 (main.go:9)	RET
```

deferproc 函数的流程比较简单，主要分为 3 个步骤，如下所示：

- 计算 deferproc 调用者的 SP、PC 寄存器值及参数存放在栈中的位置。
- 在堆内存中分配新的 _defer 结构体，并将其插入当前协程记录 _defer 的链表头部。
- 将 SP、PC 寄存器值记录到新的 defer 结构体中，并将栈上的参数复制到堆区。

```go
// Create a new deferred function fn, which has no arguments and results.
// The compiler turns a defer statement into a call to this.
func deferproc(fn func()) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}

	d := newdefer()
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.link = gp._defer
	gp._defer = d
	d.fn = fn
	d.pc = getcallerpc()
	// We must not be preempted between calling getcallersp and
	// storing it to d.sp because getcallersp's result is a
	// uintptr stack pointer.
	d.sp = getcallersp()

	// deferproc returns 0 normally.
	// a deferred func that stops a panic
	// makes the deferproc return 1.
	// the code the compiler generates always
	// checks the return value and jumps to the
	// end of the function if deferproc returns != 0.
	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}
```

以下面的代码为例，defer add 需要传递两个参数。

```go
package main

func add(a, b int) int {
	return a + b
}

func f() {
	for i := 0; i < 2; i++ {
		defer add(3, 4)
	}
}

func main() {
	f()
}
```

当执行到 defer 语句时，调用运行时 deferproc 函数，其在栈中的结构如图 10-1 所示。

![](../../../assets/images/docs/internal/defer/underlying_principle/图10-1%20deferproc函数在栈中的结构.png)

每个协程都对应着一个结构体 g，deferproc 函数新建的 _defer 结构最终会被放置到当前协程存储 _defer 结构的链表中。

```go
type g struct {
	_defer    *_defer // innermost defer
}

// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in deferProcStack.
// This struct must match the code in cmd/compile/internal/ssagen/ssa.go:deferstruct
// and cmd/compile/internal/ssagen/ssa.go:(*state).call.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
	started bool
	heap    bool
	// openDefer indicates that this _defer is for a frame with open-coded
	// defers. We have only one defer record for the entire frame (which may
	// currently have 0, 1, or more defers active).
	openDefer bool
	sp        uintptr // sp at time of defer
	pc        uintptr // pc at time of defer
	fn        func()  // can be nil for open-coded defers
	_panic    *_panic // panic that is running defer
	link      *_defer // next defer on G; can point to either heap or stack!

	// If openDefer is true, the fields below record values about the stack
	// frame and associated function that has the open-coded defer(s). sp
	// above will be the sp for the frame, and pc will be address of the
	// deferreturn call in the function.
	fd   unsafe.Pointer // funcdata for the function associated with the frame
	varp uintptr        // value of varp for the stack frame
	// framepc is the current pc associated with the stack frame. Together,
	// with sp above (which is the sp associated with the stack frame),
	// framepc/sp can be used as pc/sp pair to continue a stack trace via
	// gentraceback().
	framepc uintptr
}
```

新加入的 _defer 结构会被放置到当前链表的头部，从而保证在后续执行 defer 函数时能以先入后出的顺序执行，如图 10-2 所示。

![](../../../assets/images/docs/internal/defer/underlying_principle/图10-2%20defer链先入后出的添加顺序.png)

runtime.newdefer 在堆中申请具体的 _defer 结构体，每个逻辑处理器 P 中都有局部缓存（deferpool），在全局中也有一个缓存池（schedt.deferpool），图 10-3 显示了 defer 全局与局部缓存池的交互。defer 根据结构的大小分为 5 个等级，以方便快速地找到最适合当前分配的 _defer 结构体。

![](../../../assets/images/docs/internal/defer/underlying_principle/图10-3%20defer全局与局部缓存池交互.png)

当在全局和局部缓存池中都搜索不到对象时，需要在堆区分配指定大小的 defer。

```go
// Allocate a Defer, usually using per-P pool.
// Each defer must be released with freedefer.  The defer is not
// added to any defer chain yet.
func newdefer() *_defer {
	var d *_defer
	mp := acquirem()
	pp := mp.p.ptr()
	if len(pp.deferpool) == 0 && sched.deferpool != nil {
		lock(&sched.deferlock)
		for len(pp.deferpool) < cap(pp.deferpool)/2 && sched.deferpool != nil {
			d := sched.deferpool
			sched.deferpool = d.link
			d.link = nil
			pp.deferpool = append(pp.deferpool, d)
		}
		unlock(&sched.deferlock)
	}
	if n := len(pp.deferpool); n > 0 {
		d = pp.deferpool[n-1]
		pp.deferpool[n-1] = nil
		pp.deferpool = pp.deferpool[:n-1]
	}
	releasem(mp)
	mp, pp = nil, nil

	if d == nil {
		// Allocate new defer.
		d = new(_defer)
	}
	d.heap = true
	return d
}
```

当 defer 执行完毕被销毁后，会重新回到局部缓存池中，当局部缓存池容纳了足够的对象时，会将 _defer 结构体放入全局缓存池。存储在全局和局部缓存池中的对象如果没有被使用，则最终在垃圾回收阶段被销毁。

## 遍历调用

当函数正常结束时，其递归调用了runtime.deferreturn 函数遍历 defer 链，并调用存储在 defer 中的函数。

```go
// deferreturn runs deferred functions for the caller's frame.
// The compiler inserts a call to this at the end of any
// function which calls defer.
func deferreturn() {
	gp := getg()
	for {
		d := gp._defer
		if d == nil {
			return
		}
		sp := getcallersp()
		if d.sp != sp {
			return
		}
		if d.openDefer {
			done := runOpenDeferFrame(d)
			if !done {
				throw("unfinished open-coded defers in deferreturn")
			}
			gp._defer = d.link
			freedefer(d)
			// If this frame uses open defers, then this
			// must be the only defer record for the
			// frame, so we can just return.
			return
		}

		fn := d.fn
		d.fn = nil
		gp._defer = d.link
		freedefer(d)
		fn()
	}
}
```

在遍历 defer 链表的过程中，有两个重要的终止条件。

一个是当遍历到链表的末尾时，最终链表指针变为 nil，这时需要终止链表。

除此之外，当 defer 结构中存储的 SP 地址与当前 deferreturn 的调用者 SP 地址不同时，仍然需要终止执行。原因是协程的链表中放入了当前函数调用链所有函数的 defer 结构，但是在执行时只能执行当前函数的 defer 结构。

例如，当前函数的执行链为 a()→ b()→ c()，在执行函数 c 正常返回后，当前三个函数的 defer 结构都存储在链表中，但是当前只能够执行函数 c 中的 fc 函数。如果发现 defer 结构是其他函数的内容，则立即返回。

```go
package main

func af() {
	println("af")
}

func bf() {
	println("bf")
}

func a() {
	defer af()
	b()
}

func b() {
	defer bf()
	c()
}

func c() {
	println("c")
}

func main() {
	a()
}
```

deferreturn 获取需要执行的 defer 函数后，需要将当前 defer 函数的参数重新转移到栈中，调用 freedefer 销毁当前的结构体，并将链表指向下一个 _defer 结构体。

![](../../../assets/images/docs/internal/defer/underlying_principle/图10-4%20递归defer调用的栈帧结构.png)

## 栈分配优化

从 defer 堆分配的过程可以看出，即便有全局和局部缓存池策略，由于涉及堆与栈参数的复制等操作，堆分配仍然比直接调用效率低下。

Go 1.13 为了解决堆分配的效率问题，对于最多调用一次的 defer 语义采用了在栈中分配的策略。

在 Go 1.14 中，对于非 for 循环的结构，有两种方式可以调试 defer 的栈分配。一种方式是禁止编译器优化（go tool compile-S -N -l main.go），另一种方式是增加 defer 的数量到 8 个以上，无论采用哪种方式，当执行到 defer 语句时，调用都会变为执行运行时的 runtime.deferprocStack 函数。在函数的最后，和堆分配一样，仍然插入了 runtime.deferreturn 函数用于遍历调用链。

deferprocStack 函数如下，其传递的参数为一个 _defer 指针，该 _defer 其实已经放置在了栈中。并在执行前将 defer 的大小、参数、函数指针放置在了栈中，在 deferprocStack 中只需要获取必要的调用者 SP、PC 指针并将 defer 压入链表的头部。

```go
func deferprocStack(d *_defer) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}
	// fn is already set.
	// The other fields are junk on entry to deferprocStack and
	// are initialized here.
	d.started = false
	d.heap = false
	d.openDefer = false
	d.sp = getcallersp()
	d.pc = getcallerpc()
	d.framepc = 0
	d.varp = 0
	// The lines below implement:
	//   d.panic = nil
	//   d.fd = nil
	//   d.link = gp._defer
	//   gp._defer = d
	// But without write barriers. The first three are writes to
	// the stack so they don't need a write barrier, and furthermore
	// are to uninitialized memory, so they must not use a write barrier.
	// The fourth write does not require a write barrier because we
	// explicitly mark all the defer structures, so we don't need to
	// keep track of pointers to them with a write barrier.
	*(*uintptr)(unsafe.Pointer(&d._panic)) = 0
	*(*uintptr)(unsafe.Pointer(&d.fd)) = 0
	*(*uintptr)(unsafe.Pointer(&d.link)) = uintptr(unsafe.Pointer(gp._defer))
	*(*uintptr)(unsafe.Pointer(&gp._defer)) = uintptr(unsafe.Pointer(d))

	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}
```

```go
package main

func add(a, b int) int {
	return a + b
}

func f() {
	for i := 0; i < 2; i++ {
		defer add(3, 4)
	}
}

func main() {
	f()
}
```

函数 f 的汇编代码如下，在调用 deferprocStack 前就已经把 defer 的大小、函数指针、参数都放置到了栈上对应的位置。

当函数执行完毕后，defer 函数和堆分配策略一样，需要遍历协程中的 _defer 链表，并递归调用 deferreturn 函数直到 defer 函数全部执行结束。从中可以看出，栈分配策略相较于之前的堆分配确实更高效。借助编译时的栈组织，不再需要在运行时将 _defer 结构体分配到堆中，既减少了分配内存的时间，也减少了在执行 defer 函数时将堆中参数复制到栈上的时间。

## 内联优化

虽然 Go 1.13 中 defer 的栈策略已经有了比较大的优化，但是与直接的函数调用还是有很大差别。一种容易想到的优化策略是在编译时函数结束时直接调用 defer 函数。这样就可以省去放置到 _defer 链表和遍历 _defer 链表的时间。

```go
func f() {
    defer a()
    defer b()
    c()
}
```

```go
def f():
    a()
    b()
    c()
```

采用这种方式最大的困难在于，defer 并不一定能够执行。例如，在 if 块中的 defer 语句，必须在运行时才能判断其是否成立。

```go
if a > b {
    defer c()
}
```

为了解决这样的问题，Go 语言编译器采取了一种巧妙的方式。通过在栈中初始化 1 字节的临时变量，以位图的形式来判断函数是否需要执行。

```go
defer f1()
if cond {
    defer f2()
}
```

如上代码将会被编译为如下所示的伪代码形式。

```go
deferBits |= 1
tmpF1 = f1
tmpA = a
if cond {
    deferBits |= 1<<1
    tmpF2 = f2
    tmpB = b
}

// 函数执行完毕后，执行 defer 函数
if deferBits & 1<<1 != 0 {
    deferBits &^= 1<<1
    tmpF2(tmpB)
}
if deferBits & 1<<0 != 0 {
    deferBits &^= 1<<0
    tmpF1(tmpA)
}
```

图 10-6 为标记是否需要调用的 deferBits 位图。上例中，由于 defer 函数 f1 一定会执行，因此把 deferBits 的最后 1 位设置为 1。而函数 f2 是否执行需要根据 cond 是否成立判断。如果成立，则需要将 deferBits 的倒数第 2 位设置为 1。

![](../../../assets/images/docs/internal/defer/underlying_principle/图10-6%20标记是否需要调用的deferBits位图.png)

在函数退出（exit）时，从后向前遍历 deferBits，如果当前位为 1，则需要执行对应的函数，如果当前位为 0，则不需要执行任何操作。另外，1 字节的 deferBits 位图以最小的代价满足了大部分情况下的需求。可以通过如下方式对加锁与解锁场景的直接调用与 defer 调用进行性能测试。

```go
package main

import (
	"sync"
	"testing"
)

func f1() {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()
}

func f2() {
	var m sync.Mutex
	m.Lock()
	m.Unlock()
}

func BenchmarkDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f1()
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f2()
	}
}
```

```
BenchmarkDefer-12       40495256                26.07 ns/op            8 B/op          1 allocs/op
BenchmarkDirect-12      54967706                23.18 ns/op            8 B/op          1 allocs/op

```

执行结果如下，从结果中可以看出，直接调用与 defer 调用的时间非常已经非常接近，二者被分配的内存大小也是相同的。

在实践中，如果使用了 Go 1.14 以上版本，则可以认为 defer 与直接调用的效率相当，不用为了考虑高性能而纠结是否使用 defer。

```go

```
