package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Thing struct {
		Length int `json:"length"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Parents []string `json:"parents"`
		Thing   `json:"thing"`
	}

	person := Person{Name: "Wetness", Age: 18,
		Parents: []string{"Gomez", "Morita"},
		// 类型字段也可以用于赋值，不用定义变量
		Thing: Thing{2, 2, 2}}

	fmt.Printf("%#v\n\n", person)

	buf, _ := json.Marshal(person)
	fmt.Printf("%s\n", buf)
}
