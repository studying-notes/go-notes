package array

// 求集合的所有子集

// 位图法
func findAllSubsets(set []int) (subsets [][]int) {
	o := 1 << len(set)

	for q := 1; q < o; q++ {
		var subset []int
		for p := 0; p < len(set); p++ {
			if q>>p&1 == 1 {
				subset = append(subset, set[p])
			}
		}
		subsets = append(subsets, subset)
	}

	return
}

func findAllSubsets2(set []int) [][]int {
	subsets := [][]int{{set[0]}}
	for i := 1; i < len(set); i++ {
		l := len(subsets)
		for j := 0; j < l; j++ {
			subsets = append(subsets, append(subsets[j], set[i]))
		}
		subsets = append(subsets, []int{set[i]})
	}
	return subsets
}
