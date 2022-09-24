package array

import "fmt"

func ExampleFindCommonElements() {
	a1 := []int{2, 5, 12, 20, 45, 85}
	a2 := []int{16, 19, 20, 85, 200}
	a3 := []int{3, 4, 15, 20, 39, 72, 85, 190}

	fmt.Println(FindCommonElements(a1, a2, a3))

	// Output:
	// [20 85]
}
