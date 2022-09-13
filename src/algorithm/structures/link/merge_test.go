package link

func ExampleMergeOrderLinkedList() {
	list1 := NewCustomLinkedList([]int{1, 3, 5, 7})
	PrintLNode(&list1)

	list2 := NewCustomLinkedList([]int{2, 4, 6, 8, 9})
	PrintLNode(&list2)

	head := MergeOrderLinkedList(&list1, &list2)
	PrintLNode(head)

	// Output:
	// 1 3 5 7
	// 2 4 6 8 9
	// 1 2 3 4 5 6 7 8 9
}

func ExampleMergeSort() {
	list1 := NewCustomLinkedList([]int{1, 3, 5, 7})
	PrintLNode(&list1)

	list2 := NewCustomLinkedList([]int{2, 4, 6, 8, 9})
	PrintLNode(&list2)

	head := &LNode{Next: MergeSort(list1.Next, list2.Next)}
	PrintLNode(head)

	// Output:
	// 1 3 5 7
	// 2 4 6 8 9
	// 1 2 3 4 5 6 7 8 9
}
