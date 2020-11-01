package main

import (
	"fmt"
	"github.com/guonaihong/gout"
)

// PostJSONByMap map[string]interface{}
func PostJSONByMap() {

	var resp struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		OriginData  string `json:"origin_data"`
	}

	err := gout.POST("http://103.85.172.135:3004/exSymmDecStr").
		Debug(true).
		SetJSON(gout.H{
			"version":     "2",
			"authcode":    "ESIzRREiMRIxJREiM0URIjESMSURIjNFESIxEjElESIzRREiMRIxJQ==",
			"cipher_data": "bwUT0HpVdUa/AFZ8ardQ9Q7GtoTPDKgiEqUXThkPx/Fl2QGg6LjaRQvwkjhbRM/iEBQlxBWUfproPPf2+ZLnt4SiLFg0xuoOx01keuQiCgPzirbhuKxQZqgz/Y+qEwAmfZ2f7FxP0mPiy4+FGbAzINbxSDSN3Pq2PBOWMn1pEwc=",
			"alg_symm":    "SM4",
			"key":         "+b5xvu3br17XYCa0RLlAcg==",
			"mode":        "CBC",
			"padding":     "PKCS5PADDING",
			"iv_value":    "g+Ri5XBZy5pAZtu02b672Q==",
		}).
		BindJSON(&resp).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("%+v", resp)
}

// 结构体方式
func PostJSONByStruct() {
	type reqModel struct {
		Version  string `json:"version"`
		AuthCode string `json:"authcode"`
		RandLen  string `json:"randLen"`
	}

	var resp struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		RandData    string `json:"rand_data"`
	}

	err := gout.POST("http://103.85.172.135:3004/generateRandom").
		Debug(true).
		SetJSON(reqModel{
			Version:  "2",
			AuthCode: "ESIzRREiMRIxJREiM0URIjESMSURIjNFESIxEjElESIzRREiMRIxJQ==",
			RandLen:  "12",
		}).
		BindJSON(&resp).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("%+v", resp)
}

// JSON 字符串方式
func PostJSONByString() {
	json := `{
    "version": "2",
    "authcode": "ESIzRREiMRIxJREiM0URIjESMSURIjNFESIxEjElESIzRREiMRIxJQ==",
    "randLen": "16"
	}`

	var resp struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		RandData    string `json:"rand_data"`
	}

	err := gout.POST("http://103.85.172.135:3004/generateRandom").
		Debug(true).
		SetJSON(json).
		BindJSON(&resp).
		Do()

	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("%+v", resp)
}

func main() {
	//PostJSONByMap()
	//StructMethod()
	PostJSONByString()
}
