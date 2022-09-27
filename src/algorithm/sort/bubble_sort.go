package sort

// BubbleSort 冒泡排序
func BubbleSort(array []int) {
	length := len(array)
	for i := 1; i < length; i++ {
		for j := 0; j < length-i; j++ {
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}
