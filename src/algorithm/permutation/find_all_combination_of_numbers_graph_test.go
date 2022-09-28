package permutation

import "fmt"

func ExampleNumberCombinatorGraph() {
	numbers := []int{1, 2, 2, 3, 4, 5}
	c := NewNumberCombinatorGraph(numbers)
	r := c.getAllCombinations()
	fmt.Printf("Length: %d\n", len(r))

	// Output:
	// Length: 198
}
