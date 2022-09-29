package math

import "fmt"

func Example_getOneInBinary() {
	fmt.Println(getOneInBinary(1))
	fmt.Println(getOneInBinary(2))
	fmt.Println(getOneInBinary(3))
	fmt.Println(getOneInBinary(7))

	// Output:
	// 1
	// 1
	// 2
	// 3
}
