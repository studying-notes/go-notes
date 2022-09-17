package tree

// 判断一个数组是否是二元查找树后序遍历的序列

// VerifySequenceOfBST
// 未指定是某棵二元查找树的后序遍历序列
// 只是判断是否是二元查找树的后序遍历序列
func VerifySequenceOfBST(sequence []int) bool {
	if len(sequence) == 0 {
		return false
	}
	return verifySequenceOfBST(sequence, 0, len(sequence)-1)
}

func verifySequenceOfBST(sequence []int, start, end int) bool {
	if start >= end {
		return true
	}
	root := sequence[end]
	i := start
	for ; i < end; i++ {
		if sequence[i] > root {
			break
		}
	}
	for j := i; j < end; j++ {
		if sequence[j] < root {
			return false
		}
	}
	return verifySequenceOfBST(sequence, start, i-1) && verifySequenceOfBST(sequence, i, end-1)
}

func VerifyArrayOfBST(root *BNode, array []int) bool {
	if root == nil {
		return false
	}

	if root.Data != array[len(array)-1] {
		return false
	}

	for i := 0; i < len(array); i++ {
		if array[i] > array[len(array)-1] {
			return VerifyArrayOfBST(root.LeftChild, array[0:i]) &&
				VerifyArrayOfBST(root.RightChild, array[i:len(array)-1])
		}
	}

	return true
}
