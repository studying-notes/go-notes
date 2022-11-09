---
date: 2022-11-06T10:11:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数组相关问题"  # 文章标题
url:  "posts/go/algorithm/structures/linear_list/array_quiz"  # 设置网页永久链接
tags: [ "Go", "array-quiz" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 求解最小三元组距离

### 定义三元组距离为最大差的绝对值

```go
func abs(a int) int {
	if a > 0 {
		return a
	}
	return ^(a - 1) // ^a + 1
}

func max(a, b, c int) int {
	d := (a + b + abs(a-b)) / 2
	return (c + d + abs(c-d)) / 2
}

func distance(a, b, c int) int {
	return max(abs(a-b), abs(b-c), abs(c-a))
}

func findMinimalDistance(s1, s2, s3 []int, l1, l2, l3 int) int {
	var i1, i2, i3, dist int
	min := distance(s1[0], s2[0], s3[0])

	for i1 < l1 && i2 < l2 && i3 < l3 {
		dist = distance(s1[i1], s2[i2], s3[i3])
		if dist < min {
			min = dist
		}
		if s1[i1] < s2[i2] && s1[i1] < s3[i3] {
			i1++
		} else if s2[i2] < s1[i1] && s2[i2] < s3[i3] {
			i2++
		} else {
			i3++
		}
	}

	return min
}

func Example_findMinimalDistance() {
	s1 := []int{-1, 0, 9}
	s2 := []int{-25, -10, 10, 11}
	s3 := []int{2, 9, 17, 30, 41}

	fmt.Println(findMinimalDistance(s1, s2, s3, len(s1), len(s2), len(s3)))

	// Output:
	// 1
}
```

## 找出数组中唯一的重复元素

### 哈希表法

```go
func findTheUniqueRepeated(array []int, length int) (bool, int) {
	visited := make(map[int]bool)
	for i := 0; i < length; i++ {
		if visited[array[i]] {
			return true, array[i]
		}
		visited[array[i]] = true
	}
	return false, 0
}

func Example_findTheUniqueRepeated() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
	if ok, v := findTheUniqueRepeated(array, len(array)); ok {
		fmt.Println(v)
	}

	// Output:
	// 4
}
```

### 累计差值求和

```go
func findTheUniqueRepeated(array []int, length int) (bool, int) {
	candidate := 0
	for i := 0; i < length; i++ {
		candidate += array[i] - i
	}
	return candidate != 0, candidate
}
```

### 异或法

```go
func findTheUniqueRepeated(array []int, length int) (bool, int) {
	candidate := 0
	for i := 0; i < length; i++ {
		candidate += array[i] - i
	}
	return candidate != 0, candidate
}
```

## 找出数组中丢失的数

```go
func findTheMissing(array []int, length int) int {
	candidate := length + 1
	for i := 0; i < length; i++ {
		candidate += i + 1 - array[i]
	}
	return candidate
}

func Example_findTheMissing() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(findTheMissing(array, len(array)))

	// Output:
	// 10
}
```

## 找出数组中出现奇数次的两个元素

```go
func findOddElements(array []int, length int) (x int, y int) {
	for i := 0; i < length; i++ {
		x ^= array[i]
	}

	y = x
	bit := x &^ (x - 1)

	for i := 0; i < length; i++ {
		element := array[i]
		if bit&element != 0 {
			y ^= element
		}
	}

	return x ^ y, y
}

func Example_findOddElements() {
	array := []int{1, 2, 3, 4, 5, 3, 2, 1}
	fmt.Println(findOddElements(array, len(array)))

	// Output:
	// 5 4
}
```

## 找出旋转数组的最小元素

```go
func findTheMinimum(array []int, begin, end int) int {
	if begin == end {
		return array[begin]
	}

	mid := (begin + end) / 2

	// 正好取到边缘
	if mid > 0 && array[mid] < array[mid-1] {
		return array[mid]
	} else if mid+1 < end && array[mid] > array[mid+1] {
		return array[mid+1]
	}

	if array[mid] < array[end] {
		return findTheMinimum(array, begin, mid)
	} else if array[mid] > array[begin] {
		return findTheMinimum(array, mid+1, end)
	}

	// array[begin] == array[mid] == array[end]
	left := findTheMinimum(array, begin, mid)
	right := findTheMinimum(array, mid+1, end)
	if left < right {
		return left
	} else {
		return right
	}
}

func Example_findTheMinimum() {
	array := []int{1, 0, 1, 1, 1, 1}
	fmt.Println(findTheMinimum(array, 0, len(array)-1))

	// Output:
	// 0
}
```

## 将数组的后面 m 个数移动为前面 m 个数

```go
func Spin(array []int, n, m int) {
	l := n - m
	ReverseArray(array, 0, l-1)
	ReverseArray(array, l, n-1)
	ReverseArray(array, 0, n-1)
}

func ExampleSpin() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	Spin(array, len(array), 5)
	fmt.Println(array)

	// Output:
	// [6 7 8 9 10 1 2 3 4 5]
}
```

```go

```
