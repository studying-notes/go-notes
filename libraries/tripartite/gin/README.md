---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 学习笔记"  # 文章标题
url:  "posts/go/libraries/tripartite/gin/README"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

以下只是简单地翻译了一遍官方的文档，算是初步熟悉一下 Gin 框架。（2020.07.12）

## 目录

- [目录](#目录)
- [模块代理](#模块代理)
- [安装配置](#安装配置)
- [快速开始](#快速开始)
- [Gin 特性](#gin-特性)
- [通过 jsoniter 编译 JSON](#通过-jsoniter-编译-json)
- [测试用例](#测试用例)
- [官方示例](#官方示例)
	- [HTTP 基本方法](#http-基本方法)
	- [URL 路径参数](#url-路径参数)
	- [解析字符串查询参数](#解析字符串查询参数)
	- [解析 URL 编码的 POST 表单](#解析-url-编码的-post-表单)
	- [同时带有字符串查询参数和表单数据](#同时带有字符串查询参数和表单数据)
	- [POST 数据是映射格式](#post-数据是映射格式)
	- [上传文件](#上传文件)
		- [单个文件](#单个文件)
		- [多个文件](#多个文件)
	- [路由组](#路由组)
	- [中间件](#中间件)
		- [默认中间件](#默认中间件)
		- [自定义中间件](#自定义中间件)
		- [将中间件一分为二执行](#将中间件一分为二执行)
		- [官方示例](#官方示例-1)
	- [自定义从错误中恢复的行为](#自定义从错误中恢复的行为)
	- [写入日志文件](#写入日志文件)
	- [自定义日志格式](#自定义日志格式)
	- [控制日志输出颜色](#控制日志输出颜色)
	- [类型绑定和验证](#类型绑定和验证)
		- [自定义类型验证器](#自定义类型验证器)
		- [只绑定查询字符串](#只绑定查询字符串)
		- [绑定查询字符串或 POST 数据](#绑定查询字符串或-post-数据)
		- [绑定 URI](#绑定-uri)
		- [绑定 Header](#绑定-header)
		- [绑定 HTML 复选框](#绑定-html-复选框)
		- [绑定 URL 编码表单](#绑定-url-编码表单)
	- [尝试将 Body 绑定到不同的结构体](#尝试将-body-绑定到不同的结构体)
	- [渲染 XML、JSON、YAML](#渲染-xmljsonyaml)
		- [SecureJSON](#securejson)
		- [JSONP](#jsonp)
		- [AsciiJSON](#asciijson)
		- [PureJSON](#purejson)
	- [静态文件服务](#静态文件服务)
	- [从文件中呈现内容](#从文件中呈现内容)
	- [从 Reader 对象呈现内容](#从-reader-对象呈现内容)
	- [HTML 模板渲染](#html-模板渲染)
		- [自定义模板渲染器](#自定义模板渲染器)
		- [自定义模板函数](#自定义模板函数)
	- [重定向](#重定向)
	- [自定义中间件](#自定义中间件-1)
	- [BasicAuth() 中间件](#basicauth-中间件)
	- [中间件内部的 Goroutines](#中间件内部的-goroutines)
	- [自定义 HTTP 配置](#自定义-http-配置)
	- [支持 Let's Encrypt](#支持-lets-encrypt)
	- [同时在不同端口提供不同服务](#同时在不同端口提供不同服务)
	- [优雅地关机或重启](#优雅地关机或重启)
		- [第三方包](#第三方包)
		- [手动](#手动)
	- [构建单个包括模板的二进制文件](#构建单个包括模板的二进制文件)
	- [用自定义的结构体绑定 form-data 请求](#用自定义的结构体绑定-form-data-请求)
	- [HTTP2 服务器推送技术](#http2-服务器推送技术)
	- [定义路由日志的格式](#定义路由日志的格式)
	- [设置和获取 Cookie](#设置和获取-cookie)

## 模块代理

国内网络基本上无法访问官方网址，所以设置代理是最好的办法。

```shell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off
go env -w GO111MODULE=on
```

## 安装配置

1. 下载安装 Gin

```shell
go get -u github.com/gin-gonic/gin
```

2. 初始化工作区

```shell
go mod init your_project
go mod edit -require github.com/gin-gonic/gin@latest
go mod vendor
```

3. 设置 VS Code

我个人写 Go 其实习惯于 Goland，但 VS Code 仍是我的最爱，只是看代码的话还是 VS Code 爽。首先必须安装 Microsoft 官方的 Go 插件，其次快捷键 `Ctrl + Shift + P` 进入命令模式，键入 `go:install/update tools` ，将全部插件选中，然后点击确定开始安装。

以下是 VS Code 与 Go 插件相关的全局设置，我一般就是这样设置的：

```json
{
  "go.useLanguageServer": true,
  "go.autocompleteUnimportedPackages": true,
  "go.gocodePackageLookupMode": "go",
  "go.gotoSymbol.includeImports": true,
  "go.useCodeSnippetsOnFunctionSuggest": true,
  "go.useCodeSnippetsOnFunctionSuggestWithoutType": true,
  "go.inferGopath": true,
  "go.docsTool": "gogetdoc",
  "go.formatTool": "goimports",
  "workbench.startupEditor": "newUntitledFile",
  "go.languageServerExperimentalFeatures": {
    "format": true,
    "autoComplete": true,
    "rename": true,
    "goToDefinition": true,
    "hover": true,
    "signatureHelp": true,
    "goToTypeDefinition": true,
    "goToImplementation": true,
    "documentSymbols": true,
    "workspaceSymbols": true,
    "findReferences": true,
    "diagnostics": true,
    "documentLink": true
  },
  // 只检查新加入代码
  "go.lintFlags": ["--enable-all", "--new"],
}
```

4. 通过以下方式导入模块

```go
import "github.com/gin-gonic/gin"
```

## 快速开始

1. 创建 `example.go` 文件

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
```

2. 命令行运行

```shell
# localhost:8080/ping
go run example.go
```

3. 访问 `localhost:8080/ping`

```json
{
    "message": "pong"
}
```

## Gin 特性

- 零分配路由
- 从路由到数据写入，仍是最快的 HTTP 框架
- 完备的单元测试套件
- 可经受生产环境挑战
- 稳定的 API，新的版本发布不影响之前的代码

## 通过 jsoniter 编译 JSON

默认标准库处理 JSON 数据，可以通过配置修改为更快的 [`jsoniter`](https://github.com/json-iterator/go)

```shell
go build -tags=jsoniter .
```

## 测试用例

HTTP 测试推荐 `net/http/httptest` 标准库，以 `example.go` 为例：

```go
package main

import "github.com/gin-gonic/gin"

func pingRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}
```

可编写如下测试用例：

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := pingRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
```

## 官方示例

### HTTP 基本方法

```go
func main() {
	router := gin.Default()

	router.GET("get", get)
	router.GET("post", post)
	router.GET("put", put)
	router.GET("delete", delete)
	router.GET("patch", patch)
	router.GET("head", head)
	router.GET("options", options)

	// 默认 8080 端口
	router.Run()
	// router.Run(":80")
}
```

### URL 路径参数

```go
func main() {
	router := gin.Default()

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})

	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + "is" + action
		c.String(http.StatusOK, message)
	})

	router.POST("user/:name/*action", func(c *gin.Context) {
		_ = c.FullPath() == "/user/:name/*action"
	})
	router.Run()
}
```

```shell
$ curl localhost:8080/user/rustlekarl
hello rustlekarl
```

### 解析字符串查询参数

```go
func main() {
	router := gin.Default()
	// 示例解析：/welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
    // 等价于 c.Request.URL.Query().Get("lastname")
    lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request GET 'localhost:8080/welcome?firstname=Jane&lastname=Doe'
Hello Jane Doe
```

### 解析 URL 编码的 POST 表单

```go
func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/form_post' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'message=hello' --data-urlencode 'nick=rustlekarl'
{"message":"hello","nick":"rustlekarl","status":"posted"}
```

### 同时带有字符串查询参数和表单数据

```go
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		c.String(http.StatusOK, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/post?id=1&page=2' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'message=hello' --data-urlencode 'name=rustlekarl'
id: 1; page: 2; name: rustlekarl; message: hello
```

### POST 数据是映射格式

```go
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		// c.String(http.StatusOK, "ids: %v; names: %v", ids, names)
		c.String(http.StatusOK, "ids[a]: %v; names[first]: %v", ids["a"], names["first"])
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/post?ids[a]=1234&ids[b]=hello' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'names[first]=rustle' --data-urlencode 'names[second]=karl'
ids[a]: 1234; names[first]: rustle
```

### 上传文件

#### 单个文件

```go
func main() {
	router := gin.Default()
	// 设置上传文件的最大内存限制
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "upload/file"
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'http://localhost:8080/upload' 
--form 'file=@D:/LICENSE' --header "Content-Type: multipart/form-data"        
'LICENSE' uploaded!
```

#### 多个文件

```go
func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			c.SaveUploadedFile(file, file.Filename)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'http://localhost:8080/upload' --form 'upload[]=@D:/LICENSE' --form 'upload[]=@D:/README.rst'
2 files uploaded!
```

### 路由组

RouteGroup 可以称之为路由集合，在 Gin 中定义为结构体：

```go
type RouterGroup struct {
    Handlers HandlersChain
    basePath string
    engine   *Engine
    root     bool
}
```

在实际的项目开发中，均是模块化开发。同一模块内的功能接口，往往会有相同的接口前缀。类似这种接口前缀统一，均属于相同模块的功能接口。可以使用路由组进行分类处理。

```go
func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}
```

RouterGroup 实现了 IRoutes 中定义的方法：

```go
type IRoutes interface {
    Use(...HandlerFunc) IRoutes

    Handle(string, string, ...HandlerFunc) IRoutes
    Any(string, ...HandlerFunc) IRoutes
    GET(string, ...HandlerFunc) IRoutes
    POST(string, ...HandlerFunc) IRoutes
    DELETE(string, ...HandlerFunc) IRoutes
    PATCH(string, ...HandlerFunc) IRoutes
    PUT(string, ...HandlerFunc) IRoutes
    OPTIONS(string, ...HandlerFunc) IRoutes
    HEAD(string, ...HandlerFunc) IRoutes

    StaticFile(string, string) IRoutes
    Static(string, string) IRoutes
    StaticFS(string, http.FileSystem) IRoutes
}
```

### 中间件

在 Web 应用服务中，完整的一个业务处理在技术上包含客户端操作、服务器端处理、返回处理结果给客户端三个步骤。

在实际的业务开发和处理中，会有更负责的业务和需求场景。一个完整的系统可能要包含鉴权认证、权限管理、安全检查、日志记录等多维度的系统支持。

鉴权认证、权限管理、安全检查、日志记录等这些保障和支持系统业务属于全系统的业务，和具体的系统业务没有关联，对于系统中的所有业务都适用。

由此，在业务开发过程中，为了更好的梳理系统架构，可以将上述描述所涉及的一些通用业务单独抽离并进行开发，然后以插件化的形式进行对接。这种方式既保证了系统功能的完整，同时又有效的将具体业务和系统功能进行解耦，并且，还可以达到灵活配置的目的。

这种通用业务独立开发并灵活配置使用的组件，一般称之为"中间件"，因为其位于服务器和实际业务处理程序之间。其含义就是相当于在请求和具体的业务逻辑处理之间增加某些操作，这种以额外添加的方式不会影响编码效率，也不会侵入到框架中。中间件的位置和角色示意图如下图所示：

![](imgs/1.png)

Gin 中间件的类型定义如下所示：

```go
// HandlerFunc defines the handler used
// by gin middleware as return value.
type HandlerFunc func(*Context)
```

#### 默认中间件

```go
func Default() *Engine {
    debugPrintWARNINGDefault()
    engine := New()
    engine.Use(Logger(), Recovery())
    return engine
}
// Log 中间件
func Logger() HandlerFunc {
    return LoggerWithConfig(LoggerConfig{})
}
// Recovery 中间件
func Recovery() HandlerFunc {
    return RecoveryWithWriter(DefaultErrorWriter)
}

// Use 方法定义
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
    engine.RouterGroup.Use(middleware...)
    engine.rebuild404Handlers()
    engine.rebuild405Handlers()
    return engine
}
```

#### 自定义中间件

处理请求时，为了方便调试，通常都将请求的一些信息打印出来。有了中间件以后，为了避免代码多次重复编写，使用统一的中间件来完成。定义一个名为 RequestInfos 的中间件，在该中间件中打印请求的 URI 和类型。

```go
func RequestInfos() gin.HandlerFunc {
    return func(context *gin.Context) {
        path := context.FullPath()
        method := context.Request.Method
        fmt.Println("URL: ", path)
        fmt.Println("Method: ", method)
    }
}

func main() {

    engine := gin.Default()
    engine.Use(RequestInfos())

    engine.GET("/query", func(context *gin.Context) {
        context.JSON(200, map[string]interface{}{
            "code": 1,
            "msg":  context.FullPath(),
        })
    })
    engine.Run(":9000")
}
```

#### 将中间件一分为二执行

在上文自定义的中间件 RequestInfos 中，打印了请求信息，然后执行了正常的业务处理函数。若想输出业务处理结果的信息，可以用 `context.Next` 函数。

`context.Next` 函数可以将中间件代码的执行顺序一分为二，`Next` 函数调用之前的代码在请求处理之前之前，当程序执行到 `context.Next` 时，会中断向下执行，转而先去执行具体的业务逻辑，执行完业务逻辑处理函数之后，程序会再次回到 `context.Next` 处，继续执行中间件后续的代码。

```go
func RequestInfos() gin.HandlerFunc {
    return func(context *gin.Context) {
        path := context.FullPath()
        method := context.Request.Method
        fmt.Println("URL: ", path)
		fmt.Println("Method: ", method)
		context.Next()
		fmt.Println(context.Writer.Status())
    }
}

func main() {
    engine := gin.Default()
    engine.Use(RequestInfos())

    engine.GET("/query", func(context *gin.Context) {
        context.JSON(404, map[string]interface{}{
            "code": 1,
            "msg":  context.FullPath(),
        })
    })
    engine.Run(":9000")
}
```

#### 官方示例

```go
func main() {
  // 创建一个不带任何中间件的路由
	r := gin.New()

	// 全局中间件
	// 日志中间件将日志写入 gin.DefaultWriter 即使设置了 GIN_MODE=release
	// 默认 gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

  // 恢复中间件从任何错误中恢复状态，返回 500 状态码
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 自定义从错误中恢复的行为

```go
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("foo")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ohai")
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 写入日志文件

```go
func main() {
    // Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

### 自定义日志格式

```go
func main() {
	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
```

### 控制日志输出颜色

```go
func main() {
    // Disable log's color
    gin.DisableConsoleColor()
    
    // Creates a gin router with default middleware:
    // logger and recovery (crash-free) middleware
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    router.Run(":8080")
}
```

```go
func main() {
    // Force log's color
    gin.ForceConsoleColor()
    
    // Creates a gin router with default middleware:
    // logger and recovery (crash-free) middleware
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    router.Run(":8080")
}
```

### 类型绑定和验证

一般情况下，如果表单数据少，用 `context.PostForm` 或者 `context.DefaultPostForm` 获取即可，但是如果表单数据很多，一个一个获取写起来很烦。

因此，Gin 可以将请求 Body 绑定到不同的结构体中。目前支持绑定 JSON、XML、YAML 和标准表单形式。

注意到必须在绑定的字段上设置标签。Gin 提供了两者类型的绑定方式：

- Type - Must bind
  - Methods - `Bind`, `BindJSON`, `BindXML`, `BindQuery`, `BindYAML`, `BindHeader`

- Type - Should bind
  - Methods - `ShouldBind`, `ShouldBindJSON`, `ShouldBindXML`, `ShouldBindQuery`, `ShouldBindYAML`, `ShouldBindHeader`

当使用绑定类型时，Gin 尝试通过 `Content-Type` 参数值推断绑定的结构体。如果可以确定数据就是某种类型，可以使用 `MustBindWith` 或 `ShouldBindWith`。

1. 解析错误在 header 中写一个 400 的状态码

```go
// 内部根据 Content-Type 判断解析
c.Bind(obj interface{})

// 内部传递了一个 binding.JSON 对象解析
c.BindJSON(obj interface{})
// 其实就是下面方法的快捷方式
c.BindWith(obj interface{}, b binding.JSON)

// 自行传入哪一种绑定的类型解析
c.BindWith(obj interface{}, b binding.Binding)
```

2. 解析错误直接返回，至于要给客户端返回什么错误状态码由编写者决定

```go
// 内部根据 Content-Type 判断解析
c.ShouldBind(obj interface{})

// 内部传递了一个 binding.JSON 对象解析
c.ShouldBindJSON(obj interface{})
// 其实就是下面方法的快捷方式
c.ShouldBindWith(obj interface{}, b binding.JSON)

// 自行传入哪一种绑定的类型解析
c.ShouldBindWith(obj interface{}, b binding.Binding)
```

```go
// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>user</user>
	//		<password>123</password>
	//	</root>)
	router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
```

#### 自定义类型验证器

```go
package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Booking contains binded and validated data.
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func main() {
	route := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	route.GET("/bookable", getBookable)
	route.Run(":8085")
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
```

```shell
$ curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-17"
{"message":"Booking dates are valid!"}

$ curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-15"
{"error":"Key: 'Booking.CheckOut' Error:Field validation for 'CheckOut' failed on the 'gtfield' tag"} 
```

#### 只绑定查询字符串

`ShouldBindQuery` 函数只绑定 GET 方式的查询字符串参数，不包括 POST 数据。

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	route := gin.Default()
	route.Any("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}
```

#### 绑定查询字符串或 POST 数据

`ShouldBind` 函数可以绑定 POST 方式提交的数据。

```go
package main

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

type Person struct {
        Name       string    `form:"name"`
        Address    string    `form:"address"`
        Birthday   time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
        CreateTime time.Time `form:"createTime" time_format:"unixNano"`
        UnixTime   time.Time `form:"unixTime" time_format:"unix"`
}

func main() {
	route := gin.Default()
	route.GET("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
        if c.ShouldBind(&person) == nil {
                log.Println(person.Name)
                log.Println(person.Address)
                log.Println(person.Birthday)
                log.Println(person.CreateTime)
                log.Println(person.UnixTime)
        }

	c.String(200, "Success")
}
```

#### 绑定 URI

```go
package main

import "github.com/gin-gonic/gin"

type Person struct {
	ID string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	route := gin.Default()
	route.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})
	route.Run(":8088")
}
```

```shell
$ curl localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
{"name":"thinkerou","uuid":"987fbc97-4bed-5078-9f07-9141ba07c9f3"}

$ curl localhost:8088/thinkerou/not-uuid
{"msg":[{}]}
```

#### 绑定 Header

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type testHeader struct {
	Rate   int    `header:"Rate"`
	Domain string `header:"Domain"`
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		h := testHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(200, err)
		}

		fmt.Printf("%#v\n", h)
		c.JSON(200, gin.H{"Rate": h.Rate, "Domain": h.Domain})
	})

	r.Run()
}
```

```shell
$ curl -H "rate:300" -H "domain:music" 127.0.0.1:8080
{"Domain":"music","Rate":300}
```

#### 绑定 HTML 复选框

```go
type myForm struct {
    Colors []string `form:"colors[]"`
}

func formHandler(c *gin.Context) {
    var fakeForm myForm
    c.ShouldBind(&fakeForm)
    c.JSON(200, gin.H{"color": fakeForm.Colors})
}
```

#### 绑定 URL 编码表单

```go
type ProfileForm struct {
	Name   string                `form:"name" binding:"required"`
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`

	// or for multiple files
	// Avatars []*multipart.FileHeader `form:"avatar" binding:"required"`
}

func main() {
	router := gin.Default()
	router.POST("/profile", func(c *gin.Context) {
		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form ProfileForm
		// in this case proper binding will be automatically selected
		if err := c.ShouldBind(&form); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		err := c.SaveUploadedFile(form.Avatar, form.Avatar.Filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "unknown error")
			return
		}

		// db.Save(&form)

		c.String(http.StatusOK, "ok")
	})
	router.Run(":8080")
}
```

### 尝试将 Body 绑定到不同的结构体

一般的绑定请求数据的类型会消耗 `c.Request.Body`，**读取之后就为空了**，所以不能被多次调用。

```go
type formA struct {
  Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
  Bar string `json:"bar" xml:"bar" binding:"required"`
}

func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // This c.ShouldBind consumes c.Request.Body and it cannot be reused.
  if errA := c.ShouldBind(&objA); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // Always an error is occurred by this because c.Request.Body is EOF now.
  } else if errB := c.ShouldBind(&objB); errB == nil {
    c.String(http.StatusOK, `the body should be formB`)
  } else {
    ...
  }
}
```

这种情况，可以用 `c.ShouldBindBodyWith`：

```go
func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // This reads c.Request.Body and stores the result into the context.
  if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // At this time, it reuses body stored in the context.
  } else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
    c.String(http.StatusOK, `the body should be formB JSON`)
  // And it can accepts other formats
  } else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
    c.String(http.StatusOK, `the body should be formB XML`)
  } else {
    ...
  }
}
```

- `c.ShouldBindBodyWith` 在绑定之前将 Body 存储进了上下文中，这对性能有轻微的影响。
- 这个特性只在 JSON、XML、MsgPack、ProtoBuf 等格式时有必要，对于字符串查询、表单等可以用 c.`ShouldBind()` 调用多次而不损失性能。

### 渲染 XML、JSON、YAML

```go

func main() {
	r := gin.Default()

	// gin.H is a shortcut for map[string]interface{}
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### SecureJSON

```go
func main() {
	r := gin.Default()

	// You can also use your own secure json prefix
	// r.SecureJsonPrefix(")]}',\n")

	r.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// Will output  :   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### JSONP

```go
func main() {
	r := gin.Default()

	r.GET("/JSONP", func(c *gin.Context) {
		data := gin.H{
			"foo": "bar",
		}
		
		//callback is x
		// Will output  :   x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")

        // client
        // curl http://127.0.0.1:8080/JSONP?callback=x
}
```

#### AsciiJSON

```go
func main() {
	r := gin.Default()

	r.GET("/someJSON", func(c *gin.Context) {
		data := gin.H{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### PureJSON

```go
func main() {
	r := gin.Default()
	
	// Serves unicode entities
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// Serves literal characters
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 静态文件服务

```go
func (group *RouterGroup) Static(relativePath, root string) IRoutes {
	return group.StaticFS(relativePath, Dir(root, false))
}

func main() {
	router := gin.Default()

	// 文件夹，本地文件系统
	router.Static("/assets", "./assets")

	// 文件夹，可自定义文件系统
	router.StaticFS("/more_static", http.Dir("my_file_system"))

	// 单个文件
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.Run(":8080")
}
```

### 从文件中呈现内容

```go
func main() {
	router := gin.Default()

	router.GET("/local/file", func(c *gin.Context) {
		c.File("local/file.go")
	})

	var fs http.FileSystem = // ...
	router.GET("/fs/file", func(c *gin.Context) {
		c.FileFromFS("fs/file.go", fs)
	})
}
```

### 从 Reader 对象呈现内容

```go
func main() {
	router := gin.Default()
	router.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	router.Run(":8080")
}
```

### HTML 模板渲染

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

不同目录的相同文件名模板

```go
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.Run(":8080")
}
```

`templates/posts/index.tmpl`

```html
{{ define "posts/index.tmpl" }}
<html><h1>
	{{ .title }}
</h1>
<p>Using posts/index.tmpl</p>
</html>
{{ end }}
```

`templates/users/index.tmpl`

```html
{{ define "users/index.tmpl" }}
<html><h1>
	{{ .title }}
</h1>
<p>Using users/index.tmpl</p>
</html>
{{ end }}
```

#### 自定义模板渲染器

```go
import "html/template"

func main() {
	router := gin.Default()
	html := template.Must(template.ParseFiles("file1", "file2"))
	router.SetHTMLTemplate(html)
	router.Run(":8080")
}
```

**自定义分隔符**

```go
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("/path/to/templates")
```

#### 自定义模板函数

```go
import (
    "fmt"
    "html/template"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func main() {
    router := gin.Default()
    router.Delims("{[{", "}]}")
    router.SetFuncMap(template.FuncMap{
        "formatAsDate": formatAsDate,
    })
    router.LoadHTMLFiles("./testdata/template/raw.tmpl")

    router.GET("/raw", func(c *gin.Context) {
        c.HTML(http.StatusOK, "raw.tmpl", gin.H{
            "now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
        })
    })

    router.Run(":8080")
}
```

### 重定向

HTTP 重定向

```go
r.GET("/test", func(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})
```

```go
r.POST("/test", func(c *gin.Context) {
	c.Redirect(http.StatusFound, "/foo")
})
```

路由重定向

```go
r.GET("/test", func(c *gin.Context) {
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(200, gin.H{"hello": "world"})
})
```

### 自定义中间件

```go
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		log.Println(example)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### BasicAuth() 中间件

```go
// simulate some private data
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	r := gin.Default()

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 中间件内部的 Goroutines

当在中间件或处理函数内部启动新的 Goroutines，不应该使用原本的上下文，而必须使用只读副本。

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// since we are NOT using a goroutine, we do not have to copy the context
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 自定义 HTTP 配置

```go
func main() {
	router := gin.Default()
	http.ListenAndServe(":8080", router)
}
```

又或者：

```go
func main() {
	router := gin.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
```

### 支持 Let's Encrypt

LetsEncrypt HTTPS 服务器示例：

```go
package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
```

自定义的 autocert 管理器示例：

```go
package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	log.Fatal(autotls.RunWithManager(r, &m))
}
```

### 同时在不同端口提供不同服务

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})

	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})

	return e
}

func main() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := server01.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := server02.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
```

### 优雅地关机或重启

#### 第三方包

可以用 [fvbock/endless](https://github.com/fvbock/endless) 替代默认的 `ListenAndServe`。

```go
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

#### 手动

Go 1.8 及以上版本，可以直接用 `http.Server` 内置的 `Shutdown()` 类型优雅地关闭服务。

```go
// +build go1.8

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	log.Println("Server exiting")
}
```

### 构建单个包括模板的二进制文件

通过 [`go-assets`](https://github.com/jessevdk/go-assets) ，可以将包括模板的服务编译成单个二进制文件。

```go
func main() {
	r := gin.New()

	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/html/index.tmpl",nil)
	})
	r.Run(":8080")
}

// loadTemplate loads templates embedded by go-assets-builder
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		defer file.Close()
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
```

### 用自定义的结构体绑定 form-data 请求

```go
type StructA struct {
    FieldA string `form:"field_a"`
}

type StructB struct {
    NestedStruct StructA
    FieldB string `form:"field_b"`
}

type StructC struct {
    NestedStructPointer *StructA
    FieldC string `form:"field_c"`
}

type StructD struct {
    NestedAnonyStruct struct {
        FieldX string `form:"field_x"`
    }
    FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
    var b StructB
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStruct,
        "b": b.FieldB,
    })
}

func GetDataC(c *gin.Context) {
    var b StructC
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStructPointer,
        "c": b.FieldC,
    })
}

func GetDataD(c *gin.Context) {
    var b StructD
    c.Bind(&b)
    c.JSON(200, gin.H{
        "x": b.NestedAnonyStruct,
        "d": b.FieldD,
    })
}

func main() {
    r := gin.Default()
    r.GET("/getb", GetDataB)
    r.GET("/getc", GetDataC)
    r.GET("/getd", GetDataD)

    r.Run()
}
```

### HTTP2 服务器推送技术

```go
package main

import (
	"html/template"
	"log"
	"github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// use pusher.Push() to do server push
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})

	// Listen and Server in https://127.0.0.1:8080
	r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
```

### 定义路由日志的格式

```go
import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	r.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// Listen and Server in http://0.0.0.0:8080
	r.Run()
}
```

### 设置和获取 Cookie

```go
import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func main() {

    router := gin.Default()
    router.GET("/cookie", func(c *gin.Context) {
        cookie, err := c.Cookie("gin_cookie")
        if err != nil {
            cookie = "NotSet"
            c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
        }
        fmt.Printf("Cookie value: %s \n", cookie)
    })
    router.Run()
}
```
