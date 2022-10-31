package sort

// HillSort 希尔排序
func HillSort(array []int, length int) {
	for step := length / 2; step > 0; step /= 2 {
		for i := step; i < length; i += step {
			for j := i; j > step-1; j -= step {
				if array[j] >= array[j-step] {
					break
				}
				swap(array, j, j-step)
			}
		}
	}
}
