---
date: 2020-11-15T12:57:09+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go defer 语句理解"  # 文章标题
url:  "posts/go/abc/defer"  # 设置网页永久链接
tags: [ "go", "error" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 前言

defer 语句用于延迟函数的调用，每次 defer 都会把一个函数压入栈中，函数返回前再把延迟的函数取出并执行。

为了方便描述，我们把创建 defer 的函数称为主函数，defer 语句后面的函数称为延迟函数。

延迟函数可能有输入参数，这些参数可能来源于定义 defer 的函数，延迟函数也可能引用主函数用于返回的变量，也就是说延迟函数可能会影响主函数的一些行为，这些场景下，如果不了解 defer 的规则很容易出错。

## 热身

按照惯例，我们看几个有意思的题目，用于检验对 defer 的了解程度。

### 题目一

下面函数输出结果是什么？

```go
func deferFuncParameter() {
    var aInt = 1

    defer fmt.Println(aInt)

    aInt = 2
    return
}
```

题目说明：  

函数 deferFuncParameter() 定义一个整型变量并初始化为 1，然后使用 defer 语句打印出变量值，最后修改变量值为 2.

参考答案：  

输出 1。延迟函数 fmt.Println(aInt) 的参数在 defer 语句出现时就已经确定了，所以无论后面如何修改 aInt 变量都不会影响延迟函数。

### 题目二

下面程序输出什么？

```go
package main

import "fmt"

func printArray(array *[3]int) {
    for i := range array {
        fmt.Println(array[i])
    }
}

func deferFuncParameter() {
    var aArray = [3]int{1, 2, 3}

    defer printArray(&aArray)

    aArray[0] = 10
    return
}

func main() {
    deferFuncParameter()
}
```

函数说明：  

函数 deferFuncParameter() 定义一个数组，通过 defer 延迟函数 printArray() 的调用，最后修改数组第一个元素。printArray() 函数接受数组的指针并把数组全部打印出来。

参考答案：

输出 10、2、3 三个值。延迟函数 printArray() 的参数在 defer 语句出现时就已经确定了，即数组的地址，由于延迟函数执行时机是在 return 语句之前，所以对数组的最终修改值会被打印出来。

### 题目三

下面函数输出什么？

```go
func deferFuncReturn() (result int) {
    i := 1

    defer func() {
       result++
    }()

    return i
}
```

函数说明：

函数拥有一个具名返回值 result，函数内部声明一个变量 i，defer 指定一个延迟函数，最后返回变量 i。延迟函数中递增 result。

参考答案：

函数输出 2。函数的 return 语句并不是原子的，实际执行分为设置返回值 -->ret，defer 语句实际执行在返回前，即拥有 defer 的函数返回过程是：设置返回值 --> 执行 defer-->ret。所以 return 语句先把 result 设置为 i 的值，即 1，defer 语句中又把 result 递增 1，所以最终返回 2。

## defer 规则

Golang 官方博客里总结了 defer 的行为规则，只有三条，我们围绕这三条进行说明。

### 规则一：延迟函数的参数在 defer 语句出现时就已经确定下来了

官方给出一个例子，如下所示：

```go
func a() {
    i := 0
    defer fmt.Println(i)
    i++
    return
}
```

defer 语句中的 fmt.Println() 参数 i 值在 defer 出现时就已经确定下来，实际上是拷贝了一份。后面对变量 i 的修改不会影响 fmt.Println() 函数的执行，仍然打印" 0 "。

对于指针类型参数，规则仍然适用，只不过延迟函数的参数是一个地址值，这种情况下，defer 后面的语句对变量的修改可能会影响延迟函数。

### 规则二：延迟函数执行按后进先出顺序执行，即先出现的 defer 最后执行

这个规则很好理解，定义 defer 类似于入栈操作，执行 defer 类似于出栈操作。

设计 defer 的初衷是简化函数返回时资源清理的动作，资源往往有依赖顺序，比如先申请 A 资源，再根据 A 资源申请 B 资源，根据 B 资源申请 C 资源，即申请顺序是: A-->B-->C，释放时往往又要反向进行。这就是把 defer 设计成 LIFO 的原因。

每申请到一个用完需要释放的资源时，立即定义一个 defer 来释放资源是个很好的习惯。

### 规则三：延迟函数可能操作主函数的具名返回值

定义 defer 的函数，即主函数可能有返回值，返回值有没有名字没有关系，defer 所作用的函数，即延迟函数可能会影响到返回值。

若要理解延迟函数是如何影响主函数返回值的，只要明白函数是如何返回的就足够了。

#### 函数返回过程

有一个事实必须要了解，关键字 *return* 不是一个原子操作，实际上 *return* 只代理汇编指令 *ret*，即将跳转程序执行。比如语句 `return i`，实际上分两步进行，即将 i 值存入栈中作为返回值，然后执行跳转，而 defer 的执行时机正是跳转前，所以说 defer 执行时还是有机会操作返回值的。

举个实际的例子进行说明这个过程：

```go
func deferFuncReturn() (result int) {
    i := 1

    defer func() {
       result++
    }()

    return i
}
```

该函数的 return 语句可以拆分成下面两行：

```go
result = i
return
```

而延迟函数的执行正是在 return 之前，即加入 defer 后的执行过程如下：

```go
result = i
result++
return
```

所以上面函数实际返回 i++ 值。

关于主函数有不同的返回方式，但返回机制就如上机介绍所说，只要把 return 语句拆开都可以很好的理解，下面分别举例说明

#### 主函数拥有匿名返回值，返回字面值

一个主函数拥有一个匿名的返回值，返回时使用字面值，比如返回"1"、"2"、"Hello"这样的值，这种情况下 defer 语句是无法操作返回值的。

一个返回字面值的函数，如下所示：

```go
func foo() int {
    var i int

    defer func() {
        i++
    }()

    return 1
}
```

上面的 return 语句，直接把 1 写入栈中作为返回值，延迟函数无法操作该返回值，所以就无法影响返回值。

#### 主函数拥有匿名返回值，返回变量

一个主函数拥有一个匿名的返回值，返回使用本地或全局变量，这种情况下 defer 语句可以引用到返回值，但不会改变返回值。

一个返回本地变量的函数，如下所示：

```go
func foo() int {
    var i int

    defer func() {
        i++
    }()

    return i
}
```

上面的函数，返回一个局部变量，同时 defer 函数也会操作这个局部变量。对于匿名返回值来说，可以假定仍然有一个变量存储返回值，假定返回值变量为"anony"，上面的返回语句可以拆分成以下过程：

```go
anony = i
i++
return
```

由于 i 是整型，会将值拷贝给 anony，所以 defer 语句中修改 i 值，对函数返回值不造成影响。

#### 主函数拥有具名返回值

主函声明语句中带名字的返回值，会被初始化成一个局部变量，函数内部可以像使用局部变量一样使用该返回值。如果 defer 语句操作该返回值，可能会改变返回结果。

一个影响函返回值的例子：

```go
func foo() (ret int) {
    defer func() {
        ret++
    }()

    return 0
}
```

上面的函数拆解出来，如下所示：

```go
ret = 0
ret++
return
```

函数真正返回前，在 defer 中对返回值做了 +1 操作，所以函数最终返回 1。

## defer 实现原理

### defer 数据结构

源码包 `src/src/runtime/runtime2.go:_defer` 定义了 defer 的数据结构：

```go
type _defer struct {
    sp      uintptr   //函数栈指针
    pc      uintptr   //程序计数器
    fn      *funcval  //函数地址
    link    *_defer   //指向自身结构的指针，用于链接多个defer
}
```

我们知道 defer 后面一定要接一个函数的，所以 defer 的数据结构跟一般函数类似，也有栈地址、程序计数器、函数地址等等。

与函数不同的一点是它含有一个指针，可用于指向另一个 defer，每个 goroutine 数据结构中实际上也有一个 defer 指针，该指针指向一个 defer 的单链表，每次声明一个 defer 时就将 defer 插入到单链表表头，每次执行 defer 时就从单链表表头取出一个 defer 执行。

下图展示多个 defer 被链接的过程：

[![Di0JHO.png](https://s3.ax1x.com/2020/11/15/Di0JHO.png)](https://imgchr.com/i/Di0JHO)

从上图可以看到，新声明的 defer 总是添加到链表头部。函数返回前执行 defer 则是从链表首部依次取出执行。

一个 goroutine 可能连续调用多个函数，defer 添加过程跟上述流程一致，进入函数时添加 defer，离开函数时取出 defer，所以即便调用多个函数，也总是能保证 defer 是按 LIFO 方式执行的。

### defer 的创建和执行

源码包 `src/runtime/panic.go` 定义了两个方法分别用于创建 defer 和执行 defer。

- deferproc()： 在声明 defer 处调用，其将 defer 函数存入 goroutine 的链表中；
- deferreturn()：在 return 指令，准确的讲是在 ret 指令前调用，其将 defer 从 goroutine 链表中取出并执行。

可以简单这么理解，在编译阶段，声明 defer 处插入了函数 deferproc()，在函数 return 前插入了函数 deferreturn()。

## 小结

- defer 定义的延迟函数参数在 defer 语句出现时就已经确定下来了
- defer 定义顺序与实际执行顺序相反
- return 不是原子操作，执行过程是: 保存返回值 ( 若有 )--> 执行 defer （若有） --> 执行 ret 跳转
- 申请资源后立即使用 defer 关闭资源是好习惯
