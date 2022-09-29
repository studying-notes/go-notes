package math

// 判断 1024! 末尾有多少个 0

func getFactorialZero(n int) int {
	dividend := 5
	count := 0

	for dividend < n {
		count += n / dividend
		dividend *= 5
	}

	return count
}
