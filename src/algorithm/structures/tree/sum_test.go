package tree

import "fmt"

func ExampleSumRearOrder() {
	root := &BNode{
		Data: 1,
		LeftChild: &BNode{
			Data:       2,
			LeftChild:  &BNode{Data: 4},
			RightChild: &BNode{Data: 5},
		},
		RightChild: &BNode{
			Data:       -3,
			LeftChild:  &BNode{Data: -1},
			RightChild: &BNode{Data: -1},
		},
	}

	fmt.Println(SumRearOrder(root))
	fmt.Println(max)

	// Output:
	// 7
	// 11
}
