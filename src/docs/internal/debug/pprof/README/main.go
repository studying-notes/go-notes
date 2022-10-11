package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go http.ListenAndServe(":6060", nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 1000; i++ {
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Duration(i) * time.Second)
				}

				if i == 999 {
					cancel()
				}
			}
		}(i)
	}

	<-ctx.Done()
}
