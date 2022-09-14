package stack

import "fmt"

func ExampleReverseStack() {
	s := Stack{1, 2, 3, 4, 5, 6, 7, 8}
	ReverseStack(&s)
	fmt.Println(s)

	// Output:
	// [8 7 6 5 4 3 2 1]
}
