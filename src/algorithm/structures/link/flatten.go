package link

import (
	"fmt"
	"strings"
)

// 展开二维链表

type L2Node struct {
	Data int
	Next *L2Node
	Down *L2Node
}

func (l *L2Node) String() (s string) {
	for cur := l; cur != nil; cur = cur.Down {
		s += fmt.Sprintf("%d ", cur.Data)
	}
	return strings.TrimRight(s, " ")
}

func NewLinked2List(data []int) (head L2Node) {
	cur := &head
	for d := range data {
		cur.Down = &L2Node{Data: data[d]}
		cur = cur.Down
	}
	return *head.Down
}

func NewLinked2Node() *L2Node {
	head := &L2Node{}
	cur := head
	list1 := NewLinked2List([]int{6, 8, 31})
	cur.Next = &L2Node{Data: 3, Down: &list1}
	cur = cur.Next
	list2 := NewLinked2List([]int{21})
	cur.Next = &L2Node{Data: 11, Down: &list2}
	cur = cur.Next
	list3 := NewLinked2List([]int{22, 50})
	cur.Next = &L2Node{Data: 15, Down: &list3}
	cur = cur.Next
	list4 := NewLinked2List([]int{39, 40, 55})
	cur.Next = &L2Node{Data: 30, Down: &list4}
	return head
}

func PrintL2Node(head *L2Node) {
	fmt.Println(head)
}

func Flatten(head *L2Node) *L2Node {
	if head == nil || head.Next == nil {
		return head
	}
	cur := head.Next
	res := head.Next
	var suc *L2Node
	for cur != nil && cur.Next != nil {
		suc = cur.Next
		res = Merge(res, cur.Next)
		cur = suc
	}
	return res
}

func Merge(l1, l2 *L2Node) (res *L2Node) {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l2
	}
	if l1.Data < l2.Data {
		res = l1
		res.Down = Merge(l1.Down, l2)
	} else {
		res = l2
		res.Down = Merge(l2.Down, l1)
	}
	return res
}
