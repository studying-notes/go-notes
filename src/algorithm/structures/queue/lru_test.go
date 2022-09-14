package queue

import (
	. "algorithm/structures/set"
)

func ExampleLRU() {
	q := &LRU{size: 3, queue: &Queue{}, set: NewSet()}
	q.AccessPage(1)
	q.AccessPage(2)
	q.AccessPage(3)
	q.AccessPage(2)
	q.AccessPage(5)
	q.AccessPage(6)
	q.PrintQueue()

	// Output:
	// [6 5 2]
}
