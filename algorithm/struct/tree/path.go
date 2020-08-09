package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := Array2Tree(array, 0, len(array)-1)
	PrintLayerOrder(root)
	fmt.Println()
	FindPath(root, 7)
	FindRoad(root, 7, 0, []int{})
}

func FindPath(root *BNode, sum int) bool {
	if root == nil && sum == 0 {
		return true
	} else if root == nil {
		return false
	}
	sum -= root.Data
	if FindPath(root.LeftChild, sum) {
		fmt.Print(root.Data, " ")
		return true
	}
	if FindPath(root.RightChild, sum) {
		fmt.Print(root.Data, " ")
		return true
	}
	return false
}

func FindRoad(root *BNode, num, sum int, v []int) {
	sum += root.Data
	v = append(v, root.Data)
	if root.LeftChild == nil && root.RightChild == nil && sum == num {
		fmt.Println(v)
	}
	if root.LeftChild != nil {
		FindRoad(root.LeftChild, num, sum, v)
	}
	if root.RightChild != nil {
		FindRoad(root.RightChild, num, sum, v)
	}
	// 不知道有什么用
	//sum -= v[len(v)-1]
	//v = v[:len(v)-1]
}
