package tree

// 最小的 64 位整数
var max = -1 << 63

// SumRearOrder 后序遍历法
// 求一棵二叉树的最大子树和
func SumRearOrder(root *BNode) int {
	if root == nil {
		return 0
	}
	// 遍历左子树
	left := SumRearOrder(root.LeftChild)
	// 遍历右子树
	right := SumRearOrder(root.RightChild)
	sum := left + right + root.Data
	if sum > max {
		max = sum
	}
	return sum
}
