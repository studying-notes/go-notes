package stack

import (
	"errors"
	"fmt"

	. "algorithm/structures/link"
)

// 栈

// SliceBasedStack 基于切片实现的栈
type SliceBasedStack struct {
	Data []int // 存储数据
	Size int   // 大小
}

// IsEmpty 判断是否为空
func (s *SliceBasedStack) IsEmpty() bool {
	return s.Size == 0
}

// Len 返回栈的长度
func (s *SliceBasedStack) Len() int {
	return s.Size
}

// Top 返回栈顶元素
func (s *SliceBasedStack) Top() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	return s.Data[s.Size-1], nil
}

// Pop 弹出栈元素
func (s *SliceBasedStack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	s.Size--
	return s.Data[s.Size], nil
}

// Push 添加元素
func (s *SliceBasedStack) Push(val int) {
	s.Data = append(s.Data, val)
	s.Size++
}

// LinkedListBasedStack 基于链表实现的栈
type LinkedListBasedStack struct {
	Head *LNode // 用链表头结点存储长度数值
}

func NewLinkedStack() *LinkedListBasedStack {
	return &LinkedListBasedStack{Head: &LNode{Data: 0}}
}

// IsEmpty 判断是否为空
func (s *LinkedListBasedStack) IsEmpty() bool {
	return s.Head.Data == 0
}

// Len 返回栈的长度
func (s *LinkedListBasedStack) Len() int {
	return s.Head.Data
}

// Top 返回栈顶元素
func (s *LinkedListBasedStack) Top() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	return s.Head.Next.Data, nil
}

// Pop 弹出栈元素
func (s *LinkedListBasedStack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty stack")
	}
	val := s.Head.Next
	s.Head.Next = s.Head.Next.Next
	s.Head.Data--
	val.Next = nil // 垃圾回收
	return val.Data, nil
}

// Push 添加元素
func (s *LinkedListBasedStack) Push(val int) {
	node := &LNode{Data: val, Next: s.Head.Next}
	s.Head.Next = node
	s.Head.Data++
}

// Stack 基于切片的实现
type Stack []interface{}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Size() int {
	return len(*s)
}

func (s *Stack) Top() interface{} {
	if s.IsEmpty() {
		panic("empty stack")
	}
	return (*s)[len(*s)-1]
}

func (s *Stack) Push(val interface{}) {
	*s = append(*s, val)
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		panic("empty stack")
	}
	val := s.Top()
	*s = (*s)[:s.Size()-1]
	return val
}

func (s *Stack) String() string {
	return fmt.Sprintf("%s", *s)
}
