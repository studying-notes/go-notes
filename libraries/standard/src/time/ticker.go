package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	ticker := time.NewTicker(2 * time.Second)
	timer := time.NewTimer(2 * time.Second)

	go func(t *time.Ticker) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("Contains Ticker", time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker)

	go func(t *time.Timer) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("Contains Timer", time.Now().Format("2006-01-02 15:04:05"))
			t.Reset(2 * time.Second) // 重置定时器，可实现类似 Ticker 的功能
		}
	}(timer)

	wg.Wait()
}
