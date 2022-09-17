package tree

// Copy1 复制二叉树
func Copy1(root *BNode) *BNode {
	if root == nil {
		return nil
	}
	cp := &BNode{Data: root.Data}
	cp.LeftChild = Copy1(root.LeftChild)
	cp.RightChild = Copy1(root.RightChild)
	return cp
}

// Copy2 复制二叉树
func Copy2(root, cp *BNode) {
	if root == nil {
		return
	}
	cp.Data = root.Data
	if root.LeftChild != nil {
		cp.LeftChild = &BNode{}
		Copy2(root.LeftChild, cp.LeftChild)
	}
	if root.RightChild != nil {
		cp.RightChild = &BNode{}
		Copy2(root.RightChild, cp.RightChild)
	}
}
