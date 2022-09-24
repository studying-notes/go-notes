package array

func findOnce(array []int) int {
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

	return x ^ s
}
