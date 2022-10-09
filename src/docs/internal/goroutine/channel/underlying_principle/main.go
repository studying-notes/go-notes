package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println(<-ch)
		fmt.Println("1")
	}()

	go func() {
		defer wg.Done()
		fmt.Println(<-ch)
		fmt.Println("2")
	}()

	ch <- 3

	wg.Wait()
}
