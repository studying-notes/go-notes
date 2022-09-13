package link

func main() {
	head := NewSeqLinkedList(13)
	PrintLNode(&head)
	// FlipAdjacentNode(&head)
	// FlipAdjPointer(&head)
	FlipAdjKNode(&head, 3)
	PrintLNode(&head)
}
