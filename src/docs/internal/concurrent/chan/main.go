package main

import (
	"fmt"
	"time"
)

func Process(ch chan int) {
	time.Sleep(time.Second)

	ch <- 1
}

func main() {
	channels := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		channels[i] = make(chan int)
		go Process(channels[i])
	}

	for i, ch := range channels {
		<-ch
		fmt.Println("Routine ", i, " quit!")
	}
}
