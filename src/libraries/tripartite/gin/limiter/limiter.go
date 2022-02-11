package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	limiter := rate.NewLimiter(
		rate.Limit(1), // 每秒产生的数量
		10,            // 桶容量大小
	)
	for i := 0; i < 100; i++ {
		for limiter.Allow() {
			fmt.Println(i)
		}
		time.Sleep(time.Second)
	}
}
