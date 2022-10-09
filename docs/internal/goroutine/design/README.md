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

之前介绍过，一般的协程有 main 协程与子协程，main 协程在整个程序中只有一个。每个线程中都有一个特殊的协程 g0。

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

协程 g0 运行在操作系统线程栈上，其作用主要是**执行协程调度的一系列运行时代码**，而**一般的协程无差别地用于执行用户代码**。很显然，执行用户代码的任何协程都不适合进行全局调度。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-2%20协程g与协程g0的对应关系.png)

在**用户协程退出或者被抢占时**，意味着需要重新执行协程调度，这时需要从用户协程 g 切换到协程 g0，协程 g 与协程 g0 的对应关系如图 15-2 所示。要注意的是，**每个线程的内部都在完成这样的切换与调度循环**。

协程经历 g → g0 → g 的过程，完成了一次调度循环。和线程类似，**协程切换的过程叫作协程的上下文切换**。当某一个协程 g 执行上下文切换时需要保存当前协程的执行现场，才能够在后续切换回 g 协程时正常执行。

协程的执行现场存储在 g.gobuf 结构体中，g.gobuf 结构体主要**保存 CPU 中几个重要的寄存器值**，分别是 rsp、rip、rbp。

rsp 寄存器始终指向函数调用栈栈顶，rip 寄存器指向程序要执行的下一条指令的地址，rbp 存储了函数栈帧的起始位置。

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

特殊的协程 g0 与执行用户代码的协程 g 有显著不同，g0 作为特殊的调度协程，其执行的函数和流程相对固定，并且，为了避免栈溢出，协程 g0 的栈会重复使用。而每个执行用户代码的协程，可能都有不同的执行流程。每次上下文切换回去后，会继续执行之前的流程。

## 线程本地存储与线程绑定

线程本地存储是一种计算机编程方法，它使用线程本地的静态或全局内存。

和普通的全局变量对程序中的所有线程可见不同，**线程本地存储中的变量只对当前线程可见**。因此，**这种类型的变量可以看作是线程“私有”的**。

一般地，**操作系统使用 FS/GS 段寄存器存储线程本地变量**。

在 Go 语言中，并没有直接暴露线程本地存储的编程方式，但是 Go 语言运行时的调度器**使用线程本地存储将具体操作系统的线程与运行时代表线程的 m 结构体绑定在一起**。

如下所示，线程本地存储的实际是结构体 m 中 `m.tls` 的地址，同时 `m.tls[0]` 会存储当前线程正在运行的协程 g 的地址，因此在任意一个线程内部，通过线程本地存储，都可以在任意时刻获取绑定到当前线程上的协程 g、结构体 m、逻辑处理器 P、特殊协程 g0 等信息。

```go
	tls           [tlsSlots]uintptr // thread-local storage (for x86 extern register)
```

通过线程本地存储可以实现**结构体 m 与工作线程之间的绑定**，如图 15-3 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-3%20线程本地存储示意图.png)

## 调度循环

调度循环指从调度协程 g0 开始，找到接下来将要运行的协程 g、再从协程 g 切换到协程 g0 开始新一轮调度的过程。它和上下文切换类似，但是上下文切换关注的是具体切换的状态，而调度循环关注的是调度的流程。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-4%20调度循环.png)

图 15-4 所示为调度循环的整个流程。从协程 g0 调度到协程 g，**经历了从 schedule 函数到 execute 函数再到 gogo 函数的过程**。

其中，**schedule 函数处理具体的调度策略，选择下一个要执行的协程**； **execute 函数执行一些具体的状态转移、协程 g 与结构体 m 之间的绑定等操作**； **gogo 函数是与操作系统有关的函数，用于完成栈的切换及 CPU 寄存器的恢复**。

执行完毕后，切换到协程 g 执行。当协程 g 主动让渡、被抢占或退出后，又会切换到协程 g0 进入第二轮调度。在从协程 g 切换回协程 g0 时，mcall 函数用于保存当前协程的执行现场，并切换到协程 g0 继续执行，mcall 函数仍然是和平台有关的汇编指令。切换到协程 g0 后会根据切换原因的不同执行不同的函数，例如，如果是用户调用 Gosched 函数则主动让渡执行权，执行 gosched_m 函数，如果协程已经退出，则执行 goexit 函数，将协程 g 放入 p 的 freeg 队列，方便下次重用。执行完毕后，再次调用 schedule 函数开始新一轮的调度循环，从而形成一个完整的闭环，循环往复。

## 调度策略

调度的核心策略位于 schedule 函数中。

`src/runtime/proc.go`

```go
// One round of scheduler: find a runnable goroutine and execute it.
// Never returns.
func schedule() {
	mp := getg().m

	if mp.locks != 0 {
		throw("schedule: holding locks")
	}

	if mp.lockedg != 0 {
		stoplockedm()
		execute(mp.lockedg.ptr(), false) // Never returns.
	}

	// We should not schedule away from a g that is executing a cgo call,
	// since the cgo call is using the m's g0 stack.
	if mp.incgo {
		throw("schedule: in cgo")
	}

top:
	pp := mp.p.ptr()
	pp.preempt = false

	// Safety check: if we are spinning, the run queue should be empty.
	// Check this before calling checkTimers, as that might call
	// goready to put a ready goroutine on the local run queue.
	if mp.spinning && (pp.runnext != 0 || pp.runqhead != pp.runqtail) {
		throw("schedule: spinning with local work")
	}

	gp, inheritTime, tryWakeP := findRunnable() // blocks until work is available

	// This thread is going to run a goroutine and is not spinning anymore,
	// so if it was marked as spinning we need to reset it now and potentially
	// start a new spinning M.
	if mp.spinning {
		resetspinning()
	}

	if sched.disable.user && !schedEnabled(gp) {
		// Scheduling of this goroutine is disabled. Put it on
		// the list of pending runnable goroutines for when we
		// re-enable user scheduling and look again.
		lock(&sched.lock)
		if schedEnabled(gp) {
			// Something re-enabled scheduling while we
			// were acquiring the lock.
			unlock(&sched.lock)
		} else {
			sched.disable.runnable.pushBack(gp)
			sched.disable.n++
			unlock(&sched.lock)
			goto top
		}
	}

	// If about to schedule a not-normal goroutine (a GCworker or tracereader),
	// wake a P if there is one.
	if tryWakeP {
		wakep()
	}
	if gp.lockedm != 0 {
		// Hands off own p to the locked m,
		// then blocks waiting for a new p.
		startlockedm(gp)
		goto top
	}

	execute(gp, inheritTime)
}
```

在 schedule 函数中，**首先会检测程序是否处于垃圾回收阶段**，如果是，则检测是否需要执行后台标记协程。

之前介绍过，程序中不可能同时执行成千上万个协程，那些等待被调度执行的协程存储在运行队列中。Go 语言调度器**将运行队列分为局部运行队列与全局运行队列**。局部运行队列是每个 P 特有的长度为 256 的数组，该数组模拟了一个循环队列，其中 runqhead 标识了循环队列的开头，runqtail 标识了循环队列的末尾。每次将 G 放入本地队列时，都从循环队列的末尾插入，而获取时从循环队列的头部获取。

除此之外，**在每个 P 内部还有一个特殊的 runnext 字段标识下一个要执行的协程。如果 runnext 不为空，则会直接执行当前 runnext 指向的协程，而不会去 runq 数组中寻找**。

```go
type p struct {
	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	pcache      pageCache
	raceprocctx uintptr

	deferpool    []*_defer // pool of available defer structs (see panic.go)
	deferpoolbuf [32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	//
	// Note that while other P's may atomically CAS this to zero,
	// only the owner P can CAS it to a valid G.
	runnext guintptr

	// Available G's (status == Gdead)
	gFree struct {
		gList
		n int32
	}

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	// Cache of mspan objects from the heap.
	mspancache struct {
		// We need an explicit length here because this field is used
		// in allocation codepaths where write barriers are not allowed,
		// and eliminating the write barrier/keeping it eliminated from
		// slice updates is tricky, moreso than just managing the length
		// ourselves.
		len int
		buf [128]*mspan
	}

	tracebuf traceBufPtr

	// traceSweep indicates the sweep events should be traced.
	// This is used to defer the sweep start event until a span
	// has actually been swept.
	traceSweep bool
	// traceSwept and traceReclaimed track the number of bytes
	// swept and reclaimed by sweeping in the current sweep loop.
	traceSwept, traceReclaimed uintptr

	palloc persistentAlloc // per-P to avoid mutex

	// The when field of the first entry on the timer heap.
	// This is 0 if the timer heap is empty.
	timer0When atomic.Int64

	// The earliest kfindrunnablenown nextwhen field of a timer with
	// timerModifiedEarlier status. Because the timer may have been
	// modified again, there need not be any timer with this value.
	// This is 0 if there are no timerModifiedEarlier timers.
	timerModifiedEarliest atomic.Int64

	// Per-P GC state
	gcAssistTime         int64 // Nanoseconds in assistAlloc
	gcFractionalMarkTime int64 // Nanoseconds in fractional mark worker (atomic)

	// limiterEvent tracks events for the GC CPU limiter.
	limiterEvent limiterEvent

	// gcMarkWorkerMode is the mode for the next mark worker to run in.
	// That is, this is used to communicate with the worker goroutine
	// selected for immediate execution by
	// gcController.findRunnableGCWorker. When scheduling other goroutines,
	// this field must be set to gcMarkWorkerNotWorker.
	gcMarkWorkerMode gcMarkWorkerMode
	// gcMarkWorkerStartTime is the nanotime() at which the most recent
	// mark worker started.
	gcMarkWorkerStartTime int64

	// gcw is this P's GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

	// wbBuf is this P's GC write barrier buffer.
	//
	// TODO: Consider caching this in the running G.
	wbBuf wbBuf

	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point

	// statsSeq is a counter indicating whether this P is currently
	// writing any stats. Its value is even when not, odd when it is.
	statsSeq atomic.Uint32

	// Lock for timers. We normally access the timers while running
	// on this P, but the scheduler can also do it from a different P.
	timersLock mutex

	// Actions to take at some time. This is used to implement the
	// standard library's time package.
	// Must hold timersLock to access.
	timers []*timer

	// Number of timers in P's heap.
	numTimers atomic.Uint32

	// Number of timerDeleted timers in P's heap.
	deletedTimers atomic.Uint32

	// Race context used while executing timer functions.
	timerRaceCtx uintptr

	// maxStackScanDelta accumulates the amount of stack space held by
	// live goroutines (i.e. those eligible for stack scanning).
	// Flushed to gcController.maxStackScan once maxStackScanSlack
	// or -maxStackScanSlack is reached.
	maxStackScanDelta int64

	// gc-time statistics about current goroutines
	// Note that this differs from maxStackScan in that this
	// accumulates the actual stack observed to be used at GC time (hi - sp),
	// not an instantaneous measure of the total stack size that might need
	// to be scanned (hi - lo).
	scannedStackSize uint64 // stack size of goroutines scanned by this P
	scannedStacks    uint64 // number of goroutines scanned by this P

	// preempt is set to indicate that this P should be enter the
	// scheduler ASAP (regardless of what G is running on it).
	preempt bool

	// Padding is no longer needed. False sharing is now not a worry because p is large enough
	// that its size class is an integer multiple of the cache line size (for any of our architectures).
}
```

被所有 P 共享的全局运行队列存储在 schedt.runq 中。

```go
type schedt struct {
	goidgen   atomic.Uint64
	lastpoll  atomic.Int64 // time of last network poll, 0 if currently polling
	pollUntil atomic.Int64 // time to which current poll is sleeping

	lock mutex

	// When increasing nmidle, nmidlelocked, nmsys, or nmfreed, be
	// sure to call checkdead().

	midle        muintptr // idle m's waiting for work
	nmidle       int32    // number of idle m's waiting for work
	nmidlelocked int32    // number of locked m's waiting for work
	mnext        int64    // number of m's that have been created and next M ID
	maxmcount    int32    // maximum number of m's allowed (or die)
	nmsys        int32    // number of system m's not counted for deadlock
	nmfreed      int64    // cumulative number of freed m's

	ngsys atomic.Int32 // number of system goroutines

	pidle        puintptr // idle p's
	npidle       atomic.Int32
	nmspinning   atomic.Int32  // See "Worker thread parking/unparking" comment in proc.go.
	needspinning atomic.Uint32 // See "Delicate dance" comment in proc.go. Boolean. Must hold sched.lock to set to 1.

	// Global runnable queue.
	runq     gQueue
	runqsize int32

	// disable controls selective disabling of the scheduler.
	//
	// Use schedEnableUser to control this.
	//
	// disable is protected by sched.lock.
	disable struct {
		// user disables scheduling of user goroutines.
		user     bool
		runnable gQueue // pending runnable Gs
		n        int32  // length of runnable
	}

	// Global cache of dead G's.
	gFree struct {
		lock    mutex
		stack   gList // Gs with stacks
		noStack gList // Gs without stacks
		n       int32
	}

	// Central cache of sudog structs.
	sudoglock  mutex
	sudogcache *sudog

	// Central pool of available defer structs.
	deferlock mutex
	deferpool *_defer

	// freem is the list of m's waiting to be freed when their
	// m.exited is set. Linked through m.freelink.
	freem *m

	gcwaiting  atomic.Bool // gc is waiting to run
	stopwait   int32
	stopnote   note
	sysmonwait atomic.Bool
	sysmonnote note

	// safepointFn should be called on each P at the next GC
	// safepoint if p.runSafePointFn is set.
	safePointFn   func(*p)
	safePointWait int32
	safePointNote note

	profilehz int32 // cpu profiling rate

	procresizetime int64 // nanotime() of last change to gomaxprocs
	totaltime      int64 // ∫gomaxprocs dt up to procresizetime

	// sysmonlock protects sysmon's actions on the runtime.
	//
	// Acquire and hold this mutex to block sysmon from interacting
	// with the rest of the runtime.
	sysmonlock mutex

	// timeToRun is a distribution of scheduling latencies, defined
	// as the sum of time a G spends in the _Grunnable state before
	// it transitions to _Grunning.
	timeToRun timeHistogram

	// idleTime is the total CPU time Ps have "spent" idle.
	//
	// Reset on each GC cycle.
	idleTime atomic.Int64

	// totalMutexWaitTime is the sum of time goroutines have spent in _Gwaiting
	// with a waitreason of the form waitReasonSync{RW,}Mutex{R,}Lock.
	totalMutexWaitTime atomic.Int64
}
```

```go

// A gQueue is a dequeue of Gs linked through g.schedlink. A G can only
// be on one gQueue or gList at a time.
type gQueue struct {
	head guintptr
	tail guintptr
}
```

因此，之前的 GMP 模型可以改进为图 15-5。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-5%20改进后的GMP模型.png)

一般的思路是先查找每个 P 局部的运行队列，当获取不到局部运行队列时，再从全局队列中获取。但是这种方法可能存在一个问题，如果只是循环往复地执行局部运行队列中的 G，那么全局队列中的 G 可能完全不会执行。为了避免这种情况，Go 语言调度器使用了一种策略：P 中每执行 61 次调度，就需要优先从全局队列中获取一个 G 到当前 P 中，并执行下一个要执行的 G。

```go
// Finds a runnable goroutine to execute.
// Tries to steal from other P's, get g from local or global queue, poll network.
// tryWakeP indicates that the returned goroutine is not normal (GC worker, trace
// reader) so the caller should try to wake a P.
func findRunnable() (gp *g, inheritTime, tryWakeP bool) {
	mp := getg().m

	// The conditions here and in handoffp must agree: if
	// findrunnable would return a G to run, handoffp must start
	// an M.

top:
	pp := mp.p.ptr()
	if sched.gcwaiting.Load() {
		gcstopm()
		goto top
	}
	if pp.runSafePointFn != 0 {
		runSafePointFn()
	}

	// now and pollUntil are saved for work stealing later,
	// which may steal timers. It's important that between now
	// and then, nothing blocks, so these numbers remain mostly
	// relevant.
	now, pollUntil, _ := checkTimers(pp, 0)

	// Try to schedule the trace reader.
	if trace.enabled || trace.shutdown {
		gp := traceReader()
		if gp != nil {
			casgstatus(gp, _Gwaiting, _Grunnable)
			traceGoUnpark(gp, 0)
			return gp, false, true
		}
	}

	// Try to schedule a GC worker.
	if gcBlackenEnabled != 0 {
		gp, tnow := gcController.findRunnableGCWorker(pp, now)
		if gp != nil {
			return gp, false, true
		}
		now = tnow
	}

	// Check the global runnable queue once in a while to ensure fairness.
	// Otherwise two goroutines can completely occupy the local runqueue
	// by constantly respawning each other.
	if pp.schedtick%61 == 0 && sched.runqsize > 0 {
		lock(&sched.lock)
		gp := globrunqget(pp, 1)
		unlock(&sched.lock)
		if gp != nil {
			return gp, false, false
		}
	}

	// Wake up the finalizer G.
	if fingStatus.Load()&(fingWait|fingWake) == fingWait|fingWake {
		if gp := wakefing(); gp != nil {
			ready(gp, 0, true)
		}
	}
	if *cgo_yield != nil {
		asmcgocall(*cgo_yield, nil)
	}

	// local runq
	if gp, inheritTime := runqget(pp); gp != nil {
		return gp, inheritTime, false
	}

	// global runq
	if sched.runqsize != 0 {
		lock(&sched.lock)
		gp := globrunqget(pp, 0)
		unlock(&sched.lock)
		if gp != nil {
			return gp, false, false
		}
	}

	// Poll network.
	// This netpoll is only an optimization before we resort to stealing.
	// We can safely skip it if there are no waiters or a thread is blocked
	// in netpoll already. If there is any kind of logical race with that
	// blocked thread (e.g. it has already returned from netpoll, but does
	// not set lastpoll yet), this thread will do blocking netpoll below
	// anyway.
	if netpollinited() && netpollWaiters.Load() > 0 && sched.lastpoll.Load() != 0 {
		if list := netpoll(0); !list.empty() { // non-blocking
			gp := list.pop()
			injectglist(&list)
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false, false
		}
	}

	// Spinning Ms: steal work from other Ps.
	//
	// Limit the number of spinning Ms to half the number of busy Ps.
	// This is necessary to prevent excessive CPU consumption when
	// GOMAXPROCS>>1 but the program parallelism is low.
	if mp.spinning || 2*sched.nmspinning.Load() < gomaxprocs-sched.npidle.Load() {
		if !mp.spinning {
			mp.becomeSpinning()
		}

		gp, inheritTime, tnow, w, newWork := stealWork(now)
		if gp != nil {
			// Successfully stole.
			return gp, inheritTime, false
		}
		if newWork {
			// There may be new timer or GC work; restart to
			// discover.
			goto top
		}

		now = tnow
		if w != 0 && (pollUntil == 0 || w < pollUntil) {
			// Earlier timer to wait for.
			pollUntil = w
		}
	}

	// We have nothing to do.
	//
	// If we're in the GC mark phase, can safely scan and blacken objects,
	// and have work to do, run idle-time marking rather than give up the P.
	if gcBlackenEnabled != 0 && gcMarkWorkAvailable(pp) && gcController.addIdleMarkWorker() {
		node := (*gcBgMarkWorkerNode)(gcBgMarkWorkerPool.pop())
		if node != nil {
			pp.gcMarkWorkerMode = gcMarkWorkerIdleMode
			gp := node.gp.ptr()
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false, false
		}
		gcController.removeIdleMarkWorker()
	}

	// wasm only:
	// If a callback returned and no other goroutine is awake,
	// then wake event handler goroutine which pauses execution
	// until a callback was triggered.
	gp, otherReady := beforeIdle(now, pollUntil)
	if gp != nil {
		casgstatus(gp, _Gwaiting, _Grunnable)
		if trace.enabled {
			traceGoUnpark(gp, 0)
		}
		return gp, false, false
	}
	if otherReady {
		goto top
	}

	// Before we drop our P, make a snapshot of the allp slice,
	// which can change underfoot once we no longer block
	// safe-points. We don't need to snapshot the contents because
	// everything up to cap(allp) is immutable.
	allpSnapshot := allp
	// Also snapshot masks. Value changes are OK, but we can't allow
	// len to change out from under us.
	idlepMaskSnapshot := idlepMask
	timerpMaskSnapshot := timerpMask

	// return P and block
	lock(&sched.lock)
	if sched.gcwaiting.Load() || pp.runSafePointFn != 0 {
		unlock(&sched.lock)
		goto top
	}
	if sched.runqsize != 0 {
		gp := globrunqget(pp, 0)
		unlock(&sched.lock)
		return gp, false, false
	}
	if !mp.spinning && sched.needspinning.Load() == 1 {
		// See "Delicate dance" comment below.
		mp.becomeSpinning()
		unlock(&sched.lock)
		goto top
	}
	if releasep() != pp {
		throw("findrunnable: wrong p")
	}
	now = pidleput(pp, now)
	unlock(&sched.lock)

	// Delicate dance: thread transitions from spinning to non-spinning
	// state, potentially concurrently with submission of new work. We must
	// drop nmspinning first and then check all sources again (with
	// #StoreLoad memory barrier in between). If we do it the other way
	// around, another thread can submit work after we've checked all
	// sources but before we drop nmspinning; as a result nobody will
	// unpark a thread to run the work.
	//
	// This applies to the following sources of work:
	//
	// * Goroutines added to a per-P run queue.
	// * New/modified-earlier timers on a per-P timer heap.
	// * Idle-priority GC work (barring golang.org/issue/19112).
	//
	// If we discover new work below, we need to restore m.spinning as a
	// signal for resetspinning to unpark a new worker thread (because
	// there can be more than one starving goroutine).
	//
	// However, if after discovering new work we also observe no idle Ps
	// (either here or in resetspinning), we have a problem. We may be
	// racing with a non-spinning M in the block above, having found no
	// work and preparing to release its P and park. Allowing that P to go
	// idle will result in loss of work conservation (idle P while there is
	// runnable work). This could result in complete deadlock in the
	// unlikely event that we discover new work (from netpoll) right as we
	// are racing with _all_ other Ps going idle.
	//
	// We use sched.needspinning to synchronize with non-spinning Ms going
	// idle. If needspinning is set when they are about to drop their P,
	// they abort the drop and instead become a new spinning M on our
	// behalf. If we are not racing and the system is truly fully loaded
	// then no spinning threads are required, and the next thread to
	// naturally become spinning will clear the flag.
	//
	// Also see "Worker thread parking/unparking" comment at the top of the
	// file.
	wasSpinning := mp.spinning
	if mp.spinning {
		mp.spinning = false
		if sched.nmspinning.Add(-1) < 0 {
			throw("findrunnable: negative nmspinning")
		}

		// Note the for correctness, only the last M transitioning from
		// spinning to non-spinning must perform these rechecks to
		// ensure no missed work. However, the runtime has some cases
		// of transient increments of nmspinning that are decremented
		// without going through this path, so we must be conservative
		// and perform the check on all spinning Ms.
		//
		// See https://go.dev/issue/43997.

		// Check all runqueues once again.
		pp := checkRunqsNoP(allpSnapshot, idlepMaskSnapshot)
		if pp != nil {
			acquirep(pp)
			mp.becomeSpinning()
			goto top
		}

		// Check for idle-priority GC work again.
		pp, gp := checkIdleGCNoP()
		if pp != nil {
			acquirep(pp)
			mp.becomeSpinning()

			// Run the idle worker.
			pp.gcMarkWorkerMode = gcMarkWorkerIdleMode
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false, false
		}

		// Finally, check for timer creation or expiry concurrently with
		// transitioning from spinning to non-spinning.
		//
		// Note that we cannot use checkTimers here because it calls
		// adjusttimers which may need to allocate memory, and that isn't
		// allowed when we don't have an active P.
		pollUntil = checkTimersNoP(allpSnapshot, timerpMaskSnapshot, pollUntil)
	}

	// Poll network until next timer.
	if netpollinited() && (netpollWaiters.Load() > 0 || pollUntil != 0) && sched.lastpoll.Swap(0) != 0 {
		sched.pollUntil.Store(pollUntil)
		if mp.p != 0 {
			throw("findrunnable: netpoll with p")
		}
		if mp.spinning {
			throw("findrunnable: netpoll with spinning")
		}
		// Refresh now.
		now = nanotime()
		delay := int64(-1)
		if pollUntil != 0 {
			delay = pollUntil - now
			if delay < 0 {
				delay = 0
			}
		}
		if faketime != 0 {
			// When using fake time, just poll.
			delay = 0
		}
		list := netpoll(delay) // block until new work is available
		sched.pollUntil.Store(0)
		sched.lastpoll.Store(now)
		if faketime != 0 && list.empty() {
			// Using fake time and nothing is ready; stop M.
			// When all M's stop, checkdead will call timejump.
			stopm()
			goto top
		}
		lock(&sched.lock)
		pp, _ := pidleget(now)
		unlock(&sched.lock)
		if pp == nil {
			injectglist(&list)
		} else {
			acquirep(pp)
			if !list.empty() {
				gp := list.pop()
				injectglist(&list)
				casgstatus(gp, _Gwaiting, _Grunnable)
				if trace.enabled {
					traceGoUnpark(gp, 0)
				}
				return gp, false, false
			}
			if wasSpinning {
				mp.becomeSpinning()
			}
			goto top
		}
	} else if pollUntil != 0 && netpollinited() {
		pollerPollUntil := sched.pollUntil.Load()
		if pollerPollUntil == 0 || pollerPollUntil > pollUntil {
			netpollBreak()
		}
	}
	stopm()
	goto top
}
```

调度协程的优先级与顺序如图 15-6 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-6%20调度协程的优先级与顺序.png)

排除从全局队列中获取这种情况，每个 P 在执行调度时，都会先尝试从 runnext 中获取下一个执行的 G，如果 runnext 为空，则继续从当前 P 中的局部运行队列 runq 中获取需要执行的 G ；如果局部运行队列为空，则尝试从全局运行队列中获取需要执行的 G ；如果全局队列也没有找到要执行的 G，则会尝试从其他的 P 中窃取可用的协程。到这一步，正常的程序基本都能获取到要运行的 G，如果窃取不到任务，那么当前的 P 会解除与 M 的绑定，P 会被放入空闲 P 队列中，而与 P 绑定的 M 没有任务可做，进入休眠状态。

### 获取本地运行队列

调度器首先查看 runnext 成员是否为空，如果不为空则返回对应的 G，如果为空则继续从局部运行队列中寻找。当循环队列的头（runqhead）和尾（runqtail）相同时，意味着循环队列中没有任何要运行的协程。否则，意味着存在可用的协程，从循环队列头部获取一个协程返回。需要注意的是，虽然在大部分情况下只有当前 G 访问局部运行队列，但是可能存在其他 P 窃取任务造成同时访问的情况，因此，在这里访问时需要加锁。

```go
// Get g from local runnable queue.
// If inheritTime is true, gp should inherit the remaining time in the
// current time slice. Otherwise, it should start a new time slice.
// Executed only by the owner P.
func runqget(pp *p) (gp *g, inheritTime bool) {
	// If there's a runnext, it's the next G to run.
	next := pp.runnext
	// If the runnext is non-0 and the CAS fails, it could only have been stolen by another P,
	// because other Ps can race to set runnext to 0, but only the current P can set it to non-0.
	// Hence, there's no need to retry this CAS if it fails.
	if next != 0 && pp.runnext.cas(next, 0) {
		return next.ptr(), true
	}

	for {
		h := atomic.LoadAcq(&pp.runqhead) // load-acquire, synchronize with other consumers
		t := pp.runqtail
		if t == h {
			return nil, false
		}
		gp := pp.runq[h%uint32(len(pp.runq))].ptr()
		if atomic.CasRel(&pp.runqhead, h, h+1) { // cas-release, commits consume
			return gp, false
		}
	}
}
```

### 获取全局运行队列

当 P 每执行 61 次调度，或者局部运行队列中不存在可用的协程时，都需要从全局运行队列中查找一批协程分配给本地运行队列，如图 15-7 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-7%20全局运行队列转移到本地运行队列.png)

```go
// Try get a batch of G's from the global runnable queue.
// sched.lock must be held.
func globrunqget(pp *p, max int32) *g {
	assertLockHeld(&sched.lock)

	if sched.runqsize == 0 {
		return nil
	}

	n := sched.runqsize/gomaxprocs + 1
	if n > sched.runqsize {
		n = sched.runqsize
	}
	if max > 0 && n > max {
		n = max
	}
	if n > int32(len(pp.runq))/2 {
		n = int32(len(pp.runq)) / 2
	}

	sched.runqsize -= n

	gp := sched.runq.pop()
	n--
	for ; n > 0; n-- {
		gp1 := sched.runq.pop()
		runqput(pp, gp1, false)
	}
	return gp
}
```

全局运行队列的数据结构是一根链表。由于每个 P 都共享了全局运行队列，因此为了保证公平，先根据 P 的数量平分全局运行队列中的 G，同时，要转移的数量不能超过局部队列容量的一半（当前是 256/2 = 128 个），再通过循环调用 runqput 将全局队列中的 G 放入 P 的局部运行队列中。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-8%20本地队列转移到全局队列.png)

如图 15-8 所示，如果本地运行队列满了，那么调度器会将本地运行队列的一半放入全局运行队列。这保证了当程序中有很多协程时，每个协程都有执行的机会。

### 获取准备就绪的网络协程

虽然很少见，但是局部运行队列和全局运行队列都找不到可用协程的情况仍有可能发生。这时，调度器会寻找当前是否有已经准备好运行的**网络协程**。

网络协程指的是在网络 I/O 中阻塞的协程，当网络 I/O 完成后，会将协程放入全局运行队列中，等待被调度。

Go 语言中的网络模型其实是对不同平台上 I/O 多路复用技术（epoll/kqueue/iocp）的封装。runtime.netpoll 函数获取当前可运行的协程列表，返回第一个可运行的协程。并通过 injectglist 函数将其余协程放入全局运行队列等待被调度。

`src/runtime/proc.go`

```go
	// Poll network.
	// This netpoll is only an optimization before we resort to stealing.
	// We can safely skip it if there are no waiters or a thread is blocked
	// in netpoll already. If there is any kind of logical race with that
	// blocked thread (e.g. it has already returned from netpoll, but does
	// not set lastpoll yet), this thread will do blocking netpoll below
	// anyway.
	if netpollinited() && netpollWaiters.Load() > 0 && sched.lastpoll.Load() != 0 {
		if list := netpoll(0); !list.empty() { // non-blocking
			gp := list.pop()
			injectglist(&list)
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			return gp, false, false
		}
	}
```

### 协程窃取

当局部运行队列、全局运行队列以及准备就绪的网络列表中都找不到可用协程时，需要从其他 P 的本地队列中窃取可用的协程执行。

所有的 P 都存储在全局的 `allp[]*p` 中，一种可以想到的简单方法是循环遍历 allp，找到可用的协程，但是这种方法缺少公平性。为了既保证随机性，又保证 allp 数组中的每个 P 都能被依次遍历，Go 语言采取了一种独特的方式，其代码位于 stealWork 函数中。

```go
// stealWork attempts to steal a runnable goroutine or timer from any P.
//
// If newWork is true, new work may have been readied.
//
// If now is not 0 it is the current time. stealWork returns the passed time or
// the current time if now was passed as 0.
func stealWork(now int64) (gp *g, inheritTime bool, rnow, pollUntil int64, newWork bool) {
	pp := getg().m.p.ptr()

	ranTimer := false

	const stealTries = 4
	for i := 0; i < stealTries; i++ {
		stealTimersOrRunNextG := i == stealTries-1

		for enum := stealOrder.start(fastrand()); !enum.done(); enum.next() {
			if sched.gcwaiting.Load() {
				// GC work may be available.
				return nil, false, now, pollUntil, true
			}
			p2 := allp[enum.position()]
			if pp == p2 {
				continue
			}

			// Steal timers from p2. This call to checkTimers is the only place
			// where we might hold a lock on a different P's timers. We do this
			// once on the last pass before checking runnext because stealing
			// from the other P's runnext should be the last resort, so if there
			// are timers to steal do that first.
			//
			// We only check timers on one of the stealing iterations because
			// the time stored in now doesn't change in this loop and checking
			// the timers for each P more than once with the same value of now
			// is probably a waste of time.
			//
			// timerpMask tells us whether the P may have timers at all. If it
			// can't, no need to check at all.
			if stealTimersOrRunNextG && timerpMask.read(enum.position()) {
				tnow, w, ran := checkTimers(p2, now)
				now = tnow
				if w != 0 && (pollUntil == 0 || w < pollUntil) {
					pollUntil = w
				}
				if ran {
					// Running the timers may have
					// made an arbitrary number of G's
					// ready and added them to this P's
					// local run queue. That invalidates
					// the assumption of runqsteal
					// that it always has room to add
					// stolen G's. So check now if there
					// is a local G to run.
					if gp, inheritTime := runqget(pp); gp != nil {
						return gp, inheritTime, now, pollUntil, ranTimer
					}
					ranTimer = true
				}
			}

			// Don't bother to attempt to steal if p2 is idle.
			if !idlepMask.read(enum.position()) {
				if gp := runqsteal(pp, p2, stealTimersOrRunNextG); gp != nil {
					return gp, false, now, pollUntil, ranTimer
				}
			}
		}
	}

	// No goroutines found to steal. Regardless, running a timer may have
	// made some goroutine ready that we missed. Indicate the next timer to
	// wait for.
	return nil, false, now, pollUntil, ranTimer
}
```

第 2 层 for 循环表示随机遍历 allp 数组，找到可窃取的 P 就立即窃取并返回。当遍历了一次没有找到时，再遍历一次，第 1 层的 4 个循环表示这个操作会重复四次，第 2 层的循环操作涉及数学上的一些特性。我们用一个例子来说明，假设一共有 8 个 P，第 1 步，fastrand 函数选择一个随机数，然后对 8 取模，算法选择了一个 0～8 之间的随机数，假设为 6。

第 2 步，找到一个比 8 小且与 8 互质的数。比 8 小且与 8 互质的数有 4 个：`coprimes = [1, 3, 5, 7]`，代码中取 `coprimes[6%4] = 5`，这 4 个数中任取一个都有相同的数学特性。

可以看到，这里将上一个计算的结果作为下一个计算的条件，这样的计算过程保证了一定会遍历到 allp 中的所有元素。

找到要窃取的 P 之后就正式开始窃取了，其核心代码位于 runqsteal 函数。窃取的核心逻辑比较简单，如图 15-9 所示，将要窃取的 P 本地运行队列中 Goroutine 个数的一半放入自己的运行队列中。

```go
// Steal half of elements from local runnable queue of p2
// and put onto local runnable queue of p.
// Returns one of the stolen elements (or nil if failed).
func runqsteal(pp, p2 *p, stealRunNextG bool) *g {
	t := pp.runqtail
	n := runqgrab(p2, &pp.runq, t, stealRunNextG)
	if n == 0 {
		return nil
	}
	n--
	gp := pp.runq[(t+n)%uint32(len(pp.runq))].ptr()
	if n == 0 {
		return gp
	}
	h := atomic.LoadAcq(&pp.runqhead) // load-acquire, synchronize with consumers
	if t-h+n >= uint32(len(pp.runq)) {
		throw("runqsteal: runq overflow")
	}
	atomic.StoreRel(&pp.runqtail, t+n) // store-release, makes the item available for consumption
	return gp
}
```

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-9%20窃取其他P中的协程.png)

## 调度时机

可以根据调度方式的不同，将调度时机分为主动、被动和抢占调度。

### 主动调度

协程可以选择主动让渡自己的执行权利，这主要是通过用户在代码中执行 runtime.Gosched 函数实现的。在大多数情况下，用户并不需要执行此函数，因为 Go 语言编译器会在调用函数之前插入检查代码，判断该协程是否需要被抢占。

但是有一些特殊的情况，例如一个密集计算，无限 for 循环的场景，这种场景由于没有抢占的时机，在 Go 1.14 版本之前是无法被抢占的。

Go 1.14 之后的版本对于长时间执行的协程使用了操作系统的信号机制进行强制抢占。这种方式需要进入操作系统的内核，速度比不上用户直接调度的 runtime.Gosched 函数。

主动调度的原理比较简单，需要先**从当前协程切换到协程 g0**，**取消 G 与 M 之间的绑定关系**，**将 G 放入全局运行队列**，并**调用 schedule 函数开始新一轮的循环**。

`src/runtime/proc.go`

```go
func goschedImpl(gp *g) {    
	status := readgstatus(gp) // 获取当前协程的状态
	if status&^_Gscan != _Grunning { // 如果当前协程不是运行状态，抛出异常
		dumpgstatus(gp) // 打印当前协程的状态
		throw("bad g status")
	}
    // 将当前协程的状态设置为 _Grunnable
	casgstatus(gp, _Grunning, _Grunnable)
    // 取消 G 与 M 之间的绑定关系
	dropg() 
	lock(&sched.lock)
    // 将 G 放入全局运行队列
	globrunqput(gp)
	unlock(&sched.lock)

    // 调用 schedule 函数开始新一轮的循环
	schedule()
}
```

### 被动调度

被动调度指协程在休眠、channel 通道堵塞、网络 I/O 堵塞、执行垃圾回收而暂停时，被动让渡自己执行权利的过程。

被动调度具有重要的意义，可以保证最大化利用 CPU 的资源。根据被动调度的原因不同，调度器可能执行一些特殊的操作。由于被动调度仍然是协程发起的操作，因此其调度的时机相对明确。

和主动调度类似的是，被动调度需要先从当前协程切换到协程 g0，更新协程的状态并解绑与 M 的关系，重新调度。

和主动调度不同的是，被动调度不会将 G 放入全局运行队列，因为当前 G 的状态不是 _Grunnable 而是 _Gwaiting，所以，被动调度需要一个额外的唤醒机制。

下面以通道的堵塞为例说明被动调度的过程。在该例中，通道 c 一直会等待通道中的消息。

```go
func main() {
    c := make(chan int)
    go func() {
        for {
            <-c
        }
    }()
    for {
    }
}
```

当通道中暂时没有数据时，会调用 gopark 函数完成被动调度，gopark 函数是被动调度的核心逻辑。

`src/runtime/proc.go`

```go
// Puts the current goroutine into a waiting state and calls unlockf on the
// system stack.
//
// If unlockf returns false, the goroutine is resumed.
//
// unlockf must not access this G's stack, as it may be moved between
// the call to gopark and the call to unlockf.
//
// Note that because unlockf is called after putting the G into a waiting
// state, the G may have already been readied by the time unlockf is called
// unless there is external synchronization preventing the G from being
// readied. If unlockf returns false, it must guarantee that the G cannot be
// externally readied.
//
// Reason explains why the goroutine has been parked. It is displayed in stack
// traces and heap dumps. Reasons should be unique and descriptive. Do not
// re-use reasons, add new ones.
func gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceEv byte, traceskip int) {
	if reason != waitReasonSleep {
		checkTimeouts() // timeouts may expire while two goroutines keep the scheduler busy
	}
	mp := acquirem()
	gp := mp.curg
	status := readgstatus(gp)
	if status != _Grunning && status != _Gscanrunning {
		throw("gopark: bad g status")
	}
	mp.waitlock = lock
	mp.waitunlockf = unlockf
	gp.waitreason = reason
	mp.waittraceev = traceEv
	mp.waittraceskip = traceskip
	releasem(mp)
	// can't do anything that might move the G between Ms here.
	mcall(park_m)
}
```

gopark 函数最后会调用 park_m，该函数会解除 G 和 M 之间的关系，根据执行被动调度的原因不同，执行不同的 waitunlockf 函数，并开始新一轮调度。

```go
// park continuation on g0.
func park_m(gp *g) {
	mp := getg().m

	if trace.enabled {
		traceGoPark(mp.waittraceev, mp.waittraceskip)
	}

	// N.B. Not using casGToWaiting here because the waitreason is
	// set by park_m's caller.
	casgstatus(gp, _Grunning, _Gwaiting)
	dropg()

	if fn := mp.waitunlockf; fn != nil {
		ok := fn(gp, mp.waitlock)
		mp.waitunlockf = nil
		mp.waitlock = nil
		if !ok {
			if trace.enabled {
				traceGoUnpark(gp, 2)
			}
			casgstatus(gp, _Gwaiting, _Grunnable)
			execute(gp, true) // Schedule it back, never returns.
		}
	}
	schedule()
}
```

在上面的例子中，通道 c 会一直等待通道中的消息，当通道中有消息时，会调用 goready 函数唤醒被调度的协程。

`src/runtime/chan.go`

```go
// send processes a send operation on an empty channel c.
// The value ep sent by the sender is copied to the receiver sg.
// The receiver is then woken up to go on its merry way.
// Channel c must be empty and locked.  send unlocks c with unlockf.
// sg must already be dequeued from c.
// ep must be non-nil and point to the heap or the caller's stack.
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if raceenabled {
		if c.dataqsiz == 0 {
			racesync(c, sg)
		} else {
			// Pretend we go through the buffer, even though
			// we copy directly. Note that we need to increment
			// the head/tail locations only when raceenabled.
			racenotify(c, c.recvx, nil)
			racenotify(c, c.recvx, sg)
			c.recvx++
			if c.recvx == c.dataqsiz {
				c.recvx = 0
			}
			c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
		}
	}
	if sg.elem != nil {
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	goready(gp, skip+1)
}

func goready(gp *g, traceskip int) {
	systemstack(func() {
		ready(gp, traceskip, true)
	})
}
```

如果当前协程需要被唤醒，那么会先将协程的状态从 _Gwaiting 转换为 _Grunnable，并添加到当前 P 的局部运行队列中。

```go
// Mark gp ready to run.
func ready(gp *g, traceskip int, next bool) {
	if trace.enabled {
		traceGoUnpark(gp, traceskip)
	}

	status := readgstatus(gp)

	// Mark runnable.
	mp := acquirem() // disable preemption because it can be holding p in a local var
	if status&^_Gscan != _Gwaiting {
		dumpgstatus(gp)
		throw("bad g->status in ready")
	}

	// status is Gwaiting or Gscanwaiting, make Grunnable and put on runq
	casgstatus(gp, _Gwaiting, _Grunnable)
	runqput(mp.p.ptr(), gp, next)
	wakep()
	releasem(mp)
}
```

### 抢占调度

为了让每个协程都有执行的机会，并且最大化利用 CPU 资源，Go 语言在初始化时会启动一个特殊的线程来执行系统监控任务。

系统监控在一个独立的 M 上运行，不用绑定逻辑处理器 P，系统监控每隔 10ms 会检测是否有准备就绪的网络协程，并放置到全局队列中。和抢占调度相关的是，系统监控服务会判断当前协程是否运行时间过长，或者处于系统调用阶段，如果是，则会抢占当前 G 的执行。

其核心逻辑位于 runtime.retake 函数中。

```go
func retake(now int64) uint32 {
	n := 0
	// Prevent allp slice changes. This lock will be completely
	// uncontended unless we're already stopping the world.
	lock(&allpLock)
	// We can't use a range loop over allp because we may
	// temporarily drop the allpLock. Hence, we need to re-fetch
	// allp each time around the loop.
	for i := 0; i < len(allp); i++ {
		pp := allp[i]
		if pp == nil {
			// This can happen if procresize has grown
			// allp but not yet created new Ps.
			continue
		}
		pd := &pp.sysmontick
		s := pp.status
		sysretake := false
		if s == _Prunning || s == _Psyscall {
			// Preempt G if it's running for too long.
			t := int64(pp.schedtick)
			if int64(pd.schedtick) != t {
				pd.schedtick = uint32(t)
				pd.schedwhen = now
			} else if pd.schedwhen+forcePreemptNS <= now {
				preemptone(pp)
				// In case of syscall, preemptone() doesn't
				// work, because there is no M wired to P.
				sysretake = true
			}
		}
		if s == _Psyscall {
			// Retake P from syscall if it's there for more than 1 sysmon tick (at least 20us).
			t := int64(pp.syscalltick)
			if !sysretake && int64(pd.syscalltick) != t {
				pd.syscalltick = uint32(t)
				pd.syscallwhen = now
				continue
			}
			// On the one hand we don't want to retake Ps if there is no other work to do,
			// but on the other hand we want to retake them eventually
			// because they can prevent the sysmon thread from deep sleep.
			if runqempty(pp) && sched.nmspinning.Load()+sched.npidle.Load() > 0 && pd.syscallwhen+10*1000*1000 > now {
				continue
			}
			// Drop allpLock so we can take sched.lock.
			unlock(&allpLock)
			// Need to decrement number of idle locked M's
			// (pretending that one more is running) before the CAS.
			// Otherwise the M from which we retake can exit the syscall,
			// increment nmidle and report deadlock.
			incidlelocked(-1)
			if atomic.Cas(&pp.status, s, _Pidle) {
				if trace.enabled {
					traceGoSysBlock(pp)
					traceProcStop(pp)
				}
				n++
				pp.syscalltick++
				handoffp(pp)
			}
			incidlelocked(1)
			lock(&allpLock)
		}
	}
	unlock(&allpLock)
	return uint32(n)
}
```

在 Go1.14 中，如果当前协程的执行时间超过了 10ms，则需要执行抢占。如果一个协程在系统调用中超过了 20 微秒，则仍然需要抢占调度。接下来，我们分别分析这两种不同的情况。

#### 执行时间过长的抢占调度

在 Go 1.14 之前，虽然仍然有系统监控抢占时间过长的 G，调用 preemptone 函数，但是抢占的时机却不太一样。preemptone 函数会将当前的 preempt 字段设置为 true，并将 stackguard0 设置为 stackPreempt。stackPreempt 常量 0xfffffffffffffade 是一个非常大的数，设置 stackguard0 使调度器能够处理抢占请求。

```go
// Tell the goroutine running on processor P to stop.
// This function is purely best-effort. It can incorrectly fail to inform the
// goroutine. It can inform the wrong goroutine. Even if it informs the
// correct goroutine, that goroutine might ignore the request if it is
// simultaneously executing newstack.
// No lock needs to be held.
// Returns true if preemption request was issued.
// The actual preemption will happen at some point in the future
// and will be indicated by the gp->status no longer being
// Grunning
func preemptone(pp *p) bool {
	mp := pp.m.ptr()
	if mp == nil || mp == getg().m {
		return false
	}
	gp := mp.curg
	if gp == nil || gp == mp.g0 {
		return false
	}

	gp.preempt = true

	// Every call in a goroutine checks for stack overflow by
	// comparing the current stack pointer to gp->stackguard0.
	// Setting gp->stackguard0 to StackPreempt folds
	// preemption into the normal stack overflow check.
	gp.stackguard0 = stackPreempt

	// Request an async preemption of this P.
	if preemptMSupported && debug.asyncpreemptoff == 0 {
		pp.preempt = true
		preemptM(mp)
	}

	return true
}
```

调度发生的时机主要在执行函数调用阶段。函数调用是一个比较安全的检查点，Go 语言编译器会在函数调用前判断 stackguard0 的大小，从而选择是否调用 runtime.morestack_noctxt 函数。morestack_noctxt 为汇编函数，函数的执行流程如下：morestack_noctxt()→ morestack()→ newstack()。

newstack 函数中的一般核心逻辑是判断 G 中 stackguard0 字段的大小，并调用 gopreempt_m 函数切换到 g0，取消 G 与 M 之间的绑定关系，将 G 的状态转换为 _Grunnable，将 G 放入全局运行队列，并调用 schedule 函数开始新一轮调度循环。

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
	thisg := getg()
	// TODO: double check all gp. shouldn't be getg().
	if thisg.m.morebuf.g.ptr().stackguard0 == stackFork {
		throw("stack growth after fork")
	}
	if thisg.m.morebuf.g.ptr() != thisg.m.curg {
		print("runtime: newstack called from g=", hex(thisg.m.morebuf.g), "\n"+"\tm=", thisg.m, " m->curg=", thisg.m.curg, " m->g0=", thisg.m.g0, " m->gsignal=", thisg.m.gsignal, "\n")
		morebuf := thisg.m.morebuf
		traceback(morebuf.pc, morebuf.sp, morebuf.lr, morebuf.g.ptr())
		throw("runtime: wrong goroutine in newstack")
	}

	gp := thisg.m.curg

	if thisg.m.curg.throwsplit {
		// Update syscallsp, syscallpc in case traceback uses them.
		morebuf := thisg.m.morebuf
		gp.syscallsp = morebuf.sp
		gp.syscallpc = morebuf.pc
		pcname, pcoff := "(unknown)", uintptr(0)
		f := findfunc(gp.sched.pc)
		if f.valid() {
			pcname = funcname(f)
			pcoff = gp.sched.pc - f.entry()
		}
		print("runtime: newstack at ", pcname, "+", hex(pcoff),
			" sp=", hex(gp.sched.sp), " stack=[", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n",
			"\tmorebuf={pc:", hex(morebuf.pc), " sp:", hex(morebuf.sp), " lr:", hex(morebuf.lr), "}\n",
			"\tsched={pc:", hex(gp.sched.pc), " sp:", hex(gp.sched.sp), " lr:", hex(gp.sched.lr), " ctxt:", gp.sched.ctxt, "}\n")

		thisg.m.traceback = 2 // Include runtime frames
		traceback(morebuf.pc, morebuf.sp, morebuf.lr, gp)
		throw("runtime: stack split at bad time")
	}

	morebuf := thisg.m.morebuf
	thisg.m.morebuf.pc = 0
	thisg.m.morebuf.lr = 0
	thisg.m.morebuf.sp = 0
	thisg.m.morebuf.g = 0

	// NOTE: stackguard0 may change underfoot, if another thread
	// is about to try to preempt gp. Read it just once and use that same
	// value now and below.
	stackguard0 := atomic.Loaduintptr(&gp.stackguard0)

	// Be conservative about where we preempt.
	// We are interested in preempting user Go code, not runtime code.
	// If we're holding locks, mallocing, or preemption is disabled, don't
	// preempt.
	// This check is very early in newstack so that even the status change
	// from Grunning to Gwaiting and back doesn't happen in this case.
	// That status change by itself can be viewed as a small preemption,
	// because the GC might change Gwaiting to Gscanwaiting, and then
	// this goroutine has to wait for the GC to finish before continuing.
	// If the GC is in some way dependent on this goroutine (for example,
	// it needs a lock held by the goroutine), that small preemption turns
	// into a real deadlock.
	preempt := stackguard0 == stackPreempt
	if preempt {
		if !canPreemptM(thisg.m) {
			// Let the goroutine keep running for now.
			// gp->preempt is set, so it will be preempted next time.
			gp.stackguard0 = gp.stack.lo + _StackGuard
			gogo(&gp.sched) // never return
		}
	}

	if gp.stack.lo == 0 {
		throw("missing stack in newstack")
	}
	sp := gp.sched.sp
	if goarch.ArchFamily == goarch.AMD64 || goarch.ArchFamily == goarch.I386 || goarch.ArchFamily == goarch.WASM {
		// The call to morestack cost a word.
		sp -= goarch.PtrSize
	}
	if stackDebug >= 1 || sp < gp.stack.lo {
		print("runtime: newstack sp=", hex(sp), " stack=[", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n",
			"\tmorebuf={pc:", hex(morebuf.pc), " sp:", hex(morebuf.sp), " lr:", hex(morebuf.lr), "}\n",
			"\tsched={pc:", hex(gp.sched.pc), " sp:", hex(gp.sched.sp), " lr:", hex(gp.sched.lr), " ctxt:", gp.sched.ctxt, "}\n")
	}
	if sp < gp.stack.lo {
		print("runtime: gp=", gp, ", goid=", gp.goid, ", gp->status=", hex(readgstatus(gp)), "\n ")
		print("runtime: split stack overflow: ", hex(sp), " < ", hex(gp.stack.lo), "\n")
		throw("runtime: split stack overflow")
	}

	if preempt {
		if gp == thisg.m.g0 {
			throw("runtime: preempt g0")
		}
		if thisg.m.p == 0 && thisg.m.locks == 0 {
			throw("runtime: g is running but p is not")
		}

		if gp.preemptShrink {
			// We're at a synchronous safe point now, so
			// do the pending stack shrink.
			gp.preemptShrink = false
			shrinkstack(gp)
		}

		if gp.preemptStop {
			preemptPark(gp) // never returns
		}

		// Act like goroutine called runtime.Gosched.
		gopreempt_m(gp) // never return
	}

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
}
```

这种抢占的方式面临着一定的问题，当执行过程中没有函数调用，而只有类似如下代码时，协程将没有被抢占的机会。

```go
for {
   i++
}
```

为了解决这一问题，Go 1.14 之后引入了信号强制抢占的机制。

这需要借助图 15-10 中的类 UNIX 操作系统信号处理机制，信号是发送给进程的各种通知，以便将各种重要的事件通知给进程。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-10%20类UNIX操作系统信号处理机制.png)

最常见的是用户发送给进程的信号，例如时常使用的 CTRL+C 键，或者在命令行中输入的 `kill-<signal><PID>` 指令。通过信号，借助操作系统中断当前程序，保存程序的执行状态和寄存器值，并切换到内核态处理信号。在内核态处理完信号后，还会返回到用户态执行程序注册的信号处理函数，之后再回到内核，恢复程序原始的栈和寄存器值，并切换到用户态继续执行程序。

Go 语言借助用户态在信号处理时完成协程的上下文切换的操作，需要借助进程对特定的信号进行处理。并不是所有的信号都可以被处理，例如 SIGKILL 与 SIGSTOP 信号用于终止或暂停程序，不能被程序捕获处理。Go 程序在初始化时会初始化信号表，并注册信号处理函数。在这里，我们关注抢占时的信号处理。

在抢占时，调度器通过向线程中发送 sigPreempt 信号，触发信号处理。在 UNIX 操作系统中，sigPreempt 为 _SIGURG 信号，由于该信号不会被用户程序和调试器使用，因此 Go 语言使用它作为安全的抢占信号，关于信号具体的选择过程，可以参考 Go 源码中对 _SIGURG 信号的注释。

```go
func preemptM(mp *m) {
	if mp == getg().m {
		throw("self-preempt")
	}

	// Synchronize with external code that may try to ExitProcess.
	if !atomic.Cas(&mp.preemptExtLock, 0, 1) {
		// External code is running. Fail the preemption
		// attempt.
		mp.preemptGen.Add(1)
		return
	}

	// Acquire our own handle to mp's thread.
	lock(&mp.threadLock)
	if mp.thread == 0 {
		// The M hasn't been minit'd yet (or was just unminit'd).
		unlock(&mp.threadLock)
		atomic.Store(&mp.preemptExtLock, 0)
		mp.preemptGen.Add(1)
		return
	}
	var thread uintptr
	if stdcall7(_DuplicateHandle, currentProcess, mp.thread, currentProcess, uintptr(unsafe.Pointer(&thread)), 0, 0, _DUPLICATE_SAME_ACCESS) == 0 {
		print("runtime.preemptM: duplicatehandle failed; errno=", getlasterror(), "\n")
		throw("runtime.preemptM: duplicatehandle failed")
	}
	unlock(&mp.threadLock)

	// Prepare thread context buffer. This must be aligned to 16 bytes.
	var c *context
	var cbuf [unsafe.Sizeof(*c) + 15]byte
	c = (*context)(unsafe.Pointer((uintptr(unsafe.Pointer(&cbuf[15]))) &^ 15))
	c.contextflags = _CONTEXT_CONTROL

	// Serialize thread suspension. SuspendThread is asynchronous,
	// so it's otherwise possible for two threads to suspend each
	// other and deadlock. We must hold this lock until after
	// GetThreadContext, since that blocks until the thread is
	// actually suspended.
	lock(&suspendLock)

	// Suspend the thread.
	if int32(stdcall1(_SuspendThread, thread)) == -1 {
		unlock(&suspendLock)
		stdcall1(_CloseHandle, thread)
		atomic.Store(&mp.preemptExtLock, 0)
		// The thread no longer exists. This shouldn't be
		// possible, but just acknowledge the request.
		mp.preemptGen.Add(1)
		return
	}

	// We have to be very careful between this point and once
	// we've shown mp is at an async safe-point. This is like a
	// signal handler in the sense that mp could have been doing
	// anything when we stopped it, including holding arbitrary
	// locks.

	// We have to get the thread context before inspecting the M
	// because SuspendThread only requests a suspend.
	// GetThreadContext actually blocks until it's suspended.
	stdcall2(_GetThreadContext, thread, uintptr(unsafe.Pointer(c)))

	unlock(&suspendLock)

	// Does it want a preemption and is it safe to preempt?
	gp := gFromSP(mp, c.sp())
	if gp != nil && wantAsyncPreempt(gp) {
		if ok, newpc := isAsyncSafePoint(gp, c.ip(), c.sp(), c.lr()); ok {
			// Inject call to asyncPreempt
			targetPC := abi.FuncPCABI0(asyncPreempt)
			switch GOARCH {
			default:
				throw("unsupported architecture")
			case "386", "amd64":
				// Make it look like the thread called targetPC.
				sp := c.sp()
				sp -= goarch.PtrSize
				*(*uintptr)(unsafe.Pointer(sp)) = newpc
				c.set_sp(sp)
				c.set_ip(targetPC)

			case "arm":
				// Push LR. The injected call is responsible
				// for restoring LR. gentraceback is aware of
				// this extra slot. See sigctxt.pushCall in
				// signal_arm.go, which is similar except we
				// subtract 1 from IP here.
				sp := c.sp()
				sp -= goarch.PtrSize
				c.set_sp(sp)
				*(*uint32)(unsafe.Pointer(sp)) = uint32(c.lr())
				c.set_lr(newpc - 1)
				c.set_ip(targetPC)

			case "arm64":
				// Push LR. The injected call is responsible
				// for restoring LR. gentraceback is aware of
				// this extra slot. See sigctxt.pushCall in
				// signal_arm64.go.
				sp := c.sp() - 16 // SP needs 16-byte alignment
				c.set_sp(sp)
				*(*uint64)(unsafe.Pointer(sp)) = uint64(c.lr())
				c.set_lr(newpc)
				c.set_ip(targetPC)
			}
			stdcall2(_SetThreadContext, thread, uintptr(unsafe.Pointer(c)))
		}
	}

	atomic.Store(&mp.preemptExtLock, 0)

	// Acknowledge the preemption.
	mp.preemptGen.Add(1)

	stdcall1(_ResumeThread, thread)
	stdcall1(_CloseHandle, thread)
}
```

进程进行信号处理的核心逻辑位于 sighandler 函数中，在进行信号处理时，当遇到 sigPreempt 抢占信号时，触发运行时的异步抢占机制。

```go
// May run during STW, so write barriers are not allowed.
//
//go:nowritebarrierrec
func sighandler(_ureg *ureg, note *byte, gp *g) int {
	gsignal := getg()
	mp := gsignal.m

	var t sigTabT
	var docrash bool
	var sig int
	var flags int
	var level int32

	c := &sigctxt{_ureg}
	notestr := gostringnocopy(note)

	// The kernel will never pass us a nil note or ureg so we probably
	// made a mistake somewhere in sigtramp.
	if _ureg == nil || note == nil {
		print("sighandler: ureg ", _ureg, " note ", note, "\n")
		goto Throw
	}
	// Check that the note is no more than ERRMAX bytes (including
	// the trailing NUL). We should never receive a longer note.
	if len(notestr) > _ERRMAX-1 {
		print("sighandler: note is longer than ERRMAX\n")
		goto Throw
	}
	if isAbortPC(c.pc()) {
		// Never turn abort into a panic.
		goto Throw
	}
	// See if the note matches one of the patterns in sigtab.
	// Notes that do not match any pattern can be handled at a higher
	// level by the program but will otherwise be ignored.
	flags = _SigNotify
	for sig, t = range sigtable {
		if hasPrefix(notestr, t.name) {
			flags = t.flags
			break
		}
	}
	if flags&_SigPanic != 0 && gp.throwsplit {
		// We can't safely sigpanic because it may grow the
		// stack. Abort in the signal handler instead.
		flags = (flags &^ _SigPanic) | _SigThrow
	}
	if flags&_SigGoExit != 0 {
		exits((*byte)(add(unsafe.Pointer(note), 9))) // Strip "go: exit " prefix.
	}
	if flags&_SigPanic != 0 {
		// Copy the error string from sigtramp's stack into m->notesig so
		// we can reliably access it from the panic routines.
		memmove(unsafe.Pointer(mp.notesig), unsafe.Pointer(note), uintptr(len(notestr)+1))
		gp.sig = uint32(sig)
		gp.sigpc = c.pc()

		pc := c.pc()
		sp := c.sp()

		// If we don't recognize the PC as code
		// but we do recognize the top pointer on the stack as code,
		// then assume this was a call to non-code and treat like
		// pc == 0, to make unwinding show the context.
		if pc != 0 && !findfunc(pc).valid() && findfunc(*(*uintptr)(unsafe.Pointer(sp))).valid() {
			pc = 0
		}

		// IF LR exists, sigpanictramp must save it to the stack
		// before entry to sigpanic so that panics in leaf
		// functions are correctly handled. This will smash
		// the stack frame but we're not going back there
		// anyway.
		if usesLR {
			c.savelr(c.lr())
		}

		// If PC == 0, probably panicked because of a call to a nil func.
		// Not faking that as the return address will make the trace look like a call
		// to sigpanic instead. (Otherwise the trace will end at
		// sigpanic and we won't get to see who faulted).
		if pc != 0 {
			if usesLR {
				c.setlr(pc)
			} else {
				sp -= goarch.PtrSize
				*(*uintptr)(unsafe.Pointer(sp)) = pc
				c.setsp(sp)
			}
		}
		if usesLR {
			c.setpc(abi.FuncPCABI0(sigpanictramp))
		} else {
			c.setpc(abi.FuncPCABI0(sigpanic0))
		}
		return _NCONT
	}
	if flags&_SigNotify != 0 {
		if ignoredNote(note) {
			return _NCONT
		}
		if sendNote(note) {
			return _NCONT
		}
	}
	if flags&_SigKill != 0 {
		goto Exit
	}
	if flags&_SigThrow == 0 {
		return _NCONT
	}
Throw:
	mp.throwing = throwTypeRuntime
	mp.caughtsig.set(gp)
	startpanic_m()
	print(notestr, "\n")
	print("PC=", hex(c.pc()), "\n")
	print("\n")
	level, _, docrash = gotraceback()
	if level > 0 {
		goroutineheader(gp)
		tracebacktrap(c.pc(), c.sp(), c.lr(), gp)
		tracebackothers(gp)
		print("\n")
		dumpregs(_ureg)
	}
	if docrash {
		crash()
	}
Exit:
	goexitsall(note)
	exits(note)
	return _NDFLT // not reached
}
```

doSigPreempt 函数是平台相关的汇编函数。其中的重要一步是修改了原程序中 rsp、rip 寄存器中的值，从而在从内核态返回后，执行新的函数路径。在 Go 语言中，内核返回后执行新的 asyncPreempt 函数。asyncPreempt 函数会保存当前程序的寄存器值，并调用 asyncPreempt2 函数。当调用 asyncPreempt2 函数时，根据 preemptPark 函数或者 gopreempt_m 函数重新切换回调度循环，从而打断密集循环的继续执行。

```go
//go:nosplit
func asyncPreempt2() {
	gp := getg()
	gp.asyncSafePoint = true
	if gp.preemptStop {
		mcall(preemptPark)
	} else {
		mcall(gopreempt_m)
	}
	gp.asyncSafePoint = false
}
```

抢占调度的执行流程如图 15-11 所示。

![](../../../../assets/images/docs/internal/goroutine/design/README/图15-11%20抢占调度执行流程.png)

当发生系统调用时，当前正在工作的线程会陷入等待状态，等待内核完成系统调用并返回。当发生下面 3 种情况之一时，需要抢占调度：

- 当前局部运行队列中有等待运行的 G。在这种情况下，抢占调度只是为了让局部运行队列中的协程有执行的机会，因为其一般是当前 P 私有的。
- 当前没有空闲的 P 和自旋的 M。如果有空闲的 P 和自旋的 M，说明当前比较空闲，那么释放当前的 P 也没有太大意义。
- 当前系统调用的时间已经超过了 10ms，这和执行时间过长一样，需要立即抢占。

系统调用时的抢占原理主要是将 P 的状态转化为 _Pidle，这仅仅是完成了第 1 步。我们的目的是让 M 接管 P 的执行，主要的逻辑位于 handoffp 函数中，该函数需要判断是否需要找到一个新的 M 来接管当前的 P。当发生如下条件之一时，需要启动一个 M 来接管：

- 本地运行队列中有等待运行的 G。
- 需要处理一些垃圾回收的后台任务。
- 所有其他 P 都在运行 G，并且没有自旋的 M。
- 全局运行队列不为空。
- 需要处理网络 socket 读写等事件。

当这些条件都不满足时，才会将当前的 P 放入空闲队列中。

当寻找可用的 M 时，需要先在 M 的空闲列表中查找是否有闲置的 M，如果没有，则向操作系统申请一个新的 M，即线程。不管是唤醒闲置的线程还是新启动一个线程，都会开始新一轮调度。

这里有一个重要的问题——工作线程的 P 被抢占，系统调用的工作线程从内核返回后会怎么办呢？这涉及系统调用之前和之后执行的一系列逻辑。在执行实际操作系统调用之前，运行时调用了 reentersyscall 函数。该函数会保存当前 G 的执行环境，并解除 P 与 M 之间的绑定，将 P 放置到 oldp 中。解除绑定是为了系统调用返回后，当前的线程能够绑定不同的 P，但是会优先选择 oldp（如果 oldp 可以被绑定）。

```go
// The goroutine g is about to enter a system call.
// Record that it's not using the cpu anymore.
// This is called only from the go syscall library and cgocall,
// not from the low-level system calls used by the runtime.
//
// Entersyscall cannot split the stack: the save must
// make g->sched refer to the caller's stack segment, because
// entersyscall is going to return immediately after.
//
// Nothing entersyscall calls can split the stack either.
// We cannot safely move the stack during an active call to syscall,
// because we do not know which of the uintptr arguments are
// really pointers (back into the stack).
// In practice, this means that we make the fast path run through
// entersyscall doing no-split things, and the slow path has to use systemstack
// to run bigger things on the system stack.
//
// reentersyscall is the entry point used by cgo callbacks, where explicitly
// saved SP and PC are restored. This is needed when exitsyscall will be called
// from a function further up in the call stack than the parent, as g->syscallsp
// must always point to a valid stack frame. entersyscall below is the normal
// entry point for syscalls, which obtains the SP and PC from the caller.
//
// Syscall tracing:
// At the start of a syscall we emit traceGoSysCall to capture the stack trace.
// If the syscall does not block, that is it, we do not emit any other events.
// If the syscall blocks (that is, P is retaken), retaker emits traceGoSysBlock;
// when syscall returns we emit traceGoSysExit and when the goroutine starts running
// (potentially instantly, if exitsyscallfast returns true) we emit traceGoStart.
// To ensure that traceGoSysExit is emitted strictly after traceGoSysBlock,
// we remember current value of syscalltick in m (gp.m.syscalltick = gp.m.p.ptr().syscalltick),
// whoever emits traceGoSysBlock increments p.syscalltick afterwards;
// and we wait for the increment before emitting traceGoSysExit.
// Note that the increment is done even if tracing is not enabled,
// because tracing can be enabled in the middle of syscall. We don't want the wait to hang.
//
//go:nosplit
func reentersyscall(pc, sp uintptr) {
	gp := getg()

	// Disable preemption because during this function g is in Gsyscall status,
	// but can have inconsistent g->sched, do not let GC observe it.
	gp.m.locks++

	// Entersyscall must not call any function that might split/grow the stack.
	// (See details in comment above.)
	// Catch calls that might, by replacing the stack guard with something that
	// will trip any stack check and leaving a flag to tell newstack to die.
	gp.stackguard0 = stackPreempt
	gp.throwsplit = true

	// Leave SP around for GC and traceback.
	save(pc, sp)
	gp.syscallsp = sp
	gp.syscallpc = pc
	casgstatus(gp, _Grunning, _Gsyscall)
	if gp.syscallsp < gp.stack.lo || gp.stack.hi < gp.syscallsp {
		systemstack(func() {
			print("entersyscall inconsistent ", hex(gp.syscallsp), " [", hex(gp.stack.lo), ",", hex(gp.stack.hi), "]\n")
			throw("entersyscall")
		})
	}

	if trace.enabled {
		systemstack(traceGoSysCall)
		// systemstack itself clobbers g.sched.{pc,sp} and we might
		// need them later when the G is genuinely blocked in a
		// syscall
		save(pc, sp)
	}

	if sched.sysmonwait.Load() {
		systemstack(entersyscall_sysmon)
		save(pc, sp)
	}

	if gp.m.p.ptr().runSafePointFn != 0 {
		// runSafePointFn may stack split if run on this stack
		systemstack(runSafePointFn)
		save(pc, sp)
	}

	gp.m.syscalltick = gp.m.p.ptr().syscalltick
	gp.sysblocktraced = true
	pp := gp.m.p.ptr()
	pp.m = 0
	gp.m.oldp.set(pp)
	gp.m.p = 0
	atomic.Store(&pp.status, _Psyscall)
	if sched.gcwaiting.Load() {
		systemstack(entersyscall_gcwait)
		save(pc, sp)
	}

	gp.m.locks--
}
```

当操作系统内核返回系统调用后，被堵塞的协程继续执行，调用 exitsyscall 函数以便协程重新执行。

```go
// The goroutine g exited its system call.
// Arrange for it to run on a cpu again.
// This is called only from the go syscall library, not
// from the low-level system calls used by the runtime.
//
// Write barriers are not allowed because our P may have been stolen.
//
// This is exported via linkname to assembly in the syscall package.
//
//go:nosplit
//go:nowritebarrierrec
//go:linkname exitsyscall
func exitsyscall() {
	gp := getg()

	gp.m.locks++ // see comment in entersyscall
	if getcallersp() > gp.syscallsp {
		throw("exitsyscall: syscall frame is no longer valid")
	}

	gp.waitsince = 0
	oldp := gp.m.oldp.ptr()
	gp.m.oldp = 0
	if exitsyscallfast(oldp) {
		// When exitsyscallfast returns success, we have a P so can now use
		// write barriers
		if goroutineProfile.active {
			// Make sure that gp has had its stack written out to the goroutine
			// profile, exactly as it was when the goroutine profiler first
			// stopped the world.
			systemstack(func() {
				tryRecordGoroutineProfileWB(gp)
			})
		}
		if trace.enabled {
			if oldp != gp.m.p.ptr() || gp.m.syscalltick != gp.m.p.ptr().syscalltick {
				systemstack(traceGoStart)
			}
		}
		// There's a cpu for us, so we can run.
		gp.m.p.ptr().syscalltick++
		// We need to cas the status and scan before resuming...
		casgstatus(gp, _Gsyscall, _Grunning)

		// Garbage collector isn't running (since we are),
		// so okay to clear syscallsp.
		gp.syscallsp = 0
		gp.m.locks--
		if gp.preempt {
			// restore the preemption request in case we've cleared it in newstack
			gp.stackguard0 = stackPreempt
		} else {
			// otherwise restore the real _StackGuard, we've spoiled it in entersyscall/entersyscallblock
			gp.stackguard0 = gp.stack.lo + _StackGuard
		}
		gp.throwsplit = false

		if sched.disable.user && !schedEnabled(gp) {
			// Scheduling of this goroutine is disabled.
			Gosched()
		}

		return
	}

	gp.sysexitticks = 0
	if trace.enabled {
		// Wait till traceGoSysBlock event is emitted.
		// This ensures consistency of the trace (the goroutine is started after it is blocked).
		for oldp != nil && oldp.syscalltick == gp.m.syscalltick {
			osyield()
		}
		// We can't trace syscall exit right now because we don't have a P.
		// Tracing code can invoke write barriers that cannot run without a P.
		// So instead we remember the syscall exit time and emit the event
		// in execute when we have a P.
		gp.sysexitticks = cputicks()
	}

	gp.m.locks--

	// Call the scheduler.
	mcall(exitsyscall0)

	// Scheduler returned, so we're allowed to run now.
	// Delete the syscallsp information that we left for
	// the garbage collector during the system call.
	// Must wait until now because until gosched returns
	// we don't know for sure that the garbage collector
	// is not running.
	gp.syscallsp = 0
	gp.m.p.ptr().syscalltick++
	gp.throwsplit = false
}
```

由于在系统调用前，M 与 P 解除了绑定关系，因此现在 exitsyscall 函数希望能够重新绑定 P。寻找 P 的过程分为三个步骤：

- 尝试能否使用之前的 oldp，如果当前的 P 处于 _Psyscall 状态，则说明可以安全地绑定此 P。
- 当 P 不可使用时，说明其已经被系统监控线程分配给了其他的 M，此时加锁从全局空闲队列中寻找空闲的 P。
- 如果空闲队列中没有空闲的 P，则需要将当前的 G 放入全局运行队列，当前工作线程进入睡眠状态。当休眠被唤醒后，才能继续开始调度循环。

```go

```
