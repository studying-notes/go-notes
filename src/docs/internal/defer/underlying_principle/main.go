package main

func add(a, b int) int {
	return a + b
}

func f() {
	for i := 0; i < 2; i++ {
		defer add(3, 4)
	}
}

func main() {
	f()
}
