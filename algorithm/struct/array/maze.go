package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	matrix := [][]int{
		{1, 1, 1, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
		{1, 1, 0, 1},
	}
	// [[0 0] [1 0] [1 1] [1 2] [1 3] [2 3] [3 3]]
	road := MazeSolver(matrix)
	fmt.Println(road)

	roadMatrix := InitMatrix(len(matrix), len(matrix[0]))
	for _, v := range road {
		i, j := v[0], v[1]
		roadMatrix[i][j] = 1
	}
	PrintMatrix(roadMatrix)
}

func MazeSolver(matrix [][]int) (road [][2]int) {
	direct := [][2]int{{0, 0}}
	idxs := []int{0}
	length := len(matrix) - 1
	i, j := 0, 0
	for i < length || j < length {
		if matrix[i][j] == 1 {
			road = append(road, [2]int{i, j})
		}
		if i < length && j < length && matrix[i+1][j] == 1 && matrix[i][j+1] == 1 {
			direct = append(direct, [2]int{i, j + 1})
			idxs = append(idxs, len(road))
		}
		if i < length && matrix[i+1][j] == 1 {
			i++
		} else if j < length && matrix[i][j+1] == 1 {
			j++
		}
		if j < length && (i < length && matrix[i+1][j] == 0 || i == length && matrix[i][j+1] == 0) && matrix[i][j+1] == 0 {
			rear := direct[len(direct)-1]
			idx := idxs[len(idxs)-1]
			i, j = rear[0], rear[1]
			road = road[:idx]
			idxs = idxs[:len(idxs)-1]
			direct = direct[:len(direct)-1]
		}
	}
	return append(road, [2]int{length, length})
}
