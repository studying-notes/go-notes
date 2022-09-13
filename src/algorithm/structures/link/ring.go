package link

// DetectRingHashSet HashSet 法
func DetectRingHashSet(head *LNode) (node *LNode, isRing bool) {
	if head == nil || head.Next == nil {
		return nil, false
	}
	// 简易 HashSet
	set := make(map[*LNode]bool)
	cur := head.Next
	for cur != nil {
		if set[cur] {
			return cur, true
		}
		set[cur] = true
		cur = cur.Next
	}
	return nil, false
}

// DetectRing2Pointer 快慢指针法
func DetectRing2Pointer(head *LNode) (node *LNode, isRing bool) {
	if head == nil || head.Next == nil {
		return nil, false
	}
	// 相同起始点
	fast := head.Next
	slow := head.Next

	for fast != nil && fast.Next != nil {
		if slow == fast.Next || slow == fast.Next.Next {
			return slow, true
		}
		slow = slow.Next
		fast = fast.Next.Next
	}
	return nil, false
}

// RingEntryNode 快慢指针法找出环的入口点
func RingEntryNode(head *LNode) *LNode {
	if head == nil || head.Next == nil {
		return nil
	}
	// 相同起始点
	fast := head.Next
	slow := head.Next
	for fast != nil && fast.Next != nil {
		if slow.Next == fast.Next.Next {
			break
		}
		slow = slow.Next
		fast = fast.Next.Next
	}

	cur := head.Next
	for cur != slow.Next {
		cur = cur.Next
		slow = slow.Next
	}
	return cur
}
