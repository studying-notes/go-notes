package string

// 对字符串进行反转

func ReverseString(s string) string {
	rs := []rune(s)

	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}

	return string(rs)
}

func ReverseStringXOR(s string) string {
	rs := []rune(s)

	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i] ^= rs[j] // i ^ j
		rs[j] ^= rs[i] // i ^ j ^ j = i
		rs[i] ^= rs[j] // i ^ j ^ i = j
	}

	return string(rs)
}
