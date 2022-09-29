package math

import (
	"math/rand"
)

func rand7() int {
	return rand.Intn(7) + 1
}

func rand49() int {
	// 1 2 3 4 5 6 7
	// 0 7 14 21 28 35 42
	// 1 ~ 49
	return rand7() + (rand7()-1)*7
}

func rand10() (n int) {
	for n == 0 || n > 40 {
		n = rand49()
	}
	return n%10 + 1
}
