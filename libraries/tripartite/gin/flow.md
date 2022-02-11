---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 处理流程分析"  # 文章标题
url:  "posts/gin/abc/flow"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 简单示例

```go
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

## 处理流程图

![](../imgs/flow.png)

## gin.Default

```go
func Default() *Engine {
    debugPrintWARNINGDefault()
    engine := New()
    engine.Use(Logger(), Recovery())
    return engine
}
```

通过调用 gin.Default 方法创建默认的 Engine 实例，它会在初始化阶段引入 Logger 和 Recovery 中间件，保障应用程序最基本的运作。这两个中间件的作用如下。

- Logger：输出请求日志，并标准化日志的格式。
- Recovery：异常捕获，也就是针对每次请求处理进行 Recovery 处理，包括恢复现场和写入 `500` 状态码，防止因为出现 `panic` 导致服务崩溃，同时将异常日志的格式标准化。

另外，在调用 debugPrintWARNINGDefault 方法时，首先会检查 Go 版本是否达到 gin 的最低要求，然后再调试日志 `[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.` 的输出，以此提醒开发人员框架内部已经开始检查并集成了默认值。

## gin.New

对 Engine 实例执行初始化动作并返回。

```go
func New() *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		AppEngine:              defaultAppEngine,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 9),
		delims:                 render.Delims{Left: "{{", Right: "}}"},
		secureJsonPrefix:       "while(1);",
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}
```

- RouterGroup：路由组。所有的路由规则都由 ` * RouterGroup` 所属的方法进行管理。在 gin 中，路由组和 `Engine` 实例形成了一个重要的关联组件。
- RedirectTrailingSlash：是否自动重定向。如果启用，在无法匹配当前路由的情况下，则自动重定向到带有或不带斜杠的处理程序中。例如，当外部请求了 ` / tour / ` 路由，但当前并没有注册该路由规则，而只有 ` / tour` 的路由规则时，将会在内部进行判定。若是 HTTP GET 请求，则会通过 HTTP Code 301 重定向到 ` / tour` 的处理程序中；若是其他类型的 HTTP 请求，则会以 HTTP Code 307 重定向，通过指定的 HTTP 状态码重定向到 ` / tour` 路由的处理程序中。
- RedirectFixedPath：是否尝试修复当前请求路径，也就是在开启的情况下，gin 会尽可能地找到一个相似的路由规则，并在内部重定向。RedirectFixedPath 的主要功能是对当前的请求路径进行格式清除（删除多余的斜杠）和不区分大小写的路由查找等。
- HandleMethodNotAllowed：判断当前路由是否允许调用其他方法，如果当前请求无法路由，则返回 Method Not Allowed （ HTTP Code 405 ）的响应结果。如果既无法路由，也不支持重定向到其他方法，则交由 NotFound Hander 进行处理。
- ForwardedByClientIP：如果开启，则尽可能地返回真实的客户端 IP 地址，先从 X-Forwarded-For 中取值，如果没有，则再从 X-Real-Ip 中取值。
- UseRawPath：如果开启，则使用 url.RawPath 来获取请求参数；如果不开启，则还是按 url.Path 来获取请求参数。
- UnescapePathValues：是否对路径值进行转义处理。
- MaxMultipartMemory：对应 http.Request ParseMultipartForm 方法，用于控制最大的文件上传大小。
- trees：多个压缩字典树（ Radix Tree ），每个树都对应一种 HTTP Method。可以这样理解，每当添加一个新路由规则时，就会往 HTTP Method 对应的树里新增一个 node 节点，以此形成关联关系。
- delims：用于 HTML 模板的左右定界符。

总体来讲，Engine实例就像引擎一样，与整个应用的运行、路由、对象、模板等管理和调度都有关联。另外，通过上述解析可以发现，其实 gin 在初始化时默认已经做了很多事情，可以说是既定了一些默认运行基础。

## r.GET

```go
// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}
```

- 计算路由的绝对路径，即 group.basePath 与我们定义的路由路径组装。group 是什么呢？实际上，在 gin 中存在组别路由的概念。
- 合并现有的和新注册的 Handler，并创建一个函数链 HandlersChain。
- 将当前注册的路由规则（含 HTTP Method、Path 和 Handlers）追加到对应的树中。

这类方法主要针对路由的各类计算和注册行为，并输出路由注册的调试信息，如运行时的路由信息

```go
func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}
```

在 combineHandlers 方法中，最终函数链 HandlersChain 是由 group.Handlers 和外部传入的 handlers 组成的，从拷贝的顺序来看，group.Handlers 的优先级高于外部传入的 handlers。
以此再结合 Use 方法来看，很显然是在 gin.Default 方法中注册的中间件影响了这个结果。因为中间件也属于 group.Handlers 的一部分，也就是在调用 gin.Use 时，就已经注册进去了

## r.Run

```go
func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)
	return
}
```

该方法会通过解析地址，再调用 http.ListenAndServe，将 Engine 实例作为 Handler 注册进去，然后启动服务，开始对外提供 HTTP 服务。

```go
// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}
```

在 gin 中，Engine 结构体实现了 ServeHTTP 方法，即符合 http.Handler 接口标准

- 从 sync.Pool 对象池中获取一个上下文对象。
- 重新初始化取出来的上下文对象。
- 处理外部的 HTTP 请求。
- 处理完毕，将取出的上下文对象返回给对象池。

在这里，上下文的池化主要是为了防止频繁反复生成上下文对象，相对地提高性能，并且针对 gin 本身的处理逻辑进行二次封装处理。
