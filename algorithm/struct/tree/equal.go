package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array1 := []int{-1, 3, 9, 6, -7}
	root1 := Array2Tree(array1, 0, len(array1)-1)
	array2 := []int{-1, 3, 8, 6, -7}
	root2 := Array2Tree(array2, 0, len(array2)-1)
	fmt.Println(IsEqual(root1, root2))
}

func IsEqual(root1, root2 *BNode) bool {
	if root1 == nil && root2 == nil {
		return true
	} else if root1 == nil || root2 == nil {
		return false
	}
	return root1.Data == root2.Data && IsEqual(root1.LeftChild,
		root2.LeftChild) && IsEqual(root1.RightChild, root1.RightChild)
}
