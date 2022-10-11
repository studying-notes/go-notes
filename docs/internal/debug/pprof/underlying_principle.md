---
date: 2022-10-11T14:21:10+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "pprof 底层原理"  # 文章标题
url:  "posts/go/docs/internal/debug/pprof/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

pprof 分为采样和分析两个阶段。采样指一段时间内某种类型的样本数据，pprof 并不会像 trace 一样记录每个事件，因此其相对于 trace 收集到的文件要小得多。

## 堆内存样本

对堆内存采样时，并不是每次调用 mallocgc 分配堆内存都会被记录下来，这里有一个指标—— MemProfileRate，当多次内存分配累积到该指标以上时，才记录一次。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-9%20pprof样本存储到链表中.png)

记录下来的每个样本都是一个 bucket，如图 21-9 所示，该 bucket 会存储到全局 mbuckets 链表中，mbuckets 链表中的对象不会被 GC 扫描，因为它加入了 span 中的 special 序列。bucket 中保留的重要数据除了当前分配的内存大小，还包括当前哪一个函数触发了内存分配以及该函数的调用链，这借助栈追踪实现。有了这一数据，才能实现 list、tree、top 等命令。

不需要将每个样本都记录为一个 bucket，如果栈追踪后发现当前样本有相同的调用链，那么不用重复记录，直接在之前的 bucket 上加上对应的内存大小即可。为了实现这样的功能，使用简单的哈希表存储调用链上的指针，如图 21-10 所示，Go 语言对栈调用链上的指针进行哈希，并采用简单的拉链法解决哈希冲突。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-10%20使用哈希表存储调用链上的指针.png)

如果找到了相同的 bucket() 即函数调用链上的指针都是相同的），那么只需要增加该 bucket 中的内存。如果在当前哈希表中没有查找到相同的 bucket，则不仅需要创建新的 bucket，还需要将该 bucket 记录到此哈希表中，方便下次查找。

在函数栈扫描过程中，需要根据函数调用链上的指针计算出所有的栈帧信息并保存。Go 语言在此处做了一些合理的优化，例如对于调用链 A() → B() → C()，如果有协程调用了函数 A，那么它一定调用了函数 B、函数 C，因此不需要重新扫描函数 B 和函数 C，也不再需要记录函数 B 和函数 C 的行号、文件名等原始数据。这些缓存信息存储在了 locs 中，locs 是一个 map 结构，key 代表标识函数的指针（例如 A），而对应的 value 存储了该函数调用链上的其他指针（例如 B、C）。

`src/runtime/pprof/proto.go`

```go
// A profileBuilder writes a profile incrementally from a
// stream of profile samples delivered by the runtime.
type profileBuilder struct {
	start      time.Time
	end        time.Time
	havePeriod bool
	period     int64
	m          profMap

	// encoding state
	w         io.Writer
	zw        *gzip.Writer
	pb        protobuf
	strings   []string
	stringMap map[string]int
	locs      map[uintptr]locInfo // list of locInfo starting with the given PC.
	funcs     map[string]int      // Package path-qualified function name to Function.ID
	mem       []memMap
	deck      pcDeck
}
```

```go
type locInfo struct {
	// location id assigned by the profileBuilder
	id uint64

	// sequence of PCs, including the fake PCs returned by the traceback
	// to represent inlined functions
	// https://github.com/golang/go/blob/d6f2f833c93a41ec1c68e49804b8387a06b131c5/src/runtime/traceback.go#L347-L368
	pcs []uintptr

	// firstPCFrames and firstPCSymbolizeResult hold the results of the
	// allFrames call for the first (leaf-most) PC this locInfo represents
	firstPCFrames          []runtime.Frame
	firstPCSymbolizeResult symbolizeFlag
}
```

## 协程栈样本收集原理

获取协程栈样本数据和获取堆内存样本非常类似。因为他们都需要保存样本的栈追踪数据，并且使用了和堆内存样本收集相同的缓存优化手段。不太相同的是，堆内存样本关注的是函数分配的内存大小，而协程栈关注的是当前有多少协程，以及大部分协程正在执行哪个函数。另外，堆内存样本不会 STW，但每次获取协程栈样本都需要启动 STW，以获取当前所有协程的快照，pprof 协程栈样本获取流程如图 21-11 所示。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-11%20pprof协程栈样本获取流程.png)

## CPU 样本收集原理

CPU profile 分析的能力令人惊讶，其可以得到在某段时间内每个函数执行的时间，而不必修改原始程序。这是如何实现的呢？其实，和调度器的抢占类似，这需要借助程序中断的功能为分析和调试提供时机，在类 UNIX 操作系统中，会通过调用操作系统库函数 setitimer 实现。setitimer 将按照设定好的频率中断当前程序，并进入操作系统内核处理中断事件，这显然进行了线程的上下文切换。操作系统从内核态返回用户态，进入之前注册好的信号处理函数，从而为分析提供时机。

图 21-12 展示了类 UNIX 操作系统信号处理的一般流程。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-12%20类UNIX操作系统信号处理的一般流程.png)

当调用 pprof 获取 CPU 样本接口时，程序会为 setitimer 函数设置中断频率为 100Hz，即每秒中断 100 次。这是深思熟虑的选择，由于中断也会花费时间成本，所以中断的频率不可过高。另外，中断的频率也不可过低，否则我们将无法准确地计算出函数花费的时间。

调用 setitimer 函数时，中断的信号为 ITIMER_PROF。当内核态返回到用户态调用注册好的 sighandler 函数，sighandler 函数识别到信号为 _SIGPROF 时，执行 sigprof 函数记录该 CPU 样本。

```go
	if sig == _SIGPROF {
		// Some platforms (Linux) have per-thread timers, which we use in
		// combination with the process-wide timer. Avoid double-counting.
		if !delayedSignal && validSIGPROF(mp, c) {
			sigprof(c.sigpc(), c.sigsp(), c.siglr(), gp, mp)
		}
		return
	}
```

Go 语言处理中断信号的具体执行流程如图 21-13 所示，在处理过程中调用了 sigprof 函数。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-13%20Go语言处理中断信号的流程.png)

sigprof 的核心功能是记录当前的栈追踪，其实现如下所示。

`src/runtime/proc.go`

```go
// Called if we receive a SIGPROF signal.
// Called by the signal handler, may run during STW.
//
//go:nowritebarrierrec
func sigprof(pc, sp, lr uintptr, gp *g, mp *m) {
	if prof.hz.Load() == 0 {
		return
	}

	// If mp.profilehz is 0, then profiling is not enabled for this thread.
	// We must check this to avoid a deadlock between setcpuprofilerate
	// and the call to cpuprof.add, below.
	if mp != nil && mp.profilehz == 0 {
		return
	}

	// On mips{,le}/arm, 64bit atomics are emulated with spinlocks, in
	// runtime/internal/atomic. If SIGPROF arrives while the program is inside
	// the critical section, it creates a deadlock (when writing the sample).
	// As a workaround, create a counter of SIGPROFs while in critical section
	// to store the count, and pass it to sigprof.add() later when SIGPROF is
	// received from somewhere else (with _LostSIGPROFDuringAtomic64 as pc).
	if GOARCH == "mips" || GOARCH == "mipsle" || GOARCH == "arm" {
		if f := findfunc(pc); f.valid() {
			if hasPrefix(funcname(f), "runtime/internal/atomic") {
				cpuprof.lostAtomic++
				return
			}
		}
		if GOARCH == "arm" && goarm < 7 && GOOS == "linux" && pc&0xffff0000 == 0xffff0000 {
			// runtime/internal/atomic functions call into kernel
			// helpers on arm < 7. See
			// runtime/internal/atomic/sys_linux_arm.s.
			cpuprof.lostAtomic++
			return
		}
	}

	// Profiling runs concurrently with GC, so it must not allocate.
	// Set a trap in case the code does allocate.
	// Note that on windows, one thread takes profiles of all the
	// other threads, so mp is usually not getg().m.
	// In fact mp may not even be stopped.
	// See golang.org/issue/17165.
	getg().m.mallocing++

	var stk [maxCPUProfStack]uintptr
	n := 0
	if mp.ncgo > 0 && mp.curg != nil && mp.curg.syscallpc != 0 && mp.curg.syscallsp != 0 {
		cgoOff := 0
		// Check cgoCallersUse to make sure that we are not
		// interrupting other code that is fiddling with
		// cgoCallers.  We are running in a signal handler
		// with all signals blocked, so we don't have to worry
		// about any other code interrupting us.
		if mp.cgoCallersUse.Load() == 0 && mp.cgoCallers != nil && mp.cgoCallers[0] != 0 {
			for cgoOff < len(mp.cgoCallers) && mp.cgoCallers[cgoOff] != 0 {
				cgoOff++
			}
			copy(stk[:], mp.cgoCallers[:cgoOff])
			mp.cgoCallers[0] = 0
		}

		// Collect Go stack that leads to the cgo call.
		n = gentraceback(mp.curg.syscallpc, mp.curg.syscallsp, 0, mp.curg, 0, &stk[cgoOff], len(stk)-cgoOff, nil, nil, 0)
		if n > 0 {
			n += cgoOff
		}
	} else {
		n = gentraceback(pc, sp, lr, gp, 0, &stk[0], len(stk), nil, nil, _TraceTrap|_TraceJumpStack)
	}

	if n <= 0 {
		// Normal traceback is impossible or has failed.
		// See if it falls into several common cases.
		n = 0
		if usesLibcall() && mp.libcallg != 0 && mp.libcallpc != 0 && mp.libcallsp != 0 {
			// Libcall, i.e. runtime syscall on windows.
			// Collect Go stack that leads to the call.
			n = gentraceback(mp.libcallpc, mp.libcallsp, 0, mp.libcallg.ptr(), 0, &stk[0], len(stk), nil, nil, 0)
		}
		if n == 0 && mp != nil && mp.vdsoSP != 0 {
			n = gentraceback(mp.vdsoPC, mp.vdsoSP, 0, gp, 0, &stk[0], len(stk), nil, nil, _TraceTrap|_TraceJumpStack)
		}
		if n == 0 {
			// If all of the above has failed, account it against abstract "System" or "GC".
			n = 2
			if inVDSOPage(pc) {
				pc = abi.FuncPCABIInternal(_VDSO) + sys.PCQuantum
			} else if pc > firstmoduledata.etext {
				// "ExternalCode" is better than "etext".
				pc = abi.FuncPCABIInternal(_ExternalCode) + sys.PCQuantum
			}
			stk[0] = pc
			if mp.preemptoff != "" {
				stk[1] = abi.FuncPCABIInternal(_GC) + sys.PCQuantum
			} else {
				stk[1] = abi.FuncPCABIInternal(_System) + sys.PCQuantum
			}
		}
	}

	if prof.hz.Load() != 0 {
		// Note: it can happen on Windows that we interrupted a system thread
		// with no g, so gp could nil. The other nil checks are done out of
		// caution, but not expected to be nil in practice.
		var tagPtr *unsafe.Pointer
		if gp != nil && gp.m != nil && gp.m.curg != nil {
			tagPtr = &gp.m.curg.labels
		}
		cpuprof.add(tagPtr, stk[:n])

		gprof := gp
		var pp *p
		if gp != nil && gp.m != nil {
			if gp.m.curg != nil {
				gprof = gp.m.curg
			}
			pp = gp.m.p.ptr()
		}
		traceCPUSample(gprof, pp, stk[:n])
	}
	getg().m.mallocing--
}
```

添加的 CPU 样本会写入叫作 data 的 buf 中，每个样本都包含该样本的长度、时间戳、hdrsize、栈追踪指针。hdrsize 和 hz 有关，用于计算持续时间，由于中断周期固定为 100Hz，所以当前的 hdrsize 也固定为 1。最后的空间的长度是可变的，存储了栈追踪指针。

添加样本时所有数据都会被写入 data 缓存，同时会有专门的协程用于获取 data 中的数据，在读取样本的过程中，记录 data 中的读取位置 r 和写入位置 w，因此 w-r 表明当前可以读取的样本数量，如图 21-14 所示。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-14%20读取与写入位于data缓存中的CPU样本.png)

## 分析原理

所有 pprof 的样本数据最后都会以 Protocol Buffers 格式序列化数据并通过 gzip 压缩后写入文件。用户获取该文件后最终将使用 go tool pprof 对样本文件进行解析。go tool pprof 将文件解码并还原为 Protocol Buffers 格式，如下。Profile 代表一系列样本的集合，主要包含样本类型、样本数组 Sample，以及表示函数、行号、文件名等调试信息的 Location 字段。

```go
message Profile {
    repeated ValueType sample_type = 1;
    repeated Sample sample = 2;
    repeated Location location = 3;
    repeated Function function = 4;
    repeated Mapping mapping = 5;
}
```

每个 Sample 样本都对应一个 Location id 数组，代表函数调用链上的函数信息，正如之前讲到的，样本中的函数有可能重复，而每个 Location id 对应一个函数可以避免记录重复的信息。value 是一个数组，代表调用链上函数对应的值，该值和样本的类型有关，例如 CPU 样本是持续时间，而内存样本是内存的大小。

```go
message Sample {
    repeated int64 location_id = 1;
    repeated int64 value = 2;
    repeated Label label = 3;
}
```

pprof 的重要功能是统计搜集到的样本，包括 flat 与 cum 这两个重要的指标，除此之外，还能够以图的形式表示函数的调用链，相同的函数是图中同一个节点，图中的调用关系由两部分决定——父函数和子函数，其中箭头的方向表示父函数调用子函数。

下面举例分析内存分配特征文件的分析原理，图 21-15 所示为单样本，假设其函数调用链为 A() → B() → C()，所有函数都会被分配内存，对应的 flat 值为 A、B、C。当计算 cum 的值时，需要从 A 函数开始从上到下遍历调用链，当遍历到节点 B() 时，需要对其父节点 A() 的 cum 字段加上当前 flat 的值 B，当遍历到叶子节点 C() 时，父节点 B() 与 A() 的 cum 字段都需要加上当前 flat 的值 C。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-15%20单样本pprof分析原理.png)

再来看多样本的情况，我们假设其函数调用链为 D() → B() → C()，相同的调用链将会合并，其 pprof 分析原理如图 21-16 所示。

![](../../../../assets/images/docs/internal/debug/pprof/underlying_principle/图21-16%20多样本pprof分析原理.png)

```go

```
