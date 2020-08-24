package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	//array := []int{6, 7, 8, 9, 1, 2, 3, 4, 5}
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//fmt.Println(BinarySearch(array, 0, len(array)-1))

	fmt.Println(SpinArrayAppend(array, 5))

	SpinArrayReverse(array, 5)
	fmt.Println(array)

	fmt.Println(Range(100))
}

// 旋转数组/循环移位
func SpinArrayReverse(array []int, idx int) {
	ReverseArray(array[:idx])
	ReverseArray(array[idx:])
	ReverseArray(array)
}

func SpinArrayAppend(array []int, idx int) []int {
	return append(array[idx:], array[:idx]...)
}

// 找出旋转数组的最小元素
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
