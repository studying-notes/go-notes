package permutation

import (
	"math/rand"
)

func func1() int {
	if rand.Int()%2 == 0 {
		return 0
	}
	return 1
}

func func2() int {
	if func1() == 0 && func1() == 0 {
		return 0
	}
	return 1
}
