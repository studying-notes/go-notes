# Go 性能测试

- [Go 性能测试](#go-性能测试)
	- [基准测试](#基准测试)
	- [编写基准测试](#编写基准测试)
	- [性能对比](#性能对比)

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
