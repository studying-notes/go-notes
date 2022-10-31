package sort

import "fmt"

func ExampleBubbleSort() {
	array := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	BubbleSort(array, len(array))
	fmt.Println(array)

	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10]
}
