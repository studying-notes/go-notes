---
date: 2022-10-10T11:08:44+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "锁"  # 文章标题
url:  "posts/go/docs/internal/goroutine/lock/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 原子锁

正如在数据争用中看到的，即便是简单如 count++ 这样的操作，在底层也经历了读取数据、更新 CPU 缓存、存入内存这一系列操作。这些操作如果并发进行，那么可能出现严重错误。

像 count++ 这样的操作是非原子性的，初学时可能觉得解决这样的并发问题用自定义的一个锁即可。

![](../../../../assets/images/docs/internal/goroutine/lock/README/图17-6%20利用flag模拟锁会面临的数据争用问题.png)

如图 17-6 所示，定义一个 flag 标志，当 flag 为 true 时进入区域，并立即让 flag 为 false 阻止其他协程进入。

但很快就会发现，这种方式仍然会面临数据争用的问题，flag 还是可能被同时读到 true。

还有一些更加复杂的场景需要用到原子操作，许多编译器（在编译时）和 CPU 处理器（在运行时）通过调整指令顺序进行优化，因此指令执行顺序可能与代码中显示的不同。例如，如果已知有两个内存引用将到达同一位置，并且没有中间写入会影响该位置，那么编译器可能只使用最初获取的值。又如，在编译时，a+b+c 并不能用一条 CPU 指令执行，因此按照加法结合律可能拆分为 b+c 再 +a 的形式。

在同一个协程中，代码在编译之后的顺序也可能发生变化。下例中的 setup 函数在代码中会修改 a 的值，并设置 done 为 true。main 函数会通过 done 的值来判断协程是否修改了 a，如果没有修改则暂时让渡执行权力。因此 main 函数中预期看到 a 的结果为 hello, world，但实际上这是不确定的。

setup 协程可能被编译器修改为如下：

```go
func setup() {
	done = true
	a = "hello, world"
	if done {
		fmt.Println(a)
	}
}
```

done 和 a 的赋值顺序可能被调整，因此 main 函数中的 a 可能是空值。

在 Go 语言的内存模型中，只保证了同一个协程的执行顺序，这意味着即便是编译器的重排，在同一协程执行的结果也和原始代码一致，就好像并没有发生重排（如在 setup 函数中，最后打印出长度 12）一样，但在不同协程中，观察到的写入顺序是不固定的，不同的编译器版本可能有不同的编译执行结果。

在 CPU 执行过程中，不仅可能发生编译器执行顺序混乱，也可能发生和程序中执行顺序不同的内存访问。例如，许多处理器包含存储缓冲区，该缓冲区接收对内存的挂起写操作，写缓冲区基本上是 `< 地址，数据 >` 的队列。通常，这些写操作可以按顺序执行，但是如果随后的写操作地址已经存在于写缓冲区中，则可以将其与先前的挂起写操作组合在一起。还有一种情况涉及处理器高速缓存未命中，这时在等待该指令从主内存中获取数据时，为了最大化利用资源，许多处理器将继续执行后续指令。

需要有一种机制解决并发访问时数据冲突及内存操作乱序的问题，即提供一种原子性的操作。这通常依赖硬件的支持，例如 X86 指令集中的 LOCK 指令，对应 Go 语言中的 sync/atomic 包。下例使用了 atomic.AddInt64 函数将变量加 1，这种原子操作不会发生并发时的数据争用问题。

```go
var count int64

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            atomic.AddInt64(&count, 1)
        }()
    }
    time.Sleep(time.Millisecond)
    fmt.Println(count)
}
```

在 sync/atomic 包中还有一个重要的操作 CompareAndSwap，与元素值进行对比并替换。下例判断 flag 变量的值是否为 0，如果是，则将 flag 的值设置为 1。这一系列操作都是原子性的，不会发生数据争用，也不会出现内存操作乱序问题。

通过 sync/atomic 包中的原子操作，我们能构建起一种自旋锁，只有获取该锁，才能执行区域中的代码。如下所示，使用一个 for 循环不断轮询原子操作，直到原子操作成功才获取该锁。

> 自旋锁是计算机科学用于多线程同步的一种锁，线程反复检查锁变量是否可用。由于线程在这一过程中保持执行，因此是一种忙等待。一旦获取了自旋锁，线程会一直保持该锁，直至显式释放自旋锁。
> 自旋锁避免了进程上下文的调度开销，因此对于线程只会阻塞很短时间的场合是有效的。

```go
package main

import (
	"sync/atomic"
)

var flag int64
var count int64

func add() {
	for {
		if atomic.CompareAndSwapInt64(&flag, 0, 1) {
			count++
			atomic.StoreInt64(&flag, 0)
			return
		}
	}
}
```

这种自旋锁的形式在 Go 源代码中随处可见，原子操作是底层最基础的同步保证，通过原子操作可以构建起许多同步原语，例如自旋锁、信号量、互斥锁等。

## 互斥锁

通过原子操作构建起的互斥锁，虽然高效而且简单，但不是万能的。例如，当某一个协程长时间霸占锁，其他协程抢占锁时将无意义地消耗 CPU 资源。同时当有许多正在获取锁的协程时，可能有协程一直抢占不到锁。

为了解决这种问题，操作系统的锁接口提供了终止与唤醒的机制，例如 Linux 中的 pthread mutex，避免了频繁自旋造成的浪费。在操作系统内部会构建起锁的等待队列，以便之后依次被唤醒。调用操作系统级别的锁会锁住整个线程使之无法运行，另外锁的抢占还会涉及线程之间的上下文切换。

Go 语言拥有比线程更加轻量级的协程，在协程的基础上实现了一种比传统操作系统级别的锁更加轻量级的互斥锁，其使用方式如下所示。

```go
package main

import (
	"sync"
)

var count int64
var m sync.Mutex

func add() {
	m.Lock()
	count++
	m.Unlock()
}
```

sync.Mutex 构建起了互斥锁，在同一时刻，只会有一个获取锁的协程继续执行，而其他的协程将陷入等待状态，这和自旋锁的功能是类似的，但是其提供了更加复杂的机制避免自旋锁的争用问题。

### 互斥锁实现原理

互斥锁是一种混合锁，其实现方式包含了自旋锁，同时参考了操作系统锁的实现。

sync.Mutex 结构比较简单，其包含了表示当前锁状态的 state 及信号量 sema。

```go
// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m.
// A successful call to TryLock is equivalent to a call to Lock.
// A failed call to TryLock does not establish any “synchronizes before”
// relation at all.
type Mutex struct {
	state int32
	sema  uint32
}
```

state 通过位图的形式存储了当前锁的状态，如图 17-7 所示，其中包含锁是否为锁定状态、正在等待被锁唤醒的协程数量、两个和饥饿模式有关的标志。

![](../../../../assets/images/docs/internal/goroutine/lock/README/图17-7%20互斥锁位图.png)

为了解决某一个协程可能长时间无法获取锁的问题，Go 1.9 之后使用了饥饿模式。在饥饿模式下，unlock 会唤醒最先申请加速的协程，从而保证公平。sema 是互质锁中实现的信号量。

### 互斥锁的加锁

互斥锁的第 1 个阶段是使用原子操作快速抢占锁，如果抢占成功则立即返回，如果抢占失败则调用 lockSlow 方法。

```go
// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

`lockSlow` 方法在正常情况下会自旋尝试抢占锁一段时间，而不会立即进入休眠状态，这使得互斥锁在频繁加锁与释放时也能良好工作。锁只有在正常模式下才能够进入自旋状态，`runtime_canSpin` 函数会判断当前是否能进入自旋状态，不同平台有不同的实现。

```go
func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		// Don't try to acquire starving mutex, new arriving goroutines must queue.
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// The current goroutine switches mutex to starvation mode.
		// But if the mutex is currently unlocked, don't do the switch.
		// Unlock expects that starving mutex has waiters, which will not
		// be true in this case.
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// If we were already waiting before, queue at the front of the queue.
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
}
```

在下面 4 种情况下，自旋状态立即终止：

1. 程序在单核 CPU 上运行。
2. 逻辑处理器 P 小于或等于 1。
3. 当前协程所在的逻辑处理器 P 的本地队列上有其他协程待运行。
4. 自旋次数超过了设定的阈值。

进入自旋状态后，`runtime_doSpin` 函数调用的 `procyield` 函数是一段汇编代码，会执行 30 次 PAUSE 指令占用 CPU 时间。

当长时间未获取到锁时，就进入互斥锁的第 2 个阶段，使用信号量进行同步。如果加锁操作进入信号量同步阶段，则信号量计数值减 1。如果解锁操作进入信号量同步阶段，则信号量计数值加 1。当信号量计数值大于 0 时，意味着有其他协程执行了解锁操作，这时加锁协程可以直接退出。当信号量计数值等于 0 时，意味着当前加锁协程需要陷入休眠状态。在互斥锁第 3 个阶段，所有锁的信息都会根据锁的地址存储在全局 `semtable` 哈希表中。

```go
var semtable semTable

// Prime to not correlate with any user patterns.
const semTabSize = 251

type semTable [semTabSize]struct {
	root semaRoot
	pad  [cpu.CacheLinePadSize - unsafe.Sizeof(semaRoot{})]byte
}
```

哈希函数为根据信号量地址简单取模。

```go
func (t *semTable) rootFor(addr *uint32) *semaRoot {
	return &t[(uintptr(unsafe.Pointer(addr))>>3)%semTabSize].root
}
```

图 17-8 为互斥锁加入等待队列中的示意图，先根据哈希函数查找当前锁存储在哪一个哈希桶（bucket）中。哈希结果相同的多个锁可能存储在同一个哈希桶中，哈希桶中通过一根双向链表解决哈希冲突问题。

![](../../../../assets/images/docs/internal/goroutine/lock/README/图17-8%20互斥锁加入等待队列.png)

哈希桶中的链表还被构造成了特殊的 treap 树，如图 17-9 所示。treap 树是一种引入了随机数的二叉搜索树，其实现简单，引入的随机数及必要时的旋转保证了比较好的平衡性。

将哈希桶中锁的数据结构设计为二叉搜索树的主要目的是快速查找到当前哈希桶中是否存在已经存在过的锁，这时能够以 log2N 的时间复杂度进行查找。如果已经查找到存在该锁，则将当前的协程添加到等待队列的尾部。

![](../../../../assets/images/docs/internal/goroutine/lock/README/图17-9%20桶中链表以二叉搜索树的形式排列.png)

如果不存在该锁，则需要向当前 treap 树中添加一个新的元素。值得注意的是，由于在访问哈希表时，仍然可能面临并发的数据争用，因此这里也需要加锁，但是此处的锁和互斥锁有所不同，其实现方式为先自旋一定次数，如果还没有获取到锁，则调用操作系统级别的锁，在 Linux 中为 pthread mutex 互斥锁。所以 Go 语言中的互斥锁算一种混合锁，它结合了原子操作、自旋、信号量、全局哈希表、等待队列、操作系统级别锁等多种技术，在正常情况下是基本不会进入操作系统级别的锁。

锁被放置到全局的等待队列中并等待被唤醒，唤醒的顺序为从前到后，遵循先入先出的准则，这样保证了公平性。当长时间无法获取锁时，当前的互斥锁会进入饥饿模式。在饥饿模式下，为了保证公平性，新申请锁的协程不会进入自旋状态，而是直接放入等待队列中。放入等待队列中的协程会切换自己的执行状态，让渡执行权利并进入新的调度循环，这不会暂停线程的运行。

### 互斥锁的解锁

```go
// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new)
	}
}
```

互斥锁的释放和互斥锁的锁定相对应，其步骤如下：

1. 如果当前锁处于普通的锁定状态，即没有进入饥饿状态和唤醒状态，也没有多个协程因为抢占锁陷入堵塞，则 Unlock 方法在修改 mutexLocked 状态后立即退出（快速路径）。否则，进入慢路径调用 unlockSlow 方法。
2. 判断锁是否重复释放。锁不能重复释放，否则会在运行时报错。
3. 如果锁当前处于饥饿状态，则进入信号量同步阶段，到全局哈希表中寻找当前锁的等待队列，以先入先出的顺序唤醒指定协程。
4. 如果锁当前未处于饥饿状态且当前 `mutexWoken` 已设置，则表明有其他申请锁的协程准备从正常状态退出，这时锁释放后不用去当前锁的等待队列中唤醒其他协程，而是直接退出。如果唤醒了等待队列中的协程，则将唤醒的协程放入当前协程所在逻辑处理器 P 的 runnext 字段中，存储到 runnext 字段中的协程会被优先调度。如果在饥饿模式下，则当前协程会让渡自己的执行权利，让被唤醒的协程直接运行，这是通过将 `runtime_Semrelease` 函数第 2 个参数设置为 true 实现的。

```go
func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		fatal("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
			// If there are no waiters or a goroutine has already
			// been woken or grabbed the lock, no need to wake anyone.
			// In starvation mode ownership is directly handed off from unlocking
			// goroutine to the next waiter. We are not part of this chain,
			// since we did not observe mutexStarving when we unlocked the mutex above.
			// So get off the way.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// Starving mode: handoff mutex ownership to the next waiter, and yield
		// our time slice so that the next waiter can start to run immediately.
		// Note: mutexLocked is not set, the waiter will set it after wakeup.
		// But mutex is still considered locked if mutexStarving is set,
		// so new coming goroutines won't acquire it.
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```

## 读写锁

在同一时间内只能有一个协程获取互斥锁并执行操作，在多读少写的情况下，如果长时间没有写操作，那么读取到的会是完全相同的值，完全不需要通过互斥的方式获取，这是读写锁产生的背景。

读写锁通过两种锁来实现，一种为读锁，另一种为写锁。当进行读取操作时，需要加读锁，而进行写入操作时需要加写锁。多个协程可以同时获得读锁并执行。如果此时有协程申请了写锁，那么该写锁会等待所有的读锁都释放后才能获取写锁继续执行。如果当前的协程申请读锁时已经存在写锁，那么读锁会等待写锁释放后再获取锁继续执行。

总之，读锁必须能观察到上一次写锁写入的值，写锁要等待之前的读锁释放才能写入。可能有多个协程获得读锁，但只有一个协程获得写锁。

举一个简单的例子，哈希表并不是并发安全的，它只能够并发读取，并发写入时会出现冲突。一种简单的规避方式如下所示，可以在获取 map 中的数据时加入 RLock 读锁，在写入数据时使用 Lock 写锁。

```go
package main

import "sync"

type Stat struct {
	counters map[string]int64
	mutex    sync.RWMutex
}

func (s *Stat) getCounter(name string) int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.counters[name]
}

func (s *Stat) setCounter(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.counters[name]++
}
```

## 读写锁原理

读写锁位于 sync 标准库中，其结构如下。读写锁复用了互斥锁及信号量这两种机制。

```go
// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m, just as for Mutex.
// For any call to RLock, there exists an n such that
// the n'th call to Unlock “synchronizes before” that call to RLock,
// and the corresponding call to RUnlock “synchronizes before”
// the n+1'th call to Lock.
type RWMutex struct {
	w           Mutex  // held if there are pending writers
	writerSem   uint32 // semaphore for writers to wait for completing readers
	readerSem   uint32 // semaphore for readers to wait for completing writers
	readerCount int32  // number of pending readers
	readerWait  int32  // number of departing readers
}
```

读取操作先通过原子操作将 readerCount 加 1，如果 `readerCount≥0` 就直接返回，所以如果只有获取读取锁的操作，那么其成本只有一个原子操作。当 readerCount<0 时，说明当前有写锁，当前协程将借助信号量陷入等待状态，如果获取到信号量则立即退出，没有获取到信号量时的逻辑与互斥锁的逻辑相似。

```go
// Happens-before relationships are indicated to the race detector via:
// - Unlock  -> Lock:  readerSem
// - Unlock  -> RLock: readerSem
// - RUnlock -> Lock:  writerSem
//
// The methods below temporarily disable handling of race synchronization
// events in order to provide the more precise model above to the race
// detector.
//
// For example, atomic.AddInt32 in RLock should not appear to provide
// acquire-release semantics, which would incorrectly synchronize racing
// readers, thus potentially missing races.

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// A writer is pending, wait for it.
		runtime_SemacquireMutex(&rw.readerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
	}
}
```

读锁解锁时，如果当前没有写锁，则其成本只有一个原子操作并直接退出。

如果当前有写锁正在等待，则调用 rUnlockSlow 判断当前是否为最后一个被释放的读锁，如果是则需要增加信号量并唤醒写锁。

```go
// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	if race.Enabled {
		_ = rw.w.state
		race.ReleaseMerge(unsafe.Pointer(&rw.writerSem))
		race.Disable()
	}
	if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		// Outlined slow-path to allow the fast-path to be inlined
		rw.rUnlockSlow(r)
	}
	if race.Enabled {
		race.Enable()
	}
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	if r+1 == 0 || r+1 == -rwmutexMaxReaders {
		race.Enable()
		fatal("sync: RUnlock of unlocked RWMutex")
	}
	// A writer is pending.
	if atomic.AddInt32(&rw.readerWait, -1) == 0 {
		// The last reader unblocks the writer.
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}
```

读写锁申请写锁时要调用 Lock 方法，必须先获取互斥锁，因为它复用了互斥锁的功能。接着 readerCount 减去 rwmutexMaxReaders 阻止后续的读操作。

但获取互斥锁并不一定能直接获取写锁，如果当前已经有其他 Goroutine 持有互斥锁的读锁，那么当前协程会加入全局等待队列并进入休眠状态，当最后一个读锁被释放时，会唤醒该协程。

```go
// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	if race.Enabled {
		_ = rw.w.state
		race.Release(unsafe.Pointer(&rw.readerSem))
		race.Disable()
	}

	// Announce to readers there is no active writer.
	r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		race.Enable()
		fatal("sync: Unlock of unlocked RWMutex")
	}
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}
	// Allow other writers to proceed.
	rw.w.Unlock()
	if race.Enabled {
		race.Enable()
	}
}
```

可以看出，读写锁在写操作时的性能与互斥锁类似，但是在只有读操作时效率要高很多，因为读锁可以被多个协程获取。

```go

```
