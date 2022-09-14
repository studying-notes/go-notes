package stack

// 求栈中最小元素

// BubbleSort 冒泡排序
func BubbleSort(s *Stack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop().(int)
	if !s.IsEmpty() {
		BubbleSort(s)
		top2 := s.Pop().(int)
		if top1 > top2 {
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

// ExtStack 双栈法
type ExtStack struct {
	stack *Stack
	ext   *Stack
}

func NewExtStack(val ...int) *ExtStack {
	s := ExtStack{&Stack{}, &Stack{}}
	for idx := range val {
		s.Push(val[idx])
	}
	return &s
}

func (s *ExtStack) Push(val int) {
	s.stack.Push(val)
	if s.ext.IsEmpty() {
		s.ext.Push(val)
	} else {
		if val <= s.ext.Top().(int) {
			s.ext.Push(val)
		}
	}
}

func (s *ExtStack) Pop() int {
	top := s.stack.Pop().(int)
	if top == s.Min() {
		s.ext.Pop()
	}
	return top
}

func (s *ExtStack) Min() int {
	return s.ext.Top().(int)
}
