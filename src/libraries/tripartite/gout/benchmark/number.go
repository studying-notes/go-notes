package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {
	benchNumber := 1000000   // 压测次数
	benchConcurrent := 60000 // 并发数

	err := gout.
		POST(":18888").
		SetJSON(gout.H{"hello": "world"}).
		Filter().
		Bench().
		Concurrent(benchConcurrent).
		Number(benchNumber).
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
