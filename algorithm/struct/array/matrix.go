package main

import "fmt"

func main() {
	// A x B x C x D
	//p := []int{40, 20, 30, 10, 30}
	p := []int{1, 5, 2, 4, 6} // 42
	fmt.Println(bestMatrixChainOrder(p, 1, len(p)-1))
}

func bestMatrixChainOrder(p []int, i, j int) int {
	if i == j {
		return 0
	}
	best := 1<<63 - 1
	for k := i; k < j; k++ {
		count := bestMatrixChainOrder(p, i, k) + bestMatrixChainOrder(p, k+1, j) + p[i-1]*p[k]*p[j]
		if count < best {
			best = count
		}
	}
	return best
}
