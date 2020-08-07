/*
	Go 语言中的 json 包在序列化空接口存放的数字类
	型（整型、浮点型等）都序列化成 float64 类型;
	
	在 JSON 协议中是没有整型和浮点型之分的，它
	们统称为 number。JSON 字符串中的数字经过 Go 
	语言中的 JSON 包反序列化之后都会成为 float64 类型。
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type s struct {
	data map[string]interface{}
}

func main() {
	var s1 = s{data: make(map[string]interface{}, 8)}
	s1.data["count"] = 1
	ret, err := json.Marshal(s1.data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", string(ret)) // "{\"count\":1}"

	var s2 = s{data: make(map[string]interface{}, 8)}
	err = json.Unmarshal(ret, &s2.data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s2) // {map[count:1]}

	for _, v := range s2.data {
		fmt.Printf("value: %v, type:%T\n", v, v) // value: 1, type:float64
	}
}
