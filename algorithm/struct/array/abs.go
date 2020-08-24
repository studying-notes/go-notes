package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/utils"
)

func main() {
	array := []int{-10, -5, -2, 0, 1, 7, 15, 20}
	fmt.Println(findMinAbs(array))
	fmt.Println(findMinBinary(array, 0, 5))
}

// 二分查找法
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

// 求升序数组中绝对值最小的数
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
