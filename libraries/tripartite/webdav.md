---
date: 2022-01-17T10:37:27+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "webdav - 简单的 WebDAV 服务实现"  # 文章标题
url:  "posts/go/libraries/tripartite/webdav"  # 设置网页链接，默认使用文件名
tags: [ "go", "webdav" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

WebDAV（Web-based Distributed Authoring and Versioning）是一种基于 HTTP 1.1 协议的通信协议。它扩展了 HTTP 1.1，在 GET、POST、HEAD 等几个 HTTP 标准方法以外添加了一些新的方法，使应用程序可对 Web Server 直接读写，并支持写文件锁定 (Locking) 及解锁 (Unlock)，还可以支持文件的版本控制。

使用 WebDAV 可以完成的工作包括：

- 特性（元数据）处理。可以使用 WebDAV 中的 PROPFIND 和 PROPPATCH 方法可创建、删除和查询有关文件的信息，例如作者和创建日期。
- 集合和资源的管理。可以使用 GET、PUT、DELETE 和 MKCOL 方法创建文档集合并检索分层结构成员列表（类似于文件系统中的目录）。
- 锁定。可以禁止多人同时对一个文档进行操作。这将有助于防止出现“丢失更新”（更改被覆盖）的问题。
- 名称空间操作。您可以使用 COPY 和 MOVE 方法让服务器复制和删除相关资源。

本节我们尝试用 Go 语言实现自己的 WebDAV 服务。

- [WebDAV 对 HTTP 的扩展](#webdav-对-http-的扩展)
- [最简的 WebDAV 服务](#最简的-webdav-服务)
- [只读的 WebDAV 服务](#只读的-webdav-服务)
- [密码认证 WebDAV 服务](#密码认证-webdav-服务)
- [浏览器视图](#浏览器视图)
- [集成进 Gin](#集成进-gin)

## WebDAV 对 HTTP 的扩展

WebDAV 扩展了 HTTP/1.1 协议。它定义了新的 HTTP 标头，客户机可以通过这些新标头传递 WebDAV 特有的资源请求。这些标头为：

- Destination:
- Lock-Token:
- Timeout:
- DAV:
- If:
- Depth:
- Overwrite:

同时，WebDAV 标准还引入了若干新 HTTP 方法，用于告知启用了 WebDAV 的服务器如何处理请求。这些方法是对现有方法（例如 GET、PUT 和 DELETE）的补充，可用来执行 WebDAV 事务。下面是这些新 HTTP 方法的介绍：

- LOCK。锁定资源，使用 Lock-Token : 标头。
- UNLOCK。解除锁定，使用 Lock-Token : 标头。
- PROPPATCH。设置、更改或删除单个资源的特性。
- PROPFIND。用于获取一个或多个资源的一个或多个特性信息。该请求可能会包含一个值为 0、1 或 infinity 的 Depth : 标头。其中，0 表示指定将获取指定 URI 处的集合的特性（也就是该文件或目录）； 1 表示指定将获取该集合以及位于该指定 URI 之下与其紧邻的资源的特性（非嵌套的子目录或子文件）； infinity 表示指定将获取全部子目录或子文件（深度过大会加重对服务器的负担）。
- COPY。复制资源，可以使用 Depth : 标头移动资源，使用 Destination : 标头指定目标。如果需要，COPY 方法也使用 Overwrite : 标头。
- MOVE。移动资源，可以使用 Depth : 标头移动资源，使用 Destination : 标头指定目标。如果需要，MOVE 方法也使用 Overwrite : 标头。
- MKCOL。用于创建新集合（对应目录）。

## 最简的 WebDAV 服务

Go 语言扩展包 `golang.org/x/net/webdav` 提供了 WebDAV 服务的支持。其中 webdav.Handler 实现了 http.Handle 接口，用处理 WebDAV 特有的 http 请求。要构造 webdav.Handler 对象的话，我们至少需要指定一个文件系统和锁服务。其中 webdav.Dir 将本地的文件系统映射为 WebDAV 的文件系统，webdav.NewMemLS 则是基于本机内存构造一个锁服务。

下面是最简单的 WebDAV 服务实现：

```go
package main

import (
	"golang.org/x/net/webdav"
	"net/http"
)

func main() {
	_ = http.ListenAndServe(
		":8080",
		&webdav.Handler{
			FileSystem: webdav.Dir("."),
			LockSystem: webdav.NewMemLS(),
		},
	)
}
```

运行之后，当前目录就可以通过 WebDAV 方式访问了。

## 只读的 WebDAV 服务

前面实现的 WebDAV 服务默认不需要任何密码就可以访问文件系统，任何匿名的用户可以添加、修改、删除文件，这对于网络服务来说太不安全了。

为了防止被用户无意或恶意修改，我们可以关闭 WebDAV 的修改功能。参考 WebDAV 协议规范可知，修改相关的操作主要涉及 PUT/DELETE/PROPPATCH/MKCOL/COPY/MOVE 等几个方法。我们只要将这几个方法屏蔽了就可以实现一个只读的 WebDAV 服务。

```go
package main

import (
	"golang.org/x/net/webdav"
	"net/http"
)

func main() {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	}
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
				http.Error(w, "WebDAV: Read Only!!!", http.StatusForbidden)
				return
			}
			fs.ServeHTTP(w, req)
		},
	)
	_ = http.ListenAndServe(":8080", nil)
}
```

我们通过 http.HandleFunc 重新包装了 fs.ServeHTTP 方法，然后将和更新相关的操作屏蔽掉。这样我们就实现了一个只读的 WebDAV 服务。

## 密码认证 WebDAV 服务

WebDAV 是基于 HTTP 协议扩展的标准，我们可以通过 HTTP 的基本认证机制设置用户名和密码。

```go
package main

import (
	"golang.org/x/net/webdav"
	"net/http"
)

func main() {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// 获取用户名/密码
		username, password, ok := req.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 验证用户名/密码
		if username != "user" || password != "123456" {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		fs.ServeHTTP(w, req)
	})

	http.ListenAndServe(":8080", nil)
}
```

我们通过 req.BasicAuth 来获取用户名和密码，然后进行验证。如果没有设置用户名和密码，则返回一个 http.StatusUnauthorized 状态，HTTP 客户端会弹出让用户输入密码的窗口。

由于 HTTP 协议并没有加密，因此用户名和密码也是明文传输。为了更安全，我们可以选择用 HTTPS 协议提供 WebDAV 服务。为此，我们需要准备一个证书文件（crypto/tls 包中的 `generate_cert.go` 程序可以生成证书），然后用 http.ListenAndServeTLS 来启动 https 服务。

同时需要注意的是，从 Windows Vista 起，微软就禁用了 http 形式的基本 WebDAV 验证形式 (KB841215)，默认必须使用 https 连接。可以在 Windows Vista/7/8 中，改注册表:

```
HKEY_LOCAL_MACHINE>>SYSTEM>>CurrentControlSet>>Services>>WebClient>>Parameters>>BasicAuthLevel
```

把这个值从 1 改为 2，然后进控制面板 / 服务，把 WebClient 服务重启。

## 浏览器视图

WebDAV 是基于 HTTP 协议，理论上从浏览器访问 WebDAV 服务器会更简单。但是，当我们在浏览器中访问 WebDAV 服务的根目录之后，收到了“Method Not Allowed”错误信息。

这是因为，根据 WebDAV 协议规范，http 的 GET 方法只能用于获取文件。在 Go 语言实现的 webdav 库中，如果用 GET 访问一个目录，会返回一个 http.StatusMethodNotAllowed 状态码，对应“Method Not Allowed”错误信息。

为了支持浏览器删除目录列表，我们对针对目录的 GET 操作单独生成 html 页面，其中，handleDirList 函数用于处理目录列表，然后返回 ture。实现如下：

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
)

func handleDirList(fs webdav.FileSystem, w http.ResponseWriter, req *http.Request) bool {
	ctx := context.Background()
	f, err := fs.OpenFile(ctx, req.URL.Path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(w, "Error reading directory", http.StatusInternalServerError)
		return false
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", name, name)
	}
	fmt.Fprintf(w, "</pre>\n")
	return true
}

func main() {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" && handleDirList(fs.FileSystem, w, req) {
			return
		}
		fs.ServeHTTP(w, req)
	})

	http.ListenAndServe(":8080", nil)
}
```

现在可以通过浏览器来访问 WebDAV 目录列表了。

## 集成进 Gin

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

const (
	MethodHead      = "HEAD"
	MethodGet       = "GET"
	MethodPut       = "PUT"
	MethodPost      = "POST"
	MethodPatch     = "PATCH"
	MethodDelete    = "DELETE"
	MethodOptions   = "OPTIONS"
	MethodMkcol     = "MKCOL"
	MethodCopy      = "COPY"
	MethodMove      = "MOVE"
	MethodLock      = "LOCK"
	MethodUnlock    = "UNLOCK"
	MethodPropfind  = "PROPFIND"
	MethodProppatch = "PROPPATCH"
)

// WebDAV handles any requests to /originals|import/*
func WebDAV(path string, router *gin.RouterGroup) {
	if router == nil {
		log.Println("webdav: router is nil")
		return
	}

	f := webdav.Dir(path)

	srv := &webdav.Handler{
		Prefix:     router.BasePath(),
		FileSystem: f,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				switch r.Method {
				case MethodPut, MethodPost, MethodPatch, MethodDelete, MethodCopy, MethodMove:
					log.Printf("webdav: %s in %s %s", (err.Error()), (r.Method), (r.URL.String()))
				case MethodPropfind:
					log.Printf("webdav: %s in %s %s", (err.Error()), (r.Method), (r.URL.String()))
				default:
					log.Printf("webdav: %s in %s %s", (err.Error()), (r.Method), (r.URL.String()))
				}

			} else {
				log.Printf("webdav: %s %s", (r.Method), (r.URL.String()))
			}
		},
	}

	handler := func(c *gin.Context) {
		w := c.Writer
		r := c.Request

		srv.ServeHTTP(w, r)
	}

	router.Handle(MethodHead, "/*path", handler)
	router.Handle(MethodGet, "/*path", handler)
	router.Handle(MethodPut, "/*path", handler)
	router.Handle(MethodPost, "/*path", handler)
	router.Handle(MethodPatch, "/*path", handler)
	router.Handle(MethodDelete, "/*path", handler)
	router.Handle(MethodOptions, "/*path", handler)
	router.Handle(MethodMkcol, "/*path", handler)
	router.Handle(MethodCopy, "/*path", handler)
	router.Handle(MethodMove, "/*path", handler)
	router.Handle(MethodLock, "/*path", handler)
	router.Handle(MethodUnlock, "/*path", handler)
	router.Handle(MethodPropfind, "/*path", handler)
	router.Handle(MethodProppatch, "/*path", handler)
}

func main() {
	e := gin.Default()
	router := e.Group("/file")

	WebDAV(".", router)

	log.Fatal(e.Run(":8080"))
}
```
