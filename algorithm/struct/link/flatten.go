package main

import "fmt"

func main() {
	head := NewLinked2Node()
	PrintL2Node(head.Next)
	PrintL2Node(head.Next.Next)
	PrintL2Node(head.Next.Next.Next)
	PrintL2Node(head.Next.Next.Next.Next)
	//res := Merge(head.Next, head.Next.Next)
	res := Flatten(head)
	PrintL2Node(res)
}

func PrintL2Node(head *L2Node) {
	for cur := head; cur != nil; cur = cur.Down {
		fmt.Print(cur.Data, " ")
	}
	fmt.Println()
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

type L2Node struct {
	Data int
	Next *L2Node
	Down *L2Node
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
