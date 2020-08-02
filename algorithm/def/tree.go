package def

import "fmt"

// 二叉树定义
type BNode struct {
	Data       int
	LeftChild  *BNode
	RightChild *BNode
}

func Array2Tree(array []int, start, end int) *BNode {
	// 必须是大于而不能是等于
	// mid:=(2+3)/2=2
	// 这种情况就会让 mid-1 小于 start
	if start > end {
		return nil
	}
	mid := (start + end + 1) / 2
	//mid := (start + end) / 2  // 这两种区别不大
	root := &BNode{Data: array[mid]}
	root.LeftChild = Array2Tree(array, start, mid-1)
	root.RightChild = Array2Tree(array, mid+1, end)
	return root
}

// 二叉树高度递归法（另一种层序遍历法）
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

// 遍历指定层
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

// 先求高度，再遍历指定层
func PrintLevel(root *BNode) {
	depth := root.Depth()
	for level := 0; level < depth; level++ {
		PrintAtLevel(root, level)
	}
}

// 层序遍历 队列法
func PrintLayerOrder(root *BNode) {
	if root == nil {
		return
	}
	q := &Queue{}
	q.EnQueue(root)
	var cur *BNode
	for !q.IsEmpty() {
		cur = q.DeQueue().(*BNode)
		fmt.Print(cur.Data, " ")
		if cur.LeftChild != nil {
			q.EnQueue(cur.LeftChild)
		}
		if cur.RightChild != nil {
			q.EnQueue(cur.RightChild)
		}
	}
}

// 前序遍历二叉树
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

// 中序遍历二叉树
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

// 后序遍历二叉树
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
