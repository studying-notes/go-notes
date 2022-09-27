package string

type LongestWord struct {
}

// 判断字符串是否在字符串数组中
func (w LongestWord) find(words []string, word string) bool {
	for i := range words {
		if words[i] == word {
			return true
		}
	}
	return false
}

func (w LongestWord) isContains(words []string, word string, length int) bool {
	return true
}
