package main

import (
	"fmt"
)

func main() {
	array := []int{10, 2, 3, 4, 15, 6, 7, 8, 9}
	fmt.Println(GetMaxAndMin(array, 0, len(array)-1))
}

// 分治法查找数组中元素的最大值和最小值
// 只是求最大值或最小值就直接依次比较即可
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
