---
date: 2022-10-02T19:43:58+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "字符串底层原理"  # 文章标题
url:  "posts/go/docs/internal/string/underlying_principle"  # 设置网页永久链接
tags: [ "Go", "underlying-principle" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 字符串解析

### 声明字符串

整数是全为数字的常量，浮点数是带小数点的常量。字符串也有特殊标识，它有两种声明方式：

```go
var a string = "hello, world"
var b string = `hello, world`
```

字符串构建过程是先根据字符串构建 stringStruct，再转换成 string。转换的源码如下：

```go
// 根据字符串地址构建 string
func gostringnocopy(str *byte) string {
    // 先构造 stringStruct
	ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)}
    // 再将 stringStruct 转换成 string
	s := *(*string)(unsafe.Pointer(&ss))
	return s
}
```

string 在 runtime 包中就是 stringStruct，对外呈现叫做 string。

### 编译阶段

字符串常量在词法解析阶段最终会被标记成 StringLit 类型的 Token 并被传递到编译的下一个阶段。在语法分析阶段，采取递归下降的方式读取 Uft-8 字符，单撇号或双引号是字符串的标识。分析的逻辑位于 syntax/scanner.go 文件中。

```go
	case '"':
		s.stdString()

	case '`':
		s.rawString()
```

如果在代码中识别到单撇号，则调用 rawString 函数；如果识别到双引号，则调用 stdString 函数，两者的处理略有不同。

对于单撇号的处理比较简单：一直循环向后读取，直到寻找到配对的单撇号，如下所示。

```go
func (s *scanner) rawString() {
	ok := true
	s.nextch()

	for {
		if s.ch == '`' {
			s.nextch()
			break
		}
		if s.ch < 0 {
			s.errorAtf(0, "string not terminated")
			ok = false
			break
		}
		s.nextch()
	}
	// We leave CRs in the string since they are part of the
	// literal (even though they are not part of the literal
	// value).

	s.setLit(StringLit, ok)
}
```

双引号调用 stdString 函数，如果出现另一个双引号则直接退出；如果出现了 `\\`，则对后面的字符进行转义。

```go
func (s *scanner) stdString() {
	ok := true
	s.nextch()

	for {
		if s.ch == '"' {
			s.nextch()
			break
		}
		if s.ch == '\\' {
			s.nextch()
			if !s.escape('"') {
				ok = false
			}
			continue
		}
		if s.ch == '\n' {
			s.errorf("newline in string")
			ok = false
			break
		}
		if s.ch < 0 {
			s.errorAtf(0, "string not terminated")
			ok = false
			break
		}
		s.nextch()
	}

	s.setLit(StringLit, ok)
}
```

在双引号中不能出现换行符，以下代码在编译时会报错：`newline in string`。这是通过对每个字符判断 `r=='\\n'` 实现的。

最后将解析到的字节转换为字符串，这种转换会在字符串左、右两边加上双引号，因此 "hello" 会被解析为 ""hello""。在抽象语法树阶段，无论是 import 语句中包的路径、结构体中的字段标签还是字符串常量，都会调用 strconv.Unquote(s) 去掉字符串两边的引号等干扰，还原其本来的面目，例如将 ""hello"" 转换为 "hello"。

```go
func (check *Checker) tag(t *syntax.BasicLit) string {
	// If t.Bad, an error was reported during parsing.
	if t != nil && !t.Bad {
		if t.Kind == syntax.StringLit {
			if val, err := strconv.Unquote(t.Value); err == nil {
				return val
			}
		}
		check.errorf(t, 0, invalidAST+"incorrect tag syntax: %q", t.Value)
	}
	return ""
}
```

## 字符串拼接

在 Go 语言中，可以方便地通过加号操作符（+）对字符串进行拼接。

```go
str := "Str1" + "Str2" + "Str3"
```

很显然，由于数字的加法操作也使用加号操作符，因此需要编译时识别具体为何种操作。

当加号操作符两边是字符串时，编译时抽象语法树阶段具体操作的 Op 会被解析为 OADDSTR。

### 字符串常量拼接

对两个及以上字符串常量进行拼接时会在语法分析阶段调用 noder.sum 函数。例如对于 "a"+"b"+"c" 的场景，noder.sum 函数先将所有的字符串常量放到字符串数组中，然后调用 strings.Join 函数完成对字符串常量数组的拼接。

### 字符串变量拼接

如果涉及字符串变量的拼接，那么其拼接操作最终是在运行时完成的。

## 运行时字符拼接

运行时字符串的拼接原理如图5-1所示，其并不是简单地将一个字符串合并到另一个字符串中，而是找到一个更大的空间，并通过内存复制的形式将字符串复制到其中。

![](../../../assets/images/docs/internal/string/underlying_principle/图5-1%20字符串拼接原理.png)

运行时具体的拼接代码如下，其实无论使用 concatstring{2, 3, 4, 5} 函数中的哪一个，最终都会调用 runtime.concatstrings 函数。

```go
func concatstring2(buf *tmpBuf, a0, a1 string) string {
	return concatstrings(buf, []string{a0, a1})
}

func concatstring3(buf *tmpBuf, a0, a1, a2 string) string {
	return concatstrings(buf, []string{a0, a1, a2})
}

func concatstring4(buf *tmpBuf, a0, a1, a2, a3 string) string {
	return concatstrings(buf, []string{a0, a1, a2, a3})
}

func concatstring5(buf *tmpBuf, a0, a1, a2, a3, a4 string) string {
	return concatstrings(buf, []string{a0, a1, a2, a3, a4})
}
```

concatstrings 函数会先对传入的切片参数进行遍历，过滤空字符串并计算拼接后字符串的长度。

一个拼接语句的字符串编译时都会被存放到一个切片中，拼接过程需要遍历两次切片，第一次遍历获取总的字符串长度，据此申请内存，第二次遍历会把字符串逐个拷贝过去。


`src/runtime/string.go`

```go
// The constant is known to the compiler.
// There is no fundamental theory behind this number.
const tmpStringBufSize = 32

type tmpBuf [tmpStringBufSize]byte

// concatstrings implements a Go string concatenation x+y+z+...
// The operands are passed in the slice a.
// If buf != nil, the compiler has determined that the result does not
// escape the calling function, so the string data can be stored in buf
// if small enough.
func concatstrings(buf *tmpBuf, a []string) string {
	idx := 0
	l := 0
	count := 0
	for i, x := range a {
		n := len(x)
		if n == 0 {
			continue
		}
		if l+n < l {
			throw("string concatenation too long")
		}
		l += n
		count++
		idx = i
	}
	if count == 0 {
		return ""
	}

	// If there is just one string and either it is not on the stack
	// or our result does not escape the calling frame (buf != nil),
	// then we can return that string directly.
	if count == 1 && (buf != nil || !stringDataOnStack(a[idx])) {
		return a[idx]
	}
	s, b := rawstringtmp(buf, l)
	for _, x := range a {
		copy(b, x)
		b = b[len(x):]
	}
	return s
}

// stringDataOnStack reports whether the string's data is
// stored on the current goroutine's stack.
func stringDataOnStack(s string) bool {
	ptr := uintptr(unsafe.Pointer(unsafe.StringData(s)))
	stk := getg().stack
	return stk.lo <= ptr && ptr < stk.hi
}
```

拼接的过程位于 rawstringtmp 函数中，当拼接后的字符串小于 32 字节时，会有一个临时的缓存供其使用。当拼接后的字符串大于 32 字节时，堆区会开辟一个足够大的内存空间，并将多个字符串存入其中，期间会涉及内存的复制（copy）。

```go
func rawstringtmp(buf *tmpBuf, l int) (s string, b []byte) {
	if buf != nil && l <= len(buf) {
		b = buf[:l]
		s = slicebytetostringtmp(&b[0], len(b))
	} else {
		s, b = rawstring(l)
	}
	return
}
```

因为 string 是无法直接修改的，所以这里使用 rawstring() 方法初始化一个指定大小的 string，同时返回一个切片，二者共享同一块内存空间，后面向切片中拷贝数据，也就间接修改了 string。

rawstring() 源代码如下：

```go
// 生成一个新的 string，返回的 string 和切片共享相同的空间
func rawstring(size int) (s string, b []byte) { 
	p := mallocgc(uintptr(size), nil, false)

	stringStructOf(&s).str = p
	stringStructOf(&s).len = size

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

	return
}
```

字符常量的拼接发生在编译时，而字符串变量的拼接发生在运行时。

当拼接后的 s 字符串小于 32 字节时，会有一个临时的缓存供其使用。当拼接后的字符串大于 32 字节时，会请求在堆区分配内存。

```go

```
