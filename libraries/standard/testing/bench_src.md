---
date: 2020-11-16T18:30:34+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 基准测试源码分析"  # 文章标题
url:  "posts/go/libraries/standard/testing/bench_src"  # 设置网页永久链接
tags: [ "Go", "bench-src" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [数据结构](#数据结构)
- [关键函数](#关键函数)
  - [B.StartTimer()](#bstarttimer)
  - [B.StopTimer()](#bstoptimer)
  - [B.ResetTimer()](#bresettimer)
  - [B.SetBytes(n int64)](#bsetbytesn-int64)
  - [报告内存信息](#报告内存信息)
- [性能测试是如何启动的](#性能测试是如何启动的)
- [B.N是如何调整的？](#bn是如何调整的)
- [内存是如何统计的？](#内存是如何统计的)

## 数据结构

源码包 `src/testing/benchmark.go:B` 定义了性能测试的数据结构：

```go
type B struct {
	common                         // 与testing.T共享的testing.common，负责记录日志、状态等
	importPath       string // import path of the package containing the benchmark
	context          *benchContext
	N                int            // 目标代码执行次数，不需要用户了解具体值，会自动调整
	previousN        int           // number of iterations in the previous run
	previousDuration time.Duration // total duration of the previous run
	benchFunc        func(b *B)   // 性能测试函数
	benchTime        time.Duration // 性能测试函数最少执行的时间，默认为1s，可以通过参数'-benchtime 10s'指定
	bytes            int64         // 每次迭代处理的字节数
	missingBytes     bool // one of the subbenchmarks does not have bytes set.
	timerOn          bool // 计时启动标志，默认为 false，启动计时为 true
	showAllocResult  bool
	result           BenchmarkResult // 测试结果
	parallelism      int // RunParallel creates parallelism*GOMAXPROCS goroutines
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64  // 计时开始时堆中分配的对象总数
	startBytes  uint64  // 计时开始时时堆中分配的字节总数
	netAllocs uint64 // 计时结束时，堆中增加的对象总数
	netBytes  uint64 // 计时结束时，堆中增加的字节总数
}
```

## 关键函数

### B.StartTimer()

StartTimer() 负责启动计时并初始化内存相关计数，测试执行时会自动调用，一般不需要用户启动。

```go
func (b *B) StartTimer() {
	if !b.timerOn {
		runtime.ReadMemStats(&memStats)     // 读取当前堆内存分配信息
		b.startAllocs = memStats.Mallocs    // 记录当前堆内存分配的对象数
		b.startBytes = memStats.TotalAlloc  // 记录当前堆内存分配的字节数
		b.start = time.Now()                // 记录测试启动时间
		b.timerOn = true                   // 标记计时标志
	}
}
```

StartTimer() 负责启动计时，并记录当前内存分配情况，不管是否有 “-benchmem” 参数，内存都会被统计，参数只决定是否要在结果中输出。

### B.StopTimer()

StopTimer() 负责停止计时，并累加相应的统计值。
```go
func (b *B) StopTimer() {
	if b.timerOn {
		b.duration += time.Since(b.start)                   // 累加测试耗时
		runtime.ReadMemStats(&memStats)                     // 读取当前堆内存分配信息
		b.netAllocs += memStats.Mallocs - b.startAllocs     // 累加堆内存分配的对象数
		b.netBytes += memStats.TotalAlloc - b.startBytes    // 累加堆内存分配的字节数
		b.timerOn = false                                  // 标记计时标志
	}
}
```

需要注意的是，StopTimer() 并不一定是测试结束，一个测试中有可能有多个统计阶段，所以其统计值是累加的。

### B.ResetTimer()

ResetTimer() 用于重置计时器，相应的也会把其他统计值也重置。

```go
func (b *B) ResetTimer() {
	if b.timerOn {
		runtime.ReadMemStats(&memStats)     // 读取当前堆内存分配信息
		b.startAllocs = memStats.Mallocs    // 记录当前堆内存分配的对象数
		b.startBytes = memStats.TotalAlloc  // 记录当前堆内存分配的字节数
		b.start = time.Now()                // 记录测试启动时间
	}
	b.duration = 0                          // 清空耗时
	b.netAllocs = 0                         // 清空内存分配对象数
	b.netBytes = 0                          // 清空内存分配字节数
}
```

ResetTimer() 比较常用，典型使用场景是一个测试中，初始化部分耗时较长，初始化后再开始计时。

### B.SetBytes(n int64)

```go
// SetBytes records the number of bytes processed in a single operation.
// If this is called, the benchmark will report ns/op and MB/s.
func (b *B) SetBytes(n int64) {
	b.bytes = n
}
```

这是一个比较含糊的函数，通过其函数说明很难明白其作用。

其实它是用来设置单次迭代处理的字节数，一旦设置了这个字节数，那么输出报告中将会呈现 “xxx MB/s” 的信息，用来表示待测函数处理字节的性能。待测函数每次处理多少字节数只有用户清楚，所以需要用户设置。

举个例子，待测函数每次执行处理 1M 数据，如果我们想看待测函数处理数据的性能，那么我们在测试中设置 SetBytes(1024 * 1024)，假如待测函数需要执行 1s 的话，那么结果中将会出现 “1 MB/s” （约等于）的信息。示例代码如下所示：

```go
func BenchmarkSetBytes(b *testing.B) {
    b.SetBytes(1024 * 1024)
    for i := 0; i < b.N; i++ {
        time.Sleep(1 * time.Second) // 模拟待测函数
    }
}
```

打印结果：

```bash
go test -bench SetBytes benchmark_test.go
```

```
BenchmarkSetBytes-4            1        1010392800 ns/op           1.04 MB/s
PASS
ok      command-line-arguments  1.412s
```

可以看到测试执行了一次，花费时间约 1S，数据处理能力约为 1MB/s。

### 报告内存信息

```go
func (b *B) ReportAllocs() {
	b.showAllocResult = true
}
```

ReportAllocs() 用于设置是否打印内存统计信息，与命令行参数 “-benchmem” 一致，但本方法只作用于单个测试函数。

## 性能测试是如何启动的

性能测试要经过多次迭代，每次迭代可能会有不同的 b.N 值，每次迭代执行测试函数一次，根据此次迭代的测试结果来分析要不要继续下一次迭代。

我们先看一下每次迭代时所用到的方法，runN():

```go
func (b *B) runN(n int) {
	b.N = n                       // 指定B.N
	b.ResetTimer()                // 清空统计数据
	b.StartTimer()                // 开始计时
	b.benchFunc(b)                // 执行测试
	b.StopTimer()                 // 停止计时
}
```

该方法指定 b.N 的值，执行一次测试函数。

与 T.Run() 类似，B.Run() 也用于启动一个子测试，实际上用户编写的任何一个测试都是使用 Run() 方法启动的，我们看下 B.Run() 的伪代码：

```go
func (b *B) Run(name string, f func(b *B)) bool {
	sub := &B{                          // 新建子测试数据结构
		common: common{
			signal:  make(chan bool),
			name:    name,
			parent:  &b.common,
		},
		benchFunc:  f,
	}
	if sub.run1() { // 先执行一次子测试，如果子测试不出错且子测试没有子测试的话继续执行sub.run()
		sub.run()       // run()里决定要执行多少次runN()
	}
	b.add(sub.result) // 累加统计结果到父测试中
	return !sub.failed
}
```

所有的测试都是先使用 run1() 方法执行一次测试，run1() 方法中实际上调用了 runN(1)，执行一次后再决定要不要继续迭代。

测试结果实际上以最后一次迭代的数据为准，当然，最后一次迭代往往意味着 b.N 更大，测试准确性相对更高。

## B.N是如何调整的？

B.launch() 方法里最终决定 B.N 的值。我们看下伪代码：

```go
func (b *B) launch() { // 此方法自动测算执行次数，但调用前必须调用run1以便自动计算次数
	d := b.benchTime
	for n := 1; !b.failed && b.duration < d && n < 1e9; { // 最少执行b.benchTime（默认为1s）时间，最多执行1e9次
		last := n
		n = int(d.Nanoseconds()) // 预测接下来要执行多少次，b.benchTime/每个操作耗时
		if nsop := b.nsPerOp(); nsop != 0 {
			n /= int(nsop)
		}
		n = max(min(n+n/5, 100*last), last+1) // 避免增长较快，先增长20%，至少增长1次
		n = roundUp(n) // 下次迭代次数向上取整到10的指数，方便阅读
		b.runN(n)
	}
}
```

不考虑程序出错，而且用户没有主动停止测试的场景下，每个性能测试至少要执行 b.benchTime 长的秒数，默认为 1s。先执行一遍的意义在于看用户代码执行一次要花费多长时间，如果时间较短，那么 b.N 值要足够大才可以测得更精确，如果时间较长，b.N 值相应的会减少，否则会影响测试效率。

最终的 b.N 会被定格在某个 10 的指数级，是为了方便阅读测试报告。

## 内存是如何统计的？

我们知道在测试开始时，会把当前内存值记入到 b.startAllocs 和 b.startBytes 中，测试结束时，会用最终内存值与开始时的内存值相减，得到净增加的内存值，并记入到 b.netAllocs 和 b.netBytes 中。

每个测试结束，会把结果保存到 BenchmarkResult 对象里，该对象里保存了输出报告所必需的统计信息：

```go
type BenchmarkResult struct {
	N         int           // 用户代码执行的次数
	T         time.Duration // 测试耗时
	Bytes     int64         // 用户代码每次处理的字节数，SetBytes()设置的值
	MemAllocs uint64        // 内存对象净增加值
	MemBytes  uint64        // 内存字节净增加值
}
```

其中 MemAllocs 和 MemBytes 分别对应 b.netAllocs 和 b.netBytes。

那么最终统计时只需要把净增加值除以 b.N 即可得到每次新增多少内存了。

每个操作内存对象新增值：

```go
func (r BenchmarkResult) AllocsPerOp() int64 {
	return int64(r.MemAllocs) / int64(r.N)
}
```

每个操作内存字节数新增值：

```go
func (r BenchmarkResult) AllocedBytesPerOp() int64 {
	if r.N <= 0 {
		return 0
	}
	return int64(r.MemBytes) / int64(r.N)
}
```

```go

```
