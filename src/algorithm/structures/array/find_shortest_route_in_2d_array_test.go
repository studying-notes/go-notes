package array

import "fmt"

func Example_findShortestRouteIn2dArray() {
	fmt.Println(findShortestRouteIn2dArray([][]int{
		{1, 4, 3},
		{8, 7, 5},
		{2, 1, 5},
	}))

	// Output:
	// 17
}
