package string

import "algorithm/structures/array"

type Palindrome struct {
	begin, length int
}

func (p *Palindrome) getLongestPalindrome(s string) string {
	if s == "" {
		return ""
	}

	ll := len(s)
	p.begin = 0
	p.length = 1

	// 申请额外的存储空间记录查找的历史信息
	history := array.InitMatrix(ll, ll)

	// 初始化长度为1的回文字符串信息
	for i := 0; i < ll; i++ {
		history[i][i] = 1
	}

	// 初始化长度为2的回文字符串信息
	for i := 0; i < ll-1; i++ {
		if s[i] == s[i+1] {
			history[i][i+1] = 1
			p.begin = i
			p.length = 2
		}
	}

	// 查找从长度为3开始的回文字符串
	for pl := 3; pl <= ll; pl++ {
		for i := 0; i < ll-pl; i++ {
			j := i + pl - 1
			if s[i] == s[j] && history[i+1][j-1] == 1 {
				history[i][j] = 1
				p.begin = i
				p.length = pl
			}
		}
	}

	// array.PrintMatrix(history)

	return s[p.begin : p.begin+p.length]
}

// 对字符串str，以c1和c2为中心向两侧扩展寻找回文子串
func (p *Palindrome) expandBothSide(s string, c1, c2 int) {
	ll := len(s)

	for c1 >= 0 && c2 < ll && s[c1] == s[c2] {
		c1--
		c2++
	}

	beginTemp := c1 + 1
	lengthTemp := c2 - beginTemp
	if lengthTemp > p.length {
		p.length = lengthTemp
		p.begin = beginTemp
	}
}

// 找出字符串最长的回文子串
func (p *Palindrome) getLongestPalindromeExpand(s string) string {
	if s == "" {
		return ""
	}

	for i := range s {
		// 找回文字符串长度为奇数的情况（从第i个字符向两边扩展）
		p.expandBothSide(s, i, i)
		// 找回文字符串长度为偶数的情况（从第ⅰ和计1两个字符字符向两边扩展）
		p.expandBothSide(s, i, i+1)
	}

	return s[p.begin : p.begin+p.length]
}
