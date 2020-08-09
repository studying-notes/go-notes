package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

var pHead, pEnd *BNode

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root := Array2Tree(array, 0, len(array)-1)
	InOrderBSTree(root)
	// 正向遍历
	for cur := pHead; cur != nil; cur = cur.RightChild {
		fmt.Print(cur.Data, " ")
	}
	fmt.Println()
	for cur := pEnd; cur != nil; cur = cur.LeftChild {
		fmt.Print(cur.Data, " ")
	}
}

func InOrderBSTree(root *BNode) {
	if root == nil {
		return
	}
	InOrderBSTree(root.LeftChild)
	root.LeftChild = pEnd
	if pEnd != nil {
		pEnd.RightChild = root
	} else {
		pHead = root
	}
	pEnd = root
	InOrderBSTree(root.RightChild)
}
