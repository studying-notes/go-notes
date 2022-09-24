package array

import . "algorithm/math"

// 求数组中两个元素的最小距离

// 双重遍历法
func findDist(array []int, n1, n2 int) (dist int) {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1
	}
	dist = 1<<63 - 1
	d := 0
	for i, v := range array {
		if v == n1 {
			for j, v := range array {
				if v == n2 {
					d = Abs(i - j)
					if d < dist {
						dist = d
					}
				}
			}
		}
	}
	return dist
}

// 动态规划法
func findDistDyn(array []int, n1, n2 int) (dist int) {
	if array == nil || len(array) == 0 {
		return 1<<63 - 1
	}

	// 记录当前最小距离
	dist = 1<<63 - 1

	d := 0

	// Go 中移位运算符优先级大于减号
	// Python 中减号优先级大于移位运算符

	x, y := -1<<31, 1<<31-1 // 防溢出

	for i, v := range array {
		if v == n1 {
			x = i
		} else if v == n2 {
			y = i
		}
		d = Abs(x - y)
		if d < dist {
			dist = d
		}
	}

	return dist
}

// 求解最小三元组距离

func findDist3Array(a1, a2, a3 []int) int {
	b1, b2, b3 := 0, 0, 0
	dist := 1<<63 - 1
	for b1 < len(a1) && b2 < len(a2) && b3 < len(a3) {
		// 通过数学等式变形可以得出以下两种计算结果等价
		// d := 1/2 * (Abs(a1[b1]-a2[b2]) + Abs(a1[b1]-a3[b3]) + Abs(a2[b2]-a3[b3]))
		if d := MaxN(Abs(a1[b1]-a2[b2]), Abs(a1[b1]-a3[b3]), Abs(a2[b2]-a3[b3])); d < dist {
			dist = d
		}
		switch MinN(a1[b1], a2[b2], a3[b3]) {
		case a1[b1]:
			b1++
		case a2[b2]:
			b2++
		case a3[b3]:
			b3++
		}
	}
	return dist
}
