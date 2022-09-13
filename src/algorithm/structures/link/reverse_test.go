package link

func ExampleReverse() {
	LinkedList := NewSeqLinkedList(12)
	head := &LinkedList

	Reverse(head)
	PrintLNode(head)
	// Output:
	// 12 11 10 9 8 7 6 5 4 3 2 1
}

func ExampleReversePrint() {
	LinkedList := NewSeqLinkedList(12)
	head := &LinkedList
	ReversePrint(head)

	// Output:
	// 12 11 10 9 8 7 6 5 4 3 2 1
}
