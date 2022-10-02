---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 错误处理"  # 文章标题
url:  "posts/go/docs/grammar/error"  # 设置网页永久链接
tags: [ "Go", "error" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 概述

对于那些将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，通常是最后一个来传递错误信息。如果导致失败的原因只有一个，额外的返回值可以是一个布尔值，通常被命名为 `ok`。比如，当从一个 `map` 查询一个结果时，可以通过额外的布尔值判断是否成功：

```go
if v, ok := m["key"]; ok {
    return v
}
```

在 C 语言中，默认采用一个整数类型的 `errno` 来表达错误，这样就可以根据需要定义多种错误类型。在 Go 语言中，`syscall.Errno` 就是对应C语言中 `errno` 类型的错误。在 `syscall` 包中的接口，如果有返回错误的话，底层也是 `syscall.Errno` 错误类型。

比如我们通过 `syscall` 包的接口来修改文件的模式时，如果遇到错误我们可以通过将 `err` 强制断言为 `syscall.Errno` 错误类型来处理：

```go
err := syscall.Chmod(":invalid path:", 0666)
if err != nil {
    log.Fatal(err.(syscall.Errno))
}
```

当返回的错误值不是 `nil` 时，我们可以通过调用 `error` 接口类型的 `Error` 方法来获得字符串类型的错误信息。

在 Go 语言中，**错误被认为是一种可以预期的结果**；而**异常则是一种非预期的结果**，发生异常可能表示程序中存在 BUG 或发生了其它不可控的问题。Go 语言推荐使用 `recover` 函数**将内部异常转为错误处理**，这使得用户可以真正的关心业务相关的错误处理。

如果某个接口简单地将所有普通的错误当做异常抛出，将会使错误信息杂乱且没有价值。就像在 `main` 函数中直接捕获全部一样，是没有意义的：

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Fatal(r)
        }
    }()

    ...
}
```

捕获异常不是最终的目的。如果异常不可预测，直接输出异常信息是最好的处理方式。

## 错误处理策略

示例：函数需要打开两个文件，然后将其中一个文件的内容复制到另一个文件。

```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }

    written, err = io.Copy(dst, src)
    dst.Close()
    src.Close()
    return
}
```

上面的代码虽然能够工作，但是隐藏一个 bug。如果第一个 `os.Open` 调用成功，但是第二个 `os.Create` 调用失败，那么会在没有释放 `src` 文件资源的情况下返回。虽然我们可以通过在第二个返回语句前添加 `src.Close()` 调用来修复这个 BUG；但是当代码变得复杂时，类似的问题将很难被发现和修复。我们可以通过 `defer` 语句来确保每个被正常打开的文件都能被正常关闭：

```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}
```

`defer` 语句可以让我们在打开文件时马上思考如何关闭文件。不管函数如何返回，文件关闭语句始终会被执行。同时 `defer` 语句可以保证，即使 `io.Copy` 发生了异常，文件依然可以安全地关闭。

Go 语言中的导出函数一般不抛出异常，一个未受控的异常可以看作是程序的 BUG。但是对于那些提供类似 Web 服务的框架而言；它们经常需要接入第三方的中间件。因为第三方的中间件是否存在 BUG 是否会抛出异常，Web 框架本身是不能确定的。为了提高系统的稳定性，Web 框架一般会通过 `recover` 来防御性地**捕获所有处理流程中可能产生的异常，然后将异常转为普通的错误返回**。

以 JSON 解析器为例，说明 `recover` 的使用场景。考虑到 JSON 解析器的复杂性，即使某个语言解析器目前工作正常，也无法肯定它没有漏洞。因此，当某个异常出现时，我们不会选择让解析器崩溃，而是会将 `panic` 异常当作普通的解析错误，并附加额外信息提醒用户报告此错误。

```go
func ParseJSON(input string) (s *Syntax, err error) {
    defer func() {
        if p := recover(); p != nil {
            err = fmt.Errorf("JSON: internal error: %v", p)
        }
    }()
    // ...parser...
}
```

标准库中的 `json` 包，在内部递归解析 JSON 数据的时候如果遇到错误，会通过**抛出异常的方式来快速跳出深度嵌套的函数调用**，然后**由最外一级的接口**通过 `recover` 捕获 `panic`，然后返回相应的错误信息。

Go 语言库的实现习惯: 即使在包内部使用了`panic`，但是在导出函数时会被转化为明确的错误值。

## 获取错误的上下文

有时候为了方便上层用户理解；底层实现者会将底层的错误重新包装为新的错误类型返回给用户：

```go
if _, err := html.Parse(resp.Body); err != nil {
    return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
```

上层用户在遇到错误时，可以很容易从业务层面理解错误发生的原因。但是鱼和熊掌总是很难兼得，在上层用户获得新的错误的同时，我们也丢失了底层最原始的错误类型（只剩下错误描述信息了）。

为了记录这种错误类型在包装的变迁过程中的信息，我们一般会定义一个辅助的 `WrapError` 函数，用于包装原始的错误，同时保留完整的原始错误类型。为了问题定位的方便，同时也为了能记录错误发生时的函数调用状态，我们很多时候希望在出现致命错误的时候保存完整的函数调用信息。同时，为了支持 RPC 等跨网络的传输，我们可能要需要将错误序列化为类似 JSON 格式的数据，然后再从这些数据中将错误解码恢复出来。

为此，我们可以定义自己的 `errors` 包，里面是以下的错误类型：

```go
type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

type Error interface {
	Caller() []CallerInfo
	Wrapped() []error
	Code() int
	error
}

```

参考 `github.com/chai2010/errors`。在 Go 语言中，错误处理也有一套独特的编码风格。检查某个子函数是否失败后，我们通常**将处理失败的逻辑代码放在处理成功的代码之前**。如果某个错误会导致函数返回，那么成功时的逻辑代码不应放在 `else` 语句块中，而应直接放在函数体中。

```go
f, err := func()(int, error){}
if err != nil {
    // 失败的情形, 马上返回错误
}

// 正常的处理流程
```

Go 语言中大部分函数的代码结构几乎相同，首先是一系列的初始检查，用于防止错误发生，之后是函数的实际逻辑。

## 错误的错误返回

Go 语言中的错误是一种接口类型。接口信息中包含了原始类型和原始的值。只有**当接口的类型和原始的值都为空**的时候，接口的值才对应 `nil`。其实**当接口中类型为空的时候，原始值必然也是空的**；反之，当**接口对应的原始值为空的时候，接口对应的原始类型并不一定为空的**。

在下面的例子中，试图返回自定义的错误类型，当没有错误的时候返回 `nil`：

```go
func returnsError() error {
    var p *MyError = nil
    if bad() {
        p = ErrBad
    }
    return p // Will always return a non-nil error.
}
```

但是，最终返回的结果其实并非是 `nil`：是一个正常的错误，错误的值是一个 `MyError` 类型的空指针。下面是改进的 `returnsError`：

```go
func returnsError() error {
    if bad() {
        return (*MyError)(err)
    }
    return nil
}
```

因此，在处理错误返回值的时候，没有错误的返回值最好直接写为 `nil`。

Go 语言作为一个强类型语言，不同类型之间必须要显式的转换（而且必须有相同的基础类型）。但是，Go 语言中 `interface` 是一个例外：非接口类型到接口类型，或者是接口类型之间的转换都是隐式的。这是为了支持鸭子类型，当然会牺牲一定的安全性。

```go

```
