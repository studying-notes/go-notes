---
date: 2022-10-02T20:41:09+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "字符串转换数组"  # 文章标题
url:  "posts/go/docs/internal/string/convert"  # 设置网页永久链接
tags: [ "Go", "convert" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [字符串和 `[]rune` 类型的相互转换](#字符串和-rune-类型的相互转换)
- [字节数组与字符串相互转换](#字节数组与字符串相互转换)
  - [string 转为 []byte](#string-转为-byte)
  - [[N]byte 转为 string](#nbyte-转为-string)
- [字符串内置操作模拟](#字符串内置操作模拟)
  - [`for range` 遍历模拟](#for-range-遍历模拟)
  - [`[]byte(s)` 转换模拟](#bytes-转换模拟)
  - [`string(bytes)` 转换模拟](#stringbytes-转换模拟)
  - [`[]rune(s)` 转换模拟](#runes-转换模拟)
  - [`string(runes)` 转换模拟](#stringrunes-转换模拟)

## 字符串和 `[]rune` 类型的相互转换

Go 语言除了 `for range` 语法对 UTF-8 字符串提供了特殊支持外，还对字符串和 `[]rune` 类型的相互转换提供了特殊的支持。

```go
func main() {
	fmt.Printf("%#v\n", []rune("世界"))             // []int32{19990, 30028}
	fmt.Printf("%#v\n", string([]rune{'世', '界'})) // 世界
}
```

从上面代码的输出结果可以发现 `[]rune` 其实是 `[]int32` 类型，这里的 `rune` 只是 `int32` 类型的别名，并不是重新定义的类型。`rune` 用于表示每个 Unicode 码点，目前只使用了 21 个位。

字符串相关的强制类型转换主要涉及 `[]byte` 和 `[]rune` 两种类型。每个转换都可能隐含重新分配内存的代价，最坏的情况下它们运算的时间复杂度都是 `O(n)`。

不过字符串和 `[]rune` 的转换要更为特殊一些，因为一般这种强制类型转换要求两个类型的底层内存结构要尽量一致，显然它们底层对应的 `[]byte` 和 `[]int32` 类型是完全不同的内存结构，因此这种转换可能隐含重新分配内存的操作。

## 字节数组与字符串相互转换

string 不能直接和 byte 数组转换 string 可以和 byte 的切片转换。

### string 转为 []byte

```go
[]byte("string")
```

字节数组转换为字符串**在运行时**调用了 slicebytetostring 函数。需要注意的是，字节数组与字符串的相互转换并不是简单的指针引用，而是涉及了复制。当字符串大于 32 字节时，还需要申请堆内存，因此在涉及一些密集的转换场景时，需要评估这种转换带来的性能损耗。

但**在编译时**，如果是临时场景，编译器会将字符串转换为字节数组的操作优化为直接引用字符串的底层字节数组，而不是复制。

当字符串转换为字节数组时，**在运行时**需要调用 stringtoslicebyte 函数，其和 slicebytetostring 函数非常类似，需要新的足够大小的内存空间。当字符串小于 32 字节时，可以直接使用缓存 buf。当字符串大于 32 字节时，rawbyteslice 函数需要向堆区申请足够的内存空间。最后使用 copy 函数完成内存复制。

### [N]byte 转为 string

```go
var buf [10]byte
string(buf[:])
```

## 字符串内置操作模拟

简单模拟 Go 语言对字符串内置的一些操作，这样对每个操作的处理的时间复杂度和空间复杂度都会有较明确的认识。

### `for range` 遍历模拟

`for range` 对字符串的迭代模拟实现如下：

```go
func forRangeString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}

func main() {
	s := "hello, 世界"
	forRangeString(s, func(i int, r rune) {
		fmt.Printf("%d %c\t", i, r)
	})
}
```

`for range` 迭代字符串时，每次解码一个 Unicode 字符，然后进入 `for` 循环体，遇到崩溃的编码并不会导致迭代停止。

### `[]byte(s)` 转换模拟

`[]byte(s)` 转换模拟实现如下：

```go
func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		p[i] = s[i]
	}
	return p
}
```

string 转换成 byte 切片，也需要一次内存拷贝，其过程如下：

* 申请切片内存空间
* 将 string 拷贝到切片

[![DiNtsS.png](../../../assets/images/docs/internal/string/convert/DiNtsS.png)](https://imgchr.com/i/DiNtsS)

模拟实现中新创建了一个切片，然后**将字符串的数组逐一复制到切片中**，这是为了**保证字符串只读的语义**。当然，在将字符串转换为 `[]byte` 时，如果转换后的变量没有被修改，编译器可能会直接返回原始的字符串对应的底层数据。

### `string(bytes)` 转换模拟

这种转换需要一次内存拷贝。

`string(bytes)` 转换模拟实现如下：

```go
func bytes2str(b []byte) (s string) {
    // 为了性能，有时候会省略复制
	data := make([]byte, len(b))
	for i, c := range b {
		data[i] = c
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(&data[0]))
	hdr.Len = len(b)
	return s
}
```

转换过程如下：

1. 根据切片的长度申请内存空间，假设内存地址为 p，切片长度为 len(b) ；
2. 构建 string （ string.str = p ； string.len = len ；）
3. 拷贝数据 ( 切片中数据拷贝到新申请的内存空间 )

[![DiN8RP.png](../../../assets/images/docs/internal/string/convert/DiN8RP.png)](https://imgchr.com/i/DiN8RP)

因为 Go 语言的字符串是只读的，无法以直接构造底层字节数组的方式生成字符串。在模拟实现中通过 `unsafe` 包获取字符串的底层数据结构，然后将切片的数据逐一复制到字符串中，这同样是为了保证字符串只读的语义不受切片的影响。如果转换后的字符串在生命周期中原始的 `[]byte` 的变量不发生变化，编译器可能会直接基于 `[]byte` 底层的数据构建字符串，而不进行复制。

### `[]rune(s)` 转换模拟

`[]rune(s)` 转换模拟实现如下：

```go
func str2runes(s string) []rune {
	// 将字符串转换为字节数组
	b := []byte(s)
	var p []int32
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		p = append(p, r)
		s = s[size:]
	}
	return p
}
```

因为底层内存结构的差异，所以字符串到 []rune 的转换必然会导致重新分配 []rune 内存空间，然后依次解码并复制对应的 Unicode 码点值。这种强制转换并不存在前面提到的字符串和字节切片转换时的优化情况。

### `string(runes)` 转换模拟

`string(runes)` 转换模拟实现如下：

```go
func runes2string(rs []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range rs {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}
```

同样因为底层内存结构的差异，**[]rune 到字符串的转换也必然会导致重新构造字符串**。这种强制转换并不存在前面提到的优化情况。

```go

```
