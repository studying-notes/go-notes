package sort

import "fmt"

func ExampleSelectionSort() {
	array := []int{10, 3, 8, 7, 0, 5, 4, 9, 2, 1, 6}
	SelectionSort(array, len(array))
	fmt.Println(array)

	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10]
}
