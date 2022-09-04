---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言接口"  # 文章标题
url:  "posts/go/docs/grammar/interface"  # 设置网页永久链接
tags: [ "Go", "interface" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

一般静态编程语言都有着严格的类型系统，这使得编译器可以深入检查程序员有没有作出什么出格的举动。但是，过于严格的类型系统却会使得编程太过烦琐，让程序员把时间都浪费在了和编译器的斗争中。

Go 语言试图让程序员能在安全和灵活的编程之间取得一个平衡。它在提供严格的类型检查的同时，通过接口类型实现了对鸭子类型的支持，使得安全动态的编程变得相对容易。

Go 的接口类型是对其他类型行为的抽象和概括，因为接口类型不会和特定的实现细节绑定在一起，通过这种抽象的方式我们可以让对象更加灵活和更具有适应能力。很多面向对象的语言都有相似的接口概念，但 Go 语言中接口类型的独特之处在于它是满足隐式实现的鸭子类型。

所谓鸭子类型说的是：只要走起路来像鸭子、叫起来也像鸭子，那么就可以把它当作鸭子。

Go 语言中的面向对象就是如此，如果一个对象只要看起来像是某种接口类型的实现，那么它就可以作为该接口类型使用。这种设计可以让你创建一个新的接口类型满足已经存在的具体类型却不用去破坏这些类型原有的定义。当使用的类型来自不受我们控制的包时这种设计尤其灵活有用。

**Go 语言的接口类型是延迟绑定**，可以实现类似虚函数的多态功能。

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

通过在结构体中嵌入匿名类型成员，可以继承匿名类型的方法。其实这个被嵌入的匿名成员不一定是普通类型，也可以是接口类型。我们可以通过嵌入匿名的 `testing.TB` 接口来伪造私有方法，因为**接口方法是延迟绑定**，所以**编译时私有方法是否真的存在并不重要**。

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

```go

```
