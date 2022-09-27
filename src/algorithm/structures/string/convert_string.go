package string

func ConvertToInteger(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	// 正负符号
	neg := 1
	pos := 0
	rs := []rune(s)

	switch rs[pos] {
	case '+':
		pos++
	case '-':
		pos++
		neg = -1
	}

	n := 0
	for ; pos < len(rs); pos++ {
		if rs[pos] < '0' || rs[pos] > '9' {
			return 0, false
		}
		n = n*10 + int(rs[pos]-'0')
	}

	return neg * n, true
}
