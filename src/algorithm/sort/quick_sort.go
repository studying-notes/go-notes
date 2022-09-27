package sort

// Partition 用于快速排序中的分割
func Partition(array []int, low, high int) int {
	pos, val := low, array[low]
	low++
	for low <= high {
		if array[low] < val {
			array[low], array[pos] = array[pos], array[low]
			low++
			pos++
		} else {
			array[high], array[low] = array[low], array[high]
			high--
		}
	}
	array[pos] = val
	return pos
}

// QuickSort 快速排序
func QuickSort(array []int, low, high int) {
	if low > high {
		return
	}
	pos := Partition(array, low, high)
	QuickSort(array, low, pos-1)
	QuickSort(array, pos+1, high)
}

// QuickSelectMedian 基于快速选择的中位数查找
func QuickSelectMedian(seq []float64, low int, high int, k int) float64 {
	if low == high {
		return seq[k]
	}
	for low < high {
		pivot := low/2 + high/2
		pivotValue := seq[pivot]
		storeIdx := low
		seq[pivot], seq[high] = seq[high], seq[pivot]
		for i := low; i < high; i++ {
			if seq[i] < pivotValue {
				seq[storeIdx], seq[i] = seq[i], seq[storeIdx]
				storeIdx++
			}
		}
		seq[high], seq[storeIdx] = seq[storeIdx], seq[high]
		if k <= storeIdx {
			high = storeIdx
		} else {
			low = storeIdx + 1
		}
	}
	if len(seq)%2 == 0 {
		return seq[k-1]/2 + seq[k]/2
	}
	return seq[k]
}
