package string

import "fmt"

func ExampleFindLongestCommonSubstring() {
	substring := FindLongestCommonSubstring([]rune("abccade"), []rune("dgcadde"))
	fmt.Println(string(substring))

	// Output:
	// 0, 0, 0, 0, 0, 0, 0, 0
	// 0, 0, 0, 0, 1, 0, 0, 0
	// 0, 0, 0, 0, 0, 0, 0, 0
	// 0, 0, 0, 1, 0, 0, 0, 0
	// 0, 0, 0, 1, 0, 0, 0, 0
	// 0, 0, 0, 0, 2, 0, 0, 0
	// 0, 1, 0, 0, 0, 3, 1, 0
	// 0, 0, 0, 0, 0, 0, 0, 2
	//
	// cad
}
