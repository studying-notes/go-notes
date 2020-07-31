package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	LinkedList := NewSeqLinkedList(12)
	head := &LinkedList
	ReversePrint(head)
}

// 实现链表的逆序打印
func ReversePrint(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	ReversePrint(head.Next)
	fmt.Print(head.Next.Data, " ")
}

// 实现链表的逆序
// 前驱 Precursor
// 后继 Successor

// 就地逆序
func Reverse(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	var pre *LNode // 前驱结点，初始化为 nil
	var cur *LNode
	suc := head.Next // 后继结点
	for suc != nil {
		cur = suc.Next
		suc.Next = pre
		pre = suc
		suc = cur
	}
	head.Next = pre
}

// 递归法
// 不带头结点逆序
func RecursiveReverseChild(node *LNode) *LNode {
	if node == nil || node.Next == nil {
		return node
	}
	head := RecursiveReverseChild(node.Next)
	// 1 -> 2 -> 3 -> nil
	// 1 -> 2 -> 3 -> 2 -> nil
	// 1 -> 2 -> 3 -> 2 -> 1 -> nil
	// 构成环，然后断掉
	node.Next.Next = node
	node.Next = nil
	return head
}

func RecursiveReverse(head *LNode) {
	first := head.Next
	tail := RecursiveReverseChild(first)
	head.Next = tail
}

// 插入法
func InsertReverse(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	var suc *LNode
	cur := head.Next.Next
	head.Next.Next = nil // 十分关键
	for cur != nil {
		// 构成环，然后断掉
		suc = cur.Next
		cur.Next = head.Next
		head.Next = cur
		cur = suc
	}
}
