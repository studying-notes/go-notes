package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := Array2Tree(array, 0, len(array)-1)
	cp := &BNode{}
	Copy(root, cp)
	PrintMidOrder(cp)
	fmt.Println()
	PrintMidOrder(Duplicate(root))
}

func Duplicate(root *BNode) *BNode {
	if root == nil {
		return nil
	}
	dup := &BNode{Data: root.Data}
	dup.LeftChild = Duplicate(root.LeftChild)
	dup.RightChild = Duplicate(root.RightChild)
	return dup
}

func Copy(root, cp *BNode) {
	if root == nil {
		return
	}
	cp.Data = root.Data
	if root.LeftChild != nil {
		cp.LeftChild = &BNode{}
		Copy(root.LeftChild, cp.LeftChild)
	}
	if root.RightChild != nil {
		cp.RightChild = &BNode{}
		Copy(root.RightChild, cp.RightChild)
	}
}
