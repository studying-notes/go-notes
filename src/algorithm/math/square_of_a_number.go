package math

// 判断一个自然数是否是某个数的平方

func isSquare1(n int) (int, bool) {
	for i := 1; i < n; i++ {
		m := i * i
		if m == n {
			return i, true
		} else if m > n {
			return 0, false
		}
	}
	return 0, false
}

func isSquare2(n int) (int, bool) {
	left, right := 0, n

	for left < right {
		mid := (left + right) / 2
		m := mid * mid
		if m < n {
			left = mid + 1
		} else if m > n {
			right = mid - 1
		} else {
			return mid, true
		}
	}
	return 0, false
}

func isSquare3(n int) bool {
	minus := 1

	for n > 0 {
		n -= minus
		if n == 0 {
			return true
		} else if n < 0 {
			return false
		} else {
			minus += 2
		}
	}

	return false
}
