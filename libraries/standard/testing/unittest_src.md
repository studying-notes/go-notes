---
date: 2020-11-16T18:20:59+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 单元测试源码分析"  # 文章标题
url:  "posts/go/libraries/standard/testing/unittest_src"  # 设置网页永久链接
tags: [ "Go", "unittest-src" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 数据结构

源码包 `src/testing/testing.go:T` 定义了其数据结构：

```go
type T struct {
	common
	isParallel bool
	context    *testContext // For running tests and subtests.
}
```

其成员简单介绍如下：

* common：即前面绍的 testing.common
* isParallel：表示当前测试是否需要并发，如果测试中执行了 t.Parallel()，则此值为 true
* context：控制测试的并发调度

因为 context 直接决定了单元测试的调度，在介绍 testing.T 支持的方法前，有必要先了解一下 context。

## testContext

源码包 `src/testing/testing.go : testContext` 定义了其数据结构：

```go
type testContext struct {
    // 匹配器，用于管理测试名称匹配、过滤等
	match *matcher

    // 互斥锁，用于控制 testContext 成员的互斥访问
	mu sync.Mutex

	// 用于通知测试可以并发执行的控制管道，测试并发达到最大限制时，需要阻塞等待该管道的通知事件
	startParallel chan bool

	// 当前并发执行的测试个数
	running int

	// 等待并发执行的测试个数，所有等待执行的测试都阻塞在 startParallel 管道处
	numWaiting int

	// 最大并发数，默认为系统 CPU 数，可以通过参数 `-parallel n` 指定
	maxParallel int
}
```

testContext 实现了两个方法用于控制测试的并发调度。

### testContext.waitParallel()

如果一个测试使用 `t.Parallel()` 启动并发，这个测试并不是立即被并发执行，需要检查当前并发执行的测试数量是否达到最大值，这个检查工作统一放在 testContext.waitParallel() 实现的。

testContext.waitParallel() 函数的源码如下：

```go
func (c *testContext) waitParallel() {
	c.mu.Lock()
	if c.running < c.maxParallel {  // 如果当前运行的测试数未达到最大值，直接返回
		c.running++
		c.mu.Unlock()
		return
	}
	c.numWaiting++                  // 如果当前运行的测试数已达最大值，需要阻塞等待
	c.mu.Unlock()
	<-c.startParallel
}
```

函数实现比较简单，如果当前运行的测试数未达最大值，将 c.running++ 后直接返回即可，否则将 c.numWaiting++ 并阻塞等待其他并发测试结束。

这里有个小细节，阻塞等待后面并没有累加 c.running，因为其他并发的测试结束后也不会递减 c.running，所以这里阻塞返回时也不用累加，一个测试结束，随即另一个测试开始，c.running 个数没有变化。

### testContext.release()

当并发测试结束后，会通过 release() 方法释放一个信号，用于启动其他等待并发测试的函数。

testContext.release() 函数的源码如下：

```go
func (c *testContext) release() {
	c.mu.Lock()
	if c.numWaiting == 0 {         // 如果没有函数在等待，直接返回
		c.running--
		c.mu.Unlock()
		return
	}
	c.numWaiting--                 // 如果有函数在等待，释放一个信号
	c.mu.Unlock()
	c.startParallel <- true // Pick a waiting test to be run.
}
```

## tRunner()

函数 tRunner 用于执行一个测试，在不考虑并发测试、子测试场景下，其处理逻辑如下：

```go
func tRunner(t *T, fn func(t *T)) {
	defer func() {
		t.duration += time.Since(t.start)
		signal := true

		t.report() // 测试执行结束后向父测试报告日志

		t.done = true
		t.signal <- signal // 向调度者发送结束信号
	}()

	t.start = time.Now()
	fn(t)

	t.finished = true
}
```

tRunner 传一个经调度者设置过的 testing.T 参数和一个测试函数，执行时记录开始时间，然后将 testing.T 参数传入测试函数并同步等待其结束。

tRunner 在 defer 语句中记录测试执行耗时，并上报日志，最后发送结束信号。

为了避免困惑，上述代码屏蔽了一些子测试和并发测试的细节，比如，defer 语句中，如果当前测试包含子测试，则需要等所有子测试结束，如果当前测试为并发测试，则需要唤醒其他等待并发的测试。更多细节，等我们分析 Parallel() 和 Run() 时再讨论。

## Run()

Run() 函数的完整函数声明为：

```
func (t *T) Run(name string, f func(t *T)) bool
```

Run() 函数启动一个单独的协程来运行名字为 `name` 的子测试 `f`，并且会阻塞等待其执行结束，
除非子测试 `f` 显式地调用 `t.Parallel()` 将自己变成一个可并行的测试，最后返回 `bool` 类型的测试结果。

比如，当在测试 `func TestXxx(t * testing.T)` 中调用 `Run(name, f)` 时，Run() 将启动一个名为 `TestXxx/name` 的子测试。

另外，需要知道的是所有的测试，包括 `func TestXxx(t * testing.T)` 自身，都是由 `TestMain` 使用 Run() 方法直接或间接启动的。

按照惯例，隐去部分代码后的 Run() 方法如下所示：

```go
func (t *T) Run(name string, f func(t *T)) bool {
	t = &T{ // 创建一个新的testing.T用于执行子测试
		common: common{
			barrier: make(chan bool),
			signal:  make(chan bool),
			name:    testName,    // 测试名字，由name及父测试名字组合而成
			parent:  &t.common,
			level:   t.level + 1, // 子测试层次+1
			chatty:  t.chatty,
		},
		context: t.context, // 子测试的context与父测试相同
	}
	go tRunner(t, f) // 启动协程执行子测试
	if !<-t.signal { // 阻塞等待子测试结束信号，子测试要么执行结束，要么以Parallel()执行。如果信号为'false'，说明出现异常退出
		runtime.Goexit()
	}
	return !t.failed // 返回子测试的执行结果
}
```

每启动一个子测试都会创建一个 testing.T 变量，该变量继承当前测试的部分属性，然后以新协程去执行，
当前测试会在子测试结束后返回子测试的结果。

子测试退出条件要么是子测试执行结束，要么是子测试设置了 Paraller()，否则是异常退出。

## Parallel()

Parallel() 方法将当前测试加入到并发队列中，其实现方法如下所示：

```go
func (t *T) Parallel() {
	t.isParallel = true

	t.duration += time.Since(t.start) // 启动并发测试有可能要等待，等待期间耗时需要剔除，此处相当于先记录当前耗时，并发执行开始后再累加

	t.parent.sub = append(t.parent.sub, t) // 将当前测试加入到父测试的列表中，由父测试调度

	t.signal <- true   // Release calling test. 当前测试即将进入并发模式，标记测试结束，以便父测试不必等待并退出Run()
	<-t.parent.barrier // Wait for the parent test to complete.  等待父测试发送子测试启动信号
	t.context.waitParallel() // 阻塞等待并发调度

	t.start = time.Now() // 开始并发执行，重新标记启动时间，这是第二段耗时
}
```

关于测试耗时统计，看过前面的 testContext 实现我们知道，启动一个并发测试时，当并发数达到最大时，新的并发测试需要等待，那么等待期间的时间消耗不能统计到测试的耗时中，所以需要先计算当前耗时，在真正被并发调度后才清空 t.start 以跳过等待时间。

看过前面的 Run() 方法实现机制后，我们知道一旦子测试以并发模式执行时，需要通知父测试，其通知机制便是向 t.signal 管道中写入一个信号，父测试便从 Run() 方法中唤醒，继续执行。

看过前面的 tRunner() 方法实现机制后，不难理解，父测试唤醒后继续执行，结束后进入 defer 流程中，在 defer 中将启动所有子测试并等待子测试执行结束。

## tRunner()

与简单版的测试执行所不同的是，defer 语句中增加了子测试、并发测试的处理逻辑，相对完整的 tRunner() 代码如下所示：

```go

func tRunner(t *T, fn func(t *T)) {
	t.runner = callerName(0)

	defer func() {
		t.duration += time.Since(t.start) // 进入defer后立即记录测试执行时间,后续流程所花费的时间不应该统计到本测试执行用时中

		if len(t.sub) > 0 { // 如果存在子测试，则启动并等待其完成
			t.context.release() // 减少运行计数
			close(t.barrier)  // 启动子测试
			for _, sub := range t.sub { // 等待所有子测试结束
				<-sub.signal
			}
			if !t.isParallel { // 如果当前测试非并发模式，则等待并发执行，类似于测试函数中执行t.Parallel()
				t.context.waitParallel()
			}
		} else if t.isParallel { // 如果当前测试是并发模式，则释放信号以启动新的测试
			t.context.release()
		}
		t.report() // 测试执行结束后向父测试报告日志

		t.done = true
		t.signal <- signal // 向父测试发送结束信号，以便结束Run()
	}()

	t.start = time.Now() // 记录测试开始时间
	fn(t)

	// code beyond here will not be executed when FailNow is invoked
	t.finished = true
}
```

测试执行结束，进入 defer 后需要启动子测试，启动方法为关闭 t.barrier 管道，然后等待所有子测试执行结束。

需要注意的是，关闭 t.barrier 管道，阻塞在 t.barrier 管道上的协程同样会被唤醒，也是发送信号的一种方式，关于管道的更多实现细节，请参考管道实现原理相关章节。

defer 中，如果检测到当前测试本身也处于并发中，那么结束后需要释放一个信号 t.context.release() 来启动一个等待的测试。

```go

```
