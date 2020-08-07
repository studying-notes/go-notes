// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"

import (
	"fmt"
)

func main() {
	C.SayHello("Hello, World\n")
	fmt.Println("test")
}

//export SayHello
func SayHello(s string) {
	for range make([][]int, 10000) {
		fmt.Print(s)
	}
}
