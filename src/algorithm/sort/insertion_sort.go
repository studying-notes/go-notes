package sort

// InsertionSort 插入排序
func InsertionSort(array []int) {
	length := len(array)
	for i := 0; i < length-1; i++ {
		for j := i + 1; j > 0 && array[j] < array[j-1]; j-- {
			// 从后往前遍历有序序列，后移比它大的元素
			array[j], array[j-1] = array[j-1], array[j]
		}
	}
}
