---
date: 2022-10-09T10:25:34+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "深入协程设计与调度原理"  # 文章标题
url:  "posts/go/docs/internal/goroutine/design/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 协程的生命周期与状态转移

Go 语言的调度器将协程分为多种状态，协程的状态与转移如图 15-1 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-1%20协程的状态与转移.png)

- `_Gidle` 为协程刚开始创建时的状态，当新创建的协程初始化后，会变为 `_Gdead` 状态，`_Gdead` 状态也是协程被销毁时的状态。
- `_Grunnable` 表示当前协程在运行队列中，正在等待运行。
- `_Grunning` 代表当前协程正在被运行，已经被分配给了逻辑处理器和线程。
- `_Gwaiting` 表示当前协程在运行时被锁定，不能执行用户代码。在垃圾回收及 channel 通信时经常会遇到这种情况。
- `_Gsyscall` 代表当前协程正在执行系统调用。
- `_Gpreempted` 代表协程 G 被强制抢占后的状态。
- `_Gcopystack` 代表在进行协程栈扫描时发现需要扩容或缩小协程栈空间，将协程中的栈转移到新栈时的状态。

还有几个状态（`_Gscan`、`_Gscanrunnable`、`_Gscanrunning` 等）涉及垃圾回收阶段。

`src/runtime/runtime2.go`

```go
// defined constants
const (
	// G status
	//
	// Beyond indicating the general state of a G, the G status
	// acts like a lock on the goroutine's stack (and hence its
	// ability to execute user code).
	//
	// If you add to this list, add to the list
	// of "okay during garbage collection" status
	// in mgcmark.go too.
	//
	// TODO(austin): The _Gscan bit could be much lighter-weight.
	// For example, we could choose not to run _Gscanrunnable
	// goroutines found in the run queue, rather than CAS-looping
	// until they become _Grunnable. And transitions like
	// _Gscanwaiting -> _Gscanrunnable are actually okay because
	// they don't affect stack ownership.

	// _Gidle means this goroutine was just allocated and has not
	// yet been initialized.
	_Gidle = iota // 0

	// _Grunnable means this goroutine is on a run queue. It is
	// not currently executing user code. The stack is not owned.
	_Grunnable // 1

	// _Grunning means this goroutine may execute user code. The
	// stack is owned by this goroutine. It is not on a run queue.
	// It is assigned an M and a P (g.m and g.m.p are valid).
	_Grunning // 2

	// _Gsyscall means this goroutine is executing a system call.
	// It is not executing user code. The stack is owned by this
	// goroutine. It is not on a run queue. It is assigned an M.
	_Gsyscall // 3

	// _Gwaiting means this goroutine is blocked in the runtime.
	// It is not executing user code. It is not on a run queue,
	// but should be recorded somewhere (e.g., a channel wait
	// queue) so it can be ready()d when necessary. The stack is
	// not owned *except* that a channel operation may read or
	// write parts of the stack under the appropriate channel
	// lock. Otherwise, it is not safe to access the stack after a
	// goroutine enters _Gwaiting (e.g., it may get moved).
	_Gwaiting // 4

	// _Gmoribund_unused is currently unused, but hardcoded in gdb
	// scripts.
	_Gmoribund_unused // 5

	// _Gdead means this goroutine is currently unused. It may be
	// just exited, on a free list, or just being initialized. It
	// is not executing user code. It may or may not have a stack
	// allocated. The G and its stack (if any) are owned by the M
	// that is exiting the G or that obtained the G from the free
	// list.
	_Gdead // 6

	// _Genqueue_unused is currently unused.
	_Genqueue_unused // 7

	// _Gcopystack means this goroutine's stack is being moved. It
	// is not executing user code and is not on a run queue. The
	// stack is owned by the goroutine that put it in _Gcopystack.
	_Gcopystack // 8

	// _Gpreempted means this goroutine stopped itself for a
	// suspendG preemption. It is like _Gwaiting, but nothing is
	// yet responsible for ready()ing it. Some suspendG must CAS
	// the status to _Gwaiting to take responsibility for
	// ready()ing this G.
	_Gpreempted // 9

	// _Gscan combined with one of the above states other than
	// _Grunning indicates that GC is scanning the stack. The
	// goroutine is not executing user code and the stack is owned
	// by the goroutine that set the _Gscan bit.
	//
	// _Gscanrunning is different: it is used to briefly block
	// state transitions while GC signals the G to scan its own
	// stack. This is otherwise like _Grunning.
	//
	// atomicstatus&~Gscan gives the state the goroutine will
	// return to when the scan completes.
	_Gscan          = 0x1000
	_Gscanrunnable  = _Gscan + _Grunnable  // 0x1001
	_Gscanrunning   = _Gscan + _Grunning   // 0x1002
	_Gscansyscall   = _Gscan + _Gsyscall   // 0x1003
	_Gscanwaiting   = _Gscan + _Gwaiting   // 0x1004
	_Gscanpreempted = _Gscan + _Gpreempted // 0x1009
)
```

## 特殊协程 g0 与协程切换

协程有 main 协程与子协程，main 协程在程序中只有一个。每个线程中都有一个特殊的协程 g0。

`src/runtime/runtime2.go`

```go
type m struct {
	g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink
	_       uint32 // align next field to 8 bytes

	// Fields not known to debuggers.
	procid        uint64            // for debuggers, but offset not hard-coded
	gsignal       *g                // signal-handling g
	goSigStack    gsignalStack      // Go-allocated signal handling stack
	sigmask       sigset            // storage for saved signal mask
	tls           [tlsSlots]uintptr // thread-local storage (for x86 extern register)
	mstartfn      func()
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a syscall
	id            int64
	mallocing     int32
	throwing      throwType
	preemptoff    string // if != "", keep curg running on this m
	locks         int32
	dying         int32
	profilehz     int32
	spinning      bool // m is out of work and is actively looking for work
	blocked       bool // m is blocked on a note
	newSigstack   bool // minit on C thread called sigaltstack
	printlock     int8
	incgo         bool   // m is executing a cgo call
	isextra       bool   // m is an extra m
	freeWait      uint32 // if == 0, safe to free g0 and delete m (atomic)
	fastrand      uint64
	needextram    bool
	traceback     uint8
	ncgocall      uint64        // number of cgo calls in total
	ncgo          int32         // number of cgo calls currently in progress
	cgoCallersUse atomic.Uint32 // if non-zero, cgoCallers in use temporarily
	cgoCallers    *cgoCallers   // cgo traceback if crashing in cgo call
	park          note
	alllink       *m // on allm
	schedlink     muintptr
	lockedg       guintptr
	createstack   [32]uintptr // stack that created this thread.
	lockedExt     uint32      // tracking for external LockOSThread
	lockedInt     uint32      // tracking for internal lockOSThread
	nextwaitm     muintptr    // next m waiting for lock
	waitunlockf   func(*g, unsafe.Pointer) bool
	waitlock      unsafe.Pointer
	waittraceev   byte
	waittraceskip int
	startingtrace bool
	syscalltick   uint32
	freelink      *m // on sched.freem

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall   libcall
	libcallpc uintptr // for cpu profiler
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall // stores syscall parameters on windows

	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
	vdsoPC uintptr // PC for traceback while in VDSO call

	// preemptGen counts the number of completed preemption
	// signals. This is used to detect when a preemption is
	// requested, but fails.
	preemptGen atomic.Uint32

	// Whether this is a pending preemption signal on this M.
	signalPending atomic.Uint32

	dlogPerM

	mOS

	// Up to 10 locks held by this m, maintained by the lock ranking code.
	locksHeldLen int
	locksHeld    [10]heldLockInfo
}
```

协程 g0 运行在操作系统线程栈上，其作用主要是**执行协程调度的一系列运行时代码**，而**普通协程无差别地用于执行用户代码**。很显然，执行用户代码的任何协程都不适合进行全局调度。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-2%20协程g与协程g0的对应关系.png)

在**用户协程退出或者被抢占时**，意味着需要重新执行协程调度，这时需要从用户协程 g 切换到协程 g0，协程 g 与协程 g0 的对应关系如图 15-2 所示。要注意的是，**每个线程的内部都在完成这样的切换与调度循环**。

协程经历 `g → g0 → g` 的过程，完成了一次调度循环。和线程类似，**协程切换的过程叫作协程的上下文切换**。当某一个协程 g 执行上下文切换时需要保存当前协程的执行现场，才能够在后续切换回 g 协程时正常执行。

协程的执行现场存储在 `gobuf` 结构体中，`gobuf` 结构体主要**保存 CPU 中几个重要的寄存器值**，分别是 `rsp`、`rip`、`rbp`。

`rsp` 寄存器始终指向函数调用栈栈顶，`rip` 寄存器指向程序要执行的下一条指令的地址，`rbp` 存储了函数栈帧的起始位置。

```go
type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it's
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn't need write barriers. It's still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr // 保存函数调用栈栈顶
	pc   uintptr // 保存程序要执行的下一条指令的地址
	g    guintptr // 保存当前协程的指针
	ctxt unsafe.Pointer // 保存当前协程的上下文
	ret  uintptr // 保存当前协程的返回值
	lr   uintptr // 保存当前协程的返回地址
	bp   uintptr // 保存函数栈帧的起始位置
}
```

特殊的协程 `g0` 与执行用户代码的协程 g 有显著不同，`g0` 作为特殊的调度协程，其执行的函数和流程相对固定，并且，为了避免栈溢出，**协程 `g0` 的栈会重复使用**。而每个执行用户代码的协程，可能都有不同的执行流程。每次上下文切换回去后，会继续执行之前的流程。

## 线程本地存储

线程本地存储是一种计算机编程方法，它使用线程本地的静态或全局内存。

**线程本地存储又叫线程局部存储，其英文为 Thread Local Storage，简称 TLS**，简而言之就是**线程私有的全局变量**。

和普通的全局变量对程序中的所有线程可见不同，**线程本地存储中的变量只对当前线程可见**。因此，**这种类型的变量可以看作是线程“私有”的**。

一般地，**操作系统使用 FS/GS 段寄存器存储线程本地变量**。

## 线程绑定

在 Go 语言中，并没有直接暴露线程本地存储的编程方式，但是 Go 语言运行时的调度器**使用线程本地存储将具体操作系统的线程与运行时代表线程的 m 结构体绑定在一起**。

如下所示，线程本地存储的实际是结构体 m 中 `m.tls` 的地址，同时 `m.tls[0]` 会存储当前线程正在运行的协程 g 的地址，因此**在任意一个线程内部，通过线程本地存储，都可以在任意时刻获取绑定到当前线程上的协程 g、结构体 m、逻辑处理器 P、特殊协程 g0 等信息。**

```go
	tls           [tlsSlots]uintptr // thread-local storage (for x86 extern register)
```

通过线程本地存储可以实现**结构体 m 与工作线程之间的绑定**，如图 15-3 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-3%20线程本地存储示意图.png)

## 调度循环

调度循环指从调度协程 g0 开始，找到接下来将要运行的协程 g、再从协程 g 切换到协程 g0 开始新一轮调度的过程。它和上下文切换类似，但是上下文切换关注的是具体切换的状态，而调度循环关注的是调度的流程。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-4%20调度循环.png)

图 15-4 所示为调度循环的整个流程。从协程 g0 调度到协程 g，**经历了从 schedule 函数到 execute 函数再到 gogo 函数的过程**。其中，

- **schedule** 函数**处理具体的调度策略**，选择下一个要执行的协程； 
- **execute** 函数执行一些具体的**状态转移、协程 g 与结构体 m 之间的绑定**等操作；
-  **gogo** 函数是与操作系统有关的函数，用于**完成栈的切换及 CPU 寄存器**的恢复。

执行完毕后，切换到协程 g 执行。当协程 g 主动让渡、被抢占或退出后，又会切换到协程 g0 进入第二轮调度。在从协程 g 切换回协程 g0 时，mcall 函数用于保存当前协程的执行现场，并切换到协程 g0 继续执行，mcall 函数仍然是和平台有关的汇编指令。切换到协程 g0 后会根据切换原因的不同执行不同的函数，例如，如果是用户调用 Gosched 函数则主动让渡执行权，执行 gosched_m 函数，如果协程已经退出，则执行 goexit 函数，将协程 g 放入 p 的 freeg 队列，方便下次重用。执行完毕后，再次调用 schedule 函数开始新一轮的调度循环，从而形成一个完整的闭环，循环往复。

## 调度策略

[调度策略](strategy.md)

## 调度时机

[调度时机](occasion.md)

```go

```
