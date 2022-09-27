package string

import "fmt"

func ExamplePermutation() {
	Permutation([]rune("abc"), 0)
	fmt.Println(result)

	// Output:
	// [abc acb bac bca cba cab]
}
