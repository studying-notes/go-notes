package tree

var pHead, pEnd *BNode

func ConvertToLinkedListByFrontOrder(root *BNode) {
	if root == nil {
		return
	}

	ConvertToLinkedListByFrontOrder(root.LeftChild)

	root.LeftChild = pEnd
	if pEnd != nil {
		pEnd.RightChild = root
	} else {
		pHead = root
	}
	pEnd = root

	ConvertToLinkedListByFrontOrder(root.RightChild)
}
