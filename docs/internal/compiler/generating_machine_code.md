---
date: 2022-10-01T14:40:36+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "机器码生成"  # 文章标题
url:  "posts/go/docs/internal/compiler/generating_machine_code"  # 设置网页永久链接
tags: [ "Go", "generating-machine-code" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

* `cmd/compile/internal/ssa` (SSA lowering and arch-specific passes)
* `cmd/internal/obj` (machine code generation)

## 汇编器

在 SSA 阶段，编译器先执行与特定指令集无关的优化，再执行与特定指令集有关的优化，并最终生成与特定指令集有关的指令和寄存器分配方式。

在 SSA lower 阶段之后，就开始执行与特定指令集有关的重写与优化，在 genssa 阶段，编译器会生成与单个指令对应的 src/cmd/internal/obj/link.go 中的 Prog 结构。

```go
type Prog struct {
	Ctxt     *Link     // linker context
	Link     *Prog     // next Prog in linked list
	From     Addr      // first source operand
	RestArgs []AddrPos // can pack any operands that not fit into {Prog.From, Prog.To}
	To       Addr      // destination operand (second is RegTo2 below)
	Pool     *Prog     // constant pool entry, for arm,arm64 back ends
	Forwd    *Prog     // for x86 back end
	Rel      *Prog     // for x86, arm back ends
	Pc       int64     // for back ends or assembler: virtual or actual program counter, depending on phase
	Pos      src.XPos  // source position of this instruction
	Spadj    int32     // effect of instruction on stack pointer (increment or decrement amount)
	As       As        // assembler opcode
	Reg      int16     // 2nd source operand
	RegTo2   int16     // 2nd destination operand
	Mark     uint16    // bitmask of arch-specific items
	Optab    uint16    // arch-specific opcode index
	Scond    uint8     // bits that describe instruction suffixes (e.g. ARM conditions)
	Back     uint8     // for x86 back end: backwards branch state
	Ft       uint8     // for x86 back end: type index of Prog.From
	Tt       uint8     // for x86 back end: type index of Prog.To
	Isize    uint8     // for x86 back end: size of the instruction in bytes
}
```

例如，最终生成的指令 MOVL R1，R2 会被 Prog 表示为 As = MOVL，From = R1，To = R2。Pcond 代表跳转指令，除此之外，还有一些与特定指令集相关的结构。

在 SSA 后，编译器将调用与特定指令集有关的汇编器（Assembler）生成 obj 文件，obj 文件作为链接器（Linker）的输入，生成二进制可执行文件。internal/obj 目录中包含了汇编与链接的核心逻辑，内部有许多与机器码生成相关的包。不同类型的指令集（amd64、arm64、mips64 等）需要使用不同的包生成。Go 语言目前能在所有常见的 CPU 指令集类型上编译运行。

汇编和链接是编译器后端与特定指令集有关的阶段。由于历史原因，Go 语言的汇编器基于了不太常见的 plan9 汇编器的输入形式。需要注意的是，输入汇编器中的汇编指令不是机器码的表现形式，其仍然是人类可读的底层抽象。在 Go 语言 runtime 及 math/big 标准库中，可以看到许多特定指令集的汇编代码，Go 语言也提供了一些方式用于查看编译器生成的汇编代码。

```go
package main

import "fmt"

func main() {
	fmt.Println("Golang is awesome!")
}
```

对于上面的简单程序，其输出的汇编代码如下所示，这段汇编代码显示了 main 函数栈帧的大小与代码的行号及其对应的汇编指令。

```
go tool compile -S main.go
```

## 链接

程序可能使用其他程序或程序库（library），正如我们在 helloworld 程序中使用的 fmt package 一样，编写的程序必须与这些程序或程序库组合在一起才能执行，链接就是将编写的程序与外部程序组合在一起的过程。

链接分为静态链接与动态链接，静态链接的特点是链接器会将程序中使用的所有库程序复制到最后的可执行文件中，而动态链接只会在最后的可执行文件中存储动态链接库的位置，并在运行时调用。因此静态链接更快，并且可移植，它不需要运行它的系统上存在该库，但是它会占用更多的磁盘和内存空间。静态链接发生在编译时的最后一步，动态链接发生在程序加载到内存时。

表 1-1 对比了静态链接与动态链接的区别。

| | 静态链接 | 动态链接 |
| -------- | -------- | -------- |
| 所在位置 | 将程序中使用的所有库模块复制到最终可执行文件的过程。加载程序后，操作系统会将包含可执行代码和数据的单个文件放入内存 | 外部库（共享库）的地址放置在最终的可执行文件中，而实际链接是在运行时将可执行文件和库都放置在内存中时进行的，动态链接让多个程序可以使用可执行模块的单个副本 |
| 发生时期 | 由被称为链接器的程序执行，是编译程序的最后一步 | 由操作系统在运行时执行 |
| 文件大小 | 由于外部程序内置在可执行文件中，文件明显更大 | 共享库中只有一个副本保留在内存中，减小了可执行程序的大小，从而节省了内存和磁盘空间 |
| 扩展性 | 如果任何外部程序已更改，则必须重新编译并重新链接它们，否则更改将不会反映在现有的可执行文件中 | 只需要更新和重新编译各个共享模块程序即可 |
| 加载时间 | 程序每次将其加载到内存中执行时，都会花费恒定的加载时间 | 如果共享库代码已存在于内存中，则可以减少加载时间 |
| 程序运行时期 | 使用静态链接库的程序通常比使用共享库的程序快 | 使用共享库的程序通常比使用静态链接库的程序慢 |
| 兼容性 | 所有代码都包含在一个可执行模块中，不会遇到兼容性问题 | 程序依赖兼容的库，如果更改了库（例如，新的编译器版本更改了库），则必须重新设计应用程序以使其与新库兼容 |

Go 语言在默认情况下是使用静态链接的，但是在一些特殊情况下，如在使用了 CGO（即引用了 C 代码）时，则会使用操作系统的动态链接库，例如，Go 语言的 net/http 包在默认情况下会使用 libpthread 与 libc 的动态链接库。Go 语言也支持在 go build 编译时通过传递参数来指定要生成的链接库的方式，可以使用 go help build 命令查看。

```
go help buildmode
```

下面我们以 helloworld 程序为例，说明 Go 语言编译与链接的过程，我们可以使用 go build 命令，-x 参数代表打印执行的过程。

```
go build -x main.go
```

由于生成的信息较长，接下来，将逐步对输出的信息进行解析。

首先创建一个临时目录，用于存放临时文件。在默认情况下，命令结束时自动删除此目录，如果需要保留则添加 -work 参数。

```
WORK=C:\Users\Admin\AppData\Local\Temp\go-build1058200208
mkdir -p $WORK\b001\
cat >$WORK\b001\importcfg << 'EOF' # internal
```

然后生成编译配置文件，主要为编译过程需要的外部依赖（如引用的其他包的函数定义）。

```
# import config
packagefile fmt=C:\Program Files\Go\pkg\windows_amd64\fmt.a
packagefile runtime=C:\Program Files\Go\pkg\windows_amd64\runtime.a
```

编译阶段会生成中间结果 $WORK/b001/_pkg_.a。

```
cd D:\OneDrive\Repositories\projects\website\content\post\go\src\docs\internal\compiler\generating_machine_code
"C:\\Program Files\\Go\\pkg\\tool\\windows_amd64\\compile.exe" -o "$WORK\\b001\\_pk
g_.a" -trimpath "$WORK\\b001=>" -p main -complete -buildid diBfjiQ_4GLigqgUDk79/diB
fjiQ_4GLigqgUDk79 -goversion go1.19.1 -c=4 -nolocalimports -importcfg "$WORK\\b001\
\importcfg" -pack "D:\\OneDrive\\Repositories\\projects\\website\\content\\post\\go\\src\\docs\\internal\\compiler\\generating_machine_code\\main.go"
"C:\\Program Files\\Go\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\_pkg_.a" # internal
cp "$WORK\\b001\\_pkg_.a" "C:\\Users\\Admin\\AppData\\Local\\go-build\\0b\\0b7034a286b3d5f6cbd36258777d935581d71eeb23d1c8c2175146f69f0c6f25-d" # internal
```

.a 类型的文件又叫目标文件（object file），是一个压缩包，其内部包含 `_.PKGDEF` 和 `go_.o` 两个文件，分别为编译目标文件和链接目标文件。

文件内容由导出的函数、变量及引用的其他包的信息组成。弄清这两个文件包含的信息需要查看 Go 语言编译器实现的相关代码，这些代码在 `gc/obj.go` 文件中。

在下面的代码中，dumpobj1 函数会生成 ar 文件，ar 文件是一种非常简单的打包文件格式，广泛用于 Linux 的静态链接库文件中，文件以字符串“! \n”开头，随后是 60 字节的文件头部（包含文件名、修改时间等信息），之后是文件内容。

因为 ar 文件格式简单，所以 Go 语言编译器直接在函数中实现了 ar 打包过程。其中，startArchiveEntry 用于预留 ar 文件头信息的位置（60 字节），finishArchiveEntry 用于写入文件头信息。因为文件头信息中包含文件大小，在写入完成之前文件大小未知，所以分两步完成。

```go
func dumpobj1(outfile string, mode int) {
	bout, err := bio.Create(outfile)
	if err != nil {
		base.FlushErrors()
		fmt.Printf("can't create %s: %v\n", outfile, err)
		base.ErrorExit()
	}
	defer bout.Close()
	bout.WriteString("!<arch>\n")

	if mode&modeCompilerObj != 0 {
		start := startArchiveEntry(bout)
		dumpCompilerObj(bout)
		finishArchiveEntry(bout, start, "__.PKGDEF")
	}
	if mode&modeLinkerObj != 0 {
		start := startArchiveEntry(bout)
		dumpLinkerObj(bout)
		finishArchiveEntry(bout, start, "_go_.o")
	}
}
```

生成链接配置文件，主要包含了需要链接的依赖项。

```
cat >$WORK\b001\importcfg.link << 'EOF' # internal
packagefile command-line-arguments=$WORK\b001\_pkg_.a
packagefile fmt=C:\Program Files\Go\pkg\windows_amd64\fmt.a
packagefile runtime=C:\Program Files\Go\pkg\windows_amd64\runtime.a
packagefile errors=C:\Program Files\Go\pkg\windows_amd64\errors.a
...
EOF
```

执行链接器，生成最终可执行文件 main，同时将可执行文件复制到当前路径下并删除临时文件。

```
mkdir -p $WORK\b001\exe\
cd .
"C:\\Program Files\\Go\\pkg\\tool\\windows_amd64\\link.exe" -o "$WORK\\b001\\exe\\a
.out.exe" -importcfg "$WORK\\b001\\importcfg.link" -buildmode=pie -buildid=zWcOjaKB
vpnLVTT5uvcF/diBfjiQ_4GLigqgUDk79/CQod4AOywYXfEcXd2uOU/zWcOjaKBvpnLVTT5uvcF -extld=gcc "$WORK\\b001\\_pkg_.a"
"C:\\Program Files\\Go\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\exe\\a.out.exe" # internal
cp $WORK\b001\exe\a.out.exe main.exe
rm -r $WORK\b001\
```

## ELF 文件解析

在 Windows 操作系统下，编译后的 Go 文本文件最终会生成以 .exe 为后缀的 PE 格式的可执行文件，而在 Linux 和类 UNIX 操作系统下，会生成 ELF 格式的可执行文件。

除机器码外，在可执行文件中还可能包含调试信息、动态链接库信息、符号表信息。ELF（Executable and Linkable Format）是类 UNIX 操作系统下最常见的可执行且可链接的文件格式。有许多工具可以完成对 ELF 文件的探索查看，如 readelf、objdump。

下面使用 readelf 查看 ELF 文件的头信息：

```
readelf -h main
```

```
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              EXEC (Executable file)
  Machine:                           Advanced Micro Devices X86-64
  Version:                           0x1
  Entry point address:               0x45bfa0
  Start of program headers:          64 (bytes into file)
  Start of section headers:          456 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           56 (bytes)
  Number of program headers:         7
  Size of section headers:           64 (bytes)
  Number of section headers:         23
  Section header string table index: 3
```

ELF 包含多个 segment 与 section。debug/elf 包中给出了一些调试 ELF 的 API，以下程序可以打印出 ELF 文件中 section 的信息。

```go
package main

import "debug/elf"

func main() {
	info, err := elf.Open("main")
	if err != nil {
		panic(err)
	}
	defer info.Close()
	for _, section := range info.Sections {
		println(section.Name)
	}
}
```

通过 readelf 工具查看 ELF 文件中 section 的信息。

```
readelf -S main
```

```
There are 23 section headers, starting at offset 0x1c8:

Section Headers:
  [Nr] Name              Type             Address           Offset
       Size              EntSize          Flags  Link  Info  Align
  [ 0]                   NULL             0000000000000000  00000000
       0000000000000000  0000000000000000           0     0     0
  [ 1] .text             PROGBITS         0000000000401000  00001000
       000000000007cfa7  0000000000000000  AX       0     0     32
  [ 2] .rodata           PROGBITS         000000000047e000  0007e000
       0000000000035104  0000000000000000   A       0     0     32
  [ 3] .shstrtab         STRTAB           0000000000000000  000b3120
       000000000000017a  0000000000000000           0     0     1
  [ 4] .typelink         PROGBITS         00000000004b32a0  000b32a0
       00000000000004c0  0000000000000000   A       0     0     32
  [ 5] .itablink         PROGBITS         00000000004b3760  000b3760
       0000000000000058  0000000000000000   A       0     0     32
  [ 6] .gosymtab         PROGBITS         00000000004b37b8  000b37b8
       0000000000000000  0000000000000000   A       0     0     1
  [ 7] .gopclntab        PROGBITS         00000000004b37c0  000b37c0
       00000000000539b8  0000000000000000   A       0     0     32
  [ 8] .go.buildinfo     PROGBITS         0000000000508000  00108000
       0000000000000110  0000000000000000  WA       0     0     16
  [ 9] .noptrdata        PROGBITS         0000000000508120  00108120
       00000000000105e0  0000000000000000  WA       0     0     32
  [10] .data             PROGBITS         0000000000518700  00118700
       0000000000007810  0000000000000000  WA       0     0     32
  [11] .bss              NOBITS           000000000051ff20  0011ff20
       000000000002ef60  0000000000000000  WA       0     0     32
  [12] .noptrbss         NOBITS           000000000054ee80  0014ee80
       00000000000051a0  0000000000000000  WA       0     0     32
  [13] .zdebug_abbrev    PROGBITS         0000000000555000  00120000
       0000000000000127  0000000000000000           0     0     1
  [14] .zdebug_line      PROGBITS         0000000000555127  00120127
       000000000001b1b3  0000000000000000           0     0     1
  [15] .zdebug_frame     PROGBITS         00000000005702da  0013b2da
       00000000000054b3  0000000000000000           0     0     1
  [16] .debug_gdb_s[...] PROGBITS         000000000057578d  0014078d
       000000000000002d  0000000000000000           0     0     1
  [17] .zdebug_info      PROGBITS         00000000005757ba  001407ba
       000000000003353a  0000000000000000           0     0     1
  [18] .zdebug_loc       PROGBITS         00000000005a8cf4  00173cf4
       000000000001a1f3  0000000000000000           0     0     1
  [19] .zdebug_ranges    PROGBITS         00000000005c2ee7  0018dee7
       0000000000008954  0000000000000000           0     0     1
  [20] .note.go.buildid  NOTE             0000000000400f9c  00000f9c
       0000000000000064  0000000000000000   A       0     0     4
  [21] .symtab           SYMTAB           0000000000000000  00196840
       000000000000c090  0000000000000018          22    85     8
  [22] .strtab           STRTAB           0000000000000000  001a28d0
       000000000000ac51  0000000000000000           0     0     1
Key to Flags:
  W (write), A (alloc), X (execute), M (merge), S (strings), I (info),
  L (link order), O (extra OS processing required), G (group), T (TLS),
  C (compressed), x (unknown), o (OS specific), E (exclude),
  D (mbind), l (large), p (processor specific)
```

segment 包含多个 section，它描述程序如何映射到内存中，如哪些 section 需要导入内存、采取只读模式还是读写模式、内存对齐大小等。以下是 section 与 segment 的对应关系。

```
readelf -lW main
```

```
Elf file type is EXEC (Executable file)
Entry point 0x45bfa0
There are 7 program headers, starting at offset 64

Program Headers:
  Type           Offset   VirtAddr           PhysAddr           FileSiz  MemSiz   Flg Align
  PHDR           0x000040 0x0000000000400040 0x0000000000400040 0x000188 0x000188 R   0x1000
  NOTE           0x000f9c 0x0000000000400f9c 0x0000000000400f9c 0x000064 0x000064 R   0x4
  LOAD           0x000000 0x0000000000400000 0x0000000000400000 0x07dfa7 0x07dfa7 R E 0x1000
  LOAD           0x07e000 0x000000000047e000 0x000000000047e000 0x089178 0x089178 R   0x1000
  LOAD           0x108000 0x0000000000508000 0x0000000000508000 0x017f20 0x04c020 RW  0x1000
  GNU_STACK      0x000000 0x0000000000000000 0x0000000000000000 0x000000 0x000000 RW  0x8
  LOOS+0x5041580 0x000000 0x0000000000000000 0x0000000000000000 0x000000 0x000000     0x8

 Section to Segment mapping:
  Segment Sections...
   00
   01     .note.go.buildid
   02     .text .note.go.buildid
   03     .rodata .typelink .itablink .gosymtab .gopclntab
   04     .go.buildinfo .noptrdata .data .bss .noptrbss
   05
   06
```

并不是所有的 section 都需要导入内存，当 Type 为 LOAD 时，代表 section 需要被导入内存。后面的 Flg 代表内存的读写模式。包含 .text 的代码区代表可以被读和执行，包含 .data 与 .bss 的全局变量可以被读写，其中，为了满足垃圾回收的需要还区分了是否包含指针的区域。包含 .rodata 常量数据的区域代表只读区，其中，.itablink 为与 Go 语言接口相关的全局符号表。.gopclntab 包含程序计数器 PC 与源代码行的对应关系。

可以看到并不是所有 section 都需要导入内存，同时，该程序包含单独存储调试信息的区域。如.note.go.buildid 包含 Go 程序唯一的 ID，可以通过 objdump 工具在.note.go.buildid 中查找到每个 Go 程序唯一的 ID。

```
objdump -s -j .note.go.buildid main
```

```
main:     file format elf64-x86-64

Contents of section .note.go.buildid:
 400f9c 04000000 53000000 04000000 476f0000  ....S.......Go..
 400fac 4c306647 745a6557 62757271 596d6a67  L0fGtZeWburqYmjg
 400fbc 6572654e 2f476b52 6972622d 7947774e  ereN/GkRirb-yGwN
 400fcc 30565067 6a784861 372f797a 54307377  0VPgjxHa7/yzT0sw
 400fdc 48677639 7034444c 46686f53 56462f53  Hgv9p4DLFhoSVF/S
 400fec 34573051 5358554b 68324856 712d545a  4W0QSXUKh2HVq-TZ
 400ffc 2d6a6d00                             -jm.
```

另外，.go.buildinfo section 包含 Go 程序的构建信息，“go version” 命令会查找该区域的信息获取 Go 语言版本号。

```go

```
