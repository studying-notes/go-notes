---
date: 2022-10-03T09:39:52+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "切片底层原理"  # 文章标题
url:  "posts/go/docs/internal/slice/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

在编译时构建抽象语法树阶段会将切片构建为如下类型。

`src/cmd/compile/internal/types/type.go`

```go
// Slice contains Type fields specific to slice types.
type Slice struct {
	Elem *Type // element type
}
```

编译时使用 NewSlice 函数新建一个切片类型，并需要传递切片元素的类型。从中可以看出，切片元素的类型 elem 是在编译期间确定的。

```go
// NewSlice returns the slice Type with element type elem.
func NewSlice(elem *Type) *Type {
	if t := elem.cache.slice; t != nil {
		if t.Elem() != elem {
			base.Fatalf("elem mismatch")
		}
		if elem.HasTParam() != t.HasTParam() || elem.HasShape() != t.HasShape() {
			base.Fatalf("Incorrect HasTParam/HasShape flag for cached slice type")
		}
		return t
	}

	t := newType(TSLICE)
	t.extra = Slice{Elem: elem}
	elem.cache.slice = t
	if elem.HasTParam() {
		t.SetHasTParam(true)
	}
	if elem.HasShape() {
		t.SetHasShape(true)
	}
	return t
}
```

## 字面量初始化

当使用形如 `[]int{1, 2, 3}` 的字面量创建新的切片时，会创建一个 array 数组(`[3]int{1, 2, 3}`)存储于静态区中，并在堆区创建一个新的切片，在程序启动时将静态区的数据复制到堆区，这样可以加快切片的初始化过程。

`src/cmd/compile/internal/walk/complit.go`

```go
func slicelit(ctxt initContext, n *ir.CompLitExpr, var_ ir.Node, init *ir.Nodes)

// recipe for var = []t{...}
// 1. make a static array
//	var vstat [...]t
// 2. assign (data statements) the constant part
//	vstat = constpart{}
// 3. make an auto pointer to array and allocate heap to it
//	var vauto *[...]t = new([...]t)
// 4. copy the static array to the auto array
//	*vauto = vstat
// 5. for each dynamic part assign to the array
//	vauto[i] = dynamic part
// 6. assign slice of allocated heap to var
//	var = vauto[:]
//
// an optimization is done if there is no constant part
//	3. var vauto *[...]t = new([...]t)
//	5. vauto[i] = dynamic part
//	6. var = vauto[:]
```

## make 初始化

对形如 make([]int,3,4) 的初始化切片。节点 Node 的 Op 操作为 OMAKESLICE，并且左节点存储长度为 3，右节点存储容量为 4。

编译时对于字面量的重要优化是判断变量应该被分配到栈中还是应该逃逸到堆区。如果 make 函数初始化了一个太大的切片，则该切片会逃逸到堆中。如果分配了一个比较小的切片，则会直接在栈中分配。

`cmd/compile/internal/ir/cfg.go`

默认为 64KB，可以通过指定编译时 smallframes 标识进行更新，因此，`make([]int64, 1023)` 与 `make([]int64, 1024)` 实现的细节是截然不同的。

```go
var (
	// maximum size variable which we will allocate on the stack.
	// This limit is for explicit variable declarations like "var x T" or "x := ...".
	// Note: the flag smallframes can update this value.
	MaxStackVarSize = int64(10 * 1024 * 1024)

	// maximum size of implicit variables that we will allocate on the stack.
	//   p := new(T)          allocating T on the stack
	//   p := &T{}            allocating T on the stack
	//   s := make([]T, n)    allocating [n]T on the stack
	//   s := []byte("...")   allocating [n]byte on the stack
	// Note: the flag smallframes can update this value.
	MaxImplicitStackVarSize = int64(64 * 1024)

	// MaxSmallArraySize is the maximum size of an array which is considered small.
	// Small arrays will be initialized directly with a sequence of constant stores.
	// Large arrays will be initialized by copying from a static temp.
	// 256 bytes was chosen to minimize generated code + statictmp size.
	MaxSmallArraySize = int64(256)
)
```

如果没有逃逸，那么切片运行时最终会被分配在栈中。而如果发生了逃逸，那么运行时调用 makesliceXX 函数会将切片分配在堆中。当切片的长度和容量小于 int 类型的最大值时，会调用 makeslice 函数，反之调用 makeslice64 函数创建切片。

makeslice64 函数最终也调用了 makeslice 函数。makeslice 函数会先判断要申请的内存大小是否超过了理论上系统可以分配的内存大小，并判断其长度是否小于容量。再调用 mallocgc 函数在堆中申请内存，申请的内存大小为`类型大小 × 容量`。

`src/runtime/slice.go`

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return mallocgc(mem, et, true)
}

func makeslice64(et *_type, len64, cap64 int64) unsafe.Pointer {
	len := int(len64)
	if int64(len) != len64 {
		panicmakeslicelen()
	}

	cap := int(cap64)
	if int64(cap) != cap64 {
		panicmakeslicecap()
	}

	return makeslice(et, len, cap)
}
```

## 切片扩容原理

append 函数在运行时调用了 runtime/slice.go 文件下的 growslice 函数：

```go
// growslice allocates new backing store for a slice.
//
// arguments:
//   oldPtr = pointer to the slice's backing array
//   newLen = new length (= oldLen + num)
//   oldCap = original slice's capacity.
//      num = number of elements being added
//       et = element type
//
// return values:
//   newPtr = pointer to the new backing store
//   newLen = same value as the argument
//   newCap = capacity of the new backing store
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
    ...

newcap := oldCap
	doublecap := newcap + newcap
	if newLen > doublecap {
		newcap = newLen
	} else {
		const threshold = 256
		if oldCap < threshold {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < newLen {
				// Transition from growing 2x for small slices
				// to growing 1.25x for large slices. This formula
				// gives a smooth-ish transition between the two.
				newcap += (newcap + 3*threshold) / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = newLen
			}
		}
	}

    ...

	return slice{p, newLen, newcap}
}
```

上面的代码显示了扩容的核心逻辑，Go 语言中切片扩容的策略为：

- 如果新申请容量（cap）大于 2 倍的旧容量（old.cap），则最终容量（newcap）是新申请的容量（cap）。
- 如果旧切片的长度小于 1024，则最终容量是旧容量的 2 倍，即 newcap = doublecap。
- 如果旧切片长度大于或等于 1024，则最终容量从旧容量开始循环增加原来的 1/4，即 `newcap=old.cap,for{newcap+=newcap/4}`，直到最终容量大于或等于新申请的容量为止，即 newcap ≥ cap。
- 如果最终容量计算值溢出，即超过了 int 的最大范围，则最终容量就是新申请容量。

Growslice 函数会根据切片的类型，分配不同大小的内存。为了对齐内存，申请的内存可能大于`实际的类型大小 × 容量大小`。

如果切片需要扩容，那么最后需要到堆区申请内存。要注意的是，扩容后新的切片不一定拥有新的地址。因此在使用 append 函数时，通常会采用 `a = append(a, T)` 的形式。根据 et.ptrdata 是否判断切片类型为指针，执行不同的逻辑。

```go
	var p unsafe.Pointer
	if et.ptrdata == 0 {
		p = mallocgc(capmem, nil, false)
		// The append() that calls growslice is going to overwrite from oldLen to newLen.
		// Only clear the part that will not be overwritten.
		memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
	} else {
		// Note: can't use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.
		p = mallocgc(capmem, et, true)
		if lenmem > 0 && writeBarrier.enabled {
			// Only shade the pointers in oldPtr since we know the destination slice p
			// only contains nil pointers because it has been cleared during alloc.
			bulkBarrierPreWriteSrcOnly(uintptr(p), uintptr(oldPtr), lenmem-et.size+et.ptrdata)
		}
	}
	memmove(p, oldPtr, lenmem)
```

当切片类型为指针，涉及垃圾回收写屏障开启时，对旧切片中指针指向的对象进行标记。

除了在切片的尾部追加，还可以在切片的开头添加元素：

```go
var a = []int{1, 2, 3}
a = append([]int{0}, a...)          // 在开头添加一个元素
a = append([]int{-3, -2, -1}, a...) // 在开头添加一个切片
```

**在开头一般都会导致内存的重新分配**，而且会导致已有的元素全部复制一次。因此，从切片的开头添加元素的性能一般要比从尾部追加元素的性能差很多。

```go

```
