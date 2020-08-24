# 数组

数组是某种类型的数据按照一定的顺序组成的数据的集合。如果将有限个类型相同的变量的集合命名，那么这个名称为数组名。组成数组的各个变量称为数组的分量，也称为数组的元素，有时也称为下标变量。用于区分数组的各个元素的数字编号称为下标。

- [数组](#数组)
	- [生成指定范围的整数切片](#生成指定范围的整数切片)
	- [找出数组中唯一的重复元素](#找出数组中唯一的重复元素)
		- [累加求和法](#累加求和法)
		- [累计差值求和](#累计差值求和)
		- [异或法](#异或法)
		- [数据映射法](#数据映射法)
		- [单链表求环入口点法](#单链表求环入口点法)
	- [找出数组中丢失的数](#找出数组中丢失的数)
		- [累计差值求和](#累计差值求和-1)
		- [异或法求值](#异或法求值)
		- [性能对比](#性能对比)
	- [找出数组中出现奇数次的数](#找出数组中出现奇数次的数)
	- [求重复出现的自然数序列](#求重复出现的自然数序列)
	- [查找数组中元素的最大值和最小值](#查找数组中元素的最大值和最小值)
		- [暴力比较](#暴力比较)
		- [分治法](#分治法)
	- [找出旋转数组的最小元素](#找出旋转数组的最小元素)
	- [实现旋转数组/循环移位](#实现旋转数组循环移位)
		- [内置函数实现](#内置函数实现)
		- [原地三次逆序](#原地三次逆序)
		- [性能对比](#性能对比-1)
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
	- [找出数组中出现1次的数](#找出数组中出现1次的数)
	- [将二维数组逆时针旋转 45° 后打印](#将二维数组逆时针旋转-45-后打印)
	- [求集合的所有子集](#求集合的所有子集)
		- [位图法](#位图法)
		- [迭代法](#迭代法)
	- [在有规律的二维数组中进行高效的数据查找](#在有规律的二维数组中进行高效的数据查找)
		- [二分查找有序数组中是否存在某个元素](#二分查找有序数组中是否存在某个元素)
		- [二分查找有序矩阵中是否存在某个元素](#二分查找有序矩阵中是否存在某个元素)
	- [寻找覆盖点最多的路径](#寻找覆盖点最多的路径)
	- [判断请求能否在给定的存储条件下完成](#判断请求能否在给定的存储条件下完成)
	- [根据规则构造新的数组](#根据规则构造新的数组)
	- [获取最好的矩阵链相乘方法](#获取最好的矩阵链相乘方法)
- [递归法](#递归法)
	- [求解迷宫问题](#求解迷宫问题)
	- [从三个有序数组中找出它们的公共元素](#从三个有序数组中找出它们的公共元素)
	- [求两个有序集合的交集](#求两个有序集合的交集)
	- [对任务进行调度](#对任务进行调度)
	- [对磁盘分区](#对磁盘分区)

## 生成指定范围的整数切片

```go
func Range(args ...int) []int {
	var start, end int
	step := 1
	if len(args) == 1 {
		end = args[0]
	} else if len(args) == 2 {
		start, end = args[0], args[1]
	} else if len(args) > 2 {
		start, end, step = args[0], args[1], args[2]
	}
	if step == 0 || start == end || (step < 0 && start < end) || (step > 0 && start > end) {
		return []int{}
	}
	s := make([]int, 0, (end-start)/step+1)
	for start != end {
		s = append(s, start)
		start += step
	}
	return s
}
```

## 找出数组中唯一的重复元素

数字 1～1000 放在含有 1001 个元素的数组中，其中只有唯一的一个元素值重复，其他数字均只出现一次。设计一个算法，将重复元素找出来，要求每个数组元素只能访问一次。

### 累加求和法

```go
func main() {
    array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}

    // 累加的数值巨大时有可能溢出
	var x, y int
	for k, v := range array {
		x += k
		y += v
	}
    fmt.Println(y - x)
}
```

### 累计差值求和

```go
func main() {
    array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
    
    // 对累加求和法的改进
	var s int
	for k, v := range array {
		s += v - k
	}
	fmt.Println(s)
}
```

### 异或法

异或运算的性质：当相同元素异或时，其运算结果为 0，当相异元素异或时，其运算结果为非 0，任何数与数字 0 进行异或运算，其运算结果为该数。

| 异或运算  | 相与运算  |
| --------- | --------- |
| 1 ^ 1 = 0 | 1 ^ 1 = 1 |
| 1 ^ 0 = 1 | 1 ^ 0 = 0 |
| 0 ^ 0 = 0 | 0 ^ 0 = 0 |

```go
func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
	var x int
	for k, v := range array {
		x ^= k ^ v
	}
	fmt.Println(x)
}
```

似乎与上面改进的累计差值求和是一样的。

### 数据映射法

访问一个元素然后将它的值作为索引访问下一个元素，同时标记已访问元素，一旦再次访问同一元素即是重复元素，可以取相反数作为标记，但这样修改了数组中元素的值。

### 单链表求环入口点法

索引方法同上，但是这种方法多次访问同一元素。


## 找出数组中丢失的数

给定一个由 n-1 个整数组成的未排序的数组序列，其元素都是 1～n 中的不同的整数。写出一个寻找数组序列中缺失整数的线性时间算法。

### 累计差值求和

```go
func main() {
	array := Range(1, 1001)
	fmt.Println(array[500])
	array[500] = 0

	// 累计差值求和
	var s int
	for k, v := range array {
		s += k + 1 - v
	}
	fmt.Println(s)
}
```

### 异或法求值

```go
func main() {
	array := Range(1, 1001)
	fmt.Println(array[500])
	array[500] = 0

	// 异或法求值
	var s int
	for k, v := range array {
		s ^= (k + 1) ^ v
	}
	fmt.Println(s)
}
```

### 性能对比

10 万数值切片运行对比：

```
BenchmarkLossXor-8         24016             48713 ns/op
BenchmarkLossSub-8         20601             57563 ns/op
```

异或方法速度更快。

## 找出数组中出现奇数次的数

数组中有 N+2 个数，其中，N 个数出现了偶数次，2 个数出现了奇数次（这两个数不相等），请用 O(1)的空间复杂度，找出这两个数。

根据异或运算的性质，任何一个数字异或它自己其结果都等于 0。所以，对于本题中的数组元素而言，如果从头到尾依次异或每一个元素，那么异或运算的结果自然也就是那个只出现奇数次的数字，因为出现偶数次的数字会通过异或运算全部消掉。

但是通过异或运算，也仅仅只是消除掉了所有出现偶数次数的数字，最后结果肯定是那两个出现了奇数次的数进行异或运算的结果。假设这两个出现奇数次的数分别为 a与 b，根据异或运算的性质，将二者异或运算的结果记为 c，由于 a 与 b 不相等，所以，c 的值自然也不会为 0，此时只需知道 c 对应的二进制数中某一个位为 1 的位数 N，例如，十进制数 44 可以由二进制 0010 1100 表示，此时可取 N=2 或者 3，或者 5，然后将 c 与数组中第 N 位为 1 的数进行异或，异或结果就是 a，b 中一个，然后用 c 异或其中一个数，就可以求出另外一个数了。

通过上述方法为什么就能得到问题的解呢？其实很简单，**因为 c 中第 N 位为 1 表示 a 或 b 中有一个数的第 N 位也为 1**，假设该数为 a，那么，当将 c 与数组中第 N 位为 1 的数进行异或时，也就是将 x 与 a 外加上其他第 N 位为 1 的出现过偶数次的数进行异或，化简即为 x 与 a 异或，**其他元素都抵消了**，结果即为 b。

```go
// 找出数组中出现奇数次的数
func main() {
	// 0001 0010 0011
	array := []int{1, 2, 3, 1, 2, 3, 1, 2}
	var s int
	for _, v := range array {
		s ^= v
	}
	fmt.Println(s)
	// 0011
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

	fmt.Println(s, x^s)
}
```

## 求重复出现的自然数序列

每次将重复数字中的一个改为靠近 N+M 的自然数，让遍历能访问到数组后面的元素，就能将整个数组遍历完。

```go
func main() {
	array := []int{1, 9, 3, 4, 5, 6, 7, 8, 9, 4, 5, 6, 7, 2, 1}
	s := NewSet()
	len1 := len(array) // 数组长度
	len2 := len1 - 9   // 重复个数
	idx := array[0]
	for len2 > 0 {
		if array[idx] < 0 { // 重复
			len2--
			s.Add(idx)
			array[idx] = len1 - len2 // 使索引向后移动一位，改变了当前数值
		}
		array[idx] *= -1 // 标记为相反数
		idx = -1 * array[idx]
		fmt.Println(array)
	}
	fmt.Println(s.List())
}
```

## 查找数组中元素的最大值和最小值

给定数组 a1， a2， a3， … an，要求找出数组中的最大值和最小值。假设数组中的值两两各不相同。

### 暴力比较

定义两个变量，遍历数组元素，每次都比较两次（2n 次），找出最大值和最小值。


### 分治法

分治法就是将一个规模为n的、难以直接解决的大问题，分割为k个规模较小的子问题，采取各个击破、分而治之的策略得到各个子问题的解，然后将各个子问题的解进行合并，从而得到原问题的解的一种方法。

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

把一个有序数组最开始的若干个元素搬到数组的末尾，称为数组的旋转。输入一个排好序的数组的一个旋转，输出旋转数组的最小元素。例如数组 {3, 4, 5, 1, 2} 为数组 {1, 2, 3, 4, 5} 的一个旋转，该数组的最小值为 1。

通过数组的特性可以发现，数组元素首先是递增的，然后突然下降到最小值，然后再递增。虽然如此，但是还有下面三种特殊情况需要注意：

1. 数组本身是没有发生过旋转的，是一个有序的数组，例如序列 {1, 2, 3, 4, 5, 6}。
2. 数组中元素值全部相等，例如序列 {1, 1, 1, 1, 1, 1}。
3. 数组中元素值大部分都相等，例如序列 {1, 0, 1, 1, 1, 1}。

```go
func BinarySearch(array []int, start, end int) int {
	if start == end {
		return array[start]
	}
    mid := (start + end) / 2
 	// 防止溢出
	if mid > 0 && array[mid] < array[mid-1] {
		return array[mid]
	} else if mid+1 < end && array[mid] > array[mid+1] {
		return array[mid+1]
	}
	if array[mid] < array[end] {
		return BinarySearch(array, start, mid)
	} else if array[mid] > array[start] {
		return BinarySearch(array, mid+1, end)
	} else { // array[start] == array[mid] == array[end]
		left := BinarySearch(array, start, mid)
		right := BinarySearch(array, mid+1, end)
		if left < right {
			return left
		} else {
			return right
		}
	}
}
```

## 实现旋转数组/循环移位

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
	array := Range(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayAppend(array, 500)
	}
}

func BenchmarkSpinArrayReverse(b *testing.B) {
	array := Range(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayReverse(array, 500)
	}
}
```

```shell
BenchmarkSpinArrayAppend-8        121742              9741 ns/op           98304 B/op          1 allocs/op
BenchmarkSpinArrayReverse-8       143238              8418 ns/op               0 B/op          0 allocs/op
```

内置 append 方法复制时消耗了大量内存空间和时间。

## 找出数组中第 k 小的数

给定一个整数数组，如何快速地求出该数组中第k小的数。假如数组为 {4, 0, 1, 0, 2, 3}，那么第3小的元素是 1。

### 排序法

最简单的方法就是首先对数组进行排序，在排序后的数组中，下标为 k-1 的值就是第 k 小的数。由于最高效的排序算法的平均时间复杂度为 O(nlogn)，因此，此时该方法的平均时间复杂度为 O(nlogn)，其中，n 为数组的长度。

### 部分排序法

由于只需要找出第k小的数，因此，没必要对数组中所有的元素进行排序，可以采用部分排序的方法。这种方法的时间复杂度为 O(n*k)。当然也可以采用堆排序进行 k 趟排序找出第k小的值。

```go
func findK(array []int, k int) (j int) {
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

快速排序的基本思想为：将数组 `array[low…high]` 中第一个元素作为划分依据，然后把数组划分为三部分：

1. `array[low…i-1]` 所有的元素的值都小于或等于 `array[i]`
2. `array[i]`
3. `array[i+1…high]` 所有的元素的值都大于 `array[i]`

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
	return i
}

func findKQuick(array []int, low, high, k int) {
	if low > high {
		return
	}
	pos := Partition(array, low, high)
	if pos+1 == k {
		return
	} else if pos+1 < k {
		findKQuick(array, pos+1, high, k)
	} else {
		findKQuick(array, low, pos-1, k)
	}
}
```

## 在不排序的情况下求数组中的中位数

把问题转化为求一列数中第 i 小的数的问题，求中位数就是求一列数的第（length/2+1）小的数的问题。

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

这种方法虽然能够在 O(n) 的时间复杂度求出前三名，但是当 k 取值很大的时候，比如求前 10 名，这种方法就不是很好了。比较经典的方法就是维护一个大小为k的堆来保存最大的 k 个数，具体思路为：维护一个大小为 k 的小顶堆用来存储最大的 k 个数，堆顶保存了，每次遍历一个数 m，如果 m 比堆顶元素小，那么说明 m 肯定不是最大的 k 个数，因此，不需要调整堆，如果 m 比堆顶与元素大，则用这个数替换堆顶元素，替换后重新调整堆为小顶堆。这种方法的时间复杂度为 O(n*logk)。这种方法适用于数据量大的情况。

## 求数组中两个元素的最小距离

给定一个数组，数组中含有重复元素，给定两个数字 num1 和 num2，求这两个数字在数组中出现的位置的最小距离。

### 双重遍历法

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

动态规划的方法是把每次遍历的结果都记录下来从而减少遍历次数。

```go
func findDistDyn(array []int, n1, n2 int) (dist int) {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1
	}
	dist = 1<<63 - 1
	d := 0
	// Go 中移位运算符优先级大于减号
	// Python 中减号优先级大于移位运算符
	x, y := -1<<31, 1<<31-1 // 防溢出
	for i, v := range array {
		if v == n1 {
			x = i
		}
		if v == n2 {
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
// MaxN 求多个数中的最大值
func MaxN(args ...int) int {
	nMax := -1 << 63
	for _, v := range args {
		if v > nMax {
			nMax = v
		}
	}
	return nMax
}

// MinN 求多个数中的最小值
func MinN(args ...int) int {
	nMin := 1<<63 - 1
	for _, v := range args {
		if v < nMin {
			nMin = v
		}
	}
	return nMin
}

// Abs 求绝对值
func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return ^(n - 1)
}
```

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

有一个升序排列的数组，数组中可能有正数、负数或 0，求数组中元素的绝对值最小的数。例如，数组{-10, -5, -2,  7,  15,  50}，该数组中绝对值最小的数是 -2。

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

## 动态规划方法

我一开始就想到了这种方法，其他的就不写了。

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

## 找出数组中出现1次的数

一个数组里，除了三个数是唯一出现的，其余的数都出现偶数次，找出这三个数中的任意一个。与前文的异或方法原理是相同的。把三个数其中的两个数看成一个整体。

```go
func main() {
	array := []int{1, 2, 4, 5, 6, 4, 2}
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
	fmt.Println(x ^ s)
}
```

## 将二维数组逆时针旋转 45° 后打印

```go
// 正常打印二维数组
func PrintMatrix(matrix [][]int) {
	for _, i := range matrix {
		for _, j := range i {
			fmt.Print(j, "  ")
		}
		fmt.Println()
	}
}

```

```go
func PrintRotateMatrix(matrix [][]int) {
	length := len(matrix)
	var i, x, k int
	j := length - 1
	for x < length {
		i = x
		for i < length && j < length {
			fmt.Print(matrix[i][j], "  ")
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
		fmt.Println()
	}
}
```

## 求集合的所有子集

有一个集合，求其全部子集（包含集合自身）。给定一个集合 s，它包含两个元素 `<a, b>`，则其全部的子集为 `<a, ab, b>`。

子集个数 Sn 与原集合元素个数 n 之间的关系满足如下等式：Sn=2^n-1。

### 位图法

1. 构造一个和集合一样大小的数组 A，分别与集合中的某个元素对应，数组 A 中的元素只有两种状态：“1” 和 “0”，分别代表每次子集输出中集合中对应元素是否要输出，这样数组 A 可以看作是原集合的一个标记位图。
2. 数组 A 模拟整数“加1”的操作，每执行“加1”操作之后，就将原集合中所有与数组 A 中值为“1”的相对应的元素输出。

设原集合为 `<a, b, c, d>`，数组 A 的某次“加1”后的状态为 `[1, 0, 1, 1]`，则本次输出的子集为 `<a, c, d>`。使用非递归的思想，如果有一个数组，大小为 n，那么就使用 n 位的二进制，如果对应的位为 1，那么就输出这个位，如果对应的位为 0，那么就不输出这个位。

以上是理论，我在实现的时候用了位运算代替数组遍历：

```go
// Pow 快速幂运算
func Pow(x, n int) int {
	ret := 1
	for n != 0 {
		if n%2 != 0 {
			ret *= x
		}
		n /= 2
		x *= x
	}
	return ret
}

func findAllSubSet(set []int) {
	for q := 1; q < Pow(2, len(set)); q++ {
		for p := 0; p < len(set); p++ {
			if q>>p&1 == 1 {
				fmt.Print(set[p], "  ")
			}
		}
		fmt.Println()
	}
}
```

时间复杂度为 O(n*2^n)，空间复杂度 O(n)。

### 迭代法

每次迭代，都是上一次迭代的结果+上次迭代结果中每个元素都加上当前迭代的元素+当前迭代的元素。

```go
func main() {
	set := []string{"1", "2", "4"}
	subs := []string{set[0]}
	for i := 1; i < len(set); i++ {
		ll := len(subs)
		for j := 0; j < ll; j++ {
			subs = append(subs, subs[j]+set[i])
		}
		subs = append(subs, set[i])
	}
	fmt.Println(subs)
}
```

时间复杂度为 O(2^n)，空间复杂度 2^(n-1)-1。

## 在有规律的二维数组中进行高效的数据查找

在一个二维数组中，每一行都按照从左到右递增的顺序排序，每一列都按照从上到下递增的顺序排序。请实现一个函数，输入这样的一个二维数组和一个整数，判断数组中是否含有该整数。

### 二分查找有序数组中是否存在某个元素

```go
func BinarySearchK(array []int, k int) bool {
	if len(array) == 0 {
		return false
	}
	start, end := 0, len(array)
	var mid int
	for start < end {
		mid = (start + end) / 2
		if array[mid] == k {
			return true
		} else if array[mid] < k {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return false
}
```

### 二分查找有序矩阵中是否存在某个元素

关键在于想到从二维数组的右上角遍历到左下角。

```go
func IsContainK(array2d [][]int, k int) bool {
	if len(array2d) == 0 {
		return false
	}
	rows, columns := len(array2d), len(array2d[0])
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

坐标轴上从左到右依次的点为 a[0]、a[1]、a[2]…a[n-1]，求满足 `a[j]-a[i]<=L&&a[j+1]-a[i]>L` 这两个条件的 j 与 i 中间的所有点个数中的最大值，即 `j-i+1` 最大。

太困了，似乎可以。

```go
func main() {
	array := []int{1, 3, 7, 8, 10, 11, 12, 13, 15, 16, 17, 19, 35}
	dst := 8
	x, y, z, final := 0, 1, 0, 1
	for _, val := range array {
		if val-array[x] > dst {
			x++
		}
		if y-x > z {
			z = y - x
			final = y
		}
		y++
	}
	fmt.Println(z, array[final-z:final])
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
```

## 根据规则构造新的数组

给定一个数组 `a[N]`，希望构造一个新的数组 `b[N]`，其中，`b[i]=a[0]*a[1]*…*a[N-1]/a[i]`。

1. 不允许使用除法
2. 要求具备 O(1) 空间复杂度和 O(n) 时间复杂度
3. 除遍历计数器与 `a[n]`、`b[n]` 外，不可以使用新的变量（包括栈临时变量、堆空间和全局静态变量等）

首先遍历一遍数组 a，在遍历的过程中对数组 b 进行赋值：`b[i]=a[i-1]*b[i-1]`，这样经过一次遍历后，数组 b 的值为 `b[i]=a[0]*a[1]*…*a[i-1]`。此时只需要将数组中的值 `b[i]` 再乘以` a[i+1]*a[i+2]*…a[N-1]`，实现方法为逆向遍历数组 a，把数组后半段值的乘积记录到 `b[0]` 中，通过 `b[i]` 与 `b[0]` 的乘积就可以得到满足题目要求的 `b[i]`，具体而言，执行 `b[i]=b[i]*b[0]`（首先执行的目的是为了保证在执行下面一个计算的时候，`b[0]` 中不包含与 `b[i]` 的乘积），接着记录数组后半段的乘积到 `b[0]` 中：`b[0]*=b[0]*a[i]`。

```go
func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b := make([]int, len(a))
	b[0] = 1// 其值顺序逆序都可用作储值
	for k := 0; k < len(a)-1; k++ {
		b[k+1] = a[k] * b[k]
	}
	b[0] = a[len(a)-1]// 其值顺序逆序都可用作储值
	for k := len(a) - 2; k > 0; k-- {
		b[k] *= b[0]
		b[0] *= a[k]
	}
	fmt.Println(b)
}
```

## 获取最好的矩阵链相乘方法

给定一个矩阵序列，找到最有效的方式将这些矩阵相乘在一起。给定表示矩阵链的数组 `p[]`，使得第 i 个矩阵 Ai 的维数为 `p[i-1]×p[i]`。

有4个大小为 40×20，20×30，30×10 和 10×30 的矩阵。假设这 4 个矩阵为 A、B、C 和 D，该函数的执行方法可以使执行乘法运算的次数最少。

4 个矩阵 A、B、C 和 D，可以有如下几种执行乘法的方法：

```
(ABC)D = (AB)(CD) = A(BCD)
```

这些方法的计算结果相同。但是，不同的方法需要执行乘法的次数是不相同的，因此效率也是不相同的。例如，假设 A 是 10×30 矩阵，B 是 30×5 矩阵，C 是 5×60 矩阵。那么，(AB)C 的执行乘法运算的次数为 (10×30×5)+(10×5×60) = 1500 + 3000 = 4500次。A(BC) 的执行乘法运算的次数为 (30×5×60)+(10×30×60)=9000+18000=27000次。

输入：p[]={40，20，30，10，30}
输出：26000

# 递归法

```go
func bestMatrixChainOrder(p []int, i, j int) int {
	if i == j {
		return 0
	}
	best := 1<<63 - 1
	for k := i; k < j; k++ {
		count := bestMatrixChainOrder(p, i, k) + bestMatrixChainOrder(p, k+1, j) + p[i-1]*p[k]*p[j]
		if count < best {
			best = count
		}
	}
	return best
}
```

矩阵运算不容易理解，另一种动态规划法暂略。

## 求解迷宫问题

给定一个大小为 N×N 的迷宫，一只老鼠需要从迷宫的左上角（对应矩阵的 `[0][0]`）走到迷宫的右下角（对应矩阵的 `[N-1][N-1]`），老鼠只能向两方向移动：向右或向下。在迷宫中，0 表示没有路（是死胡同），1 表示有路。

尝试可能的路径，遇到岔路口时保存其中一个方向，然后尝试走另一个方向，当碰到死胡同的时候，回溯到前一步，然后从前一步出发继续寻找可达的路径。

```go
// InitMatrix 初始化二维数组/矩阵
func InitMatrix(x, y int) [][]int {
	matrix := make([][]int, x)
	for i := range matrix {
		matrix[i] = make([]int, y)
	}
	return matrix
}
```

```go
func MazeSolver(matrix [][]int) (road [][2]int) {
	direct := [][2]int{{0, 0}}
	idxs := []int{0}
	length := len(matrix) - 1
	i, j := 0, 0
	for i < length || j < length {
		if matrix[i][j] == 1 {
			road = append(road, [2]int{i, j})
		}
		if i < length && j < length && matrix[i+1][j] == 1 && matrix[i][j+1] == 1 {
			direct = append(direct, [2]int{i, j + 1})
			idxs = append(idxs, len(road))
		}
		if i < length && matrix[i+1][j] == 1 {
			i++
		} else if j < length && matrix[i][j+1] == 1 {
			j++
		}
		if j < length && (i < length && matrix[i+1][j] == 0 || i == length && matrix[i][j+1] == 0) && matrix[i][j+1] == 0 {
			rear := direct[len(direct)-1]
			idx := idxs[len(idxs)-1]
			i, j = rear[0], rear[1]
			road = road[:idx]
			idxs = idxs[:len(idxs)-1]
			direct = direct[:len(direct)-1]
		}
	}
	return append(road, [2]int{length, length})
}
```

```go
func main() {
	matrix := [][]int{
		{1, 1, 1, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
	}
	// [[0 0] [1 0] [1 1] [1 2] [1 3] [2 3] [3 3]]
	road := MazeSolver(matrix)
	fmt.Println(road)

	roadMatrix := InitMatrix(len(matrix), len(matrix[0]))
	for _, v := range road {
		i, j := v[0], v[1]
		roadMatrix[i][j] = 1
	}
	PrintMatrix(roadMatrix)
}
```

## 从三个有序数组中找出它们的公共元素

给定以非递减顺序排序的三个数组，找出这三个数组中的所有公共元素。例如，给出下面三个数组：

```
a1 := []int{2, 5, 12, 20, 45, 85}
a2 := []int{16, 19, 20, 85, 200}
a3 := []int{3, 4, 15, 20, 39, 72, 85, 190}
```

那么这三个数组的公共元素为 {20, 85}。

```go
func findCommon(a1, a2, a3 []int) (com []int) {
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
```

## 求两个有序集合的交集

有两个有序的集合，集合中的每个元素都是一段范围，求其交集，例如集合 `{[4, 8]，[9, 13]}` 和 `{[6, 12]}` 的交集为 `{[6, 8]，[9, 12]}`。

```go
func findSetIntersection(s1, s2 [][]int) (com [][2]int) {
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

假设有一个中央调度机，有 n 个相同的任务需要调度到 m 台服务器上去执行，由于每台服务器的配置不一样，因此，服务器执行一个任务所花费的时间也不同。现在假设第 i 个服务器执行一个任务所花费的时间也不同。假设第 i 个服务器执行一个任务需要的时间为 t[i]。

例如：有 2 个执行机 a 与 b，执行一个任务分别需要 7min 和 10min，有 6 个任务待调度。如果平分这 6 个任务，即 a 与 b 各 3 个任务，则最短需要 30min 执行完所有。如果 a 分 4 个任务，b 分 2 个任务，则最短 28min 执行完。请设计调度算法，使得所有任务完成所需要的时间最短。

```go
func main() {
	costs := []int{7, 11}
	c1, c2 := 1, 1
	tasks := 100
	for tasks > 2 {
		if (c1+1)*costs[0] < (c2+1)*costs[1] {
			c1++
		} else {
			c2++
		}
		tasks--
	}
	fmt.Println(c1, c2, c1*costs[0], c2*costs[1])
}
```

## 对磁盘分区

有 N 个磁盘，每个磁盘大小为 `D[i]（i=0…N-1）`，现在要在这 N 个磁盘上"顺序分配" M 个分区，每个分区大小为 `P[j]（j=0…M-1）`，顺序分配的意思是：分配一个分区 P[j] 时，如果当前磁盘剩余空间足够，则在当前磁盘分配；如果不够，则尝试下一个磁盘，直到找到一个磁盘 `D[i+k]` 可以容纳该分区，分配下一个分区 `P[j+1]` 时，则从当前磁盘 `D[i+k]` 的剩余空间开始分配，不在使用 `D[i+k]` 之前磁盘未分配的空间，如果这 M 个分区不能在这 N 个磁盘完全分配，则认为分配失败。

判断给定 N 个磁盘（数组D）和 M 个分区（数组P），是否会出现分配失败的情况？举例：磁盘为 [120，120，120]，分区为 [60，60，80，20，80] 可分配，如果为 [60，80，80，20，80]，则分配失败。

```go
func main() {
	D := []int{120, 120, 120} // N
	//P := []int{60, 60, 80, 20, 80} // M
	//P := []int{60, 80, 80, 20, 80} // M
	P := []int{60, 60, 60, 60, 60, 60} // M
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
	if i == len(D) {
		fmt.Println("fail")
	} else {
		fmt.Println("success")
	}
}
```

```go

```


```go

```
