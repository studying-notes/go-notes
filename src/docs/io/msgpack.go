/*
	MessagePack 是一种高效的二进制序列化格式
*/

package main

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"log"
)

type Person struct {
	Name   string
	Age    int
	Gender string
}

func main() {
	p1 := Person{
		Name:   "萧瑟",
		Age:    18,
		Gender: "男",
	}

	// marshal
	b, err := msgpack.Marshal(p1)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal
	var p2 Person
	err = msgpack.Unmarshal(b, &p2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("p2:%#v\n", p2) // p2:main.Person{Name:"萧瑟", Age:18, Gender:"男"}
}
