package tree

func IsEqual(root1, root2 *BNode) bool {
	if root1 == nil && root2 == nil {
		return true
	} else if root1 == nil || root2 == nil {
		return false
	}
	return root1.Data == root2.Data && IsEqual(root1.LeftChild,
		root2.LeftChild) && IsEqual(root1.RightChild, root1.RightChild)
}
