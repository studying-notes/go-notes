package tree

import (
	"fmt"

	. "algorithm/structures/queue"
)

// BNode 二叉树定义
type BNode struct {
	Data       int
	LeftChild  *BNode
	RightChild *BNode
}

// PrintFrontOrder 前序遍历二叉树
func PrintFrontOrder(root *BNode) {
	if root == nil {
		return
	}
	// 打印当前值
	fmt.Print(root.Data, " ")
	// 遍历左子树
	PrintFrontOrder(root.LeftChild)
	// 遍历右子树
	PrintFrontOrder(root.RightChild)
}

// PrintMidOrder 中序遍历二叉树
func PrintMidOrder(root *BNode) {
	if root == nil {
		return
	}
	// 遍历左子树
	PrintMidOrder(root.LeftChild)
	// 打印当前值
	fmt.Print(root.Data, " ")
	// 遍历右子树
	PrintMidOrder(root.RightChild)
}

// PrintRearOrder 后序遍历二叉树
func PrintRearOrder(root *BNode) {
	if root == nil {
		return
	}
	// 遍历左子树
	PrintRearOrder(root.LeftChild)
	// 遍历右子树
	PrintRearOrder(root.RightChild)
	// 打印当前值
	fmt.Print(root.Data, " ")
}

// PrintLayerOrder 层序遍历 队列法
func PrintLayerOrder(root *BNode) {
	if root == nil {
		return
	}
	q := &Queue{}
	// 将根节点入队
	q.EnQueue(root)
	var cur *BNode
	// 队列不为空
	for !q.IsEmpty() {
		// 出队
		cur = q.DeQueue().(*BNode)
		// 打印当前值
		fmt.Print(cur.Data, " ")
		// 左子树入队
		if cur.LeftChild != nil {
			q.EnQueue(cur.LeftChild)
		}
		// 右子树入队
		if cur.RightChild != nil {
			q.EnQueue(cur.RightChild)
		}
	}
}

// Depth 层序遍历 递归法求二叉树高度
func (root *BNode) Depth() int {
	if root == nil {
		return 0
	}
	lDepth := root.LeftChild.Depth()
	rDepth := root.RightChild.Depth()
	if rDepth > lDepth {
		return rDepth + 1
	} else {
		return lDepth + 1
	}
}

// PrintAtLevel 遍历指定层
func PrintAtLevel(root *BNode, level int) int {
	if root == nil || level < 0 {
		return 0
	} else if level == 0 {
		fmt.Print(root.Data, " ")
		return 1
	} else {
		return PrintAtLevel(root.LeftChild, level-1) +
			PrintAtLevel(root.RightChild, level-1)
	}
}

// PrintLevel 先求高度，再遍历指定层
func PrintLevel(root *BNode) {
	depth := root.Depth()
	for level := 0; level < depth; level++ {
		PrintAtLevel(root, level)
	}
}
