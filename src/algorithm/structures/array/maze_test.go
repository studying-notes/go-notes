package array

import "fmt"

func ExampleMazeSolver() {
	matrix := [][]int{
		{1, 1, 1, 1},
		{1, 0, 0, 0},
		{1, 1, 1, 1},
		{1, 0, 0, 1},
	}

	fmt.Println(MazeSolver(matrix))

	// Output:
	// [[0 0] [1 0] [2 0] [2 1] [2 2] [2 3] [3 3]]
}
