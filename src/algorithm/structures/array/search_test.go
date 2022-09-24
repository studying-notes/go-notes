package array

import "fmt"

func ExampleIsContainK() {
	array2d := [][]int{
		{1, 2, 8, 10},
		{2, 4, 9, 12},
		{4, 7, 10, 13},
		{6, 8, 11, 15},
	}

	fmt.Println(IsContainK(array2d, 7))

	// Output:
	// true
}
