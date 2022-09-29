package math

// 求二进制数中 1 的个数

func getOneInBinary(n int) int {
	count := 0
	for i := 0; i < 32; i++ {
		m := 1 << i
		if n&m == m {
			count++
		}
	}
	return count
}
