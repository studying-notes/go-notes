package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

type respHeader struct {
	Total int       `header:"total"`
	Sid   string    `header:"sid"`
	Time  time.Time `header:"time"`
}

func bindHeader() {
	resp := respHeader{}
	err := gout.GET("example.com").Debug(true).BindHeader(&resp).Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
