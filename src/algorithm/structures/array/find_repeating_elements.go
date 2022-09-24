package array

// 找出数组中唯一的重复元素

// 异或法
func xor() int {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 4}
	var x int
	for k, v := range array {
		x ^= k ^ v
	}
	return x
}
