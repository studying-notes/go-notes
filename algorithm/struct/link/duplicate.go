package main

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	LinkedList := NewNoSeqLinkedList(12)
	head := &LinkedList
	PrintLNode(head)
	//RemoveDup(head)
	//RecursiveRemoveDupChild(head.Next)
	SetRemoveDup(head)
	PrintLNode(head)
}

// 利用 HashSet (map)
func SetRemoveDup(head *LNode) {
	set := NewHashSet()
	cur := head.Next
	for cur.Next != nil {
		if set.Get(cur.Next.Data) {
			cur.Next = cur.Next.Next
			continue
		}
		set.Add(cur.Data)
		cur = cur.Next
	}
}

// 递归法
// 不带头结点删除
func RecursiveRemoveDupChild(node *LNode) *LNode {
	if node == nil || node.Next == nil {
		return node
	}
	RecursiveRemoveDupChild(node.Next)
	cur := node
	for cur.Next != nil {
		if node.Data == cur.Next.Data {
			cur.Next = cur.Next.Next
			continue
		}
		cur = cur.Next
	}
	return node
}

// 顺序删除
func RemoveDup(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	outerCur := head.Next // 外层循环
	var innerCur *LNode   // 内层循环
	for outerCur != nil && outerCur.Next != nil {
		innerCur = outerCur
		// 不记录前驱而是判断后继
		for innerCur.Next != nil {
			if innerCur.Next.Data == outerCur.Data {
				innerCur.Next = innerCur.Next.Next // 将重复数据结点短路
				continue                           // 可能不止一个相同数据结点
			}
			innerCur = innerCur.Next
		}
		outerCur = outerCur.Next
	}
}

func RemoveDupSeq(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
	for cur.Next != nil {
		if cur.Next.Data == cur.Data {
			cur.Next = cur.Next.Next // 将重复数据结点短路
			continue                 // 可能不止一个相同数据结点
		}
		cur = cur.Next
	}
}
