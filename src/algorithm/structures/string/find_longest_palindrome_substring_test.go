package string

import "fmt"

func ExampleGetLongestPalindrome() {
	p := Palindrome{}
	fmt.Println(p.getLongestPalindrome("abcdefgfedxyz"))

	// Output:
	// defgfed
}
