package main

func change(c [3]int) {
	c[0] = 4
}

func main() {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 0
	change(a)
	println(a[0])
}
