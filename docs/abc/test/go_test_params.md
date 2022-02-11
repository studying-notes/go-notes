---
date: 2020-11-16T21:01:20+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "深入测试标准库之 go test 参数"  # 文章标题
url:  "posts/go/abc/memory/test/go_test_params"  # 设置网页永久链接
tags: [ "go", "test" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 前言
go test有非常丰富的参数，一些参数用于控制测试的编译，另一些参数控制测试的执行。

有关测试覆盖率、vet和pprof相关的参数先略过，我们在讨论相关内容时再详细介绍。

## 控制编译的参数
### -args
指示go test把-args后面的参数带到测试中去。具体的测试函数会根据此参数来控制测试流程。

-args后面可以附带多个参数，所有参数都将以字符串形式传入，每个参数作为一个string，并存放到字符串切片中。

```go
// TestArgs 用于演示如何解析-args参数
func TestArgs(t *testing.T) {
    if !flag.Parsed() {
        flag.Parse()
    }

    argList := flag.Args() // flag.Args() 返回 -args 后面的所有参数，以切片表示，每个元素代表一个参数
    for _, arg := range argList {
        if arg == "cloud" {
            t.Log("Running in cloud.")
        } else {
            t.Log("Running in other mode.")
        }
    }
}
```
执行测试时带入参数：
```
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -run TestArgs -v -args "cloud"
TestMain setup.
=== RUN   TestArgs
--- PASS: TestArgs (0.00s)
    unit_test.go:28: Running in cloud.
PASS
TestMain tear-down.
ok      gotest  0.353s

```
通过参数-args指定传递给测试的参数。

### -json
-json 参数用于指示go test将结果输出转换成json格式，以方便自动化测试解析使用。

示例如下：
```
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -run TestAdd -json
{"Time":"2019-02-28T15:46:50.3756322+08:00","Action":"output","Package":"gotest","Output":"TestMain setup.\n"}
{"Time":"2019-02-28T15:46:50.4228258+08:00","Action":"run","Package":"gotest","Test":"TestAdd"}
{"Time":"2019-02-28T15:46:50.423809+08:00","Action":"output","Package":"gotest","Test":"TestAdd","Output":"=== RUN   TestAdd\n"}
{"Time":"2019-02-28T15:46:50.423809+08:00","Action":"output","Package":"gotest","Test":"TestAdd","Output":"--- PASS: TestAdd (0.00s)\n"}
{"Time":"2019-02-28T15:46:50.423809+08:00","Action":"pass","Package":"gotest","Test":"TestAdd","Elapsed":0}
{"Time":"2019-02-28T15:46:50.4247922+08:00","Action":"output","Package":"gotest","Output":"PASS\n"}
{"Time":"2019-02-28T15:46:50.4247922+08:00","Action":"output","Package":"gotest","Output":"TestMain tear-down.\n"}
{"Time":"2019-02-28T15:46:50.4257754+08:00","Action":"output","Package":"gotest","Output":"ok  \tgotest\t0.465s\n"}
{"Time":"2019-02-28T15:46:50.4257754+08:00","Action":"pass","Package":"gotest","Elapsed":0.465}

```

### -o <file>
-o 参数指定生成的二进制可执行程序，并执行测试，测试结束不会删除该程序。

没有此参数时，go test生成的二进制可执行程序存放到临时目录，执行结束便删除。

示例如下：
```
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -run TestAdd -o TestAdd
TestMain setup.
PASS
TestMain tear-down.
ok      gotest  0.439s
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>TestAdd
TestMain setup.
PASS
TestMain tear-down.
```

本例中，使用-o 参数指定生成二进制文件"TestAdd"并存放到当前目录，测试执行结束后，仍然可以直接执行该二进制程序。

## 控制测试的参数
### -bench regexp
go test默认不执行性能测试，使用-bench参数才可以运行，而且只运行性能测试函数。

其中正则表达式用于筛选所要执行的性能测试。如果要执行所有的性能测试，使用参数"-bench ."或"-bench=."。

此处的正则表达式不是严格意义上的正则，而是种包含关系。

比如有如下三个性能测试：
* func BenchmarkMakeSliceWithoutAlloc(b *testing.B)
* func BenchmarkMakeSliceWithPreAlloc(b *testing.B)
* func BenchmarkSetBytes(b *testing.B)

使用参数“-bench=Slice”，那么前两个测试因为都包含"Slice"，所以都会被执行，第三个测试则不会执行。

对于包含子测试的场景下，匹配是按层匹配的。举一个包含子测试的例子：
```go
func BenchmarkSub(b *testing.B) {
    b.Run("A=1", benchSub1)
    b.Run("A=2", benchSub2)
    b.Run("B=1", benchSub3)
}
```
测试函数命名规则中，子测试的名字需要以父测试名字作为前缀并以"/"连接，上面的例子实际上是包含4个测试：
* Sub
* Sub/A=1
* Sub/A=2
* Sub/B=1

如果想执行三个子测试，那么使用参数“-bench Sub”。如果只想执行“Sub/A=1”，则使用参数"-bench Sub/A=1"。如果想执行"Sub/A=1"和“Sub/A=2”，则使用参数"-bench Sub/A="。

### -benchtime <t>s
-benchtime指定每个性能测试的执行时间，如果不指定，则使用默认时间1s。

例如，执定每个性能测试执行2s，则参数为："go test -bench Sub/A=1 -benchtime 2s"。

### -cpu 1,2,4
-cpu 参数提供一个CPU个数的列表，提供此列表后，那么测试将按照这个列表指定的CPU数设置GOMAXPROCS并分别测试。

比如“-cpu 1,2”，那么每个测试将执行两次，一次是用1个CPU执行，一次是用2个CPU执行。
例如，使用命令"go test -bench Sub/A=1 -cpu 1,2,3,4" 执行测试：
```
BenchmarkSub/A=1                    1000           1256835 ns/op
BenchmarkSub/A=1-2                  2000            912109 ns/op
BenchmarkSub/A=1-3                  2000            888671 ns/op
BenchmarkSub/A=1-4                  2000            894531 ns/op
```
测试结果中测试名后面的-2、-3、-4分别代表执行时GOMAXPROCS的数值。 如果GOMAXPROCS为1，则不显示。

### -count n
-count指定每个测试执行的次数，默认执行一次。

例如，指定测试执行2次：
```go
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -bench Sub/A=1 -count 2
TestMain setup.
goos: windows
goarch: amd64
pkg: gotest
BenchmarkSub/A=1-4                  2000            917968 ns/op
BenchmarkSub/A=1-4                  2000            882812 ns/op
PASS
TestMain tear-down.
ok      gotest  10.236s

```
可以看到结果中也将呈现两次的测试结果。

如果使用-count指定执行次数的同时还指定了-cpu列表，那么测试将在每种CPU数量下执行count指定的次数。

注意，示例测试不关心-count和-cpu参数，它总是执行一次。

### -failfast
默认情况下，go test将会执行所有匹配到的测试，并最后打印测试结果，无论成功或失败。

-failfast指定如果有测试出现失败，则立即停止测试。这在有大量的测试需要执行时，能够更快的发现问题。

### -list regexp
-list 只是列出匹配成功的测试函数，并不真正执行。而且，不会列出子函数。

例如，使用参数"-list Sub"则只会列出包含子测试的三个测试，但不会列出子测试：
```go
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -list Sub
TestMain setup.
TestSubParallel
TestSub
BenchmarkSub
TestMain tear-down.
ok      gotest  0.396s

```

### -parallel n
指定测试的最大并发数。

当测试使用t.Parallel()方法将测试转为并发时，将受到最大并发数的限制，默认情况下最多有GOMAXPROCS个测试并发，其他的测试只能阻塞等待。

### -run regexp
根据正则表达式执行单元测试和示例测试。正则匹配规则与-bench 类似。

### -timeout d
默认情况下，测试执行超过10分钟就会超时而退出。

例时，我们把超时时间设置为1s，由本来需要3s的测试就会因超时而退出：
```
E:\OpenSource\GitHub\RainbowMango\GoExpertProgrammingSourceCode\GoExpert\src\gotest>go test -timeout=1s
TestMain setup.
panic: test timed out after 1s

```
设置超时可以按秒、按分和按时：
* 按秒设置：-timeout xs或-timeout=xs
* 按分设置：-timeout xm或-timeout=xm
* 按时设置：-timeout xh或-timeout=xh

### -v
默认情况下，测试结果只打印简单的测试结果，-v 参数可以打印详细的日志。

性能测试下，总是打印日志，因为日志有时会影响性能结果。

### -benchmem
默认情况下，性能测试结果只打印运行次数、每个操作耗时。使用-benchmem则可以打印每个操作分配的字节数、每个操作分配的对象数。

```
// 没有使用-benchmem
BenchmarkMakeSliceWithoutAlloc-4            2000            971191 ns/op

// 使用-benchmem
BenchmarkMakeSliceWithoutAlloc-4            2000            914550 ns/op         4654335 B/op         30 allocs/op

```
此处，每个操作的含义是放到循环中的操作，如下示例所示：
```go
func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gotest.MakeSliceWithoutAlloc() // 一次操作
    }
}
```
