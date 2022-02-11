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
- [Cookie 的定义](#cookie-的定义)
	- [设置 Cookie](#设置-cookie)
	- [获取 Cookie](#获取-cookie)
- [Gin 框架操作 Cookie](#gin-框架操作-cookie)
	- [Session](#session)

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

## Cookie 的定义

标准库 `net/http` 中定义了 `Cookie`，它代表一个出现在 HTTP 响应头中 `Set-Cookie` 的值里或者 HTTP 请求头中 Cookie 的值的 `HTTP Cookie`。

```go
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

### 设置 Cookie

`net/http` 中提供了如下 `SetCookie` 函数，它在响应的头域中添加 `Set-Cookie` 参数。

```go
func SetCookie(w ResponseWriter, cookie *Cookie)
```

### 获取 Cookie

`Request` 对象拥有两个获取 Cookie 的方法和一个添加 Cookie 的方法：

- 获取 Cookie 的两种方法：

```go
// 解析并返回该请求的 Cookie 头设置的所有 cookie
func (r *Request) Cookies() []*Cookie

// 返回请求中名为 name 的 cookie
// 如果未找到该 cookie 会返回 nil, ErrNoCookie
func (r *Request) Cookie(name string) (*Cookie, error)
```

- 添加 Cookie 的方法：

```go
// AddCookie 向请求中添加一个 cookie
func (r *Request) AddCookie(c *Cookie)
```

## Gin 框架操作 Cookie

```go
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		cookie, _ := c.Cookie("key")
		c.SetCookie("key", "value", 3600,
			"/", "localhost", false, true)
		fmt.Println(cookie)
	})
	_ = r.Run()
}
```

### Session

Cookie 虽然在一定程度上解决了“保持状态”的需求，但是由于 Cookie 本身最大支持 4096 字节，以及 Cookie 本身保存在客户端，可能被拦截或窃取，因此就需要有一种新的东西，它能支持更多的字节，并且保存在服务器，有较高的安全性。这就是 `Session`。

问题来了，基于 HTTP 协议的无状态特征，服务器根本就不知道访问者是“谁”。那么上述的 Cookie 就起到桥接的作用。

用户登陆成功之后，我们在服务端为每个用户创建一个特定的 session 和一个唯一的标识，它们一一对应。其中：

- Session 是在服务端保存的一个数据结构，用来跟踪用户的状态，这个数据可以保存在集群、数据库、文件中；
- 唯一标识通常称为 `Session ID `会写入用户的 Cookie 中。

这样该用户后续再次访问时，请求会自动携带 Cookie 数据，其中包含了`Session ID`，服务器通过该 `Session ID` 就能找到与之对应的 Session 数据，也就知道来的人是“谁”。

Cookie 弥补了 HTTP 无状态的不足，让服务器知道来的人是“谁”；但是 Cookie 以文本的形式保存在本地，自身安全性较差；所以我们就通过 Cookie 识别不同的用户，对应的在服务端为每个用户保存一个 Session 数据，该 Session 数据中能够保存具体的用户数据信息。

另外，上述所说的Cookie和Session其实是共通性的东西，不限于语言和框架。
