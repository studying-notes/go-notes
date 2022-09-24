package array

import "fmt"

func ExampleFindThePathWithTheMostCoveragePoints() {
	array := []int{1, 3, 7, 8, 10, 11, 12, 13, 15, 16, 17, 19, 35}

	z, road := FindThePathWithTheMostCoveragePoints(array, 8)
	fmt.Println(z)
	fmt.Println(road)

	// Output:
	// 7
	// [7 8 10 11 12 13 15]
}
