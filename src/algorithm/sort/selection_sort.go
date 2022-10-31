package sort

// SelectionSort 选择排序
func SelectionSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		k := i // 记录最小值索引
		for j := i + 1; j < length; j++ {
			if array[j] < array[k] {
				k = j
			}
		}
		swap(array, i, k) // 交换最小值
	}
}
