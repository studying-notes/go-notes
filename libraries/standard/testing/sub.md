---
date: 2020-11-15T20:06:59+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 子测试"  # 文章标题
url:  "posts/go/libraries/standard/testing/sub"  # 设置网页永久链接
tags: [ "Go", "sub" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

子测试提供一种在一个测试函数中执行多个测试的能力，比如原来有 TestA、TestB 和 TestC 三个测试函数，每个测试函数执行开始都需要做些相同的初始化工作，那么可以利用子测试将这三个测试合并到一个测试中，这样初始化工作只需要做一次。

## 示例

```go
package main

import (
	"testing"
)

func Add(x, y int) int {
	return x + y
}

func sub1(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub2(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub3(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func TestSub(t *testing.T) {
    // setup
	t.Run("A=1", sub1)
	t.Run("A=2", sub2)
	t.Run("B=1", sub3)
    // tear-down
}
```

本例中 `TestSub()` 通过 `t.Run()` 依次执行三个子测试。`t.Run()` 函数声明如下：

```go
func (t *T) Run(name string, f func(t *T)) bool
```

`name` 参数为子测试的名字，`f` 为子测试函数，本例中 `Run()` 一直阻塞到 `f` 执行结束后才返回，返回值为 f 的执行结果。

`Run()` 会启动新的协程来执行 `f`，并阻塞等待 `f` 执行结束才返回，除非 `f` 中使用 `t.Parallel()` 设置子测试为并发。

本例中 `TestSub()` 把三个子测试合并起来，可以共享 setup 和 tear-down 部分的代码。

我们在命令行下，使用 `-v` 参数执行测试：

```bash
go test main_test.go -v
```

```
=== RUN   TestSub
=== RUN   TestSub/A=1
=== RUN   TestSub/A=2
=== RUN   TestSub/B=1
--- PASS: TestSub (0.00s)
    --- PASS: TestSub/A=1 (0.00s)
    --- PASS: TestSub/A=2 (0.00s)
    --- PASS: TestSub/B=1 (0.00s)
PASS
ok      command-line-arguments  0.025s
```

从输出中可以看出，三个子测试都被执行到了，而且执行次序与调用次序一致。

## 子测试命名规则

通过上面的例子我们知道 `Run()` 方法第一个参数为子测试的名字，而实际上子测试的内部命名规则为："*<父测试名字>*/*<传递给Run的名字>*"。比如，传递给 `Run()` 的名字是 “A = 1”，那么子测试名字为 “TestSub/A = 1”。这个在上面的命令行输出中也可以看出。

## 过滤筛选

通过测试的名字，可以在执行中过滤掉一部分测试。

比如，只执行上例中 “A =* ” 的子测试，那么执行时使用 `-run Sub/A = ` 参数即可：

```bash
go test main_test.go -v -run Sub/A=
```

上例中，使用参数 `-run Sub/A = ` 则只会执行 `TestSub/A = 1` 和 `TestSub/A = 2` 两个子测试。

对于子性能测试则使用 `-bench` 参数来筛选，此处的筛选不是严格的正则匹配，而是包含匹配。比如，`-run A=`那么所有测试（含子测试）的名字中如果包含“A=”则会被选中执行。

## 子测试并发

前面提到的多个子测试共享 setup 和 teardown 有一个前提是子测试没有并发，如果子测试使用 `t.Parallel()` 指定并发，那么就没办法共享 teardown 了，因为执行顺序很可能是 setup-> 子测试 1->teardown-> 子测试 2...。

如果子测试可能并发，则可以把子测试通过 `Run()` 再嵌套一层，`Run()` 可以保证其下的所有子测试执行结束后再返回。

为便于说明，我们创建文件 `subparallel_test.go` 用于说明：

```go
package main

import (
    "testing"
    "time"
)

// 并发子测试，无实际测试工作，仅用于演示
func parallelTest1(t *testing.T) {
    t.Parallel()
    time.Sleep(3 * time.Second)
    // do some testing
}

// 并发子测试，无实际测试工作，仅用于演示
func parallelTest2(t *testing.T) {
    t.Parallel()
    time.Sleep(2 * time.Second)
    // do some testing
}

// 并发子测试，无实际测试工作，仅用于演示
func parallelTest3(t *testing.T) {
    t.Parallel()
    time.Sleep(1 * time.Second)
    // do some testing
}

// TestSubParallel 通过把多个子测试放到一个组中并发执行，同时多个子测试可以共享setup和tear-down
func TestSubParallel(t *testing.T) {
    // setup
    t.Logf("Setup")

    t.Run("group", func(t *testing.T) {
        t.Run("Test1", parallelTest1)
        t.Run("Test2", parallelTest2)
        t.Run("Test3", parallelTest3)
    })

    // tear down
    t.Logf("teardown")
}
```

上面三个子测试中分别 sleep 了 3s、2s、1s 用于观察并发执行顺序。通过 `Run()` 将多个子测试 “ 封装 ” 到一个组中，可以保证所有子测试全部执行结束后再执行 tear-down。

命令行下的输出如下：

```bash
go test subparallel_test.go -v -run SubParallel
```

```
=== RUN   TestSubParallel
=== RUN   TestSubParallel/group
=== RUN   TestSubParallel/group/Test1
=== RUN   TestSubParallel/group/Test2
=== RUN   TestSubParallel/group/Test3
--- PASS: TestSubParallel (3.01s)
        subparallel_test.go:25: Setup
    --- PASS: TestSubParallel/group (0.00s)
        --- PASS: TestSubParallel/group/Test3 (1.00s)
        --- PASS: TestSubParallel/group/Test2 (2.01s)
        --- PASS: TestSubParallel/group/Test1 (3.01s)
        subparallel_test.go:34: teardown
PASS
ok      command-line-arguments  3.353s
```

通过该输出可以看出：

1. 子测试是并发执行的（ Test1 最先被执行却最后结束）
2. tear-down 在所有子测试结束后才执行

## 小结

* 子测试适用于单元测试和性能测试；
* 子测试可以控制并发；
* 子测试提供一种类似 table-driven 风格的测试；
* 子测试可以共享 setup 和 tear-down。

```go

```
