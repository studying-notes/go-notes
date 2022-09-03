package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		time.Sleep(1 * time.Second)

		fmt.Println("Goroutine 1 finished!")
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)

		fmt.Println("Goroutine 2 finished!")
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("All Goroutine finished!")
}
