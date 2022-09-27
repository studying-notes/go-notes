package sort

// MergeSort 归并排序
func MergeSort(array []int, begin, end int) {
	if begin < end {
		q := (begin + end) / 2
		MergeSort(array, begin, q)
		MergeSort(array, q+1, end)
		merge(array, begin, q, end)
	}
}

// 合并两个有序数组
func merge(array []int, begin, middle, end int) {
	n1 := middle - begin + 1
	n2 := end - middle

	left := make([]int, n1)
	right := make([]int, n2)

	for i, k := 0, begin; i < n1; i, k = i+1, k+1 {
		left[i] = array[k]
	}

	for i, k := 0, middle+1; i < n2; i, k = i+1, k+1 {
		right[i] = array[k]
	}

	var i, j int
	for i < n1 && j < n2 {
		if left[i] < right[j] {
			array[begin+i+j] = left[i]
			i++
		} else {
			array[begin+i+j] = right[j]
			j++
		}
	}

	for i < n1 {
		array[begin+i+j] = left[i]
		i++
	}

	for j < n2 {
		array[begin+i+j] = right[j]
		j++
	}
}
