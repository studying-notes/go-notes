package string

func Pmt(s string) []int {
	pmt := make([]int, len(s))
	pmt[0] = 0
	for i, j := 1, 0; i < len(s); {
		if s[i] == s[j] {
			pmt[i] = pmt[i-1] + 1
			i, j = i+1, j+1
		} else if j == 0 {
			i++
		} else {
			// 从两边扫描对称值
			x, y := 0, i
			for s[x] == s[y] {
				x, y = x+1, y-1
			}
			pmt[i] = x
			i, j = i+1, x
		}
	}
	return pmt
}

// Next 求子字符串的 next 数组
func Next(s string) []int {
	next := make([]int, len(s)+1)
	next[0] = -1
	for i, j := 0, -1; i < len(s); {
		if j == -1 || s[i] == s[j] {
			i, j = i+1, j+1
			next[i] = j
		} else {
			j = next[j]
		}
	}
	return next
}

func IsSubStringKMP(s, sub string) (start int, ret bool) {
	if len(s) < len(sub) {
		return -1, false
	}
	next := Next(sub)
	i, j := 0, 0
	for i < len(s) && j < len(sub) {
		if j == -1 || s[i] == sub[j] {
			j++
		} else {
			j = next[j]
		}
		i++
	}
	if j == len(sub) {
		return i - len(sub), true
	}
	return -1, false
}
