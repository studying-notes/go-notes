package main

import "encoding/json"

type Flags struct {
	A string "json:\"a\""
	B string "json:\"b\""
}

func main() {
	f := Flags{}
	f.A = "a"
	f.B = "b"
	bytes, err := json.Marshal(f)
	if err != nil {
		return
	}
	println(string(bytes))
}
