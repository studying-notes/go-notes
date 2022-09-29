package math

func isPowerOf21(n int) bool {
	if n < 1 {
		return false
	}

	for i := 1; i < n; i++ {
		m := 1 << i
		if m == n {
			return true
		} else if m > n {
			return false
		}
	}
	return false
}

func isPowerOf22(n int) bool {
	if n < 1 {
		return false
	}

	return n&(n-1) == 0
}
