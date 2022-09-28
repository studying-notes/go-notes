package array

import (
	"fmt"
)

func ExampleFindLongestSquareContinuousSquareWaveSignal() {
	arrays := [][]int{
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1},
	}

	for i := range arrays {
		fmt.Println(FindLongestSquareContinuousSquareWaveSignal(arrays[i]))
	}

	// Output:
	// [0 1 0 1 0]
	// [0 1 0 1 0 1 0 1 0]
	// []
}
