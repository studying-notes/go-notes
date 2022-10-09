package main

import (
	"fmt"
	"net/http"
	"sync"
)

func checkLink(link string, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, err := http.Get(link); err != nil {
		fmt.Println(link, "might be down!")
	} else {
		fmt.Println(link, "is up!")
	}
}

func main() {
	links := []string{
		"https://www.baidu.com/",
		"https://www.google.com/",
		"https://www.jd.com/",
		"https://www.taobao.com/",
		"https://www.tmall.com/",
		"https://www.sina.com.cn/",
		"https://www.sohu.com/",
		"https://www.163.com/",
	}

	wg := sync.WaitGroup{}
	for _, link := range links {
		wg.Add(1)
		go checkLink(link, &wg)
	}

	wg.Wait()
}
