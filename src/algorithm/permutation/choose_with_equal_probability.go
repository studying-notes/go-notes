package permutation

import "math/rand"

func ChooseWithEqualProbability(array []int, m, n int) {
	for i := 0; i < m; i++ {
		j := i + rand.Intn(n-i)
		array[i], array[j] = array[j], array[i]
	}
}
