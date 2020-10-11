package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type respData struct {
	resp *http.Response
}

func doCall(ctx context.Context) {
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: &transport,
	}
	respChan := make(chan *respData, 1)
	req, _ := http.NewRequest("GET", "http://google.com", nil)

	req = req.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		resp, _ := client.Do(req)
		fmt.Printf("client.do resp:%v, _:%v\n", resp, _)
		rd := &respData{
			resp: resp,
		}
		respChan <- rd
		wg.Done()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("call api timeout")
	case result := <-respChan:
		fmt.Println("call server api success")
		defer result.resp.Body.Close()
		data, _ := ioutil.ReadAll(result.resp.Body)
		fmt.Printf("resp:%v\n", string(data))
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	doCall(ctx)
}
