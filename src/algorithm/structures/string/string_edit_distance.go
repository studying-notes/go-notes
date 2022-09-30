package string

import (
	"algorithm/math"
	"algorithm/structures/array"
)

func getStringEditDistance(a, b string) int {
	al, bl := len(a), len(b)

	if al == 0 || bl == 0 {
		return al + bl
	}

	matrix := array.InitMatrix(al+1, bl+1)

	for i := 1; i <= al; i++ {
		matrix[i][0] = i
	}

	for j := 1; j <= bl; j++ {
		matrix[0][j] = j
	}

	// array.PrintMatrix(matrix)

	for i := 1; i <= al; i++ {
		for j := 1; j <= bl; j++ {
			matrix[i][j] = math.MinN(matrix[i][j-1], matrix[i-1][j], matrix[i-1][j-1]) + 1
			if a[i-1] == b[j-1] {
				matrix[i][j]--
			}
		}
	}

	// array.PrintMatrix(matrix)

	return math.MinN(matrix[al-1][bl-1])
}
