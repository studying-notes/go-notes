package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"time"
)

const (
	benchTime       = 4 * time.Second
	benchConcurrent = 3000
)

func main() {
	err := gout.
		POST(":18888").
		SetJSON(gout.H{"hello": "world"}).
		Filter().
		Bench().
		Concurrent(benchConcurrent).
		Durations(benchTime).
		Do()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
