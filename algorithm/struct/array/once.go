package main

import "fmt"

func main() {
	//array := []int{1, 2, 4, 5, 6, 4, 2}
	array := []int{6, 3, 4, 5, 9, 4, 3}
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
