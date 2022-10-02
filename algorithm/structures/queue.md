---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数据结构与算法之队列"  # 文章标题
url:  "posts/go/algorithm/structures/queue"  # 设置网页永久链接
tags: [ "algorithm", "go" ]  # 标签
categories: [ "Go 数据结构与算法"]  # 系列

weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

- [数据结构](#数据结构)
	- [数组实现的队列](#数组实现的队列)
	- [链表实现的队列](#链表实现的队列)
	- [用切片简单实现](#用切片简单实现)
- [设计一个排序系统](#设计一个排序系统)
- [实现 LRU 缓存方案](#实现-lru-缓存方案)

## 数据结构

> 源码位置 *src/algorithm/structures/queue/queue.go*

### 数组实现的队列

```go
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
```

### 链表实现的队列

```go
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
```

### 用切片简单实现

```go
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
```

## 设计一个排序系统

设计一个排队系统，能够让每个进入队伍的用户都能看到自己在队列中所处的位置和变化，队伍可能随时有人加入和退出，加入必须在队尾，但是退出可以在任意位置；当有人退出影响到用户的位置排名时需要及时反馈到用户。

> 源码位置 *src/algorithm/structures/queue/wait.go*

```go
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
```

## 实现 LRU 缓存方案

LRU 是 Least Recently Used 的缩写，它的意思是“**最近最少使用**”，LRU 1 缓存就是使用这种原理实现，简单地说就是缓存一定量的数据，当超过设定的阈值时就把一些过期的数据删除掉。

常用于页面置换算法，是为虚拟页式存储管理中常用的算法。

我们可以使用两个数据结构实现一个 LRU 缓存。

1. 使用双向链表实现的队列，队列的最大的容量为缓存的大小。在使用的过程中，**把最近使用的页面移动到队列首**，最近没有使用的页面将被放在队列尾的位置。

2. 使用一个哈希表，**把页号作为键，把缓存在队列中的结点的地址作为值**。

当引用一个页面时，所需的页面在内存中，我们需要把这个页对应的结点移动到队列的前面。如果所需的页面不在内存中，我们将它存储在内存中。

简单地说，就是**将一个新结点添加到队列的前面，并在哈希表中更新相应的结点地址。如果队列是满的，那么就从队列尾部移除一个结点，并将新结点添加到队列的前面。**

```go
type LRU struct {
	size  int // 缓存最大值
	queue *Queue
	set   *Set
}

// IsFull 判断缓存队列是否已满
func (q *LRU) IsFull() bool {
	return q.queue.Size() == q.size
}

// EnQueueLeft 把页缓存到队首，同时添加到哈希表
func (q *LRU) EnQueueLeft(page int) {
	// 队列满了就删除队尾缓存的页
	if q.IsFull() {
		q.set.Remove(q.queue.Pop())
	}
	q.queue.EnQueueLeft(page)
	// 同时添加到哈希表
	q.set.Add(page)
}

func (q *LRU) AccessPage(page int) {
	if !q.set.Contains(page) {
		q.EnQueueLeft(page)
	} else if page != q.queue.First() {
		q.queue.Remove(page)
		q.queue.EnQueueLeft(page)
	}
}

func (q *LRU) PrintQueue() {
	fmt.Println(*(q.queue))
}
```
