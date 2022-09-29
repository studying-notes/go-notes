package bit_operation

func Subtract(x, y int) int {
	return Add(x, Add(^y, 1))
}
