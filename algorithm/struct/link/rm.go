package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	head := NewSeqLinkedList(12)
	PrintLNode(&head)
	cur := head.Next.Next.Next
	//RemoveNode(&head, cur)
	QuickRmNode(cur)
	PrintLNode(&head)
}

func QuickRmNode(node *LNode) {
	if node == nil || node.Next == nil {
		return
	}
	n := node.Next
	node.Data = n.Data
	node.Next = n.Next
	n.Next = nil // 清理被删结点
}

func RemoveNode(head *LNode, node *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
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
