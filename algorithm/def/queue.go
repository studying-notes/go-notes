package def

import "errors"

// ----- 队列 -----

// 基于 Go 切片数据结构的实现
type Queue []int

// 判断是否为空
func (q Queue) IsEmpty() bool {
	return len(q) == 0
}

// 把新元素加入队列尾
func (q *Queue) EnQueue(val int) {
	*q = append(*q, val)
}

// 把新元素加入队列首
func (q *Queue) EnQueueLeft(val int) {
	*q = append(Queue{val}, *q...)
}

// 弹出队列头元素
func (q *Queue) DeQueue() int {
	if q.IsEmpty() {
		panic("empty queue")
	}
	val := (*q)[0]
	*q = (*q)[1:]
	return val
}

// 获取队列头元素的值
func (q Queue) First() int {
	if q.IsEmpty() {
		panic("empty queue")
	}
	return q[0]
}

// 弹出队列尾元素
func (q *Queue) Pop() int {
	if q.IsEmpty() {
		panic("empty queue")
	}
	val := (*q)[q.Size()-1]
	*q = (*q)[:q.Size()-1]
	return val
}

// 根据值删除唯一元素
func (q *Queue) Remove(val int) {
	for k, v := range *q {
		if v == val {
			*q = append((*q)[:k], (*q)[k+1:]...)
			return
		}
	}
}

// 队列大小
func (q Queue) Size() int {
	return len(q)
}

// 数组实现的队列（不考虑并发操作）
type ArrayQueue struct {
	Data        []int
	front, rear int
}

func NewArrayQueue() *ArrayQueue {
	return &ArrayQueue{Data: []int{}, front: 0, rear: 0}
}

// 判断是否为空
func (q *ArrayQueue) IsEmpty() bool {
	return q.front == q.rear
}

// 入队
func (q *ArrayQueue) EnQueue(val int) {
	q.rear++
	q.Data = append(q.Data, val)
}

// 出队
func (q *ArrayQueue) DeQueue() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	val := q.Data[q.front]
	q.front++
	return val, nil
}

// 首元素
func (q *ArrayQueue) Head() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	return q.Data[q.front], nil
}

// 尾元素
func (q *ArrayQueue) Tail() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	return q.Data[q.rear], nil
}

// 队列大小
func (q *ArrayQueue) Size() int {
	return q.rear - q.front
}

// 链表实现的队列（不考虑并发操作）
type LinkedQueue struct {
	front, rear *LNode
}

func NewLinkedQueue() *LinkedQueue {
	node := &LNode{}
	return &LinkedQueue{front: node, rear: node}
}

// 判断是否为空
func (q *LinkedQueue) IsEmpty() bool {
	return q.front == q.rear
}

// 入队
func (q *LinkedQueue) EnQueue(val int) {
	q.rear.Next = &LNode{Data: val}
	q.rear = q.rear.Next
}

// 出队
func (q *LinkedQueue) DeQueue() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("empty queue")
	}
	q.front = q.front.Next
	return q.front.Data, nil
}
