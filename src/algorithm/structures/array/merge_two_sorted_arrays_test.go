package array

import "fmt"

func ExampleMergeSortedArrays() {
	a := []int{1, 5, 9}
	b := []int{2, 3, 4}

	fmt.Println(MergeSortedArrays(a, b))

	// Output:
	// [1 2 3 4 5 9]
}

func ExampleMergeSortedArraysInPlace() {
	a := []int{1, 5, 9, 2, 3, 4}

	MergeSortedArraysInPlace(a, 0, 2, 5)

	fmt.Println(a)

	// Output:
	// [1 2 3 4 5 9]
}
