---
date: 2022-10-13T15:43:41+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "了解 CPU 寄存器"  # 文章标题
url:  "posts/go/docs/internal/base/register"  # 设置网页永久链接
tags: [ "Go", "register" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 简介

寄存器是 CPU 内部的存储单元，用于存放从内存读取而来的数据（包括指令）和 CPU 运算的中间结果，之所以要使用寄存器来临时存放数据而不是直接操作内存，一是因为 CPU 的工作原理决定了**有些操作运算只能在 CPU 内部进行**，二是因为 **CPU 读写寄存器的速度比读写内存的速度快得多**。

为了便于交流和使用汇编语言进行编程，CPU 厂商为每个寄存器都取了一个名字，比如 AMD64 CPU 中的 rax, rbx, rcx, rdx 等等，这样就可以很方便的在汇编代码中使用寄存器的名字来进行编程。

假设有如下 go 语言编写的代码：

```go
package main

func main() {
	var a, b, c int
	c = a + b
	println(c)
}
```

在 Windows x64 平台下，使用 Go 编译器编译它。

```bash
go tool compile -N -l -S main.go
```

提取后可得到如下汇编代码：

```
0x0014 00020 (main.go:5)	MOVQ	$1, main.a+24(SP)
0x001d 00029 (main.go:6)	MOVQ	$2, main.b+16(SP)
0x0026 00038 (main.go:7)	MOVQ	$0, main.c+8(SP)
0x002f 00047 (main.go:10)	MOVQ	main.a+24(SP), AX
0x0034 00052 (main.go:10)	ADDQ	main.b+16(SP), AX
0x0039 00057 (main.go:10)	MOVQ	AX, main.c+8(SP)
```

可以看到，上面的代码被编译成了多条汇编指令，指令中出现的 AX 和 SP 都是寄存器的名字，汇编代码所做的工作就是把数据在内存和寄存器中搬来搬去或做一些基础的数学和逻辑运算。

## 寄存器分类

不同体系结构的 CPU，其内部寄存器的数量、种类以及名称可能大不相同，目前使用最为广泛的 AMD64 这种体系结构的 CPU，这种 CPU 共有 20 多个可以直接在汇编代码中使用的寄存器，其中有几个寄存器在操作系统代码中才会见到，而应用层代码一般只会用到如下分为三类的 19 个寄存器。

1. **通用寄存器**：rax, rbx, rcx, rdx, rsi, rdi, rbp, rsp, r8, r9, r10, r11, r12, r13, r14, r15 寄存器。CPU 对这 16 个通用寄存器的用途没有做特殊规定，可以自定义其用途；

2. **程序计数寄存器（PC 寄存器）**：rip 寄存器。它用来存放下一条即将执行的指令的地址，这个寄存器决定了程序的执行流程；

3. **段寄存器**：fs 和 gs 寄存器。一般用它来实现**线程本地存储**（TLS），比如 AMD64 Linux 平台下 go 语言和 pthread 都使用 fs 寄存器来实现系统线程的 TLS（线程本地存储）。

上述这些寄存器除了 fs 和 gs 段寄存器是 16 位的，其它都是 64 位的，也就是 8 个字节，其中的 16 个通用寄存器还可以作为 32/16/8 位寄存器使用，只是使用时需要换一个名字，比如可以用 eax 这个名字来表示一个 32 位的寄存器，它使用的是 rax 寄存器的低 32 位。

为了便于查阅，下表列出这些64通用寄存器对应的 32/16/8 位寄存器的名字：

| 64 位 | 32 位 | 16 位 | 8 位 |
| :----- | :------- | :------- | :------- |
| rax | eax | ax | al/ah |
| rbx | ebx | bx | bl/bh |
| rcx | ecx | cx | cl/ch |
| rdx | edx | dx | dl/dh |
| rsi | esi | si | - |
| rdi | edi | di | - |
| rbp | ebp | bp | - |
| rsp | esp | sp | - |
| r8 ~ r15 | r8d ~ r15d | r8w ~ r15w | r8b ~ r15b |

通用寄存器主要用于临时存放数据，但有三个比较特殊的寄存器值得在这里单独提出来做一下说明：

### rip 寄存器

rip 寄存器里面存放的是 CPU 即将执行的下一条指令在内存中的地址。

**rip 寄存器的值不是正在被 CPU 执行的指令在内存中的地址，而是紧挨这条正在被执行的指令后面那一条指令的地址。**

修改 rip 寄存器的值是 CPU 自动控制的，不需要我们用指令去修改，当然 CPU 也提供了几条可以间接修改 rip 寄存器的指令。

## rsp 栈顶寄存器和 rbp 栈基址寄存器

这两个寄存器都跟函数调用栈有关，其中 rsp 寄存器一般用来存放函数调用栈的栈顶地址，而 rbp 寄存器通常用来存放函数的栈帧起始地址，编译器一般使用这两个寄存器加一定偏移的方式来访问函数局部变量或函数参数。

```go

```
