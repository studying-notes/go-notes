package bit_operation

import "fmt"

func ExampleDivide() {
	fmt.Println(Divide(12, 3))
	fmt.Println(Divide(13, 3))
	fmt.Println(Divide(11, -3))
	fmt.Println(Divide(12, -3))
	fmt.Println(Divide(13, -3))

	// Output:
	// 4
	// 4
	// -3
	// -4
	// -4
}
