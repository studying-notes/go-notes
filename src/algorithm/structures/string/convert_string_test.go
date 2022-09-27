package string

import "fmt"

func ExampleConvertToInteger() {
	fmt.Println(ConvertToInteger("-345"))
	fmt.Println(ConvertToInteger("345"))
	fmt.Println(ConvertToInteger("+345"))
	fmt.Println(ConvertToInteger("++345"))

	// Output:
	// -345 true
	// 345 true
	// 345 true
	// 0 false
}
