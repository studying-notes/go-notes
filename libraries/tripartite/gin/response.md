---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "多数据格式返回响应数据"  # 文章标题
url:  "posts/gin/doc/response"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## []byte

使用 `context.Writer.Write` 向客户端写入返回数据。

```go
engine := gin.Default()
engine.GET("/hello", func(context *gin.Context) {
    hello := "hello world!"
    context.Writer.Write([]byte(hello))
})
engine.Run()
```

`Writer` 是 Gin 框架中封装的一个 `ResponseWriter` 接口类型，`ResponseWriter` 源码定义如下所示：

```go
type ResponseWriter interface {
    http.ResponseWriter
    http.Hijacker
    http.Flusher
    http.CloseNotifier

    // Returns the HTTP response status code of the current request.
    Status() int

    // Returns the number of bytes already written into the response http body.
    // See Written()
    Size() int

    // Writes the string into the response body.
    WriteString(string) (int, error)

    // Returns true if the response body was already written.
    Written() bool

    // Forces to write the http header (status code + headers).
    WriteHeaderNow()

    // get the http.Pusher for server push
    Pusher() http.Pusher
}
```

## string

`ResponseWriter` 封装了 `WriteString` 方法返回数据。

```go
engine := gin.Default()
engine.GET("/hello", func(context *gin.Context) {
    hello := "hello world!"
    context.Writer.WriteString(hello)
})
engine.Run()
```

## JSON

Gin 框架中的 `context` 包含的 JSON 方法可以将结构体类型的数据转换成 JSON 格式的结构化数据，然后返回给客户端。

Gin 支持 Map 类型和结构体类型。

### Map 类型

```go
engine := gin.Default()
engine.GET("/hello", func(context *gin.Context) {
    context.JSON(200, map[string]interface{}{
    // context.JSON(200, gin.H{
        "code":    1,
        "message": "OK",
        "data":    "hello",
    })
})
engine.Run()
```

### 结构体类型

```go
//通用请求返回结构体定义
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"msg"`
    Data    interface{} `json:"data"`
}

engine := gin.Default()
engine.GET("/hello", func(context *gin.Context) {
    resp := Response{Code: 1, Message: "Ok", Data: "hello"}
    context.JSON(200, &resp)
})
engine.Run()
```

## HTML 模板

```go
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":8080")
}
```

`templates/index.tmpl`

```html
<html>
	<h1>
		{{ .title }}
	</h1>
</html>
```

## 加载静态资源文件

```go
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.Run(":8080")
}
```
