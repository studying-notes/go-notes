package permutation

import "fmt"

func ExampleNumberCombinator() {
	numbers := []int{1, 2, 2, 3, 4, 5}
	c := NewNumberCombinator()
	r := c.Combine(numbers)
	fmt.Printf("Length: %d\n", len(r))

	// Output:
	// Length: 198
}
