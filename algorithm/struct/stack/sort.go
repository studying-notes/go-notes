package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	s := Stack{9, 8, 1, 6, 5, 10, 3, 12, 11, 4}
	//ExchangeSort(&s)
	SortStack(&s)
	fmt.Println(s)
}

func SortStack(s *Stack) {
	if s.IsEmpty() {
		return
	}
	ExchangeSort(s)
	top := s.Pop()
	SortStack(s)
	s.Push(top)
}

// 原理其实就是冒泡排序
func ExchangeSort(s *Stack) {
	if s.IsEmpty() {
		return
	}
	top1 := s.Pop()
	if !s.IsEmpty() {
		ExchangeSort(s)
		top2 := s.Pop()
		if top1 < top2 {
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
