---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 语言字符串"  # 文章标题
url:  "posts/go/docs/grammar/string"  # 设置网页永久链接
tags: [ "Go", "string" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## string 标准概念

Go 标准库 `builtin` 给出了所有内置类型的定义。源代码位于 `src/builtin/builtin.go`，其中关于 string 的描述如下:

```go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
```

所以 string 是 8 比特字节的集合，通常是但并不一定非得是 UTF-8 编码的文本。另外，还提到了两点：

* string 可以为空（长度为 0 ），但不会是 nil ；
* string 对象不可以修改。

## 底层数据结构

一个字符串是一个**不可改变的字节序列**，字符串的**元素不可修改**，是一个**只读的字节数组**。每个字符串的长度虽然也是固定的，但是**字符串的长度并不是字符串类型的一部分**。

Go 源代码中出现的字符串面值常量一般是 UTF-8 编码的的（对于转义字符，则没有这个限制）。源代码中的文本字符串通常被解释为采用 UTF-8 编码的 Unicode 码点（rune）序列。

因为字节序列对应的是只读的字节序列，所以字符串可以包含任意的数据，包括字节值 0。

我们也可以用字符串表示 GBK 等非 UTF-8 编码的数据，不过这时候将字符串看作是**一个只读的二进制数组**更准确，因为 `for range` 等语法并不支持非 UTF-8 编码的字符串的遍历。

Go 语言字符串的底层结构在 `reflect.StringHeader` 中定义：

`src\reflect\value.go`

```go
// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
	Data uintptr  // 字符串指向的底层字节数组
	Len  int  // 字符串的字节的长度
}
```

源码包 `src/runtime/string.go:stringStruct` 中的定义：

```go
type stringStruct struct {
	str unsafe.Pointer // 字符串的首地址
	len int  // 字符串的长度
}
```

字符串结构由两个信息组成：第一个是字符串指向的底层字节数组，**这个字节数组不是 []byte**，而是内存上一块连续的区域；第二个是字符串的字节的长度。字符串其实是一个结构体，因此字符串的赋值操作也就是 `reflect.StringHeader` 结构体的复制过程，并**不会涉及底层字节数组的复制**。可以将字符串数组看作一个结构体数组。

## 声明字符串

如下代码所示，可以声明一个 string 变量变赋予初值：

```go
    var str string
    str = "Hello World"
```

字符串构建过程是先根据字符串构建 stringStruct，再转换成 string。转换的源码如下：

```go
func gostringnocopy(str *byte) string { // 根据字符串地址构建 string
	ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)} // 先构造 stringStruct
	s := *(*string)(unsafe.Pointer(&ss))                             // 再将 stringStruct 转换成 string
	return s
}
```

string 在 runtime 包中就是 stringStruct，对外呈现叫做 string。

## 内存布局

字符串 `hello, world` 对应的内存结构：

![06NZGR.png](https://s1.ax1x.com/2020/10/10/06NZGR.png)

分析可以发现，`hello, world` 字符串底层数据和以下数组是完全一致的：

```go
var data = [...]byte{
    'h', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd',
}
```

## 切片操作

字符串虽然不是切片，但是**支持切片操作**，不同位置的切片底层**访问的是同一块内存数据**（因为字符串是只读的，所以**相同的字符串面值常量通常对应同一个字符串常量**）：

```go
s := "hello, world"
hello := s[:5]
world := s[7:]

s1 := "hello, world"[:5]
s2 := "hello, world"[7:]
```

字符串和数组类似，内置的 `len()` 函数返回字符串的长度。也可以通过 `reflect.StringHeader` 结构访问字符串的长度（这里只是为了演示字符串的结构，并不是推荐的做法）：

```go
func main() {
	s := "hello, world"
	s1 := "hello, world"[:5]
	s2 := "hello, world"[7:]

	fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)   // 12
	fmt.Println("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len) // 5
	fmt.Println("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len) // 5
}
```

## 遍历与打印

假设字符串对应的是一个合法的 `UTF-8` 编码的字符序列。可以用内置的 `print` 调试函数或 `fmt.Print()` 函数直接打印，也可以用 `for range` 循环直接遍历 `UTF-8` 解码后的 Unicode 码点值。

下面的 "hello,世界" 字符串中包含了中文字符，可以通过打印转型为字节类型来查看字符底层对应的数据：

```go
fmt.Printf("%#v\n", []byte("hello, 世界"))
```

```go
[]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
```

分析可以发现，`0xe4, 0xb8, 0x96` 对应中文“世”，`0xe7, 0x95, 0x8c` 对应中文“界”。

```go
fmt.Println("\xe4\xb8\x96") // 打印“世”
fmt.Println("\xe7\x95\x8c") // 打印“界”
```

“hello, 世界” 字符串的内存结构布局：

![06NeR1.png](https://s1.ax1x.com/2020/10/10/06NeR1.png)

一个中文字符占了 3 个字节。

```go
func main() {
	s := "hello, 世界"
	for idx, v := range s {  // idx 是字符的字节位置，v 是字符的拷贝
		fmt.Printf("%02d %c\t", idx, v)
	}
}
```

```
00 h    01 e    02 l    03 l    04 o    05 ,    06      07 世   10 界
```

Go 语言的字符串中可以存放任意的二进制字节序列，而且即使是 UTF-8 字符序列也可能会遇到错误的编码。如果遇到一个错误的 UTF-8 编码输入，将生成一个特别的 Unicode 字符 '\uFFFD'，这个字符在不同的软件中的显示效果可能不太一样，在印刷中这个符号通常是一个黑色六角形或钻石形状，里面包含一个白色的问号“�”。

下面的字符串中，我们故意损坏了第一字符的第二和第三字节，因此第一字符将会打印为“�”，第二和第三字节则被忽略，后面的 “abc” 依然可以正常解码打印（错误编码不会向后扩散是 UTF-8 编码的优秀特性之一）。

```go
fmt.Println("\xe4\x00\x00\xe7\x95\x8cabc") // �界a
```

不过在 `for range` 迭代这个含有损坏的 UTF-8 字符串时，第一字符的第二和第三字节依然会被单独迭代到，不过此时迭代的值是损坏后的 0：

```go
func main() {
	for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" {
		fmt.Println(i, c)
	}
}

// 0 65533  // \uFFF，对应�
// 1 0      // 空字符
// 2 0      // 空字符
// 3 30028  // 界
// 6 97     // a
// 7 98     // b
// 8 99     // c
```

### 遍历原始的字节码

如果不想解码 UTF-8 字符串，想直接遍历原始的字节码，可以将字符串强制转为 []byte 字节序列后再进行遍历（这里的转换一般不会产生运行时开销），或者是采用传统的下标方式遍历字符串的字节数组：

```go
func main() {
	s := "世界abc"
	for i, c := range []byte(s) {
		fmt.Print(i, c, '\t')
	}
	fmt.Println()
	for i := 0; i < len(s); i++ {
		fmt.Print(i, s[i], '\t')
	}
}
```

## 字符串和 `[]rune` 类型的相互转换

Go 语言除了 `for range` 语法对 UTF-8 字符串提供了特殊支持外，还对字符串和 `[]rune` 类型的相互转换提供了特殊的支持。

```go
func main() {
	fmt.Printf("%#v\n", []rune("世界"))             // []int32{19990, 30028}
	fmt.Printf("%#v\n", string([]rune{'世', '界'})) // 世界
}
```

从上面代码的输出结果可以发现 `[]rune` 其实是 `[]int32` 类型，这里的 `rune` 只是 `int32` 类型的别名，并不是重新定义的类型。`rune` 用于表示每个 Unicode 码点，目前只使用了 21 个位。

字符串相关的强制类型转换主要涉及 `[]byte` 和 `[]rune` 两种类型。每个转换都可能隐含重新分配内存的代价，最坏的情况下它们运算的时间复杂度都是 `O(n)`。不过字符串和 `[]rune` 的转换要更为特殊一些，因为一般这种强制类型转换要求两个类型的底层内存结构要尽量一致，显然它们底层对应的 `[]byte` 和 `[]int32` 类型是完全不同的内存结构，因此这种转换可能隐含重新分配内存的操作。

## 字节数组与字符串相互转换

string 不能直接和 byte 数组转换 string 可以和 byte 的切片转换。

### string 转为 []byte

```go
[]byte("string")
```

### [N]byte 转为 string

```go
var buf [10]byte
string(buf[:])
```

## 字符串内置操作模拟

简单模拟 Go 语言对字符串内置的一些操作，这样对每个操作的处理的时间复杂度和空间复杂度都会有较明确的认识。

### `for range` 遍历模拟

`for range` 对字符串的迭代模拟实现如下：

```go
func forRangeString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}

func main() {
	s := "hello, 世界"
	forRangeString(s, func(i int, r rune) {
		fmt.Printf("%d %c\t", i, r)
	})
}
```

`for range` 迭代字符串时，每次解码一个 Unicode 字符，然后进入 `for` 循环体，遇到崩溃的编码并不会导致迭代停止。

### `[]byte(s)` 转换模拟

`[]byte(s)` 转换模拟实现如下：

```go
func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		p[i] = s[i]
	}
	return p
}
```

string 转换成 byte 切片，也需要一次内存拷贝，其过程如下：
* 申请切片内存空间
* 将 string 拷贝到切片

[![DiNtsS.png](https://s3.ax1x.com/2020/11/15/DiNtsS.png)](https://imgchr.com/i/DiNtsS)

模拟实现中新创建了一个切片，然后**将字符串的数组逐一复制到切片中**，这是为了**保证字符串只读的语义**。当然，在将字符串转换为 `[]byte` 时，如果转换后的变量没有被修改，编译器可能会直接返回原始的字符串对应的底层数据。

### `string(bytes)` 转换模拟

这种转换需要一次内存拷贝。

`string(bytes)` 转换模拟实现如下：

```go
func bytes2str(b []byte) (s string) {
    // 为了性能，有时候会省略复制
	data := make([]byte, len(b))
	for i, c := range b {
		data[i] = c
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(&data[0]))
	hdr.Len = len(b)
	return s
}
```

转换过程如下：

1. 根据切片的长度申请内存空间，假设内存地址为 p，切片长度为 len(b) ；
2. 构建 string （ string.str = p ； string.len = len ；）
3. 拷贝数据 ( 切片中数据拷贝到新申请的内存空间 )

[![DiN8RP.png](https://s3.ax1x.com/2020/11/15/DiN8RP.png)](https://imgchr.com/i/DiN8RP)

因为 Go 语言的字符串是只读的，无法以直接构造底层字节数组的方式生成字符串。在模拟实现中通过 `unsafe` 包获取字符串的底层数据结构，然后将切片的数据逐一复制到字符串中，这同样是为了保证字符串只读的语义不受切片的影响。如果转换后的字符串在生命周期中原始的 `[]byte` 的变量不发生变化，编译器可能会直接基于 `[]byte` 底层的数据构建字符串，而不进行复制。

### `[]rune(s)` 转换模拟

`[]rune(s)` 转换模拟实现如下：

```go
func str2runes(s string) []rune {
	// 将字符串转换为字节数组
	b := []byte(s)
	var p []int32
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		p = append(p, r)
		s = s[size:]
	}
	return p
}
```

因为底层内存结构的差异，所以字符串到 []rune 的转换必然会导致重新分配 []rune 内存空间，然后依次解码并复制对应的 Unicode 码点值。这种强制转换并不存在前面提到的字符串和字节切片转换时的优化情况。

### `string(runes)` 转换模拟

`string(runes)` 转换模拟实现如下：

```go
func runes2string(rs []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range rs {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}
```

同样因为底层内存结构的差异，**[]rune 到字符串的转换也必然会导致重新构造字符串**。这种强制转换并不存在前面提到的优化情况。

## 字符串拼接

字符串可以很方便的拼接，像下面这样：

```go
str := "Str1" + "Str2" + "Str3"
```

即便有非常多的字符串需要拼接，性能上也有比较好的保证，因为新字符串的内存空间是一次分配完成的，所以性能消耗主要在拷贝数据上。

一个拼接语句的字符串编译时都会被存放到一个切片中，拼接过程需要遍历两次切片，第一次遍历获取总的字符串长度，据此申请内存，第二次遍历会把字符串逐个拷贝过去。

字符串拼接伪代码如下：

```go
func concatstrings(a []string) string { // 字符串拼接
    length := 0        // 拼接后总的字符串长度

    for _, str := range a {
        length += len(str)
    }

    s, b := rawstring(length) // 生成指定大小的字符串，返回一个string和切片，二者共享内存空间

    for _, str := range a {
        copy(b, str)    // string无法修改，只能通过切片修改
        b = b[len(str):]
    }

    return s
}
```

因为 string 是无法直接修改的，所以这里使用 rawstring() 方法初始化一个指定大小的 string，同时返回一个切片，二者共享同一块内存空间，后面向切片中拷贝数据，也就间接修改了 string。

rawstring() 源代码如下：

```go
func rawstring(size int) (s string, b []byte) { // 生成一个新的string，返回的string和切片共享相同的空间
	p := mallocgc(uintptr(size), nil, false)

	stringStructOf(&s).str = p
	stringStructOf(&s).len = size

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

	return
}
```

## 为什么字符串不允许修改？

像 C++ 语言中的 string，其本身拥有内存空间，修改 string 是支持的。但 Go 的实现中，string 不包含内存空间，只有一个内存的指针，这样做的好处是 string 变得非常轻量，可以很方便的进行传递而不用担心内存拷贝。

因为 string 通常指向字符串字面量，而字符串字面量存储位置是只读段，而不是堆或栈上，所以才有了 string 不可修改的约定。

## []byte 转换成 string 一定会拷贝内存吗？

byte 切片转换成 string 的场景很多，为了性能上的考虑，有时候只是临时需要字符串的场景下，byte 切片转换成 string 时并不会拷贝内存，而是直接返回一个 string，这个 string 的指针 (string.str) 指向切片的内存。

比如，编译器会识别如下临时场景：

* 使用 m[string(b)] 来查找 map （ map 是 string 为 key，临时把切片 b 转成 string ）；
* 字符串拼接： `"<" + "string(b)" + ">"`；
* 字符串比较：`string(b) == "foo"`

因为是临时把 byte 切片转换成 string，也就避免了因 byte 切片同容改成而导致 string 引用失败的情况，所以此时可以不必拷贝内存新建一个 string。

## string 和 []byte 如何取舍

string 和 []byte 都可以表示字符串，但因数据结构不同，其衍生出来的方法也不同，要根据实际应用场景来选择。

string 擅长的场景：

* 需要字符串比较的场景；
* 不需要 nil 字符串的场景；

[]byte 擅长的场景：

* 修改字符串的场景，尤其是修改粒度为 1 个字节；
* 函数返回值，需要用 nil 表示含义的场景；
* 需要切片操作的场景；

虽然看起来 string 适用的场景不如 []byte 多，但因为 string 直观，在实际应用中还是大量存在，在偏底层的实现中 []byte 使用更多。

```go

```
