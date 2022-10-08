package main

import (
	"fmt"
	"reflect"
)

type People struct {
	Age   int
	Name  string
	Test1 string
	Test2 string
}

func NewPeople() *People {
	return &People{
		Age:   18,
		Name:  "Karl",
		Test1: "Test1",
		Test2: "Test2",
	}
}

func NewPeopleByReflect() interface{} {
	var people People
	t := reflect.TypeOf(people)
	v := reflect.New(t)
	v.Elem().Field(0).SetInt(18)
	v.Elem().Field(1).SetString("Karl")
	v.Elem().Field(2).SetString("Test1")
	v.Elem().Field(3).SetString("Test2")
	return v
}

func main() {
	p1 := NewPeople()
	p2 := NewPeopleByReflect().(*People)

	fmt.Println(p1)
	fmt.Println(p2)
}
