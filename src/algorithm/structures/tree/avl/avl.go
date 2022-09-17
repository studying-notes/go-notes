package avl

import "fmt"

type Node struct {
	Value   int
	Height  int   // Height of the node
	Balance int   // Balance factor of the node
	Left    *Node // Left child of the node
	Right   *Node // Right child of the node
}

// NewNode creates a new node with the given value
func NewNode(value int) *Node {
	return &Node{Value: value, Height: 1}
}

// Insert inserts a new node with the given value into the tree
func (n *Node) Insert(value int) *Node {
	if n == nil {
		return NewNode(value)
	}

	if value < n.Value {
		n.Left = n.Left.Insert(value)
	} else {
		n.Right = n.Right.Insert(value)
	}

	return n.rebalance()
}

// Delete deletes the node with the given value from the tree
func (n *Node) Delete(value int) *Node {
	if n == nil {
		return nil
	}

	if value < n.Value {
		n.Left = n.Left.Delete(value)
	} else if value > n.Value {
		n.Right = n.Right.Delete(value)
	} else {
		if n.Left == nil && n.Right == nil {
			return nil
		} else if n.Left == nil {
			return n.Right
		} else if n.Right == nil {
			return n.Left
		}

		min := n.Right.min()
		n.Value = min.Value
		n.Right = n.Right.Delete(min.Value)
	}

	return n.rebalance()
}

// Search searches for the node with the given value in the tree
func (n *Node) Search(value int) *Node {
	if n == nil {
		return nil
	}

	if value < n.Value {
		return n.Left.Search(value)
	} else if value > n.Value {
		return n.Right.Search(value)
	}

	return n
}

// min returns the node with the minimum value in the tree
func (n *Node) min() *Node {
	if n.Left == nil {
		return n
	}

	return n.Left.min()
}

// max returns the node with the maximum value in the tree
func (n *Node) max() *Node {
	if n.Right == nil {
		return n
	}

	return n.Right.max()
}

// rebalance the tree
func (n *Node) rebalance() *Node {
	n.updateHeight()
	n.updateBalance()

	if n.Balance == -2 {
		if n.Left.Balance <= 0 {
			return n.rotateRight()
		}

		return n.rotateLeftRight()
	} else if n.Balance == 2 {
		if n.Right.Balance >= 0 {
			return n.rotateLeft()
		}

		return n.rotateRightLeft()
	}

	return n
}

// rotateLeft rotates the tree to the left
func (n *Node) rotateLeft() *Node {
	newRoot := n.Right
	n.Right = newRoot.Left
	newRoot.Left = n
	n.updateHeight()
	newRoot.updateHeight()
	return newRoot
}

// rotateRight rotates the tree to the right
func (n *Node) rotateRight() *Node {
	newRoot := n.Left
	n.Left = newRoot.Right
	newRoot.Right = n
	n.updateHeight()
	newRoot.updateHeight()
	return newRoot
}

// rotateLeftRight rotates the tree to the left and then to the right
func (n *Node) rotateLeftRight() *Node {
	n.Left = n.Left.rotateLeft()
	return n.rotateRight()
}

// rotateRightLeft rotates the tree to the right and then to the left
func (n *Node) rotateRightLeft() *Node {
	n.Right = n.Right.rotateRight()
	return n.rotateLeft()
}

// updateHeight updates the height of the node
func (n *Node) updateHeight() {
	n.Height = max(n.Left.height(), n.Right.height()) + 1
}

// updateBalance updates the balance factor of the node
func (n *Node) updateBalance() {
	n.Balance = n.Right.height() - n.Left.height()
}

// height returns the height of the node
func (n *Node) height() int {
	if n == nil {
		return 0
	}

	return n.Height
}

// max returns the maximum of the two given values
func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Inorder traverses the tree in-order
func (n *Node) Inorder() []int {
	if n == nil {
		return []int{}
	}

	return append(append(n.Left.Inorder(), n.Value), n.Right.Inorder()...)
}

// Preorder traverses the tree preorder
func (n *Node) Preorder() []int {
	if n == nil {
		return []int{}
	}

	return append(append([]int{n.Value}, n.Left.Preorder()...), n.Right.Preorder()...)
}

// Postorder traverses the tree post-order
func (n *Node) Postorder() []int {
	if n == nil {
		return []int{}
	}

	return append(append(n.Left.Postorder(), n.Right.Postorder()...), n.Value)
}

// Levelorder traverses the tree level-order
func (n *Node) Levelorder() []int {
	if n == nil {
		return []int{}
	}

	var result []int
	queue := []*Node{n}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node.Value)

		if node.Left != nil {
			queue = append(queue, node.Left)
		}

		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}

	return result
}

// String returns a string representation of the tree
func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Inorder())
}
