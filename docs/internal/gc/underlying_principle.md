---
date: 2022-10-11T09:18:13+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "深入垃圾回收全流程"  # 文章标题
url:  "posts/go/docs/internal/gc/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

垃圾回收贯穿于程序的整个生命周期，运行时将循环不断地检测当前程序的内存使用状态并选择在合适的时机执行垃圾回收。

## 垃圾回收循环

Go 语言的垃圾回收循环大致会经历如图 20-1 所示的几个阶段。当内存到达了垃圾回收的阈值后，将触发新一轮的垃圾回收。之后会先后经历标记准备阶段、并行标记阶段、标记终止阶段及垃圾清扫阶段。在并行标记阶段引入了辅助标记技术，在垃圾清扫阶段还引入了辅助清扫、系统驻留内存清除技术。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-1%20垃圾回收循环.png)

## 标记准备阶段

标记准备阶段最重要的任务是清扫上一阶段 GC 遗留的需要清扫的对象，因为使用了懒清扫算法，所以当执行下一次 GC 时，可能还有垃圾对象没有被清扫。同时，标记准备阶段会重置各种状态和统计指标、启动专门用于标记的协程、统计需要扫描的任务数量、开启写屏障、启动标记协程等。总之，标记准备阶段是初始阶段，执行轻量级的任务。在标记准备阶段，上面大部分重要的步骤需要在 STW(Stop The World) 时进行。

标记准备阶段会为每个逻辑处理器 P 启动一个标记协程，但并不是所有的标记协程都有执行的机会，因为在标记阶段，标记协程与正常执行用户代码的协程需要并行，以减少 GC 给用户程序带来的影响。在这里，需要关注标记准备阶段两个重要的问题——如何决定需要多少标记协程以及如何调度标记协程。

### 计算标记协程的数量

在标记准备阶段，会计算当前后台需要开启多少标记协程。目前，Go 语言规定后台标记协程消耗的 CPU 应该接近 25%，其核心代码位于 startCycle 函数中。

`src/runtime/mgcpacer.go`

```go
// startCycle resets the GC controller's state and computes estimates
// for a new GC cycle. The caller must hold worldsema and the world
// must be stopped.
func (c *gcControllerState) startCycle(markStartTime int64, procs int, trigger gcTrigger) {
	c.heapScanWork.Store(0)
	c.stackScanWork.Store(0)
	c.globalsScanWork.Store(0)
	c.bgScanCredit.Store(0)
	c.assistTime.Store(0)
	c.dedicatedMarkTime.Store(0)
	c.fractionalMarkTime.Store(0)
	c.idleMarkTime.Store(0)
	c.markStartTime = markStartTime

	// TODO(mknyszek): This is supposed to be the actual trigger point for the heap, but
	// causes regressions in memory use. The cause is that the PI controller used to smooth
	// the cons/mark ratio measurements tends to flail when using the less accurate precomputed
	// trigger for the cons/mark calculation, and this results in the controller being more
	// conservative about steady-states it tries to find in the future.
	//
	// This conservatism is transient, but these transient states tend to matter for short-lived
	// programs, especially because the PI controller is overdamped, partially because it is
	// configured with a relatively large time constant.
	//
	// Ultimately, I think this is just two mistakes piled on one another: the choice of a swingy
	// smoothing function that recalls a fairly long history (due to its overdamped time constant)
	// coupled with an inaccurate cons/mark calculation. It just so happens this works better
	// today, and it makes it harder to change things in the future.
	//
	// This is described in #53738. Fix this for #53892 by changing back to the actual trigger
	// point and simplifying the smoothing function.
	heapTrigger, heapGoal := c.trigger()
	c.triggered = heapTrigger

	// Compute the background mark utilization goal. In general,
	// this may not come out exactly. We round the number of
	// dedicated workers so that the utilization is closest to
	// 25%. For small GOMAXPROCS, this would introduce too much
	// error, so we add fractional workers in that case.
	totalUtilizationGoal := float64(procs) * gcBackgroundUtilization
	dedicatedMarkWorkersNeeded := int64(totalUtilizationGoal + 0.5)
	utilError := float64(dedicatedMarkWorkersNeeded)/totalUtilizationGoal - 1
	const maxUtilError = 0.3
	if utilError < -maxUtilError || utilError > maxUtilError {
		// Rounding put us more than 30% off our goal. With
		// gcBackgroundUtilization of 25%, this happens for
		// GOMAXPROCS<=3 or GOMAXPROCS=6. Enable fractional
		// workers to compensate.
		if float64(dedicatedMarkWorkersNeeded) > totalUtilizationGoal {
			// Too many dedicated workers.
			dedicatedMarkWorkersNeeded--
		}
		c.fractionalUtilizationGoal = (totalUtilizationGoal - float64(dedicatedMarkWorkersNeeded)) / float64(procs)
	} else {
		c.fractionalUtilizationGoal = 0
	}

	// In STW mode, we just want dedicated workers.
	if debug.gcstoptheworld > 0 {
		dedicatedMarkWorkersNeeded = int64(procs)
		c.fractionalUtilizationGoal = 0
	}

	// Clear per-P state
	for _, p := range allp {
		p.gcAssistTime = 0
		p.gcFractionalMarkTime = 0
	}

	if trigger.kind == gcTriggerTime {
		// During a periodic GC cycle, reduce the number of idle mark workers
		// required. However, we need at least one dedicated mark worker or
		// idle GC worker to ensure GC progress in some scenarios (see comment
		// on maxIdleMarkWorkers).
		if dedicatedMarkWorkersNeeded > 0 {
			c.setMaxIdleMarkWorkers(0)
		} else {
			// TODO(mknyszek): The fundamental reason why we need this is because
			// we can't count on the fractional mark worker to get scheduled.
			// Fix that by ensuring it gets scheduled according to its quota even
			// if the rest of the application is idle.
			c.setMaxIdleMarkWorkers(1)
		}
	} else {
		// N.B. gomaxprocs and dedicatedMarkWorkersNeeded are guaranteed not to
		// change during a GC cycle.
		c.setMaxIdleMarkWorkers(int32(procs) - int32(dedicatedMarkWorkersNeeded))
	}

	// Compute initial values for controls that are updated
	// throughout the cycle.
	c.dedicatedMarkWorkersNeeded.Store(dedicatedMarkWorkersNeeded)
	c.revise()

	if debug.gcpacertrace > 0 {
		assistRatio := c.assistWorkPerByte.Load()
		print("pacer: assist ratio=", assistRatio,
			" (scan ", gcController.heapScan.Load()>>20, " MB in ",
			work.initialHeapLive>>20, "->",
			heapGoal>>20, " MB)",
			" workers=", dedicatedMarkWorkersNeeded,
			"+", c.fractionalUtilizationGoal, "\n")
	}
}
```

一种简单的想法是根据当前逻辑处理器 P 的数量来计算，开启的协程数量应该为 0.25P。为什么 startCycle 函数的计算过程如此复杂呢？这是因为需要处理当协程数量过小（例如 P ≤ 3）、不为整数（0.25P）的情况。

而 fractionalUtilizationGoal 是一个附加的参数，其小于 1。例如当 P = 2 时，其值为 0.25。代表每个 P 在标记阶段需要花 25% 的时间执行后台标记协程。

fractionalUtilizationGoal 是专门为 P 为 1、2、3、6 时而设计的，例如当 P = 2 时，2×0.25 = 0.5，即只能花 0.5 个 P 来执行标记任务，但如果用一个 P 来执行后台任务，这时标记的 CPU 使用量就变为了 1/2 = 50%，这和 25% 的 CPU 使用率的设计目标差距太大。

所以，当 P = 2 时，fractionalUtilizationGoal 的计算结果为 0.25。如图 20-3 所示，它表示在总的标记周期 t 内，每个 P 都需要花 25% 的时间来执行后台标记工作。这是一种基于时间的调度。当超出时间后，当前的后台标记协程可以被抢占，从而执行其他的协程。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-3%20特殊数量逻辑处理器下的25%时间调度.png)

### 切换到后台标记协程

标记准备阶段的第 2 个问题是如何调度标记协程。在标记准备阶段执行了 STW，此时暂停了所有协程。可以预料到，当关闭 STW 准备再次启动所有协程时，每个逻辑处理器 P 都会进入一轮新的调度循环，在调度循环开始时，调度器会判断程序是否处于 GC 阶段，如果是，则尝试判断当前 P 是否需要执行后台标记任务。

`src/runtime/proc.go:findRunnable`

```go
	// Try to schedule a GC worker.
	if gcBlackenEnabled != 0 {
		gp, tnow := gcController.findRunnableGCWorker(pp, now)
		if gp != nil {
			return gp, false, true
		}
		now = tnow
	}
```

在 findRunnableGCWorker 函数中，如果代表执行完整的后台标记协程的字段 dedicatedMarkWorkersNeeded 大于 0，则当前协程立即执行后台标记任务。如果参数 fractionalUtilizationGoal 大于 0，并且当前逻辑处理器 P 执行标记任务的时间小于 fractionalUtilizationGoal× 当前标记周期的时间，则仍然会执行后台标记任务，但是并不会在整个标记周期内一直执行。此时，后台标记协程的运行模式会切换为 gcMarkWorkerFractionalMode，如下所示。

`src/runtime/mgcpacer.go`

```go
// findRunnableGCWorker returns a background mark worker for pp if it
// should be run. This must only be called when gcBlackenEnabled != 0.
func (c *gcControllerState) findRunnableGCWorker(pp *p, now int64) (*g, int64) {
	if gcBlackenEnabled == 0 {
		throw("gcControllerState.findRunnable: blackening not enabled")
	}

	// Since we have the current time, check if the GC CPU limiter
	// hasn't had an update in a while. This check is necessary in
	// case the limiter is on but hasn't been checked in a while and
	// so may have left sufficient headroom to turn off again.
	if now == 0 {
		now = nanotime()
	}
	if gcCPULimiter.needUpdate(now) {
		gcCPULimiter.update(now)
	}

	if !gcMarkWorkAvailable(pp) {
		// No work to be done right now. This can happen at
		// the end of the mark phase when there are still
		// assists tapering off. Don't bother running a worker
		// now because it'll just return immediately.
		return nil, now
	}

	// Grab a worker before we commit to running below.
	node := (*gcBgMarkWorkerNode)(gcBgMarkWorkerPool.pop())
	if node == nil {
		// There is at least one worker per P, so normally there are
		// enough workers to run on all Ps, if necessary. However, once
		// a worker enters gcMarkDone it may park without rejoining the
		// pool, thus freeing a P with no corresponding worker.
		// gcMarkDone never depends on another worker doing work, so it
		// is safe to simply do nothing here.
		//
		// If gcMarkDone bails out without completing the mark phase,
		// it will always do so with queued global work. Thus, that P
		// will be immediately eligible to re-run the worker G it was
		// just using, ensuring work can complete.
		return nil, now
	}

	decIfPositive := func(val *atomic.Int64) bool {
		for {
			v := val.Load()
			if v <= 0 {
				return false
			}

			if val.CompareAndSwap(v, v-1) {
				return true
			}
		}
	}

	if decIfPositive(&c.dedicatedMarkWorkersNeeded) {
		// This P is now dedicated to marking until the end of
		// the concurrent mark phase.
		pp.gcMarkWorkerMode = gcMarkWorkerDedicatedMode
	} else if c.fractionalUtilizationGoal == 0 {
		// No need for fractional workers.
		gcBgMarkWorkerPool.push(&node.node)
		return nil, now
	} else {
		// Is this P behind on the fractional utilization
		// goal?
		//
		// This should be kept in sync with pollFractionalWorkerExit.
		delta := now - c.markStartTime
		if delta > 0 && float64(pp.gcFractionalMarkTime)/float64(delta) > c.fractionalUtilizationGoal {
			// Nope. No need to run a fractional worker.
			gcBgMarkWorkerPool.push(&node.node)
			return nil, now
		}
		// Run a fractional worker.
		pp.gcMarkWorkerMode = gcMarkWorkerFractionalMode
	}

	// Run the background mark worker.
	gp := node.gp.ptr()
	casgstatus(gp, _Gwaiting, _Grunnable)
	if trace.enabled {
		traceGoUnpark(gp, 0)
	}
	return gp, now
}
```

## 并发标记阶段

在并发标记阶段，后台标记协程可以与执行用户代码的协程并行。Go 语言的目标是后台标记协程占用 CPU 的时间为 25%，以最大限度地避免因执行 GC 而中断或减慢用户协程的执行。如图 20-4 所示，后台标记任务有 3 种不同的模式。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-4%20后台标记任务的模式.png)

DedicatedMode 代表处理器专门负责标记对象，不会被调度器抢占。

FractionalMode 代表协助后台标记，在标记阶段到达目标时间后，会自动退出。

IdleMode 代表当处理器没有查找到可以执行的协程时，执行垃圾收集的标记任务，直到被抢占。标记阶段的核心逻辑位于 gcDrain 函数，其中第 2 个参数为后台标记 flag，大部分 flag 和后台标记协程 3 种不同的模式有关。

```go
// gcDrain scans roots and objects in work buffers, blackening grey
// objects until it is unable to get more work. It may return before
// GC is done; it's the caller's responsibility to balance work from
// other Ps.
//
// If flags&gcDrainUntilPreempt != 0, gcDrain returns when g.preempt
// is set.
//
// If flags&gcDrainIdle != 0, gcDrain returns when there is other work
// to do.
//
// If flags&gcDrainFractional != 0, gcDrain self-preempts when
// pollFractionalWorkerExit() returns true. This implies
// gcDrainNoBlock.
//
// If flags&gcDrainFlushBgCredit != 0, gcDrain flushes scan work
// credit to gcController.bgScanCredit every gcCreditSlack units of
// scan work.
//
// gcDrain will always return if there is a pending STW.
//
//go:nowritebarrier
func gcDrain(gcw *gcWork, flags gcDrainFlags) {
	if !writeBarrier.needed {
		throw("gcDrain phase incorrect")
	}

	gp := getg().m.curg
	preemptible := flags&gcDrainUntilPreempt != 0
	flushBgCredit := flags&gcDrainFlushBgCredit != 0
	idle := flags&gcDrainIdle != 0

	initScanWork := gcw.heapScanWork

	// checkWork is the scan work before performing the next
	// self-preempt check.
	checkWork := int64(1<<63 - 1)
	var check func() bool
	if flags&(gcDrainIdle|gcDrainFractional) != 0 {
		checkWork = initScanWork + drainCheckThreshold
		if idle {
			check = pollWork
		} else if flags&gcDrainFractional != 0 {
			check = pollFractionalWorkerExit
		}
	}

	// Drain root marking jobs.
	if work.markrootNext < work.markrootJobs {
		// Stop if we're preemptible or if someone wants to STW.
		for !(gp.preempt && (preemptible || sched.gcwaiting.Load())) {
			job := atomic.Xadd(&work.markrootNext, +1) - 1
			if job >= work.markrootJobs {
				break
			}
			markroot(gcw, job, flushBgCredit)
			if check != nil && check() {
				goto done
			}
		}
	}

	// Drain heap marking jobs.
	// Stop if we're preemptible or if someone wants to STW.
	for !(gp.preempt && (preemptible || sched.gcwaiting.Load())) {
		// Try to keep work available on the global queue. We used to
		// check if there were waiting workers, but it's better to
		// just keep work available than to make workers wait. In the
		// worst case, we'll do O(log(_WorkbufSize)) unnecessary
		// balances.
		if work.full == 0 {
			gcw.balance()
		}

		b := gcw.tryGetFast()
		if b == 0 {
			b = gcw.tryGet()
			if b == 0 {
				// Flush the write barrier
				// buffer; this may create
				// more work.
				wbBufFlush(nil, 0)
				b = gcw.tryGet()
			}
		}
		if b == 0 {
			// Unable to get work.
			break
		}
		scanobject(b, gcw)

		// Flush background scan work credit to the global
		// account if we've accumulated enough locally so
		// mutator assists can draw on it.
		if gcw.heapScanWork >= gcCreditSlack {
			gcController.heapScanWork.Add(gcw.heapScanWork)
			if flushBgCredit {
				gcFlushBgCredit(gcw.heapScanWork - initScanWork)
				initScanWork = 0
			}
			checkWork -= gcw.heapScanWork
			gcw.heapScanWork = 0

			if checkWork <= 0 {
				checkWork += drainCheckThreshold
				if check != nil && check() {
					break
				}
			}
		}
	}

done:
	// Flush remaining scan work credit.
	if gcw.heapScanWork > 0 {
		gcController.heapScanWork.Add(gcw.heapScanWork)
		if flushBgCredit {
			gcFlushBgCredit(gcw.heapScanWork - initScanWork)
		}
		gcw.heapScanWork = 0
	}
}
```

后台标记 flag 有 4 种，用于指定后台标记协程的不同行为。

```go
type gcDrainFlags int

const (
	gcDrainUntilPreempt gcDrainFlags = 1 << iota
	gcDrainFlushBgCredit
	gcDrainIdle
	gcDrainFractional
)
```

- gcDrainUntilPreempt 标记代表当前标记协程处于可以被抢占的状态。
- gcDrainFlushBgCredit 标记会计算后台完成的标记任务量以减少并行标记期间用户协程执行辅助垃圾收集的工作量。
- gcDrainIdle 标记对应 IdleMode 模式，表示当处理器上包含其他待执行的协程时标记协程退出。
- gcDrainFractional 标记对应 FractionalMode 模式，表示后台标记协程到达目标时间后退出。

在 DedicatedMode 下，会一直执行后台标记任务，这意味着当前逻辑处理器 P 本地队列中的协程将一直得不到执行，这是不能接受的。所以 Go 语言的做法是先执行可以被抢占的后台标记任务，如果标记协程已经被其他协程抢占，那么当前的逻辑处理器 P 并不会执行其他协程，而是将其他协程转移到全局队列中，并取消 gcDrainUntilPreempt 标志，进入不能被抢占的模式。

`src/runtime/mgc.go`

```go
func gcBgMarkWorker() {
	gp := getg()

	// We pass node to a gopark unlock function, so it can't be on
	// the stack (see gopark). Prevent deadlock from recursively
	// starting GC by disabling preemption.
	gp.m.preemptoff = "GC worker init"
	node := new(gcBgMarkWorkerNode)
	gp.m.preemptoff = ""

	node.gp.set(gp)

	node.m.set(acquirem())
	notewakeup(&work.bgMarkReady)
	// After this point, the background mark worker is generally scheduled
	// cooperatively by gcController.findRunnableGCWorker. While performing
	// work on the P, preemption is disabled because we are working on
	// P-local work buffers. When the preempt flag is set, this puts itself
	// into _Gwaiting to be woken up by gcController.findRunnableGCWorker
	// at the appropriate time.
	//
	// When preemption is enabled (e.g., while in gcMarkDone), this worker
	// may be preempted and schedule as a _Grunnable G from a runq. That is
	// fine; it will eventually gopark again for further scheduling via
	// findRunnableGCWorker.
	//
	// Since we disable preemption before notifying bgMarkReady, we
	// guarantee that this G will be in the worker pool for the next
	// findRunnableGCWorker. This isn't strictly necessary, but it reduces
	// latency between _GCmark starting and the workers starting.

	for {
		// Go to sleep until woken by
		// gcController.findRunnableGCWorker.
		gopark(func(g *g, nodep unsafe.Pointer) bool {
			node := (*gcBgMarkWorkerNode)(nodep)

			if mp := node.m.ptr(); mp != nil {
				// The worker G is no longer running; release
				// the M.
				//
				// N.B. it is _safe_ to release the M as soon
				// as we are no longer performing P-local mark
				// work.
				//
				// However, since we cooperatively stop work
				// when gp.preempt is set, if we releasem in
				// the loop then the following call to gopark
				// would immediately preempt the G. This is
				// also safe, but inefficient: the G must
				// schedule again only to enter gopark and park
				// again. Thus, we defer the release until
				// after parking the G.
				releasem(mp)
			}

			// Release this G to the pool.
			gcBgMarkWorkerPool.push(&node.node)
			// Note that at this point, the G may immediately be
			// rescheduled and may be running.
			return true
		}, unsafe.Pointer(node), waitReasonGCWorkerIdle, traceEvGoBlock, 0)

		// Preemption must not occur here, or another G might see
		// p.gcMarkWorkerMode.

		// Disable preemption so we can use the gcw. If the
		// scheduler wants to preempt us, we'll stop draining,
		// dispose the gcw, and then preempt.
		node.m.set(acquirem())
		pp := gp.m.p.ptr() // P can't change with preemption disabled.

		if gcBlackenEnabled == 0 {
			println("worker mode", pp.gcMarkWorkerMode)
			throw("gcBgMarkWorker: blackening not enabled")
		}

		if pp.gcMarkWorkerMode == gcMarkWorkerNotWorker {
			throw("gcBgMarkWorker: mode not set")
		}

		startTime := nanotime()
		pp.gcMarkWorkerStartTime = startTime
		var trackLimiterEvent bool
		if pp.gcMarkWorkerMode == gcMarkWorkerIdleMode {
			trackLimiterEvent = pp.limiterEvent.start(limiterEventIdleMarkWork, startTime)
		}

		decnwait := atomic.Xadd(&work.nwait, -1)
		if decnwait == work.nproc {
			println("runtime: work.nwait=", decnwait, "work.nproc=", work.nproc)
			throw("work.nwait was > work.nproc")
		}

		systemstack(func() {
			// Mark our goroutine preemptible so its stack
			// can be scanned. This lets two mark workers
			// scan each other (otherwise, they would
			// deadlock). We must not modify anything on
			// the G stack. However, stack shrinking is
			// disabled for mark workers, so it is safe to
			// read from the G stack.
			casGToWaiting(gp, _Grunning, waitReasonGCWorkerActive)
			switch pp.gcMarkWorkerMode {
			default:
				throw("gcBgMarkWorker: unexpected gcMarkWorkerMode")
			case gcMarkWorkerDedicatedMode:
				gcDrain(&pp.gcw, gcDrainUntilPreempt|gcDrainFlushBgCredit)
				if gp.preempt {
					// We were preempted. This is
					// a useful signal to kick
					// everything out of the run
					// queue so it can run
					// somewhere else.
					if drainQ, n := runqdrain(pp); n > 0 {
						lock(&sched.lock)
						globrunqputbatch(&drainQ, int32(n))
						unlock(&sched.lock)
					}
				}
				// Go back to draining, this time
				// without preemption.
				gcDrain(&pp.gcw, gcDrainFlushBgCredit)
			case gcMarkWorkerFractionalMode:
				gcDrain(&pp.gcw, gcDrainFractional|gcDrainUntilPreempt|gcDrainFlushBgCredit)
			case gcMarkWorkerIdleMode:
				gcDrain(&pp.gcw, gcDrainIdle|gcDrainUntilPreempt|gcDrainFlushBgCredit)
			}
			casgstatus(gp, _Gwaiting, _Grunning)
		})

		// Account for time and mark us as stopped.
		now := nanotime()
		duration := now - startTime
		gcController.markWorkerStop(pp.gcMarkWorkerMode, duration)
		if trackLimiterEvent {
			pp.limiterEvent.stop(limiterEventIdleMarkWork, now)
		}
		if pp.gcMarkWorkerMode == gcMarkWorkerFractionalMode {
			atomic.Xaddint64(&pp.gcFractionalMarkTime, duration)
		}

		// Was this the last worker and did we run out
		// of work?
		incnwait := atomic.Xadd(&work.nwait, +1)
		if incnwait > work.nproc {
			println("runtime: p.gcMarkWorkerMode=", pp.gcMarkWorkerMode,
				"work.nwait=", incnwait, "work.nproc=", work.nproc)
			throw("work.nwait > work.nproc")
		}

		// We'll releasem after this point and thus this P may run
		// something else. We must clear the worker mode to avoid
		// attributing the mode to a different (non-worker) G in
		// traceGoStart.
		pp.gcMarkWorkerMode = gcMarkWorkerNotWorker

		// If this worker reached a background mark completion
		// point, signal the main GC goroutine.
		if incnwait == work.nproc && !gcMarkWorkAvailable(nil) {
			// We don't need the P-local buffers here, allow
			// preemption because we may schedule like a regular
			// goroutine in gcMarkDone (block on locks, etc).
			releasem(node.m.ptr())
			node.m.set(nil)

			gcMarkDone()
		}
	}
}
```

FractionalMode 模式和 IdleMode 模式都允许被抢占。除此之外，FractionalMode 模式加上了 gcDrainFractional 标记表明当前标记协程会在到达目标时间后退出，IdleMode 模式加上了 gcDrainIdle 标记表明在发现有其他协程可以运行时退出当前标记协程。三种模式都加上了 gcDrainFlushBgCredit 标志，用于计算后台完成的标记任务量，并唤醒之前由于分配内存太频繁而陷入等待的用户协程。

### 根对象扫描

扫描的第一阶段是扫描根对象。在最开始的标记准备阶段会统计这次 GC 一共要扫描多少对象，每个具体的序号都对应着要扫描的对象，如下所示。

```go
job := atomic.Xadd(&work.markrootNext, 1) - 1
```

work.markrootNext 必须原子增加。这是因为可能出现图 20-6 中多个后台标记协程同时访问该变量的情况，这种机制保证了多个后台标记协程能够并发执行不同的任务。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-6%20根对象的并行扫描机制.png)

那么何为根对象呢？根对象是最基本的对象，从根对象出发，可以找到所有的引用对象（即活着的对象）。在 Go 语言中，根对象包括全局变量（在 .bss 和 .data 段内存中）、span 中 finalizer 的任务数量，以及所有的协程栈。finalizer 是 Go 语言中对象绑定的析构器，当对象的内存释放后，需要调用析构器函数，从而完整释放资源。例如，os.File 对象使用析构器函数关闭操作系统文件描述符，即便用户忘记了调用 close 方法也会释放操作系统资源。

### 全局变量扫描

扫描全局变量需要编译时与运行时的共同努力。只有在运行时才能确定全局变量被分配到虚拟内存的哪一个区域，另外，如果全局变量有指针，那么在运行时其指针指向的内存可能变化。而在编译时，可以确定全局变量中哪些位置包含指针，如图 20-7 所示，信息位于位图 ptrmask 字段中。ptrmask 的每个 bit 位都对应了.data 段中一个指针的大小（8byte），bit 位为 1 代表当前位置是一个指针，这时，需要求出当前的指针在堆区的哪一个对象上，并将当前对象标记为灰色。

如何通过指针找到指针对应的对象位置呢？这靠的是 Go 语言对内存的精细化管理。如图 20-8 所示，可以先找到指针在哪一个 heapArena 中，heapArena 是内存分配时每一次向操作系统申请的最小 64MB 的区域。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-7%20编译时确定的指针位图.png)

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-8%20通过指针查找到对象在span中的位置.png)

heapArena 存储了许多元数据，其中包括每个 page（8 KB）对应的 mspan。

所以，可以进一步通过指针的位置找到其对应的 mspan，进而找到其位于 mspan 中第几个元素中。当找到此元素后，会将 gcmarkBits 位图对应元素的 bit 设置为 1，表明其已经被标记，同时将该元素（对象）放入标记队列中。

在 span 中，位图 gcmarkBits 中的每个元素都有标志位表明当前元素中的对象是否被标记，如图 20-9 所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-9%20span中表明对象是否被标记的位图.png)

### finalizer

之前提到 finalizer 是特殊的对象，其是在对象释放后会被调用的析构器，用于资源释放。析构器不会被栈上或全局变量引用，需要单独处理。

如下所示，在标记期间，后台标记协程会遍历 mspan 中的 specials 链表，扫描 finalizer 所位于的元素（对象），并扫描当前元素（对象），扫描对象的详细过程将在下一节介绍。注意在这里，并不能把 finalizer 所位于的 span 中的对象加入根对象中，否则我们将失去回收该对象的机会。同时需要扫描析构器字段 fn，因为 fn 可能指向了堆中的内存，并可能被回收。

```go
// markrootSpans marks roots for one shard of markArenas.
//
//go:nowritebarrier
func markrootSpans(gcw *gcWork, shard int) {
	// Objects with finalizers have two GC-related invariants:
	//
	// 1) Everything reachable from the object must be marked.
	// This ensures that when we pass the object to its finalizer,
	// everything the finalizer can reach will be retained.
	//
	// 2) Finalizer specials (which are not in the garbage
	// collected heap) are roots. In practice, this means the fn
	// field must be scanned.
	sg := mheap_.sweepgen

	// Find the arena and page index into that arena for this shard.
	ai := mheap_.markArenas[shard/(pagesPerArena/pagesPerSpanRoot)]
	ha := mheap_.arenas[ai.l1()][ai.l2()]
	arenaPage := uint(uintptr(shard) * pagesPerSpanRoot % pagesPerArena)

	// Construct slice of bitmap which we'll iterate over.
	specialsbits := ha.pageSpecials[arenaPage/8:]
	specialsbits = specialsbits[:pagesPerSpanRoot/8]
	for i := range specialsbits {
		// Find set bits, which correspond to spans with specials.
		specials := atomic.Load8(&specialsbits[i])
		if specials == 0 {
			continue
		}
		for j := uint(0); j < 8; j++ {
			if specials&(1<<j) == 0 {
				continue
			}
			// Find the span for this bit.
			//
			// This value is guaranteed to be non-nil because having
			// specials implies that the span is in-use, and since we're
			// currently marking we can be sure that we don't have to worry
			// about the span being freed and re-used.
			s := ha.spans[arenaPage+uint(i)*8+j]

			// The state must be mSpanInUse if the specials bit is set, so
			// sanity check that.
			if state := s.state.get(); state != mSpanInUse {
				print("s.state = ", state, "\n")
				throw("non in-use span found with specials bit set")
			}
			// Check that this span was swept (it may be cached or uncached).
			if !useCheckmark && !(s.sweepgen == sg || s.sweepgen == sg+3) {
				// sweepgen was updated (+2) during non-checkmark GC pass
				print("sweep ", s.sweepgen, " ", sg, "\n")
				throw("gc: unswept span")
			}

			// Lock the specials to prevent a special from being
			// removed from the list while we're traversing it.
			lock(&s.speciallock)
			for sp := s.specials; sp != nil; sp = sp.next {
				if sp.kind != _KindSpecialFinalizer {
					continue
				}
				// don't mark finalized object, but scan it so we
				// retain everything it points to.
				spf := (*specialfinalizer)(unsafe.Pointer(sp))
				// A finalizer can be set for an inner byte of an object, find object beginning.
				p := s.base() + uintptr(spf.special.offset)/s.elemsize*s.elemsize

				// Mark everything that can be reached from
				// the object (but *not* the object itself or
				// we'll never collect it).
				if !s.spanclass.noscan() {
					scanobject(p, gcw)
				}

				// The special itself is a root.
				scanblock(uintptr(unsafe.Pointer(&spf.fn)), goarch.PtrSize, &oneptrmask[0], gcw, nil)
			}
			unlock(&s.speciallock)
		}
	}
}
```

使用 finalizer 可以实现一些有趣的功能，例如 Go 语言的文件描述符使用了 finalizer，这样在文件描述符不再被使用时，即便用户忘记了手动关闭文件描述符，在垃圾回收时也可以自动调用 finalizer 关闭文件描述符。

另外，finalizer 可以将资源的释放托管给垃圾回收，这一点在一些高级的场景（例如 CGO）中非常有用。在 Go 语言调用 C 函数时，C 函数分配的内存不受 Go 垃圾回收的管理，这时我们常常借助 defer 在函数调用结束时手动释放内存。

如下所示，在 defer 释放 C 结构体中的指针。

```go
package main

// #include <stdio.h>
// #include <malloc.h>
// typedef struct {
// char *msg;
// } Bug;
//
// void bug(Bug *b) {
// printf("%s", b->msg);
// }
import "C"
import "unsafe"

func main() {
	bug := C.Bug{C.CString("Hello, World!")}
	defer C.free(unsafe.Pointer(bug.msg))
	C.bug(&bug)
}
```

将其修改为 finalizer 的形式如下，其中 runtime.KeepAlive 保证了 finalizer 的调用只能发生在该函数之后。另外，finalizer 函数并不一定要执行实际的内存释放，可以将当前指针存储起来，由单独的协程定时释放。

```go
func main() {
	bug := C.Bug{C.CString("Hello, World!")}
	runtime.SetFinalizer(&bug, func(bug *C.Bug) {
		C.free(unsafe.Pointer(bug.msg))
	})
	C.bug(&bug)
	runtime.KeepAlive(&bug)
}
```

### 栈扫描

栈扫描是根对象扫描中最重要的部分，因为在一个程序中，可能有成千上万个协程栈。栈扫描需要编译时与运行时的共同努力，运行时能够计算出当前协程栈的所有栈帧信息，而编译时能够得知栈上有哪些指针，以及对象中的哪一部分包含了指针。运行时首先计算出栈帧布局，每个栈帧都代表一个函数，运行时可以得知当前栈帧的函数参数、函数本地变量、寄存器 SP、BP 等一系列信息。

每个栈帧函数的参数和局部变量，都需要进行扫描，确认该对象是否仍然在使用，如果在使用则需要扫描位图判断对象中是否包含指针。

`runtime/mgcmark.go`

```go
// Scan a stack frame: local variables and function arguments/results.
//
//go:nowritebarrier
func scanframeworker(frame *stkframe, state *stackScanState, gcw *gcWork) {
	if _DebugGC > 1 && frame.continpc != 0 {
		print("scanframe ", funcname(frame.fn), "\n")
	}

	isAsyncPreempt := frame.fn.valid() && frame.fn.funcID == funcID_asyncPreempt
	isDebugCall := frame.fn.valid() && frame.fn.funcID == funcID_debugCallV2
	if state.conservative || isAsyncPreempt || isDebugCall {
		if debugScanConservative {
			println("conservatively scanning function", funcname(frame.fn), "at PC", hex(frame.continpc))
		}

		// Conservatively scan the frame. Unlike the precise
		// case, this includes the outgoing argument space
		// since we may have stopped while this function was
		// setting up a call.
		//
		// TODO: We could narrow this down if the compiler
		// produced a single map per function of stack slots
		// and registers that ever contain a pointer.
		if frame.varp != 0 {
			size := frame.varp - frame.sp
			if size > 0 {
				scanConservative(frame.sp, size, nil, gcw, state)
			}
		}

		// Scan arguments to this frame.
		if frame.arglen != 0 {
			// TODO: We could pass the entry argument map
			// to narrow this down further.
			scanConservative(frame.argp, frame.arglen, nil, gcw, state)
		}

		if isAsyncPreempt || isDebugCall {
			// This function's frame contained the
			// registers for the asynchronously stopped
			// parent frame. Scan the parent
			// conservatively.
			state.conservative = true
		} else {
			// We only wanted to scan those two frames
			// conservatively. Clear the flag for future
			// frames.
			state.conservative = false
		}
		return
	}

	locals, args, objs := getStackMap(frame, &state.cache, false)

	// Scan local variables if stack frame has been allocated.
	if locals.n > 0 {
		size := uintptr(locals.n) * goarch.PtrSize
		scanblock(frame.varp-size, size, locals.bytedata, gcw, state)
	}

	// Scan arguments.
	if args.n > 0 {
		scanblock(frame.argp, uintptr(args.n)*goarch.PtrSize, args.bytedata, gcw, state)
	}

	// Add all stack objects to the stack object list.
	if frame.varp != 0 {
		// varp is 0 for defers, where there are no locals.
		// In that case, there can't be a pointer to its args, either.
		// (And all args would be scanned above anyway.)
		for i := range objs {
			obj := &objs[i]
			off := obj.off
			base := frame.varp // locals base pointer
			if off >= 0 {
				base = frame.argp // arguments and return values base pointer
			}
			ptr := base + uintptr(off)
			if ptr < frame.sp {
				// object hasn't been allocated in the frame yet.
				continue
			}
			if stackTraceDebug {
				println("stkobj at", hex(ptr), "of size", obj.size)
			}
			state.addObject(ptr, obj)
		}
	}
}
```

什么情况下对象可能没有被使用呢？如下所示，当 foo 函数执行到调用 bar 函数时，局部对象 t 就已经没有被使用了，所以即便对象中有指针，位图中仍然全为 0，因为一个不再被使用的对象，不需要再被扫描。

```go
func foo() {
    t := T()
    t.a = 2
    bar()
}
```

### 栈对象

Go 语言早期就是通过上述方式对协程栈中的对象进行扫描的。但是这种方法在有些情况下会出现问题，例如在如下函数中，对象 t 首先被变量 p 引用，但是在之后的程序中，变量 p 的值发生了变化，这意味着对象 t 其实并没有被使用。但是由于编译器难以知道变量 p 在何时会重新赋值导致对象 t 不再被引用，因此会采取保守的算法认为对象 t 仍然存在，此时如果对象 t 中有指针指向了堆内存，就会造成内存泄漏，因为这部分内存本应该被释放。

```go
t := T{...}
p := &t
for {
    if ... {
        p = ...
    }
}
```

为了解决内存泄漏问题，Go 语言引进了栈对象（stack object）的概念。栈对象是在栈上能够被寻址的对象。例如上例中的对象 t，由于其能够被& t 的形式寻址，所以其一定在栈上有地址，这样的对象 t 就被叫作栈对象。不是所有的变量都会存储在栈上，例如存储在寄存器中的变量就是不能被寻址的。

编译器会在编译时将所有的栈对象都记录下来，同时，编译器将追踪栈中所有可能指向栈对象的指针。在垃圾回收期间，所有的栈对象都会存储到一棵二叉搜索树中。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-10%20栈对象扫描.png)

如图 20-10 所示，假设 F 为一个局部变量指针，其引用了栈帧上的栈对象 E → C → D → A，那么说明栈对象 E、C、D、A 都是存活的，需要被扫描。相反，如果栈对象 B 没有被引用，并且接下来在 foo 函数中没有使用到 B 对象，那么 B 对象将不会被扫描，从而解决了内存泄漏问题。

### 扫描灰色对象

从根对象的收集来看，全局变量、析构器、所有协程的栈都会被扫描，从而标记目前还在使用的内存对象。下一步是从这些被标记为灰色的内存对象出发，进一步标记整个堆内存中活着的对象。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-11%20局部队列扫描避免使用锁.png)

如图 20-11 所示，在进行根对象扫描时，会将标记的对象放入本地队列中，如果本地队列放不下，则放入全局队列中。这种设计最大限度地避免了使用锁，在本地缓存的队列可以被逻辑处理器 P 无锁访问。

在进行扫描时，使用相同的原理，先消费本地队列中找到的标记对象，如果本地队列为空，则加锁获取全局队列中存储的对象。
在标记期间、会循环往复地从本地标记队列获取灰色对象，灰色对象扫描到的白色对象仍然会被放入标记队列中，如果扫描到已经被标记的对象则忽略，一直到队列中的任务为空为止。

`src/runtime/mgcmark.go:gcDrain`

```go
		b := gcw.tryGetFast()
		if b == 0 {
			b = gcw.tryGet()
			if b == 0 {
				// Flush the write barrier
				// buffer; this may create
				// more work.
				wbBufFlush(nil, 0)
				b = gcw.tryGet()
			}
		}
		if b == 0 {
			// Unable to get work.
			break
		}
		scanobject(b, gcw)
```

对象的扫描过程位于 scanobject 函数中。之前介绍过，堆上的任意一个指针都能找到其对象所在 span 中的位置，并且可以通过 gcmarkBits 位图检查对象是否被扫描。但现在面对的问题是需要对所有对象的内存逐个进行扫描，查看对象内存中是否含有指针，如果对象中没有存储指针，则根本不需要花时间进行检查。为了实现更快的查找，Go 语言在内存分配时记录了对象中是否包含指针等元信息。

之前介绍过，heapArena 包含整个 64MB 的 Arena 元数据。

其中有一个重要的 bitmap 字段 (bitmap)[heapArenaBitmapBytes]byte() 用位图的形式记录了每个指针大小（8byte() 的内存中的信息。每个指针大小的内存都会有 2 个 bit 分别表示当前内存是否应该继续扫描及是否包含指针。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-12%20位图记彔内存是否需要被扫描及是否包含指针.png)

如图 20-12 所示，bitmap 位图记录了内存是否需要被扫描以及是否包含指针。bitmap 中 1 个 byte 大小的空间对应了虚拟内存中 4 个指针大小的空间。bitmap 中的前 4 位为扫描位，后 4 位为指针位。分别对应指定的指针大小的空间是否需要继续进行扫描及是否包含指针。

例如，对于一个结构体 obj，当我们知道其前 2 个字段为指针，后面的字段不包含指针时，后面的字段就不再需要被扫描。因此，扫描位可以加速对象的扫描，避免扫描无用的字段。

```go
type obj struct {
    a *int
    b *int
    c int
    d int
}
```

当需要继续扫描并且发现了当前有指针时，就需要取出指针的值，并对其进行扫描，图 20-13 为通过指针查找到的对应的 span 中对象的示意图。与之前介绍的全局变量相似，可以根据指针查找到 span 中的对象，如果发现引用的是堆中的白色对象（即还没有被标记），则标记该对象（标记为灰色）并将该对象放入本地任务队列中。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-13%20指针最终查找到span中对应的对象.png)

## 标记终止阶段

完成并发标记阶段所有灰色对象的扫描和标记后进入标记终止阶段，标记终止阶段主要完成一些指标，例如统计用时、统计强制开始 GC 的次数、更新下一次触发 GC 需要达到的堆目标、关闭写屏障等，并唤醒后台清扫的协程，开始下一阶段的清扫工作。标记终止阶段会再次进入 STW。

标记终止阶段的重要任务是计算下一次触发 GC 时需要达到的堆目标，这叫作垃圾回收的调步算法。调步算法是 Go 1.5 提出的算法，由于从 Go 1.5 开始使用并发的三色标记，在 GC 开始到结束的过程中，用户协程可能被分配了大量的内存，所以在 GC 的过程中，程序占用的内存（后简称占用内存）的大小实际上超过了我们设定的触发 GC 的目标。为了解决这样的问题，需要对程序进行估计，从而在达到内存占用量目标（后简称目标内存）之前就启动 GC，并保证在 GC 结束之后，占用内存的大小刚好在目标内存附近，如图 20-14 所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-14%20调步算法中的重要指标.png)

因此，调步算法最重要的任务是估计出下一次触发 GC 的最佳时机，而这依赖本次 GC 的阶段差额—— GC 完成后占用内存与目标内存之间的差距。如果 GC 完成后占用内存远小于目标内存，则意味着触发 GC 的时间过早。如果 GC 完成后占用内存远大于目标内存，则意味着触发 GC 的时间太迟。因此调度算法的第 1 个目标是 min() |目标占用内存 -GC 完成后的占用内存|）。除此之外，调步算法还有第 2 个目标，即预计执行标记的 CPU 占用率接近 25%。结合之前提到的 25% 的后台标记协程，这个要求是满足的，在正常情况下，只有 25% 的 CPU 会执行后台标记任务。但如果用户工作协程执行了辅助标记（将在下节介绍），那么这一前提将不再成立。如果用户协程执行了过多的辅助标记，则会导致 GC 完成后的占用内存偏小，因为用户协程将本来应该用来分配内存的时间用来了执行辅助标记。算法将先计算目标内存与 GC 完成后的占用内存的偏差，

```
偏差率=(目标增长率 - 触发率)-(实际增长率 - 触发率)
```

其中

- 目标增长率=目标内存 / 上一次 GC 完成后的标记内存 -1
- 触发率=触发 GC 时的占用内存 / 上一次 GC 完成后的标记内存 -1
- 实际增长率= GC 完成后的内存 / 上一次 GC 完成后的标记内存 -1

这其实是

```
偏差=(目标内存 - 触发 GC 时的占用内存)-(GC) 完成后的占用内存 - 触发 GC 时的占用内存)
```

的变形。为了修复辅助标记带来的偏差，计算辅助标记所用的时间，从而调整 (GC) 完成后的占用内存 - 触发 GC 时的占用内存）的大小。因此最终的偏差率调整为

```
偏差率=(目标增长率 - 触发率)- 调整率x(实际增长率 - 触发率)
```

其中，

```
调整率 = GC 标记阶段的 CPU 占用率 / 目标 CPU 占用率
```

实际代码如下：

`src/runtime/mgcpacer.go`

```go
// endCycle computes the consMark estimate for the next cycle.
// userForced indicates whether the current GC cycle was forced
// by the application.
func (c *gcControllerState) endCycle(now int64, procs int, userForced bool) {
	// Record last heap goal for the scavenger.
	// We'll be updating the heap goal soon.
	gcController.lastHeapGoal = c.heapGoal()

	// Compute the duration of time for which assists were turned on.
	assistDuration := now - c.markStartTime

	// Assume background mark hit its utilization goal.
	utilization := gcBackgroundUtilization
	// Add assist utilization; avoid divide by zero.
	if assistDuration > 0 {
		utilization += float64(c.assistTime.Load()) / float64(assistDuration*int64(procs))
	}

	if c.heapLive.Load() <= c.triggered {
		// Shouldn't happen, but let's be very safe about this in case the
		// GC is somehow extremely short.
		//
		// In this case though, the only reasonable value for c.heapLive-c.triggered
		// would be 0, which isn't really all that useful, i.e. the GC was so short
		// that it didn't matter.
		//
		// Ignore this case and don't update anything.
		return
	}
	idleUtilization := 0.0
	if assistDuration > 0 {
		idleUtilization = float64(c.idleMarkTime.Load()) / float64(assistDuration*int64(procs))
	}
	// Determine the cons/mark ratio.
	//
	// The units we want for the numerator and denominator are both B / cpu-ns.
	// We get this by taking the bytes allocated or scanned, and divide by the amount of
	// CPU time it took for those operations. For allocations, that CPU time is
	//
	//    assistDuration * procs * (1 - utilization)
	//
	// Where utilization includes just background GC workers and assists. It does *not*
	// include idle GC work time, because in theory the mutator is free to take that at
	// any point.
	//
	// For scanning, that CPU time is
	//
	//    assistDuration * procs * (utilization + idleUtilization)
	//
	// In this case, we *include* idle utilization, because that is additional CPU time that
	// the GC had available to it.
	//
	// In effect, idle GC time is sort of double-counted here, but it's very weird compared
	// to other kinds of GC work, because of how fluid it is. Namely, because the mutator is
	// *always* free to take it.
	//
	// So this calculation is really:
	//     (heapLive-trigger) / (assistDuration * procs * (1-utilization)) /
	//         (scanWork) / (assistDuration * procs * (utilization+idleUtilization)
	//
	// Note that because we only care about the ratio, assistDuration and procs cancel out.
	scanWork := c.heapScanWork.Load() + c.stackScanWork.Load() + c.globalsScanWork.Load()
	currentConsMark := (float64(c.heapLive.Load()-c.triggered) * (utilization + idleUtilization)) /
		(float64(scanWork) * (1 - utilization))

	// Update cons/mark controller. The time period for this is 1 GC cycle.
	//
	// This use of a PI controller might seem strange. So, here's an explanation:
	//
	// currentConsMark represents the consMark we *should've* had to be perfectly
	// on-target for this cycle. Given that we assume the next GC will be like this
	// one in the steady-state, it stands to reason that we should just pick that
	// as our next consMark. In practice, however, currentConsMark is too noisy:
	// we're going to be wildly off-target in each GC cycle if we do that.
	//
	// What we do instead is make a long-term assumption: there is some steady-state
	// consMark value, but it's obscured by noise. By constantly shooting for this
	// noisy-but-perfect consMark value, the controller will bounce around a bit,
	// but its average behavior, in aggregate, should be less noisy and closer to
	// the true long-term consMark value, provided its tuned to be slightly overdamped.
	var ok bool
	oldConsMark := c.consMark
	c.consMark, ok = c.consMarkController.next(c.consMark, currentConsMark, 1.0)
	if !ok {
		// The error spiraled out of control. This is incredibly unlikely seeing
		// as this controller is essentially just a smoothing function, but it might
		// mean that something went very wrong with how currentConsMark was calculated.
		// Just reset consMark and keep going.
		c.consMark = 0
	}

	if debug.gcpacertrace > 0 {
		printlock()
		goal := gcGoalUtilization * 100
		print("pacer: ", int(utilization*100), "% CPU (", int(goal), " exp.) for ")
		print(c.heapScanWork.Load(), "+", c.stackScanWork.Load(), "+", c.globalsScanWork.Load(), " B work (", c.lastHeapScan+c.lastStackScan.Load()+c.globalsScan.Load(), " B exp.) ")
		live := c.heapLive.Load()
		print("in ", c.triggered, " B -> ", live, " B (∆goal ", int64(live)-int64(c.lastHeapGoal), ", cons/mark ", oldConsMark, ")")
		if !ok {
			print("[controller reset]")
		}
		println()
		printunlock()
	}
}
```

从公式中可以看出，实际增长率和辅助标记的时长都会影响最终的偏差率。目标内存与 GC 完成后的占用内存偏离越大偏差率越大。这时，下一次 GC 的触发率会渐进调整，即每次只调整偏差的一半，公式如下：

```
下次 GC 触发率=上次 GC 触发率 +1/2× 偏差率
```

计算出下次 GC 触发率后，需要计算出目标内存大小，这是在标记终止阶段的 gcSetTriggerRatio 函数中完成的，目标内存的计算如下：

```go
goal := ^uint64(0)
if gcpercent >= 0 {
    goal = memstats.heap_marked + memstats.heap_marked*uint64(gcpercent)/100
}
```

goal 为下次 GC 完成后的目标内存，其大小取决于本次 GC 扫描后的占用内存及 gcpercent 的大小。gcpercent 可以由用户动态设置，调用标准库的 SetGCPercent 函数，可以修改 gcpercent 的大小。

gcpercent 的默认值为 100，代表目标内存是上一次 GC 目标内存的 2 倍。当 gcpercent 的值小于 0 时，将禁用 Go 的垃圾回收。

另外，也可以通过在编译或运行时添加 GOGC 环境变量的方式修改 gcpercent 的大小，其核心逻辑是在程序初始化时调用 readgogc 函数实现的。例如，GOGC = off./main 将关闭 GC。

`src/runtime/mgcpacer.go`

```go
func readGOGC() int32 {
	p := gogetenv("GOGC")
	if p == "off" {
		return -1
	}
	if n, ok := atoi32(p); ok {
		return n
	}
	return 100
}
```

明确了目标内存后，触发内存的大小可以简单定义如下：

```
触发内存=触发率 × 目标内存
```

其中，触发率不能大于 0.95，也不能小于 0.6。

## 辅助标记

Go 1.5 引入了并发标记后，带来了许多新的问题。例如，在并发标记阶段，扫描内存的同时用户协程也不断被分配内存，当用户协程的内存分配速度快到后台标记协程来不及扫描时，GC 标记阶段将永远不会结束，从而无法完成完整的 GC 周期，造成内存泄漏。

为了解决这样的问题，引入辅助标记算法。辅助标记必须在垃圾回收的标记阶段进行，由于用户协程被分配了超过限度的内存而不得不将其暂停并切换到辅助标记工作。所以一个简单的策略是让 X = M，其中，X 为后台标记协程需要多扫描的内存，M 为新分配的内存。即在并发标记期间，一旦新分配了内存 M，就必须完成 M 的扫描工作。我们之前看到过，对于 obj 这样的对象，并不需要扫描对象中所有的内存。

```go
type obj struct {
    a *int
    b *int
    c int
    d int
}
```

因此扫描策略可以调整为 X = assistWorkPerByte×M

其中，assistWorkPerByte < 1，代表每字节需要完成多少扫描工作，并且真实需要扫描的内存会少于实际的内存。

在 GC 并发标记阶段，当用户协程分配内存时，会先检查是否已经完成了指定的扫描工作。当前协程中的 gcAssistBytes 字段代表当前协程可以被分配的内存大小，类似资产池。当本地的资产池不足时（即 gcAssistBytes<0），会尝试从全局的资产池中获取。用户协程一开始是没有资产的，所有的资产都来自后台标记协程。

`src/runtime/malloc.go`

```go
// Allocate an object of size bytes.
// Small objects are allocated from the per-P cache's free lists.
// Large objects (> 32 kB) are allocated straight from the heap.
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
...
// assistG is the G to charge for this allocation, or nil if
	// GC is not currently active.
	var assistG *g
	if gcBlackenEnabled != 0 {
		// Charge the current user G for this allocation.
		assistG = getg()
		if assistG.m.curg != nil {
			assistG = assistG.m.curg
		}
		// Charge the allocation against the G. We'll account
		// for internal fragmentation at the end of mallocgc.
		assistG.gcAssistBytes -= int64(size)

		if assistG.gcAssistBytes < 0 {
			// This G is in debt. Assist the GC to correct
			// this before allocating. This must happen
			// before disabling preemption.
			gcAssistAlloc(assistG)
		}
	}
...
}
```

用户协程中的本地资产来自后台标记协程的扫描工作。后台标记协程的扫描工作会增加全局资产池的大小。之前提到，X = assistWorkPerByte×M。反过来，如果标记协程已经扫描完成的内存为 X，那么意味着全局资产池可以容忍用户协程分配的内存数量为 M = X/assistWorkPerByte。这种机制保证了在 GC 并发标记时，工作协程分配的内存数量不至于过多，也不会太少。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-15%20全局资产池与本地资产池的协调.png)

如果工作协程在分配内存时，既无法从本地资产池也无法从全局资产池获取资产，那么需要停止工作协程，并执行辅助标记协程。辅助标记协程需要额外扫描的内存大小为 assistWorkPerByte×M，当扫描完成指定工作或被抢占时会退出。当辅助标记完成后，如果本地仍然没有足够的资产，则可能是因为当前协程被抢占，也可能是因为当前逻辑处理器的工作池中没有多余的标记工作。当协程被抢占时，会调用 Gosched 函数让渡当前辅助标记的执行权利，而如果当前逻辑处理器的工作池中没有多余的标记工作可做，则会陷入休眠状态，当后台标记协程扫描了足够的任务后，会刷新全局资产池并将等待中的协程唤醒，如图 20-16 所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-16%20唤醒等待中的用户协程.png)

## 屏障技术

在并发标记中，标记协程与用户协程共同工作的模式带来了很多难题。如果说辅助标记解决的是垃圾回收正常结束与循环的问题，那么屏障技术将解决更棘手的问题——准确性。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-17%20并収标记的准确性问题.png)

如图20-17所示，假设在垃圾回收已经扫描完根对象（此时根对象为黑色）并继续扫描期间，白色对象Z正被一个灰色对象引用，但此时工作协程在执行过程中，让黑色的根对象指向了白色的对象Z。由于黑色的对象不会被扫描，这将导致白色对象Z被视为垃圾对象最终被回收。

那么是不是黑色对象一定不能指向白色对象呢？其实也不一定。如图20-18所示，即便黑色对象引用了白色对象，但只要白色对象中有一条路径始终被灰色对象引用了，此白色对象就一定能被扫描到。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-18%20只要白色对象始终被灰色对象引用就能被扫描.png)

这其实引出了保证并发标记准确性需要遵守的原则，即强、弱三色不变性。强三色不变性指所有白色对象都不能被黑色对象引用，这是一种比较严格的要求。与之对应的是弱三色不变性，弱三色不变性允许白色对象被黑色对象引用，但是白色对象必须有一条路径始终是被灰色对象引用的，这保证了该对象最终能被扫描到。
在并发标记写入和删除对象时，可能破坏三色不变性，因此必须有一种机制能够维护三色不变性，这就是屏障技术。屏障技术的原则是在写入或者删除对象时将可能活着的对象标记为灰色。上例如果能够在对象写入时将Z对象设置为灰色，那么Z对象最终将被扫描到，如图20-19所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-19%20屏障技术保证准确性.png)

上图提到的这种简单的屏障技术是Dijkstra风格的插入屏障，其实现形式如下，如果目标对象src为黑色，则将新引用的对象标记为灰色。

```go
func writeBarrier(src, dst *object) {
    if src.color == black {
        dst.color = gray
    }
}
```

还有一种常见的策略是在删除引用时做文章，如图20-20所示，Yuasa删除写屏障在对象被解除引用后，会立即将原引用对象标记为灰色。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-20%20删除屏障在取消引用时变为灰色.png)

这样即便没有写屏障，在插入操作时也不会破坏三色不变性，如图20-21所示，但是Z对象可能是垃圾对象。
插入屏障与删除屏障通过在写入和删除时重新标记颜色保证了三色不变性，解决了并发标记期间的准确性问题，但是它们都存在浮动垃圾的问题。插入屏障在删除引用时，可能标记一个已经变成垃圾的对象。而删除屏障在删除引用时可能把一个垃圾对象标记为灰色。这些是垃圾回收的精度问题，不会影响其准确性，因为浮动垃圾会在下一次垃圾回收中被回收。

插入屏障与删除屏障独立存在并能良好工作的前提是并发标记期间所有的写入操作都应用了屏障技术，但现实情况不会如此。大多数垃圾回收语言不会对栈上的操作或寄存器上的操作应用屏障技术，这是因为栈上操作是最频繁的，如果每个写入或删除操作都应用屏障技术则会大大减慢程序的速度。在Go 1.8之前，尽管使用了插入屏障，但是仍然需要在标记终止期间STW阶段重新扫描根对象，来保证三色标记的一致性。为了解决重复扫描的问题，Go 1.8之后使用了混合写屏障技术，结合了Dijkstra与Yuasa两种风格。
为了了解使用混合写屏障技术的原因，我们先来看一看单纯地插入屏障和删除屏障在现实中面临的困境。假设栈上初始状态如图20-22所示，栈上变量p指向堆区内存，如果现在垃圾回收扫描完了根对象，那么old变量是不会被扫描的。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-21%20删除屏障能够保证插入时的准确性.png)

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-22%20插入屏障案例——栈上初始状态的对象引用情况.png)

进入到并发标记阶段之后，假设并发标记阶段如图20-23所示，old对象引用了p.x，但是赋值给栈上的变量不会经过写屏障。如果下一步p.x引用了一个新的内存对象k，并把k标记为灰色，但是并不把原始对象标记为灰色，那么这时原始对象即便被栈上的对象old标记也无法被扫描到。所以，必须在p.x=&k时应用删除屏障，在取消引用时，将p.x的原值标记为灰色。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-23%20插入屏障案例——并収标记阶段的对象引用情况.png)

如果只有删除屏障而没有写屏障，那么也会面临问题。假设根对象未开始扫描，对象全为白色，栈上变量p引用堆区对象o，栈上变量a引用堆区对象k，在并发标记期间，扫描完变量p还未扫描变量a时的情形如图20-24所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-24%20删除屏障案例——栈上初始状态的对象引用情况.png)

此时，工作协程将变量a置为nil，p.x=&k将对象p指向了k。如果只存在删除屏障而不启用写屏障（不标记新的k值），那么会违背三色不变性，让黑色对象引用白色对象。导致k无法被标记，如图20-25所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-25%20单独的删除屏障违背三色不变性.png)

因此，要想在标记终止阶段不用重新扫描根对象，需要使用写屏障与删除屏障混合的屏障技术，其伪代码如下：

```go
writePointer(slot, ptr):
shade(*slot)
shade(ptr)
*slot = ptr
```

在Go语言中，混合写屏障技术的实现依赖编译时与运行时的共同努力。在标记准备阶段的STW阶段会打开写屏障，具体做法是将全局变量writeBarrier.enabled设置为true。

```go
// The compiler knows about this variable.
// If you change it, you must change builtin/runtime.go, too.
// If you change the first four bytes, you must also change the write
// barrier insertion code.
var writeBarrier struct {
	enabled bool    // compiler emits a check of this before calling write barrier
	pad     [3]byte // compiler uses 32-bit load for "enabled" field
	needed  bool    // whether we need a write barrier for current GC phase
	cgo     bool    // whether we need a write barrier for a cgo check
	alignme uint64  // guarantee alignment so that compiler can use a 32 or 64-bit load
}
```

编译器会在所有堆写入或删除操作前判断当前是否为垃圾回收标记阶段，如果是则会执行对应的混合写屏障标记对象。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-26%20写屏障指针缓存池.png)

Go语言中构建了如图20-26所示的写屏障指针缓存池，gcWriteBarrier先将所有被标记的指针放入缓存池中，并在容量满后，一次性全部刷新到扫描任务池中。最终这些被标记的指针都将被扫描。

## 垃圾清扫

垃圾标记工作完成意味着已经追踪到内存中所有活着的对象（虽然可能有一些浮动垃圾），之后进入垃圾清扫阶段，将垃圾对象的内存回收重用或返还给操作系统。
在标记结束阶段会调用gcSweep函数，该函数会将sweep.g清扫协程的状态变为running，在结束STW阶段并开始重新调度循环时优先清扫协程。

`src/runtime/mgc.go`

```go
// gcSweep must be called on the system stack because it acquires the heap
// lock. See mheap for details.
//
// The world must be stopped.
//
//go:systemstack
func gcSweep(mode gcMode) {
	assertWorldStopped()

	if gcphase != _GCoff {
		throw("gcSweep being done but phase is not GCoff")
	}

	lock(&mheap_.lock)
	mheap_.sweepgen += 2
	sweep.active.reset()
	mheap_.pagesSwept.Store(0)
	mheap_.sweepArenas = mheap_.allArenas
	mheap_.reclaimIndex.Store(0)
	mheap_.reclaimCredit.Store(0)
	unlock(&mheap_.lock)

	sweep.centralIndex.clear()

	if !_ConcurrentSweep || mode == gcForceBlockMode {
		// Special case synchronous sweep.
		// Record that no proportional sweeping has to happen.
		lock(&mheap_.lock)
		mheap_.sweepPagesPerByte = 0
		unlock(&mheap_.lock)
		// Sweep all spans eagerly.
		for sweepone() != ^uintptr(0) {
			sweep.npausesweep++
		}
		// Free workbufs eagerly.
		prepareFreeWorkbufs()
		for freeSomeWbufs(false) {
		}
		// All "free" events for this mark/sweep cycle have
		// now happened, so we can make this profile cycle
		// available immediately.
		mProf_NextCycle()
		mProf_Flush()
		return
	}

	// Background sweep.
	lock(&sweep.lock)
	if sweep.parked {
		sweep.parked = false
		ready(sweep.g, 0, true)
	}
	unlock(&sweep.lock)
}
```

清扫阶段在程序启动时调用的gcenable函数中启动。注意，程序中只有一个垃圾清扫协程，并在清扫阶段与用户协程同时运行。

```go
// gcenable is called after the bulk of the runtime initialization,
// just before we're about to start letting user code run.
// It kicks off the background sweeper goroutine, the background
// scavenger goroutine, and enables GC.
func gcenable() {
	// Kick off sweeping and scavenging.
	c := make(chan int, 2)
	go bgsweep(c)
	go bgscavenge(c)
	<-c
	<-c
	memstats.enablegc = true // now that runtime is initialized, GC is okay
}
```

当清扫协程被唤醒后，会开始垃圾清扫。垃圾清扫采取了懒清扫的策略，即执行少量清扫工作后，通过Gosched函数让渡自己的执行权利，不需要一直执行。因此当触发下一阶段的垃圾回收后，可能有没有被清理的内存，需要先将它们清理完。

```go
func bgsweep(c chan int) {
	sweep.g = getg()

	lockInit(&sweep.lock, lockRankSweep)
	lock(&sweep.lock)
	sweep.parked = true
	c <- 1
	goparkunlock(&sweep.lock, waitReasonGCSweepWait, traceEvGoBlock, 1)

	for {
		// bgsweep attempts to be a "low priority" goroutine by intentionally
		// yielding time. It's OK if it doesn't run, because goroutines allocating
		// memory will sweep and ensure that all spans are swept before the next
		// GC cycle. We really only want to run when we're idle.
		//
		// However, calling Gosched after each span swept produces a tremendous
		// amount of tracing events, sometimes up to 50% of events in a trace. It's
		// also inefficient to call into the scheduler so much because sweeping a
		// single span is in general a very fast operation, taking as little as 30 ns
		// on modern hardware. (See #54767.)
		//
		// As a result, bgsweep sweeps in batches, and only calls into the scheduler
		// at the end of every batch. Furthermore, it only yields its time if there
		// isn't spare idle time available on other cores. If there's available idle
		// time, helping to sweep can reduce allocation latencies by getting ahead of
		// the proportional sweeper and having spans ready to go for allocation.
		const sweepBatchSize = 10
		nSwept := 0
		for sweepone() != ^uintptr(0) {
			sweep.nbgsweep++
			nSwept++
			if nSwept%sweepBatchSize == 0 {
				goschedIfBusy()
			}
		}
		for freeSomeWbufs(true) {
			// N.B. freeSomeWbufs is already batched internally.
			goschedIfBusy()
		}
		lock(&sweep.lock)
		if !isSweepDone() {
			// This can happen if a GC runs between
			// gosweepone returning ^0 above
			// and the lock being acquired.
			unlock(&sweep.lock)
			continue
		}
		sweep.parked = true
		goparkunlock(&sweep.lock, waitReasonGCSweepWait, traceEvGoBlock, 1)
	}
}
```

### 懒清扫逻辑

清扫是以span为单位进行的，sweepone函数的作用是找到一个span并进行相应的清扫工作。先从mheap中的sweepSpans队列中取出需要清扫的span。

sweepSpans数组的长度为2，sweepSpans[sweepgen/2%2]保存当前正在使用的span列表，sweepSpans[1-sweepgen/2%2]保存等待清扫的span列表，由于sweepgen每次清扫时加2。

因此sweepSpans [0]、sweepSpans [1]每次清扫时互相交换身份，即本次正在使用的span列表将是下一次GC待清扫的列表。
在清扫span期间，最重要的一步是将gcmarkBits位图赋值给allocBits位图，如图20-27所示。

当前gcmarkBits是GC标记后最新的对象位图，当gcmarkBits中的bit位为1时，代表当前对象是活着的。所以，当gcmarkBits中的某一个bit位为1，但是对应的allocBits位图中的bit位为0时，代表这个对象是会被回收的垃圾对象。完成这一切换后，就可以通过位图使用已经是垃圾对象的内存了。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-27%20将gcmarkBits位图赋值给allocBits位图.png)

如果GC后gcmarkBits的全部bit位都为0，那么意味着当前所有span中的对象都不会再被其他对象引用（大对象比较特殊，因为其span内部只有一个对象）。
这时，整个span将会被mheap回收，并更新整个基数树（参考第18章），表明当前span的整个空间都可以被程序再次使用。如果当前span的整个空间并不完全为空，那么span会被重新放入sweepSpans正在使用的span列表中。

可以看出，这种回收方式并没有直接将内存释放到操作系统中，而是再次组织内存以便能在下次内存分配时利用已经被回收的内存。

### 辅助清扫

我们已经知道，清扫是通过懒清扫的形式进行的，因此，在下次触发GC时，必须将上一次GC未清扫的span全部扫描。如果剩余的未清扫span太多，那么将大大拖后下一次GC开始的时间。为了规避这一问题，Go语言使用了辅助清扫的手段，这是在Go 1.5之后，和并发GC同时推出的。
辅助清扫即工作协程必须在适当的时机执行辅助清扫工作，以避免下一次GC发生时还有大量的未清扫span。判断是否需要清扫的最好时机是在工作协程分配内存时。
目前Go语言会在两个时机判断是否需要辅助扫描。一个是在需要向mcentrel申请内存时，一个在是大对象分配时。在这两个时间会判断当前已经清扫的page数大于清理的目标page数这个条件是否成立，如果不成立则会进行辅助清扫直到条件成立。sweepPagesPerByte是一个重要的参数，其代表了工作协程每分配1 byte需要辅助清扫的page数，是一个比率。

可以看出，辅助标记策略会尽可能地保证在下次触发GC时，已经扫描了所有待扫描的span。

## 系统驻留内存清除

驻留内存（RSS）是主内存（RAM）保留的进程占用的内存部分，是从操作系统中分配的内存。
为了将系统分配的内存保持在适当的大小，同时回收不再被使用的内存，Go语言使用了单独的后台清扫协程来清除内存。后台清扫协程目前是在程序开始时启动的，并且只启动一个。

清除策略占用当前线程CPU 1%的时间进行清除，因此，在大部分时间里，该协程处于休眠状态。bgscavenge花了很大的精力来计算和调整时间，以保证实现1%CPU执行时间的目标。因此如果清除花费的时间太多，那么休眠的时间也必须相应增加。

```go
// Background scavenger.
//
// The background scavenger maintains the RSS of the application below
// the line described by the proportional scavenging statistics in
// the mheap struct.
func bgscavenge(c chan int) {
	scavenger.init()

	c <- 1
	scavenger.park()

	for {
		released, workTime := scavenger.run()
		if released == 0 {
			scavenger.park()
			continue
		}
		atomic.Xadduintptr(&mheap_.pages.scav.released, released)
		scavenger.sleep(workTime)
	}
}
```

一次只清除一个物理页。scavengeOne包含清扫的核心逻辑，其基本思路是在基数树中找到连续的没有被操作系统回收的内存，我们在介绍内存分配时提到过，基数树的叶子节点管理了一个chunk块大小的内存。
对于每个chunk，都会有位图pallocBits管理其中每个page的内存分配。之前没有提到的是，scavenged是一个额外的位图，每一位与page的对应方式和分配位图pallocBits相似，但是含义不同。当bit位为1时，代表当前page已经被操作系统回收，因此当pallocBits中的某一位为1时，其对应的scavenged位必定为0。同时，只有当pallocBits与scavenged对应的位同时为0时，才表明其对应的page可以被清扫。

位图中每个bit位管理的page是固定的8KB，但是释放回操作系统中的内存至少为一个物理页大小。因此，实际可能需要释放n个page，即需要找到位图中连续可用的n个bit位。和分配时辅助查找的searchAddr字段一样，有一个辅助清扫的scavAddr字段，在系统驻留内存清扫时会从scavAddr之后进行搜索，而忽视掉scavAddr之前的地址，如图20-28所示。

![](../../../assets/images/docs/internal/gc/underlying_principle/图20-28%20scavAddr用于辅助系统驻留内存清除.png)

在开始清除搜索时，会查找searchAddr所在的chunk块中是否存在即空闲又没有被清除的连续空间，如果查找不到，则通过基数树从上到下进行扫描，找到符合要求的区域。
当查找到满足要求的连续空间后，就将scavenged位图的相应位置设置为1，更新scavAddr地址，将内存归还给操作系统，并更新相应的统计。

从位图中快速查找符合要求的区域的核心代码位于fillAligned，假如没有连续的空间，那么fillAligned全为1，如此即可快速判断是否找到该区域。

```go

```
