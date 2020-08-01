package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {

}

type StackQueue struct {
	en, de *Stack
}

func (q *StackQueue) IsEmpty() bool {
	return q.en.Size() == 0 && q.de.Size() == 0
}

func (q *StackQueue) EnQueue(val int) {
	q.en.Push(val)
}

func (q *StackQueue) DeQueue() int {
	if q.IsEmpty() {
		return 1 << 32
	}
	if q.de.IsEmpty() {
		for !q.en.IsEmpty() {
			q.de.Push(q.en.Pop())
		}
	}
	return q.de.Pop()
}
