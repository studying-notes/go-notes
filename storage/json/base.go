package main

import (
	"encoding/json"
	"fmt"
)



type Box struct {
	Length  string  `json:"length"`
	Width   int64   `json:"width,omitempty"`
	Height  float64 `json:"-"`
	*Things `json:"things,omitempty"`
}

type Things struct {
	Weight []float64 `json:"weight"`
	Logo   string    `json:"logo"`
}

func main() {
	b1 := Box{
		Length: "120",
		Width:  75,
		Height: 4,
	}
	buf, _ := json.Marshal(b1)
	fmt.Printf("%s\n", buf)

	var b2 Box
	_ = json.Unmarshal(buf, &b2)
	fmt.Printf("%+v\n", b2)
}
