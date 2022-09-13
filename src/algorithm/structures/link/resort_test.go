package link

func ExampleReverseResort() {
	// -> 0 9 1 8 2 7 3 6 4 5
	// -> 0 8 1 7 2 6 3 5 4
	list := NewCustomLinkedList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	PrintLNode(&list)
	ReverseResort(&list)
	PrintLNode(&list)

	// Output:
	// 0 1 2 3 4 5 6 7 8
	// 0 8 1 7 2 6 3 5 4
}
