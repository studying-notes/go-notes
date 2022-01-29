package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/filter"
	"sync/atomic"
)

func main() {
	const (
		benchNumber     = 30000
		benchConcurrent = 30
	)

	i := int32(0)

	err := filter.NewBench().
		Concurrent(benchConcurrent).
		Number(benchNumber).
		Loop(func(c *gout.Context) error {

			// 下面的代码，每次生成不一样的http body 用于压测
			uid := uuid.New()
			id := atomic.AddInt32(&i, 1)

			c.POST(":8080").SetJSON(gout.H{"sid": uid.String(),
				"appkey": fmt.Sprintf("ak:%d", id),
				"text":   fmt.Sprintf("test text :%d", id)})
			return nil

		}).Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
}
