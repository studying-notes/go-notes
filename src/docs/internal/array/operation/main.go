package main

import "fmt"

func main() {
	var array [10]int

	var slice = array[:]

	slice[0] = 1

	fmt.Printf("%v", array)
	fmt.Printf("%v", slice)
}
