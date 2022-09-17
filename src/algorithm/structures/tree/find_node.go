package tree

// 找到最小值
func getMinNode(node *BNode) *BNode {
	if node == nil {
		return node
	}
	cur := node
	for cur.LeftChild != nil {
		cur = cur.LeftChild
	}
	return cur
}

// 找到最大值
func getMaxNode(node *BNode) *BNode {
	if node == nil {
		return node
	}
	cur := node
	for cur.RightChild != nil {
		cur = cur.RightChild
	}
	return cur
}

func FindMiddleNode(root *BNode) (result *BNode) {
	minNode := getMinNode(root)
	maxNode := getMaxNode(root)

	mid := (minNode.Data + maxNode.Data) / 2

	cur := root
	for cur.LeftChild != nil || root.RightChild != nil {
		for cur.LeftChild != nil {
			if cur.Data > mid {
				if cur.LeftChild.Data < mid {
					return cur
				}
				cur = cur.LeftChild
			} else {
				cur = cur.RightChild
			}
		}
		if cur.Data > mid {
			return cur
		}
		if root.RightChild != nil {
			cur = cur.RightChild
		}
	}
	return cur
}
