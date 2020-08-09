package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := Array2Tree(array, 0, len(array)-1)
	PrintMidOrder(root)
	fmt.Println()
	fmt.Println(MoreMidNode(root))
}

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

func MoreMidNode(root *BNode) (result *BNode) {
	minNode := getMinNode(root)
	maxNode := getMaxNode(root)
	mid := (minNode.Data + maxNode.Data) / 2
	cur := root
	for cur != nil {
		if root.Data <= mid {
			root = root.RightChild
		} else {
			result = root
			root = root.LeftChild
		}
	}
	//for cur.LeftChild != nil || root.RightChild != nil {
	//	for cur.LeftChild != nil {
	//		if cur.Data > mid {
	//			if cur.LeftChild.Data < mid {
	//				return cur.LeftChild
	//			}
	//			cur = cur.LeftChild
	//		} else {
	//			cur = cur.RightChild
	//		}
	//	}
	//	if cur.Data > mid {
	//		return cur
	//	}
	//	if root.RightChild != nil {
	//		cur = cur.RightChild
	//	}
	//}
	return cur
}
