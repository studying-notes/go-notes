---
date: 2020-10-05T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 工程架构"  # 文章标题
url:  "posts/gin/project/project"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 工程目录结构

基本原则：与功能密切联系不可复用的放在 app 子目录下，可以复用的放在一级目录下。

```
project
│
├─app
│  ├─dao
│  ├─model
│  ├─routers
│  └─service
│           ├─client
│           └─server
├─configs
├─docs
├─global
├─middleware
├─pkg
├─scripts
└─storage
```

- configs：配置文件。
- docs：文档集合。
- global：全局变量。
- app：项目模块。
  - dao：数据访问层（Database Access Object ），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、Elasticsearch等。
  - model：模型层，用于存放model对象。
  - routers：路由相关的逻辑。
  - service：项目核心业务逻辑。
- middleware：HTTP中间件。
- pkg：项目相关的模块包。
- storage：项目生成的临时文件。
- scripts：各类构建、安装、分析等操作的脚本。

## 路由原则

在 RESTful API 中，HTTP 方法对应的行为和动作如下：

- GET：读取和检索动作。
- POST：新增和新建动作。
- PUT：更新动作，用于更新一个完整的资源，要求为幂等，即任意多次执行所产生的影响均与一次执行的影响相同。
- PATCH：更新动作，用于更新某一个资源的一个组成部分。也就是说，当只需更新该资源的某一项时，应该使用 PATCH 而不是 PUT，可以不幂等。
- DELETE：删除动作。

## 基础组件

### 错误码标准化

### 配置管理

```shell
go get -u github.com/spf13/viper
```

### 数据库连接

```shell
go get -u github.com/jinzhu/gorm
```

### 日志系统

```shell
go get -u gopkg.in/natefinch/lumberjack.v2
```

核心功能是把日志写入滚动文件中，该库允许我们设置单日志文件的最大占用空间、最大生存周期、可保留的最多旧文件数等。如果有出现超出设置项的情况，则对日志文件进行滚动处理。

### 响应处理

事先按规范定义好响应结果。

## Swagger 接口文档

```shell
go get -u github.com/swaggo/swag/cmd/swag
get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

```shell
go get -u github.com/alecthomas/template
```

在安装完 Swagger 关联库后，就需要项目里的 API 接口编写注解，以便后续在生成时能够正确地运行。

[![B6hjot.png](https://s1.ax1x.com/2020/11/04/B6hjot.png)](https://imgchr.com/i/B6hjot)

## 接口参数校验

在开发对应的业务模块时，第一步要考虑的问题就是如何进行入参校验。我们需要将整个项目，甚至整个团队的组件给定下来，形成一个通用规范。

开源项目 go-playground/validator 是一个基于标签对结构体和字段进行值验证的一个验证器。

```shell
go get -u github.com/go-playground/validator/v10
```

常见的标签含义

[![B643wR.png](https://s1.ax1x.com/2020/11/04/B643wR.png)](https://imgchr.com/i/B643wR)

## 访问控制

目前市场上比较常见的API访问控制方案有两种，分别是 OAuth 2.0 和 JWT。

- OAuth 2.0：本质上是一个授权的行业标准协议，提供了一整套授权机制的指导标准，常用于使用第三方登录的情况。例如，在登录某些网站时，也可以用第三方站点（例如用微信、QQ、GitHub 账号）关联登录，这些往往是用 OAuth 2.0 的标准实现的。OAuth 2.0 相对会“重”一些，常常还会授予第三方应用获取对应账号的个人基本信息等。
- JWT：与 OAuth 2.0 完全不同，它常用于前后端分离的情况，能够非常便捷地给 API 接口提供安全鉴权。

## 常见中间件

### 日志中间件

[![B64JFx.png](https://s1.ax1x.com/2020/11/04/B64JFx.png)](https://imgchr.com/i/B64JFx)

### 异常捕获处理

[![B64tfK.png](https://s1.ax1x.com/2020/11/04/B64tfK.png)](https://imgchr.com/i/B64tfK)

#### 邮件报警处理

```shell
go get -u gopkg.in/gomail.v2
```

gomail 是一个用于发送电子邮件的简单且高效的第三方开源库，目前只支持使用 SMTP 服务器发送电子邮件。

### 服务信息存储

我们经常需要在进程内上下文设置一些内部信息，既可以是应用名称和应用版本号这类基本信息，也可以是业务属性信息。例如，想要根据不同的租户号获取不同的数据库实例对象，这时就需要在一个统一的地方进行处理。

[![Bc99pD.png](https://s1.ax1x.com/2020/11/04/Bc99pD.png)](https://imgchr.com/i/Bc99pD)

### 接口限流控制

[![Bc9knA.png](https://s1.ax1x.com/2020/11/04/Bc9knA.png)](https://imgchr.com/i/Bc9knA)

```shell
go get -u github.com/juju/ratelimit
```

ratelimit 提供了一个简单又高效的令牌桶实现，可以帮助我们实现限流器的逻辑。

### 统一超时控制

假设应用 A 调用应用 B，应用 B 调用应用 C，如果应用 C 出现问题，则在没有任何约束的情况下仍持续调用，就会导致应用 A、B、C 均出现问题。

[![Bc9A0I.png](https://s1.ax1x.com/2020/11/04/Bc9A0I.png)](https://imgchr.com/i/Bc9A0I)

为了避免出现这种情况，最简单的一个约束点，就是统一在应用程序中针对所有请求都进行一个最基本的超时时间控制。

[![Bc9ZAP.png](https://s1.ax1x.com/2020/11/04/Bc9ZAP.png)](https://imgchr.com/i/Bc9ZAP)

## 链路追踪

项目在不断迭代之后，它可能会涉及许许多多的接口，而这些接口很可能是分布式部署的，既存在着多份副本，又存在着相互调用，并且在各自的调用中还可能包含大量的SQL、HTTP、Redis以及应用的调用逻辑。如果对每一个都“打”调用堆栈的日志来记录，未免太多了。如果不做任何记录，那么在出现问题时，很可能会完全找不到方向。

为了更好地解决这个问题，我们使用分布式链路追踪系统，它能够有效解决可观察性上的一部分问题，即多程序部署可在多环境下调用链路的“观察”。

### OpenTracing 规范

OpenTracing 规范的出现是为了解决不同供应商的分布式追踪系统 API 互不兼容的问题，它提供了一个标准的、与供应商无关的工具框架，可以认为它是一个接入层，下面从多个维度进行分析：

- 从功能上：在 OpenTracing 规范中会提供一系列与供应商无关的 API。
- 从系统上：它能够让开发人员更便捷地对接（新增或替换）追踪系统，只需简单地更改 Tracer 的配置就可以了。
- 从语言上：OpenTracing 规范是跨语言的，不会特定涉及某类语言标准，它通过接口的设计概念去封装一系列 API 的相关功能。
- 从标准上：OpenTracing 规范并不是什么官方标准，它的主体 Cloud Native Computing Foundation（CNCF）并不是官方的标准机构。
- 
总的来说，通过 OpenTracing 规范来对接追踪系统之后，我们就可以很方便地在不同的追踪系统中进行切换，今年用 A 系统，明年用 B 系统。因为它不会与具体的某一个供应商系统产生强捆绑关系。

目前，市面上比较流行的追踪系统的思维模型均起源于 Google 的 Dapper，a Large-Scale Distributed Systems Tracing Infrastructure 论文。OpenTracing 规范也不例外，它有一系列约定的术语概念知识，追踪系统中常见的3个术语含义如表所示。

[![BjCgi9.png](https://s1.ax1x.com/2020/11/11/BjCgi9.png)](https://imgchr.com/i/BjCgi9)

### Jaeger 的使用

Jaeger 是 uber 开源的一个分布式链路追踪系统，受到了Google Dapper 和 OpenZipkin 的启发，目前由 Cloud Native Computing Foundation（CNCF）托管。它提供了分布式上下文传播、分布式交易监控、原因分析、服务依赖性分析、性能/延迟优化分析等等核心功能。

目前，市面上比较流行的分布式追踪系统都已经完全支持 OpenTracing 规范。

### 安装 Jaeger

Jaeger 官方提供了 all-in-one 的安装包，提供了Docker 或已打包好的二进制文件，直接运行即可。这里通过 Docker 的方式安装并启动：

```shell
docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one
```

[![BjPVyV.png](https://s1.ax1x.com/2020/11/11/BjPVyV.png)](https://imgchr.com/i/BjPVyV)

访问 `http://localhost:16686` 可看到 Jaeger Web UI 界面则表示已经成功。

### 在应用中注入追踪

```shell
go get -u github.com/opentracing/opentracing-go
go get -u github.com/uber/jaeger-client-go
```

### 实现日志追踪

在记录日志时，把链路的 SpanID 和 TraceID 也记录进去，这样就可以串联起该次请求的所有请求链路和日志信息情况，而且实现的方式并不难，只需在对应方法的第一个参数中传入上下文（context），并在内部解析此上下文来获取链路信息即可。

### 实现 SQL 追踪

在 Jaeger UI 可视化界面可以查看。

## 应用配置问题

配置文件这种非 .go 文件的文件类型，并不会被打包进二进制文件中。

实际上，**go run** 命令并不像 go build 命令那样可以直接编译输出当前目录，而是将其转化到**临时目录下编译并执行**，是一个相对临时的运行路径。

- go run 命令和 go build 命令的不同之处在于，一个是在临时目录下执行，另一个可手动在编译后的目录下执行，路径的处理方式不同。
- 每次执行 go run 命令之后，生成的新的二进制文件不一定在同一个地方。
- **依赖相对路径读取的文件在没有遵守约定条件时，有可能出现最终路径出错的问题**。

### 命令行参数

在 Go 语言中，可以直接通过 flag 标准库来实现该功能。实现逻辑为，如果存在命令行参数，则优先使用命令行参数，否则使用配置文件中的配置参数。

### 系统环境变量

```go
os.Getenv("env")
```

### 打包进二进制文件中

```shell
go get -u github.com/go-bindata/go-bindata
```

通过 go-bindata 库可以将数据文件转换为 Go 代码。例如，常见的配置文件、资源文件（如 Swagger UI ）等都可以打包进 Go 代码中，这样就可以 “ 摆脱 ” 静态资源文件了。接下来在项目根目录下执行生成命令：

```shell
go-bindata -o configs/config.go -pkg=configs configs/config.yaml
```

执行这条命令后，会将 configs/config.yaml 文件打包，并通过 -o 选项指定的路径输出到 configs/config.go 文件中，再通过设置的 -pkg 选项指定生成的 package name 为 configs，接下来只需执行下述代码，就可以读取对应的文件内容了：

```go
b, _ := configs.Asset("configs/config.yaml")
```

把第三方文件打包进二进制文件后，二进制文件必然增大，而且在常规方法下无法做文件的热更新和监听，必须要重启并且重新打包才能使用最新的内容，因此这种方式是有利有弊的。

## 配置热更新

开源库 fsnotify 是用 Go 语言编写的跨平台文件系统监听事件库，常用于文件监听，因此我们可以借助该库来实现这个功能。

```shell
go get -u github.com/fsnotify/fsnotify
```

## 重启和停止

在开发完成应用程序后，即可将其部署到测试、预发布或生产环境中，这时又涉及一个问题，即持续集成。简单来说，开发人员需要关注的是这个应用程序是不断地进行更新和发布的，也就是说，这个应用程序在发布时，很可能某客户正在使用这个应用。如果直接硬发布，就会造成客户的行为被中断。

为了避免这种情况的发生，我们希望在应用更新或发布时，现有正在处理既有连接的应用不要中断，要先处理完既有连接后再退出。而新发布的应用在部署上去后再开始接收新的请求并进行处理，这样即可避免原来正在处理的连接被中断的问题。

### 信号量

信号是UNIX、类UNIX，以及其他POSIX兼容的操作系统中进程间通信的一种有限制的方式。

它是一种异步的通知机制，用来提醒进程一个事件（硬件异常、程序执行异常、外部发出信号）已经发生。当一个信号发送给一个进程时，操作系统中断了进程正常的控制流程。此时，任何非原子操作都将被中断。如果进程定义了信号的处理函数，那么它将被执行，否则执行默认的处理函数。

我们可以通过 kill-l 查看系统中所支持的所有信号。

### 常用的快捷键

在终端执行特定的组合键可以使系统发送特定的信号给指定进程，并完成一系列的动作。

![](../imgs/signl.png)

因此在按组合键ctrl+c关闭服务端时，会发送希望进程结束的通知（发送SIGINT信号），如果没有进行额外处理，则该进程会直接退出，最终导致正在访问的用户出现无法访问的问题。而平时常用的kill-9 pid命令，其会发送SIGKILL信号给进程，作用是强制中断进程。

### 实现优雅重启和停止

1. 实现目的

- 不关闭现有连接（正在运行中的程序）。
- 新的进程启动并替代旧进程。
- 新的进程接管新的连接。
- 连接要随时响应用户的请求。当用户仍在请求旧进程时要保持连接，新用户应请求新进程，不可以出现拒绝请求的情况。

2. 需要达到的流程

- 替换可执行文件或修改配置文件。
- 发送信号量SIGHUP。
- 拒绝新连接请求旧进程，保证正在处理的连接正常。
- 启动新的子进程。
- 新的子进程开始Accept。
- 系统将新的请求转交新的子进程。
- 旧进程处理完所有旧连接后正常退出。
