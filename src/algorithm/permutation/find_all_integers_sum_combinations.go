package permutation

// 求正整数 n 所有可能的整数和组合

type IntegerSumCombiner struct {
	combination []int // 存储和为 n 的组合方式
	result      [][]int
}

func NewIntegerSumCombiner(n int) *IntegerSumCombiner {
	return &IntegerSumCombiner{
		combination: make([]int, n),
	}
}

func (c *IntegerSumCombiner) combine(sum, count int) {
	if sum < 0 {
		return
	} else if sum == 0 {
		c.result = append(c.result, deepcopy(c.combination[:count]))
		return
	}

	var next = 1
	if count != 0 {
		// 保证组合中的下一个数字一定不会小于前一个数字从而保证了组合的递增性
		next = c.combination[count-1]
	}

	for i := next; i <= sum; i++ {
		c.combination[count] = i
		count++
		c.combine(sum-i, count)
		count--
	}
}
