package string

// 实现单词反转

// ReverseStringXORInPlace 字符串原地逆序
func ReverseStringXORInPlace(rs []rune) {
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i] ^= rs[j] // i ^ j
		rs[j] ^= rs[i] // i ^ j ^ j = i
		rs[i] ^= rs[j] // i ^ j ^ i = j
	}
}

func ReverseWords(words string) string {
	rs := []rune(words)

	ReverseStringXORInPlace(rs)

	begin := 0
	for i := range rs {
		if rs[i] == ' ' {
			ReverseStringXORInPlace(rs[begin:i])
			begin = i + 1
		}
	}

	ReverseStringXORInPlace(rs[begin:])

	return string(rs)
}
