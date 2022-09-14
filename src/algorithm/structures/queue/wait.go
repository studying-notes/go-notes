package queue

import "fmt"

// 设计一个排序系统

func NewSeqQueue(max int) *WaitQueue {
	wq := NewWaitQueue()
	for i := 1; i <= max; i++ {
		wq.EnQueueVal(i)
	}
	return wq
}

func PrintQueue(q *WaitQueue) {
	if q == nil || q.rear.no == 0 {
		return
	}
	cur := q.front.next
	for cur != nil {
		fmt.Printf("No: %d, Val: %d\n", cur.no, cur.val)
		cur = cur.next
	}
}

type Node struct {
	val      int // 值
	no       int // 位置序号
	previous *Node
	next     *Node
}

func (n *Node) Position() int {
	return n.no
}

func (n *Node) Value() int {
	return n.val
}

func (n *Node) SetPosition(no int) {
	n.no = no
}

func (n *Node) Leave() {
	if n.previous == nil {
		panic("the head node can't leave the queue")
	}
	if n.next == nil {
		n.previous.next = nil
	} else {
		n.previous.next = n.next
		n.next.previous = n.previous
		cur := n.next.next
		for cur != nil {
			cur.no--
			cur = cur.next
		}
	}
}

type WaitQueue struct {
	front *Node // 队首
	rear  *Node // 队尾
}

func NewWaitQueue() *WaitQueue {
	node := &Node{no: 0}
	node.next = node
	node.previous = node
	return &WaitQueue{front: node, rear: node}
}

// IsEmpty 判断是否为空
func (q *WaitQueue) IsEmpty() bool {
	return q.rear.no == 0
}

// EnQueueVal 入队
func (q *WaitQueue) EnQueueVal(val int) {
	q.rear.next = &Node{val: val, no: q.rear.no + 1, previous: q.rear}
	q.rear = q.rear.next
}

func (q *WaitQueue) EnQueueValues(values ...int) {
	for i := range values {
		q.EnQueueVal(values[i])
	}
}

// EnQueueNode 以结点入队
func (q *WaitQueue) EnQueueNode(node *Node) {
	node.no = q.rear.no + 1
	node.previous = q.rear
	node.next = nil
	q.rear.next = node
	q.rear = q.rear.next
}

// DeQueue 排到出队
func (q *WaitQueue) DeQueue() int {
	if q.IsEmpty() {
		panic("empty queue")
	}
	val := q.front.next.val
	q.front.next.Leave()
	return val
}

func (q *WaitQueue) DeQueueNode(node *Node) {
	node.Leave()
}

// Size 队列大小
func (q *WaitQueue) Size() int {
	return q.rear.no
}
