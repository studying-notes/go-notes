---
date: 2020-12-02T08:57:22+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 启动定时任务"  # 文章标题
url:  "posts/go/libraries/tripartite/cron"  # 设置网页链接，默认使用文件名
tags: [ "go", "corn" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```shell
go get -u github.com/robfig/cron/v3
```

## 快速开始

```go
package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	
	c.AddFunc("30 * * * *", func() {
		fmt.Println("Every hour on the half hour")
	})

	c.AddFunc("30 3-6,20-23 * * *", func() {
		fmt.Println(".. in the range 3-6am, 8-11pm")
	})

	c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() {
		fmt.Println("Runs at 04:30 Tokyo time every day")
	})

	c.AddFunc("@hourly", func() {
		fmt.Println("Every hour, starting an hour from now")
	})

	c.AddFunc("@daily", func() { 
		fmt.Println("Every day") 
	})

	c.AddFunc("@every 1h30m", func() {
		fmt.Println("Every hour thirty, starting an hour thirty from now")
	})

	c.Start()
	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
}
```

首先创建 cron 对象，这个对象用于管理定时任务。

调用 cron 对象的 AddFunc() 方法向管理器中添加定时任务。AddFunc() 接受两个参数，参数 1 以字符串形式指定触发时间规则，参数 2 是一个无参的函数，每次触发时调用。

`@every 1s` 表示每秒触发一次，`@every` 后加一个时间间隔，表示每隔多长时间触发一次。例如 `@every 1h` 表示每小时触发一次，`@every 1m2s` 表示每隔 1 分 2 秒触发一次。`time.ParseDuration()` 支持的格式都可以用在这里。

调用 c.Start() 启动定时循环。

## 时间表达式

| 字段   | 强制 | 范围  | 特殊符号 |
| ----------   | ---------- | ----------  | ---------- |
| Minutes      | Yes        | 0-59            | * / , - |
| Hours        | Yes        | 0-23            | * / , - |
| Day of month | Yes        | 1-31            | * / , - ? |
| Month        | Yes        | 1-12 or JAN-DEC | * / , - |
| Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ? |

月份和周历名称都是不区分大小写的，SUN/Sun/sun 表示同样的含义。

## 实现秒级控制

```go
cron.New(cron.WithSeconds())
```

## 特殊字符

### Asterisk ( * )

匹配该位置字段的任意值，比如月份字段，就表示每个月。

### Slash ( / )

指定范围的步长，例如将小时域设置为 3-59/15 表示第 3 分钟触发，以后每隔 15 分钟触发一次，因此第 2 次触发为第 18 分钟，第 3 次为 33 分钟……直到分钟大于 59。

### Comma ( , )

列举一些离散的值和多个范围，例如将周域设置为 MON,WED,FRI 表示周一、三和五。

### Hyphen ( - )

表示范围，例如将小时域设置为 9-17 表示上午 9 点到下午 17 点。

### Question mark ( ? )

只能用在月和周的域中，用来代替 *，表示每月/周的任意一天。

## 预定义时间规则

| Entry | Description |Equivalent To |
| ----- | ----- | ----- |
| @yearly (or @annually) | Run once a year, midnight, Jan. 1st | 0 0 1 1 * |
| @monthly | Run once a month, midnight, first of month | 0 0 1 * * |
| @weekly | Run once a week, midnight between Sat/Sun  | 0 0 * * 0 |
| @daily (or @midnight)  | Run once a day, midnight | 0 0 * * * |
| @hourly | Run once an hour, beginning of hour | 0 * * * * |

## 固定时间间隔

```
@every <duration>
```

含义为每隔 duration 触发一次。`<duration>` 会调用 `time.ParseDuration()` 函数解析，所以 `ParseDuration` 支持的格式都可以，例如 1h30m10s。

## 时区

默认当前时区。

```go
cron.New(
    cron.WithLocation(time.UTC))
```

Individual cron schedules may also override the time zone they are to be interpreted in by providing an additional space-separated field at the beginning of the cron spec, of the form "CRON_TZ=Asia/Tokyo".

For example:

```go
// Runs at 6am in time.Local
cron.New().AddFunc("0 6 * * ?", ...)

// Runs at 6am in America/New_York
nyc, _ := time.LoadLocation("America/New_York")
c := cron.New(cron.WithLocation(nyc))
c.AddFunc("0 6 * * ?", ...)

// Runs at 6am in Asia/Tokyo
cron.New().AddFunc("CRON_TZ=Asia/Tokyo 0 6 * * ?", ...)

// Runs at 6am in Asia/Tokyo
c := cron.New(cron.WithLocation(nyc))
c.SetLocation("America/New_York")
c.AddFunc("CRON_TZ=Asia/Tokyo 0 6 * * ?", ...)
```

The prefix "TZ=(TIME ZONE)" is also supported for legacy compatibility.

Be aware that jobs scheduled during daylight-savings leap-ahead transitions will not be run!

## Job Wrappers

A Cron runner may be configured with a chain of job wrappers to add cross-cutting functionality to all submitted jobs. For example, they may be used to achieve the following effects:

```
- Recover any panics from jobs (activated by default)
- Delay a job's execution if the previous run hasn't completed yet
- Skip a job's execution if the previous run hasn't completed yet
- Log each job's invocations
```

Install wrappers for all jobs added to a cron using the `cron.WithChain` option:

```go
cron.New(cron.WithChain(
	cron.SkipIfStillRunning(logger),
))
```

Install wrappers for individual jobs by explicitly wrapping them:

```go
job = cron.NewChain(
	cron.SkipIfStillRunning(logger),
).Then(job)
```

## Thread safety

Since the Cron service runs concurrently with the calling code, some amount of care must be taken to ensure proper synchronization.

All cron methods are designed to be correctly synchronized as long as the caller ensures that invocations have a clear happens-before ordering between them.

## Logging

Cron defines a Logger interface that is a subset of the one defined in github.com/go-logr/logr. It has two logging levels (Info and Error), and parameters are key/value pairs. This makes it possible for cron logging to plug into structured logging systems. An adapter, [Verbose]PrintfLogger, is provided to wrap the standard library *log.Logger.

For additional insight into Cron operations, verbose logging may be activated which will record job runs, scheduling decisions, and added or removed jobs. Activate it with a one-off logger as follows:

```go
cron.New(
	cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))
```

```go

```

```go

```

```go

```
