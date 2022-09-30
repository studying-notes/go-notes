package string

func getLengthOfLongestIncreasingSubstring(s string) string {
	maxEnd := 0
	maxLength := 1
	currentLength := 1

	for i := 1; i < len(s); i++ {
		if s[i] > s[i-1] {
			currentLength++
		} else {
			if currentLength > maxLength {
				maxLength = currentLength
				maxEnd = i
			}
			currentLength = 1
		}
	}

	return s[maxEnd-maxLength : maxEnd]
}
