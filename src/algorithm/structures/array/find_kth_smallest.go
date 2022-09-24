package array

func findKthSmallest(array []int, k int) (j int) {
	if k > len(array) {
		panic("k must < array's length")
	}

	for k != 0 {
		for i := len(array) - 1; i > j; i-- {
			if array[i] < array[i-1] {
				array[i], array[i-1] = array[i-1], array[i]
			}
		}
		k, j = k-1, j+1
	}

	return array[j-1]
}

// 类快速排序法 利用快速排序的分割方法

// Partition 用于快速排序中的分割
func Partition(array []int, low, high int) int {
	i, j, val := low, high, array[low]
	for i < j {
		for i < j && array[j] >= val {
			j--
		}
		if i < j {
			array[i] = array[j]
		}
		for i < j && array[i] <= val {
			i++
		}
		if i < j {
			array[j] = array[i]
		}
	}
	array[i] = val
	return i // 中间值索引
}

func findSmallK(array []int, low, high, k int) {
	if low > high {
		return
	}
	pos := Partition(array, low, high)
	if pos+1 == k {
		return
	} else if pos+1 < k {
		// 取右侧
		findSmallK(array, pos+1, high, k)
	} else {
		// 取左侧
		findSmallK(array, low, pos-1, k)
	}
}
