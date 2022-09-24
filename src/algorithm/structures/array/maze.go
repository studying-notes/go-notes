package array

// 求解迷宫问题

func MazeSolver(matrix [][]int) [][2]int {
	rows := len(matrix)
	columns := len(matrix[0])

	var s [][2]int           // 岔口
	var r = [][2]int{{0, 0}} // 路径
	var i, j int

	for i < rows-1 || j < columns-1 {
		if i < rows-1 && matrix[i+1][j] == 1 {
			s = append(s, [2]int{i + 1, j})
		}

		if j < columns-1 && matrix[i][j+1] == 1 {
			s = append(s, [2]int{i, j + 1})
		}

		next := s[len(s)-1]
		s = s[:len(s)-1]
		i, j = next[0], next[1]

		for len(r) > 0 {
			last := r[len(r)-1]
			if last[0] != i && last[1] != j {
				// 清除无效路径
				r = r[:len(r)-1]
			} else {
				break
			}
		}

		r = append(r, next)
	}

	return r
}
