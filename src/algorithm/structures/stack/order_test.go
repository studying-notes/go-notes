package stack

import "fmt"

func ExampleIsPopSerial() {
	pushOrder := []int{1, 2, 3, 4, 5}
	popOrder := []int{3, 2, 5, 4, 1}
	// popOrder := []int{5, 3, 4, 1, 2}
	fmt.Println(IsPopSerial(pushOrder, popOrder))

	// Output:
	// true
}
