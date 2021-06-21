---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 性能测试"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/test/benchmark"  # 设置网页永久链接
tags: [ "go", "benchmark" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 基准测试

基准测试是通过测试 CPU 和内存的效率问题，来评估被测试代码的性能，进而找到更好的解决方案。

## 编写基准测试

```go
import (
	"fmt"
	"testing"
)

func BenchmarkSprint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", i)
	}
}
```

1. 基准测试的代码文件必须以 `_test.go` 结尾；
2. 基准测试的函数必须以 `Benchmark` 开头，必须是可导出的；
3. 基准测试函数必须接受一个指向 `Benchmark` 类型的指针作为唯一参数；
4. 基准测试函数不能有返回值；
5. `b.ResetTimer` 是重置计时器，这样可以避免 `for` 循环之前的**初始化代码**的干扰；
6. 最后的 `for` 循环很重要，被测试的代码要放到循环里；
7. `b.N` 是基准测试框架提供的，表示循环的次数，因为需要反复调用测试的代码，才可以评估性能。

```shell
$ go test -bench=".*"
goos: windows
goarch: amd64
pkg: github/fujiawei-dev/go-notes/test/benchmark
BenchmarkSprint-8       11351102               103 ns/op
PASS
ok      github/fujiawei-dev/go-notes/test/benchmark     1.338s
```

- `BenchmarkSprint-8` 中的 8 表示运行时 `GOMAXPROCS` 的值，即 CPU 的核数；
- `11351102` 表示运行 `for` 循环的次数，即调用被测试代码的次数；
- `103 ns/op` 表示执行一次花费 117 纳秒。

## 性能对比

```shell
$ go test -bench=".*"
goos: windows
goarch: amd64
pkg: github/fujiawei-dev/go-notes/test/benchmark
BenchmarkSprint-8       15825914                75.3 ns/op
BenchmarkFormat-8       377147500                3.16 ns/op
BenchmarkItoa-8         384372810                3.09 ns/op
PASS
ok      github/fujiawei-dev/go-notes/test/benchmark     4.349s
```

`-benchmem` 参数可以显示每次操作分配内存的次数，以及每次操作分配的字节数。

```shell
$ go test -bench=".*" -benchmem
goos: windows
goarch: amd64
pkg: github/fujiawei-dev/go-notes/test/benchmark
BenchmarkSprint-8       15229449                76.4 ns/op             2 B/op          1 allocs/op   
BenchmarkFormat-8       349739210                3.13 ns/op            0 B/op          0 allocs/op   
BenchmarkItoa-8         383148979                3.08 ns/op            0 B/op          0 allocs/op   
PASS
ok      github/fujiawei-dev/go-notes/test/benchmark     4.247s
```

## 切片构建方法性能对比

```go
package gotest

// MakeSliceWithPreAlloc 不预分配
func MakeSliceWithoutAlloc() []int {
    var newSlice []int

    for i := 0; i < 100000; i++ {
        newSlice = append(newSlice, i)
    }

    return newSlice
}

// MakeSliceWithPreAlloc 通过预分配Slice的存储空间构造
func MakeSliceWithPreAlloc() []int {
    var newSlice []int

    newSlice = make([]int, 0, 100000)
    for i := 0; i < 100000; i++ {
        newSlice = append(newSlice, i)
    }

    return newSlice
}
```

两个方法都会构造一个容量为 100000 的切片，所不同的是 `MakeSliceWithPreAlloc()` 会预先分配内存，而 `MakeSliceWithoutAlloc()` 不预先分配内存，二者理论上存在性能差异。

```go
package gotest_test

import (
    "testing"
    "gotest"
)

func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gotest.MakeSliceWithoutAlloc()
    }
}

func BenchmarkMakeSliceWithPreAlloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gotest.MakeSliceWithPreAlloc()
    }
}
```

```shell
go test -bench=.
```
