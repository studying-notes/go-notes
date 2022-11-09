---
date: 2022-11-06T10:11:38+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数组"  # 文章标题
url:  "posts/go/algorithm/structures/linear_list/array"  # 设置网页永久链接
tags: [ "Go", "array" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 交换数组元素

```go
func swap(array []int, i, j int) {
	array[i], array[j] = array[j], array[i]
}
```

## 产生指定范围的有序整数数组

```go
// Range 可生成指定范围的自然数切片
func Range(args ...uint) []uint {
	var (
		start, end uint
		step       uint = 1
	)

	if len(args) == 0 {
		return []uint{}
	} else if len(args) == 1 {
		end = args[0]
	} else if len(args) == 2 {
		start, end = args[0], args[1]
	} else if len(args) > 2 {
		start, end, step = args[0], args[1], args[2]
	}

	if step == 0 ||
		start == end ||
		(step < 0 && start < end) ||
		(step > 0 && start > end) {
		return []uint{}
	}

	s := make([]uint, 0, (end-start)/step+1)

	for start < end {
		s = append(s, start)
		start += step
	}

	return s
}
```

```go
func ExampleRange() {
	fmt.Println(Range(5))
	fmt.Println(Range(1, 5))
	fmt.Println(Range(1, 10, 2))

	// Output:
	// [0 1 2 3 4]
	// [1 2 3 4]
	// [1 3 5 7 9]
}
```

```go

```
