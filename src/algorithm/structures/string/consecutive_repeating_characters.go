package string

func getMaxRepeatCountLoop(s string) (int, string) {
	maxCount := 0
	maxCharIndex := 1
	currentCount := 1
	chars := []rune(s)

	for i := 1; i < len(chars); i++ {
		// fmt.Printf("%s %s %d %d\n", string(chars[i-1]), string(chars[i]), currentCount, maxCount)
		if chars[i] == chars[i-1] {
			currentCount++
		} else {
			if currentCount > maxCount {
				maxCount = currentCount
				maxCharIndex = i - 1
			}
			currentCount = 1
		}
		// fmt.Printf("%s %s %d %d\n", string(chars[i-1]), string(chars[i]), currentCount, maxCount)
	}

	return maxCount, string(chars[maxCharIndex])
}

// 用递归的方式实现一个求字符串中连续出现相同字符的最大值
func getMaxRepeatCountRecursion(s string, start, currentCount, maxCount int) int {
	chars := []rune(s)

	if start < len(chars) {
		if chars[start] == chars[start-1] {
			currentCount++
		} else {
			if currentCount > maxCount {
				maxCount = currentCount
			}
			currentCount = 1
		}
		return getMaxRepeatCountRecursion(s, start+1, currentCount, maxCount)
	}

	return maxCount
}
