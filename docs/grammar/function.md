---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言函数"  # 文章标题
url:  "posts/go/docs/grammar/function"  # 设置网页永久链接
tags: [ "Go", "function" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

函数对应操作序列，是程序的基本组成元素。

Go 语言中的函数有具名和匿名之分：**具名函数一般对应于包级的函数**，是匿名函数的一种特例。

**当匿名函数引用了外部作用域中的变量时就成了闭包函数**，闭包函数是函数式编程语言的核心。

方法是绑定到一个具体类型的特殊函数，Go 语言中的方法是依托于类型的，必须**在编译时静态绑定**。

接口定义了方法的集合，这些方法依托于运行时的接口对象，因此**接口对应的方法是在运行时动态绑定的**。Go 语言通过**隐式接口机制**实现了鸭子面向对象模型。

在 Go 语言中，函数是第一类对象，可以将函数保存到变量中。函数主要有具名和匿名之分，包级函数一般都是具名函数，具名函数是匿名函数的一种特例。当然，Go语言中每个类型还可以有自己的方法，方法其实也是函数的一种。

```go
// 具名函数
func Add(a, b int) int {
    return a+b
}

// 匿名函数
var Add = func(a, b int) int {
    return a+b
}
```

Go 语言中的函数**可以有多个参数和多个返回值**，参数和返回值都是**以传值的方式**和被调用者交换数据。在语法上，函数还支持可变数量的参数，**可变数量的参数必须是最后出现的参数**，可变数量的参数其实是一个**切片类型的参数**。

```go
// 多个参数和多个返回值
func Swap(a, b int) (int, int) {
    return b, a
}

// 可变数量的参数
// more 对应 []int 切片类型
func Sum(a int, more ...int) int {
    for _, v := range more {
        a += v
    }
    return a
}
```

当可变参数是一个空接口类型时，调用者是否解包可变参数会导致不同的结果：

```go
func Print(a ...interface{}) {
    fmt.Println(a...)
}

func main() {
    var a = []interface{}{123, "abc"}

    Print(a...) // 123 abc
    Print(a)    // [123 abc]
}
```

第一个 `Print` 调用时传入的参数是 `a...`，等价于直接调用 `Print(123, "abc")`。第二个 `Print` 调用传入的是未解包的 `a`，等价于直接调用 `Print([]interface{}{123, "abc"} )`。

不仅函数的参数可以有名字，也可以给函数的返回值命名：

```go
func Find(m map[int]int, key int) (value int, ok bool) {
    value, ok = m[key]
    return
}
```

如果返回值命名了，可以通过名字来修改返回值，也可以通过 defer 语句在 return 语句之后修改返回值：

```go
func Inc() (v int) {
    defer func(){ v++ } ()
    return 42
}
```

其中 defer 语句延迟执行了一个匿名函数，因为这个匿名函数**捕获了外部函数的局部变量** v，这种函数我们一般称为闭包。**闭包对捕获的外部变量并不是以传值方式访问，而是以引用方式访问**。

闭包的这种以**引用方式访问外部变量**的行为可能会导致一些隐含的问题：

```go
func main() {
    for i := 0; i < 3; i++ {
        defer func(){ println(i) } ()
    }
}
// Output:
// 3
// 3
// 3
```

因为是闭包，在 for 迭代语句中，每个 defer 语句延迟执行的函数引用的都是同一个 i 迭代变量，在循环结束后这个变量的值为 3，因此最终输出的都是 3。

修复的思路是在每轮迭代中为每个 defer 语句的闭包函数生成独有的变量。可以用下面两种方式：

```go
func main() {
    for i := 0; i < 3; i++ {
        i := i // 定义一个循环体内局部变量i
        defer func(){ println(i) } ()
    }
}

func main() {
    for i := 0; i < 3; i++ {
        // 通过函数传入i
        // defer 语句会马上对调用参数求值
        defer func(i int){ println(i) } (i)
    }
}
```

第一种方法是在循环体内部再定义一个局部变量，这样每次迭代 defer 语句的闭包函数捕获的都是不同的变量，这些变量的值对应迭代时的值。

第二种方式是将迭代变量通过闭包函数的参数传入，defer 语句会马上对调用参数求值。两种方式都是可以工作的。不过一般来说，在 for 循环内部执行 defer 语句并不是一个好的习惯，此处仅为示例，不建议使用。

Go 语言中，如果以切片为参数调用函数，有时候会给人一种参数采用了传引用的方式的假象：因为在被调用函数内部可以修改传入的切片的元素。其实，**任何可以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数**。

函数参数传值的规范更准确说是**只针对数据结构中固定的部分传值**，例如字符串或切片对应结构体中的指针和字符串长度结构体传值，但是并**不包含指针间接指向的内容**。将切片类型的参数替换为类似 `reflect.SliceHeader` 结构体就能很好理解切片传值的含义了：

```go
func twice(x []int) {
    for i := range x {
        x[i] *= 2
    }
}

type IntSliceHeader struct {
    Data []int
    Len  int
    Cap  int
}

func twice(x IntSliceHeader) {
    for i := 0; i < x.Len; i++ {
        x.Data[i] *= 2
    }
}
```

因为切片中的底层数组部分通过隐式指针传递（指针本身依然是传值的，但是指针指向的却是同一份的数据），所以被调用函数可以通过指针修改调用参数切片中的数据。除数据之外，切片结构还包含了切片长度和切片容量信息，这两个信息也是传值的。如果被调用函数中修改了 Len 或 Cap 信息，就无法反映到调用参数的切片中，这时候我们一般会通过返回修改后的切片来更新之前的切片。这也是内置的 `append ()` 必须要返回一个切片的原因。

Go 语言中，函数还可以直接或间接地调用自己，也就是支持递归调用。Go 语言函数的递归调用深度在逻辑上没有限制，函数调用的栈是不会出现溢出错误的，因为 Go 语言运行时会根据需要动态地调整函数栈的大小。

每个 Goroutine 刚启动时只会分配很小的栈（4 KB 或 8 KB，具体依赖实现），根据需要动态调整栈的大小，栈最大可以达到 GB 级（依赖具体实现，在目前的实现中，32 位体系结构为 250 MB，64 位体系结构为 1 GB）。

在 Go 1.4 以前，Go 的动态栈采用的是分段式的动态栈，通俗地说就是采用一个链表来实现动态栈，每个链表的节点内存位置不会发生变化。但是链表实现的动态栈对某些导致跨越链表不同节点的热点调用的性能影响较大，因为相邻的链表节点在内存位置一般不是相邻的，这会增加 CPU 高速缓存命中失败的概率。

为了解决热点调用的 CPU 缓存命中率问题，Go 1.4 之后改用连续的动态栈实现，也就是采用一个类似动态数组的结构来表示栈。不过连续动态栈也带来了新的问题：当连续栈动态增长时，需要将之前的数据移动到新的内存空间，这会导致之前栈中全部变量的地址发生变化。

虽然 Go 语言运行时会自动更新引用了地址变化的栈变量的指针，但最重要的一点是要明白 Go 语言中指针不再是固定不变的（因此不能随意将指针保存到数值变量中，Go 语言的地址也不能随意保存到不在垃圾回收器控制的环境中，因此使用 CGO 时不能在 C 语言中长期持有 Go 语言对象的地址）。

因为 Go 语言函数的栈会自动调整大小，所以普通 Go 程序员已经很少需要关心栈的运行机制了。在 Go 语言规范中甚至故意没有讲到栈和堆的概念。**我们无法知道函数参数或局部变量到底是保存在栈中还是堆中**，我们只需要知道它们能够正常工作就可以了。看看下面这个例子：

```go
func f(x int) *int {
    return &x
}

func g() int {
    x = new(int)
    return *x
}
```

第一个函数直接返回了函数参数变量的地址——这似乎是不可以的，因为如果参数变量在栈上，函数返回之后栈变量就失效了，返回的地址自然也应该失效了。但是 Go 语言的编译器和运行时比我们聪明得多，它会保证指针指向的变量在合适的地方。

第二个函数，内部虽然调用 new() 函数创建了 *int 类型的指针对象，但是依然不知道它具体保存在哪里。

不用关心 Go 语言中函数栈和堆的问题，编译器和运行时会帮我们搞定；同样不要假设变量在内存中的位置是固定不变的，指针随时可能会变化，特别是在你不期望它变化的时候。

```go

```