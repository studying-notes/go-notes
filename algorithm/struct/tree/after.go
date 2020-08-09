package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := Array2Tree(array, 0, len(array)-1)
	order := []int{1, 3, 2, 5, 7, 6, 4}
	fmt.Println(IsPostOrder(root, order))
	fmt.Println(IsAfterOrder(order, 0, len(order)-1))
}

// 未指定是某棵二元查找树
func IsAfterOrder(array []int, start int, end int) bool {
	if start > end {
		return true
	}
	if array == nil {
		return false
	}
	root := array[end]
	var i, j int
	for i = start; i < end; i++ {
		if array[i] > root {
			break
		}
	}
	for j = i; j < end; j++ {
		if array[j] < root {
			return false
		}
	}
	return IsAfterOrder(array, start, i-1) && IsAfterOrder(array, i+1, end)
}

// 指定是某棵二元查找树
func IsPostOrder(root *BNode, array []int) bool {
	if root == nil {
		return false
	}
	if root.Data != array[len(array)-1] {
		return false
	}
	for i := 0; i < len(array); i++ {
		if array[i] > array[len(array)-1] {
			return IsPostOrder(root.LeftChild, array[0:i]) &&
				IsPostOrder(root.RightChild, array[i:len(array)-1])
		}
	}
	return true
}
