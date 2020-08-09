package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

var max = -1 << 63

func main() {
	array := []int{-1, 3, 9, 6, -7}
	//array := []int{5, 2, 3, -2}
	root := Array2Tree(array, 0, len(array)-1)
	PrintMidOrder(root)
	fmt.Println()
	PrintLayerOrder(root)
	fmt.Println()
	MaxRoad(root)
	fmt.Println(max)
}

// 在二叉树中找出路径最大的和
func MaxRoad(root *BNode) (val int) {
	if root == nil {
		return 0
	}
	// 遍历左子树
	left := MaxRoad(root.LeftChild)
	// 遍历右子树
	right := MaxRoad(root.RightChild)

	var sum int
	if left <= 0 && right <= 0 {
		val = root.Data
		sum = val
	} else if left > right {
		val = left + root.Data
		if right < 0 {
			right = 0
		}
		sum = val + right
	} else {
		val = right + root.Data
		if left < 0 {
			left = 0
		}
		sum = val + left
	}
	if sum > max {
		max = sum
	}
	return val
}

// 后序遍历法 - 求一棵二叉树的最大子树和
func RearOrder(root *BNode) int {
	if root == nil {
		return 0
	}
	// 遍历左子树
	left := RearOrder(root.LeftChild)
	// 遍历右子树
	right := RearOrder(root.RightChild)
	sum := left + right + root.Data
	if sum > max {
		max = sum
	}
	return sum
}
