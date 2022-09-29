package bit_operation

import "fmt"

func ExampleAdd() {
	fmt.Println(Add(1, 2))
	fmt.Println(Add(15, 25))

	// Output:
	// 3
	// 40
}

func ExampleAddPlusPlus() {
	fmt.Println(AddPlusPlus(1, 2))
	fmt.Println(AddPlusPlus(15, 25))

	// Output:
	// 3
	// 40
}
