package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	PrintMatrix(matrix)
	PrintRotateMatrix(matrix)
}

// 将二维数组逆时针旋转 45° 后打印
func PrintRotateMatrix(matrix [][]int) {
	length := len(matrix)
	var i, x, k int
	j := length - 1
	for x < length {
		i = x
		for i < length && j < length {
			fmt.Print(matrix[i][j], "  ")
			i++
			j++
		}
		if i == length {
			x++
			j = 0
		} else {
			k++
			j = length - k - 1
		}
		fmt.Println()
	}
}
