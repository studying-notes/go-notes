package permutation

import "fmt"

func ExampleIntegerSumCombiner() {
	c := NewIntegerSumCombiner(4)
	c.combine(4, 0)
	fmt.Println(c.result)

	// Output:
	// [[1 1 1 1] [1 1 2] [1 3] [2 2] [4]]
}
