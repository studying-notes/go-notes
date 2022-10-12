---
date: 2022-10-01T09:15:55+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "解析"  # 文章标题
url:  "posts/go/docs/internal/compiler/parsing"  # 设置网页永久链接
tags: [ "Go", "parsing" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

这一阶段的相关源码位于 `cmd/compile/internal/syntax` (lexer, parser, syntax tree) 文件夹下。

在编译的第一阶段，对源代码进行标记化（词法分析）、解析（语法分析），并为每个源文件构建语法树。

每个语法树都是相应源文件的精确表示，其节点对应于源文件的各种元素，例如表达式、声明和语句。语法树还包括位置信息，用于错误报告和调试信息的创建。

编译时的入口点为 `noder.LoadPackage()` 函数，位于 `src/cmd/compile/internal/noder/noder.go` 文件中。

```go
p.file, _ = syntax.Parse(fbase, f, p.error, p.pragma, syntax.CheckBranches) // errors are tracked via p.error
```

```go
// Parse parses a single Go source file from src and returns the corresponding
// syntax tree. If there are errors, Parse will return the first error found,
// and a possibly partially constructed syntax tree, or nil.
//
// If errh != nil, it is called with each error encountered, and Parse will
// process as much source as possible. In this case, the returned syntax tree
// is only nil if no correct package clause was found.
// If errh is nil, Parse will terminate immediately upon encountering the first
// error, and the returned syntax tree is nil.
//
// If pragh != nil, it is called with each pragma encountered.
func Parse(base *PosBase, src io.Reader, errh ErrorHandler, pragh PragmaHandler, mode Mode) (_ *File, first error) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(Error); ok {
				first = err
				return
			}
			panic(p)
		}
	}()

	var p parser
	p.init(base, src, errh, pragh, mode)
	p.next()
	return p.fileOrNil(), p.first
}
```

```go
// Package files
//
// Parse methods are annotated with matching Go productions as appropriate.
// The annotations are intended as guidelines only since a single Go grammar
// rule may be covered by multiple parse methods and vice versa.
//
// Excluding methods returning slices, parse methods named xOrNil may return
// nil; all others are expected to return a valid non-nil node.

// SourceFile = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
func (p *parser) fileOrNil() *File
```

## 词法解析

在词法解析阶段，Go 语言编译器会**扫描**输入的 Go 源文件，并将其符号（token）化。

例如 + 和 - 操作符会被转换为 `_IncOp`，赋值符号 := 会被转换为 `_Define`。

这些 token 实质上是用 iota 声明的整数，定义在 `syntax/tokens.go` 中。

```go
const (
	_    token = iota
	_EOF       // EOF

	// names and literals
	_Name    // name
	_Literal // literal

	// operators and operations
	// _Operator is excluding '*' (_Star)
	_Operator // op
	_AssignOp // op=
	_IncOp    // opop
	_Assign   // =
	_Define   // :=
	_Arrow    // <-
	_Star     // *

	// delimiters
	_Lparen    // (
	_Lbrack    // [
	_Lbrace    // {
	_Rparen    // )
	_Rbrack    // ]
	_Rbrace    // }
	_Comma     // ,
	_Semi      // ;
	_Colon     // :
	_Dot       // .
	_DotDotDot // ...

	// keywords
	_Break       // break
	_Case        // case
	_Chan        // chan
	_Const       // const
	_Continue    // continue
	_Default     // default
	_Defer       // defer
	_Else        // else
	_Fallthrough // fallthrough
	_For         // for
	_Func        // func
	_Go          // go
	_Goto        // goto
	_If          // if
	_Import      // import
	_Interface   // interface
	_Map         // map
	_Package     // package
	_Range       // range
	_Return      // return
	_Select      // select
	_Struct      // struct
	_Switch      // switch
	_Type        // type
	_Var         // var

	// empty line comment to exclude it from .String
	tokenCount //
)
```

通过 stringer 工具，为 token 生成了 String() 方法：

```go
//go:generate stringer -type token -linecomment tokens.go
```

stringer 可以通过下面的命令安装：

```
go get golang.org/x/tools/cmd/stringer
```

符号化保留了 Go 语言中定义的符号，可以识别出错误的拼写。同时，字符串被转换为整数后，在后续的阶段中能够被更加高效地处理。

![](../../../assets/images/docs/internal/compiler/parsing/图1-3%20Go语言编译器词法解析示例.png)

图 1-3 为一个示例，展现了将表达式 `a:=b+c(12)` 符号化之后的情形。代码中声明的标识符、关键字、运算符和分隔符等字符串都可以转化为对应的符号。

编译器用的扫描函数在 `src/cmd/compile/internal/syntax/scanner.go` 中。

```go
type scanner struct {
	source
	mode   uint
	nlsemi bool // if set '\n' and EOF translate to ';'

	// current token, valid after calling next()
	line, col uint
	blank     bool // line is blank up to col
	tok       token
	lit       string   // valid if tok is _Name, _Literal, or _Semi ("semicolon", "newline", or "EOF"); may be malformed if bad is true
	bad       bool     // valid if tok is _Literal, true if a syntax error occurred, lit may be malformed
	kind      LitKind  // valid if tok is _Literal
	op        Operator // valid if tok is _Operator, _Star, _AssignOp, or _IncOp
	prec      int      // valid if tok is _Operator, _Star, _AssignOp, or _IncOp
}

type source struct {
	in   io.Reader
	errh func(line, col uint, msg string)

	buf       []byte // source buffer
	ioerr     error  // pending I/O error, or nil
	b, r, e   int    // buffer indices (see comment above)
	line, col uint   // source position of ch (0-based)
	ch        rune   // most recently read character
	chw       int    // width of ch
}
```

scanner 结构体中的 source 字段是一个 io.Reader 接口，用于读取源文件。errh 字段是一个函数，用于处理错误。buf 字段是一个字节切片，用于存储源文件的内容。ioerr 字段是一个 error 类型，用于存储 I/O 错误。b、r、e 字段是三个整数，用于存储 buf 字段的索引。line、col 字段是两个整数，用于存储当前字符的行号和列号。ch 字段是一个 rune 类型，用于存储当前字符。chw 字段是一个整数，用于存储当前字符的宽度。

scanner 结构体中的 mode 字段是一个整数，用于存储扫描模式。nlsemi 字段是一个布尔值，用于指示是否将换行符和文件结束符转换为分号。

scanner 结构体中的 line、col、blank、tok、lit、bad、kind、op、prec 字段是用于存储当前符号的相关信息。

scanner 结构体中的 next() 方法用于获取下一个符号，next() 方法将当前符号的信息存储到 scanner 结构体中的 line、col、blank、tok、lit、bad、kind、op、prec 字段中。这些字段的值在调用 next() 方法后会被覆盖更新，所以在调用 next() 方法之前，需要先保存或处理这些字段的值。

初始化函数将枚举量和字符串的哈希值做了映射，这样就可以通过字符串的哈希值来获取枚举量。

```go
// hash is a perfect hash function for keywords.
// It assumes that s has at least length 2.
func hash(s []byte) uint {
	return (uint(s[0])<<4 ^ uint(s[1]) + uint(len(s))) & uint(len(keywordMap)-1)
}

var keywordMap [1 << 6]token // size must be power of two

func init() {
	// populate keywordMap
	for tok := _Break; tok <= _Var; tok++ {
		h := hash([]byte(tok.String()))
		if keywordMap[h] != 0 {
			panic("imperfect hash")
		}
		keywordMap[h] = tok
	}
}
```

下面的函数就是通过哈希值来获取枚举量作为 tok 值。

```go
func (s *scanner) ident() {
	// accelerate common case (7bit ASCII)
	for isLetter(s.ch) || isDecimal(s.ch) {
		s.nextch()
	}

	// general case
	if s.ch >= utf8.RuneSelf {
		for s.atIdentChar(false) {
			s.nextch()
		}
	}

	// possibly a keyword
	lit := s.segment()
	if len(lit) >= 2 {
		if tok := keywordMap[hash(lit)]; tok != 0 && tokStrFast(tok) == string(lit) {
			s.nlsemi = contains(1<<_Break|1<<_Continue|1<<_Fallthrough|1<<_Return, tok)
			s.tok = tok
			return
		}
	}

	s.nlsemi = true
	s.lit = string(lit)
	s.tok = _Name
}
```

另外，Go 语言标准库 `go/scanner`、`go/token` 提供了接口用于非编译期扫描源代码。

在下例中，我们将使用这些接口模拟对 Go 文本文件的扫描。

```go
package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	src := []byte("cos(x) + 2i*sin(x) // Euler")

	var s scanner.Scanner
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("", fileSet.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", file.Position(pos), tok, lit)
	}
}
```

```
1:1     IDENT   "cos"
1:4     (       ""
1:5     IDENT   "x"
1:6     )       ""
1:8     +       ""
1:10    IMAG    "2i"
1:12    *       ""
1:13    IDENT   "sin"
1:16    (       ""
1:17    IDENT   "x"
1:18    )       ""
1:20    ;       "\n"
1:20    COMMENT "// Euler"
```

在上例中，src 为进行词法扫描的表达式，可以将其模拟为一个文件并调用 scanner.Scanner 词法，扫描后分别打印出 token 的位置、符号及其字符串字面量。每个标识符与运算符都被特定的 token 代替，例如 2i 被识别为复数 IMAG，注释被识别为 COMMENT。

Go 中是允许语句以分号 `;` 结尾的，其实 Go 语言的编译器会在语义上将换行符认为是分号，所以实际上每条语句在语义上仍以分号 `;` 结尾。

```go
	case '\n':
		s.nextch()
		s.lit = "newline"
		s.tok = _Semi

	case ';':
		s.nextch()
		s.lit = "semicolon"
		s.tok = _Semi
```

每个语句都必须以分号 `;` `token` 结尾。

```go
	// { ImportDecl ";" }
	for p.got(_Import) {
		f.DeclList = p.appendGroup(f.DeclList, p.importDecl)
		p.want(_Semi)
	}
```

但是，如果在一行中有多个语句，那么就必须在每个语句后面加上分号。

## 语法解析

从源码看，词法解析和语法解析是交替进行的，且由语法解析驱动。词法解析完一个 token 后，语法解析器判断是否继续。

结构体 parser 用于语法解析，其定义如下：

```go
type parser struct {
	file  *PosBase // source file handle
	errh  ErrorHandler // error reporting; or nil
	mode  Mode // parsing mode
	pragh PragmaHandler // #pragma handler; or nil
	scanner // embedded scanner

	base   *PosBase // current position base
	first  error    // first error encountered
	errcnt int      // number of errors encountered
	pragma Pragma   // pragmas

	fnest  int    // function nesting level (for error handling)
	xnest  int    // expression nesting level (for complit ambiguity resolution)
	indent []byte // tracing support
}
```

Go 语言采用了标准的**自上而下的递归下降**（Top-Down Recursive-Descent）算法，以简单高效的方式完成无须回溯的语法扫描，核心算法位于 `syntax/nodes.go` 及 `syntax/parser.go` 中。

![](../../../assets/images/docs/internal/compiler/parsing/图1-4%20Go语言编译器对文件进行语法解析的示意图.png)

图 1-4 为 Go 语言编译器对文件进行语法解析的示意图。在一个 Go 源文件中主要有包导入声明（import）、静态常量（const）、类型声明（type）、变量声明（var）及函数声明。

源文件中的每一种声明都有对应的语法，递归下降通过识别初始的标识符，采用对应的语法进行解析。这种方式能够较快地解析并识别可能出现的语法错误。

每一种声明语法在 Go 语言规范中都有定义。

```go
//包导入声明
ImportSpec = [ "." | PackageName ] ImportPath .
ImportPath = string_lit .
//静态常量
ConstSpec = IdentifierList [ [ Type ] "=" ExpressionList ] .
//类型声明
TypeSpec = identifier [ "=" ] Type .
//变量声明
VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
```

函数声明是文件中最复杂的一类语法，因为在函数体的内部可能有多种声明、赋值（例如 :=）、表达式及函数调用等。例如 defer 语法为 `defer Expression`，其后必须跟一个函数或方法。

每一种声明语法或者表达式都有对应的结构体，例如 `a:=b+f(89)` 对应的结构体为赋值声明 `AssignStmt`。Op 代表当前的操作符，即 “:=”，`X` 与 `Y` 分别代表左右两个表达式。

```go
// An AssignStmt is a simple assignment statement: X = Y.
// If Def is true, the assignment is a :=.
type AssignStmt struct {
	miniStmt
	X   Node
	Def bool
	Y   Node
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

语法解析丢弃了一些不重要的标识符，例如括号 “(”，并将语义存储到了对应的结构体中。语法声明的结构体拥有对应的层次结构，这是构建抽象语法树的基础。

图 1-5 为 `a:= b+c(12)` 语句被语法解析后转换为对应的 `AssignStmt` 结构体之后的情形。最顶层的 Op 操作符为 `Def(:=)`。X 表达式类型为标识符 Name，值为标识符 “a”。Y 表达式为 Operator 加法运算。加法运算左边为标识符 “b”，右边为函数调用表达式，类型为 `CallExpr`。其中，函数名 c 的类型为 Name，参数为常量类型 `BasicLit`，代表数字 12。

![](../../../assets/images/docs/internal/compiler/parsing/图1-5%20特定表达式的语法解析示例.png)

上面入口最后返回了 `fileOrNil` 函数的处理结果，返回值的结构体为 `File`，其包含了整个文件的所有语法声明。`File` 结构体中包含了一个 `Decls` 字段，其类型为 `Nodes`，是一个结构体数组，其中包含了所有的语法声明。

`internal/syntax/nodes.go`

```go
// package PkgName; DeclList[0], DeclList[1], ...
type File struct {
	Pragma   Pragma
	PkgName  *Name
	DeclList []Decl
	EOF      Pos
	node
}
```

```go
// SourceFile = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
func (p *parser) fileOrNil() *File {
	if trace {
		defer p.trace("file")()
	}

	f := new(File)
	f.pos = p.pos()

	// PackageClause
	if !p.got(_Package) {
		p.syntaxError("package statement must be first")
		return nil
	}
	f.Pragma = p.takePragma()
	f.PkgName = p.name()
	p.want(_Semi)

	// don't bother continuing if package clause has errors
	if p.first != nil {
		return nil
	}

	// { ImportDecl ";" }
	for p.got(_Import) {
		f.DeclList = p.appendGroup(f.DeclList, p.importDecl)
		p.want(_Semi)
	}

	// { TopLevelDecl ";" }
	for p.tok != _EOF {
		switch p.tok {
		case _Const:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.constDecl)

		case _Type:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.typeDecl)

		case _Var:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.varDecl)

		case _Func:
			p.next()
			if d := p.funcDeclOrNil(); d != nil {
				f.DeclList = append(f.DeclList, d)
			}

		default:
			if p.tok == _Lbrace && len(f.DeclList) > 0 && isEmptyFuncDecl(f.DeclList[len(f.DeclList)-1]) {
				// opening { of function declaration on next line
				p.syntaxError("unexpected semicolon or newline before {")
			} else {
				p.syntaxError("non-declaration statement outside function body")
			}
			p.advance(_Const, _Type, _Var, _Func)
			continue
		}

		// Reset p.pragma BEFORE advancing to the next token (consuming ';')
		// since comments before may set pragmas for the next function decl.
		p.clearPragma()

		if p.tok != _EOF && !p.got(_Semi) {
			p.syntaxError("after top level declaration")
			p.advance(_Const, _Type, _Var, _Func)
		}
	}
	// p.tok == _EOF

	p.clearPragma()
	f.EOF = p.pos()

	return f
}
```

该函数执行完毕后，词法解析、语法解析和抽象语法树构建全部完成，这三个步骤其实是同时进行的。

## 抽象语法树构建

编译器前端必须构建程序的中间表示形式，以便在编译器中间阶段及后端使用。抽象语法树（Abstract Syntax Tree，AST）是一种常见的树状结构的中间态。

在 Go 语言源文件中的任何一种 import、type、const、func 声明都是一个根节点，在根节点下包含当前声明的子节点。节点的接口定义如下：

```go
type Node interface {
	// Pos() returns the position associated with the node as follows:
	// 1) The position of a node representing a terminal syntax production
	//    (Name, BasicLit, etc.) is the position of the respective production
	//    in the source.
	// 2) The position of a node representing a non-terminal production
	//    (IndexExpr, IfStmt, etc.) is the position of a token uniquely
	//    associated with that production; usually the left-most one
	//    ('[' for IndexExpr, 'if' for IfStmt, etc.)
	Pos() Pos
	aNode()
}
```

```go
// Declarations

type (
	Decl interface {
		Node
		aDecl()
	}

	//              Path
	// LocalPkgName Path
	ImportDecl struct {
		Group        *Group // nil means not part of a group
		Pragma       Pragma
		LocalPkgName *Name     // including "."; nil means no rename present
		Path         *BasicLit // Path.Bad || Path.Kind == StringLit; nil means no path
		decl
	}

	// NameList
	// NameList      = Values
	// NameList Type = Values
	ConstDecl struct {
		Group    *Group // nil means not part of a group
		Pragma   Pragma
		NameList []*Name
		Type     Expr // nil means no type
		Values   Expr // nil means no values
		decl
	}

	// Name Type
	TypeDecl struct {
		Group      *Group // nil means not part of a group
		Pragma     Pragma
		Name       *Name
		TParamList []*Field // nil means no type parameters
		Alias      bool
		Type       Expr
		decl
	}

	// NameList Type
	// NameList Type = Values
	// NameList      = Values
	VarDecl struct {
		Group    *Group // nil means not part of a group
		Pragma   Pragma
		NameList []*Name
		Type     Expr // nil means no type
		Values   Expr // nil means no values
		decl
	}

	// func          Name Type { Body }
	// func          Name Type
	// func Receiver Name Type { Body }
	// func Receiver Name Type
	FuncDecl struct {
		Pragma     Pragma
		Recv       *Field // nil means regular function
		Name       *Name
		TParamList []*Field // nil means no type parameters
		Type       *FuncType
		Body       *BlockStmt // nil means no body (forward declaration)
		decl
	}
)
```

以 `a:= b+c(12)` 为例，该赋值语句最终会变为如图 1-6 所示的抽象语法树。节点之间具有从上到下的层次结构和依赖关系。

![](../../../assets/images/docs/internal/compiler/parsing/图1-6%20抽象语法树.png)

```go

```
