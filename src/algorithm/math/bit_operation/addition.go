package bit_operation

func Add(x, y int) int {
	sum := x ^ y          // 不考虑进位的加法可以用异或运算来代替
	carry := (x & y) << 1 // 进位的计算可以用与操作代替
	if carry != 0 {
		sum = Add(sum, carry)
	}
	return sum
}
