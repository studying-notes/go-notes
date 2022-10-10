---
date: 2022-10-09T09:23:09+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "协程初探"  # 文章标题
url:  "posts/go/docs/internal/goroutine/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [进程与线程](#进程与线程)
- [线程上下文切换](#线程上下文切换)
- [线程与协程](#线程与协程)
	- [调度方式](#调度方式)
	- [上下文切换的速度](#上下文切换的速度)
	- [调度策略](#调度策略)
	- [栈的大小](#栈的大小)
- [并发与并行](#并发与并行)
- [简单协程入门](#简单协程入门)
	- [主协程与子协程](#主协程与子协程)
	- [WaitGroup](#waitgroup)
- [GMP 模型](#gmp-模型)
- [常见并发模型](#常见并发模型)

## 进程与线程

线程是可以由调度程序独立管理的最小程序指令集，而进程是程序运行的实例。

进程是资源分配的最小单位，线程是程序执行的最小单位。

在大多数情况下，线程是进程的组成部分。如图 14-1 所示，一个进程中可以存在多个线程，这些线程并发执行并共享进程的内存（例如全局变量）等资源。而进程之间相对独立，不同进程具有不同的内存地址空间、代表程序运行的机器码、进程状态、操作系统资源描述符等。

![](../../../assets/images/docs/internal/goroutine/README/图14-1%20线程与进程的区别.png)

在一个进程内部，可能有多个线程被同时处理。追求高并发处理、高性能的程序或者库一般都会设计为多线程。

那为什么程序通常不采取多进程，而采取多线程的方式进行设计呢？这是因为开启一个新进程的开销要比开启一个新线程大得多，而且**进程具有独立的内存空间**，这使得多进程之间的共享通信更加困难。

操作系统调度到 CPU 中执行的最小单位是线程。在传统的单核（Core）CPU 上运行的多线程应用程序必须交织线程，交替抢占 CPU 的时间片，如图 14-2 所示。

![](../../../assets/images/docs/internal/goroutine/README/图14-2%20单核处理器与多核处理器的区别.png)

但是，现代计算机系统普遍拥有多核处理器。在多核 CPU 上，线程可以分布在多个 CPU 核心上，从而实现真正的并行处理。

## 线程上下文切换

虽然多核处理器可以保证并行计算，但是实际中程序的数量以及实际运行的线程数量会比 CPU 核心数多得多。因此，为了平衡每个线程能够被 CPU 处理的时间并最大化利用 CPU 资源，操作系统需要在适当的时间通过定时器中断（Timer Interrupt）、I/O 设备中断、系统调用时执行上下文切换（Context Switch）。

![](../../../assets/images/docs/internal/goroutine/README/图14-3%20线程上下文切换.png)

如图 14-3 所示，当发生线程上下文切换时，需要从操作系统用户态转移到内核态，记录上一个线程的重要寄存器值（例如栈寄存器 SP）、进程状态等信息，这些信息存储在操作系统线程控制块（Thread Control Block）中。当切换到下一个要执行的线程时，需要加载重要的 CPU 寄存器值，并从内核态转移到操作系统用户态。如果线程在上下文切换时属于不同的进程，那么需要更新额外的状态信息及内存地址空间，同时将新的页表（Page Tables）导入内存。

进程之间的上下文切换最大的问题在于内存地址空间的切换导致的缓存失效（例如 CPU 中用于缓存虚拟地址与物理地址之间映射的 TLB 表），所以不同进程的切换要显著慢于同一进程中线程的切换。现代的 CPU 使用了快速上下文切换（Rapid Context Switch）技术来解决不同进程切换带来的缓存失效问题。

## 线程与协程

在 Go 语言中，协程被认为是轻量级的线程。和线程不同的是，操作系统内核感知不到协程的存在，协程的管理依赖 Go 语言运行时自身提供的调度器。同时，Go 语言中的协程是从属于某一个线程的。

为什么 Go 语言需要在线程的基础上抽象出协程的概念，而不是直接操作线程？要回答这个问题，就需要深入地理解线程与协程的区别。

下面将从调度方式、上下文切换的速度、调度策略、栈的大小这四个方面分析线程与协程的不同之处。

### 调度方式

协程是用户态的。协程的管理依赖 Go 语言运行时的调度器。同时，Go 语言中的协程是从属于某一个线程的，协程与线程的对应关系为 M: N，即多对多，如图 14-4 所示。Go 语言调度器可以将多个协程调度到一个线程中，一个协程也可能切换到多个线程中执行。

![](../../../assets/images/docs/internal/goroutine/README/图14-4%20线程与协程的对应关系.png)

### 上下文切换的速度

协程的速度要快于线程，其原因在于协程切换不用经过操作系统用户态与内核态的切换，并且 Go 语言中的协程切换只需要保留极少的状态和寄存器变量值（SP/BP/PC），而线程切换会保留额外的寄存器变量值（例如浮点寄存器）。

上下文切换的速度受到诸多因素的影响，这里列出一些值得参考的量化指标：线程切换的速度大约为 1～2 微秒，Go 语言中协程切换的速度比它快数倍，为 0.2 微秒左右。

### 调度策略

线程的调度在大部分时间是抢占式的，操作系统调度器为了均衡每个线程的执行周期，会定时发出中断信号强制执行线程上下文切换。而 Go 语言中的协程在一般情况下是协作式调度的，当一个协程处理完自己的任务后，可以主动将执行权限让渡给其他协程。这意味着协程可以更好地在规定时间内完成自己的工作，而不会轻易被抢占。当一个协程运行了过长时间时，Go 语言调度器才会强制抢占其执行。

### 栈的大小

线程的栈大小一般是在创建时指定的，为了避免出现栈溢出（Stack Overflow），默认的栈会相对较大（例如 2MB），这意味着每创建 1000 个线程就需要消耗 2GB 的虚拟内存，大大限制了线程创建的数量（64 位的虚拟内存地址空间已经让这种限制变得不太严重）。而 Go 语言中的协程栈默认为 2KB，在实践中，经常会看到成千上万的协程存在。

同时，线程的栈在运行时不能更改，但是 Go 语言中的协程栈在 Go 运行时的帮助下会动态检测栈的大小，并动态地进行扩容。因此，在实践中，可以将协程看作轻量的资源。

## 并发与并行

在 Go 语言的程序设计中，有两个非常重要但容易被误解的概念，分别是并发（concurrency）与并行（parallelism）。通俗来讲，并发指同时处理多个任务的能力，这些任务是独立的执行单元。

并发并不意味着同一时刻所有任务都在执行，而是在一个时间段内，所有的任务都能执行完毕。因此，开发者对任意时刻具体执行的是哪一个任务并不关心。如图 14-5 所示，在单核处理器中，任意一个时刻只能执行一个具体的线程，而在一个时间段内，线程可能通过上下文切换交替执行。多核处理器是真正的并行执行，因为在任意时刻，可以同时有多个线程在执行。

在实际的多核处理场景中，并发与并行常常是同时存在的，即多核在并行地处理多个线程，而单核中的多个线程又在上下文切换中交替执行。

由于 Go 语言中的协程依托于线程，所以即便处理器运行的是同一个线程，在线程内 Go 语言调度器也会切换多个协程执行，这时协程是并发的。如果多个协程被分配给了不同的线程，而这些线程同时被不同的 CPU 核心处理，那么这些协程就是并行处理的。因此在多核处理场景下，Go 语言的协程是并发与并行同时存在的。

![](../../../assets/images/docs/internal/goroutine/README/图14-5%20并发与并行的区别.png)

但是，协程的并发是一种更加常见的现象，因为处理器的核心是有限的，而一个程序中的协程数量可以成千上万，这就需要依赖 Go 语言调度器合理公平地调度。

## 简单协程入门

我们通过一个程序来检查一些网站是否可以访问，构建一个 links 作为 url 列表。

```go
package main

import (
	"fmt"
	"net/http"
)

func checkLink(link string) {
	if _, err := http.Get(link); err != nil {
		fmt.Println(link, "might be down!")
	} else {
		fmt.Println(link, "is up!")
	}
}

func main() {
	links := []string{
		"https://www.baidu.com/",
		"https://www.google.com/",
		"https://www.jd.com/",
		"https://www.taobao.com/",
		"https://www.tmall.com/",
		"https://www.sina.com.cn/",
		"https://www.sohu.com/",
		"https://www.163.com/",
	}

	for _, link := range links {
		checkLink(link)
	}
}
```

默认走系统代理。

```
https://www.baidu.com/ is up!
https://www.google.com/ is up!
https://www.jd.com/ is up!
https://www.taobao.com/ is up!
https://www.tmall.com/ is up!
https://www.sina.com.cn/ is up!
https://www.sohu.com/ is up!
https://www.163.com/ is up!
```

当前程序在正常情况下能够很好地运行，但是其有严重的性能问题。该程序为线性程序，必须等待前一个请求执行完毕，后一个请求才能继续执行。如果请求的网站出现了问题，则可能需要等待很长时间。这种情况在网络访问、磁盘文件访问时经常会遇到。

### 主协程与子协程

为了能够加快程序的执行，需要将访问修改为并发执行。这样，我们不仅能使用到多核的资源，Go 语言的调度器也能够在当前协程 I/O 堵塞时，切换到其他协程执行。

在 Go 语言中，使用协程非常方便，只需在特定的函数前加上关键字 go 即可。该关键字会被 Go 语言的编译器识别，并在运行时创建一个新的协程。新创建的协程会独立运行，不需要返回值，也不会堵塞创建它的协程。

```go
func main() {
	links := []string{
		"https://www.baidu.com/",
		"https://www.google.com/",
		"https://www.jd.com/",
		"https://www.taobao.com/",
		"https://www.tmall.com/",
		"https://www.sina.com.cn/",
		"https://www.sohu.com/",
		"https://www.163.com/",
	}

	for _, link := range links {
		go checkLink(link)
	}
}
```

当执行此程序时，我们会惊讶地发现程序直接退出了，原因在于协程（Goroutine）分为了主协程（main Goroutine）与子协程（child Goroutine），如图 14-7 所示。

![](../../../assets/images/docs/internal/goroutine/README/图14-7%20协程的两种形式.png)

main 函数是一个特殊的协程，当主协程退出时，程序直接退出，这是主协程与其他协程的显著区别，如图 14-8 所示。如果其他协程还未执行完成，主协程就直接退出了，那么此时不会有任何输出。

![](../../../assets/images/docs/internal/goroutine/README/图14-8%20主协程退出时程序直接退出.png)

明白这一点后，可以设法对程序进行调整。

```go
package main

import (
	"fmt"
	"net/http"
	"sync"
)

func checkLink(link string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, err := http.Get(link); err != nil {
		fmt.Println(link, "might be down!")
	} else {
		fmt.Println(link, "is up!")
	}
}

func main() {
	links := []string{
		"https://www.baidu.com/",
		"https://www.google.com/",
		"https://www.jd.com/",
		"https://www.taobao.com/",
		"https://www.tmall.com/",
		"https://www.sina.com.cn/",
		"https://www.sohu.com/",
		"https://www.163.com/",
	}

	wg := sync.WaitGroup{}
	for _, link := range links {
		wg.Add(1)
		go checkLink(link, &wg)
	}

	wg.Wait()
}
```

### WaitGroup

这里使用了 sync 包中的 WaitGroup 类型，它可以让主协程等待所有子协程执行完成后再退出。在主协程中，我们创建了一个 WaitGroup 类型的变量 wg，然后在每个子协程中调用 wg.Done()，这样，当所有子协程执行完成后，主协程才会继续执行。

内部实现不算复杂，WaitGroup 类型的变量 wg 会维护一个计数器，每次调用 wg.Done()，计数器就会减一，当计数器为 0 时退出。

```go
// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
//
// In the terminology of the Go memory model, a call to Done
// “synchronizes before” the return of any Wait call that it unblocks.
type WaitGroup struct {
	noCopy noCopy

	// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers only guarantee that 64-bit fields are 32-bit aligned.
	// For this reason on 32 bit architectures we need to check in state()
	// if state1 is aligned or not, and dynamically "swap" the field order if
	// needed.
	state1 uint64
	state2 uint32
}

// state returns pointers to the state and sema fields stored within wg.state*.
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	if unsafe.Alignof(wg.state1) == 8 || uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		// state1 is 64-bit aligned: nothing to do.
		return &wg.state1, &wg.state2
	} else {
		// state1 is 32-bit aligned but not 64-bit aligned: this means that
		// (&state1)+4 is 64-bit aligned.
		state := (*[3]uint32)(unsafe.Pointer(&wg.state1))
		return (*uint64)(unsafe.Pointer(&state[1])), &state[0]
	}
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	statep, semap := wg.state()
	if race.Enabled {
		_ = *statep // trigger nil deref early
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	v := int32(state >> 32)
	w := uint32(state)
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(semap))
	}
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	*statep = 0
	for ; w != 0; w-- {
		runtime_Semrelease(semap, false, 0)
	}
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	statep, semap := wg.state()
	if race.Enabled {
		_ = *statep // trigger nil deref early
		race.Disable()
	}
	for {
		state := atomic.LoadUint64(statep)
		v := int32(state >> 32)
		w := uint32(state)
		if v == 0 {
			// Counter is 0, no need to wait.
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
		// Increment waiters count.
		if atomic.CompareAndSwapUint64(statep, state, state+1) {
			if race.Enabled && w == 0 {
				// Wait must be synchronized with the first Add.
				// Need to model this is as a write to race with the read in Add.
				// As a consequence, can do the write only for the first waiter,
				// otherwise concurrent Waits will race with each other.
				race.Write(unsafe.Pointer(semap))
			}
			runtime_Semacquire(semap)
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
	}
}
```

## GMP 模型

Go 进程中的众多协程其实依托于线程，借助操作系统将线程调度到 CPU 执行，从而最终执行协程。

在 GMP 模型中，G 代表的是 Go 语言中的协程（Goroutine），M 代表的是实际的线程，而 P 代表的是 Go 逻辑处理器（Process），Go 语言为了方便协程调度与缓存，抽象出了逻辑处理器。

G、M、P 之间的对应关系如图 14-9 所示。在任一时刻，一个 P 可能在其本地包含多个 G，同时，一个 P 在任一时刻只能绑定一个 M。

![](../../../assets/images/docs/internal/goroutine/README/图14-9%20GMP模型.png)

图 14-9 中没有涵盖的信息是：一个 G 并不是固定绑定同一个 P 的，有很多情况（例如 P 在运行时被销毁）会导致一个 P 中的 G 转移到其他的 P 中。同样的，一个 P 只能对应一个 M，但是具体对应的是哪一个 M 也是不固定的。一个 M 可能在某些时候转移到其他的 P 中执行。

详细的 GMP 模型可以参考 [Go 语言的 GPM 调度器](gpm.md)。

## 常见并发模型

[常见并发模型](model.md)

```go

```
