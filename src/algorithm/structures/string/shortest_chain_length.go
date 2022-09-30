package string

func isAdjacent(a, b string) bool {
	difference := 0
	for i := range a {
		if a[i] != b[i] {
			difference++
		}
		if difference > 1 {
			return false
		}
	}
	return difference == 1
}

func findShortestChainLength(words []string, start, target string) int {

	return 0
}
