package sort

// SelectionSort 选择排序
func SelectionSort(array []int) {
	length := len(array)
	for i := 0; i < length; i++ {
		// 记录最小值及其索引
		m, n := array[i], i
		for j := i + 1; j < length; j++ {
			if array[j] < m {
				m, n = array[j], j
			}
		}
		// 交换最小值
		array[i], array[n] = array[n], array[i]
	}
}
