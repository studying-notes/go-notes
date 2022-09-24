package array

// GetMaxAndMin 查找数组中元素的最大值和最小值
func GetMaxAndMin(array []int, start, end int) (min, max int) {
	if start == end { // 这里不会越界
		return array[start], array[start]
	}

	mid := (start + end) / 2
	l1, l2 := GetMaxAndMin(array, start, mid)
	r1, r2 := GetMaxAndMin(array, mid+1, end)

	if l1 < r1 {
		min = l1
	} else {
		min = r2
	}

	if l2 > r2 {
		max = l2
	} else {
		max = r2
	}

	return min, max
}

// 找出旋转数组的最小元素

func findTheSmallestElementOfARotatedArray(array []int, start, end int) int {
	if start == end {
		return array[start]
	}

	mid := (start + end) / 2

	// 正好取到边缘
	if mid > 0 && array[mid] < array[mid-1] {
		return array[mid]
	} else if mid+1 < end && array[mid] > array[mid+1] {
		return array[mid+1]
	}

	if array[mid] < array[end] {
		return findTheSmallestElementOfARotatedArray(array, start, mid)
	} else if array[mid] > array[start] {
		return findTheSmallestElementOfARotatedArray(array, mid+1, end)
	} else { // array[start] == array[mid] == array[end]
		left := findTheSmallestElementOfARotatedArray(array, start, mid)
		right := findTheSmallestElementOfARotatedArray(array, mid+1, end)
		if left < right {
			return left
		} else {
			return right
		}
	}
}
