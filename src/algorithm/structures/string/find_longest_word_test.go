package string

import "fmt"

func ExampleLongestWord() {
	w := NewLongestWord([]string{"test", "tester", "testertest",
		"testing", "apple", "seattle", "banana", "batting", "ngcat",
		"batti", "bat", "testingtester", "testbattingcat"})

	fmt.Println(w.GetLongestWord())

	// Output:
	// testbattingcat
}
