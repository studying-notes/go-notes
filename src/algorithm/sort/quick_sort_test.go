package sort

import "fmt"

func ExamplePartition() {
	array := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	Partition(array, 0, len(array)-1)
	fmt.Println(array)

	// Output:
	// [0 9 8 7 6 5 4 3 2 1 10]
}

func ExampleQuickSort() {
	array := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	QuickSort(array, 0, len(array)-1)
	fmt.Println(array)

	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10]
}
