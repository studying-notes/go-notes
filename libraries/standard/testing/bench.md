---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 基准测试"  # 文章标题
url:  "posts/go/libraries/standard/testing/bench"  # 设置网页永久链接
tags: [ "Go", "bench" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

基准测试是通过测试 CPU 和内存的效率问题，来评估被测试代码的性能，进而找到更好的解决方案。

## 编写测试

```go
package main

import (
	"fmt"
	"testing"
)

func BenchmarkSprint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprint(i)
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

```bash
go test -bench=".*"
```

```
goos: windows
goarch: amd64
pkg: example
cpu: Intel(R) Core(TM) i5-8500 CPU @ 3.00GHz
BenchmarkSprint-6       14685240                81.15 ns/op
PASS
ok      example 1.308s
```

- `BenchmarkSprint-8` 中的 8 表示运行时 `GOMAXPROCS` 的值，即 CPU 的核数；
- `14685240` 表示运行 `for` 循环的次数，即调用被测试代码的次数；
- `81.15 ns/op` 表示执行一次花费的纳秒时间。

## 性能对比

```go
package main

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkSprint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprint(i)
	}
}

func BenchmarkItoa(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}
```

```bash
go test -bench=".*"
```

```
goos: windows
goarch: amd64
pkg: example
cpu: Intel(R) Core(TM) i5-8500 CPU @ 3.00GHz
BenchmarkSprint-6       14691478                79.50 ns/op
BenchmarkItoa-6         42859285                26.58 ns/op
PASS
ok      example 2.451s
```

`-benchmem` 参数可以显示每次操作分配内存的次数，以及每次操作分配的字节数。

```bash
go test -bench=".*" -benchmem
```

```
goos: windows
goarch: amd64
pkg: example
cpu: Intel(R) Core(TM) i5-8500 CPU @ 3.00GHz
BenchmarkSprint-6       14935942                80.86 ns/op           16 B/op                     1 allocs/op
BenchmarkItoa-6         48002112                26.84 ns/op            7 B/op                     0 allocs/op
PASS
ok      example 2.647s
```

```go

```
