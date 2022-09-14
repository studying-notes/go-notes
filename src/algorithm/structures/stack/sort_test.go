package stack

import "fmt"

func ExampleSortStack() {
	s := Stack{9, 8, 1, 6, 5, 10, 3, 12, 11, 4}
	SortStack(&s)
	fmt.Println(s)

	// Output:
	// [1 3 4 5 6 8 9 10 11 12]
}
