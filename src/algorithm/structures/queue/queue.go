package queue

import (
	"errors"

	. "algorithm/structures/link"
)

// 队列

// Queue 基于切片的实现
type Queue []interface{}

// IsEmpty 判断是否为空
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

// EnQueue 把新元素加入队列尾
func (q *Queue) EnQueue(val interface{}) {
	*q = append(*q, val)
}

// EnQueueLeft 把新元素加入队列首
func (q *Queue) EnQueueLeft(val interface{}) {
	*q = append(Queue{val}, *q...)
}

// DeQueue 弹出队列头元素
func (q *Queue) DeQueue() interface{} {
	if q.IsEmpty() {
		panic("empty queue")
	}
	val := (*q)[0]
	*q = (*q)[1:]
	return val
}

// First 获取队列头元素的值
func (q *Queue) First() interface{} {
	if q.IsEmpty() {
		panic("empty queue")
	}
	return (*q)[0]
}

// Pop 弹出队列尾元素
func (q *Queue) Pop() interface{} {
	if q.IsEmpty() {
		panic("empty queue")
	}
	val := (*q)[q.Size()-1]
	*q = (*q)[:q.Size()-1]
	return val
}

// Remove 根据值删除唯一元素
func (q *Queue) Remove(val interface{}) {
	for k, v := range *q {
		if v == val {
			*q = append((*q)[:k], (*q)[k+1:]...)
			return
		}
	}
}

// Size 队列大小
func (q *Queue) Size() int {
	return len(*q)
}

// SliceBasedQueue 数组实现的队列
type SliceBasedQueue struct {
	Data        []int
	front, rear int
}

func NewArrayQueue() *SliceBasedQueue {
	return &SliceBasedQueue{Data: []int{}, front: 0, rear: 0}
}

// IsEmpty 判断是否为空
func (q *SliceBasedQueue) IsEmpty() bool {
	return q.front == q.rear
}

// EnQueue 入队
func (q *SliceBasedQueue) EnQueue(val int) {
	q.rear++
	q.Data = append(q.Data, val)
}

// DeQueue 出队
func (q *SliceBasedQueue) DeQueue() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	val := q.Data[q.front]
	q.front++
	return val, nil
}

// Head 首元素
func (q *SliceBasedQueue) Head() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	return q.Data[q.front], nil
}

// Tail 尾元素
func (q *SliceBasedQueue) Tail() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	return q.Data[q.rear], nil
}

// Size 队列大小
func (q *SliceBasedQueue) Size() int {
	return q.rear - q.front
}

// LinkedQueue 链表实现的队列
type LinkedQueue struct {
	front, rear *LNode
}

func NewLinkedQueue() *LinkedQueue {
	node := &LNode{}
	return &LinkedQueue{front: node, rear: node}
}

// IsEmpty 判断是否为空
func (q *LinkedQueue) IsEmpty() bool {
	return q.front == q.rear
}

// EnQueue 入队
func (q *LinkedQueue) EnQueue(val int) {
	q.rear.Next = &LNode{Data: val}
	q.rear = q.rear.Next
}

// DeQueue 出队
func (q *LinkedQueue) DeQueue() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	q.front = q.front.Next
	return q.front.Data, nil
}
