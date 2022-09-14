package queue

func ExampleWaitQueue() {
	wq := NewSeqQueue(4)
	node := &Node{val: 8}
	wq.EnQueueNode(node)
	wq.EnQueueValues(5, 6, 7)
	PrintQueue(wq)
	node.Leave()
	PrintQueue(wq)

	// Output:
	// No: 1, Val: 1
	// No: 2, Val: 2
	// No: 3, Val: 3
	// No: 4, Val: 4
	// No: 5, Val: 8
	// No: 6, Val: 5
	// No: 7, Val: 6
	// No: 8, Val: 7
	// No: 1, Val: 1
	// No: 2, Val: 2
	// No: 3, Val: 3
	// No: 4, Val: 4
	// No: 6, Val: 5
	// No: 6, Val: 6
	// No: 7, Val: 7
}
