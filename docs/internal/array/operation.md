---
date: 2022-10-03T09:15:33+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数组基本操作"  # 文章标题
url:  "posts/go/docs/internal/array/operation"  # 设置网页永久链接
tags: [ "Go", "operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [创建](#创建)
- [复制](#复制)
- [比较](#比较)
- [扩容](#扩容)
- [遍历](#遍历)
- [排序](#排序)
- [查找](#查找)
- [删除](#删除)
- [插入](#插入)
- [追加](#追加)

## 创建

```go
func make(t *_type, len int) unsafe.Pointer {
	// 1. 计算数组的大小
	size := int(t.size) * len
	// 2. 分配内存
	p := mallocgc(size, t, true)
	// 3. 初始化数组
	memclrNoHeapPointers(p, uintptr(size))
	return p
}
```

```go
func newobject(t *_type) unsafe.Pointer {
	// 1. 分配内存
	p := mallocgc(t.size, t, true)
	// 2. 初始化数组
	memclrNoHeapPointers(p, t.size)
	return p
}
```

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

## 扩容

数组的扩容是通过创建新的数组来实现的，然后将旧数组的元素复制到新数组中。

```go
func growslice(et *_type, old slice, cap int) slice {
	// 1. 计算新数组的长度
	newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if newcap < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
	// 2. 创建新数组
	newlen := old.len
	if cap < newlen {
		panic(errorString("growslice: cap out of range"))
	}
	newbase := mallocgc(uintptr(newcap)*et.size, et, true)
	// 3. 复制旧数组的元素到新数组
	memmove(newbase, old.array, uintptr(newlen)*et.size)
	// 4. 返回新数组
	return slice{newbase, newlen, newcap}
}
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

## 排序

```go
func sort(a slice, less func(i, j int) bool) {
	// 1. 创建排序器
	s := newSorter(a, less)
	// 2. 调用排序器的排序方法
	s.sort(0, a.len)
}
```

## 查找

```go
func search(a slice, x interface{}, less func(i, j int) bool) int {
	// 1. 创建查找器
	s := newSorter(a, less)
	// 2. 调用查找器的查找方法
	return s.search(x)
}
```

## 删除

```go
func deleteslice(a slice, i int) slice {
	// 1. 计算删除元素后的数组长度
	n := a.len - 1
	// 2. 计算删除元素后的数组容量
	c := a.cap
	// 3. 计算删除元素后的数组地址
	p := a.array
	// 4. 计算删除元素的地址
	q := add(p, uintptr(i)*a.elemtype.size)
	// 5. 计算删除元素后的数组长度
	r := add(q, a.elemtype.size)
	// 6. 计算删除元素后的数组长度
	s := add(p, uintptr(n)*a.elemtype.size)
	// 7. 计算删除元素后的数组长度
	t := add(p, uintptr(c)*a.elemtype.size)
	// 8. 计算删除元素后的数组长度
	memmove(q, r, uintptr(n-i)*a.elemtype.size)
	// 9. 计算删除元素后的数组长度
	memclrNoHeapPointers(s, uintptr(t-s))
	// 10. 返回删除元素后的数组
	return slice{p, n, c}
}
```

## 插入

```go
func insertslice(a slice, i int, x interface{}) slice {
	// 1. 计算插入元素后的数组长度
	n := a.len + 1
	// 2. 计算插入元素后的数组容量
	c := a.cap
	// 3. 计算插入元素后的数组地址
	p := a.array
	// 4. 计算插入元素的地址
	q := add(p, uintptr(i)*a.elemtype.size)
	// 5. 计算插入元素后的数组长度
	r := add(q, a.elemtype.size)
	// 6. 计算插入元素后的数组长度
	s := add(p, uintptr(n)*a.elemtype.size)
	// 7. 计算插入元素后的数组长度
	t := add(p, uintptr(c)*a.elemtype.size)
	// 8. 计算插入元素后的数组长度
	memmove(r, q, uintptr(n-i)*a.elemtype.size)
	// 9. 计算插入元素后的数组长度
	memclrNoHeapPointers(q, uintptr(t-s))
	// 10. 返回插入元素后的数组
	return slice{p, n, c}
}
```

## 追加

```go
func appendslice(t slice, x slice) slice {
	// 1. 计算追加元素后的数组长度
	n := t.len + x.len
	// 2. 计算追加元素后的数组容量
	c := t.cap
	// 3. 计算追加元素后的数组地址
	p := t.array
	// 4. 计算追加元素的地址
	q := add(p, uintptr(t.len)*t.elemtype.size)
	// 5. 计算追加元素后的数组长度
	r := add(p, uintptr(n)*t.elemtype.size)
	// 6. 计算追加元素后的数组长度
	s := add(p, uintptr(c)*t.elemtype.size)
	// 7. 计算追加元素后的数组长度
	memmove(q, x.array, uintptr(x.len)*t.elemtype.size)
	// 8. 计算追加元素后的数组长度
	memclrNoHeapPointers(r, uintptr(s-r))
	// 9. 返回追加元素后的数组
	return slice{p, n, c}
}
```

```go

```
