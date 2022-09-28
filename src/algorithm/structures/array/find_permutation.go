package array

// 复制切片
func deepcopy(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

// 比较数组的大小
func moreThan(left, right []int) bool {
	leftLength, rightLength := len(left), len(right)

	if leftLength != rightLength {
		return leftLength > rightLength
	}

	for i := 0; i < leftLength; i++ {
		if left[i] != right[i] {
			return left[i] > right[i]
		}
	}

	return false
}

type Permutation struct {
	result [][]int
	array  []int
}

func NewPermutation(n int) *Permutation {
	array := make([]int, n)
	for i := range array {
		array[i] = i + 1
	}

	return &Permutation{array: array}
}

func (p *Permutation) Perform(start int) {
	if start == len(p.array)-1 {
		p.result = append(p.result, deepcopy(p.array))
	} else {
		for i := start; i < len(p.array); i++ {
			p.array[start], p.array[i] = p.array[i], p.array[start]
			p.Perform(start + 1)
			p.array[start], p.array[i] = p.array[i], p.array[start]
		}
	}
}

func (p *Permutation) Sort() {
	// 已经基本有序，所以用一遍冒泡排序即可
	for i := 1; i < len(p.result); i++ {
		if moreThan(p.result[i-1], p.result[i]) {
			p.result[i-1], p.result[i] = p.result[i], p.result[i-1]
		}
	}
}
