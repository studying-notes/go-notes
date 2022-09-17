package avl

import "algorithm/math"

type AVLNode struct {
	Data   int
	Height int
	Count  int
	Left   *AVLNode
	Right  *AVLNode
}

func NewAVLNode(data int) *AVLNode {
	return &AVLNode{Data: data}
}

func GetHeight(n *AVLNode) int {
	if n == nil {
		return 0
	}
	return n.Height
}

// LeftRotate 左旋
func LeftRotate(v *AVLNode) (u *AVLNode) {
	u = v.Right
	v.Right, u.Left = u.Left, v
	v.Height = math.Max(GetHeight(v.Left), GetHeight(v.Right)) + 1
	u.Height = math.Max(GetHeight(u.Left), GetHeight(u.Right)) + 1
	return u // 该子树新的顶点
}

// RightRotate 右旋
func RightRotate(v *AVLNode) (u *AVLNode) {
	u = v.Left
	v.Left, u.Right = u.Right, v
	v.Height = math.Max(GetHeight(v.Left), GetHeight(v.Right)) + 1
	u.Height = math.Max(GetHeight(u.Left), GetHeight(u.Right)) + 1
	return u
}

// LeftRightRotate 左旋然后右旋
func LeftRightRotate(v *AVLNode) (u *AVLNode) {
	v.Left = LeftRotate(v.Left)
	return RightRotate(v)
}

// RightLeftRotate 右旋然后左旋
func RightLeftRotate(v *AVLNode) (u *AVLNode) {
	v.Right = RightRotate(v.Right)
	return LeftRotate(v)
}

// InsertAVL 插入新结点
func InsertAVL(root *AVLNode, data int) *AVLNode {
	if root == nil {
		root = NewAVLNode(data)
	} else if data < root.Data {
		root.Left = InsertAVL(root.Left, data)
		if GetHeight(root.Left)-GetHeight(root.Right) == 2 {
			if data < root.Left.Data {
				root = LeftRotate(root)
			} else {
				root = LeftRightRotate(root)
			}
		}
	} else if data > root.Data {
		root.Right = InsertAVL(root.Right, data)
		if GetHeight(root.Right)-GetHeight(root.Left) == 2 {
			if data < root.Right.Data {
				root = RightLeftRotate(root)
			} else {
				root = RightRotate(root)
			}
		}
	}
	root.Height = math.Max(GetHeight(root.Left), GetHeight(root.Right)) + 1
	return root
}
