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
	Mirror(root)
	PrintMidOrder(root)
}

func Mirror(root *BNode) {
	if root == nil {
		return
	}
	root.LeftChild, root.RightChild = root.RightChild, root.LeftChild
	Mirror(root.LeftChild)
	Mirror(root.RightChild)
}
