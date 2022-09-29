package math

func getPower(d, n int) (m int) {
	if n == 0 {
		return 1
	} else if d == 0 {
		return 0
	} else if n == 1 {
		return d
	}

	m = getPower(d, Abs(n/2))

	if n > 0 {
		if n%2 == 1 {
			return m * m * d
		}
		return m * m
	}

	if n%2 == 1 {
		return 1 / (m * m * d)
	}
	return 1 / (m * m)
}
