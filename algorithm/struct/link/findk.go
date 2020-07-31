package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	list := NewSeqLinkedList(12)
	PrintLNode(&list)
	//nodeK := FindLastK(&list, 6)
	//fmt.Printf("%+v", nodeK)

	SpinLastK(&list, 6)
	PrintLNode(&list)
}

// 快慢指针法
func FindLastK(head *LNode, k int) *LNode {
	if head == nil || head.Next == nil {
		return nil
	}
	fast := head
	slow := head
	for i := 1; i <= k && fast != nil; i++ {
		fast = fast.Next
	}
	// 链表长度小于 k 情况
	if fast == nil {
		return nil
	}
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}
	return slow
}

func SpinLastK(head *LNode, k int) {
	if head == nil || head.Next == nil {
		return
	}
	fast := head
	slow := head
	for i := 1; i <= k && fast != nil; i++ {
		fast = fast.Next
	}
	// 链表长度小于 k 情况
	if fast == nil {
		return
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	defer func() {
		slow.Next = nil
	}()
	fast.Next = head.Next
	head.Next = slow.Next
}
