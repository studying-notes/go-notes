package sort

// MergeSort 归并排序
func MergeSort(array []int, begin, end int) {
	// begin 表示开始索引 end 表示结束索引，可以取到
	if begin < end {
		mid := (begin + end) / 2
		MergeSort(array, begin, mid)
		MergeSort(array, mid+1, end)
		merge(array, begin, mid, end)
	}
}

// 合并两个有序数组
func merge(array []int, begin, mid, end int) {
	leftLength := mid - begin + 1
	rightLength := end - mid

	left := make([]int, leftLength)
	right := make([]int, rightLength)

	for i := 0; i < leftLength; i++ {
		left[i] = array[begin+i]
	}

	for j := 0; j < rightLength; j++ {
		right[j] = array[mid+j+1]
	}

	var i, j int
	for i < leftLength && j < rightLength {
		if left[i] < right[j] {
			array[begin+i+j] = left[i]
			i++
		} else {
			array[begin+i+j] = right[j]
			j++
		}
	}

	for i < leftLength {
		array[begin+i+j] = left[i]
		i++
	}

	for j < rightLength {
		array[begin+i+j] = right[j]
		j++
	}
}
