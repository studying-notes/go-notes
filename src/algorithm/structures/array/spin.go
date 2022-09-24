package array

// 旋转数组/循环移位

func SpinArrayAppend(array []int, idx int) []int {
	return append(array[idx:], array[:idx]...)
}

// SpinArrayReverse 原地三次逆序
func SpinArrayReverse(array []int, idx int) {
	ReverseArray(array[:idx])
	ReverseArray(array[idx:])
	ReverseArray(array)
}
