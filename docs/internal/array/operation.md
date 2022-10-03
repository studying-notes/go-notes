---
date: 2022-10-03T09:15:33+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数组的基本操作"  # 文章标题
url:  "posts/go/docs/internal/array/operation"  # 设置网页永久链接
tags: [ "Go", "operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [复制](#复制)
- [比较](#比较)
- [遍历](#遍历)

## 复制

```go
func slicecopy(to slice, fr slice) int {
	// 1. 计算需要复制的元素个数
	n := fr.len
	if n > to.len {
		n = to.len
	}
	// 2. 复制元素
	memmove(to.array, fr.array, uintptr(n)*fr.elemtype.size)
	return n
}
```

```go
// 汇编实现
func memmove(to, from unsafe.Pointer, n uintptr)
```

将 from 指向的内存区域的前 n 个字节拷贝到 to 指向的内存区域中。

与 C 语言中的数组显著不同的是，Go 语言中的数组在赋值和函数调用时的形参都是值复制。

如下所示，无论是赋值的 b 还是函数调用中的 c，都是值复制的。这意味着不管是修改 b 还是 c 的值，都不会影响 a 的值，因为他们是完全不同的数组。每个数组在内存中的位置都是不相同的。

```go
package main

func change(c [3]int) {
	c[0] = 4
}

func main() {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 0
	change(a)
	println(a[0])
}
```

## 比较

```go
// 汇编实现
func memequal(x, y unsafe.Pointer, n uintptr) bool
```

## 遍历

遍历数组的 3 种方式：

```go
func main() {
	var a = [...]int{1, 2, 3}

	for i := range a {
		fmt.Printf("a[%d]: %d\n", i, a[i])
	}
	for i, v := range a {
		fmt.Printf("a[%d]: %d\n", i, v)
	}
	for i := 0; i < len(a); i++ {
		fmt.Printf("a[%d]: %d\n", i, a[i])
	}
}
```

用 `for range` 方式迭代的性能可能会更好一些，因为这种迭代可以**保证不会出现数组越界**的情形，每轮迭代对数组元素的访问时可以省去对下标越界的判断。

用 `for range` 方式迭代，还可以忽略迭代时的下标：

```go
func main() {
	var times [5][0]int
	for range times {
		fmt.Println("hello")
	}

	for range make([][]int, 5){
		fmt.Println("world")
	}
}
```

其中 times 对应一个 `[5][0]int` 类型的数组，虽然第一维数组有长度，但是数组的元素 `[0]int` 大小是 0，因此整个数组占用的内存大小依然是 0。**不用付出额外的内存代价**，我们就通过 `for range` 方式实现 `times` 次快速迭代。

```go

```
