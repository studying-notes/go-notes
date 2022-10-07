package main

import "os"

func foo() error {
	var err *os.PathError
	return err
}

func main() {
	err := foo()
	println(err == nil)
}
