---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "优雅地读取 HTTP 请求数据"  # 文章标题
url:  "posts/gin/doc/readbody"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

从 `http.Request.Body` 或 `http.Response.Body` 中读取数据方法或许很多，标准库中大多数使用 `ioutil.ReadAll` 方法一次读取所有数据， `json` 格式的数据还可以使用 `json.NewDecoder` 从 `io.Reader` 创建一个解析器，假使使用 `pprof` 来分析程序总是会发现 `bytes.makeSlice` 分配了大量内存，且总是排行第一，今天就这个问题来说一下如何高效优雅地读取 `http` 中的数据。

## 标准库读取

```go
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}
```

标准库 `ioutil.ReadAll` 每次会创建一个 `var buf bytes.Buffer` 并且初始化 `buf.Grow(int(capacity))` 的大小为 `bytes.MinRead`, 这个值是 `512`，按这个 `buffer` 的大小读取一次数据需要分配 2~16 次内存。

## 优化读取方法

```go
buffer := bytes.NewBuffer(make([]byte, 4096))
_, _ = io.Copy(buffer, c.Request.Body)
```

大部分请求数据都是比 `4096` 小的。由于每次请求都必须创建 `buffer`，可以用 `sync.Pool` 建立一个缓冲池。

```go
var pool *sync.Pool

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 4096))
		},
	}
}

func main() {
	c := http.Request{}
	
	buffer := pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			pool.Put(buffer)
			buffer = nil
		}
	}()
	_, _ = io.Copy(buffer, c.Body)
}
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




