package bit_operation

func Divide(x, y int) int {
	if y == 0 {
		panic("division by zero")
	} else if x == 0 {
		return 0
	}

	// 记录结果正负号
	negative := false

	// 将负数取相反数参与运算
	if x < 0 {
		negative = !negative
		x = Add(^x, 1)
	}

	if y < 0 {
		negative = !negative
		y = Add(^y, 1)
	}

	// 计算商
	p := 0

	for Multiply(y, p) <= x {
		p++
	}

	p = Subtract(p, 1)

	// 负数则取相反数
	if negative {
		return Add(^p, 1)
	}

	return p
}
