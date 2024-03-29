---
date: 2022-10-01T10:41:31+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "类型检查"  # 文章标题
url:  "posts/go/docs/internal/compiler/type_checking"  # 设置网页永久链接
tags: [ "Go", "type-checking" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

完成抽象语法树的初步构建后，就进入类型检查阶段遍历节点树并决定节点的类型。

```go
// Use types2 to type-check and generate IR.
check2(noders)
```

每个 noder 包含了语法解析后得到的 `file *syntax.File`。

```go
// check2 type checks a Go package using types2, and then generates IR
// using the results.
func check2(noders []*noder) {
	m, pkg, info := checkFiles(noders)

	g := irgen{
		target: typecheck.Target,
		self:   pkg,
		info:   info,
		posMap: m,
		objs:   make(map[types2.Object]*ir.Name),
		typs:   make(map[types2.Type]*types.Type),
	}
	g.generate(noders)
}

func (g *irgen) generate(noders []*noder)
```

这其中包括了语法中明确指定的类型，例如 var a int，也包括了需要通过编译器类型推断得到的类型。例如，a:=1 中的变量 a 与常量 1 都未直接声明类型，编译器会自动推断出节点常量 1 的类型为 TINT()，并自动推断出 a 的类型为 TINT()。

在类型检查阶段，会对一些类型做特别的语法或语义检查。例如：引用的结构体字段是否是大写可导出的？数组字面量的访问是否超过了其长度？数组的索引是不是正整数？

除此之外，在类型检查阶段还会进行其他工作。例如：计算编译时常量、将标识符与声明绑定等。

```go

```
