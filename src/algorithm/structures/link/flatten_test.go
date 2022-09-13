package link

func ExampleFlatten() {
	head := NewLinked2Node()
	PrintL2Node(head.Next)
	PrintL2Node(head.Next.Next)
	PrintL2Node(head.Next.Next.Next)
	PrintL2Node(head.Next.Next.Next.Next)
	res := Flatten(head)
	PrintL2Node(res)

	// Output:
	// 3 6 8 31
	// 11 21
	// 15 22 50
	// 30 39 40 55
	// 3 6 8 11 15 21 22 30 31 39 40 50 55
}
