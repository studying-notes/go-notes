package bit_operation

func Multiply(x, y int) int {
	if x == 0 || y == 0 {
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

	p, q := 0, 1
	for i := 0; i < 32; i++ {
		q = 1 << i
		if y&q == q { // 判断该为是 1
			p = Add(p, x<<i)
		}
	}

	// 负数则取相反数
	if negative {
		return Add(^p, 1)
	}

	return p
}
