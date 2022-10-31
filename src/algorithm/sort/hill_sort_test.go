package sort

import "fmt"

func ExampleHillSort() {
	array := []int{26, 53, 67, 48, 57, 13, 48, 32, 60, 50}
	HillSort(array, len(array))
	fmt.Println(array)

	// Output:
	// [13 26 32 48 48 50 53 57 60 67]
}
