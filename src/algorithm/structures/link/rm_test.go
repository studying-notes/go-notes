package link

func ExampleQuickRmNode() {
	head := NewSeqLinkedList(12)
	PrintLNode(&head)
	cur := head.Next.Next.Next
	QuickRmNode(cur)
	PrintLNode(&head)

	// Output:
	// 1 2 3 4 5 6 7 8 9 10 11 12
	// 1 2 4 5 6 7 8 9 10 11 12
}
