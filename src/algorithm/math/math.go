package math

// Pow 快速幂运算 即求 x 的 n 次方
func Pow(x, n int) int {
	ret := 1
	for n != 0 {
		if n%2 != 0 {
			ret *= x
		}
		n /= 2
		x *= x
	}
	return ret
}

// MaxN 求多个数中的最大值
func MaxN(args ...int) int {
	nMax := -1 << 63
	for _, v := range args {
		if v > nMax {
			nMax = v
		}
	}
	return nMax
}

// MinN 求多个数中的最小值
func MinN(args ...int) int {
	nMin := 1<<63 - 1
	for _, v := range args {
		if v < nMin {
			nMin = v
		}
	}
	return nMin
}

// Max 求最大值
func Max(n1, n2 int) int {
	if n1 < n2 {
		return n2
	}
	return n1
}

// Min 求最小值
func Min(n1, n2 int) int {
	if n1 > n2 {
		return n2
	}
	return n1
}

// Abs 求绝对值
func Abs(n int) int {
	if n >= 0 {
		return n
	}
	// return -1
	return ^(n - 1) // 实际性能与取负数是差不多的
}

// Opposite 求相反数
func Opposite(n int) int {
	return ^(n - 1)
}
