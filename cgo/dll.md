---
date: 2020-08-07T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "CGO 调用 DLL 动态库"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/cgo/dll"  # 设置网页永久链接
tags: [ "go", "cgo" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 最简 CGO 程序

```go
// test.go
package main

import "C"

func hello(msg string) string {
	return "Hello! " + msg
}

func main() {
}
```

## Go 程序编译成动态库

### Windows

```shell
go build -buildmode=c-shared -o test.dll test.go
```

### Linux/Unix/macOS

```shell
go build -buildmode=c-shared -o test.so test.go
```

## Go 调用 Windows DLL

### syscall.Syscall 系列方法

```go
syscall.Syscall
syscall.Syscall6
syscall.Syscall9
syscall.Syscall12
syscall.Syscall15
```

分别对应 3/6/9/12/15 个参数或以下的调用。参数都形如：

```go
syscall.Syscall(trap, nargs, a1, a2, a3)
```

第二个参数, `nargs` 即参数的个数，一旦传错，轻则调用失败，重者直接 `APPCARSH`，多余的参数，用 0 代替，跟 `Syscall` 系列一样，`Call` 方法最多 15 个参数。`Must` 开头的方法，如不存在，会 `panic`。

```go
// 根据函数名字加载 DLL
func LoadDLL(name string) (*DLL, error)
func MustLoadDLL(name string) *DLL

// 表示 DLL 抽象
type DLL struct {
    // ...
}

// 根据函数名字加载 DLL 中的函数
func (d *DLL) FindProc(name string) (proc *Proc, err error)
func (d *DLL) MustFindProc(name string) *Proc

// 表示 DLL 中的函数
type Proc struct {
    // ...
}

// Call 方法调用 syscall.Proc 表示的函数，参数为 uintptr 切片
// error 表示错误，该错误永远不为 nil，它是由 GetLastError() 构成的
func (p *Proc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error)

// 将字符串转换为 byte 指针
func StringBytePtr(s string) *byte
```

### 调用 user32.dll 试验

```go
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	YesNoCancel = 0x00000003
)

func main() {
	fmt.Println("---------Start---------")

	ShowMessageA()
	ShowMessageB()
	ShowMessageC()
	ShowMessageD()

	fmt.Println("---------Stop---------")
}

// The first DLL method call
func ShowMessageA() int {
	user32, _ := syscall.LoadLibrary("user32.dll")
	msgBox, _ := syscall.GetProcAddress(user32, "MessageBoxW")
	defer syscall.FreeLibrary(user32)

	ret, _, err := syscall.Syscall9(msgBox, 4, 0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The first DLL method call"))),
		YesNoCancel, 0, 0, 0, 0, 0,
	)
	if err != 0 {
		panic(err.Error())
	}
	return int(ret)
}

// The second DLL method call
func ShowMessageB() {
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The second DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}

// The third DLL method call
func ShowMessageC() {
	user32, _ := syscall.LoadDLL("user32.dll")
	MessageBoxW, _ := user32.FindProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The third DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}

// The fourth DLL method call
func ShowMessageD() {
	user32 := syscall.MustLoadDLL("user32.dll")
	MessageBoxW := user32.MustFindProc("MessageBoxW")
	MessageBoxW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("The fourth DLL method call"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("CGO DLL"))),
		uintptr(0),
	)
}
```

## 调用自己编译的 DLL

### 准备文件

```c++
// test.h
#ifndef TEST_H
#define TEST_H

#ifdef TEST_DLL_EXPORT
#define TEST_API __declspec(dllexport)
#else
#define TEST_API __declspec(dllimport)
#endif

TEST_API void greet(char *n);

#endif
```

```c
// test.c
#define TEST_DLL_EXPORT

#include "test.h"
#include <stdio.h>

void greet(char *n) {
    printf("Hello, %s!\n", n);
}

// gcc test.c -shared -o test.dll
```

### 编译成 DLL 文件

```shell
gcc test.c -shared -o test.dll
```

### 在 Go 中载入 DLL

DLL 文件和 Go 源文件在同一文件夹下即可。

```go
package main

import (
    "syscall"
    "unsafe"
)

func main() {
    dll := syscall.MustLoadDLL("test.dll")

    procGreet := dll.MustFindProc("greet")
    _, _, _ = procGreet.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr("World"))))
}
```

从 C 中返回的指针是由 malloc 动态分配的，Go 中不会对此指针进行引用计数，不会被垃圾回收，因此会造成内存泄漏。解决方法是在 DLL 中提供释放资源的接口，在 Go 中调用此接口。
