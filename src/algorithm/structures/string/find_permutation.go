package string

// 保存结果列表
var result []string

func Permutation(s []rune, start int) {
	if s == nil {
		return
	}

	if start == len(s)-1 {
		result = append(result, string(s))
	} else {
		for i := start; i < len(s); i++ {
			s[start], s[i] = s[i], s[start]
			Permutation(s, start+1)
			s[start], s[i] = s[i], s[start]
		}
	}
}
