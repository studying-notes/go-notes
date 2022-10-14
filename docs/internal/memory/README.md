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

## 前言

C 语言程序标准库通过 malloc() 方法动态申请内存，其中内存分配器使用的是 glibc 提供的 ptmalloc2。

除了 glibc，业界比较出名的内存分配器有 Google 的 tcmalloc 和 Facebook 的 jemalloc。二者在避免内存碎片和性能上均比 glibc 有比较大的优势，在多线程环境中效果更明显。

## 内存分配全局

Go 语言吸收了 TCMalloc 的思想来进行内存分配，将对象分为微小对象、小对象、大对象，使用**三级管理结构** mcache、mcentral、mheap **管理、缓存加速 span 对象的访问和分配**，使用**位图****管理已分配的和未分配的对象及对象的大小**。

简单的说就是维护一块大的全局内存，每个逻辑处理器 P 维护一块小的私有内存，私有内存不足再从全局申请。

## span 与元素

Go 语言将内存分成了 67 个级别的 span，其中，**0 级代表特殊的大对象，其大小是不固定的**。当具体的对象需要分配内存时，并不是直接分配 span，而是**分配不同级别的 span 中的元素**。因此 span 的级别不是以每个 span 的大小为依据的，而是以 span 中元素的大小为依据的。

根据对象大小，划分了一系列级别 class，每个 class 都代表一个固定大小的对象，以及每个 span 的大小。如下表所示：

`src/runtime/sizeclasses.go`

```go
// class  bytes/obj  bytes/span  objects  tail waste  max waste  min align
//     1          8        8192     1024           0     87.50%          8
//     2         16        8192      512           0     43.75%         16
//     3         24        8192      341           8     29.24%          8
//     4         32        8192      256           0     21.88%         32
//     5         48        8192      170          32     31.52%         16
//     6         64        8192      128           0     23.44%         64
//     7         80        8192      102          32     19.07%         16
//     8         96        8192       85          32     15.95%         32
//     9        112        8192       73          16     13.56%         16
//    10        128        8192       64           0     11.72%        128
//    11        144        8192       56         128     11.82%         16
//    12        160        8192       51          32      9.73%         32
//    13        176        8192       46          96      9.59%         16
//    14        192        8192       42         128      9.25%         64
//    15        208        8192       39          80      8.12%         16
//    16        224        8192       36         128      8.15%         32
//    17        240        8192       34          32      6.62%         16
//    18        256        8192       32           0      5.86%        256
//    19        288        8192       28         128     12.16%         32
//    20        320        8192       25         192     11.80%         64
//    21        352        8192       23          96      9.88%         32
//    22        384        8192       21         128      9.51%        128
//    23        416        8192       19         288     10.71%         32
//    24        448        8192       18         128      8.37%         64
//    25        480        8192       17          32      6.82%         32
//    26        512        8192       16           0      6.05%        512
//    27        576        8192       14         128     12.33%         64
//    28        640        8192       12         512     15.48%        128
//    29        704        8192       11         448     13.93%         64
//    30        768        8192       10         512     13.94%        256
//    31        896        8192        9         128     15.52%        128
//    32       1024        8192        8           0     12.40%       1024
//    33       1152        8192        7         128     12.41%        128
//    34       1280        8192        6         512     15.55%        256
//    35       1408       16384       11         896     14.00%        128
//    36       1536        8192        5         512     14.00%        512
//    37       1792       16384        9         256     15.57%        256
//    38       2048        8192        4           0     12.45%       2048
//    39       2304       16384        7         256     12.46%        256
//    40       2688        8192        3         128     15.59%        128
//    41       3072       24576        8           0     12.47%       1024
//    42       3200       16384        5         384      6.22%        128
//    43       3456       24576        7         384      8.83%        128
//    44       4096        8192        2           0     15.60%       4096
//    45       4864       24576        5         256     16.65%        256
//    46       5376       16384        3         256     10.92%        256
//    47       6144       24576        4           0     12.48%       2048
//    48       6528       32768        5         128      6.23%        128
//    49       6784       40960        6         256      4.36%        128
//    50       6912       49152        7         768      3.37%        256
//    51       8192        8192        1           0     15.61%       8192
//    52       9472       57344        6         512     14.28%        256
//    53       9728       49152        5         512      3.64%        512
//    54      10240       40960        4           0      4.99%       2048
//    55      10880       32768        3         128      6.24%        128
//    56      12288       24576        2           0     11.45%       4096
//    57      13568       40960        3         256      9.99%        256
//    58      14336       57344        4           0      5.35%       2048
//    59      16384       16384        1           0     12.49%       8192
//    60      18432       73728        4           0     11.11%       2048
//    61      19072       57344        3         128      3.57%        128
//    62      20480       40960        2           0      6.87%       4096
//    63      21760       65536        3         256      6.25%        256
//    64      24576       24576        1           0     11.45%       8192
//    65      27264       81920        3         128     10.00%        128
//    66      28672       57344        2           0      4.91%       4096
//    67      32768       32768        1           0     12.50%       8192

// alignment  bits  min obj size
//         8     3             8
//        16     4            32
//        32     5           256
//        64     6           512
//       128     7           768
//      4096    12         28672
//      8192    13         32768
```

上表中每列含义如下：

- class：class ID，每个 span 结构中都有一个 class ID, 表示该 span 可处理的对象类型
- bytes/obj：该 class 代表对象的字节数
- bytes/span：每个 span 占用堆的字节数，也即`页数*页大小`
- objects : 每个 span 可分配的对象个数，也即 `(bytes/spans)/(bytes/obj)`
- waste bytes : 每个 span 产生的内存碎片，也即 `(bytes/spans)%(bytes/obj)`

上表可见，最小的对象是 8B，**最大的对象是 32KB**，超过 32KB 大小的由特殊的 class 表示，该 class ID 为 0，每个 class 只包含一个对象。

**每个具体的对象在分配时都需要对齐到指定的大小**，例如分配 17 字节的对象，会对应分配到比 17 字节大并最接近它的元素级别，即第 3 级，这导致最终分配了 32 字节。因此，这种分配方式会不可避免地带来内存的浪费。

`src/runtime/mheap.go`

```go
//go:notinheap
type mspan struct {
	next *mspan     // next span in list, or nil if none
	prev *mspan     // previous span in list, or nil if none
	list *mSpanList // For debugging. TODO: Remove.

	startAddr uintptr // address of first byte of span aka s.base()
	npages    uintptr // number of pages in span

	manualFreeList gclinkptr // list of free objects in mSpanManual spans

	// freeindex is the slot index between 0 and nelems at which to begin scanning
	// for the next free object in this span.
	// Each allocation scans allocBits starting at freeindex until it encounters a 0
	// indicating a free object. freeindex is then adjusted so that subsequent scans begin
	// just past the newly discovered free object.
	//
	// If freeindex == nelem, this span has no free objects.
	//
	// allocBits is a bitmap of objects in this span.
	// If n >= freeindex and allocBits[n/8] & (1<<(n%8)) is 0
	// then object n is free;
	// otherwise, object n is allocated. Bits starting at nelem are
	// undefined and should never be referenced.
	//
	// Object n starts at address n*elemsize + (start << pageShift).
	freeindex uintptr
	// TODO: Look up nelems from sizeclass and remove this field if it
	// helps performance.
	nelems uintptr // number of object in the span.

	// Cache of the allocBits at freeindex. allocCache is shifted
	// such that the lowest bit corresponds to the bit freeindex.
	// allocCache holds the complement of allocBits, thus allowing
	// ctz (count trailing zero) to use it directly.
	// allocCache may contain bits beyond s.nelems; the caller must ignore
	// these.
	allocCache uint64

	// allocBits and gcmarkBits hold pointers to a span's mark and
	// allocation bits. The pointers are 8 byte aligned.
	// There are three arenas where this data is held.
	// free: Dirty arenas that are no longer accessed
	//       and can be reused.
	// next: Holds information to be used in the next GC cycle.
	// current: Information being used during this GC cycle.
	// previous: Information being used during the last GC cycle.
	// A new GC cycle starts with the call to finishsweep_m.
	// finishsweep_m moves the previous arena to the free arena,
	// the current arena to the previous arena, and
	// the next arena to the current arena.
	// The next arena is populated as the spans request
	// memory to hold gcmarkBits for the next GC cycle as well
	// as allocBits for newly allocated spans.
	//
	// The pointer arithmetic is done "by hand" instead of using
	// arrays to avoid bounds checks along critical performance
	// paths.
	// The sweep will free the old allocBits and set allocBits to the
	// gcmarkBits. The gcmarkBits are replaced with a fresh zeroed
	// out memory.
	allocBits  *gcBits
	gcmarkBits *gcBits

	// sweep generation:
	// if sweepgen == h->sweepgen - 2, the span needs sweeping
	// if sweepgen == h->sweepgen - 1, the span is currently being swept
	// if sweepgen == h->sweepgen, the span is swept and ready to use
	// if sweepgen == h->sweepgen + 1, the span was cached before sweep began and is still cached, and needs sweeping
	// if sweepgen == h->sweepgen + 3, the span was swept and then cached and is still cached
	// h->sweepgen is incremented by 2 after every GC

	sweepgen              uint32
	divMul                uint32        // for divide by elemsize
	allocCount            uint16        // number of allocated objects
	spanclass             spanClass     // size class and noscan (uint8)
	state                 mSpanStateBox // mSpanInUse etc; accessed atomically (get/set methods)
	needzero              uint8         // needs to be zeroed before allocation
	allocCountBeforeCache uint16        // a copy of allocCount that is stored just before this span is cached
	elemsize              uintptr       // computed from sizeclass or from npages
	limit                 uintptr       // end of data in span
	speciallock           mutex         // guards specials list
	specials              *special      // linked list of special records sorted by offset.
}
```

## 三级对象管理

[三级对象管理](object.md)

## 四级内存管理

根据对象的大小，Go 语言将堆内存分成了图 18-3 所示的 HeapArena、chunk、span 与 page 4 种内存块进行管理。

其中，**HeapArena** 内存块最大，其大小与平台相关，在 64 位操作系统中占据 64MB。**chunk** 占据了 512KB，**span** 根据级别大小的不同而不同，就是之前讲的 67 级别，但必须是 **page** 的倍数。而 1 个 page 占据 8KB。不同的内存块用于不同的场景，便于高效地对内存进行管理。

![](../../../assets/images/docs/internal/memory/README/图18-3%20内存块管理结构.png)


## 对象分配

[对象分配](object_malloc.md)

```go

```
