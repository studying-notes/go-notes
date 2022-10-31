package sort

// InsertionSort 插入排序
func InsertionSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		for j := i + 1; j > 0; j-- {
			if array[j] < array[j-1] {
				break // 已经有序则中断
			}
			// 从后往前遍历有序序列，后移比它大的元素
			swap(array, j, j-1)
		}
	}
}
