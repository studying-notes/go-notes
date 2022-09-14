package stack

import "fmt"

func ExampleNewExtStack() {
	// s := Stack{9, 8, 1, 6, 5, 10, 3, 12, 11, 4}
	// BubbleSort(&s)
	s := NewExtStack(9, 8, 1, 6, 5, 0, 3, 12, 11, 4)
	fmt.Println(s.Min())

	// Output:
	// 0
}
