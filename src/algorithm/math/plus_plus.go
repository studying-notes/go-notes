package math

func add(x, y int) int {
	var i, j, z int

	// 选择比较小的数进行循环

	if x > y {
		z, j = x, y
	} else {
		z, j = y, x
	}

	for i = 0; i < j; i++ {
		z++
	}

	return z
}
