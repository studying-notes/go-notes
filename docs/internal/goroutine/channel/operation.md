---
date: 2022-10-09T16:48:26+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "通道基本使用方式"  # 文章标题
url:  "posts/go/docs/internal/goroutine/channel/operation"  # 设置网页永久链接
tags: [ "Go", "operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 通道声明与初始化

```go
var ch chan T
```

其中，ch 代表 chan 的名字，为用户自定义的；`chan T` 代表通道的类型，T 代表通道中的元素类型。在声明时，channel 必须与一个实际的类型 T 绑定在一起，代表通道中能够读取和传递的元素类型。

通道的表示形式有如下有三种：

```go
chan T
chan <- T
<- chan T
```

具体来说，不带“<-”的通道可读可写，而带“<-”的类型限制了通道的读写。

例如，`chan <- float` 代表该通道只能写入浮点数，`<- chan string` 代表该通道只能读取字符串。

一个仅声明还未初始化的通道会被预置为 nil。

```go
var ch chan int
fmt.Println(ch) // <nil>
```

一个未初始化的通道在编译时和运行时并不会报错，不过无法向通道中写入或读取任何数据。要对通道进行操作，需要使用 make 操作符，make 会初始化通道，在内存中分配通道的空间。

```go
ch := make(chan int)
```

## 通道写入数据

可以通过如下简单的方式向通道中写入数据：

```go
ch <- 1
```

对于无缓冲通道，能够向通道写入数据的前提是必须有另一个协程在读取通道。否则，当前的协程会陷入休眠状态，直到能够向通道中成功写入数据。

**无缓冲通道的读与写应该位于不同的协程中**，否则，程序将陷入死锁的状态，如下所示。

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int)
    ch <- 1
    fmt.Println(<-ch)
}
```

## 通道读取数据

```go
fmt.Println(<-ch)
```

和写入数据一样，如果不能直接读取通道的数据，那么当前的读取协程将陷入堵塞，直到有协程写入通道为止。

读取通道还有两种返回值的形式，借助编译时将该形式转换为不同的处理函数。第 1 个返回值仍然为通道读取到的数据，第 2 个返回值为布尔类型，返回值为 false 代表当前通道已经关闭。

```go
v, ok := <-ch
```

## 关闭通道

关闭通道的语法如下：

```go
close(ch)
```

在正常读取的情况下，通道返回的 ok 为 true。**通道在关闭时仍然会返回，但是 data 为其类型的零值**，ok 也变为了 false。和通道读取不同的是，不能向已经关闭的通道中写入数据。

**通道关闭会通知所有正在读取通道的协程**，相当于向所有读取协程中都写入了数据。

如下所示，有两个协程正在等待通道中的数据，当 main 协程关闭通道后，两个协程都会收到通知。

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var c = make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		data, ok := <-c
		fmt.Println("goroutine 1", data, ok)
	}()

	go func() {
		defer wg.Done()
		data, ok := <-c
		fmt.Println("goroutine 2", data, ok)
	}()

	close(c)

	wg.Wait()
}
```

两个协程都读取到了结果，但是结果都是零值。

要注意的是，如果读取通道是一个循环操作，关闭通道并不能终止循环，依然会收到一个永无休止的零值序列。

```go
go func() {
    for {
        data, ok := <-c
        fmt.Println("goroutine 1", data, ok)
    }
}
```

因此，在实践中会通过第二个返回的布尔值来判断通道是否已经关闭，如果已经关闭，那么退出循环是一种比较常见的操作。

```go
go func() {
    for {
        data, ok := <-c
        if !ok {
            break
        }
        fmt.Println("goroutine 1", data, ok)
    }
}
```

试图**重复关闭一个 channel 将导致 panic 异常，试图关闭一个 nil 值的 channel 也将导致 panic 异常。**

在实践中，并不需要关心是否所有的通道都已关闭，**当通道没有被引用时将被 Go 语言的垃圾自动回收器回收**。

关闭通道会**触发所有通道读取操作被唤醒**的特性，被使用在了很多重要的场景中，例如一个协程退出后，其创建的一系列子协程能够快速退出的场景。

## 通道作为参数和返回值

通道作为一等公民，可以作为参数和返回值。通道是协程之间交流的方式，不管是将数据读取还是写入通道，都需要将代表通道的变量通过函数传递到所在的协程中去。

如下所示，代表协程执行的 worker 函数以通道 (chan int) 作为参数。

```go
func worker(ch chan int) {
    // ...
}
```

通道作为返回值一般用于创建通道的阶段，下例中的 createWorker 函数创建了通道 c，并新建了一个 worker 协程，最后返回的通道可能继续传递给其他的消费者使用。

```go
func createWorker() chan int {
    c := make(chan int)
    go worker(c)
    return c
}
```

由于通道是 Go 语言中的引用类型而不是值类型，因此传递到其他协程中的通道，实际引用了同一个通道，它在运行时是 `*hchan` 类型指针，这跟切片不太一样。

## 单方向通道

一般来说，一个协程在大多数情况下只会读取或者写入通道，为了表达这种语义并防止通道被误用，Go 语言的类型系统提供了单方向的通道类型。

例如 chan <- float 代表该通道只能写入而不能读取浮点数，<- chan string 代表该通道只能读取而不能写入字符串。

上例中 worker 函数的作用主要是读取通道的信息，因此可以将其函数签名改写如下，而不影响其任何功能。

```go
func worker(ch <- chan int) {
    // ...
}
```

普通的通道具有读和写的功能，普通的通道类型能够隐式地转换为单通道的类型。

```go
var c chan int
var c1 chan <- int = c
```

反之，单通道的类型不能转换为普通的通道类型。

## select 多路复用

在实践中使用通道时，更多的时候会与 select 结合，因为时常会出现多个通道与多个协程进行通信的情况，我们当然不希望由于一个通道的读写陷入堵塞，影响其他通道的正常读写。select 正是为了解决这一问题诞生的。

select 是 go 在语言层面提供的多路 IO 复用的机制，其可以检测多个 channel 是否 ready (即是否可读或可写)。

在使用方法上，select 的语法类似 switch，形式如下：

```go
select {
    case <-ch1:
        // ...
    case ch2 <- 1:
        // ...
    default:
        // ...
}
```

和 switch 不同的是，每个 case 语句都必须对应通道的读写操作。select 语句会陷入堵塞，直到一个或多个通道能够正常读写才恢复。

### 随机选择机制

如下所示，向通道 c 中写入数据 1，虽然两个 case 都能够读取到通道的内容，但是当我们多次执行程序时会发现，程序有时会输出 random 01，有时会输出 random 02。

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    c := make(chan int)
    go func() {
        for {
            select {
            case <-c:
                fmt.Println("random 01", n)
            case <-c:
                fmt.Println("random 02")
            }
        }
    }()

    c <- 1
}
```

这就是 select 的特性之一，case 是随机选取的，所以当 select 有两个通道同时准备好时，会随机执行不同的 case。每次执行 select 语句时，会将 case 语句的顺序打乱，然后按照打乱后的顺序依次执行。其中用到了类似洗牌算法的方式，将序列打散。通过引入随机数的方式给序列带来了随机性。

### 堵塞与控制

如果 select 中没有任何的通道准备好，那么当前 select 所在的协程会永远陷入等待，直到有一个 case 中的通道准备好为止。

在实践中，为了避免这种情况发生，有时会加上 default 分支。default 分支的作用是当所有的通道都陷入堵塞时，正常执行 default 分支。

除了 default，还有一些其他的选择。例如，如果我们希望一段时间后仍然没有通道准备好则超时退出，可以选择 select 与定时器或者超时器配套使用。

```go
select {
    case <-ch1:
        // ...
    case ch2 <- 1:
        // ...
    case <-time.After(3 * time.Second):
        // ...
}
```

### 空语句

```go
package main

func main() {
    select {
    }
}
```

上面程序中只有一个空的 select 语句。

对于空的 select 语句，程序会被阻塞，准确的说是当前协程被阻塞，同时 go 自带死锁检测机制，当发现当前协程再也没有机会被唤醒时，则会 panic。所以上述程序会 panic。

### 循环

很多时候，我们不希望 select 执行完一个分支就退出，而是循环往复执行 select 中的内容，因此需要将 for 与 select 进行组合，如下所示。

```go
for {
    select {
        case <-ch1:
            // ...
        case ch2 <- 1:
            // ...
    }
}
```

for 与 select 组合后，可以向 select 中加入一些定时任务。下例中的 tick 每隔 1s 就会向 tick 通道中写入数据，从而完成一些定时任务。

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    tick := time.Tick(1 * time.Second)
    for {
        select {
        case <-tick:
            fmt.Println("tick")
        }
    }
}
```

需要注意的是，定时器 time.Tick 与 time.After 是有本质不同的。time.After 并不会定时发送数据到通道中，而只是在时间到了后发送一次数据。当其放入 for+select 后，新一轮的 select 语句会重置 time.After，这意味着第 2 次 select 语句依然需要等待 800ms 才执行超时。

如果在 800ms 之前，其他的通道就已经执行好了，那么 time.After 的 case 将永远得不到执行。而定时器 tick 不同，由于 tick 在 for 循环的外部，因此其不重置，只会累积时间，实现定时执行任务的功能。

## select 与 nil

**一个为 nil 的通道，不管是读取还是写入都将陷入堵塞状态**。当 select 语句的 case 对 nil 通道进行操作时，case 分支将永远得不到执行。

nil 通道的这种特性，可以用于设计一些特别的模式。例如，假设有 a、b 两个通道，我们希望交替地向 a、b 通道中发送消息，那么可以用如下方式：

```go
package main

import "fmt"

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			select {
			case a <- i:
				a = nil
				b = make(chan int)
			case b <- i:
				b = nil
				a = make(chan int)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		select {
		case v := <-a:
			fmt.Println("a:", v)
		case v := <-b:
			fmt.Println("b:", v)
		}
	}
}
```

上例的协程中，一旦写入通道后，就将该通道置为 nil，导致再也没有机会执行该 case。从而达到交替写入 a、b 通道的目的。

这里是交替写入两个通道，不是用两个协程交替打印。下面的例子中，我们用两个协程交替打印。

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	a := make(chan bool)
	b := make(chan bool)
	c := make(chan bool)

	wg := sync.WaitGroup{}
	wg.Add(2)

	n := 0

	go func() {
		for {
			<-a

			for i := 0; i < 5; i++ {
				n++
				fmt.Println("a:", n)
			}

			if n == 100 {
				close(c)
			} else {
				b <- true
			}
		}
	}()

	go func() {
		for {
			<-b

			for i := 0; i < 5; i++ {
				n++
				fmt.Println("b:", n)
			}

			if n == 100 {
				close(c)
			} else {
				a <- true
			}
		}
	}()

	a <- true
	<-c
}
```

### range 读取数据

通过 range 可以持续从 channel 中读出数据，好像在遍历一个数组一样，当 channel 中没有数据时会阻塞当前 goroutine，与读 channel 时阻塞处理机制一样。

```go
func chanRange(chanName chan int) {
    for e := range chanName {
        fmt.Printf("Get element from chan: %d\n", e)
    }
}
```

如果向此 channel 写数据的 goroutine 退出时，系统检测到这种情况后会 panic，否则 range 将会永久阻塞。

关闭 channel 后，range 会正常退出，不会再阻塞。

```go

```
