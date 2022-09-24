package array

// ReverseArray 逆序数组
func ReverseArray(array []int) {
	start, end := 0, len(array)-1
	for start < end {
		array[start], array[end] = array[end], array[start]
		start, end = start+1, end-1
	}
}
