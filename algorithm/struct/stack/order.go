package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	pushOrder := []int{1, 2, 3, 4, 5}
	popOrder := []int{3, 2, 5, 4, 1}
	//popOrder := []int{5, 3, 4, 1, 2}
	fmt.Println(IsPopOrder(pushOrder, popOrder))
	fmt.Println(IsPopSerial(pushOrder, popOrder))
}

func IsPopSerial(push, pop []int) bool {
	if len(push) != len(pop) {
		return false
	}
	pushIndex := 0
	popIndex := 0
	stack := Stack{}
	for pushIndex < len(push) {
		stack.Push(push[pushIndex])
		pushIndex++
		for !stack.IsEmpty() {
			val := stack.Top()
			if val == pop[popIndex] {
				stack.Pop()
				popIndex++
			} else {
				break
			}
		}
	}
	if stack.IsEmpty() && popIndex == len(pop) {
		return true
	}
	return false
}

func IsPopOrder(pushOrder, popOrder []int) bool {
	if len(pushOrder) != len(popOrder) {
		return false
	}
	s := Stack{}
	for len(popOrder) != 0 {
		if !s.IsEmpty() {
			val := s.Top()
			if val == popOrder[0] {
				popOrder = popOrder[1:]
				s.Pop()
			} else if len(pushOrder) == 0 {
				return false
			} else {
				s.Push(pushOrder[0])
				pushOrder = pushOrder[1:]
			}
		} else if pushOrder[0] == popOrder[0] {
			popOrder = popOrder[1:]
			pushOrder = pushOrder[1:]
		} else {
			s.Push(pushOrder[0])
			pushOrder = pushOrder[1:]
		}
	}
	return true
}
