package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	var m1 = make(map[string]interface{}, 1)
	m1["count"] = 1
	buf, _ := json.Marshal(m1)
	fmt.Printf("%s\n", buf)  // {"count":1}

	var m2 map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(buf))
	decoder.UseNumber()
	_ = decoder.Decode(&m2)
	fmt.Printf("%T\n", m2["count"])  // json.Number
	count, _ := m2["count"].(json.Number).Int64()
	fmt.Printf("%T\n", int(count))  // int
}
