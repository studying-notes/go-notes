package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

func main() {
	var resp struct {
		ErrorCode int    `json:"errorCode"`
		Message   string `json:"message"`
		Result    bool   `json:"result"`
	}

	err := gout.POST("http://183.134.197.66:13027/device/platform/unit").
		Debug(true).
		SetHeader(gout.H{
			"Access-Token": "device",
		}).
		SetJSON(`{
    "areaCode": "310000",
    "code": "12330109470452028N",
    "contact": "苏颖",
    "name": "杭州市萧山区第二人民医院",
    "openingModule": "110",
    "parentID": 0,
    "address": "杭州市萧山区第二人民医院",
    "phone": "13522223020",
    "platform": 3,
    "status": 1,
    "ukey": "535748f27069c3167273"
}`).
		BindJSON(&resp).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("%+v", resp)
}
