package array

import "fmt"

func NewArrayByRule(a []int) {
	b := make([]int, len(a))

	b[0] = 1

	for i := 1; i < len(a); i++ {
		b[i] = a[i-1] * b[i-1]
	}

	for j := len(a) - 1; j > 0; j-- {
		b[j] *= b[0]
		b[0] *= a[j]
	}

	fmt.Println(b)
}
