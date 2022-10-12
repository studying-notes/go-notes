---
date: 2022-10-12T11:06:16+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "构建中间表示"  # 文章标题
url:  "posts/go/docs/internal/compiler/ir"  # 设置网页永久链接
tags: [ "Go", "ir" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 源码分布

这一阶段的相关源码分布在：

* `cmd/compile/internal/types` (compiler types)
* `cmd/compile/internal/ir` (compiler AST)
* `cmd/compile/internal/typecheck` (AST transformations)
* `cmd/compile/internal/noder` (create compiler AST)

## 简介

编译器中间端必须用它自己的 AST 定义和曾经用 C 编写时继承的 Go 类型表示。它的所有代码都是根据这些编写的，所以类型检查之后的下一步是将**语法**和 **types2 表示**转换到 ir 和 types。这个过程被称为 "noding"。

当前有两种实现：

1. irgen (aka "-G=3" or "noder2")，从 1.18 开始。
2. Unified IR，正在开发中的实现 (通过 `GOEXPERIMENT=unified` 启用)

```go
	if base.Debug.Unified != 0 {
		unified(noders)
		return
	}

	// Use types2 to type-check and generate IR.
	check2(noders)
```

在 Go 1.18 之前，有第三种实现 ("noder" or "-G=0")，它直接将预先类型检查的语法表示转换为 IR，然后调用包类型检查的类型检查器。这个实现在 Go 1.18 之后被删除，所以现在包类型检查只用于 IR 转换。

## IR 转换

```go
func (g *irgen) generate(noders []*noder) {
	types.LocalPkg.Name = g.self.Name()
	types.LocalPkg.Height = g.self.Height()
	typecheck.TypecheckAllowed = true

	// Prevent size calculations until we set the underlying type
	// for all package-block defined types.
	types.DeferCheckSize()

	// At this point, types2 has already handled name resolution and
	// type checking. We just need to map from its object and type
	// representations to those currently used by the rest of the
	// compiler. This happens in a few passes.

	// 1. Process all import declarations. We use the compiler's own
	// importer for this, rather than types2's gcimporter-derived one,
	// to handle extensions and inline function bodies correctly.
	//
	// Also, we need to do this in a separate pass, because mappings are
	// instantiated on demand. If we interleaved processing import
	// declarations with other declarations, it's likely we'd end up
	// wanting to map an object/type from another source file, but not
	// yet have the import data it relies on.
	declLists := make([][]syntax.Decl, len(noders))
Outer:
	for i, p := range noders {
		g.pragmaFlags(p.file.Pragma, ir.GoBuildPragma)
		for j, decl := range p.file.DeclList {
			switch decl := decl.(type) {
			case *syntax.ImportDecl:
				g.importDecl(p, decl)
			default:
				declLists[i] = p.file.DeclList[j:]
				continue Outer // no more ImportDecls
			}
		}
	}

	// 2. Process all package-block type declarations. As with imports,
	// we need to make sure all types are properly instantiated before
	// trying to map any expressions that utilize them. In particular,
	// we need to make sure type pragmas are already known (see comment
	// in irgen.typeDecl).
	//
	// We could perhaps instead defer processing of package-block
	// variable initializers and function bodies, like noder does, but
	// special-casing just package-block type declarations minimizes the
	// differences between processing package-block and function-scoped
	// declarations.
	for _, declList := range declLists {
		for _, decl := range declList {
			switch decl := decl.(type) {
			case *syntax.TypeDecl:
				g.typeDecl((*ir.Nodes)(&g.target.Decls), decl)
			}
		}
	}
	types.ResumeCheckSize()

	// 3. Process all remaining declarations.
	for i, declList := range declLists {
		old := g.haveEmbed
		g.haveEmbed = noders[i].importedEmbed
		g.decls((*ir.Nodes)(&g.target.Decls), declList)
		g.haveEmbed = old
	}
	g.exprStmtOK = true

	// 4. Run any "later" tasks. Avoid using 'range' so that tasks can
	// recursively queue further tasks. (Not currently utilized though.)
	for len(g.laterFuncs) > 0 {
		fn := g.laterFuncs[0]
		g.laterFuncs = g.laterFuncs[1:]
		fn()
	}

	if base.Flag.W > 1 {
		for _, n := range g.target.Decls {
			s := fmt.Sprintf("\nafter noder2 %v", n)
			ir.Dump(s, n)
		}
	}

	for _, p := range noders {
		// Process linkname and cgo pragmas.
		p.processPragmas()

		// Double check for any type-checking inconsistencies. This can be
		// removed once we're confident in IR generation results.
		syntax.Crawl(p.file, func(n syntax.Node) bool {
			g.validate(n)
			return false
		})
	}

	if base.Flag.Complete {
		for _, n := range g.target.Decls {
			if fn, ok := n.(*ir.Func); ok {
				if fn.Body == nil && fn.Nname.Sym().Linkname == "" {
					base.ErrorfAt(fn.Pos(), "missing function body")
				}
			}
		}
	}

	// Check for unusual case where noder2 encounters a type error that types2
	// doesn't check for (e.g. notinheap incompatibility).
	base.ExitIfErrors()

	typecheck.DeclareUniverse()

	// Create any needed instantiations of generic functions and transform
	// existing and new functions to use those instantiations.
	BuildInstantiations()

	// Remove all generic functions from g.target.Decl, since they have been
	// used for stenciling, but don't compile. Generic functions will already
	// have been marked for export as appropriate.
	j := 0
	for i, decl := range g.target.Decls {
		if decl.Op() != ir.ODCLFUNC || !decl.Type().HasTParam() {
			g.target.Decls[j] = g.target.Decls[i]
			j++
		}
	}
	g.target.Decls = g.target.Decls[:j]

	base.Assertf(len(g.laterFuncs) == 0, "still have %d later funcs", len(g.laterFuncs))
}
```

将 Decl 接口转换为 Node 接口：

```go
// A Node is the abstract interface to an IR node.
type Node interface {
	// Formatting
	Format(s fmt.State, verb rune)

	// Source position.
	Pos() src.XPos
	SetPos(x src.XPos)

	// For making copies. For Copy and SepCopy.
	copy() Node

	doChildren(func(Node) bool) bool
	editChildren(func(Node) Node)

	// Abstract graph structure, for generic traversals.
	Op() Op
	Init() Nodes

	// Fields specific to certain Ops only.
	Type() *types.Type
	SetType(t *types.Type)
	Name() *Name
	Sym() *types.Sym
	Val() constant.Value
	SetVal(v constant.Value)

	// Storage for analysis passes.
	Esc() uint16
	SetEsc(x uint16)

	// Typecheck values:
	//  0 means the node is not typechecked
	//  1 means the node is completely typechecked
	//  2 means typechecking of the node is in progress
	//  3 means the node has its type from types2, but may need transformation
	Typecheck() uint8
	SetTypecheck(x uint8)
	NonNil() bool
	MarkNonNil()
}
```

```go

```
