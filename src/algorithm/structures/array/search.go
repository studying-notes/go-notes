package array

// BinarySearchK 二分查找有序数组中是否存在某个元素
func BinarySearchK(array []int, k int) bool {
	if len(array) == 0 {
		return false
	}
	start, end := 0, len(array)
	var mid int
	for start < end {
		mid = (start + end) / 2
		if array[mid] == k {
			return true
		} else if array[mid] < k {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return false
}

// IsContainK 二分查找有序矩阵中是否存在某个元素
func IsContainK(array2d [][]int, k int) bool {
	if len(array2d) == 0 {
		return false
	}

	// 行列
	rows, columns := len(array2d), len(array2d[0])

	// 从右上角开始查找
	for i, j := 0, columns-1; i < rows && j > 0; {
		if array2d[i][j] == k {
			return true
		} else if array2d[i][j] > k {
			j--
		} else {
			i++
		}
	}
	return false
}
