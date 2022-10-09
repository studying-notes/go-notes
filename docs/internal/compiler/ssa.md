---
date: 2022-10-01T13:44:36+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "SSA 生成"  # 文章标题
url:  "posts/go/docs/internal/compiler/ssa"  # 设置网页永久链接
tags: [ "Go", "ssa" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

* `cmd/compile/internal/ssa` (SSA passes and rules)
* `cmd/compile/internal/ssagen` (converting IR to SSA)

在此阶段，IR 被转换为静态单一分配 (SSA) 形式，这是一种具有特定属性的较低级别的中间表示，可以更轻松地实现优化并最终从中生成机器代码。

在此转换期间，将应用函数内在函数。这些是编译器被教导的特殊功能，可以根据具体情况用高度优化的代码替换。

在 AST 到 SSA 的转换过程中，某些节点也被降低为更简单的组件，以便编译器的其余部分可以使用它们。例如，内置的 copy 被内存移动替换，range 循环被重写为 for 循环。由于历史原因，其中一些目前发生在转换为 SSA 之前，但长期计划是将它们全部移到这个阶段。

然后，应用一系列与机器无关的通行证和规则。这些不涉及任何单一的计算机架构，因此可以在所有“GOARCH”变体上运行。这些通行证包括消除死代码、删除不需要的 nil 检查和删除未使用的分支。通用重写规则主要关注表达式，比如用常量值替换一些表达式，优化乘法和浮点运算。

-------------------------------

在 SSA 生成阶段，每个变量在声明之前都需要被定义，并且，每个变量只会被赋值一次。

```go
y := 1
y := 2
x = y
```

例如，在上面的代码中，变量 y 被赋值了两次，不符合 SSA 的规则，很容易看出，y:=1 这条语句是无效的。可以转化为如下形式：

```go
y1 := 1
y2 := 2
x1 :=y2
```

通过 SSA，很容易识别出 y1 是无效的代码并将其清除。

条件判断等多个分支的情况会稍微复杂一些，如下所示，假如我们将第一个 x 变为 x_1，条件变量括号内的 x 变为 x_2，那么 f(x) 中的 x 应该是 x_1 还是 x_2 呢？

```go
x = 1
if condition {
    x =2
}
f(x)
```

为了解决以上问题，在 SSA 生成阶段需要引入额外的函数Φ接收 x_1 和 x_2 产生新的变量 x_v，x_v 的大小取决于代码运行的路径，如图 1-8 所示。

![](../../../assets/images/docs/internal/compiler/ssa/图1-8%20SSA生成阶段处理多分支下的单一变量名.png)

SSA 生成阶段是编译器进行后续优化的保证，例如常量传播（Constant Propagation）、无效代码清除、消除冗余、强度降低（Strength Reduction）等。

Go 语言提供了工具查看 SSA 初始及其后续优化阶段生成的代码片段，可以通过在编译时指定 GOSSAFUNC=main 实现。

```go
package main

var d uint8

func main() {
	var a uint8 = 1
	a = 2
	if true {
		a = 3
	}
	d = a
}
```

```
GOSSAFUNC=main GOOS=linux GOARCH=amd64 go tool compile main.go
```

通过浏览器打开 ssa.html 文件，将看到图 1-9 所示的许多代码片段，其中一些片段是隐藏的。这些是 SSA 的初始阶段、优化阶段、最终阶段的代码片段。

![](../../../assets/images/docs/internal/compiler/ssa/图1-9%20SSA所有优化阶段的代码片段.png)

以如下最初生成 SSA 代码的初始（start）阶段为例，其中，bN 代表不同的执行分支，例如 b1、b2、b3。vN 代表变量，每个变量只能被分配一次，变量后的 Op 操作代表不同的语义，与特定的机器无关。例如 Addr 代表取值操作，Const8 代表常量，后接要操作的类型； Store 代表赋值是与内存有关的操作。Go 语言编译器采取了特殊的方式处理内存操作，例如 v11 中 Store 的第三个参数代表内存的状态，用于确定内存的依赖关系，从而避免编译器内存的重排。另外，v8 的取值取决于判断语句是否为 true，这就是之前介绍的函数Φ。

```
start
number lines [520504 ns]
b1:-
v1 (?) = InitMem <mem>
v2 (?) = SP <uintptr>
v3 (?) = SB <uintptr>
v4 (?) = Const8 <uint8> [1] (a[uint8])
v5 (?) = Const8 <uint8> [2] (a[uint8])
v6 (?) = Const8 <uint8> [3] (a[uint8])
v7 (?) = Addr <*uint8> {"".d} v3
v8 (+11) = Store <mem> {uint8} v7 v6 v1
v9 (+12) = MakeResult <mem> v8
Ret v9 (12)
name a[uint8]: v4 v5 v6
```

初始阶段结束后，编译器将根据生成的 SSA 进行一系列重写和优化。SSA 最终的阶段叫作 genssa，在上例的 genssa 阶段中，编译器清除了无效的代码及不会进入的 if 分支，并且将常量 Op 操作变为了 amd64 下特定的 MOVBstoreconst 操作。

```go

```
