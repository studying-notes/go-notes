package link

func ExampleSum2LinkedList() {
	list1 := NewCustomLinkedList([]int{3, 1, 5})
	PrintLNode(&list1)

	list2 := NewCustomLinkedList([]int{5, 9, 4, 9, 9, 9, 9})
	PrintLNode(&list2)

	sum := Sum2LinkedList(&list1, &list2)
	PrintLNode(sum)

	// Output:
	// 3 1 5
	// 5 9 4 9 9 9 9
	// 8 0 0 0 0 0 0 1
}
