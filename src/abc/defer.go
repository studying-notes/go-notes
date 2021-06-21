package main

import (
	"fmt"
)

func deferFuncReturn() (result int) {
	i := 1

	defer func() {
		result++
	}()

	return i
}

func main() {
	fmt.Println(deferFuncReturn())
}
