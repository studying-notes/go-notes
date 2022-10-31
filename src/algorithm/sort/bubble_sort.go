package sort

// BubbleSort 冒泡排序
func BubbleSort(array []int, length int) {
	for i := 0; i < length-1; i++ {
		for j := length - 1; j > i; j-- {
			if array[j] < array[j-1] {
				swap(array, j, j-1)
			}
		}
	}
}
