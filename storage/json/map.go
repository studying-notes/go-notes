package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Parents []string `json:"parents"`
	Thing   `json:"thing"`
}

type Thing struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func main() {
	b := []byte(`{"type":"1","content":{"name":"Wetness","age":18,"parents":["Gomez","Morita"],"thing":{"length":1,"width":2,"height":3}}}`)

	model := make(map[string]interface{})
	_ = json.Unmarshal(b, &model)
	t := model["type"].(string) // 根据类型判断解析
	switch t {
	case "1":
		//p := Message{Content: Person{}}
		//_ = json.Unmarshal(b, &p)
		//fmt.Printf("%#v\n\n", p.Content)
		// 不是期望得到的
		// map[string]interface {}{"age":18, "name":"Wetness", "parents":[]interface {}{"Gomez", "Morita"}}

		p2 := struct {
			Type    string `json:"type"`
			Content Person `json:"content"`
		}{}
		_ = json.Unmarshal(b, &p2)
		fmt.Printf("%#v\n\n", p2.Content)
		// main.Person{Name:"Wetness", Age:18, Parents:[]string{"Gomez", "Morita"}, Thing:main.Thing{Length:1, Width:2, Height:3}}
	}
}
