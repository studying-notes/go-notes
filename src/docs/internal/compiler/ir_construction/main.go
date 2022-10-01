package main

import "fmt"

var o *int

func do() {
	a := 1
	func1(&a)
}

func func1(a *int) {
	fmt.Println(*a)
	*a = 2
}

func main() {
	l := new(int)
	*l = 42
	m := &l
	n := &m // &&l
	o = **n // l
}
