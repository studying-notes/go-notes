---
date: 2022-03-19T08:43:59+08:00
author: "Rustle Karl"

title: "go1.18 新特性详解"
url:  "posts/go/quickstart/release/go1.18"  # 永久链接
tags: [ "go", "README" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

toc: true  # 目录
draft: false  # 草稿
---

## Module 工作区模式

### 缘起

假设本地有两个项目，分别是：project-json 和 project-example。

```shell
mkdir project-json
mkdir project-example
```

```shell
cd project-json
go mod init github.com/fujiawei-dev/project-json
code json.go
```

```go
package pjson

func Json(){
	println("json: for example")
}
```

```shell
cd ../project-example
go mod init github.com/fujiawei-dev/project-example
code main.go
```

```go
package main

import (
    "github.com/fujiawei-dev/project-json"
)

func main() {
    pjson.Json()
}
```

这时候，如果我们运行 `go mod tidy`，肯定会报错，因为我们的 project-json 包根本没有提交到 github 上，肯定找不到。

我们当然可以提交 project-json 到 github，但我们没修改一次 project-json，就需要提交，否则 project-example 中就没法使用上最新的。

针对这种情况，在 go1.18 以前是建议通过 replace 来解决，即在 example 中的 go.mod 增加如下 replace：

```shell
github.com/fujiawei-dev/project-example

go 1.17

require github.com/fujiawei-dev/project-json v1.0.0

replace github.com/fujiawei-dev/project-json => ../project-json
```

当都开发完成时，我们需要手动删除 replace，并执行 `go mod tidy` 后提交，否则别人使用就报错了。

这还是挺不方便的，如果本地有多个 module，每一个都得这么处理。

### 工作区模式

难道还没？

```go

```

```shell

```


## 字符串 Clone

Clone 返回 s 的新副本。它保证将 s 复制到一个新分配的副本中，当只保留一个很大的字符串中的一个小子字符串时，这一点很重要。使用克隆可以帮助这些程序使用更少的内存。当然，由于使用克隆制作拷贝，过度使用克隆会使程序使用更多内存。通常，只有在分析表明需要克隆时，才谨慎使用克隆。对于长度为零的字符串，将返回字符串 ""，不进行内存分配。

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s := "abcdefghijklmn"
	sSlice := s[:4]

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sSliceHeader := (*reflect.StringHeader)(unsafe.Pointer(&sSlice))
	
	fmt.Println(sHeader.Len == sSliceHeader.Len)
	fmt.Println(sHeader.Data == sSliceHeader.Data)

	// Output:
	// false
	// true
}
```

在上面示例场景中，如果 s 很大，而之后我们只需要使用它的某个短子串，这会导致内存的浪费，因为子串和原字符串的 Data 部分指向相同的内存，因此整个字符串并不会被 GC 回收。

strings.Clone 函数就是为了解决这个问题的：

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	s := "abcdefghijklmn"
	// sSlice := s[:4]
	sSlice := strings.Clone(s[:4])

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sSliceHeader := (*reflect.StringHeader)(unsafe.Pointer(&sSlice))

	fmt.Println(sHeader.Len == sSliceHeader.Len)
	fmt.Println(sHeader.Data == sSliceHeader.Data)

	// Output:
	// false
	// false
}
```

通过克隆得到 s2，从最后输出结果看，Data 已经不同了，原始的长字符串就可以被垃圾回收了。

### 内部实现

```go
func Clone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
    // 实现 []byte 到 string 的零内存拷贝转换。
	return *(*string)(unsafe.Pointer(&b))
}
```

## strings.Title 被废弃

strings.Title 会将每个单词的首字母变成大写字母。strings 中还有一个函数：ToTitle，它的作用和 ToUpper 类似，所有字符全部变成大写，而不只是首字母。不过 ToTitle 和 ToUpper 的区别特别微小，Stackoverflow 上有相关讨论，它们的区别是 Unicode 规定的区别。

那 strings.Title 为什么废弃呢？strings.Title 的规则是使用单词边界，不能正确处理 Unicode 标点。

```go
fmt.Println(strings.Title("here comes o'brian"))
```

期望输出：`Here Comes O'brian`，但 strings.Title 的结果是：`Here Comes O'Brian`。

### 继任者

在 strings.Title 中提到，可以使用 golang.org/x/text/cases 代替 strings.Title，具体来说就是 cases.Title。

该包提供了通用和特定于语言的 case map，其中有一个 Title 函数，签名如下：

```go
func Title(t language.Tag, opts ...Option) Caser
```

第一个参数是 language.Tag 类型，表示 BCP 47 种语言标记。它用于指定特定语言或区域设置的实例。所有语言标记值都保证格式良好。

第二个参数是不定参数，类型是 Option，这是一个函数类型：

```go
type Option func(o options) options
```

它被用来修改 Caser 的行为，cases 包可以找到相关 Option 的实例。

cases.Title 的返回类型是 Caser，这是一个结构体，这里我们只关心它的 String 方法，它接收一个字符串，并返回一个经过 Caser 处理过后的字符串。

所以，针对上文 strings.Title 的场景，可以改为 cases.Title 实现。

```go
caser := cases.Title(language.English)
caser.String("here comes o'brian")
```

## 新增 strings.Cut

### strings.Index 系列函数

常用与字符串的切割。

strings 包中，Index 相关函数有好几个：

```go
func Index(s, substr string) int
func IndexAny(s, chars string) int
func IndexByte(s string, c byte) int
func IndexFunc(s string, f func(rune) bool) int
func IndexRune(s string, r rune) int
func LastIndex(s, substr string) int
func LastIndexAny(s, chars string) int
func LastIndexByte(s string, c byte) int
func LastIndexFunc(s string, f func(rune) bool) int
```

Go 官方统计了 Go 源码中使用相关函数的代码：

- 311 Index calls outside examples and testdata.
- 20 should have been Contains
- 2 should have been 1 call to IndexAny
- 2 should have been 1 call to ContainsAny
- 1 should have been TrimPrefix
- 1 should have been HasSuffix

相关需求是这么多，而 Index 显然不是处理类似需求最好的方式。

于是 Russ Cox 提议，在 strings 包中新增一个函数 Cut，专门处理类似的常见。

### 新增的 Cut 函数

Cut 函数的签名如下：

```go
func Cut(s, sep string) (before, after string, found bool)
```

将字符串 s 在第一个 sep 处切割为两部分，分别存在 before 和 after 中。如果 s 中没有 sep，返回 s,"",false。

```go
addr := "192.168.1.1:8080"
ip, port, ok := strings.Cut(addr, ":")
```

## 泛型

### 类型约束

```go
type Addable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string
}
```

```go
func add[T int | float64](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(add(1, 2))
	fmt.Println(add(1.2, 2.3))
}
```

在 Go 语言中，基于某类型定义新类型，有时可能希望泛型约束是某类型的所有衍生类型。看一个具体例子：

```go
package main

import (
	"fmt"
)

func add[T ~string](x, y T) T {
	return x + y
}

type MyString string

func main() {
	var x string = "ab"
	var y MyString = "cd"
	fmt.Println(add(x, x))
	fmt.Println(add(y, y))
}

// Output:
// abab
// cdcd
```

约束 ~string  表示支持 string 类型以及底层是 string 类型的类型，因此 MyString 类型值也可以传递给 add。

### 约束形式的多样性

```go
// 没有任何约束
func add[T any](x, y T) T
// 约束 Addble (需要单独定义)
func add[T Addble](x, y T) T
// 约束允许 int 或 float64 类型
func add[T int|float64](x, y T) T
// 约束允许底层类型是 string 的类型（包括 string 类型）
func add[T ~string](x, y T) T
```

```go
func MakeChan[T chan bool | chan int](c T) {
	_ = make(T) // 错误
	_ = new(T)  // 正确
	_ = len(c)  // 正确
}

// 以下代码无法编译：
// cannot range over c (variable of type T constrained by []string|map[int]string) (T has no structural type)
func ForEach[T []string | map[int]string](c T, f func(int, string)) {
	for i, v := range c {
		f(i, v)
	}
}
```

## 新 IP 包

net/netip

https://polarisxu.studygolang.com/posts/go/dynamic/go1.18-ip/

## 原生支持 Fuzzing 测试

为什么Go要支持Fuzzing？

无论是在过去的单机时代，还是在今天的云计算时代，亦或是已经出现苗头的人工智能时代，安全都是程序员在构建软件过程中必不可少且日益重要的考量因素。

同时，安全也是 Go 语言设计者们在语言设计伊始就为 Go 设定的一个重要目标。在语言层面，Go 提供了很多“安全保障”特性，比如：

- 显式类型转换，不允许隐式转换；
- 针对数组、切片的下标越界访问的检查；
- 不安全代码隔离到 unsafe 包中，并提供安全使用 unsafe 包的几条 rules ；
- go module 构建模式内置包 hash 校验，放置恶意包对程序的攻击；
- 雇佣安全专家，提供高质量且及时更新的 crypto 包，尽量防止使用第三方加解密包带来的不安全性；
- 支持纯静态编译，避免动态连接到恶意动态库；
- 原生提供方便测试的工具链，并支持测试覆盖率统计。

进入云原生时代后，Go 语言成为了云原生基础设施与云原生服务的头部语言，由 Go 语言建造的基础设施、中间件以及应用服务支撑着这个世界的很多重要系统的运行，这些系统对安全性的要求不言而喻。尤其是在“输入”方面，这些系统都会被要求：无论用户使用什么本文数据、二进制数据，无论用户如何构造协议包、无论用户提供的文件包含何种内容，系统都应该是安全健壮的，不会因为用户的故意攻击而出现处理异常、被操控，甚至崩溃。

这就需要我们的系统面对任何输入的情况下都能正确处理，但传统的代码 review、静态分析、人工测试和自动化的单元测试无法穷尽所有输入组合，尤其是难于模拟一些无效的、意料之外的、随机的、边缘的输入数据。

于是，Go 社区的一些安全方面的专家就尝试将业界在解决这方面问题的优秀实践引入 Go，Fuzzing 技术就是其中最重要的一个。

### 什么是 Fuzzing ？

Fuzzing，又叫 fuzz testing，中文叫做模糊测试或随机测试。其本质上是一种自动化测试技术，更具体一点，它是一种基于随机输入的自动化测试技术，常被用于发现处理用户输入的代码中存在的 bug 和问题。

在具体实现上，Fuzzing 不需要像单元测试那样使用预先定义好的数据集作为程序输入，而是会通过数据构造引擎自行构造或基于开发人员提供的初始数据构造一些随机数据，并作为输入提供给我们的程序，然后监测程序是否出现 panic、断言失败、无限循环等。这些构造出来的随机数据被称为语料 (corpus)。另外 Fuzz testing 不是一次性执行的测试，如果不限制执行次数和执行时间，Fuzz testing 会一直执行下去，因此它也是一种持续测试的技术。

Fuzzing 是对其他形式的测试、代码审查和静态分析的补充，它通过生成一个有趣的输入语料库，而这些输入几乎不可能用手去想去写出来，因此极易被传统类型的测试所遗漏。Fuzzing 可以帮助开发人员发现难以发现的稳定性、逻辑性甚至是安全性方面的错误，特别是当被测系统变得更加复杂时。

更为关键的，也是 Fuzzing 技术能够流行开来的原因是构建一个 Fuzzing 测试足够简单。有了上述 Fuzzing 相关开源工具和开源库的支持，Fuzzing test 不再是需要专业知识才能成功使用的技术，并且现代 Fuzzing 引擎可以更有策略的、更快的、更有效地找到有用的输入语料。
