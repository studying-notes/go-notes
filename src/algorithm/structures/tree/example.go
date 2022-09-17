package tree

var FullBinaryTree = &BNode{
	Data: 1,
	LeftChild: &BNode{
		Data:       2,
		LeftChild:  &BNode{Data: 4},
		RightChild: &BNode{Data: 5},
	},
	RightChild: &BNode{
		Data:       3,
		LeftChild:  &BNode{Data: 6},
		RightChild: &BNode{Data: 7},
	},
}
