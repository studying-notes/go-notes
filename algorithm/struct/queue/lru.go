package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	q := &LRU{size: 3, queue: &Queue{}, set: NewHashSet()}
	q.AccessPage(1)
	q.AccessPage(2)
	q.AccessPage(3)
	q.AccessPage(2)
	q.AccessPage(5)
	q.AccessPage(6)
	q.PrintQueue()
}

func (q LRU) PrintQueue() {
	fmt.Println(*(q.queue))
}

type LRU struct {
	size  int // 缓存最大值
	queue *Queue
	set   *HashSet
}

// 判断缓存队列是否已满
func (q LRU) IsFull() bool {
	return q.queue.Size() == q.size
}

// 把页缓存到队首，同时添加到哈希表
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
