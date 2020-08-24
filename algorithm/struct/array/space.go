package main

import "fmt"

func main() {
	R := []int{10, 15, 23, 20, 6, 9, 7, 16}
	O := []int{2, 7, 8, 4, 5, 8, 6, 8}

	sort(R, O)
	fmt.Println(R)
	fmt.Println(O)
}

func schedule(R, O []int, M int) bool {
	sort(R, O)
	left := M
	for idx := range R {
		if R[idx] > left {
			return false
		}
		left -= O[idx]
	}
	return true
}

func swap(array []int, i, j int) {
	array[i], array[j] = array[j], array[i]
}

func sort(R, O []int) {
	for i := len(R); i > 0; i-- {
		for j := 1; j < i; j++ {
			if R[j]-O[j] > R[j-1]-O[j-1] {
				swap(R, j-1, j)
				swap(O, j-1, j)
			}
		}
	}
}
