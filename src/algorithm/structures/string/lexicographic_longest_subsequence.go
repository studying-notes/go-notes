package string

import "bytes"

func getLexicographicLongestSubstring(m string) string {
	var i, j int
	var n bytes.Buffer

	for i < len(m) {
		maxIndex, maxValue := i, m[i]
		for j = i + 1; j < len(m); j++ {
			if m[j] > maxValue {
				maxIndex, maxValue = j, m[j]
			}
		}
		i = maxIndex + 1
		n.WriteByte(maxValue)
	}

	return n.String()
}
