---
date: 2022-10-14T10:27:33+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "对象分配过程"  # 文章标题
url:  "posts/go/docs/internal/memory/malloc"  # 设置网页永久链接
tags: [ "Go", "find" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## mcache 缓存位图

在查找空闲元素空间时，**首先需要从 mcache 中找到对应级别的 mspan**，mspan 中拥有 allocCache 字段，其作为一个位图，用于标记 span 中的元素是否被分配。由于 allocCache 元素为 uint64，因此其最多一次缓存 64 字节。

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

![](../../../assets/images/docs/internal/memory/malloc/图18-6%20allocCache位图标记span中的元素是否被分配.png)

有时候，span 中元素的个数大于 64，因此需要专门有一个字段 freeindex 标识当前 span 中的元素被分配到了哪里。如图 18-7 所示，span 中小于 freeindex 序号的元素都已经被分配了，将从 freeindex 开始继续分配。

![](../../../assets/images/docs/internal/memory/malloc/图18-7%20freeindex之前的元素都已被分配.png)

因此，只要从 allocCache 开始找到哪一位为 0 即可。假如 X 位为 0，那么 X+freeindex 为当前 span 中可用的元素序号。当 allocCache 中的 bit 位全部被标记为 1 后，需要移动 freeindex，并更新 allocCache，一直到 span 中元素的末尾为止。

## mcentral 遍历 span

如果当前的 span 中没有可以使用的元素，这时就需要从 mcentral 中**加锁查找**。之前介绍过，mcentral 中有两种类型的 span 链表，分别是有空闲元素的 nonempty 链表和没有空闲元素的 empty 链表。在 mcentral 查找时，会分**别遍历这两个链表**，查找是否有可用的 span。

既然是没有空闲元素的 empty 链表，为什么还需要遍历呢？这是由于可能**有些 span 虽然被垃圾回收器标记为空闲了，但是还没有来得及清理，这些 span 在清扫后仍然是可以使用的，因此需要遍历**。

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

图 18-8 为查找 mcentral 中可用 span 并分配到 mcache 中的示意图。**如果在 mcentral 中查找到有空闲元素的 span，则将其赋值到 mcache 中，并更新 allocCache，同时需要将 span 添加到 mcentral 的 empty 链表中去。**

![](../../../assets/images/docs/internal/memory/malloc/图18-8%20查找mcentral中的可用span并分配到mcache中.png)

```go
// refill acquires a new span of span class spc for c. This span will
// have at least one free object. The current span in c must be full.
//
// Must run in a non-preemptible context since otherwise the owner of
// c could change.
func (c *mcache) refill(spc spanClass) {
	// Return the current cached span to the central lists.
	s := c.alloc[spc]

	if uintptr(s.allocCount) != s.nelems {
		throw("refill of span with free space remaining")
	}
	if s != &emptymspan {
		// Mark this span as no longer cached.
		if s.sweepgen != mheap_.sweepgen+3 {
			throw("bad sweepgen in refill")
		}
		mheap_.central[spc].mcentral.uncacheSpan(s)

		// Count up how many slots were used and record it.
		stats := memstats.heapStats.acquire()
		slotsUsed := int64(s.allocCount) - int64(s.allocCountBeforeCache)
		atomic.Xadd64(&stats.smallAllocCount[spc.sizeclass()], slotsUsed)

		// Flush tinyAllocs.
		if spc == tinySpanClass {
			atomic.Xadd64(&stats.tinyAllocCount, int64(c.tinyAllocs))
			c.tinyAllocs = 0
		}
		memstats.heapStats.release()

		// Count the allocs in inconsistent, internal stats.
		bytesAllocated := slotsUsed * int64(s.elemsize)
		gcController.totalAlloc.Add(bytesAllocated)

		// Clear the second allocCount just to be safe.
		s.allocCountBeforeCache = 0
	}

	// Get a new cached span from the central lists.
	s = mheap_.central[spc].mcentral.cacheSpan()
	if s == nil {
		throw("out of memory")
	}

	if uintptr(s.allocCount) == s.nelems {
		throw("span has no free space")
	}

	// Indicate that this span is cached and prevent asynchronous
	// sweeping in the next sweep phase.
	s.sweepgen = mheap_.sweepgen + 3

	// Store the current alloc count for accounting later.
	s.allocCountBeforeCache = s.allocCount

	// Update heapLive and flush scanAlloc.
	//
	// We have not yet allocated anything new into the span, but we
	// assume that all of its slots will get used, so this makes
	// heapLive an overestimate.
	//
	// When the span gets uncached, we'll fix up this overestimate
	// if necessary (see releaseAll).
	//
	// We pick an overestimate here because an underestimate leads
	// the pacer to believe that it's in better shape than it is,
	// which appears to lead to more memory used. See #53738 for
	// more details.
	usedBytes := uintptr(s.allocCount) * s.elemsize
	gcController.update(int64(s.npages*pageSize)-int64(usedBytes), int64(c.scanAlloc))
	c.scanAlloc = 0

	c.alloc[spc] = s
}
```

## mheap 缓存查找

如果在 mcentral 中找不到可以使用的 span，就需要在 mheap 中查找。

每个逻辑处理器 P 中都维护了一份 page cache。

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

![](../../../assets/images/docs/internal/memory/malloc/图18-9%20cache位图标记内存缓存页是否分配.png)

例如要分配 n pages，就需要查找 cache 中是否有连续 n 个为 1 的位。如果存在，则说明在缓存中查找到了合适的内存，用于构建 span。

## mheap 基数树查找

如果要分配的 page 过大或者在逻辑处理器 P 的 cache 中没有找到可用的 page，就需要对 mheap 加锁，并在整个 mheap 管理的虚拟地址空间的位图中查找是否有可用的 page，这涉及 Go 语言对线性地址空间的位图管理。

管理线性地址空间的位图结构叫作基数树（radix tree），其结构如图 18-10 所示。该结构和一般的基数树结构不太一样，会有这个名字很大一部分是由于父节点包含了子节点的若干信息。

![](../../../assets/images/docs/internal/memory/malloc/图18-10%20内存管理基数树结构.png)

该树中的每个节点都对应一个 pallocSum，最底层的叶子节点对应的 pallocSum 包含一个 chunk 的信息（512×8KB），除叶子节点外的节点都包含连续 8 个子节点的内存信息。例如，倒数第 2 层的节点包含连续 8 个叶子节点（即 8×chunk）的内存信息。因此，越上层的节点对应的内存越多。

![](../../../assets/images/docs/internal/memory/malloc/图18-11%20pallocSum的内部结构.png)

pallocSum 是一个简单的 uint64，分为开头（start）、中间（max）、末尾（end）3 部分，其结构如图 18-11 所示。pallocSum 的开头与末尾部分各占 21bit，中间部分占 22bit，它们分别包含了这个区域中连续空闲内存页的信息，包括开头有多少连续内存页，最多有多少连续内存页，末尾有多少连续内存页。对于最顶层的节点，由于其 max 位为 22bit，因此一棵完整的基数树最多代表 2^21 pages = 16GB 内存。

不需要每一次查找都从根节点开始。在 Go 语言中，存储了一个特别的字段 searchAddr，顾名思义是用于搜索可用内存的。如图 18-12 所示，利用 searchAddr 可以加速内存查找。searchAddr 有一个重要的设定是它前面的地址一定是已经分配过的，因此在查找时，只需要向 searchAddr 地址的后方查找即可跳过已经查找的节点，减少查找的时间。

![](../../../assets/images/docs/internal/memory/malloc/图18-12%20利用searchAddr加速内存查找.png)

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

![](../../../assets/images/docs/internal/memory/malloc/图18-13%20更大的可用内存可能跨越了多个pallocSum.png)

每一次从基数树中查找到内存，或者事后从操作系统分配到内存时，都需要更新基数树中每个节点的 pallocSum。

## 操作系统内存申请

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

```go

```
