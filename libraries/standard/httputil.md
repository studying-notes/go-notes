---
date: 2022-01-29T18:46:11+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Golang 网络包 httputil"  # 文章标题
# description: "文章描述"
url:  "posts/go/libraries/standard/httputil"  # 设置网页永久链接
tags: [ "go", "http", "io" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 打印网络请求

这个工具一般用在调试时，它的执行效率不高。

### 打印服务端收到的请求

Golang 提供了一个 DumpRequest 工具可以用来输出请求的内容。

```go
import (
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true) // 第二个参数是否显示 body
		log.Printf("代理请求数据: %s", dump)

		/* 设置应答头 */
		w.Header().Set("Content-Type", "application/json")

		/* 设置状态码 */
		w.WriteHeader(201)

		/* 设置应答体 */
		w.Write([]byte("{}"))
	})

	http.ListenAndServe(":1280", nil)
}
```

请求

```shell
curl -d'login=emma＆password=123' -X POST http://localhost:1280
```

Golang 控制台输出

```shell
2021/12/20 14:18:01 代理请求数据: POST / HTTP/1.1
Host: localhost:1280
Accept: */*
Content-Length: 25
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.64.0

login=emma＆password=123
```

### 打印客户端请求

上面的 DumpRequest 用作打印服务端收到的的 Request,当需要自己去请求时则通过这个 DumpRequestOut 打印

```go
func main() {
	req, err := http.NewRequest(http.MethodGet, "https://www.google.com", nil)
	if err != nil {
		return
	}
	requestDump, _ := httputil.DumpRequestOut(req, false)
	log.Printf("打印请求头: %s", requestDump)
}
```

Golang 控制台输出

```shell
打印请求头: GET / HTTP/1.1
Host: www.google.com
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip
```

DumpRequest 与 DumpRequestOut 的区别在于：

- 前者仅打印 Request 当前存在的值
- 前者实际上到最终发出请求还会自动添加一些字段
- 后者打印 Request 发出请求前最终的值

### 打印客户端响应

```go
func main() {
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dump, _ := httputil.DumpResponse(resp, false)
	log.Printf("打印响应头: %s", dump)
}
```

Golang 控制台输出

```shell
2021/12/20 14:58:35 打印响应头: HTTP/1.1 200 OK
Content-Length: 227
Accept-Ranges: bytes
Cache-Control: no-cache
Connection: keep-alive
Content-Type: text/html
Date: Mon, 20 Dec 2021 06:58:35 GMT
P3p: CP=" OTI DSP COR IVA OUR IND COM "
P3p: CP=" OTI DSP COR IVA OUR IND COM "
Pragma: no-cache
Server: BWS/1.1
Set-Cookie: BD_NOT_HTTPS=1; path=/; Max-Age=300
Set-Cookie: BIDUPSID=90CA49C500F16CCDF8320520F1A2093E; expires=Thu, 31-Dec-37 23:55:55 GMT; max-age=2147483647; path=/; domain=.baidu.com
Set-Cookie: PSTM=1639983515; expires=Thu, 31-Dec-37 23:55:55 GMT; max-age=2147483647; path=/; domain=.baidu.com
Set-Cookie: BAIDUID=90CA49C500F16CCD83C727248FF791B7:FG=1; max-age=31536000; expires=Tue, 20-Dec-22 06:58:35 GMT; domain=.baidu.com; path=/; version=1; comment=bd
Strict-Transport-Security: max-age=0
Traceid: 1639983515045629082614963470251705031207
X-Frame-Options: sameorigin
X-Ua-Compatible: IE=Edge,chrome=1
```

## 反向代理

正向代理：隐藏真正的客户端

反向代理：隐藏后端服务器

### 简单使用例

编写一个被代理的服务

```go
package main

import (
	"log"
	"net/http"
)

// 这里创建一个类型是为了实现 Handler 接口
type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!\n"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:7070", &s)
}
```

启动这个服务后

```shell
curl http://127.0.0.1:7070
```

再编写一个反向代理服务

```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	// 返回一个单主机代理对象
	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	// 返回一个代理方法
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// 初始化反向代理并传入真正后端服务的地址（被代理的服务器）
	proxy, err := NewProxy("http://127.0.0.1:7070")
	if err != nil {
		panic(err)
	}

	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

这时访问 `:8080` 就能代理到 `:7070` 那里去了

```shell
curl http://127.0.0.1:8080
```

关键的代码就是 NewSingleHostReverseProxy 这个方法，该方法返回了一个 ReverseProxy 对象，在 ReverseProxy 中的 ServeHTTP 方法实现了这个具体的过程，主要是对源 http 包头进行重新封装，而后发送到后端服务器。

### 修改响应

修改响应，例如添加一个请求头字段

```go
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = modifyResponse()
	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "Magical")
		return nil
	}##
}
```

这个 ReverseProxy 提供了一个用来修改响应的执行

在 modifyResponse 中，可以返回一个错误（如果你在处理响应发生了错误）， 如果你设置了 proxy.ErrorHandler, modifyResponse 返回错误时会自动调用 ErrorHandler 进行错误处理。

```go
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil##
}

// 抛出一个异常
func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return errors.New("response body is invalid")
	}
}

// 捕获异常，这里可以进行错误处理
func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		// 重新把请求响应回去
		w.Write([]byte(err.Error() + "\n"))
	}
}
```

响应结果：

```sh
curl http://127.0.0.1:8080
```

### 修改请求

```go
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
    url, err := url.Parse(targetHost)
    if err != nil {
        return nil, err
    }

    proxy := httputil.NewSingleHostReverseProxy(url)

    originalDirector := proxy.Director // 原本执行的流程

    proxy.Director = func(req *http.Request) {
        originalDirector(req) // 需要执行原本的代理流程
        modifyRequest(req)
    }

    proxy.ModifyResponse = modifyResponse()
    proxy.ErrorHandler = errorHandler()
    return proxy, nil
}

func modifyRequest(req *http.Request) {
    req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}
```
