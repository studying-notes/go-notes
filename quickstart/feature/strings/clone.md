---
date: 2022-03-19T08:43:59+08:00
author: "Rustle Karl"  # 作者

title: "字符串克隆 Clone"  # 文章标题
url:  "posts/go/quickstart/feature/strings/clone"  # 设置网页永久链接
tags: [ "Go", "clone" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 功能

Clone 返回字符串的新副本。它保证将字符串复制到一个新分配的副本中，当只保留一个很大的字符串中的一个小子字符串时，这一点很重要。

使用克隆可以帮助这些程序使用更少的内存。当然，由于使用克隆制作拷贝，过度使用克隆会使程序使用更多内存。通常，只有在需要克隆时，才谨慎使用克隆。

对于长度为零的字符串，将返回字符串 ""，不进行内存分配。

## 原因

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s := "abcdefghijklmn"
	sSlice := s[:4]

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sSliceHeader := (*reflect.StringHeader)(unsafe.Pointer(&sSlice))

	fmt.Println(sHeader.Len == sSliceHeader.Len)
	fmt.Println(sHeader.Data == sSliceHeader.Data)

	// Output:
	// false
	// true
}
```

在上面示例场景中，如果 s 很大，而之后我们只需要使用它的某个短子串，这会导致内存的浪费，因为子串和原字符串的 Data 部分指向相同的内存，因此整个字符串并不会被 GC 回收。

strings.Clone 函数就是为了解决这个问题的。

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	s := "abcdefghijklmn"
	sSlice := strings.Clone(s[:4])

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sSliceHeader := (*reflect.StringHeader)(unsafe.Pointer(&sSlice))

	fmt.Println(sHeader.Len == sSliceHeader.Len)
	fmt.Println(sHeader.Data == sSliceHeader.Data)

	// Output:
	// false
	// false
}
```

通过克隆得到 sSlice，从最后输出结果看，Data 已经不同了，原始的长字符串就可以被垃圾回收了。

### 实现

```go
func Clone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
    // 实现 []byte 到 string 的零内存拷贝转换。
	return *(*string)(unsafe.Pointer(&b))
}
```

```go

```
