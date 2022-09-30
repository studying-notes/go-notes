package string

import "algorithm/structures/array"

func getLongestRepeatingSubstring(s string) string {
	length := len(s)
	maxLength := 1
	maxEnd := 0

	matrix := array.InitMatrix(length+1, length+1)

	for i := 1; i <= length; i++ {
		for j := 1; j <= length; j++ {
			if s[i-1] == s[j-1] && i != j {
				matrix[i][j] = matrix[i-1][j-1] + 1
				if matrix[i][j] > maxLength {
					maxLength = matrix[i][j]
					maxEnd = i
				}
			}
		}
	}

	// array.PrintMatrix(matrix)

	return s[maxEnd-maxLength : maxEnd]
}
