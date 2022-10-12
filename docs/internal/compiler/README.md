---
date: 2022-09-30T19:34:16+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "深入 Go 语言编译器"  # 文章标题
url:  "posts/go/docs/internal/compiler/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

以 .go 为后缀的 UTF-8 格式的 Go 文本文件最终能被编译成特定机器上的可执行文件，离不开 Go 语言编译器的复杂工作。

Go 语言编译器不仅能准确地翻译高级语言，也能进行代码优化。

## 为什么要了解 Go 语言编译器

编译器是一个大型且复杂的系统，一个好的编译器会很好地结合形式语言理论、算法、人工智能、系统设计、计算机体系结构及编程语言理论。

Go 语言的编译器遵循了主流编译器采用的经典策略及相似的处理流程和优化规则，例如**经典的递归下降的语法解析**、**抽象语法树的构建**。

另外，Go 语言编译器有一些特殊的设计，例如**内存的逃逸**等。

通过了解 Go 语言编辑器，不仅可以了解大部分高级语言编译器的一般性流程与规则，也能指导我们写出更加优秀的程序。

可以通过禁用编译器的优化以及内联函数等特性调试和查看代码的执行流程。

很多 Go 语言的语法特性都离不开编译时与运行时的共同作用。

另外，如果读者希望开发 go import、go fmt、go lint 等扫描源代码的工具，那么同样离不开编译器的知识和 Go 语言提供的 API。

## 三阶段编译器

![](../../../assets/images/docs/internal/compiler/README/图1-1%20三阶段编译器.png)

如图 1 -1 所示，在经典的编译原理中，一般将编译器分为**编译器前端**、**优化器**和**编译器后端**。这种编译器被称为**三阶段编译器**（three-phase compiler）。

- 编译器前端主要专注于理解源程序、扫描解析源程序并进行精准的语义表达。
- 编译器的中间阶段（Intermediate Representation，IR）可能有多个，编译器会使用多个 IR 阶段、多种数据结构表示代码，并在中间阶段对代码进行多次**优化**。例如，识别冗余代码、识别内存逃逸等。编译器的中间阶段离不开编译器前端记录的细节。
- 编译器后端专注于**生成特定目标机器上的程序**，这种程序可能是可执行文件，也可能是需要进一步处理的中间形态 obj 文件、汇编语言等。

编译器优化并不是一个非常明确的概念。优化的主要目的一般是降低程序资源的消耗，比较常见的是降低内存与 CPU 的使用率。但在很多时候，这些目标可能是相互冲突的，对一个目标的优化可能降低另一个目标的效率。同时，理论已经表明有一些代码优化存在着 NP 难题，这意味着随着代码的增加，优化的难度将越来越大，需要花费的时间呈指数增长。因为这些原因，编译器无法进行最佳的优化，所以通常采用一种折中的方案。

## Go 语言的编译器

Go 语言编译器一般缩写为小写的 gc（go compiler），需要和大写的 GC（垃圾回收）进行区分。

Go 语言编译器的执行流程可细化为多个阶段，包括词法解析、语法解析、抽象语法树构建、类型检查、变量捕获、函数内联、逃逸分析、闭包重写、遍历函数、SSA 生成、机器码生成，如图 1 -2 所示。

![](../../../assets/images/docs/internal/compiler/README/图1-2%20Go语言编译器执行流程.png)

官方文档 `cmd/compile/README.md` 将编译过程分为 7 个阶段：

- Parsing：词法解析、语法解析和抽象语法树构建。
- Type checking：类型检查。
- IR construction：构建中间表示。
- Middle end：中间阶段，包括变量捕获、函数内联、逃逸分析、闭包重写。
- Walk：遍历函数。
- Generic SSA：生成 SSA。
- Generating machine code：生成机器码。

## 编译入口和基本流程

和 Go 语言编译器有关的代码主要位于 `src/cmd/compile` 目录下，在后面分析中给出的文件路径均默认位于该目录中。

入口程序为 `src/cmd/compile/main.go`，该模块可以单独编译。进入该目录，执行 `go build` 命令即可编译出可执行文件 `compile.exe`。

```go
var archInits = map[string]func(*ssagen.ArchInfo){
	"386":      x86.Init,
	"amd64":    amd64.Init,
	"arm":      arm.Init,
	"arm64":    arm64.Init,
	"loong64":  loong64.Init,
	"mips":     mips.Init,
	"mipsle":   mips.Init,
	"mips64":   mips64.Init,
	"mips64le": mips64.Init,
	"ppc64":    ppc64.Init,
	"ppc64le":  ppc64.Init,
	"riscv64":  riscv64.Init,
	"s390x":    s390x.Init,
	"wasm":     wasm.Init,
}

func main() {
	// disable timestamps for reproducible output
	log.SetFlags(0)
	log.SetPrefix("compile: ")

	buildcfg.Check()
	archInit, ok := archInits[buildcfg.GOARCH]
	if !ok {
		fmt.Fprintf(os.Stderr, "compile: unknown architecture %q\n", buildcfg.GOARCH)
		os.Exit(2)
	}

	gc.Main(archInit)
	base.Exit(0)
}
```

### 架构初始化函数

`main()` 函数根据配置选择编译目标架构的初始化函数。以 amd64 为例：

`internal/amd64/galign.go`

```go
func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &x86.Linkamd64
	arch.REGSP = x86.REGSP
	arch.MAXWIDTH = 1 << 50

	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
	arch.LoadRegResult = loadRegResult
	arch.SpillArgReg = spillArgReg
}
```

然后调用了 `gc.Main()` 函数，该函数位于 `src/cmd/compile/internal/gc/main.go` 文件中，该文件中定义了 `Main()` 函数。

```go
// Main parses flags and Go source files specified in the command-line
// arguments, type-checks the parsed Go package, compiles functions to machine
// code, and finally writes the compiled package definition to disk.
func Main(archInit func(*ssagen.ArchInfo))
```

首先 `Main()` 函数会初始化参数，然后解析命令行参数等。

### 创建包结构体

```go
typecheck.Target = new(ir.Package)
```

```go
// A Package holds information about the package being compiled.
type Package struct {
	// Imports, listed in source order.
	// See golang.org/issue/31636.
	Imports []*types.Pkg

	// Init functions, listed in source order.
	Inits []*Func

	// Top-level declarations.
	Decls []Node

	// Extern (package global) declarations.
	Externs []Node

	// Assembly function declarations.
	Asms []*Name

	// Cgo directives.
	CgoPragmas [][]string

	// Variables with //go:embed lines.
	Embeds []*Name

	// Exported (or re-exported) symbols.
	Exports []*Name
}
```

### 解析和类型检查

```go
// Parse and typecheck input.
noder.LoadPackage(flag.Args())
```

`noder.LoadPackage()` 函数位于 `src/cmd/compile/internal/noder/noder.go` 文件中，该函数的主要功能是解析和类型检查输入的 Go 源文件。

```go
// noder transforms package syntax's AST into a Node tree.
type noder struct {
	posMap

	file           *syntax.File
	linknames      []linkname
	pragcgobuf     [][]string
	err            chan syntax.Error
	importedUnsafe bool
	importedEmbed  bool
}

func LoadPackage(filenames []string) {
	base.Timer.Start("fe", "parse")

	// Limit the number of simultaneously open files.
	sem := make(chan struct{}, runtime.GOMAXPROCS(0)+10)

	noders := make([]*noder, len(filenames))
	for i := range noders {
		p := noder{
			err: make(chan syntax.Error),
		}
		noders[i] = &p
	}

	// Move the entire syntax processing logic into a separate goroutine to avoid blocking on the "sem".
	go func() {
		for i, filename := range filenames {
			filename := filename
			p := noders[i]
			sem <- struct{}{}  // acquire token
			go func() {
				defer func() { <-sem }()  // release token
				defer close(p.err)// close error channel when done
				fbase := syntax.NewFileBase(filename)

				f, err := os.Open(filename)
				if err != nil {
					p.error(syntax.Error{Msg: err.Error()})
					return
				}
				defer f.Close()

				// Parse parses a single Go source file from src and returns the corresponding
				// syntax tree. If there are errors, Parse will return the first error found,
				// and a possibly partially constructed syntax tree, or nil.
				p.file, _ = syntax.Parse(fbase, f, p.error, p.pragma, syntax.CheckBranches) // errors are tracked via p.error
			}()
		}
	}()

	var lines uint
	for _, p := range noders {
		// Wait for the goroutine to finish.
		for e := range p.err {
			p.errorAt(e.Pos, "%s", e.Msg)
		}
		if p.file == nil {
			base.ErrorExit()
		}
		lines += p.file.EOF.Line()
	}
	base.Timer.AddEvent(int64(lines), "lines")

	if base.Debug.Unified != 0 {
		unified(noders)
		return
	}

	// Use types2 to type-check and generate IR.
	check2(noders)
}
```

[解析](parsing.md)

[类型检查](type_checking.md)

[构建中间表示](ir.md)

### 创建初始化函数

```go
// Create "init" function for package-scope variable initialization
// statements, if any.
//
// Note: This needs to happen early, before any optimizations. The
// Go spec defines a precise order than initialization should be
// carried out in, and even mundane optimizations like dead code
// removal can skew the results (e.g., #43444).
pkginit.MakeInit()
```

`MakeInit()` 函数位于 `src/cmd/compile/internal/pkginit/init.go` 文件中，该函数的主要功能是创建用于初始化包级变量的 `init` 函数。

```go
// MakeInit creates a synthetic init function to handle any
// package-scope initialization statements.
//
// TODO(mdempsky): Move into noder, so that the types2-based frontends
// can use Info.InitOrder instead.
func MakeInit() {
	nf := initOrder(typecheck.Target.Decls)
	if len(nf) == 0 {
		return
	}

	// Make a function that contains all the initialization statements.
	base.Pos = nf[0].Pos() // prolog/epilog gets line number of first init stmt
	initializers := typecheck.Lookup("init")
	fn := typecheck.DeclFunc(initializers, nil, nil, nil)
	for _, dcl := range typecheck.InitTodoFunc.Dcl {
		dcl.Curfn = fn
	}
	fn.Dcl = append(fn.Dcl, typecheck.InitTodoFunc.Dcl...)
	typecheck.InitTodoFunc.Dcl = nil

	// Suppress useless "can inline" diagnostics.
	// Init functions are only called dynamically.
	fn.SetInlinabilityChecked(true)

	fn.Body = nf
	typecheck.FinishFuncBody()

	typecheck.Func(fn)
	ir.WithFunc(fn, func() {
		typecheck.Stmts(nf)
	})
	typecheck.Target.Decls = append(typecheck.Target.Decls, fn)

	// Prepend to Inits, so it runs first, before any user-declared init
	// functions.
	typecheck.Target.Inits = append([]*ir.Func{fn}, typecheck.Target.Inits...)

	if typecheck.InitTodoFunc.Dcl != nil {
		// We only generate temps using InitTodoFunc if there
		// are package-scope initialization statements, so
		// something's weird if we get here.
		base.Fatalf("InitTodoFunc still has declarations")
	}
	typecheck.InitTodoFunc = nil
}
```

`MakeInit()` 函数首先调用 `initOrder()` 函数，该函数的主要功能是将包级变量的初始化语句按照初始化顺序进行排序，然后将排序后的初始化语句放入 `init` 函数中，最后将 `init` 函数添加到 `typecheck.Target.Decls` 中。

```go
// initOrder computes initialization order for a list l of
// package-level declarations (in declaration order) and outputs the
// corresponding list of statements to include in the init() function
// body.
func initOrder(l []ir.Node) []ir.Node {
	var res ir.Nodes
	o := InitOrder{
		blocking: make(map[ir.Node][]ir.Node),
		order:    make(map[ir.Node]int),
	}

	// Process all package-level assignment in declaration order.
	for _, n := range l {
		switch n.Op() {
		case ir.OAS, ir.OAS2DOTTYPE, ir.OAS2FUNC, ir.OAS2MAPR, ir.OAS2RECV:
			o.processAssign(n)
			o.flushReady(func(n ir.Node) { res.Append(n) })
		case ir.ODCLCONST, ir.ODCLFUNC, ir.ODCLTYPE:
			// nop
		default:
			base.Fatalf("unexpected package-level statement: %v", n)
		}
	}

	// Check that all assignments are now Done; if not, there must
	// have been a dependency cycle.
	for _, n := range l {
		switch n.Op() {
		case ir.OAS, ir.OAS2DOTTYPE, ir.OAS2FUNC, ir.OAS2MAPR, ir.OAS2RECV:
			if o.order[n] != orderDone {
				// If there have already been errors
				// printed, those errors may have
				// confused us and there might not be
				// a loop. Let the user fix those
				// first.
				base.ExitIfErrors()

				o.findInitLoopAndExit(firstLHS(n), new([]*ir.Name), new(ir.NameSet))
				base.Fatalf("initialization unfinished, but failed to identify loop")
			}
		}
	}

	// Invariant consistency check. If this is non-zero, then we
	// should have found a cycle above.
	if len(o.blocking) != 0 {
		base.Fatalf("expected empty map: %v", o.blocking)
	}

	return res
}
```

### 中间端

[中间端](middle_end.md)

### 编译函数

```go
// Compile top level functions.
// Don't use range--walk can add functions to Target.Decls.
base.Timer.Start("be", "compilefuncs")
fcount := int64(0)
for i := 0; i < len(typecheck.Target.Decls); i++ {
if fn, ok := typecheck.Target.Decls[i].(*ir.Func); ok {
	// Don't try compiling dead hidden closure.
	if fn.IsDeadcodeClosure() {
		continue
	}
	enqueueFunc(fn)
	fcount++
}
}
base.Timer.AddEvent(fcount, "funcs")

compileFunctions()
```

### 写入目标文件

```go
// Write object data to disk.
base.Timer.Start("be", "dumpobj")
dumpdata()
base.Ctxt.NumberSyms()
dumpobj()
if base.Flag.AsmHdr != "" {
  dumpasmhdr()
}
```

```go

```
