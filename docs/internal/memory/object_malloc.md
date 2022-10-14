---
date: 2022-10-14T09:56:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "对象分配"  # 文章标题
url:  "posts/go/docs/internal/memory/object_malloc"  # 设置网页永久链接
tags: [ "Go", "block" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 对象分配

在运行时分配对象的逻辑主要位于 mallocgc 函数中 `src/runtime/malloc.go`，malloc 代表分配，gc 代表垃圾回收（GC），此函数除了分配内存还会为垃圾回收做一些位图标记工作。

内存分配时，将对象按照大小不同划分为微小对象、小对象、大对象。

```go
const (
	maxTinySize   = _TinySize
	tinySizeClass = _TinySizeClass
	maxSmallSize  = _MaxSmallSize

	// Tiny allocator parameters, see "Tiny allocator" comment in malloc.go.
	_TinySize      = 16
	_TinySizeClass = int8(2)
)
```

微小对象的分配流程最长，逻辑链路最复杂。

```go
// Allocate an object of size bytes.
// Small objects are allocated from the per-P cache's free lists.
// Large objects (> 32 kB) are allocated straight from the heap.
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
	if gcphase == _GCmarktermination {
		throw("mallocgc called with gcphase == _GCmarktermination")
	}

	if size == 0 {
		return unsafe.Pointer(&zerobase)
	}

	// It's possible for any malloc to trigger sweeping, which may in
	// turn queue finalizers. Record this dynamic lock edge.
	lockRankMayQueueFinalizer()

	userSize := size
	if asanenabled {
		// Refer to ASAN runtime library, the malloc() function allocates extra memory,
		// the redzone, around the user requested memory region. And the redzones are marked
		// as unaddressable. We perform the same operations in Go to detect the overflows or
		// underflows.
		size += computeRZlog(size)
	}

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

	// Set mp.mallocing to keep from being preempted by GC.
	mp := acquirem()
	if mp.mallocing != 0 {
		throw("malloc deadlock")
	}
	if mp.gsignal == getg() {
		throw("malloc during signal")
	}
	mp.mallocing = 1

	shouldhelpgc := false
	dataSize := userSize
	c := getMCache(mp)
	if c == nil {
		throw("mallocgc called without a P or outside bootstrapping")
	}
	var span *mspan
	var x unsafe.Pointer
	noscan := typ == nil || typ.ptrdata == 0
	// In some cases block zeroing can profitably (for latency reduction purposes)
	// be delayed till preemption is possible; delayedZeroing tracks that state.
	delayedZeroing := false
	if size <= maxSmallSize {
		if noscan && size < maxTinySize {
			// Tiny allocator.
			//
			// Tiny allocator combines several tiny allocation requests
			// into a single memory block. The resulting memory block
			// is freed when all subobjects are unreachable. The subobjects
			// must be noscan (don't have pointers), this ensures that
			// the amount of potentially wasted memory is bounded.
			//
			// Size of the memory block used for combining (maxTinySize) is tunable.
			// Current setting is 16 bytes, which relates to 2x worst case memory
			// wastage (when all but one subobjects are unreachable).
			// 8 bytes would result in no wastage at all, but provides less
			// opportunities for combining.
			// 32 bytes provides more opportunities for combining,
			// but can lead to 4x worst case wastage.
			// The best case winning is 8x regardless of block size.
			//
			// Objects obtained from tiny allocator must not be freed explicitly.
			// So when an object will be freed explicitly, we ensure that
			// its size >= maxTinySize.
			//
			// SetFinalizer has a special case for objects potentially coming
			// from tiny allocator, it such case it allows to set finalizers
			// for an inner byte of a memory block.
			//
			// The main targets of tiny allocator are small strings and
			// standalone escaping variables. On a json benchmark
			// the allocator reduces number of allocations by ~12% and
			// reduces heap size by ~20%.
			off := c.tinyoffset
			// Align tiny pointer for required (conservative) alignment.
			if size&7 == 0 {
				off = alignUp(off, 8)
			} else if goarch.PtrSize == 4 && size == 12 {
				// Conservatively align 12-byte objects to 8 bytes on 32-bit
				// systems so that objects whose first field is a 64-bit
				// value is aligned to 8 bytes and does not cause a fault on
				// atomic access. See issue 37262.
				// TODO(mknyszek): Remove this workaround if/when issue 36606
				// is resolved.
				off = alignUp(off, 8)
			} else if size&3 == 0 {
				off = alignUp(off, 4)
			} else if size&1 == 0 {
				off = alignUp(off, 2)
			}
			if off+size <= maxTinySize && c.tiny != 0 {
				// The object fits into existing tiny block.
				x = unsafe.Pointer(c.tiny + off)
				c.tinyoffset = off + size
				c.tinyAllocs++
				mp.mallocing = 0
				releasem(mp)
				return x
			}
			// Allocate a new maxTinySize block.
			span = c.alloc[tinySpanClass]
			v := nextFreeFast(span)
			if v == 0 {
				v, span, shouldhelpgc = c.nextFree(tinySpanClass)
			}
			x = unsafe.Pointer(v)
			(*[2]uint64)(x)[0] = 0
			(*[2]uint64)(x)[1] = 0
			// See if we need to replace the existing tiny block with the new one
			// based on amount of remaining free space.
			if !raceenabled && (size < c.tinyoffset || c.tiny == 0) {
				// Note: disabled when race detector is on, see comment near end of this function.
				c.tiny = uintptr(x)
				c.tinyoffset = size
			}
			size = maxTinySize
		} else {
			var sizeclass uint8
			if size <= smallSizeMax-8 {
				sizeclass = size_to_class8[divRoundUp(size, smallSizeDiv)]
			} else {
				sizeclass = size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]
			}
			size = uintptr(class_to_size[sizeclass])
			spc := makeSpanClass(sizeclass, noscan)
			span = c.alloc[spc]
			v := nextFreeFast(span)
			if v == 0 {
				v, span, shouldhelpgc = c.nextFree(spc)
			}
			x = unsafe.Pointer(v)
			if needzero && span.needzero != 0 {
				memclrNoHeapPointers(unsafe.Pointer(v), size)
			}
		}
	} else {
		shouldhelpgc = true
		// For large allocations, keep track of zeroed state so that
		// bulk zeroing can be happen later in a preemptible context.
		span = c.allocLarge(size, noscan)
		span.freeindex = 1
		span.allocCount = 1
		size = span.elemsize
		x = unsafe.Pointer(span.base())
		if needzero && span.needzero != 0 {
			if noscan {
				delayedZeroing = true
			} else {
				memclrNoHeapPointers(x, size)
				// We've in theory cleared almost the whole span here,
				// and could take the extra step of actually clearing
				// the whole thing. However, don't. Any GC bits for the
				// uncleared parts will be zero, and it's just going to
				// be needzero = 1 once freed anyway.
			}
		}
	}

	if !noscan {
		var scanSize uintptr
		heapBitsSetType(uintptr(x), size, dataSize, typ)
		if dataSize > typ.size {
			// Array allocation. If there are any
			// pointers, GC has to scan to the last
			// element.
			if typ.ptrdata != 0 {
				scanSize = dataSize - typ.size + typ.ptrdata
			}
		} else {
			scanSize = typ.ptrdata
		}
		c.scanAlloc += scanSize
	}

	// Ensure that the stores above that initialize x to
	// type-safe memory and set the heap bits occur before
	// the caller can make x observable to the garbage
	// collector. Otherwise, on weakly ordered machines,
	// the garbage collector could follow a pointer to x,
	// but see uninitialized memory or stale heap bits.
	publicationBarrier()

	// Allocate black during GC.
	// All slots hold nil so no scanning is needed.
	// This may be racing with GC so do it atomically if there can be
	// a race marking the bit.
	if gcphase != _GCoff {
		gcmarknewobject(span, uintptr(x), size)
	}

	if raceenabled {
		racemalloc(x, size)
	}

	if msanenabled {
		msanmalloc(x, size)
	}

	if asanenabled {
		// We should only read/write the memory with the size asked by the user.
		// The rest of the allocated memory should be poisoned, so that we can report
		// errors when accessing poisoned memory.
		// The allocated memory is larger than required userSize, it will also include
		// redzone and some other padding bytes.
		rzBeg := unsafe.Add(x, userSize)
		asanpoison(rzBeg, size-userSize)
		asanunpoison(x, userSize)
	}

	if rate := MemProfileRate; rate > 0 {
		// Note cache c only valid while m acquired; see #47302
		if rate != 1 && size < c.nextSample {
			c.nextSample -= size
		} else {
			profilealloc(mp, x, size)
		}
	}
	mp.mallocing = 0
	releasem(mp)

	// Pointerfree data can be zeroed late in a context where preemption can occur.
	// x will keep the memory alive.
	if delayedZeroing {
		if !noscan {
			throw("delayed zeroing on data that may contain pointers")
		}
		memclrNoHeapPointersChunked(size, x) // This is a possible preemption point: see #47302
	}

	if debug.malloc {
		if debug.allocfreetrace != 0 {
			tracealloc(x, size, typ)
		}

		if inittrace.active && inittrace.id == getg().goid {
			// Init functions are executed sequentially in a single goroutine.
			inittrace.bytes += uint64(size)
		}
	}

	if assistG != nil {
		// Account for internal fragmentation in the assist
		// debt now that we know it.
		assistG.gcAssistBytes -= int64(size - dataSize)
	}

	if shouldhelpgc {
		if t := (gcTrigger{kind: gcTriggerHeap}); t.test() {
			gcStart(t)
		}
	}

	if raceenabled && noscan && dataSize < maxTinySize {
		// Pad tinysize allocations so they are aligned with the end
		// of the tinyalloc region. This ensures that any arithmetic
		// that goes off the top end of the object will be detectable
		// by checkptr (issue 38872).
		// Note that we disable tinyalloc when raceenabled for this to work.
		// TODO: This padding is only performed when the race detector
		// is enabled. It would be nice to enable it if any package
		// was compiled with checkptr, but there's no easy way to
		// detect that (especially at compile time).
		// TODO: enable this padding for all allocations, not just
		// tinyalloc ones. It's tricky because of pointer maps.
		// Maybe just all noscan objects?
		x = add(x, size-dataSize)
	}

	return x
}
```

### 微小对象

Go 语言将小于 16 字节的对象划分为微小对象。划分微小对象的主要目的是处理极小的字符串和独立的转义变量。

微小对象会被放入 class 为 2 的 span 中，我们已经知道，在 class 为 2 的 span 中元素的大小为 16 字节。首先对微小对象按照 2、4、8 的规则进行字节对齐。例如，字节为 1 的元素会被分配 2 字节，字节为 7 的元素会被分配 8 字节。

查看之前分配的元素中是否有空余的空间，图 18-4 所示为微小对象分配示意图。如果当前对象要分配 8 字节，并且正在分配的元素可以容纳 8 字节，则返回 tiny+offset 的地址，意味着当前地址往后 8 字节都是可以被分配的。

![](../../../assets/images/docs/internal/memory/object_malloc/图18-4%20微小对象分配.png)

如图 18-5 所示，分配完成后 offset 的位置也需要相应增加，为下一次分配做准备。

![](../../../assets/images/docs/internal/memory/object_malloc/图18-5%20tiny%20offset代表当前已分配内存的偏移量.png)

如果当前要分配的元素空间不够，将尝试从 mcache 中查找 span 中下一个可用的元素。因此，tiny 分配的第一步是尝试利用分配过的前一个元素的空间，达到节约内存的目的。

### 小对象分配

当对象不属于微小对象时，在内存分配时会继续判断其是否为小对象，**小对象指小于 32KB 的对象**。Go 语言会计算小对象对应哪一个等级的 span，并在指定等级的 span 中查找。

此后和微小对象的分配一样，小对象分配经历 **mcache → mcentral → mheap 位图查找→ mheap 基数树查找→操作系统分配**的过程。

GC 性能与对象数量负相关，**对象越多 GC 性能越差**，对程序影响越大。

所以 GC 性能优化的思路之一就是**减少对象分配个数**，比如对象复用或使用大对象组合多个小对象等等。

通常小对象过多会导致三色法消耗过多的 CPU。

### 大对象分配

**大对象指大于 32KB 的对象**，内存分配时不与 mcache 和 mcentral 沟通，直接通过 mheap 进行分配。大对象分配经历 mheap 基数树查找→操作系统分配的过程。每个大对象都是一个特殊的 span，其 class 为 0。

## 对象分配过程

[对象分配过程](malloc.md)

```go

```
