package sort

import "fmt"

func Example_merge() {
	a := []int{1, 5, 9, 2, 3, 4}
	merge(a, 0, 2, 5)
	fmt.Println(a)

	// Output:
	// [1 2 3 4 5 9]
}

func ExampleMergeSort() {
	array := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	MergeSort(array, 0, len(array)-1)
	fmt.Println(array)

	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10]
}
