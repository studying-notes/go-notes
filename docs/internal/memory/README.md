---
date: 2022-10-10T20:42:48+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "内存分配管理"  # 文章标题
url:  "posts/go/docs/internal/memory/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [内存分配全局](#内存分配全局)
	- [span 与元素](#span-与元素)
	- [三级对象管理](#三级对象管理)
	- [四级内存块管理](#四级内存块管理)
- [对象分配](#对象分配)
	- [微小对象](#微小对象)
	- [mcache 缓存位图](#mcache-缓存位图)
	- [mcentral 遍历 span](#mcentral-遍历-span)
	- [mheap 缓存查找](#mheap-缓存查找)
	- [mheap 基数树查找](#mheap-基数树查找)
	- [操作系统内存申请](#操作系统内存申请)
	- [小对象分配](#小对象分配)
	- [大对象分配](#大对象分配)

## 内存分配全局

Go 语言采用现代内存分配 TCMalloc 算法的思想来进行内存分配，将对象分为微小对象、小对象、大对象，使用三级管理结构 mcache、mcentral、mheap 用于管理、缓存加速 span 对象的访问和分配，使用精准的位图管理已分配的和未分配的对象及对象的大小。

Go 语言运行时依靠细微的对象切割、极致的多级缓存、精准的位图管理实现了对内存的精细化管理以及快速的内存访问，同时减少了内存的碎片。

### span 与元素

Go 语言将内存分成了大大小小 67 个级别的 span，其中，0 级代表特殊的大对象，其大小是不固定的。当具体的对象需要分配内存时，并不是直接分配 span，而是分配不同级别的 span 中的元素。因此 span 的级别不是以每个 span 的大小为依据的，而是以 span 中元素的大小为依据的。

根据对象大小，划分了一系列级别 class，每个 class 都代表一个固定大小的对象，以及每个 span 的大小。如下表所示：

```go
// class  bytes/obj  bytes/span  objects  waste bytes
//     1          8        8192     1024            0
//     2         16        8192      512            0
//     3         32        8192      256            0
//     4         48        8192      170           32
//     5         64        8192      128            0
//     6         80        8192      102           32
//     7         96        8192       85           32
//     8        112        8192       73           16
//     9        128        8192       64            0
//    10        144        8192       56          128
//    11        160        8192       51           32
//    12        176        8192       46           96
//    13        192        8192       42          128
//    14        208        8192       39           80
//    15        224        8192       36          128
//    16        240        8192       34           32
//    17        256        8192       32            0
//    18        288        8192       28          128
//    19        320        8192       25          192
//    20        352        8192       23           96
//    21        384        8192       21          128
//    22        416        8192       19          288
//    23        448        8192       18          128
//    24        480        8192       17           32
//    25        512        8192       16            0
//    26        576        8192       14          128
//    27        640        8192       12          512
//    28        704        8192       11          448
//    29        768        8192       10          512
//    30        896        8192        9          128
//    31       1024        8192        8            0
//    32       1152        8192        7          128
//    33       1280        8192        6          512
//    34       1408       16384       11          896
//    35       1536        8192        5          512
//    36       1792       16384        9          256
//    37       2048        8192        4            0
//    38       2304       16384        7          256
//    39       2688        8192        3          128
//    40       3072       24576        8            0
//    41       3200       16384        5          384
//    42       3456       24576        7          384
//    43       4096        8192        2            0
//    44       4864       24576        5          256
//    45       5376       16384        3          256
//    46       6144       24576        4            0
//    47       6528       32768        5          128
//    48       6784       40960        6          256
//    49       6912       49152        7          768
//    50       8192        8192        1            0
//    51       9472       57344        6          512
//    52       9728       49152        5          512
//    53      10240       40960        4            0
//    54      10880       32768        3          128
//    55      12288       24576        2            0
//    56      13568       40960        3          256
//    57      14336       57344        4            0
//    58      16384       16384        1            0
//    59      18432       73728        4            0
//    60      19072       57344        3          128
//    61      20480       40960        2            0
//    62      21760       65536        3          256
//    63      24576       24576        1            0
//    64      27264       81920        3          128
//    65      28672       57344        2            0
//    66      32768       32768        1            0
```

上表中每列含义如下：

- class：class ID，每个 span 结构中都有一个 class ID, 表示该 span 可处理的对象类型
- bytes/obj：该 class 代表对象的字节数
- bytes/span：每个 span 占用堆的字节数，也即页数*页大小
- objects : 每个 span 可分配的对象个数，也即 (bytes/spans)/(bytes/obj)
- waste bytes : 每个 span 产生的内存碎片，也即 (bytes/spans)%(bytes/obj)

上表可见最大的对象是 32K 大小，超过 32K 大小的由特殊的 class 表示，该 class ID 为 0，每个 class 只包含一个对象。

span 的大小虽然不固定，但其是 8KB 或更大的连续内存区域。

每个具体的对象在分配时都需要对齐到指定的大小，例如分配 17 字节的对象，会对应分配到比 17 字节大并最接近它的元素级别，即第 3 级，这导致最终分配了 32 字节。因此，这种分配方式会不可避免地带来内存的浪费。

### 三级对象管理

为了能够方便地对 span 进行管理，加速 span 对象的访问和分配，Go 语言采取了三级管理结构，分别为 mcache、mcentral、mheap。

Go 语言采用了现代 TCMalloc 内存分配算法的思想，每个逻辑处理器 P 都存储了一个本地 span 缓存，称作 mcache。如果协程需要内存可以直接从 mcache 中获取，由于在同一时间只有一个协程运行在逻辑处理器 P 上，所以中间不需要加锁。mcache 包含所有大小规格的 mspan，但是每种规格大小只包含一个。除 class() 外，mcache 的 span 都来自 mcentral。

- mcentral 是被所有逻辑处理器 P 共享的。
- mcentral 对象收集所有给定规格大小的 span。每个 mcentral 都包含两个 mspan 的链表：empty mspanList 表示没有空闲对象或 span 已经被 mcache 缓存的 span 链表，nonempty mspanList 表示有空闲对象的 span 链表。

做这种区分是为了更快地分配 span 到 mcache 中。图 18-1 为三级内存对象管理的示意图。除了级别 0，每个级别的 span 都会有一个 mcentral 用于管理 span 链表。而所有级别的这些 mcentral，其实都是一个数组，由 mheap 进行管理。

![](../../../assets/images/docs/internal/memory/README/图18-1%20三级内存对象管理.png)

mheap 的作用不只是管理 central，大对象也会直接通过 mheap 进行分配。如图 18-2 所示，mheap 实现了对虚拟内存线性地址空间的精准管理，建立了 span 与具体线性地址空间的联系，保存了分配的位图信息，是管理内存的最核心单元。后面还会看到，堆区的内存被分成了 HeapArea 大小进行管理。对 Heap 进行的操作必须全局加锁，而 mcache、mcentral 可以被看作某种形式的缓存。

![](../../../assets/images/docs/internal/memory/README/图18-2%20mheap管理虚拟内存线性地址空间.png)

### 四级内存块管理

根据对象的大小，Go 语言将堆内存分成了图 18-3 所示的 HeapArea、chunk、span 与 page 4 种内存块进行管理。

其中，HeapArea 内存块最大，其大小与平台相关，在 UNIX 64 位操作系统中占据 64MB。chunk 占据了 512KB，span 根据级别大小的不同而不同，但必须是 page 的倍数。而 1 个 page 占据 8KB。

不同的内存块用于不同的场景，便于高效地对内存进行管理。

![](../../../assets/images/docs/internal/memory/README/图18-3%20内存块管理结构.png)

## 对象分配

在运行时分配对象的逻辑主要位于 mallocgc 函数中 `src/runtime/malloc.go`，malloc 代表分配，gc 代表垃圾回收（GC），此函数除了分配内存还会为垃圾回收做一些位图标记工作。

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

	if debug.malloc {
		if debug.sbrk != 0 {
			align := uintptr(16)
			if typ != nil {
				// TODO(austin): This should be just
				//   align = uintptr(typ.align)
				// but that's only 4 on 32-bit platforms,
				// even if there's a uint64 field in typ (see #599).
				// This causes 64-bit atomic accesses to panic.
				// Hence, we use stricter alignment that matches
				// the normal allocator better.
				if size&7 == 0 {
					align = 8
				} else if size&3 == 0 {
					align = 4
				} else if size&1 == 0 {
					align = 2
				} else {
					align = 1
				}
			}
			return persistentalloc(size, align, &memstats.other_sys)
		}

		if inittrace.active && inittrace.id == getg().goid {
			// Init functions are executed sequentially in a single goroutine.
			inittrace.allocs += 1
		}
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

内存分配时，将对象按照大小不同划分为微小（tiny）对象、小对象、大对象。微小对象的分配流程最长，逻辑链路最复杂。

### 微小对象

Go 语言将小于 16 字节的对象划分为微小对象。划分微小对象的主要目的是处理极小的字符串和独立的转义变量。对 json 的基准测试表明，使用微小对象减少了 12% 的分配次数和 20% 的堆大小。

微小对象会被放入 class 为 2 的 span 中，我们已经知道，在 class 为 2 的 span 中元素的大小为 16 字节。首先对微小对象按照 2、4、8 的规则进行字节对齐。例如，字节为 1 的元素会被分配 2 字节，字节为 7 的元素会被分配 8 字节。

查看之前分配的元素中是否有空余的空间，图 18-4 所示为微小对象分配示意图。如果当前对象要分配 8 字节，并且正在分配的元素可以容纳 8 字节，则返回 tiny+offset 的地址，意味着当前地址往后 8 字节都是可以被分配的。

![](../../../assets/images/docs/internal/memory/README/图18-4%20微小对象分配.png)

如图 18-5 所示，分配完成后 offset 的位置也需要相应增加，为下一次分配做准备。

![](../../../assets/images/docs/internal/memory/README/图18-5%20tiny%20offset代表当前已分配内存的偏移量.png)

如果当前要分配的元素空间不够，将尝试从 mcache 中查找 span 中下一个可用的元素。因此，tiny 分配的第一步是尝试利用分配过的前一个元素的空间，达到节约内存的目的。

### mcache 缓存位图

在查找空闲元素空间时，首先需要从 mcache 中找到对应级别的 mspan，mspan 中拥有 allocCache 字段，其作为一个位图，用于标记 span 中的元素是否被分配。由于 allocCache 元素为 uint64，因此其最多一次缓存 64 字节。

```go
// nextFreeFast returns the next free object if one is quickly available.
// Otherwise it returns 0.
func nextFreeFast(s *mspan) gclinkptr {
	theBit := sys.Ctz64(s.allocCache) // Is there a free object in the allocCache?
	if theBit < 64 {
		result := s.freeindex + uintptr(theBit)
		if result < s.nelems {
			freeidx := result + 1
			if freeidx%64 == 0 && freeidx != s.nelems {
				return 0
			}
			s.allocCache >>= uint(theBit + 1)
			s.freeindex = freeidx
			s.allocCount++
			return gclinkptr(result*s.elemsize + s.base())
		}
	}
	return 0
}
```

allocCache 使用图 18-6 中的小端模式标记 span 中的元素是否被分配。allocCache 中的最后 1 bit 对应的是 span 中的第 1 个元素是否被分配。当 bit 位为 1 时代表当前对应的 span 中的元素已经被分配。

![](../../../assets/images/docs/internal/memory/README/图18-6%20allocCache位图标记span中的元素是否被分配.png)

有时候，span 中元素的个数大于 64，因此需要专门有一个字段 freeindex 标识当前 span 中的元素被分配到了哪里。如图 18-7 所示，span 中小于 freeindex 序号的元素都已经被分配了，将从 freeindex 开始继续分配。

![](../../../assets/images/docs/internal/memory/README/图18-7%20freeindex之前的元素都已被分配.png)

因此，只要从 allocCache 开始找到哪一位为 0 即可。假如 X 位为 0，那么 X+freeindex 为当前 span 中可用的元素序号。当 allocCache 中的 bit 位全部被标记为 1 后，需要移动 freeindex，并更新 allocCache，一直到 span 中元素的末尾为止。

### mcentral 遍历 span

如果当前的 span 中没有可以使用的元素，这时就需要从 mcentral 中加锁查找。之前介绍过，mcentral 中有两种类型的 span 链表，分别是有空闲元素的 nonempty 链表和没有空闲元素的 empty 链表。在 mcentral 查找时，会分别遍历这两个链表，查找是否有可用的 span。

既然是没有空闲元素的 empty 链表，为什么还需要遍历呢？这是由于可能有些 span 虽然被垃圾回收器标记为空闲了，但是还没有来得及清理，这些 span 在清扫后仍然是可以使用的，因此需要遍历。

`src/runtime/mcentral.go`

```go
// Allocate a span to use in an mcache.
func (c *mcentral) cacheSpan() *mspan {
	// Deduct credit for this span allocation and sweep if necessary.
	spanBytes := uintptr(class_to_allocnpages[c.spanclass.sizeclass()]) * _PageSize
	deductSweepCredit(spanBytes, 0)

	traceDone := false
	if trace.enabled {
		traceGCSweepStart()
	}

	// If we sweep spanBudget spans without finding any free
	// space, just allocate a fresh span. This limits the amount
	// of time we can spend trying to find free space and
	// amortizes the cost of small object sweeping over the
	// benefit of having a full free span to allocate from. By
	// setting this to 100, we limit the space overhead to 1%.
	//
	// TODO(austin,mknyszek): This still has bad worst-case
	// throughput. For example, this could find just one free slot
	// on the 100th swept span. That limits allocation latency, but
	// still has very poor throughput. We could instead keep a
	// running free-to-used budget and switch to fresh span
	// allocation if the budget runs low.
	spanBudget := 100

	var s *mspan
	var sl sweepLocker

	// Try partial swept spans first.
	sg := mheap_.sweepgen
	if s = c.partialSwept(sg).pop(); s != nil {
		goto havespan
	}

	sl = sweep.active.begin()
	if sl.valid {
		// Now try partial unswept spans.
		for ; spanBudget >= 0; spanBudget-- {
			s = c.partialUnswept(sg).pop()
			if s == nil {
				break
			}
			if s, ok := sl.tryAcquire(s); ok {
				// We got ownership of the span, so let's sweep it and use it.
				s.sweep(true)
				sweep.active.end(sl)
				goto havespan
			}
			// We failed to get ownership of the span, which means it's being or
			// has been swept by an asynchronous sweeper that just couldn't remove it
			// from the unswept list. That sweeper took ownership of the span and
			// responsibility for either freeing it to the heap or putting it on the
			// right swept list. Either way, we should just ignore it (and it's unsafe
			// for us to do anything else).
		}
		// Now try full unswept spans, sweeping them and putting them into the
		// right list if we fail to get a span.
		for ; spanBudget >= 0; spanBudget-- {
			s = c.fullUnswept(sg).pop()
			if s == nil {
				break
			}
			if s, ok := sl.tryAcquire(s); ok {
				// We got ownership of the span, so let's sweep it.
				s.sweep(true)
				// Check if there's any free space.
				freeIndex := s.nextFreeIndex()
				if freeIndex != s.nelems {
					s.freeindex = freeIndex
					sweep.active.end(sl)
					goto havespan
				}
				// Add it to the swept list, because sweeping didn't give us any free space.
				c.fullSwept(sg).push(s.mspan)
			}
			// See comment for partial unswept spans.
		}
		sweep.active.end(sl)
	}
	if trace.enabled {
		traceGCSweepDone()
		traceDone = true
	}

	// We failed to get a span from the mcentral so get one from mheap.
	s = c.grow()
	if s == nil {
		return nil
	}

	// At this point s is a span that should have free slots.
havespan:
	if trace.enabled && !traceDone {
		traceGCSweepDone()
	}
	n := int(s.nelems) - int(s.allocCount)
	if n == 0 || s.freeindex == s.nelems || uintptr(s.allocCount) == s.nelems {
		throw("span has no free objects")
	}
	freeByteBase := s.freeindex &^ (64 - 1)
	whichByte := freeByteBase / 8
	// Init alloc bits cache.
	s.refillAllocCache(whichByte)

	// Adjust the allocCache so that s.freeindex corresponds to the low bit in
	// s.allocCache.
	s.allocCache >>= s.freeindex % 64

	return s
}
```

图 18-8 为查找 mcentral 中可用 span 并分配到 mcache 中的示意图。如果在 mcentral 中查找到有空闲元素的 span，则将其赋值到 mcache 中，并更新 allocCache，同时需要将 span 添加到 mcentral 的 empty 链表中去。

![](../../../assets/images/docs/internal/memory/README/图18-8%20查找mcentral中的可用span并分配到mcache中.png)

### mheap 缓存查找

如果在 mcentral 中找不到可以使用的 span，就需要在 mheap 中查找。Go 1.12 采用 treap 结构进行内存管理，treap 是一种引入了随机数的二叉搜索树，其实现简单，引入的随机数及必要时的旋转保证了比较好的平衡性。Michael Knyszek 提出，这种方式有扩展性的问题，由于这棵树是 mheap 管理的，所以在操作它时需要维持一个 lock。这在密集的对象分配及逻辑处理器 P 过多时，会导致更长的等待时间。Michael Knyszek 提出用 bitmap 来管理内存页，因此在 Go 1.14 之后，我们会看到每个逻辑处理器 P 中都维护了一份 page cache，这就是现在 Go 语言实现的方式。

```go
// pageCache represents a per-p cache of pages the allocator can
// allocate from without a lock. More specifically, it represents
// a pageCachePages*pageSize chunk of memory with 0 or more free
// pages in it.
type pageCache struct {
	base  uintptr // base address of the chunk
	cache uint64  // 64-bit bitmap representing free pages (1 means free)
	scav  uint64  // 64-bit bitmap representing scavenged pages (1 means scavenged)
}
```

mheap 会首先查找每个逻辑处理器 P 中 pageCache 字段的 cache。如图 18-9 所示，cache 也是一个位图，每一位都代表了一个 page（8 KB）。由于 cache 为 uint64，因此一共可以提供 64×8 = 512KB 的连续虚拟内存。在 cache 中，1 代表未分配的内存，0 代表已分配的内存。base 代表该虚拟内存的基地址。当需要分配的内存小于 512/4 = 128KB 时，需要首先从 cache 中分配。

![](../../../assets/images/docs/internal/memory/README/图18-9%20cache位图标记内存缓存页是否分配.png)

例如要分配 n pages，就需要查找 cache 中是否有连续 n 个为 1 的位。如果存在，则说明在缓存中查找到了合适的内存，用于构建 span。

### mheap 基数树查找

如果要分配的 page 过大或者在逻辑处理器 P 的 cache 中没有找到可用的 page，就需要对 mheap 加锁，并在整个 mheap 管理的虚拟地址空间的位图中查找是否有可用的 page，这涉及 Go 语言对线性地址空间的位图管理。

管理线性地址空间的位图结构叫作基数树（radix tree），其结构如图 18-10 所示。该结构和一般的基数树结构不太一样，会有这个名字很大一部分是由于父节点包含了子节点的若干信息。

![](../../../assets/images/docs/internal/memory/README/图18-10%20内存管理基数树结构.png)

该树中的每个节点都对应一个 pallocSum，最底层的叶子节点对应的 pallocSum 包含一个 chunk 的信息（512×8KB），除叶子节点外的节点都包含连续 8 个子节点的内存信息。例如，倒数第 2 层的节点包含连续 8 个叶子节点（即 8×chunk）的内存信息。因此，越上层的节点对应的内存越多。

![](../../../assets/images/docs/internal/memory/README/图18-11%20pallocSum的内部结构.png)

pallocSum 是一个简单的 uint64，分为开头（start）、中间（max）、末尾（end）3 部分，其结构如图 18-11 所示。pallocSum 的开头与末尾部分各占 21bit，中间部分占 22bit，它们分别包含了这个区域中连续空闲内存页的信息，包括开头有多少连续内存页，最多有多少连续内存页，末尾有多少连续内存页。对于最顶层的节点，由于其 max 位为 22bit，因此一棵完整的基数树最多代表 2^21 pages = 16GB 内存。

不需要每一次查找都从根节点开始。在 Go 语言中，存储了一个特别的字段 searchAddr，顾名思义是用于搜索可用内存的。如图 18-12 所示，利用 searchAddr 可以加速内存查找。searchAddr 有一个重要的设定是它前面的地址一定是已经分配过的，因此在查找时，只需要向 searchAddr 地址的后方查找即可跳过已经查找的节点，减少查找的时间。

![](../../../assets/images/docs/internal/memory/README/图18-12%20利用searchAddr加速内存查找.png)

在第 1 次查找时，会从当前 searchAddr 的 chunk 块中查找是否有对应大小的连续空间，这种优化主要针对比较小的内存（至少小于 512KB）分配。Go 语言对于内存有非常精细的管理，chunk 块的每个 page（8 KB）都有位图表明其是否已经被分配。

每个 chunk 都有一个 pallocData 结构，其中 pallocBits 管理其分配的位图。pallocBits 是 uint64，有 8 字节，由于其每一位对应一个 page，因此 pallocBits 一共对应 64×8 = 512KB，恰好是一个 chunk 块的大小。位图的对应方式和之前是一样的。

```go
// pallocData encapsulates pallocBits and a bitmap for
// whether or not a given page is scavenged in a single
// structure. It's effectively a pallocBits with
// additional functionality.
//
// Update the comment on (*pageAlloc).chunks should this
// structure change.
type pallocData struct {
	pallocBits
	scavenged pageBits
}

// pageBits is a bitmap representing one bit per page in a palloc chunk.
type pageBits [pallocChunkPages / 64]uint64

const(
	// The size of a bitmap chunk, i.e. the amount of bits (that is, pages) to consider
	// in the bitmap at once.
	pallocChunkPages    = 1 << logPallocChunkPages
	pallocChunkBytes    = pallocChunkPages * pageSize
	logPallocChunkPages = 9
	logPallocChunkBytes = logPallocChunkPages + pageShift
)
```

当内存分配过大或者当前 chunk 块没有连续的 npages 空间时，需要到基数树中从上到下进行查找。基数树有一个特性——要分配的内存越大，它能够越快地查找到当前的基数树中是否有连续的满足需求的空间。

在查找基数树的过程中，需要从上到下、从左到右地查找每个节点是否符合要求。先计算 pallocSum 的开头有多少连续的内存空间，如果大于或等于 npages，则说明找到了可用的空间和地址。如果小于 npages，则会计算 pallocSum 字段的 max，即中间有多少连续的内存空间。如果 max 大于或等于 npages，那么需要继续向基数树当前节点对应的下一级查找，原因在于，max 大于 npages，表明当前一定有连续的空间大于或等于 npages，但是并不知道具体在哪一个位置，必须查找下一级才能找到可用的地址。如果 max 也不满足，那么是不是就不满足了呢？不一定，如图 18-13 所示，有可能两个节点可以合并起来组成一个更大的连续空间。因此还需要将当前 pallocSum 计算的 end 与后一个节点的 start 加起来查看是否能够组合成大于 npages 的连续空间。

![](../../../assets/images/docs/internal/memory/README/图18-13%20更大的可用内存可能跨越了多个pallocSum.png)

每一次从基数树中查找到内存，或者事后从操作系统分配到内存时，都需要更新基数树中每个节点的 pallocSum。

### 操作系统内存申请

当在基数树中查找不到可用的连续内存时，需要从操作系统中获取内存。从操作系统获取内存的代码是平台独立的。

```go
// sysReserve transitions a memory region from None to Reserved. It reserves
// address space in such a way that it would cause a fatal fault upon access
// (either via permissions or not committing the memory). Such a reservation is
// thus never backed by physical memory.
//
// If the pointer passed to it is non-nil, the caller wants the
// reservation there, but sysReserve can still choose another
// location if that one is unavailable.
//
// NOTE: sysReserve returns OS-aligned memory, but the heap allocator
// may use larger alignment, so the caller must be careful to realign the
// memory obtained by sysReserve.
func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer {
	return sysReserveOS(v, n)
}
```

Go 语言规定，每一次向操作系统申请的内存大小必须为 heapArena 的倍数。heapArena 是和平台有关的内存大小，在 64 位 UNIX 操作系统中，其大小为 64MB。这意味着即便需要的内存很小，最终也至少要向操作系统申请 64MB 内存。多申请的内存可以用于下次分配。

Go 语言中对于 heapArena 有精准的管理，精准到每个指针大小的内存信息，每个 page 对应的 mspan 信息都有记录。

```go
// A heapArena stores metadata for a heap arena. heapArenas are stored
// outside of the Go heap and accessed via the mheap_.arenas index.
type heapArena struct {
	_ sys.NotInHeap

	// bitmap stores the pointer/scalar bitmap for the words in
	// this arena. See mbitmap.go for a description.
	// This array uses 1 bit per word of heap, or 1.6% of the heap size (for 64-bit).
	bitmap [heapArenaBitmapWords]uintptr

	// If the ith bit of noMorePtrs is true, then there are no more
	// pointers for the object containing the word described by the
	// high bit of bitmap[i].
	// In that case, bitmap[i+1], ... must be zero until the start
	// of the next object.
	// We never operate on these entries using bit-parallel techniques,
	// so it is ok if they are small. Also, they can't be bigger than
	// uint16 because at that size a single noMorePtrs entry
	// represents 8K of memory, the minimum size of a span. Any larger
	// and we'd have to worry about concurrent updates.
	// This array uses 1 bit per word of bitmap, or .024% of the heap size (for 64-bit).
	noMorePtrs [heapArenaBitmapWords / 8]uint8

	// spans maps from virtual address page ID within this arena to *mspan.
	// For allocated spans, their pages map to the span itself.
	// For free spans, only the lowest and highest pages map to the span itself.
	// Internal pages map to an arbitrary span.
	// For pages that have never been allocated, spans entries are nil.
	//
	// Modifications are protected by mheap.lock. Reads can be
	// performed without locking, but ONLY from indexes that are
	// known to contain in-use or stack spans. This means there
	// must not be a safe-point between establishing that an
	// address is live and looking it up in the spans array.
	spans [pagesPerArena]*mspan

	// pageInUse is a bitmap that indicates which spans are in
	// state mSpanInUse. This bitmap is indexed by page number,
	// but only the bit corresponding to the first page in each
	// span is used.
	//
	// Reads and writes are atomic.
	pageInUse [pagesPerArena / 8]uint8

	// pageMarks is a bitmap that indicates which spans have any
	// marked objects on them. Like pageInUse, only the bit
	// corresponding to the first page in each span is used.
	//
	// Writes are done atomically during marking. Reads are
	// non-atomic and lock-free since they only occur during
	// sweeping (and hence never race with writes).
	//
	// This is used to quickly find whole spans that can be freed.
	//
	// TODO(austin): It would be nice if this was uint64 for
	// faster scanning, but we don't have 64-bit atomic bit
	// operations.
	pageMarks [pagesPerArena / 8]uint8

	// pageSpecials is a bitmap that indicates which spans have
	// specials (finalizers or other). Like pageInUse, only the bit
	// corresponding to the first page in each span is used.
	//
	// Writes are done atomically whenever a special is added to
	// a span and whenever the last special is removed from a span.
	// Reads are done atomically to find spans containing specials
	// during marking.
	pageSpecials [pagesPerArena / 8]uint8

	// checkmarks stores the debug.gccheckmark state. It is only
	// used if debug.gccheckmark > 0.
	checkmarks *checkmarksMap

	// zeroedBase marks the first byte of the first page in this
	// arena which hasn't been used yet and is therefore already
	// zero. zeroedBase is relative to the arena base.
	// Increases monotonically until it hits heapArenaBytes.
	//
	// This field is sufficient to determine if an allocation
	// needs to be zeroed because the page allocator follows an
	// address-ordered first-fit policy.
	//
	// Read atomically and written with an atomic CAS.
	zeroedBase uintptr
}
```

### 小对象分配

当对象不属于微小对象时，在内存分配时会继续判断其是否为小对象，小对象指小于 32KB 的对象。Go 语言会计算小对象对应哪一个等级的 span，并在指定等级的 span 中查找。

此后和微小对象的分配一样，小对象分配经历 mcache → mcentral → mheap 位图查找→ mheap 基数树查找→操作系统分配的过程。

### 大对象分配

大对象指大于 32KB 的对象，内存分配时不与 mcache 和 mcentral 沟通，直接通过 mheap 进行分配。大对象分配经历 mheap 基数树查找→操作系统分配的过程。每个大对象都是一个特殊的 span，其 class 为 0。

```go

```
