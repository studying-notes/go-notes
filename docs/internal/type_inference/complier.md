---
date: 2022-10-02T15:00:55+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "类型推断原理"  # 文章标题
url:  "posts/go/docs/internal/type_inference/complier"  # 设置网页永久链接
tags: [ "Go", "complier" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 解析阶段

在词法解析阶段，会将赋值语句右边的常量解析为一个未定义的类型，例如，ImagLit 代表复数，FloatLit 代表浮点数，IntLit 代表整数。

Go 语言源代码采用 UTF-8 的编码方式，在进行词法解析时，当遇到需要赋值的常量操作时，会逐个读取后面常量的 UTF-8 字符。字符串的首字符为 "，数字的首字符为 '0' ～ '9'。具体实现位于 syntax.next 函数中。

`cmd/compile/internal/syntax/sxanner.go`

```go
func (s *scanner) next() {
...
	switch s.ch {
	case -1:
		if nlsemi {
			s.lit = "EOF"
			s.tok = _Semi
			break
		}
		s.tok = _EOF

	case '\n':
		s.nextch()
		s.lit = "newline"
		s.tok = _Semi

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number(false)

	case '"':
		s.stdString()

	case '`':
		s.rawString()
...
}
```

因此对于整数、小数等常量的识别就显得非常简单。

如图 3-2 所示，整数就是字符中全是 0～9 的数字，浮点数就是字符中有“.”号的数字，字符串的首字符为 " 或 '。

![](../../../assets/images/docs/internal/type_inference/complier/图3-2%20词法解析阶段解析未定义的常量示例.png)

下面列出的 number 函数为语法分析阶段处理数字的具体实现。数字首先会被分为小数部分与整数部分，通过字符. 进行区分。如果整数部分是以 0 开头的，则可能有不同的含义，例如 0x 代表十六进制数、0b 代表二进制数。

```go
	// integer part
	if !seenPoint {
		if s.ch == '0' {
			s.nextch()
			switch lower(s.ch) {
			case 'x':
				s.nextch()
				base, prefix = 16, 'x'
			case 'o':
				s.nextch()
				base, prefix = 8, 'o'
			case 'b':
				s.nextch()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				digsep = 1 // leading 0
			}
		}
		digsep |= s.digits(base, &invalid)
		if s.ch == '.' {
			if prefix == 'o' || prefix == 'b' {
				s.errorf("invalid radix point in %s literal", baseName(base))
				ok = false
			}
			s.nextch()
			seenPoint = true
		}
	}

	// fractional part
	if seenPoint {
		kind = FloatLit
		digsep |= s.digits(base, &invalid)
	}
```

以赋值语句 `a := 333` 为例，完成词法解析与语法分析时，此赋值语句将以 AssignStmt 结构表示。

`internal/syntax/nodes.go`

```go
AssignStmt struct {
	Op       Operator // 0 means no operation
	Lhs, Rhs Expr     // Rhs == nil means Lhs++ (Op == Add) or Lhs-- (Op == Sub)
	simpleStmt
}
```

其中 Op 代表操作符，在这里是赋值操作 OAS。Lhs 与 Rhs 分别代表左右两个表达式，左边代表变量 a，右边代表常量 333，此时其类型为 intLit。

## 类型检查与构建中间表示

完成解析后，进入类型检查与构建中间表示阶段。在该阶段会将解析阶段生成的 AssignStmt 结构解析为一个 Node 接口。

```go
func (g *irgen) decls(res *ir.Nodes, decls []syntax.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *syntax.ConstDecl:
			g.constDecl(res, decl)
		case *syntax.FuncDecl:
			g.funcDecl(res, decl)
		case *syntax.TypeDecl:
			if ir.CurFunc == nil {
				continue // already handled in irgen.generate
			}
			g.typeDecl(res, decl)
		case *syntax.VarDecl:
			g.varDecl(res, decl)
		default:
			g.unhandled("declaration", decl)
		}
	}
}
```

```go
// An AssignListStmt is an assignment statement with
// more than one item on at least one side: Lhs = Rhs.
// If Def is true, the assignment is a :=.
type AssignListStmt struct {
	miniStmt
	Lhs Nodes
	Def bool
	Rhs Nodes
}

// A miniStmt is a miniNode with extra fields common to statements.
type miniStmt struct {
	miniNode
	init Nodes
}

// A miniNode is a minimal node implementation,
// meant to be embedded as the first field in a larger node implementation,
// at a cost of 8 bytes.
//
// A miniNode is NOT a valid Node by itself: the embedding struct
// must at the least provide:
//
//	func (n *MyNode) String() string { return fmt.Sprint(n) }
//	func (n *MyNode) rawCopy() Node { c := *n; return &c }
//	func (n *MyNode) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
//
// The embedding struct should also fill in n.op in its constructor,
// for more useful panic messages when invalid methods are called,
// instead of implementing Op itself.
type miniNode struct {
	pos  src.XPos // uint32
	op   Op       // uint8
	bits bitset8
	esc  uint16
}
```

AssignListStmt 实现了 Node 接口，其中 Lhs 与 Rhs 分别代表左右两个表达式，左边代表变量 a，右边代表常量 333，此时其类型为 intLit。其中 op 操作为 OLITERAL。
