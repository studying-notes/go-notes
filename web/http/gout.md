---
date: 2020-10-30T14:38:37+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go HTTP 客户端 gout"  # 文章标题
# description: "文章描述"
url:  "posts/go/web/http/gout"  # 设置网页永久链接
tags: [ "go", "http", "web" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## POST 请求带 JSON

### 映射方式

```go
// MapMethod map[string]interface{}
func MapMethod() {

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
```

### 结构体方式

```go
func StructMethod() {
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
```

### JSON 字符串方式

```go
func JSONStringMethod() {
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
```

## GET 请求带查询字符串

### 映射方式

```go
func QueryByMap() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(gout.H{
			"name":     "user",
			"age":      18,
			"weight":   50.4,
			"birthday": time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 数组方式

```go
func QueryByArray() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(gout.A{
			"name", "user",
			"age", 18,
			"weight", 50.4,
			"birthday", time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 结构体方式

```go
func QueryByStruct() {
	err := gout.GET("example.com").
		Debug(true).
		SetQuery(Query{
			Name:     "user",
			Age:      18,
			Weight:   50.4,
			Birthday: time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### JSON 字符串方式

```go
func QueryByString() {
	err := gout.GET("example.com").Debug(true).
		SetQuery("name=user&age=18&weight=50.5&birthday=2020-1-20").Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

## 设置请求 Header

### 映射方式

```go
func SetHeaderByMap() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(gout.H{
			"name":     "user",
			"age":      18,
			"weight":   50.4,
			"birthday": time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 数组方式

```go
func SetHeaderByArray() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(gout.A{
			"name", "user",
			"age", 18,
			"weight", 50.4,
			"birthday", time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 结构体方式

```go
func SetHeaderByStruct() {
	err := gout.GET("example.com").
		Debug(true).
		SetHeader(Header{
			Name:     "user",
			Age:      18,
			Weight:   50.4,
			Birthday: time.Now(),
		}).
		Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 绑定方式

```go
type respHeader struct {
	Total int       `header:"total"`
	Sid   string    `header:"sid"`
	Time  time.Time `header:"time"`
}

func bindHeader() {
	resp := respHeader{}
	err := gout.GET("example.com").Debug(true).BindHeader(&resp).Do()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
```

## GET 请求设置 body

### 字符串内容

```go
func SetBodyByString() {
	err := gout.POST("example.com").Debug(true).SetBody("string").Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### io.Reader

```go
func SetBodyByReader() {
	err := gout.POST("example.com").Debug(true).SetBody(strings.NewReader("io.Reader")).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 基础类型

```go
func SetBodyByBaseType() {
	err := gout.POST("example.com").Debug(true).SetBody(3.14).Do()

	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
```

### 结构体方式
```go

```

### JSON 字符串方式

```go

```

```go

```

## GET 请求带查询字符串
### 映射方式
### 结构体方式
### JSON 字符串方式

```go

```

```go

```
## GET 请求带查询字符串
### 映射方式
### 结构体方式
### JSON 字符串方式

```go

```

```go

```
## GET 请求带查询字符串
### 映射方式
### 结构体方式
### JSON 字符串方式

```go

```

```go

```
```go

```

```go

```

```go

```


```go

```

```go

```

```go

```

```go

```

```go

```


