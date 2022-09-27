package array

func MergeSortedArraysInPlace(array []int, u, v, w int) {
	// [1 5 9 ] [ 2 3 4]
}

// MergeSortedArrays 允许申请新的数组空间
func MergeSortedArrays(a, b []int) (result []int) {
	aLeft, bLeft := 0, 0
	aRight, bRight := len(a), len(b)
	result = make([]int, aRight+bRight)

	for aLeft < aRight && bLeft < bRight {
		if a[aLeft] > b[bLeft] {
			result[aLeft+bLeft] = b[bLeft]
			bLeft++
		} else {
			result[aLeft+bLeft] = a[aLeft]
			aLeft++
		}
	}

	// 处理多余元素
	for aLeft < aRight {
		result[aLeft+bLeft] = a[aLeft]
		aLeft++
	}

	for bLeft < bRight {
		result[aLeft+bLeft] = b[bLeft]
		bLeft++
	}

	return
}
