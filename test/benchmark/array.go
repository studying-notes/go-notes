package benchmark

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func LossXor(array []int) (s int) {
	for k, v := range array {
		s ^= (k + 1) ^ v
	}
	return s
}

func LossSub(array []int) (s int) {
	for k, v := range array {
		s += k + 1 - v
	}
	return s
}

func SpinArrayReverse(array []int, idx int) {
	ReverseArray(array[:idx])
	ReverseArray(array[idx:])
	ReverseArray(array)
}

func SpinArrayAppend(array []int, idx int) []int {
	return append(array[idx:], array[:idx]...)
}
