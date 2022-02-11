---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "context - 上下文管理"  # 文章标题
url:  "posts/go/libraries/standard/context"  # 设置网页链接，默认使用文件名
tags: [ "go", "context", "goroutine" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

在 Go http  标准库的 Server 中，每一个请求在都有一个对应的 goroutine 去处理。请求处理函数通常会启动额外的 goroutine 用来访问后端服务，比如数据库和 RPC 服务。用来处理一个请求的 goroutine 通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的 token、请求的截止时间。 当一个请求被取消或超时时，所有用来处理该请求的 goroutine 都应该迅速退出，然后系统才能释放这些 goroutine 占用的资源。

- [为什么需要 Context](#为什么需要-context)
	- [全局变量方式](#全局变量方式)
	- [通道方式](#通道方式)
	- [官方版方式](#官方版方式)
- [Context 接口](#context-接口)
	- [Background() 和 TODO()](#background-和-todo)
- [With 系列函数](#with-系列函数)
	- [WithCancel](#withcancel)
	- [WithDeadline](#withdeadline)
	- [WithTimeout](#withtimeout)
	- [WithValue](#withvalue)
- [使用 Context 的注意事项](#使用-context-的注意事项)
- [客户端超时取消示例](#客户端超时取消示例)

## 为什么需要 Context

思考如下情况：如何接收外部命令实现退出？

```go
var wg sync.WaitGroup

func worker() {
	for {
		fmt.Println("working")
		time.Sleep(time.Second)
	}
	// 思考如何接收外部命令实现退出
	wg.Done()
}

func main() {
	wg.Add(1)
	go worker()
	wg.Wait()
	fmt.Println("over")
}
```

### 全局变量方式

```go
var wg sync.WaitGroup
var exit bool

func worker() {
	for {
		fmt.Println("working")
		time.Sleep(time.Second)
		if exit {
			break
		}
	}
	wg.Done()
}

func main() {
	wg.Add(1)
	go worker()
	time.Sleep(time.Second * 3)
	exit = true
	wg.Wait()
	fmt.Println("over")
}
```

存在问题：

- 作为模块被调用时不容易统一；
- 若函数中再次启动 goroutine 就难以控制了。

### 通道方式

```go
var wg sync.WaitGroup

func worker(exit chan bool) {
LOOP:
	for {
		fmt.Println("working")
		time.Sleep(time.Second)
		select {
		case <-exit:
			break LOOP
		default:
		}
	}
	wg.Done()
}

func main() {
	wg.Add(1)
	exit := make(chan bool)
	go worker(exit)
	time.Sleep(time.Second * 3)
	exit <- true
	close(exit)
	wg.Wait()
	fmt.Println("over")
}
```

存在问题：

- 作为模块被调用时不容易统一；
- 必须维护一个公共的 channel。

### 官方版方式

```go
var wg sync.WaitGroup

func worker(ctx context.Context) {
LOOP:
	for {
		fmt.Println("working")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	wg.Done()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 3)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
```

当子 goroutine 又开启另外一个 goroutine 时，只需要将 ctx 传入即可。

## Context 接口

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

- `Deadline` 方法需要返回当前 `Context` 被取消的时间，也就是完成工作的截止时间；
- `Done` 方法需要返回一个 `Channel`，这个 `Channel` 会在当前工作完成或者上下文被取消之后关闭，多次调用 `Done` 方法会返回同一个 `Channel`；
- `Err` 方法会返回当前 `Context` 结束的原因，它只会在 `Done` 返回的 `Channel` 被关闭时才会返回非空的值；
  - 如果当前 `Context` 被取消就会返回 `Canceled` 错误；
  - 如果当前 `Context` 超时就会返回 `DeadlineExceeded` 错误；
- `Value` 方法会从 `Context` 中返回键对应的值，对于同一个上下文来说，多次调用 `Value` 并传入相同的 `Key` 会返回相同的结果，该方法仅用于传递跨 API 和进程间跟请求域的数据。

### Background() 和 TODO()

Go 内置两个函数：`Background()` 和 `TODO()`，这两个函数分别返回一个实现了 `Context` 接口的 `background` 和 `todo`。

`Background()` 主要用于 `main` 函数、初始化以及测试代码中，作为 `Context` 这个树结构的最顶层的 `Context`，也就是根 `Context`。

`TODO()`，它目前还不知道具体的使用场景，如果我们不知道该使用什么 `Context` 的时候，可以使用这个。

`background` 和 `todo` 本质上都是 `emptyCtx` 结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的 `Context`。

## With 系列函数

### WithCancel

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

`WithCancel` 返回带有新 `Done` 通道的父节点的副本。当调用返回的 `cancel` 函数或当关闭父上下文的 `Done` 通道时，将关闭返回上下文的 `Done` 通道，无论先发生什么情况。

取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用 `cancel`。

```go
func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
```

上面的示例代码中，`gen` 函数在单独的 goroutine 中生成整数并将它们发送到返回的通道。 `gen` 的调用者在使用生成的整数之后需要取消上下文，以免 `gen` 启动的内部 goroutine 发生泄漏。

### WithDeadline

```go
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
```

返回父上下文的副本，并将 deadline 调整为不迟于 d。如果父上下文的 deadline 已经早于 d，则 `WithDeadline(parent, d)` 在语义上等同于父上下文。当截止日过期时，当调用返回的 cancel 函数时，或者当父上下文的 Done 通道关闭时，返回上下文的 Done 通道将被关闭，以最先发生的情况为准。

取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用 cancel。

```go
func main() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
```

上面的代码中，定义了一个 50 毫秒之后过期的 deadline，然后调用 `context.WithDeadline(context.Background(), d)` 得到一个上下文 `ctx` 和一个取消函数 `cancel`，使用一个 `select` 让主程序陷入等待，等待 1 秒后打印 `overslept` 退出或者等待 `ctx` 过期后退出。 因为 `ctx` 50 秒后就过期，所以 `ctx.Done()` 会先接收到值，上面的代码会打印 `ctx.Err()`。

### WithTimeout

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithTimeout` 返回 `WithDeadline(parent, time.Now().Add(timeout))`。

取消上下文将释放与其相关的资源，因此代码应该在此上下文中运行的操作完成后立即调用 `cancel`，通常用于数据库或者网络连接的超时控制。

```go
var wg sync.WaitGroup

func worker(ctx context.Context) {
LOOP:
	for {
		fmt.Println("connecting ...")
		time.Sleep(time.Millisecond * 10)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	fmt.Println("done")
	wg.Done()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
```

### WithValue

将请求作用域的数据与 Context 对象建立关系

```go
func WithValue(parent Context, key, val interface{}) Context
```

`WithValue` 返回父节点的副本，其中与 key 关联的值为 val。

仅对 API 和进程间传递请求域的数据使用上下文值，而不是使用它来传递可选参数给函数。

所提供的键必须是可比较的，并且不应该是 `string` 类型或任何其他内置类型，以避免使用上下文在包之间发生冲突。`WithValue` 的用户应该为键定义自己的类型。为了避免在分配给 `interface{}` 时进行分配，上下文键通常具有具体类型 `struct{}`。或者，导出的上下文关键变量的静态类型应该是指针或接口。

```go
type TraceCode string

var wg sync.WaitGroup

func worker(ctx context.Context) {
	key := TraceCode("TRACE_CODE")
	traceCode, ok := ctx.Value(key).(string)
	if !ok {
		fmt.Println("invalid trace code")
	}
LOOP:
	for {
		fmt.Printf("worker, trace code:%s\n", traceCode)
		time.Sleep(time.Millisecond * 10)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	fmt.Println("done")
	wg.Done()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "12512312234")
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
```

## 使用 Context 的注意事项

- 推荐以参数的方式显示传递 `Context`；
- 以 Context 作为参数的函数方法，应该把 `Context` 作为第一个参数；
- 给一个函数方法传递 `Context` 的时候，不要传递 `nil`，如果不知道传递什么，就使用 `context.TODO()`；
- `Context` 的 `Value` 相关方法应该传递请求域的必要数据，不应该用于传递可选参数；
- `Context` 是线程安全的，可以放心的在多个 `goroutine` 中传递。

## 客户端超时取消示例

**client.go**

```go
type respData struct {
	resp *http.Response
}

func doCall(ctx context.Context) {
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := http.Client{
		Transport: &transport,
	}
	respChan := make(chan *respData, 1)
	req, _ := http.NewRequest("GET", "http://google.com", nil)

	req = req.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		resp, _ := client.Do(req)
		fmt.Printf("client.do resp:%v, _:%v\n", resp, _)
		rd := &respData{
			resp: resp,
		}
		respChan <- rd
		wg.Done()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("call api timeout")
	case result := <-respChan:
		fmt.Println("call server api success")
		defer result.resp.Body.Close()
		data, _ := ioutil.ReadAll(result.resp.Body)
		fmt.Printf("resp:%v\n", string(data))
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	doCall(ctx)
}
```
