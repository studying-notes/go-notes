package tree

// 把一个有序整数数组放到二叉树中

func ConvertArrayToTree(array []int, start, end int) *BNode {
	// 必须是大于而不能是等于
	// mid:=(2+3)/2=2
	// 这种情况就会让 mid-1 小于 start
	if start > end {
		return nil
	}
	mid := (start + end + 1) / 2
	// mid := (start + end) / 2  // 这两种实际没区别
	root := &BNode{Data: array[mid]}
	root.LeftChild = ConvertArrayToTree(array, start, mid-1)
	root.RightChild = ConvertArrayToTree(array, mid+1, end)
	return root
}
