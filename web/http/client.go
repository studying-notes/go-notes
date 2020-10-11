package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: nil},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}
	resp, _ = client.Get(apiUrl)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
