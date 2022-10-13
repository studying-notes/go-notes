---
date: 2022-10-07T13:50:29+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "接口与程序设计模式"  # 文章标题
url:  "posts/go/docs/internal/interface/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 接口的用途

### 隐藏细节

接口可以对对象进行必要的抽象，外接设备只要满足相应标准（例如 USB 协议），就可与主设备对接，应用程序只要满足操作系统规定的系统调用方式，就可以使用操作系统提供的强大功能，而不必关注对方具体的实现细节。

### 控制系统复杂性

通过接口，我们能够以模块化的方式构建起复杂、庞大的系统。通过**将复杂的功能拆分成彼此独立的模块**，不仅有助于更好地并行开发系统、加速系统开发效率，也能在设计系统时以全局的视野看待整个系统。另外，模块拆分有助于快速排查、定位和解决问题。

### 权限控制

接口是系统与外界交流的唯一途径，例如 Go 语言对于垃圾回收暴露出的 GOGC 参数及 runtime.GC 方法。USB 接口有标准的接口协议，如果外界不满足这种协议，就无法与指定的系统进行交流。因此系统可以**通过接口控制接入方式和接入方行为，降低安全风险**。

## Go 语言中的接口

编程语言中的接口具有不同的表现形式，类似一种用于沟通的“共享边界”。

在面向对象的编程语言中，**接口指相互独立的两个对象之间的交流方式**。

Go 语言采用一种不寻常的方法进行面向对象编程，在 Go 语言中可以为任何自定义的类型添加方法，而不仅仅是类。没有任何形式的基于类型的继承，取而代之的是使用接口来实现扁平化、**面向组合**的设计模式。

**Go 语言的接口类型是延迟绑定**，可以实现类似虚函数的多态功能。

在 Go 语言中，接口是一种特殊的类型，**是其他类型可以实现的方法签名的集合**。**方法签名只包含方法名、输入参数和返回值**，下列为 `ReadCloser` 接口。

```go
// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadCloser interface {
	Reader
	Closer
}

// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Closer is the interface that wraps the basic Close method.
//
// The behavior of Close after the first call is undefined.
// Specific implementations may document their own behavior.
type Closer interface {
	Close() error
}
```

定义的接口可以是其他接口的组合。

一个接口包含越多的方法，其抽象性就越低，表达的行为就越具体。

Go 语言的设计者认为，对于传统拥有类型继承的面向对象语言，必须尽早设计其层次结构，一旦开始编写程序，早期决策就很难改变。这种方式导致了早期的过度设计，因为开发者试图预测程序所有可能的行为，增加了不必要的类型和抽象层。

Go 语言在设计之初，就鼓励开发者**使用组合而不是继承的方式来编写程序**。通常使用一种方法的接口来定义琐碎的行为，这些行为充当组件之间清晰、可理解的边界。Go 语言中的接口可以使程序自然、优雅、安全地增长，接口的更改仅影响实现接口的直接类型。

## 接口动态类型

一个接口类型的变量能够接收任何实现了此接口的用户自定义类型。一般将存储在接口中的类型称为接口的动态类型，而将接口本身的类型称为接口的静态类型。

## 接口类型断言

### 直接断言

在确信类型是正确的情况下可以直接断言：

```go
t := i.(T)
```

中 i 代表接口，T 代表实现此接口的动态类型。

这个表达式可以断言一个接口对象 `i` 里不是 `nil`，并且接口对象 `i` 存储的值的类型是 `T`，如果断言成功，就会返回值给 `t`，如果断言失败，就会触发 `panic`。

在编译时会保证类型 Type 一定是实现了接口 i 的类型，否则编译不会通过。

虽然 Go 语言在编译时已经防止了此类错误，但是仍然需要在运行时判断一次，这是由于在类型断言方法 `m=i.(Type)` 中，当 Type 实现了接口 i，而接口内部没有任何**动态类型**（其结构体中的 `_type` 字段代表接口存储的动态类型，此时为 nil）时，在运行时会直接 panic，因为 nil 无法调用任何方法。

### 判别断言

为了避免运行时触发 `panic`，类型转换还有第二种接口类型断言语法。

```go
t, ok:= i.(T)
```

断言成功就会返回其类型给 `t`，并且此时 `ok` 的值为 `true`，表示断言成功；断言失败，不会触发 `panic`，而是将 `ok` 的值设为 `false` ，表示断言失败，此时 `t` 为 `T` 的零值。

## 空接口

空接口是指没有定义任何方法的接口。因此任何类型都实现了空接口。

```go
type interface{}{}
```

空接口类型的变量可以存储任意类型的变量。

```go
var i interface{}
i = 42
i = "hello"
```

空接口类型的变量可以作为函数的参数。

```go
func describe(i interface{}) {
    fmt.Printf("Type = %T, value = %v\n", i, i)
}
```

Go 语言中提供了一种获取空接口中动态类型的方法。其语法是：

```go
i.(type)
```

其中，i 代表接口变量，type 是固定的关键字，不可与带方法接口的断言混淆。同时，此语法仅在 switch 语句中有效。

例如在 `fmt.Println` 源码中，使用 switch 语句嵌套这一语法可以获取空接口中的动态类型，并根据动态类型的不同进行不同的格式化输出。

`src/fmt/print.go`

```go
func (p *pp) printArg(arg any, verb rune) {
	// Some types can be done without reflection.
	switch f := arg.(type) {
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
		p.fmtFloat(f, 64, verb)
	case complex64:
		p.fmtComplex(complex128(f), 64, verb)
	case complex128:
		p.fmtComplex(f, 128, verb)
    }
}
```

## 接口的比较性

接口类型的变量可以使用 == 或 != 进行比较。

```go
var i interface{}
var j interface{}
fmt.Println(i == j) // true
```

接口的比较规则如下：

- 动态类型字段 `_type` 为 nil 的接口变量总是相等的。
- 如果只有 1 个接口为 nil，那么比较结果总是 false。
- 如果两个接口不为 nil 且接口变量具有相同的**动态类型和动态类型值**，那么两个接口是相同的。
- 如果接口存储的**动态类型值**是不可比较的，那么在运行时会报错。

## 接口转换

Go 语言中，对于基础类型（非接口类型）不支持隐式的转换，我们无法将一个 int 类型的值直接赋值给 int64 类型的变量，也无法将 int 类型的值赋值给底层是 int 类型的新定义命名类型的变量。Go 语言对基础类型的类型一致性要求可谓是非常的严格，但是 Go 语言对于接口类型的转换则非常灵活。对象和接口之间的转换、接口和接口之间的转换都可能是隐式的转换。可以看下面的例子：

```go
var (
	a io.ReadCloser = (*os.File)(f) // 隐式转换，*os.File 满足 io.ReadCloser接 口
	b io.Reader     = a             // 隐式转换，io.ReadCloser 满足 io.Reader 接口
	c io.Closer     = a             // 隐式转换，io.ReadCloser 满足 io.Closer 接口
	d io.Reader     = c.(io.Reader) // 显式转换，io.Closer 不满足 io.Reader 接口
)
```

有时候对象和接口之间太灵活了，需要人为地限制这种无意之间的适配。常见的做法是定义一个特殊方法来区分接口。例如 `runtime` 包中的 `Error` 接口就定义了一个特有的 `RuntimeError()` 方法，用于避免其他类型无意中适配了该接口：

`src/runtime/error.go`

```go
// The Error interface identifies a run time error.
type Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	RuntimeError()
}
```

在 `Protobuf` 中，`Message` 接口也采用了类似的方法，也定义了一个特有的 `ProtoMessage`，用于避免其他类型无意中适配了该接口：

```go
type proto.Message interface {
    Reset()
    String() string
    ProtoMessage()
}
```

不过这种做法只是“君子协定”，如果有人故意伪造一个 `proto.Message` 接口也是很容易的。再严格一点的做法是**给接口定义一个私有方法**。只有满足了这个私有方法的对象才可能满足这个接口，而**私有方法的名字是包含包的绝对路径名的**，因此**只有在包内部实现这个私有方法才能满足这个接口**。测试包中的 `testing.TB` 接口就是采用类似的技术：

`src/testing/testing.go`

```go
// TB is the interface common to T, B, and F.
type TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    ...

    // A private method to prevent users implementing the
    // interface and so future additions to it will not
    // violate Go 1 compatibility.
    private()
}
```

不过这种通过私有方法禁止外部对象实现接口的做法也是有代价的：首先是**这个接口只能在包内部使用**，外部包在正常情况下是无法直接创建满足该接口对象的；其次，这种防护措施也不是绝对的，用户**依然可以绕过这种保护机制**。

通过在结构体中嵌入匿名类型成员，可以继承匿名类型的方法。其实这个被**嵌入的匿名成员不一定是普通类型，也可以是接口类型**。我们可以通过嵌入匿名的 `testing.TB` 接口来伪造私有方法，因为**接口方法是延迟绑定**，所以**编译时私有方法是否真的存在并不重要**。

```go
package main

import (
    "fmt"
    "testing"
)

type TB struct {
    testing.TB
}

func (p *TB) Fatal(args ...interface{}) {
    fmt.Println("TB.Fatal disabled!")
}

func main() {
    var tb testing.TB = new(TB)
    tb.Fatal("Hello, playground")
}
```

我们在自己的 TB 结构体类型中重新实现了 `Fatal()` 方法，然后通过将对象隐式转换为 `testing.TB` 接口类型（因为内嵌了匿名的 `testing.TB` 对象，所以是满足 `testing.TB` 接口的），再通过 `testing.TB` 接口来调用自己的 `Fatal()` 方法。

这种**通过嵌入匿名接口或嵌入匿名指针对象来实现继承**的做法其实是一种**纯虚继承**，**继承的只是接口指定的规范**，**真正的实现在运行的时候才被注入**。

## 接口底层原理

[接口底层原理](underlying_principle.md)

```go

```
