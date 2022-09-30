package array

import "algorithm/math"

const MaxInt32 = 1<<31 - 1

func findShortestRouteIn2dArray(routes [][]int) int {
	rows, columns := len(routes), len(routes[0])

	matrix := InitMatrix(rows+1, columns+1)

	for i := 2; i <= rows; i++ {
		matrix[i][0] = MaxInt32
	}

	for j := 2; j <= columns; j++ {
		matrix[0][j] = MaxInt32
	}

	// PrintMatrix(matrix)

	for i := 1; i <= rows; i++ {
		for j := 1; j <= columns; j++ {
			matrix[i][j] = math.MinN(matrix[i-1][j], matrix[i][j-1]) + routes[i-1][j-1]
		}
	}

	// PrintMatrix(matrix)

	return matrix[rows][columns]
}
