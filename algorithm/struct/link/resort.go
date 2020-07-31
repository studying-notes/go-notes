package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	// -> 0 9 1 8 2 7 3 6 4 5
	// -> 0 8 1 7 2 6 3 5 4
	list := NewCustomLinkedList([]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	PrintLNode(&list)
	ReverseResort(&list)
	PrintLNode(&list)
	//LoopResort(&list)
	//PrintLNode(&list)
}

// 对链表进行重新排序

// 分割逆序组合法
func ReverseResort(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	// 分割成两部分，对第二部分逆序
	part1 := head.Next
	part2 := ReverseChild(Split2Parts(head))
	var cur *LNode // 保存拆下来的结点
	for part2 != nil {
		cur = part2
		part2 = part2.Next
		cur.Next = part1.Next
		part1.Next = cur
		part1 = part1.Next.Next
	}
}

// 2 倍速指针找中间结点，然后截断链表
// 注意返回的链表不带头结点
func Split2Parts(head *LNode) *LNode {
	if head == nil || head.Next == nil {
		return nil
	}
	// 将头结点作为相同起点
	cur := head
	fast := head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next // 2 倍速指针
		cur = cur.Next
	}
	// 最后清除前驱指针
	defer func() {
		cur.Next = nil
	}()
	return cur.Next
}

// from reverse.go
func ReverseChild(node *LNode) *LNode {
	if node == nil || node.Next == nil {
		return node
	}
	head := ReverseChild(node.Next)
	node.Next.Next = node
	node.Next = nil
	return head
}

// 两层循环，每次都找最后一个结点
func LoopResort(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
	var pre, tail *LNode // 最后一个结点及其前驱
	for cur != nil && cur.Next != nil {
		pre, tail = cur, cur.Next
		for tail.Next != nil {
			pre, tail = tail, tail.Next // 编译器会先计算再赋值
		}
		pre.Next = nil
		tail.Next = cur.Next
		cur.Next = tail
		cur = cur.Next.Next
	}
}
