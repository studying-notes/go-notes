package sort

import "fmt"

// HillSort 希尔排序
func HillSort(array []int) {
	length := len(array)

	for h := length / 2; h > 0; h /= 2 {
		for i := h; i < length; i++ {
			for j := i - h; j >= 0 && array[j] > array[j+h]; j -= h {
				array[j], array[j+h] = array[j+h], array[j]
			}
		}

		fmt.Println(array)
	}
}
