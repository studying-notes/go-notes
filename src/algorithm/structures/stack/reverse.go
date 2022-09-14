package stack

// 翻转栈的所有元素

func ReverseStack(s *Stack) {
	if s.IsEmpty() {
		return
	}
	MoveBottom2Top(s)
	top := s.Pop()
	ReverseStack(s)
	s.Push(top)
}

func MoveBottom2Top(s *Stack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		MoveBottom2Top(s)
		top2 := s.Pop()
		s.Push(top1)
		s.Push(top2)
	} else {
		s.Push(top1)
	}
}
