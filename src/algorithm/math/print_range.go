package math

func printRange(n int) {
	if n > 0 {
		printRange(n - 1)
		println(n)
	}
}
