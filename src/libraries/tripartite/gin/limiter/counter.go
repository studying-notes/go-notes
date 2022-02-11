package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 普通计数限流器

type CountLimiter struct {
	counter  int64         // 计数器
	max      int64         // 最大数量
	interval time.Duration // 间隔时间
	last     time.Time     // 上一次时间
}

func NewCountLimiter(max int64, interval time.Duration) *CountLimiter {
	return &CountLimiter{
		counter:  0,
		max:      max,
		interval: interval,
		//last:     time.Now(), // 可能不必要
	}
}

func (c *CountLimiter) Allow() bool {
	current := time.Now()
	// 超过时间计数清零
	if current.After(c.last.Add(c.interval)) {
		atomic.StoreInt64(&c.counter, 1)
		c.last = current
		return true
	}
	// 取出一个
	atomic.AddInt64(&c.counter, 1)
	// 判断是否超过限流个数
	if c.counter <= c.max {
		return true
	}
	return false
}

func main() {
	limiter := NewCountLimiter(3, time.Second)
	for i := 0; i < 10; i++ {
		for limiter.Allow() {
			fmt.Println(i)
		}
		time.Sleep(time.Second)
	}
}
