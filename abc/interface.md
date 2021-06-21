---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 语言接口"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/abc/interface"  # 设置网页永久链接
tags: [ "go", "interface" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 接口指针问题

定义时是接口，导致无法动态修改。

```go
package main

import "fmt"

type One interface {
	Value() string
}

type one struct{}

func (o one) Value() string {
	return "one"
}

type two struct{}

func (t two) Value() string {
	return "two"
}

func main() {
	var oo One

	oo = one{}
	xx := &oo
	oo = two{}

	fmt.Println((*xx).Value())
}
```

```go
package main

import "fmt"

type One interface {
Value() string
}

type one struct{ v string }

func (o one) Value() string {
return o.v
}

type two struct{}

func (t two) Value() string {
return "two"
}

func main() {
var oo one
oo = one{v: "1"}
xx := &oo
oo = one{v: "2"}

	fmt.Println(xx.Value())
}
```


一般静态编程语言都有着严格的类型系统，这使得编译器可以深入检查程序员有没有作出什么出格的举动。但是，过于严格的类型系统却会使得编程太过烦琐，让程序员把时间都浪费在了和编译器的斗争中。

Go 语言试图让程序员能在安全和灵活的编程之间取得一个平衡。它在提供严格的类型检查的同时，通过接口类型实现了对鸭子类型的支持，使得安全动态的编程变得相对容易。

Go 的接口类型是对其他类型行为的抽象和概括，因为接口类型不会和特定的实现细节绑定在一起，通过这种抽象的方式我们可以让对象更加灵活和更具有适应能力。很多面向对象的语言都有相似的接口概念，但 Go 语言中接口类型的独特之处在于它是满足隐式实现的鸭子类型。

所谓鸭子类型说的是：只要走起路来像鸭子、叫起来也像鸭子，那么就可以把它当作鸭子。Go 语言中的面向对象就是如此，如果一个对象只要看起来像是某种接口类型的实现，那么它就可以作为该接口类型使用。这种设计可以让你创建一个新的接口类型满足已经存在的具体类型却不用去破坏这些类型原有的定义。当使用的类型来自不受我们控制的包时这种设计尤其灵活有用。

**Go 语言的接口类型是延迟绑定**，可以实现类似虚函数的多态功能。

接口在 Go 语言中无处不在，在 “Hello, World” 的例子中，`fmt.Printf()` 函数的设计就是完全基于接口的，它的真正功能由 `fmt.Fprintf()` 函数完成。用于表示错误的 `error` 类型更是内置的接口类型。在 C 语言中，`printf` 只能将几种有限的基础数据类型打印到文件对象中。但是 Go 语言由于灵活的接口特性，`fmt.Fprintf` 可以向任何自定义的输出流对象打印，可以打印到文件或标准输出，也可以打印到网络，甚至可以打印到一个压缩文件；同时，打印的数据也不仅局限于语言内置的基础类型，任意隐式满足 `fmt.Stringer` 接口的对象都可以打印，不满足 `fmt.Stringer` 接口的依然可以通过反射的技术打印。`fmt.Fprintf()` 函数的签名如下：

```go
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
```

其中 `io.Writer` 是用于输出的接口，`error` 是内置的错误接口，它们的定义如下：

```go
type io.Writer interface {
    Write(p []byte) (n int, err error)
}

type error interface {
    Error() string
}
```

我们可以通过定制自己的输出对象，将每个字符转换为大写字符后输出：

```go
type UpperWriter struct {
    io.Writer
}

func (p *UpperWriter) Write(data []byte) (n int, err error) {
    return p.Writer.Write(bytes.ToUpper(data))
}

func main() {
    fmt.Fprintln(&UpperWriter{os.Stdout}, "hello, world")
}
```

当然，我们也可以定义自己的打印格式来实现将每个字符转换为大写字符后输出的效果。对于每个要打印的对象，如果满足了 `fmt.Stringer` 接口，则默认使用对象的 `String()` 方法返回的结果打印：

```go
type UpperString string

func (s UpperString) String() string {
    return strings.ToUpper(string(s))
}

type fmt.Stringer interface {
    String() string
}

func main() {
    fmt.Fprintln(os.Stdout, UpperString("hello, world"))
}
```

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

```go
type runtime.Error interface {
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

```go
type testing.TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    ...

    // A private method to prevent users implementing the
    // interface and so future additions to it will not
    // violate Go 1 compatibility.
    private()
}
```

不过这种通过私有方法禁止外部对象实现接口的做法也是有代价的：首先是**这个接口只能在包内部使用**，外部包在正常情况下是无法直接创建满足该接口对象的；其次，这种防护措施也不是绝对的，恶意的用户依然可以绕过这种保护机制。

通过在结构体中嵌入匿名类型成员，可以继承匿名类型的方法。其实这个被嵌入的匿名成员不一定是普通类型，也可以是接口类型。我们可以通过嵌入匿名的 `testing.TB` 接口来伪造私有方法，因为接口方法是延迟绑定，所以**编译时私有方法是否真的存在并不重要**。

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

这种**通过嵌入匿名接口或嵌入匿名指针对象来实现继承**的做法其实是一种**纯虚继承**，**继承的只是接口指定的规范**，**真正的实现在运行的时候才被注入**。例如，可以模拟实现一个 `gRPC` 的插件：

构造的 `grpcPlugin` 类型对象必须满足 `generate.Plugin` 接口：

```go
type Plugin interface {
    // Name identifies the plugin.
    Name() string
    // Init is called once after data structures are built but before
    // code generation begins.
    Init(g *Generator)
    // Generate produces the code generated by the plugin for this file,
    // except for the imports, by calling the generator's methods
    // P, In, and Out.
    Generate(file *FileDescriptor)
    // GenerateImports produces the import declarations for this file.
    // It is called after Generate.
    GenerateImports(file *FileDescriptor)
}
```

`generate.Plugin` 接口对应的 `grpcPlugin` 类型的 `GenerateImports()` 方法中使用的 `p.P(...)` 函数，是通过 `Init()` 函数注入的 `generator.Generator` 对象实现。这里的 `generator.Generator` 对应一个具体类型，但是如果 `generator.Generator` 是接口类型，我们甚至可以传入直接的实现。
 
```go
import "github.com/golang/protobuf/protoc-gen-go/generator"

type grpcPlugin struct {
    *generator.Generator
}

func (p *grpcPlugin) Name() string { return "grpc" }

func (p *grpcPlugin) Init(g *generator.Generator) {
    p.Generator = g
}

func (p *grpcPlugin) GenerateImports(file *generator.FileDescriptor) {
    if len(file.Service) == 0 {
        return
    }

    p.P(`import "google.golang.org/grpc"`)
    // ...
}
```
