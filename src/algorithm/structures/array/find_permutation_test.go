package array

import (
	"fmt"
)

func Example_moreThan() {
	left := []int{1, 3, 3}
	right := []int{1, 2, 3}
	fmt.Println(moreThan(left, right))

	// Output:
	// true
}

func ExamplePermutation() {
	p := NewPermutation(3)
	p.Perform(0)
	fmt.Println(p.result)
	p.Sort()
	fmt.Println(p.result)

	// Output:
	// [[1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 2 1] [3 1 2]]
	// [[1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 1 2] [3 2 1]]
}
