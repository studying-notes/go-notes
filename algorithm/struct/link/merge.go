package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	list1 := NewCustomLinkedList([]int{1, 3, 5, 7})
	PrintLNode(&list1)
	list2 := NewCustomLinkedList([]int{2, 4, 6, 8, 9})
	PrintLNode(&list2)

	//head := MergeOrderLinkedList(&list1, &list2)
	head := &LNode{Next: MergeSort(list1.Next, list2.Next)}
	PrintLNode(head)
}

// 归并排序法，不带头结点
func MergeSort(l1, l2 *LNode) (res *LNode) {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l2
	}
	if l1.Data < l2.Data {
		res = l1
		res.Next = MergeSort(l1.Next, l2)
	} else {
		res = l2
		res.Next = MergeSort(l2.Next, l1)
	}
	return res
}

// 合并有序链表
func MergeOrderLinkedList(head1, head2 *LNode) (head *LNode) {
	if head1 == nil || head1.Next == nil {
		return head2
	} else if head2 == nil || head2.Next == nil {
		return head1
	}
	head = head1 // 合并到链表 1
	cur := head
	cur1 := head1.Next
	cur2 := head2.Next
	head2.Next = nil // 清理回收
	for cur1 != nil && cur2 != nil {
		if cur1.Data < cur2.Data {
			cur.Next = cur1
			cur1 = cur1.Next
		} else {
			cur.Next = cur2
			cur2 = cur2.Next
		}
		cur.Next.Next = nil
		cur = cur.Next
	}
	// 链表长度不一致情况
	if cur1 != nil {
		cur.Next = cur1
	}
	if cur2 != nil {
		cur.Next = cur2
	}
	return head
}
