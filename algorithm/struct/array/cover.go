package main

import "fmt"

func main() {
	array := []int{1, 3, 7, 8, 10, 11, 12, 13, 15, 16, 17, 19, 35}
	dst := 8
	x, y, z, final := 0, 1, 0, 1
	for _, val := range array {
		if val-array[x] > dst {
			x++
		}
		if y-x > z {
			z = y - x
			final = y
		}
		y++
	}
	fmt.Println(z, array[final-z:final])
}