---
date: 2022-02-23T11:01:37+08:00
author: "Rustle Karl"

title: "Fiber，超高性能的Go Web框架"
url:  "posts/go/libraries/tripartite/fiber"  # 永久链接
tags: [ "go", "README" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

toc: true  # 目录
draft: false  # 草稿
---

- [安装](#安装)
- [快速入门](#快速入门)
- [特点](#特点)
- [哲学](#哲学)
- [示例](#示例)
    - [**基础路由**](#基础路由)
    - [**静态文件**服务](#静态文件服务)
    - [**中间件**和[**Next**](https://docs.gofiber.io/api/ctx#next)](#中间件和next)
  - [模版引擎](#模版引擎)
  - [组合路由链](#组合路由链)
  - [日志中间件](#日志中间件)
  - [跨域资源共享(CORS)中间件](#跨域资源共享cors中间件)
  - [自定义 404 响应](#自定义-404-响应)
  - [JSON 响应](#json-响应)
  - [升级到 WebSocket](#升级到-websocket)
  - [恢复(panic)中间件](#恢复panic中间件)
- [内置中间件](#内置中间件)
- [外部中间件](#外部中间件)
- [第三方中间件](#第三方中间件)

**Fiber**，一个受[Express](https://github.com/expressjs/express)启发的Golang **Web框架**，建立在[Fasthttp](https://github.com/valyala/fasthttp) 的基础之上。旨在**简化**、**零内存分配**和**高性能**，以及**快速**开发。

## 安装

确保已安装 ([下载](https://golang.org/dl/)) `1.14` 或更高版本的 Go。

通过创建文件夹并在文件夹内运行 `go mod init github.com/your/repo` ([了解更多](https://blog.golang.org/using-go-modules)) 来初始化项目，然后使用 [`go get`](https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them) 命令安装 Fiber：

```
go get -u github.com/gofiber/fiber/v2
```

## 快速入门

```go
package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })

    app.Listen(":3000")
}
```

## 特点

- 强大的[路由](https://docs.gofiber.io/routing)
- [静态文件](https://docs.gofiber.io/api/app#static)服务
- 极限[表现](https://docs.gofiber.io/benchmarks)
- [内存占用低](https://docs.gofiber.io/benchmarks)
- [API 接口](https://docs.gofiber.io/api/ctx)
- [中间件](https://docs.gofiber.io/middleware)和[Next](https://docs.gofiber.io/api/ctx#next)支持
- [快速](https://dev.to/koddr/welcome-to-fiber-an-express-js-styled-fastest-web-framework-written-with-on-golang-497)服务器端编程
- [模版引擎](https://github.com/gofiber/template)
- [WebSocket 支持](https://github.com/gofiber/websocket)
- [频率限制器](https://docs.gofiber.io/middleware#limiter)
- [15 种语言](https://docs.gofiber.io/)
- 以及更多请[探索文档](https://docs.gofiber.io/)

## 哲学

从[Node.js](https://nodejs.org/en/about/)切换到[Go](https://golang.org/doc/)的新`gopher`在开始构建`Web`应用程序或微服务之前正在应对学习曲线。 `Fiber`作为一个**Web 框架** ，是按照**极简主义**的思想并遵循**UNIX 方式**创建的，因此新的`gopher`可以在热烈和可信赖的欢迎中迅速进入`Go`的世界。

`Fiber`受到了互联网上最流行的`Web`框架`Express`的**启发** 。我们结合了`Express`的**易用性**和`Go`的**原始性能** 。如果您曾经在`Node.js`上实现过`Web`应用程序(*使用 Express 或类似工具*)，那么许多方法和原理对您来说应该**非常易懂**。

我们**关注** *整个互联网* 用户在[issues](https://github.com/gofiber/fiber/issues)和 Discord [channel](https://gofiber.io/discord)的消息，为了创建一个**迅速**，**灵活**以及**友好**的`Go web`框架，满足**任何**任务，**最后期限**和开发者**技能**。就像`Express`在`JavaScript`世界中一样。

## 示例

下面列出了一些常见示例。如果您想查看更多代码示例，请访问我们的[Recipes](https://github.com/gofiber/recipes)代码库或[API 文档](https://docs.gofiber.io/) 。

#### [**基础路由**](https://docs.gofiber.io/#basic-routing)

```go
func main() {
    app := fiber.New()
    // GET /john
    app.Get("/:name", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
        return c.SendString(msg) // => Hello john 👋!
    })
    // GET /john/75
    app.Get("/:name/:age", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
        return c.SendString(msg) // => 👴 john is 75 years old
    })
    // GET /dictionary.txt
    app.Get("/:file.:ext", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
        return c.SendString(msg) // => 📃 dictionary.txt
    })
    // GET /flights/LAX-SFO
    app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
        return c.SendString(msg) // => 💸 From: LAX, To: SFO
    })
    // GET /api/register
    app.Get("/api/*", func(c *fiber.Ctx) error {
        msg := fmt.Sprintf("✋ %s", c.Params("*"))
        return c.SendString(msg) // => ✋ register
    })
    log.Fatal(app.Listen(":3000"))
}
```

#### [**静态文件**](https://docs.gofiber.io/api/app#static)服务

```go
func main() {
    app := fiber.New()
    app.Static("/", "./public")
    // => http://localhost:3000/js/script.js
    // => http://localhost:3000/css/style.css
    app.Static("/prefix", "./public")
    // => http://localhost:3000/prefix/js/script.js
    // => http://localhost:3000/prefix/css/style.css
    app.Static("*", "./public/index.html")
    // => http://localhost:3000/any/path/shows/index/html
    log.Fatal(app.Listen(":3000"))
}
```

#### [**中间件**](https://docs.gofiber.io/middleware)和[**Next**](https://docs.gofiber.io/api/ctx#next)

```go
func main() {
    app := fiber.New()
    // Match any route
    app.Use(func(c *fiber.Ctx) error {
        fmt.Println("🥇 First handler")
        return c.Next()
    })

    // Match all routes starting with /api
    app.Use("/api", func(c *fiber.Ctx) error {
        fmt.Println("🥈 Second handler")
        return c.Next()
    })

    // GET /api/register
    app.Get("/api/list", func(c *fiber.Ctx) error {
        fmt.Println("🥉 Last handler")
        return c.SendString("Hello, World 👋!")
    })

    log.Fatal(app.Listen(":3000"))
}
```

📚 展示更多代码示例

### 模版引擎

📖 [配置](https://docs.gofiber.io/fiber#config) 📖 [模版引擎](https://github.com/gofiber/template) 📖 [渲染](https://docs.gofiber.io/context#render)

如果未设置模版引擎，则`Fiber`默认使用[html/template](https://golang.org/pkg/html/template/)。

如果您要执行部分模版或使用其他引擎，例如[amber](https://github.com/eknkc/amber)，[handlebars](https://github.com/aymerick/raymond)，[mustache](https://github.com/cbroglie/mustache)或者[pug](https://github.com/Joker/jade)等等…

请查看我们的[Template](https://github.com/gofiber/template)包，该包支持多个模版引擎。

```go
package main
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/pug"
)
func main() {
    // You can setup Views engine before initiation app:
    app := fiber.New(fiber.Config{
        Views: pug.New("./views", ".pug"),
    })

    // And now, you can call template `./views/home.pug` like this:
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("home", fiber.Map{
            "title": "Homepage",
            "year":  1999,
        })
    })

    log.Fatal(app.Listen(":3000"))
}
```

### 组合路由链

📖 [路由分组](https://docs.gofiber.io/application#group)

```go
func middleware(c *fiber.Ctx) error {
    fmt.Println("Don't mind me!")
    return c.Next()
}

func handler(c *fiber.Ctx) error {
    return c.SendString(c.Path())
}

func main() {
    app := fiber.New()
    // Root API route
    api := app.Group("/api", middleware) // /api
    // API v1 routes
    v1 := api.Group("/v1", middleware) // /api/v1
    v1.Get("/list", handler)           // /api/v1/list
    v1.Get("/user", handler)           // /api/v1/user
    // API v2 routes
    v2 := api.Group("/v2", middleware) // /api/v2
    v2.Get("/list", handler)           // /api/v2/list
    v2.Get("/user", handler)           // /api/v2/user
    // ...
}
```

### 日志中间件

📖 [Logger](https://docs.gofiber.io/middleware/logger)

```go
package main
import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)
func main() {
    app := fiber.New()
    app.Use(logger.New())
    // ...
    log.Fatal(app.Listen(":3000"))
}
```

### 跨域资源共享(CORS)中间件

📖 [CORS](https://docs.gofiber.io/middleware/cors)

```go
import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)
func main() {
    app := fiber.New()
    app.Use(cors.New())
    // ...
    log.Fatal(app.Listen(":3000"))
}
```

通过在请求头中设置`Origin`传递任何域来检查 CORS：

```
curl -H "Origin: http://example.com" --verbose http://localhost:3000
```

### 自定义 404 响应

📖 [HTTP Methods](https://docs.gofiber.io/application#http-methods)

```go
func main() {
    app := fiber.New()
    app.Static("/", "./public")
    app.Get("/demo", func(c *fiber.Ctx) error {
        return c.SendString("This is a demo!")
    })
    app.Post("/register", func(c *fiber.Ctx) error {
        return c.SendString("Welcome!")
    })
    // Last middleware to match anything
    app.Use(func(c *fiber.Ctx) error {
        return c.SendStatus(404)
        // => 404 "Not Found"
    })
    log.Fatal(app.Listen(":3000"))
}
```

### JSON 响应

📖 [JSON](https://docs.gofiber.io/ctx#json)

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
func main() {
    app := fiber.New()
    app.Get("/user", func(c *fiber.Ctx) error {
        return c.JSON(&User{"John", 20})
        // => {"name":"John", "age":20}
    })
    app.Get("/json", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Hi John!",
        })
        // => {"success":true, "message":"Hi John!"}
    })
    log.Fatal(app.Listen(":3000"))
}
```

### 升级到 WebSocket

📖 [Websocket](https://github.com/gofiber/websocket)

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/websocket"
)
func main() {
  app := fiber.New()
  app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    for {
      mt, msg, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        break
      }
      log.Printf("recv: %s", msg)
      err = c.WriteMessage(mt, msg)
      if err != nil {
        log.Println("write:", err)
        break
      }
    }
  }))
  log.Fatal(app.Listen(":3000"))
  // ws://localhost:3000/ws
}
```

### 恢复(panic)中间件

📖 [Recover](https://docs.gofiber.io/middleware/recover)

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/recover"
)
func main() {
    app := fiber.New()
    app.Use(recover.New())
    app.Get("/", func(c *fiber.Ctx) error {
        panic("normally this would crash your app")
    })
    log.Fatal(app.Listen(":3000"))
}
```

## 内置中间件

以下为`fiber`框架的内置中间件：

| 中间件                                                       | 描述                                                     |
| :----------------------------------------------------------- | :------------------------------------------------------- |
| [basicauth](https://github.com/gofiber/fiber/tree/master/middleware/basicauth) | basicauth中间件提供HTTP基本身份验证                      |
| [compress](https://github.com/gofiber/fiber/tree/master/middleware/compress) | Fiber的压缩中间件，它支持deflate，gzip 和 brotli（默认） |
| [cache](https://github.com/gofiber/fiber/tree/master/middleware/cache) | 拦截和响应缓存                                           |
| [cors](https://github.com/gofiber/fiber/tree/master/middleware/cors) | 跨域处理                                                 |
| [csrf](https://github.com/gofiber/fiber/tree/master/middleware/csrf) | CSRF攻击防护                                             |
| [filesystem](https://github.com/gofiber/fiber/tree/master/middleware/filesystem) | Fiber的文件系统中间件                                    |
| [favicon](https://github.com/gofiber/fiber/tree/master/middleware/favicon) | favicon图标                                              |
| [limiter](https://github.com/gofiber/fiber/tree/master/middleware/limiter) | `请求频率限制`中间件，用于控制API请求频率                |
| [logger](https://github.com/gofiber/fiber/tree/master/middleware/logger) | HTTP请求与响应日志记录器                                 |
| [pprof](https://github.com/gofiber/fiber/tree/master/middleware/pprof) | pprof 中间件                                             |
| [proxy](https://github.com/gofiber/fiber/tree/master/middleware/proxy) | 请求代理                                                 |
| [requestid](https://github.com/gofiber/fiber/tree/master/middleware/requestid) | 为每个请求添加一个requestid。                            |
| [recover](https://github.com/gofiber/fiber/tree/master/middleware/recover) | `Recover`中间件将程序从`panic`状态中恢复过来             |
| [timeout](https://github.com/gofiber/fiber/tree/master/middleware/timeout) | 添加请求的最大时间，如果超时，则转发给ErrorHandler。     |

## 外部中间件

有`fiber`团队维护的外部中间件

| 中间件                                            | 描述                                      |
| :------------------------------------------------ | :---------------------------------------- |
| [adaptor](https://github.com/gofiber/adaptor)     | `net/http` 与 `Fiber`请求的相互转换适配器 |
| [helmet](https://github.com/gofiber/helmet)       | 可设置各种HTTP Header来保护您的应用       |
| [jwt](https://github.com/gofiber/jwt)             | JSON Web Token (JWT) 中间件               |
| [keyauth](https://github.com/gofiber/keyauth)     | 提供基于密钥的身份验证                    |
| [rewrite](https://github.com/gofiber/rewrite)     | URL路径重写                               |
| [session](https://github.com/gofiber/session)     | Session中间件                             |
| [template](https://github.com/gofiber/template)   | 模板引擎                                  |
| [websocket](https://github.com/gofiber/websocket) | Fasthttp WebSocket 中间件                 |

## 第三方中间件

这是由`Fiber`社区创建的中间件列表，如果您想看到自己的中间件，请创建`PR`。

- [arsmn/fiber-casbin](https://github.com/arsmn/fiber-casbin)
- [arsmn/fiber-introspect](https://github.com/arsmn/fiber-introspect)
- [arsmn/fiber-swagger](https://github.com/arsmn/fiber-swagger)
- [arsmn/gqlgen](https://github.com/arsmn/gqlgen)
- [codemicro/fiber-cache](https://github.com/codemicro/fiber-cache)
- [sujit-baniya/fiber-boilerplate](https://github.com/sujit-baniya/fiber-boilerplate)
- [juandiii/go-jwk-security](https://github.com/juandiii/go-jwk-security)
- [kiyonlin/fiber_limiter](https://github.com/kiyonlin/fiber_limiter)
- [shareed2k/fiber_limiter](https://github.com/shareed2k/fiber_limiter)
- [shareed2k/fiber_tracing](https://github.com/shareed2k/fiber_tracing)
- [thomasvvugt/fiber-boilerplate](https://github.com/thomasvvugt/fiber-boilerplate)
- [ansrivas/fiberprometheus](https://github.com/ansrivas/fiberprometheus)
- [LdDl/fiber-long-poll](https://github.com/LdDl/fiber-long-poll)
- [K0enM/fiber_vhost](https://github.com/K0enM/fiber_vhost)
- [theArtechnology/fiber-inertia](https://github.com/theArtechnology/fiber-inertia)
