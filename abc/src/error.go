package main

import "fmt"

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

type Error interface {
	Caller() []CallerInfo
	Wrapped() []error
	Code() int
	error
}

func main() {
	var err Error
	fmt.Println(err)
}
