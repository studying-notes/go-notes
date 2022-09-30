package string

import "fmt"

func Example_getLengthOfLongestIncreasingSubstring() {
	fmt.Println(getLengthOfLongestIncreasingSubstring("xbcdza"))
	fmt.Println(getLengthOfLongestIncreasingSubstring("xbcdze"))

	// Output:
	// bcdz
	// bcdz
}
