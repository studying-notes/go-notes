package main

import "fmt"

func main() {
	D := []int{120, 120, 120} // N
	//P := []int{60, 60, 80, 20, 80} // M
	//P := []int{60, 80, 80, 20, 80} // M
	P := []int{60, 60, 60, 60, 60, 60} // M
	i, j := 0, 0
	for i < len(D) && j < len(P) {
		for i < len(D) && P[j] > D[i] {
			i++
		}
		if i == len(D) {
			break
		}
		D[i] -= P[j]
		j++
	}
	if i == len(D) {
		fmt.Println("fail")
	} else {
		fmt.Println("success")
	}
}
