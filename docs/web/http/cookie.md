---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Cookie 和 Session"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/web/http/cookie"  # 设置网页永久链接
tags: [ "go", "http", "cookie" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: false  # 是否自动生成目录
draft: false  # 草稿
---

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

# Gin 框架操作 Cookie

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

## Session

Cookie 虽然在一定程度上解决了“保持状态”的需求，但是由于 Cookie 本身最大支持 4096 字节，以及 Cookie 本身保存在客户端，可能被拦截或窃取，因此就需要有一种新的东西，它能支持更多的字节，并且保存在服务器，有较高的安全性。这就是 `Session`。

问题来了，基于 HTTP 协议的无状态特征，服务器根本就不知道访问者是“谁”。那么上述的 Cookie 就起到桥接的作用。

用户登陆成功之后，我们在服务端为每个用户创建一个特定的 session 和一个唯一的标识，它们一一对应。其中：

- Session 是在服务端保存的一个数据结构，用来跟踪用户的状态，这个数据可以保存在集群、数据库、文件中；
- 唯一标识通常称为 `Session ID `会写入用户的 Cookie 中。

这样该用户后续再次访问时，请求会自动携带 Cookie 数据，其中包含了`Session ID`，服务器通过该 `Session ID` 就能找到与之对应的 Session 数据，也就知道来的人是“谁”。

Cookie 弥补了 HTTP 无状态的不足，让服务器知道来的人是“谁”；但是 Cookie 以文本的形式保存在本地，自身安全性较差；所以我们就通过 Cookie 识别不同的用户，对应的在服务端为每个用户保存一个 Session 数据，该 Session 数据中能够保存具体的用户数据信息。

另外，上述所说的Cookie和Session其实是共通性的东西，不限于语言和框架。
