package math

func xor(x, y int) int {
	return (x | y) & (^x | ^y)
}
