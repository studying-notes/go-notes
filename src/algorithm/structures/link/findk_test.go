package link

import "fmt"

func ExampleFindLastK() {
	list := NewSeqLinkedList(12)
	PrintLNode(&list)

	nodeK := FindLastK(&list, 6)
	fmt.Println(nodeK.Data)

	// Output:
	// 1 2 3 4 5 6 7 8 9 10 11 12
	// 7
}

func ExampleSpinLastK() {
	list := NewSeqLinkedList(12)
	SpinLastK(&list, 6)
	PrintLNode(&list)

	// Output:
	// 7 8 9 10 11 12 1 2 3 4 5 6
}
