---
date: 2022-10-05T16:18:07+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "异常与异常捕获"  # 文章标题
url:  "posts/go/docs/internal/panic/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

异常退出可能是由于用户代码中错误的状态引起的；或者是运行时数组越界、哈希表读写冲突等引起的；也可能是访问无效内存引起的。有时候，我们不希望程序异常退出，而是希望捕获异常并让函数正常执行，这涉及 defer、recover 的结合使用。

## panic 函数使用方法

有以下两个内置函数可以处理程序的异常情况：

```go
func panic(v interface{})
func recover() interface{}
```

panic 函数传递的参数为空接口 interface{}，其可以存储任何形式的错误信息并进行传递。在异常退出时会打印出来。

Go 语言函数调用的正常流程是函数执行返回语句返回结果，在这个流程中是没有异常的，因此在这个流程中执行 `recover` 异常捕获函数始终是返回 `nil`。

另一种是异常流程: 当函数调用 `panic` 抛出异常，函数将停止执行后续的普通语句，但是之前注册的 `defer` 函数调用仍然保证会被正常执行，然后再返回到调用者。对于当前函数的调用者，因为处理异常状态还没有被捕获，和直接调用 `panic` 函数的行为类似。在异常发生时，如果在 `defer` 中执行 `recover` 调用，它可以捕获触发 `panic` 时的参数，并且恢复到正常的执行流程。

`panic` 会停掉当前正在执行的程序，但是与 `os.Exit(-1)` 这种退出不同，`panic` 的撤退比较有秩序，先处理完**当前** `goroutine` 已经 `defer` 挂上去的任务，即在 `panic` 语句之前注册的 `defer` 语句，执行完毕后再退出整个程序。

`panic` 的原则是：执行且只执行当前 `goroutine` 的 `defer`，且在 `panic` 之前注册。

`panic` 仅保证当前 goroutine 下的 `defer` 都会被执行到，但不保证其他协程的 `defer` 也会执行到。如果是在同一 `goroutine` 下的调用者的 `defer`，那么可以一路回溯回去执行；但如果是不同 `goroutine`，那就不做保证了。

`recover` 只在 `defer` 的函数中有效，如果不是在 `defer` 上下文中调用，`recover` 会直接返回 `nil`。

Go 程序在 panic 时并不会像大多数人想象的一样导致程序异常退出，而是会终止当前函数的正常执行，执行 defer 函数并逐级返回。例如，对于函数调用链 a()→ b()→ c()，当函数 c 发生 panic 后，会返回函数 b。此时，函数 b 也像发生了 panic 一样，返回函数 a。在函数 c、b、a 中的 defer 函数都将正常执行。

```go
package main

import "fmt"

func a() {
	defer fmt.Println("defer a")
	b()
	fmt.Println("after a")
}

func b() {
	defer fmt.Println("defer b")
	c()
	fmt.Println("after b")
}

func c() {
	defer fmt.Println("defer c")
	panic("panic c")
	fmt.Println("after c")
}

func main() {
	a()
}
```

当函数 c 触发了 panic 后，所有函数中的 defer 语句都将被正常调用，并且在 panic 时打印出堆栈信息。

除了手动触发 panic，在 Go 语言运行时的一些阶段也会检查并触发 panic，例如数组越界(runtime error: index out of range) 及 map 并发冲突(fatal error: concurrent map read and map write)。

## 异常捕获

为了让程序在 panic 时仍然能够执行后续的流程，Go 语言提供了内置的 recover 函数用于异常恢复。recover 函数一般与 defer 函数结合使用才有意义，其返回值是 panic 中传递的参数。由于 panic 会调用 defer 函数，因此，在 defer 函数中可以加入 recover 起到让函数恢复正常执行的作用。

在非 `defer` 语句中执行 `recover` 调用是初学者常犯的错误:

```go
func main() {
    if r := recover(); r != nil {
        log.Fatal(r)
    }

    panic(123)

    if r := recover(); r != nil {
        log.Fatal(r)
    }
}
```

上面程序中两个 `recover` 调用都不能捕获任何异常。在第一个 `recover` 调用执行时，函数必然是在正常的非异常执行流程中，这时候 `recover` 调用将返回 `nil`。发生异常时，第二个 `recover` 调用将没有机会被执行到，因为 `panic` 调用会导致函数马上执行已经注册 `defer` 的函数后返回。

### 条件


其实 `recover` 函数调用有着更严格的要求：我们必须在 `defer` 函数中直接调用 `recover`。如果 `defer` 中调用的是 `recover` 函数的包装函数的话，异常的捕获工作将失败！比如，有时候我们可能希望包装自己的 `MyRecover` 函数，在内部增加必要的日志信息然后再调用 `recover`，这是错误的做法：

```go
func main() {
    defer func() {
        // 无法捕获异常
        if r := MyRecover(); r != nil {
            fmt.Println(r)
        }
    }()
    panic(1)
}

func MyRecover() interface{} {
    log.Println("trace...")
    return recover()
}
```

同样，如果是在嵌套的 `defer` 函数中调用 `recover` 也将导致无法捕获异常：

```go
func main() {
    defer func() {
        defer func() {
            // 无法捕获异常
            if r := recover(); r != nil {
                fmt.Println(r)
            }
        }()
    }()
    panic(1)
}
```

2 层嵌套的 `defer` 函数中直接调用 `recover` 和 1 层 `defer` 函数中调用包装的 `MyRecover` 函数一样，都是经过了 2 个函数帧才到达真正的 `recover` 函数，这个时候 Goroutine 的对应上一级栈帧中已经没有异常信息。

如果我们直接在 `defer` 语句中调用 `MyRecover` 函数又可以正常工作了：

```go
func MyRecover() interface{} {
    return recover()
}

func main() {
    // 可以正常捕获异常
    defer MyRecover()
    panic(1)
}
```

但是，如果 `defer` 语句直接调用 `recover` 函数，依然不能正常捕获异常：

```go
func main() {
    // 无法捕获异常
    defer recover()
    panic(1)
}
```

必须要和有异常的栈帧只隔一个栈帧，`recover` 函数才能正常捕获异常。换言之，`recover` 函数捕获的是祖父一级调用函数栈帧的异常（刚好可以跨越一层 `defer` 函数）！

当然，为了避免 `recover` 调用者不能识别捕获到的异常, 应该避免用 `nil` 为参数抛出异常:

```go
func main() {
    defer func() {
        if r := recover(); r != nil { ... }
        // 虽然总是返回nil, 但是可以恢复异常状态
    }()

    // 警告: 用`nil`为参数抛出异常
    panic(nil)
}
```

当希望将捕获到的异常转为错误时，如果希望忠实返回原始的信息，需要针对不同的类型分别处理：

```go
func foo() (err error) {
    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case string:
                err = errors.New(x)
            case error:
                err = x
            default:
                err = fmt.Errorf("Unknown panic: %v", r)
            }
        }
    }()

    panic("TODO")
}
```

以下三个条件会让 recover() 返回 `nil` :

1. panic 时指定的参数为 `nil` ；(一般 panic 语句如 `panic( "xxx failed... " )`)
2. 当前协程没有发生 panic ；
3. recover 没有被 defer 方法**直接调用**；

### 改进

对上面的例子进行改进后的 recover 版本如下。

```go
func b() {
	defer func() {
		fmt.Println("defer b")
		if err := recover(); err != nil {
			fmt.Println("recover b")
		}
	}()
	c()
	fmt.Println("after b")
}
```

函数 c 在触发了 panic 之后，会调用 defer 函数，接着返回函数 b，执行函数 b 中的 defer 函数。由于 defer 函数中加入了 recover 函数进行异常捕获，因此，当函数 b 结束返回函数 a 后，函数 b 就像是正常退出一样，函数 a 继续正常执行其后的流程。

panic 会遍历 defer 链并调用，那么如果在 defer 函数中发生了 panic 会怎么样呢？如下所示，函数 a 触发了 panic 后调用 defer 函数 b，而函数 b 触发了 panic 调用 defer 函数 c，函数 c 同样触发了 panic。

```go
package main

func a() {
	defer b()
	panic("panic a")
}

func b() {
	defer c()
	panic("panic b")
}

func c() {
	panic("panic c")
}

func main() {
	a()
}
```

最终程序输出如下，先打印最早出现的 panic，再打印其他的 panic。后面会看到，每一次 panic 调用都新建了一个 _panic 结构体，并用一个链表进行了存储。

```
panic: panic a
        panic: panic b
        panic: panic c

```

嵌套 panic 不会陷入死循环，每个 defer 函数都只会被调用一次。当嵌套的 panic 遇到了 recover 时，情况变得更加复杂。将上面的程序稍微改进一下，让 main 函数捕获嵌套的 panic。

```go
func catch() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
func main() {
	defer catch()
	a()
}
```

recover 函数最终捕获的是最近发生的 panic，即便有多个 panic 函数，在最上层的函数也只需要一个 recover 函数就能让函数按照正常的流程执行。

## panic 函数底层原理

panic 函数在编译时会被解析为调用运行时 runtime.gopanic 函数，如下所示。

`src/runtime/panic.go`

```go
// getg returns the pointer to the current g.
// The compiler rewrites calls to this function into instructions
// that fetch the g directly (from TLS or from the dedicated register).
func getg() *g

// The implementation of the predeclared function panic.
func gopanic(e any) {
	gp := getg()

    var p _panic
	p.arg = e
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	runningPanicDefers.Add(1)

	// By calculating getcallerpc/getcallersp here, we avoid scanning the
	// gopanic frame (stack scanning is slow...)
	addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))
}
```

每调用一次 panic 都会创建一个 _panic 结构体。和 _defer 结构体一样，_panic 也会被放置到当前协程的链表中，原因是 panic 可能发生嵌套，例如 panic → defer → panic，因此可能同时存在多个 _panic 结构体。

首先查看在正常情况下 panic 的执行流程，其简单遍历协程中的 defer 链表，对于通过堆分配或者栈分配实现的 defer 语句，通过反射的方式调用 defer 中的函数。panic 通过 reflectcall 调用 defered 函数而不是直接调用的原因在于，直接调用 defered 函数需要在当前栈帧中为它准备参数，而不同 defered 函数的参数大小可能有很大差异，然而 gopanic 函数的栈帧大小固定而且很小，所以可能没有足够的空间来存放 defered 函数的参数。

在正常情况下，当 defer 链表遍历完毕后，panic 会退出。但是这里有一个例外，Go 1.14 之后的版本通过内联汇编实现的 defer 并不会被放置到链表中存储，而是被放置到了栈上，那么如何保证 defer 在 panic 时内联 defer 函数还能正常执行呢？

```go
func gopanic(e any) {
	for {
		d := gp._defer
		if d == nil {
			break
		}

		// If defer was started by earlier panic or Goexit (and, since we're back here, that triggered a new panic),
		// take defer off list. An earlier panic will not continue running, but we will make sure below that an
		// earlier Goexit does continue running.
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			if !d.openDefer {
				// For open-coded defers, we need to process the
				// defer again, in case there are any other defers
				// to call in the frame (not including the defer
				// call that caused the panic).
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}
    }
```

这需要借助编译时与运行时的共同努力。注意上方的 addOneOpenDeferFrame 函数，该函数将调用 gentraceback 函数进行栈扫描，从调用 panic 的当前函数栈帧开始扫描，直到找到第一个包含内联 defer 的函数帧，并构建一个新的 _defer 结构体存储到协程的 _defer 链表中。

例如对于如下构造的代码：函数 a 中的函数 fa 与函数 c 中的函数 fc 都将被放入协程的 _defer 链表，但是函数 b 中的 defer 由于内联优化并不会被放入链表。

```go
package main

func a() {
	for i := 0; i < 3; i++ {
		defer fa()
	}
	b()
}

func fa() {
	println("fa")
}

func b() {
	defer fb1()
	if true {
		defer fb2()
	}
	c()
}

func fb1() {
	println("fb1")
}

func fb2() {
	println("fb2")
}

func c() {
	for i := 0; i < 3; i++ {
		defer fc()
	}

	panic("c")
}

func fc() {
	println("fc")
}
```

当函数 c 发生 panic 后，runtime.addOneOpenDeferFrame 函数会尝试遍历函数帧，当其遍历到 b 函数的函数栈帧时，发现了内联函数，将创建一个新的 _defer 结构体并加入协程 _defer 链表 defer a 函数与 defer c 函数中间，这种顺序保证了之后的 defer 能够按照先入后出的顺序排列，如图 11-1 所示。

![](../../../assets/images/docs/internal/panic/README/图11-1%20panic保证defer先入后出的机制.png)

addOneOpenDeferFrame 函数每次只会将扫描到的一个栈帧加入 defer 链表，_defer 结构体中专门有一个字段 fd 存储了栈帧的元数据，用于在运行时查找对应的内联 defer 的一系列函数指针、参数及 defer 位图。当遍历 _defer 链表的过程中发现 d.openDeferw 为 true 时，会调用 runtime.runOpenDeferFrame 方法执行某一个函数中所有需要被执行的 defer 函数，除非在这期间发生了 recovered。

```go
// addOneOpenDeferFrame scans the stack (in gentraceback order, from inner frames to
// outer frames) for the first frame (if any) with open-coded defers. If it finds
// one, it adds a single entry to the defer chain for that frame. The entry added
// represents all the defers in the associated open defer frame, and is sorted in
// order with respect to any non-open-coded defers.
//
// addOneOpenDeferFrame stops (possibly without adding a new entry) if it encounters
// an in-progress open defer entry. An in-progress open defer entry means there has
// been a new panic because of a defer in the associated frame. addOneOpenDeferFrame
// does not add an open defer entry past a started entry, because that started entry
// still needs to finished, and addOneOpenDeferFrame will be called when that started
// entry is completed. The defer removal loop in gopanic() similarly stops at an
// in-progress defer entry. Together, addOneOpenDeferFrame and the defer removal loop
// ensure the invariant that there is no open defer entry further up the stack than
// an in-progress defer, and also that the defer removal loop is guaranteed to remove
// all not-in-progress open defer entries from the defer chain.
//
// If sp is non-nil, addOneOpenDeferFrame starts the stack scan from the frame
// specified by sp. If sp is nil, it uses the sp from the current defer record (which
// has just been finished). Hence, it continues the stack scan from the frame of the
// defer that just finished. It skips any frame that already has a (not-in-progress)
// open-coded _defer record in the defer chain.
//
// Note: All entries of the defer chain (including this new open-coded entry) have
// their pointers (including sp) adjusted properly if the stack moves while
// running deferred functions. Also, it is safe to pass in the sp arg (which is
// the direct result of calling getcallersp()), because all pointer variables
// (including arguments) are adjusted as needed during stack copies.
func addOneOpenDeferFrame(gp *g, pc uintptr, sp unsafe.Pointer) {
	var prevDefer *_defer
	if sp == nil {
		prevDefer = gp._defer
		pc = prevDefer.framepc
		sp = unsafe.Pointer(prevDefer.sp)
	}
	systemstack(func() {
		gentraceback(pc, uintptr(sp), 0, gp, 0, nil, 0x7fffffff,
			func(frame *stkframe, unused unsafe.Pointer) bool {
				if prevDefer != nil && prevDefer.sp == frame.sp {
					// Skip the frame for the previous defer that
					// we just finished (and was used to set
					// where we restarted the stack scan)
					return true
				}
				f := frame.fn
				fd := funcdata(f, _FUNCDATA_OpenCodedDeferInfo)
				if fd == nil {
					return true
				}
				// Insert the open defer record in the
				// chain, in order sorted by sp.
				d := gp._defer
				var prev *_defer
				for d != nil {
					dsp := d.sp
					if frame.sp < dsp {
						break
					}
					if frame.sp == dsp {
						if !d.openDefer {
							throw("duplicated defer entry")
						}
						// Don't add any record past an
						// in-progress defer entry. We don't
						// need it, and more importantly, we
						// want to keep the invariant that
						// there is no open defer entry
						// passed an in-progress entry (see
						// header comment).
						if d.started {
							return false
						}
						return true
					}
					prev = d
					d = d.link
				}
				if frame.fn.deferreturn == 0 {
					throw("missing deferreturn")
				}

				d1 := newdefer()
				d1.openDefer = true
				d1._panic = nil
				// These are the pc/sp to set after we've
				// run a defer in this frame that did a
				// recover. We return to a special
				// deferreturn that runs any remaining
				// defers and then returns from the
				// function.
				d1.pc = frame.fn.entry() + uintptr(frame.fn.deferreturn)
				d1.varp = frame.varp
				d1.fd = fd
				// Save the SP/PC associated with current frame,
				// so we can continue stack trace later if needed.
				d1.framepc = frame.pc
				d1.sp = frame.sp
				d1.link = d
				if prev == nil {
					gp._defer = d1
				} else {
					prev.link = d1
				}
				// Stop stack scanning after adding one open defer record
				return false
			},
			nil, 0)
	})
}
```

## recover 底层原理

内置的 recover 函数在运行时将被转换为调用 runtime.gorecover 函数。

```go
// The implementation of the predeclared function recover.
// Cannot split the stack because it needs to reliably
// find the stack segment of its caller.
//
// TODO(rsc): Once we commit to CopyStackAlways,
// this doesn't need to be nosplit.
//
//go:nosplit
func gorecover(argp uintptr) any {
	// Must be in a function running as part of a deferred call during the panic.
	// Must be called from the topmost function of the call
	// (the function used in the defer statement).
	// p.argp is the argument pointer of that topmost deferred function call.
	// Compare against argp reported by caller.
	// If they match, the caller is the one who can recover.
	gp := getg()
	p := gp._panic
	if p != nil && !p.goexit && !p.recovered && argp == uintptr(p.argp) {
		p.recovered = true
		return p.arg
	}
	return nil
}
```

gorecover 函数的参数 argp 为其调用者函数的参数地址，而 p.argp 为发生 panic 时 defer 函数的参数地址。语句 argp == uintptr（p.argp）可用于判断 panic 和 recover 是否匹配，这是因为内层 recover 不能捕获外层的 panic。

gorecover 并没有进行任何异常处理，真正的处理发生在 runtime.gopanic 函数中，在遍历 defer 链表执行的过程中，一旦发现 p.recovered 为 true，就代表当前 defer 中调用了 recover 函数，会删除当前链表中为内联 defer 的 _defer 结构。原因是之后程序将恢复正常的流程，内联 defer 直接通过内联的方式执行。

```go
func gopanic(e any) {
		if p.recovered {
			gp._panic = p.link
			if gp._panic != nil && gp._panic.goexit && gp._panic.aborted {
				// A normal recover would bypass/abort the Goexit.  Instead,
				// we return to the processing loop of the Goexit.
				gp.sigcode0 = uintptr(gp._panic.sp)
				gp.sigcode1 = uintptr(gp._panic.pc)
				mcall(recovery)
				throw("bypassed recovery failed") // mcall should not return
			}
			runningPanicDefers.Add(-1)

			// After a recover, remove any remaining non-started,
			// open-coded defer entries, since the corresponding defers
			// will be executed normally (inline). Any such entry will
			// become stale once we run the corresponding defers inline
			// and exit the associated stack frame. We only remove up to
			// the first started (in-progress) open defer entry, not
			// including the current frame, since any higher entries will
			// be from a higher panic in progress, and will still be
			// needed.
			d := gp._defer
			var prev *_defer
			if !done {
				// Skip our current frame, if not done. It is
				// needed to complete any remaining defers in
				// deferreturn()
				prev = d
				d = d.link
			}
			for d != nil {
				if d.started {
					// This defer is started but we
					// are in the middle of a
					// defer-panic-recover inside of
					// it, so don't remove it or any
					// further defer entries
					break
				}
				if d.openDefer {
					if prev == nil {
						gp._defer = d.link
					} else {
						prev.link = d.link
					}
					newd := d.link
					freedefer(d)
					d = newd
				} else {
					prev = d
					d = d.link
				}
			}

			gp._panic = p.link
			// Aborted panics are marked but remain on the g.panic list.
			// Remove them from the list.
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // must be done with signal
				gp.sig = 0
			}
			// Pass information about recovering frame to recovery.
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			mcall(recovery)
			throw("recovery failed") // mcall should not return
		}
}
```

gopanic 函数最后会调用 mcall（recovery），mcall 将切换到 g0 栈执行 recovery 函数。recovery 函数接受从 defer 中传递进来的 SP、PC 寄存器并借助 gogo 函数切换到当前协程继续执行。gogo 函数是与操作系统有关的函数，用于完成栈的切换及 CPU 寄存器的恢复。只不过当前协程的 SP、PC 寄存器已经被修改为了从 defer 中传递进来的 SP、PC 寄存器，从而改变了函数的执行路径。

```go
// Unwind the stack after a deferred function calls recover
// after a panic. Then arrange to continue running as though
// the caller of the deferred function returned normally.
func recovery(gp *g) {
	// Info about defer passed in G struct.
	sp := gp.sigcode0
	pc := gp.sigcode1

	// d's arguments need to be in the stack.
	if sp != 0 && (sp < gp.stack.lo || gp.stack.hi < sp) {
		print("recover: ", hex(sp), " not in [", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
		throw("bad recovery")
	}

	// Make the deferproc for this d return again,
	// this time returning 1. The calling function will
	// jump to the standard return epilogue.
	gp.sched.sp = sp
	gp.sched.pc = pc
	gp.sched.lr = 0
	gp.sched.ret = 1
	gogo(&gp.sched)
}
```

对于 defer 栈分配与堆分配，PC 地址其实对应的是指令 TESTL AX，AX。如下例汇编代码所示，CALL runtime.deferprocStack 的下一条指令为 TESTL AX，AX。在正常情况下，AX 寄存器为 0，因此不会跳转，而是执行正常的流程。在 panic+recover 异常捕获后，AX 寄存器置为 1，因此程序将发生跳转，直接调用 runtime.deferreturn 函数。

对于内联 defer，addOneOpenDeferFrame 函数会直接将 PC 的地址设置为 deferreturn 函数的地址。

总之，recover 通过修改协程 SP、PC 寄存器值使函数重新执行 deferreturn 函数。deferreturn 函数的作用是继续执行剩余的 defer 函数（因为有一部分 defer 函数可能已经在 gopanic 函数中得到了执行），并返回到调用者函数，就像程序并没有 panic 一样。

```go

```
