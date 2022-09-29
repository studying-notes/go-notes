package bit_operation

import "fmt"

func ExampleMultiply() {
	fmt.Println(Multiply(11, 14))
	fmt.Println(Multiply(3, 0))
	fmt.Println(Multiply(0, 4))
	fmt.Println(Multiply(3, -4))

	// Output:
	// 154
	// 0
	// 0
	// -12
}
