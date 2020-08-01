package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	//s := Stack{9, 8, 1, 6, 5, 10, 3, 12, 11, 4}
	//BubbleSort(&s)
	s := NewExtStack(9, 8, 1, 6, 5, 0, 3, 12, 11, 4)
	fmt.Println(s.Min())
}

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
		if val <= s.ext.Top() {
			s.ext.Push(val)
		}
	}
}

func (s *ExtStack) Pop() int {
	top := s.stack.Pop()
	if top == s.Min() {
		s.ext.Pop()
	}
	return top
}

func (s *ExtStack) Min() int {
	return s.ext.Top()
}

// 冒泡排序
func BubbleSort(s *Stack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		BubbleSort(s)
		top2 := s.Pop()
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
