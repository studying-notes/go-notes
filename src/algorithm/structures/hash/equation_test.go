package hash

import "fmt"

func ExampleFindEquation() {
	list := []int{3, 4, 7, 10, 2, 9, 8}
	fmt.Println(FindEquation(list))

	// Output:
	// [3 8] [4 7]
}
