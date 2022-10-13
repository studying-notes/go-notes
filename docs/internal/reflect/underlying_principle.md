---

date: 2022-10-08T13:32:04+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "反射底层原理"  # 文章标题
url:  "posts/go/docs/internal/reflect/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## reflect.Type 详解

通过如下 reflect.TypeOf 函数对于 reflect.Type 的构建过程可以发现，其实现原理为将传递进来的接口变量转换为底层的实际空接口 emptyInterface，并获取空接口的类型值。reflect.Type 实质上是空接口结构体中的 typ 字段，其是 rtype 类型，Go 语言中任何具体类型的底层结构都包含这一类型。

```go
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
```

生成 reflect.Value 的原理也可以从 reflect.ValueOf 函数的生成方法中看出端倪。reflect.ValueOf 函数的核心是调用了 unpackEface 函数。

```go
// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func ValueOf(i any) Value {
	if i == nil {
		return Value{}
	}

	// TODO: Maybe allow contents of a Value to live on the stack.
	// For now we make the contents always escape to the heap. It
	// makes life easier in a few places (see chanrecv/mapassign
	// comment below).
	escapes(i)

	return unpackEface(i)
}

// unpackEface converts the empty interface i to a Value.
func unpackEface(i any) Value {
	e := (*emptyInterface)(unsafe.Pointer(&i))
	// NOTE: don't read e.word until we know whether it is really a pointer or not.
	t := e.typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}
}
```

reflect.Value 包含了接口中存储的值及类型，除此之外还包含了特殊的 flag 标志。

如图 13-1 所示，flag 标记以位图的形式存储了反射类型的元数据。

![](../../../assets/images/docs/internal/reflect/underlying_principle/图13-1%20反射flag位图.png)

其中，flag 的低 5 位存储了类型的标志，利用 flag.kind 方法有助于快速知道反射中存储的类型。

`reflect/value.go`

```go
func (f flag) kind() Kind {
	return Kind(f & flagKindMask)
}
```

低 6～10 位代表了字段的一些特征，例如该字段是否是可以外部访问的、是否可以寻址、是否是方法等。具体含义如下：

```go
const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5 // 结构体未导出字段，不是嵌入字段
	flagEmbedRO     flag = 1 << 6 // 结构体未导出字段，是嵌入字段
	flagIndir       flag = 1 << 7 // 间接的，val 存储了可以寻址的指针
	flagAddr        flag = 1 << 8 // 可寻址的
	flagMethod      flag = 1 << 9 // 方法
	flagMethodShift      = 10   // 方法的偏移量
	flagRO          flag = flagStickyRO | flagEmbedRO // 结构体未导出字段
)
```

flag 的其余位存储了方法的 index 序号，代表第几个方法。只有在当前的 value 是方法类型时才会用到。例如第 5 号方法，其存储的位置为 5<<10。

其中，flagIndir 是最让人困惑的标志，代表间接的。我们知道存储在反射或接口中的值都是指针，如下面这个简单的例子，虽然看似存储的是 int 值 3，但在反射中实际上存储的是指针。

```go
k := 3
v := reflect.ValueOf(k)
```

因此，为了和如下实际存储的是指针的场景进行区别，需要使用 flagIndir 来标识当前存储的值是间接的，是需要根据当前的指针进行寻址的。

```go
k := 3
v := reflect.ValueOf(&k)
```

另外，容器类型如切片、哈希表、通道也被认为是间接的，因为它们也需要当前容器的指针间接找到存储在其内部的元素。

## Interface 方法原理

知道了空接口的构成以及 reflect.Value 存储了空接口的类型和值，就可以理解 Interface 方法将 reflect.Value 转换为空接口是非常容易的。Interface 核心方法调用了 packEface 函数。

```go
// packEface converts v to the empty interface.
func packEface(v Value) any {
	t := v.typ
	var i any
	e := (*emptyInterface)(unsafe.Pointer(&i))
	// First, fill in the data portion of the interface.
	switch {
	case ifaceIndir(t):
		if v.flag&flagIndir == 0 {
			panic("bad indir")
		}
		// Value is indirect, and so is the interface we're making.
		ptr := v.ptr
		if v.flag&flagAddr != 0 {
			// TODO: pass safe boolean from valueInterface so
			// we don't need to copy if safe==true?
			c := unsafe_New(t)
			typedmemmove(t, c, ptr)
			ptr = c
		}
		e.word = ptr
	case v.flag&flagIndir != 0:
		// Value is indirect, but interface is direct. We need
		// to load the data at v.ptr into the interface data word.
		e.word = *(*unsafe.Pointer)(v.ptr)
	default:
		// Value is direct, and so is the interface.
		e.word = v.ptr
	}
	// Now, fill in the type portion. We're very careful here not
	// to have any operation between the e.word and e.typ assignments
	// that would let the garbage collector observe the partially-built
	// interface value.
	e.typ = t
	return i
}

// unpackEface converts the empty interface i to a Value.
func unpackEface(i any) Value {
	e := (*emptyInterface)(unsafe.Pointer(&i))
	// NOTE: don't read e.word until we know whether it is really a pointer or not.
	t := e.typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}
}
```

`e: =(*emptyInterface)(unsafe.Pointer(&i))` 构建了一个空接口，e.typ = t 将 reflect.Value 中的类型赋值给空接口中的类型。但是对于接口中的值 e.word 的处理仍然有所区别，原因和之前介绍的一样，有些值是间接获得的。

假如有一个 int 切片，通过 vvv.Index(1) 可以得到切片中序号为 1 的值，即 23。

```go
vvv := reflect.ValueOf([]int{1, 23, 4})
vvv.Index(1).Interface().(int) // 23
```

但实际上当前 reflect.Value 中存储的是当前数字在切片中的地址。如图 13-2 所示，我们构造的 interface.word 应该是当前的 value.ptr 地址吗？显然不是，因为我们希望返回的类型是一个新的副本，这样不会对原始的切片造成任何干扰。当出现这种情况时，case ifaceIndir(t) 为 true，会生成一个新的值。

![](../../../assets/images/docs/internal/reflect/underlying_principle/图13-2%20切片反射示例1.png)

而如果我们面对的是如图13-3所示的情形：

```go
aa, bb, cc := 13, 23, 33
a := []*int{&aa, &bb, &cc}
b := reflect.ValueOf(a)
b.Index(1).Interface()
```

![](../../../assets/images/docs/internal/reflect/underlying_principle/图13-3%20切片反射示例2.png)

那么构造的 interface.word 应该是当前的 value.ptr 地址吗？显然也不是，而应该是存储在内部的实际指向数据的指针。

```go
e.word = *(*unsafe.Pointer)(v.ptr)
```

## Elem 方法

```go
package main

import "reflect"

func main() {
	kk := 3
	kkk := reflect.ValueOf(kk)
	kkk.SetInt(4)
}
```

```go
panic: reflect: reflect.Value.SetInt using unaddressable value
```

Value.SetInt 函数会在一开始检查 Value 的 flag 中是否有 flagAddr 字段。

而在通过 ValueOf 构建 reflect.Value 时，只判断了变量是否为间接的地址，因此会报错。

```go
// mustBeAssignable panics if f records that the value is not assignable,
// which is to say that either it was obtained using an unexported field
// or it is not addressable.
func (f flag) mustBeAssignable() {
	if f&flagRO != 0 || f&flagAddr == 0 {
		f.mustBeAssignableSlow()
	}
}
```

```go
// SetInt sets v's underlying value to x.
// It panics if v's Kind is not Int, Int8, Int16, Int32, or Int64, or if CanSet() is false.
func (v Value) SetInt(x int64) {
	v.mustBeAssignable()
	switch k := v.kind(); k {
	default:
		panic(&ValueError{"reflect.Value.SetInt", v.kind()})
	case Int:
		*(*int)(v.ptr) = int(x)
	case Int8:
		*(*int8)(v.ptr) = int8(x)
	case Int16:
		*(*int16)(v.ptr) = int16(x)
	case Int32:
		*(*int32)(v.ptr) = int32(x)
	case Int64:
		*(*int64)(v.ptr) = x
	}
}
```

Elem 的功能是返回接口内部包含的或指针指向的数据值。

```go
// Elem returns the value that the interface v contains
// or that the pointer v points to.
// It panics if v's Kind is not Interface or Pointer.
// It returns the zero Value if v is nil.
func (v Value) Elem() Value {
	k := v.kind()
	switch k {
	case Interface:
		var eface any
		if v.typ.NumMethod() == 0 {
			eface = *(*any)(v.ptr)
		} else {
			eface = (any)(*(*interface {
				M()
			})(v.ptr))
		}
		x := unpackEface(eface)
		if x.flag != 0 {
			x.flag |= v.flag.ro()
		}
		return x
	case Pointer:
		ptr := v.ptr
		if v.flag&flagIndir != 0 {
			if ifaceIndir(v.typ) {
				// This is a pointer to a not-in-heap object. ptr points to a uintptr
				// in the heap. That uintptr is the address of a not-in-heap object.
				// In general, pointers to not-in-heap objects can be total junk.
				// But Elem() is asking to dereference it, so the user has asserted
				// that at least it is a valid pointer (not just an integer stored in
				// a pointer slot). So let's check, to make sure that it isn't a pointer
				// that the runtime will crash on if it sees it during GC or write barriers.
				// Since it is a not-in-heap pointer, all pointers to the heap are
				// forbidden! That makes the test pretty easy.
				// See issue 48399.
				if !verifyNotInHeapPtr(*(*uintptr)(ptr)) {
					panic("reflect: reflect.Value.Elem on an invalid notinheap pointer")
				}
			}
			ptr = *(*unsafe.Pointer)(ptr)
		}
		// The returned value's address is v's value.
		if ptr == nil {
			return Value{}
		}
		tt := (*ptrType)(unsafe.Pointer(v.typ))
		typ := tt.elem
		fl := v.flag&flagRO | flagIndir | flagAddr
		fl |= flag(typ.Kind())
		return Value{typ, ptr, fl}
	}
	panic(&ValueError{"reflect.Value.Elem", v.kind()})
}
```

对于指针来说，如果 flag 标识了 reflect.Value 是间接的，则会返回数据真实的地址 `(*unsafe.Pointer)(ptr)`，而对于直接的指针，则返回本身即可，并且会将 flag 修改为 flagAddr，即可赋值的。

## 动态调用剖析

反射提供的核心动能是动态的调用方法或函数，这在 RPC 远程过程调用中使用频繁。MethodByName 方法可以根据方法名找到代表方法的 reflect.Value 对象。

```go
// MethodByName returns a function value corresponding to the method
// of v with the given name.
// The arguments to a Call on the returned function should not include
// a receiver; the returned function will always use v as the receiver.
// It returns the zero Value if no method was found.
func (v Value) MethodByName(name string) Value {
	if v.typ == nil {
		panic(&ValueError{"reflect.Value.MethodByName", Invalid})
	}
	if v.flag&flagMethod != 0 {
		return Value{}
	}
	m, ok := v.typ.MethodByName(name)
	if !ok {
		return Value{}
	}
	return v.Method(m.Index)
}

func (t *rtype) MethodByName(name string) (m Method, ok bool) {
	if t.Kind() == Interface {
		tt := (*interfaceType)(unsafe.Pointer(t))
		return tt.MethodByName(name)
	}
	ut := t.uncommon()
	if ut == nil {
		return Method{}, false
	}
	// TODO(mdempsky): Binary search.
	for i, p := range ut.exportedMethods() {
		if t.nameOff(p.name).name() == name {
			return t.Method(i), true
		}
	}
	return Method{}, false
}

// MethodByName method with the given name in the type's method set.
func (t *interfaceType) MethodByName(name string) (m Method, ok bool) {
	if t == nil {
		return
	}
	var p *imethod
	for i := range t.methods {
		p = &t.methods[i]
		if t.nameOff(p.name).name() == name {
			return t.Method(i), true
		}
	}
	return
}
```

动态调用的核心方法是 Call 方法，其参数为 reflect.Value 数组，返回的也是 reflect.Value 数组，由于代码较长，下面将对代码流程逐一进行分析。

```go
// Call calls the function v with the input arguments in.
// For example, if len(in) == 3, v.Call(in) represents the Go call v(in[0], in[1], in[2]).
// Call panics if v's Kind is not Func.
// It returns the output results as Values.
// As in Go, each input argument must be assignable to the
// type of the function's corresponding input parameter.
// If v is a variadic function, Call creates the variadic slice parameter
// itself, copying in the corresponding values.
func (v Value) Call(in []Value) []Value {
	v.mustBe(Func)
	v.mustBeExported()
	return v.call("Call", in)
}
```

Call 方法的第 1 步是获取函数的指针，对于方法的调用要略微复杂一些，会调用 methodReceiver 方法获取调用者的实际类型、函数类型，以及函数指针的位置。

```go
	if v.flag&flagMethod != 0 {
		rcvr = v
缓存		rcvrtype, t, fn = methodReceiver(op, v, int(v.flag)>>flagMethodShift)
	} else if v.flag&flagIndir != 0 {
		fn = *(*unsafe.Pointer)(v.ptr)
	} else {
		fn = v.ptr
	}
```

第 2 步是进行有效性验证，例如函数的输入大小和个数是否与传入的参数匹配，传入的参数能否赋值给函数参数等。

第 3 步是调用 funcLayout 函数，用于构建函数参数及返回值的栈帧布局，其中 frametype 代表调用时需要的内存大小，用于内存分配。retOffset 用于标识函数参数及返回值在内存中的位置。

```go
	// Compute frame type.
	frametype, framePool, abid := funcLayout(t, rcvrtype)
```

framePool 是一个内存缓存池，用于在没有返回值的场景中复用内存。但是如果函数中有返回值，则不能复用内存，这是为了防止发生内存泄漏。

```go
	// Allocate a chunk of memory for frame if needed.
	var stackArgs unsafe.Pointer
	if frametype.size != 0 {
		if nout == 0 {
			stackArgs = framePool.Get().(unsafe.Pointer)
		} else {
			// Can't use pool if the function has return values.
			// We will leak pointer to args in ret, so its lifetime is not scoped.
			stackArgs = unsafe_New(frametype)
		}
	}
```

如果是方法调用，那么栈中的第一个参数是接收者的指针。

```go
	// Copy inputs into args.

	// Handle receiver.
	inStart := 0
	if rcvrtype != nil {
		// Guaranteed to only be one word in size,
		// so it will only take up exactly 1 abiStep (either
		// in a register or on the stack).
		switch st := abid.call.steps[0]; st.kind {
		case abiStepStack:
			storeRcvr(rcvr, stackArgs)
		case abiStepPointer:
			storeRcvr(rcvr, unsafe.Pointer(&regArgs.Ptrs[st.ireg]))
			fallthrough
		case abiStepIntReg:
			storeRcvr(rcvr, unsafe.Pointer(&regArgs.Ints[st.ireg]))
		case abiStepFloatReg:
			storeRcvr(rcvr, unsafe.Pointer(&regArgs.Floats[st.freg]))
		default:
			panic("unknown ABI parameter kind")
		}
		inStart = 1
	}
```

然后将输入参数放入栈中，

```go
	// Handle arguments.
	for i, v := range in {
		v.mustBeExported()
		targ := t.In(i).(*rtype)
		// TODO(mknyszek): Figure out if it's possible to get some
		// scratch space for this assignment check. Previously, it
		// was possible to use space in the argument frame.
		v = v.assignTo("reflect.Value.Call", targ, nil)
	stepsLoop:
		for _, st := range abid.call.stepsForValue(i + inStart) {
			switch st.kind {
			case abiStepStack:
				// Copy values to the "stack."
				addr := add(stackArgs, st.stkOff, "precomputed stack arg offset")
				if v.flag&flagIndir != 0 {
					typedmemmove(targ, addr, v.ptr)
				} else {
					*(*unsafe.Pointer)(addr) = v.ptr
				}
				// There's only one step for a stack-allocated value.
				break stepsLoop
			case abiStepIntReg, abiStepPointer:
				// Copy values to "integer registers."
				if v.flag&flagIndir != 0 {
					offset := add(v.ptr, st.offset, "precomputed value offset")
					if st.kind == abiStepPointer {
						// Duplicate this pointer in the pointer area of the
						// register space. Otherwise, there's the potential for
						// this to be the last reference to v.ptr.
						regArgs.Ptrs[st.ireg] = *(*unsafe.Pointer)(offset)
					}
					intToReg(&regArgs, st.ireg, st.size, offset)
				} else {
					if st.kind == abiStepPointer {
						// See the comment in abiStepPointer case above.
						regArgs.Ptrs[st.ireg] = v.ptr
					}
					regArgs.Ints[st.ireg] = uintptr(v.ptr)
				}
			case abiStepFloatReg:
				// Copy values to "float registers."
				if v.flag&flagIndir == 0 {
					panic("attempted to copy pointer to FP register")
				}
				offset := add(v.ptr, st.offset, "precomputed value offset")
				floatToReg(&regArgs, st.freg, st.size, offset)
			default:
				panic("unknown ABI part kind")
			}
		}
	}

	// TODO(mknyszek): Remove this when we no longer have
	// caller reserved spill space.
	frameSize = align(frameSize, goarch.PtrSize)
	frameSize += abid.spill
```

```go
// align returns the result of rounding x up to a multiple of n.
// n must be a power of two.
func align(x, n uintptr) uintptr {
	return (x + n - 1) &^ (n - 1)
}
```

`off=(off+a-1)&^(a-1)` 是计算内存对齐的标准方式，在结构体内存对齐中使用频繁。

调用 Call 汇编函数完成调用逻辑，Call 函数需要传递内存布局类型(frametype)、函数指针(fn)、内存地址(args)、栈大小(frametype.size)、输入参数与返回值的内存间隔(retOffset)。

完成调用后，如果函数没有返回，则将 args 内部全部清空为 0，并再次放入 framePool 中。如果有返回值，则清空 args 中输入参数部分，并将输出包装为 ret 切片后返回。

```go
	var ret []Value
	if nout == 0 {
		if stackArgs != nil {
			typedmemclr(frametype, stackArgs)
			framePool.Put(stackArgs)
		}
	} else {
		if stackArgs != nil {
			// Zero the now unused input area of args,
			// because the Values returned by this function contain pointers to the args object,
			// and will thus keep the args object alive indefinitely.
			typedmemclrpartial(frametype, stackArgs, 0, abid.retOffset)
		}

		// Wrap Values around return values in args.
		ret = make([]Value, nout)
		for i := 0; i < nout; i++ {
			tv := t.Out(i)
			if tv.Size() == 0 {
				// For zero-sized return value, args+off may point to the next object.
				// In this case, return the zero value instead.
				ret[i] = Zero(tv)
				continue
			}
			steps := abid.ret.stepsForValue(i)
			if st := steps[0]; st.kind == abiStepStack {
				// This value is on the stack. If part of a value is stack
				// allocated, the entire value is according to the ABI. So
				// just make an indirection into the allocated frame.
				fl := flagIndir | flag(tv.Kind())
				ret[i] = Value{tv.common(), add(stackArgs, st.stkOff, "tv.Size() != 0"), fl}
				// Note: this does introduce false sharing between results -
				// if any result is live, they are all live.
				// (And the space for the args is live as well, but as we've
				// cleared that space it isn't as big a deal.)
				continue
			}

			// Handle pointers passed in registers.
			if !ifaceIndir(tv.common()) {
				// Pointer-valued data gets put directly
				// into v.ptr.
				if steps[0].kind != abiStepPointer {
					print("kind=", steps[0].kind, ", type=", tv.String(), "\n")
					panic("mismatch between ABI description and types")
				}
				ret[i] = Value{tv.common(), regArgs.Ptrs[steps[0].ireg], flag(tv.Kind())}
				continue
			}

			// All that's left is values passed in registers that we need to
			// create space for and copy values back into.
			//
			// TODO(mknyszek): We make a new allocation for each register-allocated
			// value, but previously we could always point into the heap-allocated
			// stack frame. This is a regression that could be fixed by adding
			// additional space to the allocated stack frame and storing the
			// register-allocated return values into the allocated stack frame and
			// referring there in the resulting Value.
			s := unsafe_New(tv.common())
			for _, st := range steps {
				switch st.kind {
				case abiStepIntReg:
					offset := add(s, st.offset, "precomputed value offset")
					intFromReg(&regArgs, st.ireg, st.size, offset)
				case abiStepPointer:
					s := add(s, st.offset, "precomputed value offset")
					*((*unsafe.Pointer)(s)) = regArgs.Ptrs[st.ireg]
				case abiStepFloatReg:
					offset := add(s, st.offset, "precomputed value offset")
					floatFromReg(&regArgs, st.freg, st.size, offset)
				case abiStepStack:
					panic("register-based return value has stack component")
				default:
					panic("unknown ABI part kind")
				}
			}
			ret[i] = Value{tv.common(), s, flagIndir | flag(tv.Kind())}
		}
	}

	return ret
```

```go

```
