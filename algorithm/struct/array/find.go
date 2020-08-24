package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

// 找出数组中出现奇数次的数
func odd() {
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

// 寻找重复值

func tag() {
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

// 异或法
func xor() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
	var x int
	for k, v := range array {
		x ^= k ^ v
	}
	fmt.Println(x)
}

func add() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}

	// 但如果累加的数值巨大时，就很有可能溢出了
	var x, y int
	for k, v := range array {
		x += k
		y += v
	}
	fmt.Println(y - x)

	// 改进 - 累计差值求和
	var s int
	for k, v := range array { // 从 0 开始个数正好
		s += v - k
	}
	fmt.Println(s)
}
