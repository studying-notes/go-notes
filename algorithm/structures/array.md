---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数据结构与算法之数组"  # 文章标题
url:  "posts/go/algorithm/structures/array"  # 设置网页永久链接
tags: [ "algorithm", "go" ]  # 标签
categories: [ "Go 数据结构与算法"]  # 系列

weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: false  # 是否自动生成目录
draft: false  # 草稿
---

## 数组

数组是某种类型的数据按照一定的顺序组成的数据的集合。如果将有限个类型相同的变量的集合命名，那么这个名称为数组名。组成数组的各个变量称为数组的分量，也称为数组的元素，有时也称为下标变量。用于区分数组的各个元素的数字编号称为下标。

- [假设该数为 a，x 中必然有 a，所以抵消，](#假设该数为-ax-中必然有-a所以抵消)
- [其他元素因为都是偶数次，所以也都抵消了](#其他元素因为都是偶数次所以也都抵消了)
	- [求自然数序列中重复出现的自然数](#求自然数序列中重复出现的自然数)
	- [查找数组中元素的最大值和最小值](#查找数组中元素的最大值和最小值)
		- [暴力比较](#暴力比较)
		- [分治法](#分治法)
	- [找出旋转数组的最小元素](#找出旋转数组的最小元素)
	- [实现旋转数组/循环移位](#实现旋转数组循环移位)
		- [内置函数实现](#内置函数实现)
		- [原地三次逆序](#原地三次逆序)
		- [性能对比](#性能对比)
	- [找出数组中第 k 小的数](#找出数组中第-k-小的数)
		- [排序法](#排序法)
		- [部分排序法](#部分排序法)
		- [类快速排序法](#类快速排序法)
	- [在不排序的情况下求数组中的中位数](#在不排序的情况下求数组中的中位数)
	- [在 O(n) 时间复杂度内查找数组中前三名](#在-on-时间复杂度内查找数组中前三名)
	- [求数组中两个元素的最小距离](#求数组中两个元素的最小距离)
		- [双重遍历法](#双重遍历法)
		- [动态规划法](#动态规划法)
	- [求解最小三元组距离](#求解最小三元组距离)
		- [最小距离法](#最小距离法)
	- [求数组中绝对值最小的数](#求数组中绝对值最小的数)
		- [顺序遍历法](#顺序遍历法)
		- [二分查找法](#二分查找法)
	- [求数组连续最大和](#求数组连续最大和)
	- [动态规划方法](#动态规划方法)
	- [确定最大子数组的位置](#确定最大子数组的位置)
	- [找出数组中出现 1 次的数](#找出数组中出现-1-次的数)
	- [将二维数组逆时针旋转 45° 后打印](#将二维数组逆时针旋转-45-后打印)
	- [求集合的所有子集](#求集合的所有子集)
		- [位图法](#位图法)
		- [迭代法](#迭代法)
	- [在有规律的二维数组中进行高效的数据查找](#在有规律的二维数组中进行高效的数据查找)
		- [二分查找法](#二分查找法-1)
	- [寻找覆盖点最多的路径](#寻找覆盖点最多的路径)
	- [判断请求能否在给定的存储条件下完成](#判断请求能否在给定的存储条件下完成)
	- [根据规则构造新的数组](#根据规则构造新的数组)
	- [求解迷宫问题](#求解迷宫问题)
	- [从三个有序数组中找出它们的公共元素](#从三个有序数组中找出它们的公共元素)
	- [求两个有序集合的交集](#求两个有序集合的交集)
	- [对任务进行调度](#对任务进行调度)
	- [对磁盘分区](#对磁盘分区)
	- [合并两个有序数组](#合并两个有序数组)
	- [求一个数组的所有排列](#求一个数组的所有排列)
	- [最长方连续方波信号](#最长方连续方波信号)
	- [在二维数组中寻找最短路线](#在二维数组中寻找最短路线)

## 生成指定范围的整数切片

> 源码位置 *src/algorithm/structures/array/range.go*

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

## 找出数组中唯一的重复元素

数字 1～1000 放在含有 1001 个元素的数组中，其中只有唯一的一个元素值重复，其他数字均只出现一次。设计一个算法，将重复元素找出来，要求每个数组元素只能访问一次。

> 源码位置 *src/algorithm/structures/array/find_repeating_elements.go*

### 数据映射法/哈希表

访问一个元素然后将它的值作为索引访问下一个元素，同时标记已访问元素，一旦再次访问同一元素即是重复元素，可以取相反数作为标记，但这样修改了数组中元素的值。

### 累加求和法

将 1001 个数求和然后减去 1~1000 的和就是多出来的那个数。

### 累计差值求和

这种方法是对累加求和法的改进。就是一边加一边减，而不是全部加起来减，因为如果累加的数值巨大时，全部求和就很有可能溢出了。

### 异或法

异或运算的性质：当相同元素异或时，其运算结果为 0，当相异元素异或时，其运算结果为非 0，任何数与数字 0 进行异或运算，其运算结果为该数。

| 异或运算  | 相与运算  |
| --------- | --------- |
| 1 ^ 1 = 0 | 1 ^ 1 = 1 |
| 1 ^ 0 = 1 | 1 ^ 0 = 0 |
| 0 ^ 0 = 0 | 0 ^ 0 = 0 |

```go
// 异或法
func xor() int {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
	var x int
	for k, v := range array {
		x ^= k ^ v
	}
	return x
}
```

## 找出数组中丢失的数

给定一个由 n-1 个整数组成的未排序的数组序列，其元素都是 1～n 中的不同的整数。写出一个寻找数组序列中缺失整数的线性时间算法。

方法同上。

## 找出数组中出现奇数次的数

数组中有 N+2 个数，其中，N 个数出现了偶数次，2 个数出现了奇数次（这两个数不相等），请用 O(1) 的空间复杂度，找出这两个数。

根据异或运算的性质，任何一个数字异或它自己其结果都等于 0。所以，对于本题中的数组元素而言，如果从头到尾依次异或每一个元素，那么异或运算的结果自然也就是那个只出现奇数次的数字，因为出现偶数次的数字会通过异或运算全部消掉。

但是通过异或运算，也仅仅只是消除掉了所有出现偶数次数的数字，最后结果肯定是那两个出现了奇数次的数进行异或运算的结果。

假设这两个出现奇数次的数分别为 a与 b，根据异或运算的性质，将二者异或运算的结果记为 c，由于 a 与 b 不相等，所以，c 的值自然也不会为 0。

此时只需知道 c 对应的二进制数中某一个位为 1 的位数 N，例如，十进制数 44 可以由二进制 0010 1100 表示，此时可取 N=2 或者 3，或者 5，然后将 c 与数组中第 N 位为 1 的数进行异或（N 取其中任意一个即可），异或结果就是 a，b 中一个，然后用 c 异或其中一个数，就可以求出另外一个数了。

通过上述方法为什么就能得到问题的解呢？

**因为 c 中第 N 位为 1 表示 a 或 b 中有一个数的第 N 位也为 1**，假设该数为 a，那么，当将 c 与数组中第 N 位为 1 的数进行异或时，也就是将 x 与 a 外加上其他第 N 位为 1 的出现过偶数次的数进行异或，化简即为 x 与 a 异或，**其他元素因为都是偶数次，所以都抵消了**，结果即为 b。

````bash
a ^ b = c
a ^ b ^ x = c ^ x
# 假设该数为 a，x 中必然有 a，所以抵消，
# 其他元素因为都是偶数次，所以也都抵消了
b = c ^ x
````

> 源码位置 *src/algorithm/structures/array/find_odd_elements.go*

```go
func odd(array []int) [2]int {
	var s int

	for _, v := range array {
		s ^= v
	}

	x := s
	pos := 0

	// 找到为 1 的位置
	for i := s; i&1 == 0; i = i >> 1 {
		pos++
	}

	// 与所有 pos 位置为 1 的数异或
	for _, v := range array {
		if (v>>pos)&1 == 1 {
			s ^= v
		}
	}

	return [2]int{s, x ^ s}
}
```

## 求自然数序列中重复出现的自然数

不对序列做任何修改。

最小值 1 最大值 N 的自然数序列，序列长度为 M，其中存在重复出现的自然数。

每次将重复数字中的一个改为靠近 N+M 的自然数，让遍历能访问到数组后面的元素，就能将整个数组遍历完。

哈希表即可。

## 查找数组中元素的最大值和最小值

给定数组 a1， a2， a3， … an，要求找出数组中的最大值和最小值。假设数组中的值两两各不相同。

### 暴力比较

定义两个变量，遍历数组元素，每次都比较两次（2n 次），找出最大值和最小值。

### 分治法

分治法就是将一个规模为n的、难以直接解决的大问题，分割为k个规模较小的子问题，采取各个击破、分而治之的策略得到各个子问题的解，然后将各个子问题的解进行合并，从而得到原问题的解的一种方法。

不断二分，找出各自中的最值。

> 源码位置 *src/algorithm/structures/array/find_max_min.go*

```go
func GetMaxAndMin(array []int, start, end int) (min, max int) {
	if start == end { // 这里不会越界
		return array[start], array[start]
	}

	mid := (start + end) / 2
	l1, l2 := GetMaxAndMin(array, start, mid)
	r1, r2 := GetMaxAndMin(array, mid+1, end)

	if l1 < r1 {
		min = l1
	} else {
		min = r2
	}

	if l2 > r2 {
		max = l2
	} else {
		max = r2
	}

	return min, max
}
```

比较次数为 3n/2-2。

## 找出旋转数组的最小元素

把一个有序数组最开始的若干个元素搬到数组的末尾，称为**数组的旋转**。

输入一个排好序的数组的一个旋转，输出旋转数组的最小元素。

例如数组 {3, 4, 5, 1, 2} 为数组 {1, 2, 3, 4, 5} 的一个旋转，该数组的最小值为 1。

通过数组的特性可以发现，数组元素首先是递增的，然后突然下降到最小值，然后再递增。虽然如此，但是还有下面三种特殊情况需要注意：

1. 数组本身是没有发生过旋转的，是一个有序的数组，例如序列 {1, 2, 3, 4, 5, 6}。
2. 数组中元素值全部相等，例如序列 {1, 1, 1, 1, 1, 1}。
3. 数组中元素值大部分都相等，例如序列 {1, 0, 1, 1, 1, 1}。

> 源码位置 *src/algorithm/structures/array/find_max_min.go*

```go
func findTheSmallestElementOfARotatedArray(array []int, start, end int) int {
	if start == end {
		return array[start]
	}

	mid := (start + end) / 2

	// 正好取到边缘
	if mid > 0 && array[mid] < array[mid-1] {
		return array[mid]
	} else if mid+1 < end && array[mid] > array[mid+1] {
		return array[mid+1]
	}

	if array[mid] < array[end] {
		return findTheSmallestElementOfARotatedArray(array, start, mid)
	} else if array[mid] > array[start] {
		return findTheSmallestElementOfARotatedArray(array, mid+1, end)
	} else { // array[start] == array[mid] == array[end]
		left := findTheSmallestElementOfARotatedArray(array, start, mid)
		right := findTheSmallestElementOfARotatedArray(array, mid+1, end)
		if left < right {
			return left
		} else {
			return right
		}
	}
}
```

## 实现旋转数组/循环移位

> 源码位置 *src/algorithm/structures/array/spin.go*

### 内置函数实现

```go
// 利用内置函数
func SpinArrayAppend(array []int, idx int) []int {
	return append(array[idx:], array[:idx]...)
}
```

### 原地三次逆序

```go
// 逆序数组
func ReverseArray(array []int) {
	start, end := 0, len(array)-1
	for start < end {
		array[start], array[end] = array[end], array[start]
		start, end = start+1, end-1
	}
}

// 多次逆序
func SpinArrayReverse(array []int, idx int) {
	ReverseArray(array[:idx])
	ReverseArray(array[idx:])
	ReverseArray(array)
}
```

### 性能对比

```go
func BenchmarkSpinArrayAppend(b *testing.B) {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayAppend(array, 5)
	}
}

func BenchmarkSpinArrayReverse(b *testing.B) {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayReverse(array, 5)
	}
}
```

内置 append 方法复制时消耗了大量内存空间和时间。

## 找出数组中第 k 小的数

给定一个整数数组，如何快速地求出该数组中第 k 小的数。假如数组为 {4, 0, 1, 0, 2, 3}，那么第 3 小的元素是 1。

> 源码位置 *src/algorithm/structures/array/find_kth_smallest.go*

### 排序法

最简单的方法就是首先对数组进行排序，在排序后的数组中，下标为 k-1 的值就是第 k 小的数。由于最高效的排序算法的平均时间复杂度为 O(nlogn)，因此，此时该方法的平均时间复杂度为 O(nlogn)，其中，n 为数组的长度。

### 部分排序法

由于只需要找出第 k 小的数，因此，没必要对数组中所有的元素进行排序，可以采用部分排序的方法。这种方法的时间复杂度为 O(n*k)。

当然也可以采用堆排序进行 k 趟排序找出第 k 小的值。

下面是冒泡排序的方法：

```go
func findKthSmallest(array []int, k int) (j int) {
	if k > len(array) {
		panic("k must < array's length")
	}

	for k != 0 {
		for i := len(array) - 1; i > j; i-- {
			if array[i] < array[i-1] {
				array[i], array[i-1] = array[i-1], array[i]
			}
		}
		k, j = k-1, j+1
	}

	return array[j-1]
}
```

### 类快速排序法

快速排序就是将数组 `array[low…high]` 中**第一个元素作为划分依据**，然后**把数组划分为三部分**：

1. `array[low…i-1]` 所有的元素的值都小于或等于 `array[i]`
2. `array[i]`
3. `array[i+1…high]` 所有的元素的值都大于 `array[i]`

类似二叉平衡搜索树的原理。

```go
// Partition 用于快速排序中的分割
func Partition(array []int, low, high int) int {
	i, j, val := low, high, array[low]
	for i < j {
		for i < j && array[j] >= val {
			j--
		}
		if i < j {
			array[i] = array[j]
		}
		for i < j && array[i] <= val {
			i++
		}
		if i < j {
			array[j] = array[i]
		}
	}
	array[i] = val
	return i // 中间值索引
}

func findSmallK(array []int, low, high, k int) {
	if low > high {
		return
	}
	pos := Partition(array, low, high)
	if pos+1 == k {
		return
	} else if pos+1 < k {
		// 取右侧
		findSmallK(array, pos+1, high, k)
	} else {
		// 取左侧
		findSmallK(array, low, pos-1, k)
	}
}
```

## 在不排序的情况下求数组中的中位数

把问题转化为求一列数中第 i 小的数的问题，求中位数就是求一列数的第（length/2+1）小的数的问题。

用部分排序法。

## 在 O(n) 时间复杂度内查找数组中前三名

```go
func findTop3(array []int) (r1, r2, r3 int) {
	if array == nil || len(array) < 3 {
		return
	}
	for _, v := range array {
		if v > r1 {
			r3 = r2
			r2 = r1
			r1 = v
		} else if v > r2 && v != r1 {
			r3 = r2
			r2 = v
		} else if v > r3 && v != r2 {
			r3 = v
		}
	}
	return
}
```

这种方法虽然能够在 O(n) 的时间复杂度求出前三名，但是当 k 取值很大的时候，比如求前 10 名，这种方法就不是很好了。

比较经典的方法就是维护一个大小为 k 的堆来保存最大的 k 个数。

最大堆，又称大根堆（大顶堆）是指根结点（亦称为堆顶）的关键字是堆里所有结点关键字中最大者，属于二叉堆的两种形式之一。

维护一个大小为 k 的小顶堆用来存储最大的 k 个数，堆顶保存了，每次遍历一个数 m，如果 m 比堆顶元素小，那么说明 m 肯定不是最大的 k 个数，因此，不需要调整堆，如果 m 比堆顶与元素大，则用这个数替换堆顶元素，替换后重新调整堆为小顶堆。这种方法的时间复杂度为 O(n * logk)。这种方法适用于数据量大的情况。

具体看堆相关的数据结构。

## 求数组中两个元素的最小距离

给定一个数组，数组中含有重复元素，给定两个数字 num1 和 num2，求这两个数字在数组中出现的位置的最小距离。

### 双重遍历法

最简单的，暴力比较。

```go
func findDist(array []int, n1, n2 int) (dist int) {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1
	}
	dist = 1<<63 - 1
	d := 0
	for i, v := range array {
		if v == n1 {
			for j, v := range array {
				if v == n2 {
					d = Abs(i - j)
					if d < dist {
						dist = d
					}
				}
			}
		}
	}
	return dist
}
```

### 动态规划法

其实就是把遍历到的最近的点记录下来。

```go
func findDistDyn(array []int, n1, n2 int) (dist int) {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1
	}

	// 记录当前最小距离
	dist = 1<<63 - 1

	d := 0

	// Go 中移位运算符优先级大于减号
	// Python 中减号优先级大于移位运算符

	x, y := -1<<31, 1<<31-1 // 防溢出

	for i, v := range array {
		if v == n1 {
			x = i
		} else if v == n2 {
			y = i
		}
		d = Abs(x - y)
		if d < dist {
			dist = d
		}
	}

	return dist
}
```

## 求解最小三元组距离

已知三个升序整数数组 `a[l]`，`b[m]` 和 `c[n]`，请在三个数组中各找一个元素，使得组成的三元组距离最小。三元组距离的定义是：假设 `a[i]`、`b[j]` 和 `c[k]` 是一个三元组，那么距离为：`Distance=max(|a[i]-b[j]|`，`|a[i]-c[k]|`，`|b[j]-c[k]|`)，请设计一个求最小三元组距离的最优算法。

### 最小距离法

从三个数组的第一个元素开始，首先求出它们的最小距离 Dist，接着找出这三个数中最小数所在的数组，只对这个数组的下标往后移一个位置，接着求三个数组中当前遍历元素的距离，如果比 Dist 小，则把当前距离赋值给 Dist，依此类推，直到遍历完其中一个数组为止。

```go
// 求解最小三元组距离
func findDist3Array(a1, a2, a3 []int) int {
	b1, b2, b3 := 0, 0, 0
	dist := 1<<63 - 1
	for b1 < len(a1) && b2 < len(a2) && b3 < len(a3) {
		// 通过数学等式变形可以得出以下两种计算结果等价
		// d := 1/2 * (Abs(a1[b1]-a2[b2]) + Abs(a1[b1]-a3[b3]) + Abs(a2[b2]-a3[b3]))
		if d := MaxN(Abs(a1[b1]-a2[b2]), Abs(a1[b1]-a3[b3]), Abs(a2[b2]-a3[b3])); d < dist {
			dist = d
		}
		switch MinN(a1[b1], a2[b2], a3[b3]) {
		case a1[b1]:
			b1++
		case a2[b2]:
			b2++
		case a3[b3]:
			b3++
		}
	}
	return dist
}
```

## 求数组中绝对值最小的数

有一个升序排列的数组，数组中可能有正数、负数或 0，求数组中元素的绝对值最小的数。例如，数组 {-10, -5, -2,  7,  15,  50}，该数组中绝对值最小的数是 -2。

跟找出旋转数组的最小元素一样的思路。

### 顺序遍历法

```go
func findMinAbs(array []int) int {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1 // 需要表示负数和整数，所以最大值就是它，0 有 +0（全1） 和 -0（全0）
	}
	if array[len(array)-1] <= 0 {
		return array[len(array)-1]
	} else if array[0] >= 0 {
		return array[0]
	}
	for k, v := range array {
		if v > 0 {
			if Abs(array[k-1]) > Abs(array[k]) {
				return Abs(array[k])
			}
			return array[k-1]
		}
	}
	return 1<<63 - 1 // 实际上已经没有其他可能性了
}
```

### 二分查找法

```go
func findMinBinary(array []int, start, end int) int {
	if start == end {
		return array[start]
	}
	if array[end] <= 0 {
		return array[end]
	} else if array[start] >= 0 {
		return array[start]
	}
	mid := (start + end) / 2
	if array[mid] == 0 {
		return 0
	} else if array[mid] > 0 {
		return findMinBinary(array, start, mid)
	} else {
		return findMinBinary(array, mid+1, end)
	}
}
```

## 求数组连续最大和

一个有 n 个元素的数组，这 n 个元素既可以是正数也可以是负数，数组中连续的一个或多个元素可以组成一个连续的子数组，一个数组可能有多个这种连续的子数组，求子数组和的最大值。

> 源码位置 *src/algorithm/structures/array/find_the_maximum_consecutive_sum_of_an_array.go*

## 动态规划方法

```go
func MaxSubArraySumDyn(array []int) (sum int) {
	if array == nil || len(array) == 0 {
		return -1 << 63
	}
	cur := 0
	for _, v := range array {
		cur = Max(cur+v, 0)
		sum = Max(cur, sum)
	}
	return sum
}
```

## 确定最大子数组的位置

> 源码位置 *src/algorithm/structures/array/find_the_maximum_consecutive_sum_of_an_array.go*

方法与上面相同，不过增加一些条件后处理。

```go
func findSubArray(array []int) []int {
	start, end := 0, 1
	sum, cur := 0, 0
	for k, v := range array {
		if cur+v >= 0 {
			cur += v
		} else {
			start, end, cur = k+1, k+2, 0
		}
		if cur > sum {
			end = k + 1
			sum = cur
		}
	}
	return array[start:end]
}
```

## 找出数组中出现 1 次的数

一个数组里，除了三个数是唯一出现的，其余的数都出现偶数次，找出这三个数中的**任意一个**。

与前文的异或方法原理是相同的。**把三个数其中的两个数看成一个整体**。

> 源码位置 *src/algorithm/structures/array/find_once.go*

```go
func findOnce(array []int) int {
	var s int
	
	for _, v := range array {
		s ^= v
	}

	x := s
	pos := 0

	// 找到为 1 的位置
	for i := s; i&1 == 0; i = i >> 1 {
		pos++
	}

	// 与所有 pos 位置为 1 的数异或
	for _, v := range array {
		if (v>>pos)&1 == 1 {
			s ^= v
		}
	}

	return x ^ s
}
```

## 将二维数组逆时针旋转 45° 后打印

相当于沿着对角线打印。

> 源码位置 *src/algorithm/structures/array/matrix.go*

```go
// PrintMatrix 打印二维数组/矩阵
func PrintMatrix(matrix [][]int) {
	repr := ""

	for _, i := range matrix {
		for n, j := range i {
			if n > 0 {
				repr += ", "
			}
			repr += strconv.Itoa(j)
		}

		repr += "\n"
	}

	fmt.Println(repr)
}
```

```go
func ExamplePrintMatrix() {
	matrix := InitMatrix(3, 3)

	matrix[0][0] = 1
	matrix[0][1] = 2
	matrix[0][2] = 3
	matrix[1][0] = 4
	matrix[1][1] = 5
	matrix[1][2] = 6
	matrix[2][0] = 7
	matrix[2][1] = 8
	matrix[2][2] = 9

	PrintMatrix(matrix)

	// Output:
	// 1, 2, 3
	// 4, 5, 6
	// 7, 8, 9
}
```

```go
func PrintRotateMatrix(matrix [][]int) {
	repr := ""
	var i, x, k int

	length := len(matrix)
	j := length - 1
	for x < length {
		i = x
		for i < length && j < length {
			if i > x {
				repr += ", "
			}
			repr += fmt.Sprint(matrix[i][j])
			i++
			j++
		}
		if i == length {
			x++
			j = 0
		} else {
			k++
			j = length - k - 1
		}

		repr += "\n"
	}

	fmt.Println(repr)
}
```

```go
func ExamplePrintRotateMatrix() {
	matrix := InitMatrix(3, 3)

	matrix[0][0] = 1
	matrix[0][1] = 2
	matrix[0][2] = 3
	matrix[1][0] = 4
	matrix[1][1] = 5
	matrix[1][2] = 6
	matrix[2][0] = 7
	matrix[2][1] = 8
	matrix[2][2] = 9

	PrintRotateMatrix(matrix)

	// Output:
	// 3
	// 2, 6
	// 1, 5, 9
	// 4, 8
	// 7
}
```

## 求集合的所有子集

有一个集合，求其全部子集（包含集合自身）。

给定一个集合 s，它包含两个元素 `<a, b>`，则其全部的子集为 `<a, ab, b>`。

子集个数 Sn 与原集合元素个数 n 之间的关系满足如下等式：Sn=2^n-1。

> 源码位置 *src/algorithm/structures/array/find_all_subsets.go*

### 位图法

1. 构造一个和集合一样大小的数组 A，分别与集合中的某个元素对应，数组 A 中的元素只有两种状态：“1” 和 “0”，分别代表每次子集输出中集合中对应元素是否要输出，这样数组 A 可以看作是原集合的一个标记位图。

2. 数组 A 模拟整数“加 1”的操作，每执行“加 1”操作之后，就将原集合中所有与数组 A 中值为“1”的相对应的元素输出。

设原集合为 `<a, b, c, d>`，数组 A 的某次“加1”后的状态为 `[1, 0, 1, 1]`，则本次输出的子集为 `<a, c, d>`。使用非递归的思想，如果有一个数组，大小为 n，那么就使用 n 位的二进制，如果对应的位为 1，那么就输出这个位，如果对应的位为 0，那么就不输出这个位。

实际中可以用位运算代替数组：

```go
func findAllSubsets(set []int) (subsets [][]int) {
	o := 1 << len(set)

	for q := 1; q < o; q++ {
		var subset []int
		for p := 0; p < len(set); p++ {
			if q>>p&1 == 1 {
				subset = append(subset, set[p])
			}
		}
		subsets = append(subsets, subset)
	}

	return
}
```

时间复杂度为 O(n*2^n)，空间复杂度 O(n)。

### 迭代法

每次迭代，都是上一次迭代的结果+上次迭代结果中每个元素都加上当前迭代的元素+当前迭代的元素。

```go
func findAllSubsets2(set []int) [][]int {
	subsets := [][]int{{set[0]}}
	for i := 1; i < len(set); i++ {
		l := len(subsets)
		for j := 0; j < l; j++ {
			subsets = append(subsets, append(subsets[j], set[i]))
		}
		subsets = append(subsets, []int{set[i]})
	}
	return subsets
}
```

时间复杂度为 O(2^n)，空间复杂度 2^(n-1)-1。

## 在有规律的二维数组中进行高效的数据查找

在一个二维数组中，每一行都按照从左到右递增的顺序排序，每一列都按照从上到下递增的顺序排序。请实现一个函数，输入这样的一个二维数组和一个整数，判断数组中是否含有该整数。

> 源码位置 *src/algorithm/structures/array/search.go*

### 二分查找法

关键在于想到从二维数组的右上角遍历到左下角。

右对角线的位置就是二分的位置，其左边都更小，下边都更大。

```go
// IsContainK 二分查找有序矩阵中是否存在某个元素
func IsContainK(array2d [][]int, k int) bool {
	if len(array2d) == 0 {
		return false
	}
	
	// 行列
	rows, columns := len(array2d), len(array2d[0])

	// 从右上角开始查找
	for i, j := 0, columns-1; i < rows && j > 0; {
		if array2d[i][j] == k {
			return true
		} else if array2d[i][j] > k {
			j--
		} else {
			i++
		}
	}
	return false
}
```

## 寻找覆盖点最多的路径

坐标轴上从左到右依次的点为 `a[0]、a[1]、a[2]…a[n-1]`，求满足：

```
a[j]-a[i]<=L
a[j+1]-a[i]>L
``` 

这两个条件的 `j` 与 `i` 中间的所有点个数中的最大值，即 `j-i+1` 最大。

```go
// FindThePathWithTheMostCoveragePoints 寻找覆盖点最多的路径
func FindThePathWithTheMostCoveragePoints(array []int, threshold int) (int, []int) {
	// begin 为起始点，end 为终点，count 为覆盖点数
	begin, end, count := 0, 0, 0

	// 从第一个点开始遍历
	for idx := range array {
		// 如果当前点与上一个起始点的差值大于阈值，则起始点向后移动一位
		for begin < idx && array[idx]-array[begin] > threshold {
			begin++
		}

		// 更新覆盖点数
		if idx-begin > count && array[idx+1]-array[begin] > threshold {
			count = idx - begin + 1 // 覆盖点数
			end = idx + 1 // end 取不到最后一个点，所以要加 1
		}
	}

	return count, array[end-count : end]
}
```

## 判断请求能否在给定的存储条件下完成

给定一台有 m 个存储空间的机器，有 n 个请求需要在这台机器上运行，第 i 个请求计算时需要占 `R[i]` 空间，计算结果需要占 `O[i]` 个空间 `O[i] < R[i]`。请设计一个算法，判断这 n 个请求能否全部完成？若能，给出这 n 个请求的安排顺序。

首先对请求按照 `R[i]-O[i]` 由大到小进行排序，然后按照由大到小的顺序进行处理，如果按照这个顺序能处理完，则这 n 个请求能被处理完，否则处理不完。

```go
func schedule(R, O []int, M int) bool {
	sort(R, O)
	left := M
	for idx := range R {
		if R[idx] > left {
			return false
		}
		left -= O[idx]
	}
	return true
}

func swap(array []int, i, j int) {
	array[i], array[j] = array[j], array[i]
}

// 冒泡排序
func sort(R, O []int) {
	for i := len(R); i > 0; i-- {
		for j := 1; j < i; j++ {
			if R[j]-O[j] > R[j-1]-O[j-1] {
				swap(R, j-1, j)
				swap(O, j-1, j)
			}
		}
	}
}

func ExampleScheduler() {
	R := []int{10, 15, 23, 20, 6, 9, 7, 16} // 计算时占用的空间
	O := []int{2, 7, 8, 4, 5, 8, 6, 8}      // 计算结果占用的空间

	sort(R, O)

	fmt.Println(R)
	fmt.Println(O)

	// Output:
	// [20 23 10 15 16 6 9 7]
	// [4 8 2 7 8 5 8 6]
}
```

## 根据规则构造新的数组

给定一个数组 `a[N]`，希望构造一个新的数组 `b[N]`，其中，`b[i]=a[0]*a[1]*…*a[N-1]/a[i]`。

1. 不允许使用除法
2. 要求具备 O(1) 空间复杂度和 O(n) 时间复杂度
3. 除遍历计数器与 `a[n]`、`b[n]` 外，不可以使用新的变量（包括栈临时变量、堆空间和全局静态变量等）

首先遍历一遍数组 a，在遍历的过程中对数组 b 进行赋值：

```
b[i]=a[i-1]*b[i-1]
```

这样经过一次遍历后，数组 b 的值为：

```
b[i]=a[0]*a[1]*…*a[i-1]
```

此时只需要将数组中的值 `b[i]` 再乘以` a[i+1]*a[i+2]*…a[N-1]`。

实现方法为逆向遍历数组 a，把数组后半段值的乘积记录到 `b[0]` 中（这个值在此之前只能是 1，所以可以用作临时存储），通过 `b[i]` 与 `b[0]` 的乘积就可以得到满足题目要求的 `b[i]`。

具体而言，执行 `b[i]=b[i]*b[0]`（首先执行的目的是为了保证在执行下面一个计算的时候，`b[0]` 中不包含与 `b[i]` 的乘积），接着记录数组后半段的乘积到 `b[0]` 中：`b[0]*=b[0]*a[i]`。

```go
func NewArrayByRule(a []int) {
	b := make([]int, len(a))

	b[0] = 1

	for i := 1; i < len(a); i++ {
		b[i] = a[i-1] * b[i-1]
	}

	for j := len(a) - 1; j > 0; j-- {
		b[j] *= b[0]
		b[0] *= a[j]
	}

	fmt.Println(b)
}
```

## 求解迷宫问题

给定一个大小为 N×N 的迷宫，一只老鼠需要从迷宫的左上角（对应矩阵的 `[0][0]`）走到迷宫的右下角（对应矩阵的 `[N-1][N-1]`），老鼠只能向两方向移动：向右或向下。在迷宫中，0 表示没有路（是死胡同），1 表示有路。

尝试可能的路径，遇到岔路口时保存其中一个方向，然后尝试走另一个方向，当碰到死胡同的时候，回溯到前一步，然后从前一步出发继续寻找可达的路径。

> 源码位置 *src/algorithm/structures/array/maze.go*

```go
func MazeSolver(matrix [][]int) [][2]int {
	rows := len(matrix)
	columns := len(matrix[0])

	var s [][2]int           // 岔口
	var r = [][2]int{{0, 0}} // 路径
	var i, j int

	for i < rows-1 || j < columns-1 {
		if i < rows-1 && matrix[i+1][j] == 1 {
			s = append(s, [2]int{i + 1, j})
		}

		if j < columns-1 && matrix[i][j+1] == 1 {
			s = append(s, [2]int{i, j + 1})
		}

		next := s[len(s)-1]
		s = s[:len(s)-1]
		i, j = next[0], next[1]

		for len(r) > 0 {
			last := r[len(r)-1]
			if last[0] != i && last[1] != j {
				// 清除无效路径
				r = r[:len(r)-1]
			} else {
				break
			}
		}

		r = append(r, next)
	}

	return r
}
```

```go
func ExampleMazeSolver() {
	matrix := [][]int{
		{1, 1, 1, 1},
		{1, 0, 0, 0},
		{1, 1, 1, 1},
		{1, 0, 0, 1},
	}

	fmt.Println(MazeSolver(matrix))

	// Output:
	// [[0 0] [1 0] [2 0] [2 1] [2 2] [2 3] [3 3]]
}
```

## 从三个有序数组中找出它们的公共元素

给定以**非递减顺序排序**的三个数组，找出这三个数组中的所有公共元素。例如，给出下面三个数组：

```go
a1 := []int{2, 5, 12, 20, 45, 85}
a2 := []int{16, 19, 20, 85, 200}
a3 := []int{3, 4, 15, 20, 39, 72, 85, 190}
```

那么这三个数组的公共元素为 {20, 85}。

```go
func FindCommonElements(a1, a2, a3 []int) (com []int) {
	i, j, k := 0, 0, 0
	for i < len(a1) && j < len(a2) && k < len(a3) {
		if a1[i] == a2[j] && a2[j] == a3[k] {
			com = append(com, a1[i])
			i, j, k = i+1, j+1, k+1
		} else {
			min := MinN(a1[i], a2[j], a3[k])
			if min == a1[i] && i < len(a1) {
				i++
			} else if min == a2[j] {
				j++
			} else {
				k++
			}
		}
	}
	return com
}

func ExampleFindCommonElements() {
	a1 := []int{2, 5, 12, 20, 45, 85}
	a2 := []int{16, 19, 20, 85, 200}
	a3 := []int{3, 4, 15, 20, 39, 72, 85, 190}

	fmt.Println(FindCommonElements(a1, a2, a3))

	// Output:
	// [20 85]
}
```

## 求两个有序集合的交集

有两个有序的集合，集合中的每个元素都是一段范围，求其交集，例如集合 `{[4, 8]，[9, 13]}` 和 `{[6, 12]}` 的交集为 `{[6, 8]，[9, 12]}`。

```go
func FindSetIntersection(s1, s2 [][]int) (com [][2]int) {
	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		start1, end1 := s1[i][0], s1[i][len(s1[i])-1]
		start2, end2 := s2[j][0], s2[j][len(s2[j])-1]
		if start1 <= start2 && end1 >= start2 {
			com = append(com, [2]int{start2, end1})
			i++
		} else if start2 <= start1 && end2 >= start1 {
			com = append(com, [2]int{start1, end2})
			j++
		} else if end1 < start2 {
			i++
		} else if end2 < start1 {
			j++
		}
	}
	return com
}
```

## 对任务进行调度

假设有一个中央调度机，有 n 个相同的任务需要调度到 m 台服务器上去执行，由于每台服务器的配置不一样，因此，服务器执行一个任务所花费的时间也不同。现在假设第 i 个服务器执行一个任务所花费的时间也不同。假设第 i 个服务器执行一个任务需要的时间为 `t[i]`。

例如，有 2 个执行机 a 与 b，执行一个任务分别需要 7min 和 10min，有 6 个任务待调度。如果平分这 6 个任务，即 a 与 b 各 3 个任务，则最短需要 30min 执行完所有。如果 a 分 4 个任务，b 分 2 个任务，则最短 28min 执行完。请设计调度算法，使得所有任务完成所需要的时间最短。

```go
func FindMinCost(count int) int {
	servers := [2]int{7, 10}
	costs := [2]int{0, 0}

	for ; count > 0; count-- {
		if costs[1]+servers[1] < costs[0]+servers[0] {
			costs[1] += servers[1]
		} else {
			costs[0] += servers[0]
		}
	}

	return Max(costs[0], costs[1])
}

func ExampleFindMinCost() {
	fmt.Println(FindMinCost(6))

	// Output:
	// 28
}
```

## 对磁盘分区

有 N 个磁盘，每个磁盘大小为 `D[i]（i=0…N-1）`，现在要在这 N 个磁盘上"顺序分配" M 个分区，每个分区大小为 `P[j]（j=0…M-1）`，顺序分配的意思是：分配一个分区 P[j] 时，如果当前磁盘剩余空间足够，则在当前磁盘分配；如果不够，则尝试下一个磁盘，直到找到一个磁盘 `D[i+k]` 可以容纳该分区，分配下一个分区 `P[j+1]` 时，则从当前磁盘 `D[i+k]` 的剩余空间开始分配，不在使用 `D[i+k]` 之前磁盘未分配的空间，如果这 M 个分区不能在这 N 个磁盘完全分配，则认为分配失败。

判断给定 N 个磁盘（数组D）和 M 个分区（数组P），是否会出现分配失败的情况？举例：磁盘为 [120，120，120]，分区为 [60，60，80，20，80] 可分配，如果为 [60，80，80，20，80]，则分配失败。

```go
func PartitionDisk(D, P []int) bool {
	// D 磁盘
	// P 分区

	i, j := 0, 0
	for i < len(D) && j < len(P) {
		for i < len(D) && P[j] > D[i] {
			i++
		}
		if i == len(D) {
			break
		}
		D[i] -= P[j]
		j++
	}

	return i != len(D)
}

func ExamplePartitionDisk() {
	D := []int{120, 120, 120} // N

	fmt.Println(PartitionDisk(D, []int{60, 60, 80, 20, 80}))
	fmt.Println(PartitionDisk(D, []int{60, 80, 80, 20, 80}))

	// Output:
	// true
	// false
}
```

## 合并两个有序数组

一般是允许申请新的数组空间，如果不允许就是普通的排序，没有任何已排序带来的优势。

```go
// MergeSortedArrays 允许申请新的数组空间
func MergeSortedArrays(a, b []int) (result []int) {
	aLeft, bLeft := 0, 0
	aRight, bRight := len(a), len(b)
	result = make([]int, aRight+bRight)

	for aLeft < aRight && bLeft < bRight {
		if a[aLeft] > b[bLeft] {
			result[aLeft+bLeft] = b[bLeft]
			bLeft++
		} else {
			result[aLeft+bLeft] = a[aLeft]
			aLeft++
		}
	}

	// 处理多余元素
	for aLeft < aRight {
		result[aLeft+bLeft] = a[aLeft]
		aLeft++
	}

	for bLeft < bRight {
		result[aLeft+bLeft] = b[bLeft]
		bLeft++
	}

	return
}
```

## 求一个数组的所有排列

给定参数 n，从 1 到 n 会有 n 个整数：1,2,3,n，这 n 个数字共有 n！种排列。

按大小顺序升序列出所有排列情况，并一一标记，当= 3 时，所有排列如下：

"123" "132" "213" "231" "312" "321" 

给定 n 和 k，返回第 k 个排列。

> 源码位置 *src/algorithm/structures/array/find_permutation.go*

```go
// 复制切片
func deepcopy(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

// 比较数组的大小
func moreThan(left, right []int) bool {
	leftLength, rightLength := len(left), len(right)

	if leftLength != rightLength {
		return leftLength > rightLength
	}

	for i := 0; i < leftLength; i++ {
		if left[i] != right[i] {
			return left[i] > right[i]
		}
	}

	return false
}

type Permutation struct {
	result [][]int
	array  []int
}

func NewPermutation(n int) *Permutation {
	array := make([]int, n)
	for i := range array {
		array[i] = i + 1
	}

	return &Permutation{array: array}
}

func (p *Permutation) Perform(start int) {
	if start == len(p.array)-1 {
		p.result = append(p.result, deepcopy(p.array))
	} else {
		for i := start; i < len(p.array); i++ {
			p.array[start], p.array[i] = p.array[i], p.array[start]
			p.Perform(start + 1)
			p.array[start], p.array[i] = p.array[i], p.array[start]
		}
	}
}

func (p *Permutation) Sort() {
	// 已经基本有序，所以用一遍冒泡排序即可
	for i := 1; i < len(p.result); i++ {
		if moreThan(p.result[i-1], p.result[i]) {
			p.result[i-1], p.result[i] = p.result[i], p.result[i-1]
		}
	}
}
```

## 最长方连续方波信号

输入一串方波信号，求取最长的完全连续交替方波信号，并将其输出，如果有相同长度的交替方波信号，输出任一即可，方波信号高位用 1 标识，低位用 0 标识，如图：

- 一个完整的信号一定以 0 开始然后以 0 结尾，即 010 是一个完整信号，但 101，1010，0101 不是。
- 输入的一串方波信号是由一个或多个完整信号组成。
- 两个相邻信号之间可能有 0 个或多个低位，如 0110010，011000010。
- 同一个信号中可以有连续的高位，如 01110101011110001010，前 14 位是一个具有连续高位的信号。
- 完全连续交替方波是指 10 交替，如 01010。

```go
const (
	low int = iota
	high
)

func FindLongestSquareContinuousSquareWaveSignal(array []int) []int {
	var (
		begin, end       int // 当前信号
		maxBegin, maxEnd int // 最长的信号
	)

	for i := 1; i < len(array); i++ {
		if array[i-1] == low && array[i] == low {
			// 00
			if begin == 0 || begin == i-1 {
				begin = i
			} else if end == 0 {
				end = i
			}
		} else if array[i-1] == high && array[i] == high {
			// 11
			begin, end = 0, 0
		} else if i == len(array)-1 && array[i-1] == high && array[i] == low {
			// 10
			end = i
		}

		if begin != 0 && end != 0 && end-begin > maxEnd-maxBegin {
			maxBegin, maxEnd = begin, end
			begin, end = 0, 0
		}
	}

	return array[maxBegin:maxEnd]
}
```

```go
func ExampleFindLongestSquareContinuousSquareWaveSignal() {
	arrays := [][]int{
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1},
	}

	for i := range arrays {
		fmt.Println(FindLongestSquareContinuousSquareWaveSignal(arrays[i]))
	}

	// Output:
	// [0 1 0 1 0]
	// [0 1 0 1 0 1 0 1 0]
	// []
}
```

## 在二维数组中寻找最短路线

寻找一条从左上角（`arr[0][0]`）到右下角（`arr[m-1][n-1]`）的路线，使得沿途经过的数组中的整数的和最小。

动态规划法：

```go
const MaxInt32 = 1<<31 - 1

func findShortestRouteIn2dArray(routes [][]int) int {
	rows, columns := len(routes), len(routes[0])

	matrix := InitMatrix(rows+1, columns+1)

	for i := 2; i <= rows; i++ {
		matrix[i][0] = MaxInt32
	}

	for j := 2; j <= columns; j++ {
		matrix[0][j] = MaxInt32
	}

	// PrintMatrix(matrix)

	for i := 1; i <= rows; i++ {
		for j := 1; j <= columns; j++ {
			matrix[i][j] = math.MinN(matrix[i-1][j], matrix[i][j-1]) + routes[i-1][j-1]
		}
	}

	// PrintMatrix(matrix)

	return matrix[rows][columns]
}
```

这种方法对二维数组进行了一次遍历，因此，其时间复杂度为 `O(m*n)`。此外由于这种方法同样申请了一个二维数组来保存中间结果，因此，其空间复杂度也为 `O(m*n)`。

```go

```
