package array

func schedule(R, O []int, M int) bool {
	sort(R, O)
	left := M
	for idx := range R {
		if R[idx] > left {
			return false
		}
		left -= O[idx]
	}
	return true
}

func swap(array []int, i, j int) {
	array[i], array[j] = array[j], array[i]
}

// 冒泡排序
func sort(R, O []int) {
	for i := len(R); i > 0; i-- {
		for j := 1; j < i; j++ {
			if R[j]-O[j] > R[j-1]-O[j-1] {
				swap(R, j-1, j)
				swap(O, j-1, j)
			}
		}
	}
}
