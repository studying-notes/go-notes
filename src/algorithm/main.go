package main

import (
	"algorithm/structures/set"
	"fmt"
)

func main() {
	array := []int{1, 9, 3, 4, 5, 6, 7, 8, 9, 4, 5, 6, 7, 2, 1}
	s := set.NewSet()
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
