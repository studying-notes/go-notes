package array

import "fmt"

// 将二维数组逆时针旋转 45° 后打印

func PrintRotateMatrix(matrix [][]int) {
	repr := ""
	var i, x, k int

	length := len(matrix)
	j := length - 1
	for x < length {
		i = x
		for i < length && j < length {
			if i > x {
				repr += ", "
			}
			repr += fmt.Sprint(matrix[i][j])
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

		repr += "\n"
	}

	fmt.Println(repr)
}
