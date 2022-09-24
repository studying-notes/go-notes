package array

import "fmt"

func ExampleFindAllSubsets() {
	subsets := findAllSubsets([]int{1, 2, 3})
	for _, subset := range subsets {
		fmt.Println(subset)
	}

	// Output:
	// [1]
	// [2]
	// [1 2]
	// [3]
	// [1 3]
	// [2 3]
	// [1 2 3]
}

func ExampleFindAllSubsets2() {
	subsets := findAllSubsets2([]int{1, 2, 3})
	for _, subset := range subsets {
		fmt.Println(subset)
	}

	// Output:
	// [1]
	// [1 2]
	// [2]
	// [1 3]
	// [1 2 3]
	// [2 3]
	// [3]
}
