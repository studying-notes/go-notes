package link

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// LNode 定义链表结点
type LNode struct {
	Data int // 默认数组类型
	Next *LNode
}

// String 打印链表
func (l *LNode) String() (s string) {
	for cur := l.Next; cur != nil; cur = cur.Next {
		s += fmt.Sprintf("%d ", cur.Data)
	}
	return strings.TrimRight(s, " ")
}

// PrintLNode 打印链表
func PrintLNode(head *LNode) {
	fmt.Println(head)
}

// NewLinkedList 创建空链表
func NewLinkedList() *LNode {
	return &LNode{}
}

// NewCustomLinkedList 创建自定义链表
func NewCustomLinkedList(data []int) (head LNode) {
	cur := &head
	for d := range data {
		cur.Next = &LNode{Data: data[d]}
		cur = cur.Next
	}
	return head
}

// NewSeqLinkedList 创建有序链表
func NewSeqLinkedList(max int) (head LNode) {
	cur := &head
	for i := 1; i <= max; i++ {
		cur.Next = &LNode{Data: i}
		cur = cur.Next
	}
	return head
}

// NewNoSeqLinkedList 创建无序链表
func NewNoSeqLinkedList(max int) (head LNode) {
	cur := &head
	rand.Seed(time.Now().Unix())
	for range make([][]int, max) {
		cur.Next = &LNode{Data: rand.Intn(max)}
		cur = cur.Next
	}
	return head
}

// NewRingLinkedList 创建带环链表，k 表示环的位置
func NewRingLinkedList(max, k int) (head *LNode) {
	list := NewSeqLinkedList(max)
	head = &list
	fast := head
	slow := head
	for i := 1; i <= k && fast != nil; i++ {
		fast = fast.Next
	}
	// 链表长度小于 k 情况
	if fast == nil {
		fast = head.Next
		return head
	}
	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}
	// 将链表尾结点指向倒数第 k 个结点
	fast.Next = slow.Next
	return head
}
