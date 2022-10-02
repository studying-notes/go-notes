---
date: 2020-11-14T22:21:59+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 切片 append 方法的陷阱"  # 文章标题
url:  "posts/go/docs/grammar/slice/append"  # 设置网页永久链接
tags: [ "Go", "append" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 函数原型

```go
func append(slice []Type, elems ...Type) []Type
```

## 会改变切片的地址

`append` 的本质是向切片中追加数据，而随着切片中元素逐渐增加，当切片底层的数组将满时，切片会发生扩容，扩容会导致产生一个新的切片（拥有容量更大的底层数组）。

`append` 每个追加元素，都有可能触发切片扩容，也即有可能返回一个新的切片，这也是 `append` 函数声明中返回值为切片的原因。实际使用中应该总是接收该返回值。

不管初始切片长度为多少，不接收 `append` 返回都是有极大风险的。

```go
package main

import (
	"fmt"
)

func AddElement(slice []int, e int) []int {
	return append(slice, e)
}

func main() {
	var slice []int
	//fmt.Println("length of slice: ", len(slice))
	//fmt.Println("capacity of slice: ", cap(slice))

	slice = append(slice, 1, 2, 3)
	//fmt.Println("length of slice: ", len(slice))
	//fmt.Println("capacity of slice: ", cap(slice))

	newSlice := AddElement(slice, 4)
	//fmt.Println("length of slice: ", len(newSlice))
	//fmt.Println("capacity of slice: ", cap(newSlice))

	fmt.Println(&slice[0] == &newSlice[0])
}
```

函数 AddElement() 接受一个切片和一个元素，把元素 append 进切片中，并返回切片。main() 函数中定义一个切片，并向切片中 append 3 个元素，接着调用 AddElement() 继续向切片 append 进第 4 个元素同时定义一个新的切片 newSlice。最后判断新切片 newSlice 与旧切片 slice 是否共用一块存储空间。

append 函数执行时会判断切片容量是否能够存放新增元素，如果不能，则会重新申请存储空间，新存储空间将是原来的 2 倍或 1.25 倍（取决于扩展原空间大小），本例中实际执行了两次 append 操作，第一次空间增长到 4，所以第二次 append 不会再扩容，所以新旧两个切片将共用一块存储空间。程序会输出 "true"。

## 可以追加 nil 值

向切片中追加一个 `nil` 值是完全不会报错的，如下代码所示：

```
slice := append(slice, nil)
```

经过追加后，slice 的长度递增 1。

实际上 `nil` 是一个预定义的值，即空值，所以完全有理由向切片中追加。

```go

```
