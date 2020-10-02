# Goroutine 并发模型

- [Goroutine 并发模型](#goroutine-并发模型)
  - [Goroutine 和系统线程](#goroutine-和系统线程)
  - [原子操作](#原子操作)
    - [sync.Mutex](#syncmutex)
    - [sync/atomic](#syncatomic)
  - [顺序一致性内存模型](#顺序一致性内存模型)
  - [初始化顺序](#初始化顺序)
  - [Goroutine 的创建](#goroutine-的创建)
  - [基于通道的通信](#基于通道的通信)
  - [不靠谱的同步](#不靠谱的同步)

常见的并行编程有多种模型，主要有多线程、消息传递等。从理论上来看，**多线程和基于消息的并发编程是等价的**。由于**多线程并发模型**可以自然对应到多核的处理器，主流的操作系统因此也都提供了系统级的多线程支持，同时从概念上讲多线程似乎也更直观，因此多线程编程模型逐步被吸纳到主流的编程语言特性或语言扩展库中。

主流编程语言对**基于消息的并发编程模型**支持则相对较少。Go 语言是基于消息并发模型的集大成者，它将**基于 CSP 模型的并发编程**内置到了语言中，通过一个 go 关键字就可以轻易地启动一个 Goroutine，Go 语言的 **Goroutine 之间是共享内存的**。

## Goroutine 和系统线程

Goroutine 是 Go 语言特有的并发体，是一种**轻量级的线程**，由 go 关键字启动。在真实的 Go 语言的实现中，**Goroutine 和系统线程也不是等价的**，两者的区别实际上只是一个量的区别。

首先，每个**系统级线程**都会有一个**固定大小的栈**（一般默认可能是2 MB），这个栈主要用来**保存函数递归调用时的参数和局部变量**。固定了栈的大小导致了两个问题：

- 一是对于很多只需要很小的栈空间的线程是一个巨大的浪费；
- 二是对于少数需要巨大栈空间的线程又面临栈溢出的风险。

针对这两个问题的解决方案是：要么降低固定的栈大小，提升空间的利用率；要么增大栈的大小以允许更深的函数递归调用，但这两者是无法兼得的。相反，一个 Goroutine 会以一个很小的栈启动（可能是2 KB或4 KB），当遇到深度递归导致当前栈空间不足时，Goroutine 会**根据需要动态地伸缩栈的大小**（主流实现中栈的最大值可达到1 GB）。因为**启动的代价很小**，所以我们可以轻易地启动成千上万个 Goroutine。

Go 的运行时还包含了其自己的调度器，这个调度器使用了一些技术手段，可以在 n 个操作系统线程上多工调度 m 个 Goroutine。Go 调度器的工作原理和内核的调度是相似的，但是这个调度器只关注单独的 Go 程序中的 Goroutine。

Goroutine 采用的是**半抢占式的协作调度**，只有在当前 Goroutine **发生阻塞时才会导致调度**；同时发生在用户态，调度器会根据具体函数只保存必要的寄存器，切换的代价要比系统线程低得多。运行时有一个 `runtime.GOMAXPROCS` 变量，用于控制当前运行正常非阻塞 Goroutine 的系统线程数目。
在 Go 语言中启动一个 Goroutine 不仅和调用函数一样简单，而且 Goroutine 之间调度代价也很低。

## 原子操作

原子操作是**并发编程中“最小的且不可并行化”的操作**。通常，如果多个并发体对同一个共享资源进行的操作是原子的话，那么同一时刻最多只能有一个并发体对该资源进行操作。从线程角度看，在当前线程修改共享资源期间，其他线程是不能访问该资源的。原子操作对多线程并发编程模型来说，不会发生有别于单线程的意外情况，共享资源的完整性可以得到保证。

### sync.Mutex

一般情况下，原子操作都是通过“互斥”访问来保证的，通常由特殊的 CPU 指令提供保护。当然，如果仅仅是想模拟粗粒度的原子操作，可以借助于 `sync.Mutex` 来实现：

```go
import (
    "sync"
)

var total struct {
    sync.Mutex
    value int
}

func worker(wg *sync.WaitGroup) {
    defer wg.Done()

    for i := 0; i <= 10000; i++ {
        total.Lock()
        total.value += i
        total.Unlock()
    }
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go worker(&wg)
    go worker(&wg)
    wg.Wait()

    fmt.Println(total.value)
}
```

在 `worker` 的循环中，为了保证 `total.value += i` 的原子性，我们通过 `sync.Mutex` 加锁和解锁来保证该语句在同一时刻只被一个线程访问。对多线程模型的程序而言，**进出临界区前后进行加锁和解锁都是必需的**。如果没有锁的保护，`total` 的最终值将由于多线程之间的竞争而可能不正确。

### sync/atomic

**用互斥锁来保护一个数值型的共享资源麻烦且效率低下**。标准库的 `sync/atomic` 包对原子操作提供了丰富的支持。我们可以重新实现上面的例子：

```go
import (
	"fmt"
	"sync"
	"sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()
	fmt.Println(total)
}
```

`atomic.AddUint64()` 函数调用保证了 `total` 的读取、更新和保存是一个原子操作，因此在多线程中访问也是安全的。

原子操作配合互斥锁可以实现非常高效的单件模式。**互斥锁的代价比普通整数的原子读写高很多**，在性能敏感的地方可以增加一个数字型的标志位，通过原子检测标志位状态降低互斥锁的使用次数来提高性能。

```go
import (
	"sync"
	"sync/atomic"
)

type singleton struct{}

var (
	instance    *singleton
	initialized uint32
	mu          sync.Mutex
)

func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}
```

我们将通用的代码提取出来，就成了标准库中 `sync.Once` 的实现：

```go
type Once struct {
    m    Mutex
    done uint32
}

func (o *Once) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }

    o.m.Lock()
    defer o.m.Unlock()

    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

基于 sync.Once 重新实现单件（singleton）模式：

```go
var (
    instance *singleton
    once     sync.Once
)

func Instance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

`sync/atomic` 包对基本数值类型及复杂对象的读写都提供了原子操作的支持。`atomic.Value` 原子对象提供了 `Load()` 和 `Store()` 两个原子方法，分别用于加载和保存数据，返回值和参数都是 `interface{}` 类型，因此可以用于任意的自定义复杂类型。

```go
var config atomic.Value // 保存当前配置信息

// 初始化配置信息
config.Store(loadConfig())

// 启动一个后台线程，加载更新后的配置信息
go func() {
    for {
        time.Sleep(time.Second)
        config.Store(loadConfig())
    }
}()

// 用于处理请求的工作者线程始终采用最新的配置信息
for i := 0; i < 10; i++ {
    go func() {
        for r := range requests() {
            c := config.Load()
            // ...
        }
    }()
}
```

这是一个简化的生产者消费者模型：后台线程生成最新的配置信息；前台多个工作者线程获取最新的配置信息。所有线程共享配置信息资源。

## 顺序一致性内存模型

如果只是想简单地在线程之间进行数据同步的话，原子操作已经为编程人员提供了一些同步保障。不过这种保障有一个前提：顺序一致性的内存模型。要了解顺序一致性，先看一个简单的例子：

```go
var a string
var done bool

func setup() {
    a = "hello, world"
    done = true
}

func main() {
    go setup()
    for !done {}
    print(a)
}
```

我们创建了 `setup` 线程，用于对字符串 `a` 的初始化工作，初始化完成之后设置 `done` 标志为 `true`。`main()` 函数所在的主线程中，通过 `for !done {}` 检测 `done` 变为 `true` 时，认为字符串初始化工作完成，然后进行字符串的打印工作。

但是 Go 语言并不保证在 `main()` 函数中观测到的对 `done` 的写入操作发生在对字符串 `a` 的写入操作之后，因此程序很可能打印一个空字符串。更糟糕的是，因为两个线程之间没有同步事件，`setup` 线程对 `done` 的写入操作甚至无法被 `main` 线程看到，`main()` 函数有可能陷入死循环中。

在 Go 语言中，**同一个 Goroutine 线程内部，顺序一致性的内存模型是得到保证的**。但是不同的 Goroutine 之间，并不满足顺序一致性的内存模型，需要通过明确定义的同步事件来作为同步的参考。**如果两个事件不可排序，那么就说这两个事件是并发的**。为了最大化并行，Go 语言的编译器和处理器在不影响上述规定的前提下可能会对执行语句重新排序（CPU也会对一些指令进行乱序执行）。

因此，如果在一个 Goroutine 中顺序执行 `a = 1; b = 2;` 这两个语句，虽然在当前的 Goroutine 中可以认为 `a = 1;` 语句先于 `b = 2;` 语句执行，但是在另一个 Goroutine 中 `b = 2;` 语句可能会先于 `a = 1;` 语句执行，甚至在另一个 Goroutine 中无法看到它们的变化（可能始终在寄存器中）。也就是说在另一个 Goroutine 看来, `a = 1; b = 2;` 这两个语句的执行顺序是不确定的。如果一个并发程序无法确定事件的顺序关系，那么程序的运行结果往往会有不确定的结果。例如，下面这个程序：

```go
func main() {
    go println("你好, 世界")
}
```

根据 Go 语言规范，`main()` 函数退出时程序结束，不会等待任何后台线程。因为 Goroutine 的执行和 `main()` 函数的返回事件是并发的，谁都有可能先发生，所以什么时候打印、能否打印都是未知的。

用前面的原子操作并不能解决问题，因为我们无法确定两个原子操作之间的顺序。解决问题的办法就是**通过同步原语来给两个事件明确排序**：

```go
func main() {
    done := make(chan int)

    go func(){
        println("你好, 世界")
        done <- 1
    }()

    <-done
}
```

当 `<-done` 执行时，必然要求 `done <- 1` 也已经执行。根据同一个 Goroutine 依然满足顺序一致性规则，可以判断当 `done <- 1` 执行时，`println(" 你好 , 世界 ")` 语句必然已经执行完成了。因此，现在的程序确保可以正常打印结果。

当然，通过 `sync.Mutex` 互斥量也是可以实现同步的：

```go
func main() {
    var mu sync.Mutex

    mu.Lock()
    go func(){
        println("你好, 世界")
        mu.Unlock()
    }()

    mu.Lock()
}
```

可以确定，后台线程的 `mu.Unlock()` 必然在 `println(" 你好 , 世界 ")` 完成后发生（同一个线程满足顺序一致性），`main()` 函数的第二个 `mu.Lock()` 必然在后台线程的 `mu.Unlock()` 之后发生（`sync.Mutex` 保证），此时后台线程的打印工作已经顺利完成了。

## 初始化顺序

Go 程序的初始化和执行总是从 `main.main()` 函数开始的。但是如果 `main` 包里**导入了其他的包**，则会按照顺序将它们包含到 `main` 包里（这里的导入顺序依赖具体实现，一般可能是以文件名或包路径名的字符串顺序导入）。如果某个包被多次导入，那么在执行的时候只会导入一次。当一个包被导入时，如果它还导入了其他的包，则先**将其他的包包含进来**，然后**创建和初始化这个包的常量和变量**。再**调用包里的 `init()` 函数**，如果一个包有多个 `init()` 函数，实现可能是以文件名的顺序调用，那么同一个文件内的多个 `init()` 是以出现的顺序依次调用的（`init()` 不是普通函数，可以定义多个，但是**不能被其他函数调用**）。最终，**在 `main` 包的所有包常量、包变量被创建和初始化，并且只有在 `init()` 函数被执行后，才会进入 `main.main()` 函数，程序开始正常执行**。

Go 程序函数启动顺序的示意图：

![](imgs/main.png)

在 `main.main()` 函数执行之前所有代码都运行在同一个 Goroutine 中，也是运行在程序的主系统线程中。如果某个 `init()` 函数内部用 `go` 关键字启动了新的 Goroutine，那么新的 Goroutine 和 `main.main()` 函数是并发执行的。

因为所有的 `init()` 函数和 `main()` 函数都是在主线程完成，它们也是满足顺序一致性模型的。

## Goroutine 的创建

`go` 语句会在当前 Goroutine 对应函数返回前创建新的 Goroutine。例如:

```go
var a string

func f() {
    print(a)
}

func hello() {
    a = "hello, world"
    go f()
}
```

执行 `go f()` 语句创建 Goroutine 和 `hello()` 函数是在同一个 Goroutine 中执行，根据语句的书写顺序可以确定 Goroutine 的创建发生在 `hello()` 函数返回之前，但是新创建 Goroutine 对应的 `f()` 的执行事件和 `hello()` 函数返回的事件则是不可排序的，也就是并发的。调用 `hello()` 可能会在将来的某一时刻打印 “hello, world”，也很可能是在 `hello()` 函数执行完成后才打印。

## 基于通道的通信

通道（channel）是在 Goroutine 之间进行同步的主要方法。在无缓存的通道上的每一次发送操作都有与其对应的接收操作相匹配，发送和接收操作通常发生在不同的 Goroutine 上（在同一个 Goroutine 上执行两个操作很容易导致死锁）。无缓存的通道上的发送操作总在对应的接收操作完成前发生。

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "你好, 世界"
    done <- true
}

func main() {
    go aGoroutine()
    <-done
    println(msg)
}
```

可保证打印出“你好, 世界”。该程序首先对 `msg` 进行写入，然后在 `done` 通道上发送同步信号，随后从 `done` 接收对应的同步信号，最后执行 `println()` 函数。

若**在关闭通道后继续从中接收数据，接收者就会收到该通道返回的零值**。因此在这个例子中，用 `close( done)` 关闭通道代替 `done <- true` 依然能保证该程序产生相同的行为。

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "你好, 世界"
    close(done)
}

func main() {
    go aGoroutine()
    <-done
    println(msg)
}
```

对于从无缓存通道进行的接收，发生在对该通道进行的发送完成之前。

基于上面这个规则可知，交换两个 Goroutine 中的接收和发送操作也是可以的（但是很危险）：

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "hello, world"
    <-done
}
func main() {
    go aGoroutine()
    // 不带缓存，必须接收方准备好了才能发送
    done <- true
    println(msg)
}
```

也可保证打印出 “hello, world”。因为 `main` 线程中 `done <- true` 发送完成前，后台线程 `<-done` 接收已经开始，这保证 `msg = "hello, world"` 被执行了，所以之后 `println(msg)` 的 `msg` 已经被赋值过了。简而言之，后台线程首先对 `msg` 进行写入，然后从 `done` 中接收信号，随后 `main` 线程向 `done` 发送对应的信号，最后执行 `println()` 函数完成。但是，若该通道为带缓存的（例如，`done = make(chan bool, 1)`），`main` 线程的 `done <- true` 接收操作将不会被后台线程的 `<-done` 接收操作阻塞，该程序将无法保证打印出 “hello, world”。

对于带缓存的通道，对于通道中的第 `K` 个接收完成操作发生在第 `K+C` 个发送操作完成之前，其中 `C` 是管道的缓存大小。如果将 `C` 设置为 0 自然就对应无缓存的通道，也就是第 `K` 个接收完成在第 `K` 个发送完成之前。因为无缓存的通道只能同步发 1 个，所以也就简化为前面无缓存通道的规则：对于从无缓存通道进行的接收，发生在对该通道进行的发送完成之前。

我们可以根据控制通道的缓存大小来控制并发执行的 Goroutine 的最大数目，例如：

```go
var limit = make(chan int, 3)

func main() {
    for _, w := range work {
        go func() {
            limit <- 1
            w()
            <-limit
        }()
    }
    select{}
}
```

最后一句 `select{}` 是一个空的通道选择语句，该语句会导致 `main` 线程阻塞，从而避免程序过早退出。还有 `for{}`、`<-make(chan int)` 等诸多方法可以达到类似的效果。因为 `main` 线程被阻塞了，如果需要程序正常退出的话，可以通过调用 `os.Exit(0)` 实现。

## 不靠谱的同步

```go
func main() {
    go println("hello, world")
    time.Sleep(time.Second)
}
```

因为主线程休眠了1秒，所以这个程序很大概率上是可以正常输出结果的。因此，很多人会觉得这个程序已经没有问题了。但是这个程序是不稳健的，依然有失败的可能性。我们先假设程序是可以稳定输出结果的。因为 Go 线程的启动是非阻塞的，main 线程显式休眠了1秒退出导致程序结束，我们可以近似地认为程序总共执行了1秒多时间。现在假设 println() 函数内部实现休眠的时间大于 main 线程休眠的时间，就会导致矛盾：后台线程既然先于 main 线程完成打印，那么执行时间肯定是小于 main 线程执行时间的。当然这是不可能的。

严谨的并发程序的正确性不应该是依赖于 CPU 的执行速度和休眠时间等不靠谱的因素的。严谨的并发也应该是可以静态推导出结果的：根据线程内顺序一致性，结合通道或 sync 事件的可排序性来推导，最终完成各个线程各段代码的偏序关系排序。如果两个事件无法根据此规则来排序，那么它们就是并发的，也就是执行先后顺序不可靠的。

解决同步问题的思路是相同的：使用显式的同步。