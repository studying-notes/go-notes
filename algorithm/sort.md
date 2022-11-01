---
date: 2022-09-26T09:09:46+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "排序算法"  # 文章标题
url:  "posts/go/algorithm/sort"  # 设置网页永久链接
tags: [ "Go", "sort" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 选择排序

```go
// SelectionSort 选择排序
func SelectionSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		k := i // 记录最小值索引
		for j := i + 1; j < length; j++ {
			if array[j] < array[k] {
				k = j
			}
		}
		swap(array, i, k) // 交换最小值
	}
}
```

## 插入排序

```go
// InsertionSort 插入排序
func InsertionSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		for j := i + 1; j > 0; j-- {
			if array[j] < array[j-1] {
				break // 已经有序则中断
			}
			// 从后往前遍历有序序列，后移比它大的元素
			swap(array, j, j-1)
		}
	}
}
```

## 冒泡排序

```go
func BubbleSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		for j := 1; j < length-i; j++ {
			if array[j-1] > array[j] {
				swap(array, j, j-1)
			}
		}
	}
}
```

```go
// BubbleSort 冒泡排序
func BubbleSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		for j := length - 1; j > i; j-- {
			if array[j] < array[j-1] {
				swap(array, j, j-1)
			}
		}
	}
}
```

## 归并排序

```go
// MergeSort 归并排序
func MergeSort(array []int, begin, end int) {
	// begin 表示开始索引 end 表示结束索引，可以取到
	if begin < end {
		mid := (begin + end) / 2
		MergeSort(array, begin, mid)
		MergeSort(array, mid+1, end)
		merge(array, begin, mid, end)
	}
}

// 合并两个有序数组
func merge(array []int, begin, mid, end int) {
	leftLength := mid - begin + 1
	rightLength := end - mid

	left := make([]int, leftLength)
	right := make([]int, rightLength)

	for i := 0; i < leftLength; i++ {
		left[i] = array[begin+i]
	}

	for j := 0; j < rightLength; j++ {
		right[j] = array[mid+j+1]
	}

	var i, j int
	for i < leftLength && j < rightLength {
		if left[i] < right[j] {
			array[begin+i+j] = left[i]
			i++
		} else {
			array[begin+i+j] = right[j]
			j++
		}
	}

	for i < leftLength {
		array[begin+i+j] = left[i]
		i++
	}

	for j < rightLength {
		array[begin+i+j] = right[j]
		j++
	}
}
```

## 快速排序

```go
// Partition 用于快速排序中的分割
func Partition(array []int, begin, end int) int {
	pos, val := begin, array[begin]
	begin++
	for begin <= end {
		if array[begin] < val {
			array[begin], array[pos] = array[pos], array[begin]
			begin++
			pos++
		} else {
			array[end], array[begin] = array[begin], array[end]
			end--
		}
	}
	array[pos] = val
	return pos
}

// QuickSort 快速排序
func QuickSort(array []int, begin, end int) {
	if begin < end {
		mid := Partition(array, begin, end)
		QuickSort(array, begin, mid-1)
		QuickSort(array, mid+1, end)
	}
}
```

上面的分割算法位置交换太频繁了。

```go
func Partition(array []int, begin, end int) int {
	pivot := array[begin]
	for begin < end {
		for begin < end && array[begin] < pivot {
			begin++
		}
		for begin < end && array[end] > pivot {
			end--
		}
		swap(array, begin, end)
	}
	array[begin] = pivot
	return begin
}
```

## 希尔排序

```go
// HillSort 希尔排序
func HillSort(array []int, length int) {
	for step := length / 2; step > 0; step /= 2 {
		for i := step; i < length; i += step {
			for j := i; j > step-1; j -= step {
				if array[j] >= array[j-step] {
					break
				}
				swap(array, j, j-step)
			}
		}
	}
}
```
