---
date: 2022-10-13T15:01:12+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "调度时机"  # 文章标题
url:  "posts/go/docs/internal/goroutine/design/occasion"  # 设置网页永久链接
tags: [ "Go", "occasion" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

可以根据调度方式的不同，将调度时机分为**主动、被动和抢占调度**。

## 主动调度

协程可以选择主动让渡自己的执行权利，这主要是通过用户在代码中执行 `runtime.Gosched` 函数实现的。在大多数情况下，用户并不需要执行此函数，因为 Go 语言编译器会在调用函数之前插入检查代码，判断该协程是否需要被抢占。

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

## 被动调度

被动调度指协程**在休眠、channel 通道堵塞、网络 I/O 堵塞、执行垃圾回收而暂停时，被动让渡自己执行权利的过程**。

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

## 抢占调度

为了让每个协程都有执行的机会，并且最大化利用 CPU 资源，Go 语言**在初始化时会启动一个特殊的线程来执行系统监控任务**。

**系统监控在一个独立的 M 上运行，不用绑定逻辑处理器 P，系统监控每隔 10ms 会检测是否有准备就绪的网络协程，并放置到全局队列中。和抢占调度相关的是，系统监控服务会判断当前协程是否运行时间过长，或者处于系统调用阶段，如果是，则会抢占当前 G 的执行。**

其核心逻辑位于 `runtime.retake` 函数中。

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

### 执行时间过长的抢占调度

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

**调度发生的时机主要在执行函数调用阶段。**函数调用是一个比较安全的检查点，Go 语言编译器会在函数调用前判断 stackguard0 的大小，从而选择是否调用 runtime.morestack_noctxt 函数。morestack_noctxt 为汇编函数，函数的执行流程如下：morestack_noctxt()→ morestack()→ newstack()。

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

这种抢占的方式面临着一定的问题，**当执行过程中没有函数调用，而只有类似如下代码时，协程将没有被抢占的机会。**

```go
for {
   i++
}
```

为了解决这一问题，Go 1.14 之后引入了**信号强制抢占的机制**。

### 信号强制抢占的机制

这需要借助图 15-10 中的类 UNIX 操作系统信号处理机制，信号是发送给进程的各种通知，以便将各种重要的事件通知给进程。

![](../../../../assets/images/docs/internal/goroutine/design/occasion/图15-10%20类UNIX操作系统信号处理机制.png)

最常见的是用户发送给进程的信号，例如时常使用的 CTRL+C 键，或者在命令行中输入的 `kill-<signal><PID>` 指令。通过信号，借助操作系统中断当前程序，保存程序的执行状态和寄存器值，并切换到内核态处理信号。在内核态处理完信号后，还会返回到用户态执行程序注册的信号处理函数，之后再回到内核，恢复程序原始的栈和寄存器值，并切换到用户态继续执行程序。

Go 语言**借助用户态在信号处理时完成协程的上下文切换的操作**，需要**借助进程对特定的信号进行处理**。并不是所有的信号都可以被处理，例如 `SIGKILL` 与 `SIGSTOP` 信号用于终止或暂停程序，不能被程序捕获处理。**Go 程序在初始化时会初始化信号表，并注册信号处理函数。**

在抢占时，调度器通过向线程中发送 **`sigPreempt`** 信号，触发信号处理。在 UNIX 操作系统中，`sigPreempt` 为 `_SIGURG` 信号，由于该信号不会被用户程序和调试器使用，因此 Go 语言使用它作为安全的抢占信号，关于信号具体的选择过程，可以参考 Go 源码中对 `_SIGURG` 信号的注释。

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

进程进行信号处理的核心逻辑位于 `sighandler` 函数中，在进行信号处理时，当遇到 `sigPreempt` 抢占信号时，触发运行时的异步抢占机制。

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

`doSigPreempt` 函数是平台相关的汇编函数。其中的重要一步是修改了原程序中 `rsp`、`rip` 寄存器中的值，从而在从内核态返回后，执行新的函数路径。在 Go 语言中，内核返回后执行新的 `asyncPreempt` 函数。`asyncPreempt` 函数会保存当前程序的寄存器值，并调用 `asyncPreempt2` 函数。当调用 `asyncPreempt2` 函数时，根据 `preemptPark` 函数或者 `gopreempt_m` 函数重新切换回调度循环，从而打断密集循环的继续执行。

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

![](../../../../assets/images/docs/internal/goroutine/design/occasion/图15-11%20抢占调度执行流程.png)

当发生系统调用时，当前正在工作的线程会陷入等待状态，等待内核完成系统调用并返回。当发生下面 3 种情况之一时，需要抢占调度：

- 当前局部运行队列中有等待运行的 G。在这种情况下，抢占调度只是为了让局部运行队列中的协程有执行的机会，因为其一般是当前 P 私有的。
- 当前没有空闲的 P 和自旋的 M。如果有空闲的 P 和自旋的 M，说明当前比较空闲，那么释放当前的 P 也没有太大意义。
- 当前系统调用的时间已经超过了 10ms，这和执行时间过长一样，需要立即抢占。

系统调用时的抢占原理主要是将 P 的状态转化为 `_Pidle`，这仅仅是完成了第 1 步。我们的目的是让 M 接管 P 的执行，主要的逻辑位于 `handoffp` 函数中，该函数需要判断是否需要找到一个新的 M 来接管当前的 P。当发生如下条件之一时，需要启动一个 M 来接管：

- 本地运行队列中有等待运行的 G。
- 需要处理一些垃圾回收的后台任务。
- 所有其他 P 都在运行 G，并且没有自旋的 M。
- 全局运行队列不为空。
- 需要处理网络 socket 读写等事件。

当这些条件都不满足时，才会将当前的 P 放入空闲队列中。

当寻找可用的 M 时，需要先在 M 的空闲列表中查找是否有闲置的 M，如果没有，则向操作系统申请一个新的 M，即线程。不管是唤醒闲置的线程还是新启动一个线程，都会开始新一轮调度。

这里有一个重要的问题——工作线程的 P 被抢占，系统调用的工作线程从内核返回后会怎么办呢？这涉及系统调用之前和之后执行的一系列逻辑。在执行实际操作系统调用之前，运行时调用了 `reentersyscall` 函数。该函数会保存当前 G 的执行环境，并解除 P 与 M 之间的绑定，将 P 放置到 `oldp` 中。解除绑定是为了系统调用返回后，当前的线程能够绑定不同的 P，但是会优先选择 `oldp`（如果 `oldp` 可以被绑定）。

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

当操作系统内核返回系统调用后，被堵塞的协程继续执行，调用 `exitsyscall` 函数以便协程重新执行。

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

由于在系统调用前，M 与 P 解除了绑定关系，因此现在 `exitsyscall` 函数希望能够重新绑定 P。寻找 P 的过程分为三个步骤：

- 尝试能否使用之前的 `oldp`，如果当前的 P 处于 `_Psyscall` 状态，则说明可以安全地绑定此 P。
- 当 P 不可使用时，说明其已经被系统监控线程分配给了其他的 M，此时加锁从全局空闲队列中寻找空闲的 P。
- 如果空闲队列中没有空闲的 P，则需要将当前的 G 放入全局运行队列，当前工作线程进入睡眠状态。当休眠被唤醒后，才能继续开始调度循环。

```go

```
