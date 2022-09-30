package string

import "fmt"

func Example_getStringEditDistance() {
	fmt.Println(getStringEditDistance("bcd", "abc"))
	fmt.Println(getStringEditDistance("bcd", "abcd"))

	// Output:
	// 2
	// 1
}
