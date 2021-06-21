package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var datas []string

func main() {
	go func() {
		for {
			log.Println(Add("pprof"))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	_ = http.ListenAndServe(":6060", nil)
}

func Add(s string) int {
	data := []byte(s)
	datas = append(datas, string(data))
	return len(datas)
}
