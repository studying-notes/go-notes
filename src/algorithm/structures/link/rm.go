package link

// 在只给定单链表中某个结点指针的情况下删除该结点

// QuickRmNode 复制数据域
func QuickRmNode(node *LNode) {
	if node == nil || node.Next == nil {
		return // 尾结点无法删除
	}
	n := node.Next
	node.Data = n.Data
	node.Next = n.Next
	n.Next = nil // 清理被删结点
}

// RemoveNode 给出链表和被删结点
func RemoveNode(head *LNode, node *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
	// 找到被删结点的前驱
	for cur != nil && cur.Next != node {
		cur = cur.Next
	}
	if cur == nil {
		return
	}
	n := cur.Next
	cur.Next = cur.Next.Next
	n.Next = nil // 清理删掉的结点
}
