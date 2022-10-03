---
date: 2022-10-03T08:19:31+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数组底层原理"  # 文章标题
url:  "posts/go/docs/internal/array/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 编译时数组解析

数组形如 `[n]T`，在编译时就需要确定其长度和类型。数组在编译时构建抽象语法树阶段的数据类型为 TARRAY，通过 NewArray 函数进行创建，AST 节点的 Op 操作为 OARRAYLIT。

`src/cmd/compile/internal/types/type.go`

```go
// NewArray returns a new fixed-length array Type.
func NewArray(elem *Type, bound int64) *Type {
	if bound < 0 {
		base.Fatalf("NewArray: invalid bound %v", bound)
	}
	t := newType(TARRAY)
	t.extra = &Array{Elem: elem, Bound: bound}
	if elem.HasTParam() {
		t.SetHasTParam(true)
	}
	if elem.HasShape() {
		t.SetHasShape(true)
	}
	return t
}
```

TARRAY 内部的 Array 结构存储了数组中元素的类型及数组的大小。

```go
// Array contains Type fields specific to array types.
type Array struct {
	Elem  *Type // element type
	Bound int64 // number of elements; <0 if unknown yet
}
```

## 数组字面量初始化原理

数组字面量的初始化在编译时类型检查阶段进行，通过 typecheckarraylit 函数循环字面量并分别进行赋值。

`src/cmd/compile/internal/typecheck/typecheck.go`

```go
// typecheckarraylit type-checks a sequence of slice/array literal elements.
func typecheckarraylit(elemType *types.Type, bound int64, elts []ir.Node, ctx string) int64 {
	// If there are key/value pairs, create a map to keep seen
	// keys so we can check for duplicate indices.
	var indices map[int64]bool
	for _, elt := range elts {
		if elt.Op() == ir.OKEY {
			indices = make(map[int64]bool)
			break
		}
	}

	var key, length int64
	for i, elt := range elts {
		ir.SetPos(elt)
		r := elts[i]
		var kv *ir.KeyExpr
		if elt.Op() == ir.OKEY {
			elt := elt.(*ir.KeyExpr)
			elt.Key = Expr(elt.Key)
			key = IndexConst(elt.Key)
			if key < 0 {
				base.Fatalf("invalid index: %v", elt.Key)
			}
			kv = elt
			r = elt.Value
		}

		r = Expr(r)
		r = AssignConv(r, elemType, ctx)
		if kv != nil {
			kv.Value = r
		} else {
			elts[i] = r
		}

		if key >= 0 {
			if indices != nil {
				if indices[key] {
					base.Errorf("duplicate index in %s: %d", ctx, key)
				} else {
					indices[key] = true
				}
			}

			if bound >= 0 && key >= bound {
				base.Errorf("array index %d out of bounds [0:%d]", key, bound)
				bound = -1
			}
		}

		key++
		if key > length {
			length = key
		}
	}

	return length
}
```

## 数组字面量编译时内存优化

在编译时还会进行重要的优化。在函数 walk 遍历阶段，anylit 函数用于处理各种类型的字面量。当数组的长度小于 4 时，在运行时数组会被放置在栈中，当数组的长度大于 4 时，数组会被放置到内存的静态只读区。fixedlit 函数将执行数组初始化与赋值的逻辑。

```go
func anylit(n ir.Node, var_ ir.Node, init *ir.Nodes) {
	t := n.Type()
	switch n.Op() {
	default:
		base.Fatalf("anylit: not lit, op=%v node=%v", n.Op(), n)
    ...
	case ir.OSTRUCTLIT, ir.OARRAYLIT:
		n := n.(*ir.CompLitExpr)
		if !t.IsStruct() && !t.IsArray() {
			base.Fatalf("anylit: not struct/array")
		}

		if isSimpleName(var_) && len(n.List) > 4 {
			// lay out static data
			vstat := readonlystaticname(t)

			ctxt := inInitFunction
			if n.Op() == ir.OARRAYLIT {
				ctxt = inNonInitFunction
			}
			fixedlit(ctxt, initKindStatic, n, vstat, init)

			// copy static to var
			appendWalkStmt(init, ir.NewAssignStmt(base.Pos, var_, vstat))

			// add expressions to automatic
			fixedlit(inInitFunction, initKindDynamic, n, var_, init)
			break
		}

		var components int64
		if n.Op() == ir.OARRAYLIT {
			components = t.NumElem()
		} else {
			components = int64(t.NumFields())
		}
		// initialization of an array or struct with unspecified components (missing fields or arrays)
		if isSimpleName(var_) || int64(len(n.List)) < components {
			appendWalkStmt(init, ir.NewAssignStmt(base.Pos, var_, nil))
		}

		fixedlit(inInitFunction, initKindLocalCode, n, var_, init)
    ...
	}
}
```

## 数组索引与访问越界原理

数组访问越界是非常严重的错误，Go 语言中对越界的判断有一部分是在编译期间类型检查阶段完成的，tcIndex 函数会对访问数组的索引进行验证。

```go
// tcIndex typechecks an OINDEX node.
func tcIndex(n *ir.IndexExpr) ir.Node {
	n.X = Expr(n.X)
	n.X = DefaultLit(n.X, nil)
	n.X = implicitstar(n.X)
	l := n.X
	n.Index = Expr(n.Index)
	r := n.Index
	t := l.Type()
	if t == nil || r.Type() == nil {
		n.SetType(nil)
		return n
	}
	switch t.Kind() {
	default:
		base.Errorf("invalid operation: %v (type %v does not support indexing)", n, t)
		n.SetType(nil)
		return n

	case types.TSTRING, types.TARRAY, types.TSLICE:
		n.Index = indexlit(n.Index)
		if t.IsString() {
			n.SetType(types.ByteType)
		} else {
			n.SetType(t.Elem())
		}
		why := "string"
		if t.IsArray() {
			why = "array"
		} else if t.IsSlice() {
			why = "slice"
		}

		if n.Index.Type() != nil && !n.Index.Type().IsInteger() {
			base.Errorf("non-integer %s index %v", why, n.Index)
			return n
		}

		if !n.Bounded() && ir.IsConst(n.Index, constant.Int) {
			x := n.Index.Val()
			if constant.Sign(x) < 0 {
				base.Errorf("invalid %s index %v (index must be non-negative)", why, n.Index)
			} else if t.IsArray() && constant.Compare(x, token.GEQ, constant.MakeInt64(t.NumElem())) {
				base.Errorf("invalid array index %v (out of bounds for %d-element array)", n.Index, t.NumElem())
			} else if ir.IsConst(n.X, constant.String) && constant.Compare(x, token.GEQ, constant.MakeInt64(int64(len(ir.StringVal(n.X))))) {
				base.Errorf("invalid string index %v (out of bounds for %d-byte string)", n.Index, len(ir.StringVal(n.X)))
			} else if ir.ConstOverflow(x, types.Types[types.TINT]) {
				base.Errorf("invalid %s index %v (index too large)", why, n.Index)
			}
		}

	case types.TMAP:
		n.Index = AssignConv(n.Index, t.Key(), "map index")
		n.SetType(t.Elem())
		n.SetOp(ir.OINDEXMAP)
		n.Assigned = false
	}
	return n
}
```

使用未命名常量索引访问数组时，例如 `a[3]`，数组的一些简单越界错误能够在编译期间被发现。但是如果使用变量去访问数组或者字符串，编译器无法发现对应的错误，因为变量的值随时可能变化。

Go 语言运行时在发现数组、切片和字符串的越界操作时，会由运行时的 panicIndex 和 runtime.goPanicIndex 函数触发程序的运行时错误并导致崩溃。

```go
// Note: these declarations are just for wasm port.
// Other ports call assembly stubs instead.
func goPanicIndex(x int, y int)

func goPanicIndex(x int, y int) {
	panicCheck1(getcallerpc(), "index out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsIndex})
}
```

如果数组的索引是命名常量，那么仍然能够在编译时通过优化检测出越界并在运行时报错。例如对于一个简单的代码：

```go
package main

func main() {
	a := [3]int{1, 2, 3}
	b := 8
	_ = a[b]
}
```

可以通过如下 `go tool compile` 命令生成 ssa.html 文件，显示整个 SSA 生成与优化阶段。

```
GOSSAFUNC=main GOOS=linux GOARCH=amd64 go tool compile main.go
```

```go

```
