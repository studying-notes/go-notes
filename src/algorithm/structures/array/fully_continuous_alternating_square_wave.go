package array

const (
	low int = iota
	high
)

func FindLongestSquareContinuousSquareWaveSignal(array []int) []int {
	var (
		begin, end       int // 当前信号
		maxBegin, maxEnd int // 最长的信号
	)

	for i := 1; i < len(array); i++ {
		if array[i-1] == low && array[i] == low {
			// 00
			if begin == 0 || begin == i-1 {
				begin = i
			} else if end == 0 {
				end = i
			}
		} else if array[i-1] == high && array[i] == high {
			// 11
			begin, end = 0, 0
		} else if i == len(array)-1 && array[i-1] == high && array[i] == low {
			// 10
			end = i
		}

		if begin != 0 && end != 0 && end-begin > maxEnd-maxBegin {
			maxBegin, maxEnd = begin, end
			begin, end = 0, 0
		}
	}

	return array[maxBegin:maxEnd]
}
