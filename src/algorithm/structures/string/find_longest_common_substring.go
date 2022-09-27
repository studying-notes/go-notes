package string

import (
	"algorithm/structures/array"
)

// 求两个字符串的最长公共子串

// FindLongestCommonSubstring 动态规划法
func FindLongestCommonSubstring(s1, s2 []rune) []rune {
	l1, l2 := len(s1), len(s2)
	var end, length int // 记录公共字符串结束位置和长度

	// 申请新的空间来记录公共字串长度信息
	matrix := array.InitMatrix(l1+1, l2+1)

	for i := 0; i < len(s1); i++ {
		for j := 0; j < len(s2); j++ {
			if s1[i] == s2[j] {
				matrix[i+1][j+1] = matrix[i][j] + 1
				if matrix[i+1][j+1] > length {
					length = matrix[i+1][j+1]
					end = i + 1 // 因为是开区间
				}
			}
		}
	}

	array.PrintMatrix(matrix)
	return s1[end-length : end]
}

// LongestSubString 滑动比较法
func LongestSubString(s1, s2 []rune) (sub []rune) {
	p, q := 0, 0 // 记录最长子序列的开始和结束
	for k := 0; k < len(s2); k++ {
		i, j, x, y := 0, 0, 0, k         // 记录子序列的开始和结束
		for x < len(s1) && y < len(s2) { // 字符序列当前下标
			if s1[x] != s2[y] {
				if i != j { // 说明之前有过相同元素
					y = k // 必须重置起始位置
				} else {
					x++
				}
				i, j = x, x
			} else {
				j, x, y = j+1, x+1, y+1
			}
			// 判断当前子序列是否更长
			if j-i > q-p {
				q, p = j, i
			}
		}
	}
	return s1[p:q]
}
