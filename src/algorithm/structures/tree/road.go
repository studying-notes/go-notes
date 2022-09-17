package tree

// MaxRoad 在二叉树中找出路径最大的和
func MaxRoad(root *BNode) (val int) {
	if root == nil {
		return 0
	}

	// 遍历左子树找出路径最大的和
	left := MaxRoad(root.LeftChild)
	// 遍历右子树找出路径最大的和
	right := MaxRoad(root.RightChild)

	var sum int

	if left <= 0 && right <= 0 {
		// 左右子树的和都小于0则放弃这段路径
		val = root.Data
		sum = val
	} else {
		if right < 0 {
			// 右子树的和小于0则放弃这段路径
			right = 0
			val = left + root.Data
		} else if left < 0 {
			// 左子树的和小于0则放弃这段路径
			left = 0
			val = right + root.Data
		}

		// 加上左右子树和
		sum = left + root.Data + right
	}

	if sum > max {
		max = sum
	}

	return val
}
