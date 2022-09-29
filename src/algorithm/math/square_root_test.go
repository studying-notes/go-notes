package math

import "fmt"

func Example_getSquareRoot() {
	fmt.Println(getSquareRoot(15, 0.00001))
	fmt.Println(getSquareRoot(16, 0.00001))
	fmt.Println(getSquareRoot(4, 0.00001))

	// Output:
	// 1
}
