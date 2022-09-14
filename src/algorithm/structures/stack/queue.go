package stack

// 用两个栈模拟队列操作

type QueueOver2Stack struct {
	en, de *Stack
}

func (q *QueueOver2Stack) IsEmpty() bool {
	return q.en.Size() == 0 && q.de.Size() == 0
}

func (q *QueueOver2Stack) EnQueue(val int) {
	q.en.Push(val)
}

func (q *QueueOver2Stack) DeQueue() int {
	if q.IsEmpty() {
		return 1 << 32
	}
	if q.de.IsEmpty() {
		for !q.en.IsEmpty() {
			q.de.Push(q.en.Pop())
		}
	}
	return q.de.Pop().(int)
}
