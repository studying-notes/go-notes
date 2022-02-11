---
date: 2020-09-19T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 学习笔记"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/README"  # 设置网页永久链接
tags: [ "go", "README" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 10 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

# Go 学习笔记

> 纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。

## 目录结构

- `assets/images`: 笔记配图
- `assets/templates`: 笔记模板
- `docs`: 基础语法
- `libraries`: 库
  - `libraries/standard`: 标准库
  - `libraries/tripartite`: 第三方库
- `quickstart`: 基础用法
- `src`: 源码示例
  - `src/docs`: 基础语法源码示例
  - `src/libraries/standard`: 标准库源码示例
  - `src/libraries/tripartite`: 第三方库源码示例
  - `src/quickstart`: 基础用法源码示例

## 基础用法

- [Go 安装与配置指南](quickstart/install.md)
- [Go 卸载指南](quickstart/uninstall.md)

## 基础语法

## 库

## 标准库

- [binary - 二进制数据的序列化与反序列化](libraries/standard/binary.md)
- [bufio - 获取用户输入](libraries/standard/bufio.md)
- [fmt - 获取用户输入](libraries/standard/fmt.md)
- [context - 上下文管理](libraries/standard/context.md)
- [exec - 执行终端命令/外部命令](libraries/standard/exec.md)
- [flag - 命令行参数解析](libraries/standard/flag.md)
- [httputil - 网络工具包](libraries/standard/httputil.md)
- [ioutil - IO 操作函数集](libraries/standard/ioutil.md)
- [json - JSON 序列化和反序列化](libraries/standard/json.md)
- [log - 日志](libraries/standard/log.md)
- [sync.Pool - 内存池](libraries/standard/pool.md)
- [pprof - 性能剖析](libraries/standard/pprof.md)
- [rand - 随机数](libraries/standard/rand.md)
- [regexp - 正则表达式](libraries/standard/regexp.md)
- [strconv - 字符串转换其他类型](libraries/standard/strconv.md)
- [strings - 字符串操作](libraries/standard/strings.md)
- [template - 文本模板引擎](libraries/standard/template.md)
- [time - 时间标准库](libraries/standard/time.md)
- [net/http - HTTP 标准库](libraries/standard/net_http.md)
- [image - 图片处理](libraries/standard/image.md)

## 第三方库

- [urfave/cli - 构建 CLI 程序](libraries/tripartite/cli.md)
- [Cobra - 构建 CLI 程序](libraries/tripartite/cobra.md)
- [cron - 启动定时任务](libraries/tripartite/cron.md)
- [fsnotify - 监听文件系统事件](libraries/tripartite/fsnotify.md)
- [gjson - 快速提取 JSON 值](libraries/tripartite/gjson.md)
- [gopsutil - 获取系统运行信息](libraries/tripartite/gopsutil.md)
- [gorm - 数据库操作](libraries/tripartite/gorm.md)
- [logrus - 日志库](libraries/tripartite/logrus.md)
- [zap - 日志库](libraries/tripartite/zap.md)
- [service - 编写 Windows 服务](libraries/tripartite/service.md)
- [go-sqlcipher - 用 SQLCipher 加密 SQLite](libraries/tripartite/sqlcipher.md)
- [sqlx - 扩展标准库 database/sql](libraries/tripartite/sqlx.md)
- [viper - 应用程序的完整配置解决方案](libraries/tripartite/viper.md)
- [webdav - 简单的 WebDAV 服务实现](libraries/tripartite/webdav.md)
- [gout - HTTP 客户端](libraries/tripartite/gout.md)
- [xlsx - 读写 Excel 表格](libraries/tripartite/excel.md)
- [mongo - 操作 MongoDB](libraries/tripartite/mongo.md)
- [mysql - MySQL 操作示例](libraries/tripartite/mysql.md)
- [redis - 操作 Redis](libraries/tripartite/redis.md)
- [go-sqlite3 - SQLite / SQLCipher 操作示例](libraries/tripartite/sqlite.md)

## 基础语法

### 数据结构

{{<card src="posts/go/abc/array">}}
{{<card src="posts/go/abc/string">}}
{{<card src="posts/go/abc/slice">}}
{{<card src="posts/go/abc/map">}}
{{<card src="posts/go/abc/func">}}
{{<card src="posts/go/abc/struct">}}
{{<card src="posts/go/abc/method">}}
{{<card src="posts/go/abc/interface">}}
{{<card src="posts/go/abc/goroutine2">}}
{{<card src="posts/go/abc/channel">}}
{{<card src="posts/go/abc/reflect">}}
{{<card src="posts/go/abc/append">}}
{{<card src="posts/go/abc/iota">}}
{{<card src="posts/go/abc/attention">}}

### 控制结构

{{<card src="posts/go/abc/defer">}}
{{<card src="posts/go/abc/recover">}}
{{<card src="posts/go/abc/error">}}
{{<card src="posts/go/abc/select">}}
{{<card src="posts/go/abc/range">}}
{{<card src="posts/go/abc/range2">}}
{{<card src="posts/go/abc/mutex">}}
{{<card src="posts/go/abc/rwmutex">}}

### 内存管理

{{<card src="posts/go/abc/memory/alloc">}}
{{<card src="posts/go/abc/memory/gc">}}
{{<card src="posts/go/abc/memory/escape">}}

## 并发控制

我们考虑这么一种场景，协程 A 执行过程中需要创建子协程 A1、A2、A3...An，协程 A 创建完子协程后就等待子协程退出。针对这种场景，Go 提供了三种解决方案：

- Channel : 使用 channel 控制子协程
- WaitGroup : 使用信号量机制控制子协程
- Context : 使用上下文控制子协程

三种方案各有优劣，比如 Channel 优点是实现简单，清晰易懂，WaitGroup 优点是子协程个数动态可调整，Context 优点是对子协程派生出来的孙子协程的控制。

{{<card src="posts/go/abc/concurrent/goroutine">}}
{{<card src="posts/go/abc/concurrent/concurrent">}}
{{<card src="posts/go/abc/concurrent/channel">}}
{{<card src="posts/go/abc/concurrent/waitgroup">}}
{{<card src="posts/go/abc/concurrent/context">}}
{{<card src="posts/go/libraries/standard/sync/pool">}}
{{<card src="posts/go/libraries/standard/context">}}

### 语法糖

语法糖（Syntactic Sugar），Go 中最常用的语法糖莫过于赋值符 `:=`，其次，表示函数变参的 `...`。

## 测试与性能

* 单元测试 - 指对软件中的最小可测试单元进行检查和验证，比如对一个函数的测试。
* 性能测试 - 也称基准测试，可以测试一段程序的性能，可以得到时间消耗、内存使用情况的报告。
* 示例测试 - 示例测试，广泛应用于 Go 源码和各种开源框架中，用于展示某个包或某个方法的用法。

{{<card src="posts/go/doc/test/benchmark_test">}}
{{<card src="posts/go/doc/test/unit_test">}}
{{<card src="posts/go/doc/test/example_test">}}
{{<card src="posts/go/doc/test/sub_test">}}
{{<card src="posts/go/doc/test/main_test">}}
{{<card src="posts/go/libraries/standard/pprof">}}

### 深入测试标准库

{{<card src="posts/go/abc/test/common">}}
{{<card src="posts/go/abc/test/tb_interface">}}
{{<card src="posts/go/abc/test/unit">}}
{{<card src="posts/go/abc/test/benchmark">}}
{{<card src="posts/go/abc/test/example">}}
{{<card src="posts/go/abc/test/main">}}
{{<card src="posts/go/abc/test/go_test">}}
{{<card src="posts/go/abc/test/go_test_params">}}
{{<card src="posts/go/abc/test/go_test_benchstat">}}

## 依赖管理

Go 语言依赖管理经历了三个阶段：

- GOPATH；
- vendor；
- Go Module；

{{<card src="posts/go/doc/mod/gopath">}}
{{<card src="posts/go/doc/mod/vendor">}}
{{<card src="posts/go/doc/mod/gomod">}}

## CGO 编程

{{<card src="posts/go/cgo/quickstart">}}
{{<card src="posts/go/cgo/intro">}}
{{<card src="posts/go/cgo/dll">}}
{{<card src="posts/go/cgo/func">}}
{{<card src="posts/go/cgo/link">}}
{{<card src="posts/go/cgo/type">}}
{{<card src="posts/go/cgo/internal">}}
