/*
	标准库 gob 是 golang 提供的“私有”的编解码方式，它的效率
	会比 json，xml 等更高，特别适合在 Go 语言程序间传递数据
*/

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type s struct {
	data map[string]interface{}
}

func main() {
	var s1 = s{data: make(map[string]interface{}, 8)}
	s1.data["count"] = 1

	// encode
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s1.data)
	if err != nil {
		log.Fatal(err)
	}
	res := buf.Bytes()
	fmt.Println(res) // 一串字节数组

	var s2 = s{data: make(map[string]interface{}, 8)}
	// decode
	dec := gob.NewDecoder(bytes.NewBuffer(res))
	err = dec.Decode(&s2.data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s2.data) // map[count:1]

	for _, v := range s2.data {
		fmt.Printf("value: %v, type:%T\n", v, v) // value: 1, type:int
	}
}
