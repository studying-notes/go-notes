package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseModel struct {
	Result    bool        `json:"result"`
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func main() {
	var resp *http.Response
	var body []byte
	apiURL := "http://118.178.86.183:13027/device/platform/device/unit/bind"
	client := &http.Client{}
	body, _ = json.Marshal(gin.H{
		"unitCode":   "11111111111111",
		"deviceCode": "工器具柜-测试-zh",
	})
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Access-Token", "device")
	resp, _ = client.Do(req)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	var r responseModel
	_ = json.Unmarshal(body, &r)

	fmt.Println(r.ErrorCode)
	fmt.Println(string(body))
}
