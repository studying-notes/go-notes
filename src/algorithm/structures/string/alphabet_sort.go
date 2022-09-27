package string

type AlphabetComparer struct {
	m map[byte]int
}

func NewAlphabetComparer(sequence []byte) *AlphabetComparer {
	m := make(map[byte]int, len(sequence))
	for i, j := range sequence {
		m[j] = i
	}
	return &AlphabetComparer{m: m}
}

// MoreThan 字符串比较算法
func (c *AlphabetComparer) MoreThan(left, right string) bool {
	leftLength := len(left)
	rightLength := len(right)

	for i := 0; i < leftLength && i < rightLength; i++ {
		difference := c.m[left[i]] - c.m[right[i]]
		if difference != 0 {
			return difference > 0
		}
	}

	return false
}

// SortStringSequences 用插入排序方式
func (c *AlphabetComparer) SortStringSequences(sequence []string) []string {
	var i, j int
	var ll = len(sequence)

	for i = 1; i < ll; i++ {
		v := sequence[i]
		for j = i - 1; j >= 0; j-- {
			if c.MoreThan(sequence[j], v) {
				sequence[j+1] = sequence[j]
			} else {
				break
			}
		}
		sequence[j+1] = v
	}

	return sequence
}
