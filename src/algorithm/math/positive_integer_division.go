package math

func integerDivision1(m, n int) (quotient int, remainder int) {
	for remainder = m; remainder > n; remainder -= n {
		quotient++
	}

	return
}

func integerDivision2(m, n int) (quotient int, remainder int) {
	for m >= n {
		multi := 1
		for multi*n <= m>>1 {
			multi <<= 1
		}
		quotient += multi
		m -= multi * n
	}

	return quotient, m
}
