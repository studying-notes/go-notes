package main

import "fmt"

func main() {
	costs := []int{7, 11}
	c1, c2 := 1, 1
	tasks := 100
	for tasks > 2 {
		if (c1+1)*costs[0] < (c2+1)*costs[1] {
			c1++
		} else {
			c2++
		}
		tasks--
	}
	fmt.Println(c1, c2, c1*costs[0], c2*costs[1])
}
