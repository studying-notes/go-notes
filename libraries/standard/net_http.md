---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "net/http - HTTP 标准库"  # 文章标题
url:  "posts/go/libraries/standard/net_http"  # 设置网页永久链接
tags: [ "go", "http", "web" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

- [HTTP 客户端](#http-客户端)
	- [GET](#get)
	- [POST](#post)
	- [POST FORM](#post-form)
- [自定义客户端](#自定义客户端)
- [自定义 Transport](#自定义-transport)

## HTTP 客户端

### GET

```go
func main() {
	var resp *http.Response
	var body []byte
	resp, _ = http.Get("http://www.baidu.com")
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

带参数的请求必须通过 `net/url` 标准库处理：

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	data := url.Values{}
	data.Set("key", "value")
	u, _:=url.ParseRequestURI(apiUrl)
	u.RawQuery = data.Encode()
	resp, _ = http.Get(u.String())
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

### POST

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	contentType := "application/json"
	data := `{"key":"value","age":18}`
	resp, _ = http.Post(apiUrl, contentType, strings.NewReader(data))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

### POST FORM

```go
func main() {
	var resp *http.Response
	var body []byte
	resp, _ = http.PostForm("http://www.baidu.com", url.Values{"key": {"value"}})
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

## 自定义客户端

管理 HTTP 客户端的 Headers、重定向策略等，必须创建一个 Client。

```go
func main() {
	var resp *http.Response
	var body []byte
	apiUrl := "http://www.baidu.com"
	client := &http.Client{CheckRedirect: nil}
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; "+
		"Ubuntu; Linux i686 on x86_64) AppleWebKit/537.36 (KHTML, "+
		"like Gecko) Chrome/64.0.3282.140 Safari/537.36 OPR/50.0.2762.67")
	resp, _ = client.Do(req)
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

## 自定义 Transport

管理代理、TLS 配置、`keep-alive`、压缩等，必须创建一个 Transport。

```go
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
```
