package array

import (
	"fmt"
	"strconv"
)

// PrintMatrix 打印二维数组/矩阵
func PrintMatrix(matrix [][]int) {
	repr := ""

	for _, i := range matrix {
		for n, j := range i {
			if n > 0 {
				repr += ", "
			}
			repr += strconv.Itoa(j)
		}

		repr += "\n"
	}

	fmt.Println(repr)
}

// InitMatrix 初始化二维数组/矩阵
func InitMatrix(y, x int) [][]int {
	matrix := make([][]int, y)
	for i := range matrix {
		matrix[i] = make([]int, x)
	}
	return matrix
}
