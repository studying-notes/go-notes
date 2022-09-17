package tree

func Mirror(root *BNode) {
	if root == nil {
		return
	}
	root.LeftChild, root.RightChild = root.RightChild, root.LeftChild
	Mirror(root.LeftChild)
	Mirror(root.RightChild)
}
