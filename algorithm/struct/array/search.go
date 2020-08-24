package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array2d := [][]int{
		{1, 2, 8, 10},
		{2, 4, 9, 12},
		{4, 7, 10, 13},
		{6, 8, 11, 15},
	}
	PrintMatrix(array2d)
	fmt.Println(IsContainK(array2d, 9))

	array := Range(10)
	fmt.Println(BinarySearchK(array, 5))
}

// BinarySearchK 二分查找有序数组中是否存在某个元素
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

// IsContainK 二分查找有序矩阵中是否存在某个元素
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
