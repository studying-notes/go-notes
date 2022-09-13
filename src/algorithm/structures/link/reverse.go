package link

import (
	"fmt"
)

// 实现链表的逆序

// RecursiveReverseChild 递归法
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

// Reverse 实现链表的逆序
// 前驱 Precursor
// 后继 Successor
// 就地逆序，可以按头插法思考
// 用三个指针保存前驱、当前和后继。最后让头节点指向前驱。
func Reverse(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	var pre *LNode   // 前驱结点，初始化为 nil
	var cur *LNode   // 当前节点
	suc := head.Next // 后继结点
	for suc != nil {
		cur = suc.Next
		suc.Next = pre
		pre = suc
		suc = cur
	}
	head.Next = pre
}

// InsertReverse 插入法/头插法
// 从链表的第二个结点开始，把遍历到的结点插入到头结点的后面，直到遍历结束。
// 与就地逆序相比，这种方法不**需要保存前驱结点的地址**，头节点充当了前驱，
// 与递归法相比，这种方法不需要递归地调用，效率更高。
// 对不带头结点的单链表进行逆序，可以先加一个头结点，最后再去掉即可。
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

// ReversePrint 实现链表的逆序打印
// 递归法
func ReversePrint(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	ReversePrint(head.Next)
	fmt.Print(head.Next.Data, " ")
}
