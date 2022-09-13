package link

// 链表相邻元素翻转

// FlipAdjacentNode 交换相邻两个结点的数据域
func FlipAdjacentNode(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	cur := head.Next
	for cur != nil && cur.Next != nil {
		cur.Data, cur.Next.Data = cur.Next.Data, cur.Data
		cur = cur.Next.Next
	}
}

// FlipAdjPointer 交换相邻两个结点的指针域
func FlipAdjPointer(head *LNode) {
	if head == nil || head.Next == nil {
		return
	}
	pre := head
	cur := head.Next
	var suc *LNode
	for cur != nil && cur.Next != nil {
		suc = cur.Next
		cur.Next = cur.Next.Next
		suc.Next = cur
		pre.Next = suc
		pre = cur
		cur = cur.Next
	}
}

// 把链表以 k 个结点为一组进行翻转

func FlipAdjKNode(head *LNode, k int) {
	if head == nil || head.Next == nil {
		return
	}
	pre := head // 记录前驱
	cur1 := head.Next
	var cur2, suc *LNode
	for cur1 != nil {
		cur2 = cur1 // 记录首结点
		for i := 1; i < k && cur1 != nil; i++ {
			cur1 = cur1.Next
		}
		if cur1 != nil {
			suc = cur1.Next // 记住下一段的首结点
			cur1.Next = nil // 置空便于逆序
		}
		pre.Next = RecursiveReverseChild(cur2)
		if cur1 != nil {
			cur2.Next = suc // 首尾已交换，连接余下部分
			cur1 = suc      // 指向下一段的首结点
		}
		pre = cur2 // 下一段的前驱
	}
}
