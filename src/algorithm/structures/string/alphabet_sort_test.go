package string

import "fmt"

func ExampleAlphabetComparer() {
	var alphabetSequence = []byte{'d', 'g', 'e', 'c', 'f', 'b', 'o', 'a'}
	c := NewAlphabetComparer(alphabetSequence)
	fmt.Println(c.MoreThan("dog", "eye"))

	// Output:
	// false
}

func ExampleSortStringSequences() {
	var alphabetSequence = []byte{'d', 'g', 'e', 'c', 'f', 'b', 'o', 'a'}
	c := NewAlphabetComparer(alphabetSequence)
	fmt.Println(c.SortStringSequences([]string{"bed", "dog", "dear", "eye"}))

	// Output:
	// [dear dog eye bed]
}
