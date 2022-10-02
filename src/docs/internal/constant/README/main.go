package main

import "fmt"

type Numbers int8

const One Numbers = 1
const Two = 2 * One

func main() {
	fmt.Printf("type: %T value: %v", Two, Two)
}
