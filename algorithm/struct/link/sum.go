package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	list1 := NewCustomLinkedList([]int{3, 1, 5})
	PrintLNode(&list1)
	list2 := NewCustomLinkedList([]int{5, 9, 4, 9, 9, 9, 9})
	PrintLNode(&list2)

	sum := Sum2LinkedList(&list1, &list2)
	PrintLNode(sum)
}

// 链表相加法
func Sum2LinkedList(head1, head2 *LNode) (head *LNode) {
	if head1 == nil || head1.Next == nil {
		return head2
	} else if head2 == nil || head2.Next == nil {
		return head1
	}
	head = NewLinkedList()
	cur := head
	cur1 := head1.Next
	cur2 := head2.Next
	var val, forward int
	for cur1 != nil && cur2 != nil {
		val = cur1.Data + cur2.Data + forward
		cur.Next = &LNode{Data: val % 10}
		forward = val / 10
		cur = cur.Next
		cur1 = cur1.Next
		cur2 = cur2.Next
	}
	// 链表长度不一致情况
	for cur1 != nil {
		val = cur1.Data + forward
		cur.Next = &LNode{Data: val % 10}
		forward = val / 10
		cur = cur.Next
		cur1 = cur1.Next
	}
	for cur2 != nil {
		val = cur2.Data + forward
		cur.Next = &LNode{Data: val % 10}
		forward = val / 10
		cur = cur.Next
		cur2 = cur2.Next
	}
	// 最后的进位
	if forward == 1 {
		cur.Next = &LNode{Data: 1}
	}
	return head
}
