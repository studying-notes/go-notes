//
// Created by Rustle Karl on 2020.11.18 15:59.
//

package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	v := []string{"1", "2", "3"}
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))

	var val []string
	err = json.Unmarshal(buf, &val)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
