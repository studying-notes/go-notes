package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var user = os.Getenv("USER_")
	defer func() {fmt.Println("main defer")}()
	go func() {
		defer func() {
			fmt.Println("defer here")
		}()

		if user == "" {
			panic("should set user env.")
		}

		defer func() {fmt.Println("defer2")}()
	}()

	time.Sleep(1 * time.Second)
	fmt.Printf("get result %s\r\n", user)
}
