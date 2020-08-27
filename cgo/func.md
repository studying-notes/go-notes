# 函数调用

- [函数调用](#函数调用)
  - [Go 调用 C 函数](#go-调用-c-函数)
  - [C 函数的返回值](#c-函数的返回值)
    - [数值型返回值](#数值型返回值)
    - [获取错误状态码](#获取错误状态码)
    - [void 型返回值](#void-型返回值)
    - [字符串型返回值](#字符串型返回值)
    - [字符串型返回值及其长度](#字符串型返回值及其长度)
    - [结构体型返回值](#结构体型返回值)
  - [C 调用 Go 导出函数](#c-调用-go-导出函数)

## Go 调用 C 函数

对于一个启用 CGO 特性的程序，CGO 会构造一个虚拟的 C 包。通过这个虚拟的 C 包可以调用 C 语言函数。

```go
package main

/*
#include <stdio.h>

static int add(int a, int b) {
    printf("%d", a+b);
}
*/
import "C"

func main() {
	C.add(1, 1)
}
```

## C 函数的返回值

对于有返回值的 C 函数，我们可以正常获取返回值。

### 数值型返回值

```go
package main

/*
static int div(int a, int b) {
    return a/b;
}

int add(int a, int b) {
    return a+b;
}
*/
import "C"
import "fmt"

func main() {
    v := C.div(6, 3)
    fmt.Println(v)
    
    // 转换为 Go 数值
    n := int(C.add(4, 8))
  	fmt.Printf("%T\n%d", n, n)
}
```

上面的 `div` 函数实现了一个整数除法的运算，然后通过返回值返回除法的结果。数值类型可以直接转换。

### 获取错误状态码

上一个示例中，对于除数为 0 的情形并没有做特殊处理。我们希望在除数为 0 的时候返回一个错误，而其他时候返回正常的结果。因为 C 语言不支持返回多个结果，因此 `<errno.h>` 标准库提供了一个 `errno` 宏用于返回错误状态。我们可以近似地将 `errno` 看成一个线程安全的全局变量，可以用于记录最近一次错误的状态码。

改进后的 `div` 函数实现如下：

```c
#include <errno.h>

int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
```

CGO 也针对 `<errno.h>` 标准库的 `errno` 宏做了特殊支持：在 CGO 调用 C 函数时如果有两个返回值，那么第二个返回值将对应 `errno` 错误状态。

```go
package main

/*
#include <errno.h>

static int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
*/
import "C"
import "fmt"

func main() {
	v0, err0 := C.div(2, 1)
	fmt.Println(v0, err0)

	v1, err1 := C.div(1, 0)
	fmt.Println(v1, err1)
}
```

输出：

```
2 <nil>
0 The device does not recognize the command.
```

我们可以近似地将 `div` 函数看作为以下类型的函数：

```go
func C.div(a, b C.int) (C.int, [error])
```

第二个返回值是可忽略的 error 接口类型，底层对应 `syscall.Errno` 错误类型。

### void 型返回值

C 语言函数还有一种没有返回值类型的函数，用 `void` 表示返回值类型。一般情况下，我们无法获取 `void` 类型函数的返回值，因为没有返回值可以获取。

前面的例子中提到，CGO 对 `errno` 做了特殊处理，可以通过第二个返回值来获取 C 语言的错误状态。对于 `void` 类型函数，这个特性依然有效。

以下的代码是获取没有返回值函数的错误状态码：

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    _, err := C.noreturn()
    fmt.Println(err)
}
```

此时，我们忽略了第一个返回值，只获取第二个返回值对应的错误码。

我们也可以尝试获取第一个返回值，它对应的是 C 语言的 `void` 对应的 Go 语言类型：

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    v, _ := C.noreturn()
    fmt.Printf("%#v", v)
}
```

输出：

```
main._Ctype_void{}
```

我们可以看出 C 语言的 `void` 类型对应的是当前的 `main` 包中的 `_Ctype_void` 类型。其实也将 C 语言的 `noreturn` 函数看作是返回 `_Ctype_void` 类型的函数，这样就可以直接获取 `void` 类型函数的返回值：

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    fmt.Println(C.noreturn())
}
```

输出：

```
[]
```

其实在 CGO 生成的代码中，`_Ctype_void` 类型对应一个 0 长的数组类型 `[0]byte`，因此 `fmt.Println` 输出的是一个表示空数值的方括弧。

### 字符串型返回值

```go
package main

//#include <stdio.h>
//#include <stdlib.h>
//#define LENGTH 16
//
//char *RetString(char *input) {
//char *s = malloc(LENGTH * sizeof(char));
//sprintf(s, "%s", input);
//return s;
//}
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 往 C 函数传入字符串
	cs := C.CString("string")
	ret := C.RetString(cs)

	// 获取字符串返回值
	fmt.Println(C.GoString(ret))

	// 释放内存，防止泄漏
	C.free(unsafe.Pointer(cs))
	C.free(unsafe.Pointer(ret))
}
```

### 字符串型返回值及其长度

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 返回字符串型及其长度
char *RetStringInt(int len, int *rLen) {
    static const char *s = "0123456789";
    char *p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    *rLen = len;
    return p;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 获取字符串型返回值及其长度
	rLen := C.int(0)
	cStr := C.RetStringInt(C.int(10), &rLen)
	defer C.free(unsafe.Pointer(cStr))
	goStr := C.GoStringN(cStr, rLen)
	fmt.Printf("%v %v\n", rLen, goStr)
}
```

### 结构体型返回值

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct StringInfo {
    char *s;
    int len;
};

// 返回自定义结构体
struct StringInfo RetStruct(int len) {
    static const char *s = "0123456789";
    char *p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    struct StringInfo str;
    str.s = p;
    str.len = len;
    return str;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 获取结构体型返回值
	cStruct := C.RetStruct(C.int(10))
	defer C.free(unsafe.Pointer(cStruct.s))
	str := C.GoStringN(cStruct.s, cStruct.len)
	fmt.Printf("%v %v\n", cStruct.len, str)
}
```

## C 调用 Go 导出函数

```go
import "C"

//export add
func add(a, b C.int) C.int {
    return a+b
}
```

`add` 函数名以小写字母开头，对于 Go 语言来说是包内的私有函数。但是从 C 语言角度来看，导出的 `add` 函数是一个可全局访问的 C 语言函数。如果在两个不同的 Go 语言包内，都存在一个同名的要导出为 C 语言函数的 `add` 函数，那么在最终的链接阶段将会出现符号重名的问题。

CGO 生成的 `_cgo_export.h` 文件会包含导出后的 C 语言函数的声明。我们可以在纯 C 源文件中包含 `_cgo_export.h` 文件来引用导出的 `add` 函数。如果希望在当前的 CGO 文件中马上使用导出的 C 语言 add 函数，则无法引用 `_cgo_export.h` 文件。因为`_cgo_export.h` 文件的生成需要依赖当前文件可以正常构建，而如果当前文件内部循环依赖还未生成的`_cgo_export.h` 文件将会导致 cgo 命令错误。

当导出 C 语言接口时，需要保证函数的参数和返回值类型都是 C 语言友好的类型，同时返回值不得直接或间接包含 Go 语言内存空间的指针。
