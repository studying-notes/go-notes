package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/utils"
)

func main() {
	array := []int{1, -2, 4, 8, -4, 7, -1, -5}
	//array := []int{1, -2}

	fmt.Println(MaxSubArraySumDyn(array))
	fmt.Println(findSubArray(array))
}

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
