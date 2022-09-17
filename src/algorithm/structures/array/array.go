package array

import "fmt"

// PrintMatrix 打印二维数组/矩阵
func PrintMatrix(matrix [][]int) {
	for _, i := range matrix {
		for _, j := range i {
			fmt.Print(j, "  ")
		}
		fmt.Println()
	}
}

// InitMatrix 初始化二维数组/矩阵
func InitMatrix(y, x int) [][]int {
	matrix := make([][]int, y)
	for i := range matrix {
		matrix[i] = make([]int, x)
	}
	return matrix
}

// IsEqualArray 比较两个切片是否相等 reflect.DeepEqual 性能极差
func IsEqualArray(a, b []int) bool {
	if len(a) != len(b) || (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// ReverseArray 逆序数组
func ReverseArray(array []int) {
	start, end := 0, len(array)-1
	for start < end {
		array[start], array[end] = array[end], array[start]
		start, end = start+1, end-1
	}
}

// Range 可生成指定范围的整数切片
func Range(args ...int) []int {
	var start, end int
	step := 1
	if len(args) == 1 {
		end = args[0]
	} else if len(args) == 2 {
		start, end = args[0], args[1]
	} else if len(args) > 2 {
		start, end, step = args[0], args[1], args[2]
	}
	if step == 0 || start == end || (step < 0 && start < end) || (step > 0 && start > end) {
		return []int{}
	}
	s := make([]int, 0, (end-start)/step+1)
	for start != end {
		s = append(s, start)
		start += step
	}
	return s
}
