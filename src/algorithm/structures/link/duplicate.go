package link

import . "algorithm/structures/set"

// RemoveDup 顺序删除
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

// RecursiveRemoveDupChild 递归法
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

// SetRemoveDup 利用 Set
func SetRemoveDup(head *LNode) {
	set := NewSet[int]()
	cur := head.Next
	for cur.Next != nil {
		if set.Contains(cur.Next.Data) {
			cur.Next = cur.Next.Next
			continue
		}
		set.Add(cur.Data)
		cur = cur.Next
	}
}

// RemoveDupSeq 从有序链表中移除重复项
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
