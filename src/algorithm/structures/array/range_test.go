package array

import "fmt"

func ExampleRange() {
	fmt.Println(Range(5))
	fmt.Println(Range(1, 5))
	fmt.Println(Range(1, 10, 2))

	// Output:
	// [0 1 2 3 4]
	// [1 2 3 4]
	// [1 3 5 7 9]
}
