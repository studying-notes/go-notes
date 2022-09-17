package tree

import "fmt"

// 在二叉树中找出与输入整数相等的所有路径

func FindPath1(root *BNode, sum int) bool {
	if root == nil && sum == 0 {
		return true
	} else if root == nil {
		return false
	}

	sum -= root.Data

	if FindPath1(root.LeftChild, sum) {
		fmt.Println(root.Data)
		return true
	}

	if FindPath1(root.RightChild, sum) {
		fmt.Println(root.Data)
		return true
	}

	return false
}

func FindPath2(root *BNode, num, sum int, v []int) {
	sum += root.Data
	v = append(v, root.Data)

	if root.LeftChild == nil && root.RightChild == nil && sum == num {
		fmt.Println(v)
		return
	}

	if root.LeftChild != nil {
		FindPath2(root.LeftChild, num, sum, v)
	}

	if root.RightChild != nil {
		FindPath2(root.RightChild, num, sum, v)
	}
}
