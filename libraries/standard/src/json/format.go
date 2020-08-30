package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Post struct {
	PostTime time.Time `json:"post_time"`
}

func main() {
	p1 := Post{
		PostTime: time.Now(),
	}
	buf, _ := json.Marshal(p1)
	fmt.Printf("%s\n", buf)

	s := `{"post_time":"2020-07-18 14:16:32"}`
	var p2 Post
	_ = json.Unmarshal([]byte(s), &p2)
	fmt.Printf("%+v\n", p2)
}
