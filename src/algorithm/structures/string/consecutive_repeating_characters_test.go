package string

import "fmt"

func Example_getMaxRepeatCountLoop() {
	fmt.Println(getMaxRepeatCountLoop("aabbcccdd"))
	fmt.Println(getMaxRepeatCountLoop("aaabbcc"))

	// Output:
	// 3 c
	// 3 a
}

func Example_getMaxRepeatCountRecursion() {
	fmt.Println(getMaxRepeatCountRecursion("aabbcccdd", 1, 1, 0))
	fmt.Println(getMaxRepeatCountRecursion("aaabbcc", 1, 1, 0))

	// Output:
	// 3
	// 3
}
