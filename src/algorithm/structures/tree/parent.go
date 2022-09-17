package tree

import (
	"math"

	. "algorithm/structures/stack"
)

// 找出排序二叉树上任意两个结点的最近共同父结点

// 路径对比法

// GetPathFromRoot 根结点到指定结点的路径
func GetPathFromRoot(root *BNode, node *BNode, s *Stack) bool {
	if root == nil {
		return false
	}

	if root.Data == node.Data {
		s.Push(root)
		return true
	}

	if GetPathFromRoot(root.LeftChild, node, s) ||
		GetPathFromRoot(root.RightChild, node, s) {
		s.Push(root)
		return true
	}

	return false
}

// FindCommonParentNode 找到最近的公共父节点
func FindCommonParentNode(root, n1, n2 *BNode) (parent *BNode) {
	s1 := &Stack{}
	s2 := &Stack{}

	GetPathFromRoot(root, n1, s1)
	GetPathFromRoot(root, n2, s2)

	// 最远父节点在栈顶
	for s1.Top().(*BNode).Data == s2.Pop().(*BNode).Data {
		// 最后一个相同的父节点
		parent = s1.Pop().(*BNode)
	}

	return parent
}

// 结点编号法

// GetNodeNumber 按照完全二叉树中对结点编号的方式进行编号
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

// GetNodeFromNumber 根据编号获取二叉树的结点
func GetNodeFromNumber(root *BNode, num int) *BNode {
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
	return GetNodeFromNumber(root, nm1)
}

// FindParentNodeByRearOrder 后序遍历法
func FindParentNodeByRearOrder(root, node1, node2 *BNode) *BNode {
	if root == nil || root.Data == node1.Data || root.Data == node2.Data {
		return root
	}

	lChild := FindParentNodeByRearOrder(root.LeftChild, node1, node2)
	rChild := FindParentNodeByRearOrder(root.RightChild, node1, node2)

	if lChild == nil {
		return rChild
	}
	if rChild == nil {
		return lChild
	}
	return root
}
