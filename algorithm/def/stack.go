package def

import (
	"errors"
)

// ----- 栈 -----

// 基于 Go 切片数据结构的实现

type Stack []int

func (s Stack) IsEmpty() bool {
	return len(s) == 0
}

func (s Stack) Size() int {
	return len(s)
}

func (s Stack) Top() int {
	if s.IsEmpty() {
		panic("empty stack")
	}
	return s[len(s)-1]
}

func (s *Stack) Push(val int) {
	*s = append(*s, val)
}

func (s *Stack) Pop() int {
	if s.IsEmpty() {
		panic("empty stack")
	}
	val := s.Top()
	*s = (*s)[:s.Size()-1]
	return val
}

// 数组实现的栈（不考虑并发操作）
type ArrayStack struct {
	Data []int // 存储数据
	Size int   // 大小
}

// 判断是否为空
func (s *ArrayStack) IsEmpty() bool {
	return s.Size == 0
}

// 返回栈的长度
func (s *ArrayStack) Len() int {
	return s.Size
}

// 返回栈顶元素
func (s *ArrayStack) Top() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	return s.Data[s.Size-1], nil
}

// 弹出栈元素
func (s *ArrayStack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	s.Size--
	return s.Data[s.Size], nil
}

// 添加元素
func (s *ArrayStack) Push(val int) {
	s.Data = append(s.Data, val)
	s.Size++
}

// 链表实现的栈（不考虑并发操作）
type LinkedStack struct {
	Head *LNode
	// 用链表头结点存储长度数值
}

func NewLinkedStack() *LinkedStack {
	return &LinkedStack{Head: &LNode{Data: 0}}
}

// 判断是否为空
func (s *LinkedStack) IsEmpty() bool {
	return s.Head.Data == 0
}

// 返回栈的长度
func (s *LinkedStack) Len() int {
	return s.Head.Data
}

// 返回栈顶元素
func (s *LinkedStack) Top() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	return s.Head.Next.Data, nil
}

// 弹出栈元素
func (s *LinkedStack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	val := s.Head.Next
	s.Head.Next = s.Head.Next.Next
	s.Head.Data--
	val.Next = nil // 垃圾回收
	return val.Data, nil
}

// 添加元素
func (s *LinkedStack) Push(val int) {
	node := &LNode{Data: val, Next: s.Head.Next}
	s.Head.Next = node
	s.Head.Data++
}
