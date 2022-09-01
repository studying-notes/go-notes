---
date: 2020-11-15T22:33:00+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 测试标准库公共组件源码分析"  # 文章标题
url:  "posts/go/libraries/standard/testing/common_src"  # 设置网页永久链接
tags: [ "Go", "common-src" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 简介

我们知道单元测试函数需要传递一个 `testing.T` 类型的参数，而性能测试函数需要传递一个 `testing.B` 类型的参数，该参数可用于控制测试的流程，比如标记测试失败等。

`testing.T` 和 `testing.B` 属于 `testing` 包中的两个数据类型，该类型提供一系列的方法用于控制函数执行流程，考虑到二者有一定的相似性，所以 Go 实现时抽象出一个 `testing.common` 作为一个基础类型，而 `testing.T` 和 `testing.B` 则属于 `testing.common` 的扩展。

## 数据结构

```go
type common struct {
	mu      sync.RWMutex        // 读写锁，仅用于控制本数据内的成员访问
	output  []byte              // 存储当前测试产生的日志，每产生一条日志则追加到该切片中，待测试结束后再一并输出
	w       io.Writer           // 子测试执行结束需要把产生的日志输送到父测试中的 output 切片中，传递时需要考虑缩进等格式调整，通过 w 把日志传递到父测试
	ran     bool                // 仅表示是否已执行过。比如，根据某个规范筛选测试，如果没有测试被匹配到的话，则 common.ran 为 false，表示没有测试运行过
	failed  bool                // 如果当前测试执行失败，则置为 tru
	skipped bool                // 标记当前测试是否已跳过
	done    bool                // 表示当前测试及其子测试已结束，此状态下再执行 Fail() 之类的方法标记测试状态会产生 panic
	helpers map[string]struct{} // 标记当前为函数为 help 函数，其中打印的日志，在记录日志时不会显示其文件名及行号

	chatty     bool   // 对应命令行中的 -v 参数，默认为 false，true 则打印更多详细日志
	finished   bool   // 如果当前测试结束，则置为 true
	hasSub     int32  // 标记当前测试是否包含子测试，当测试使用 t.Run() 方法启动子测试时，t.hasSub 则置为 1
	raceErrors int    // 竞态检测错误数
	runner     string // 执行当前测试的函数名

	parent   *common // 如果当前测试为子测试，则置为父测试的指针
	level    int       // 测试嵌套层数，比如创建子测试时，子测试嵌套层数就会加 1
	creator  []uintptr // 测试函数调用栈
	name     string    // 记录每个测试函数名
	start    time.Time // 记录测试开始的时间
	duration time.Duration // 记录测试所花费的时间
	barrier  chan bool // 用于控制父测试和子测试执行的 channel，如果测试为 Parallel，则会阻塞等待父测试结束后再继续
	signal   chan bool // 通知当前测试结束
	sub      []*T      // 子测试列表
}
```

## 成员方法

### common.Name()

```go
// Name returns the name of the running test or benchmark.
func (c *common) Name() string {
	return c.name
}
```

该方法直接返回 common 结构体中存储的名称。

### common.Fail()

```go

// Fail marks the function as having failed but continues execution.
func (c *common) Fail() {
	if c.parent != nil {
		c.parent.Fail()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// c.done needs to be locked to synchronize checks to c.done in parent tests.
	if c.done {
		panic("Fail in goroutine after " + c.name + " has completed")
	}
	c.failed = true
}
```

Fail() 方法会标记当前测试为失败，然后继续运行，并不会立即退出当前测试。如果是子测试，则除了标记当前测试结果外还通过 `c.parent.Fail()` 来标记父测试失败。

### common.FailNow()

```go
func (c *common) FailNow() {
	c.Fail()
	c.finished = true
	runtime.Goexit()
}
```

FailNow() 内部会调用 Fail() 标记测试失败，还会标记测试结束并退出当前测试协程。

可以简单的把一个测试理解为一个协程，FailNow() 只会退出当前协程，并不会影响其他测试协程，但要保证在当前测试协程中调用 FailNow() 才有效，不可以在当前测试创建的协程中调用该方法。

### common.log()

```go
func (c *common) log(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.output = append(c.output, c.decorate(s)...)
}
```

common.log() 为内部记录日志入口，日志会统一记录到 common.output 切片中，测试结束时再统一打印出来。

日志记录时会调用 common.decorate() 进行装饰，即加上文件名和行号，还会做一些其他格式化处理。

调用 common.log() 的方法，有 Log()、Logf()、Error()、Errorf()、Fatal()、Fatalf()、Skip()、Skipf() 等。

单元测试中记录的日志只有在执行失败或指定了 `-v` 参数才会打印，否则不会打印。而在性能测试中则总是被打印出来，因为是否打印日志有可能影响性能测试结果。

### common.Log(args ...interface{})

```go
func (c *common) Log(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
}
```

common.Log() 方法用于记录简单日志，通过 fmt.Sprintln() 方法生成日志字符串后记录。

### common.Logf(format string, args ...interface{})

```go
func (c *common) Logf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
}
```

common.Logf() 方法用于格式化记录日志，通过 fmt.Sprintf() 生成字符串后记录。

### common.Error(args ...interface{})

```go
// Error is equivalent to Log followed by Fail.
func (c *common) Error(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.Fail()
}
```

common.Error() 方法等同于 common.Log()+common.Fail()，即记录日志并标记失败，但测试继续进行。

### common.Errorf(format string, args ...interface{})

```go
// Errorf is equivalent to Logf followed by Fail.
func (c *common) Errorf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.Fail()
}
```

common.Errorf() 方法等同于 common.Logf()+common.Fail()，即记录日志并标记失败，但测试继续进行。

### common.Fatal(args ...interface{})

```go
// Fatal is equivalent to Log followed by FailNow.
func (c *common) Fatal(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.FailNow()
}
```

common.Fatal() 方法等同于 common.Log()+common.FailNow()，即记录日志、标记失败并退出当前测试。

### common.Fatalf(format string, args ...interface{})

```go
// Fatalf is equivalent to Logf followed by FailNow.
func (c *common) Fatalf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.FailNow()
}
```

common.Fatalf() 方法等同于 common.Logf()+common.FailNow()，即记录日志、标记失败并退出当前测试。

### common.skip()

```go
func (c *common) skip() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.skipped = true
}
```

common.skip() 方法标记当前测试为已跳过状态，比如测试中检测到某种条件，不再继续测试。该函数仅标记测试跳过，与测试结果无关。测试结果仍然取决于 common.failed。

### common.SkipNow()

```go
func (c *common) SkipNow() {
	c.skip()
	c.finished = true
	runtime.Goexit()
}
```

common.SkipNow() 方法标记测试跳过，并标记测试结束，最后退出当前测试。

### common.Skip(args ...interface{})

```go
// Skip is equivalent to Log followed by SkipNow.
func (c *common) Skip(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.SkipNow()
}
```

common.Skip() 方法等同于 common.Log()+common.SkipNow()。

### common.Skipf(format string, args ...interface{})

```go
// Skipf is equivalent to Logf followed by SkipNow.
func (c *common) Skipf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.SkipNow()
}
```

common.Skipf() 方法等同于 common.Logf() + common.SkipNow()。

### common.Helper()

```go
// Helper marks the calling function as a test helper function.
// When printing file and line information, that function will be skipped.
// Helper may be called simultaneously from multiple goroutines.
func (c *common) Helper() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.helpers == nil {
		c.helpers = make(map[string]struct{})
	}
	c.helpers[callerName(1)] = struct{}{}
}
```

common.Helper() 方法标记当前函数为 `help` 函数，所谓 `help` 函数，即其中打印的日志，不记录 `help` 函数的函数名及行号，而是记录上一层函数的函数名和行号。

```go

```
