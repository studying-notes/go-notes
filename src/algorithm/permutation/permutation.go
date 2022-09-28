package permutation

func deepcopy(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}
