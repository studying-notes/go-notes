# CGO 功能预览

- [CGO 功能预览](#cgo-功能预览)
	- [基于 C 标准库函数输出字符串](#基于-c-标准库函数输出字符串)
		- [手动释放资源](#手动释放资源)
	- [用自己定义的 C 函数](#用自己定义的-c-函数)
		- [额外的 C 语言源文件](#额外的-c-语言源文件)
	- [C 代码的模块化](#c-代码的模块化)
	- [用 Go 重新实现 C 函数](#用-go-重新实现-c-函数)
	- [面向 C 接口的 Go 编程](#面向-c-接口的-go-编程)

## 基于 C 标准库函数输出字符串

```go
package main

//#include <stdio.h>
import "C" // 必须单独一行

func main() {
	C.puts(C.CString("Hello, World\n"))
}
```

通过 `import "C"` 语句启用 CGO 特性，同时包含 C 语言的 `<stdio.h>` 头文件。然后通过 CGO 包的 `C.CString` 函数将 Go 语言字符串转为 C 语言字符串，最后调用 CGO 包的 `C.puts` 函数向标准输出窗口打印转换后的 C 字符串。另外，`import "C"` 必须单独一行，且与上面的 C 代码之间**不可以有未注释掉的空行**。

我们没有在程序退出前释放 `C.CString` 创建的 C 语言字符串；还有我们改用 `puts` 函数直接向标准输出打印，之前是采用 `fputs` 向标准输出打印。

没有释放使用 `C.CString` 创建的 C 语言字符串会导致内存泄漏。但是对于这个小程序来说，这样是没有问题的，因为程序退出后操作系统会自动回收程序的所有资源。

### 手动释放资源

```go
package main
// #include <stdio.h>
// #include <stdlib.h>
/*
void print(char *str) {
    printf("%s\n", str);
}
*/
import "C"
import "unsafe"

func main() {
	cs := C.CString("hello")
	// 释放资源
    defer C.free(unsafe.Pointer(cs))
    C.print(cs)
}
```

## 用自己定义的 C 函数

```go
package main

/*
#include <stdio.h>

static void SayHello(const char *s) {
    puts(s);
}
*/
import "C"

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}
```

### 额外的 C 语言源文件

```c
// hello.c

#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}
```

```go
// hello.go
package main

//void SayHello(const char* s);
import "C"

func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
```

但是，这种情况下不可以运行单独文件，而必须编译该文件夹。

## C 代码的模块化

在编程过程中，抽象和模块化是将复杂问题简化的通用手段。当代码语句变多时，我们可以将相似的代码封装到一个个函数中；当程序中的函数变多时，我们将函数拆分到不同的文件或模块中。而模块化编程的核心是面向程序接口编程。

在前面的例子中，我们可以抽象一个名为 `hello` 的模块，模块的全部接口函数都在 `hello.h` 头文件定义：

```c
// hello.h
void SayHello(const char* s);
```

```c
// hello.c

#include "hello.h"
#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}
```

也可以用 C++ 实现该接口：

```c++
// hello.cpp

#include <iostream>

extern "C" {
    #include "hello.h"
}

void SayHello(const char* s) {
    std::cout << s;
}
```

## 用 Go 重新实现 C 函数

CGO 不仅仅用于 Go 语言中调用 C 语言函数，还可以用于导出 Go 语言函数给 C 语言函数调用。

```c
// hello.h
void SayHello(/*const*/ char* s);
```

用 Go 实现 C 接口：

```go
// hello.go
package main

import "C"
import "fmt"

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

通过 cgo 的 `//export SayHello` 指令将 Go 语言实现的函数 `SayHello` 导出为 C 语言函数。为了适配 cgo 导出的 C 语言函数，我们禁止了在函数的声明语句中的 `const` 修饰符。需要注意的是，这里其实有两个版本的 `SayHello` 函数：一个 Go 语言环境的；另一个是 C 语言环境的。cgo 生成的 C 语言版本 `SayHello` 函数最终会通过桥接代码调用 Go 语言版本的 `SayHello` 函数。

```go
// main.go
package main

//#include <hello.h>
import "C"

func main() {
    C.SayHello(C.CString("Hello, World\n"))
}
```

在编译时注释掉其他语言的接口实现文件，否则报错。

## 面向 C 接口的 Go 编程

将上文中的几个文件重新合并到一个Go文件：

```go
package main

//void SayHello(char* s);
import "C"

import (
	"fmt"
)

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}
```

在 Go1.10 中 cgo 新增加了一个 `_GoString_` 预定义的 C 语言类型，用来表示 Go 语言字符串。因此可以简化上面的代码：

```go
// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"

import (
	"fmt"
)

func main() {
	C.SayHello("Hello, World\n")
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}
```

> C 函数和 Go 函数仍在同一个 Goroutine 里执行。
