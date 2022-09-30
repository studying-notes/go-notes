package string

import "sort"

type LongestWord struct {
	words []string
}

func NewLongestWord(words []string) *LongestWord {
	w := &LongestWord{words: words}
	w.sort()
	return w
}

func (w *LongestWord) sort() {
	sort.Slice(w.words, func(i, j int) bool {
		return len(w.words[i]) > len(w.words[j])
	})
}

// 判断字符串是否在字符串数组中
func (w *LongestWord) find(word string) bool {
	for i := range w.words {
		if w.words[i] == word {
			return true
		}
	}
	return false
}

func (w *LongestWord) isContains(word string, length int) bool {
	if word == "" {
		return true
	}

	ll := len(word)

	// 循环取前缀
	for i := 1; i <= ll; i++ {
		if i == length {
			return false
		}

		if w.find(word[0:i]) {
			if w.isContains(word[i:], length) {
				return true
			}
		}
	}

	return false
}

func (w *LongestWord) GetLongestWord() string {
	for i := range w.words {
		if w.isContains(w.words[i], len(w.words[i])) {
			return w.words[i]
		}
	}
	return ""
}
