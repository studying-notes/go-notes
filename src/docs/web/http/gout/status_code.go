package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func GetStatusCode() {
	var code int
	err := gout.GET("example.com").
		Debug(true).Code(&code).Do()

	if err != nil || code != 200 {
		fmt.Printf("%s:code = %d\n", err, code)
		return
	}
}
