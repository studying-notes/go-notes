package sort

import "fmt"

func ExampleHillSort() {
	array := []int{26, 53, 67, 48, 57, 13, 48, 32, 60, 50}
	fmt.Println(array)
	HillSort(array)

	// Output:
	// [26 53 67 48 57 13 48 32 60 50]
	// [13 48 32 48 50 26 53 67 60 57]
	// [13 26 32 48 50 48 53 57 60 67]
	// [13 26 32 48 48 50 53 57 60 67]
}
