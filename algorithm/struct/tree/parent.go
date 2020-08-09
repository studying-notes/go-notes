package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
	"math"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root := Array2Tree(array, 0, len(array)-1)
	n1 := root.LeftChild.LeftChild.LeftChild
	n2 := root.LeftChild.RightChild
	//s := &Stack{}
	//PathFromRoot(root, n1, s)
	//PathFromRoot(root, n2, s)
	//for !s.IsEmpty() {
	//	fmt.Println(s.Pop().(*BNode).Data)
	//}

	fmt.Println(FindParentNode(root, n1, n2).Data)

	fmt.Println(GetNodeNumber(root, n1, 1))
	fmt.Println(GetNodeNumber(root, n2, 1))

	fmt.Println(FindParent(root, n1, n2).Data)
}

func FindParentNodeReverse(root, node1, node2 *BNode) *BNode {
	if root == nil || root.Data == node1.Data || root.Data == node2.Data {
		return root
	}
	lChild := FindParentNodeReverse(root.LeftChild, node1, node2)
	rChild := FindParentNodeReverse(root.RightChild, node1, node2)
	if lChild == nil {
		return rChild
	} else if rChild == nil {
		return lChild
	} else {
		return root
	}
}

// 按照完全二叉树中对结点编号的方式进行编号
func GetNodeNumber(root, node *BNode, number int) (bool, int) {
	if root == nil {
		return false, number
	} else if root == node {
		return true, number
	} else if ok, num := GetNodeNumber(root.LeftChild, node, number<<1); ok {
		return true, num
	}
	return GetNodeNumber(root.RightChild, node, number<<1+1)
}

// 根据编号获取二叉树的结点
func GetNodeFromNum(root *BNode, num int) *BNode {
	if root == nil || num < 0 {
		return nil
	} else if num == 0 {
		return root
	}
	// 二进制数长度 - 1 后的结果
	lg := (uint)(math.Log2(float64(num)))
	// 减去根结点
	num -= 1 << lg
	for lg > 0 {
		if ((1 << (lg - 1)) & num) == 1 {
			root = root.RightChild
		} else {
			root = root.LeftChild
		}
		lg--
	}
	return root
}

func FindParent(root, n1, n2 *BNode) *BNode {
	_, nm1 := GetNodeNumber(root, n1, 1)
	_, nm2 := GetNodeNumber(root, n2, 1)
	for nm1 != nm2 {
		if nm1 > nm2 {
			nm1 /= 2
		} else {
			nm2 /= 2
		}
	}
	return GetNodeFromNum(root, nm1)
}

// 根结点到指定结点的路径
func GetPathFromRoot(root *BNode, node *BNode, s *Stack) bool {
	if root == nil {
		return false
	}
	//fmt.Println(root.Data, node.Data)
	if root.Data == node.Data {
		s.Push(root)
		return true
	}
	if GetPathFromRoot(root.LeftChild, node, s) || GetPathFromRoot(root.RightChild, node, s) {
		s.Push(root)
		return true
	}
	return false
}

// 找到最近的公共父节点
func FindParentNode(root, n1, n2 *BNode) (parent *BNode) {
	s1 := &Stack{}
	s2 := &Stack{}
	GetPathFromRoot(root, n1, s1)
	GetPathFromRoot(root, n2, s2)
	for s1.Top().(*BNode).Data == s2.Pop().(*BNode).Data {
		parent = s1.Pop().(*BNode)
	}
	return parent
}
