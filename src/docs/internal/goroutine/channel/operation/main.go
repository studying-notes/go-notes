package main

import "fmt"

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			select {
			case a <- i:
				a = nil
				b = make(chan int)
			case b <- i:
				b = nil
				a = make(chan int)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		select {
		case v := <-a:
			fmt.Println("a:", v)
		case v := <-b:
			fmt.Println("b:", v)
		}
	}
}
