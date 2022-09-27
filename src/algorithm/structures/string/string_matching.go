package string

// IsSubString 遍历字符串匹配
func IsSubString(s, sub string) (start int, ret bool) {
	if len(s) < len(sub) {
		return -1, false
	}
	i, j := 0, 0
	for i < len(s) && j < len(sub) {
		for i < len(s) && s[i] != sub[j] {
			i++
		}
		if i == len(s) {
			return -1, false
		}
		start = i
		for i < len(s) && j < len(sub) && s[i] == sub[j] {
			i, j = i+1, j+1
		}
		if j == len(sub) {
			break
		}
		if i == len(s) {
			return -1, false
		}
		i, j = start+1, 0
	}
	return start, true
}
