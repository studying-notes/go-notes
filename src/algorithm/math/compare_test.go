package math

import "fmt"

func Example_compare() {
	fmt.Println(compare(3, 2))
	fmt.Println(compare(1, 2))
	fmt.Println(compare(-1, 2))
	fmt.Println(compare(3, -2))
	fmt.Println(compare(-3, -2))

	// Output:
	// true
	// false
	// false
	// true
	// false
}
