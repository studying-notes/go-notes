package math

func getSquareRoot(n, e float64) float64 {
	m := n
	l := 1.0

	for m-l > e {
		m = (m + l) / 2
		l = n / m
	}

	return m
}
