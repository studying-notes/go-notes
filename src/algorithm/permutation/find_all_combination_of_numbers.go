package permutation

import . "algorithm/structures/set"

type NumberCombinator struct {
	s       *Set[int]
	numbers []int
}

func NewNumberCombinator() *NumberCombinator {
	return &NumberCombinator{
		s: NewSet[int](),
	}
}

func (c *NumberCombinator) combine(start int) {
	if start == len(c.numbers)-1 {
		if !c.filter(c.numbers) {
			c.s.Add(toNumber(c.numbers))
		}
	} else {
		for i := start; i < len(c.numbers); i++ {
			c.numbers[start], c.numbers[i] = c.numbers[i], c.numbers[start]
			c.combine(start + 1)
			c.numbers[start], c.numbers[i] = c.numbers[i], c.numbers[start]
		}
	}
}

func (c *NumberCombinator) filter(numbers []int) bool {
	if numbers[2] == 4 {
		return true
	}

	for i := 0; i < len(numbers)-1; i++ {
		if (numbers[i] == 3 && numbers[i+1] == 5) ||
			(numbers[i] == 5 && numbers[i+1] == 3) {
			return true
		}
	}

	return false
}

func (c *NumberCombinator) Combine(numbers []int) []int {
	c.s.Clear()
	c.numbers = numbers
	c.combine(0)
	return c.s.ToList()
}
