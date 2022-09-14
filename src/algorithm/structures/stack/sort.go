package stack

// 栈排序

// SortStack 冒泡排序外层循环
func SortStack(s *Stack) {
	if s.IsEmpty() {
		return
	}
	ExchangeSort(s)
	top := s.Pop()
	SortStack(s)
	s.Push(top)
}

// ExchangeSort 冒泡排序内层循环
func ExchangeSort(s *Stack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		ExchangeSort(s)
		top2 := s.Pop()
		if top1.(int) < top2.(int) {
			s.Push(top1)
			s.Push(top2)
		} else {
			s.Push(top2)
			s.Push(top1)
		}
	} else {
		s.Push(top1)
	}
}
