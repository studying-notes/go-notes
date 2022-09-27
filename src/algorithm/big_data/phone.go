package big_data

import "strconv"

type PhoneFilter struct {
	bitmap []int // 位图
}

func NewPhoneFilter(maxValue int) *PhoneFilter {
	return &PhoneFilter{bitmap: make([]int, maxValue/(8*strconv.IntSize))}
}

func (p *PhoneFilter) AddToBitmap(phone int) {
	p.bitmap[phone/(8*strconv.IntSize)] |= 1 << uint(phone%(8*strconv.IntSize))
}
