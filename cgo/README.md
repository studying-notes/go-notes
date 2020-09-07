# CGO 编程笔记目录

CGO 的原理其实就是由编译器识别出 `import "C"` 的位置，然后在其上的注释中提取 C 代码，最后调用 C 编译器进行分开编译。

## CGO 功能预览

- [CGO 功能预览](quick.md#cgo-功能预览)
	- [基于 C 标准库函数输出字符串](quick.md#基于-c-标准库函数输出字符串)
		- [手动释放资源](quick.md#手动释放资源)
	- [用自己定义的 C 函数](quick.md#用自己定义的-c-函数)
		- [独立的 C 语言源文件](quick.md#独立的-c-语言源文件)
	- [C 代码的模块化](quick.md#c-代码的模块化)
	- [用 Go 重新实现 C 函数](quick.md#用-go-重新实现-c-函数)
	- [面向 C 接口的 Go 编程](quick.md#面向-c-接口的-go-编程)

## CGO 入门

- [CGO 入门](start.md#cgo-入门)
  - [`import "C"` 语句](start.md#import-c-语句)
    - [不同包的引用问题](start.md#不同包的引用问题)
  - [`#cgo` 语句](start.md#cgo-语句)
    - [编译和链接参数](start.md#编译和链接参数)
    - [不同平台间的差异](start.md#不同平台间的差异)
  - [pkg-config](start.md#pkg-config)
  - [build tag 条件编译](start.md#build-tag-条件编译)

## 类型转换

- [类型转换](type.md#类型转换)
  - [C 与 Go 的数据类型](type.md#c-与-go-的数据类型)
    - [数值类型](type.md#数值类型)
      - [Go 数值转 C 数值](type.md#go-数值转-c-数值)
      - [C 数值转 Go 数值](type.md#c-数值转-go-数值)
    - [指针类型](type.md#指针类型)
      - [Go 字符串和切片的 C 语言表示](type.md#go-字符串和切片的-c-语言表示)
  - [结构体、联合和枚举类型](type.md#结构体联合和枚举类型)
  - [字符串、数组和切片](type.md#字符串数组和切片)
    - [字符串双向转换](type.md#字符串双向转换)
    - [数组双向转换](type.md#数组双向转换)
  - [指针间的转换](type.md#指针间的转换)
  - [数值和指针的转换](type.md#数值和指针的转换)
  - [切片间的转换](type.md#切片间的转换)

## 函数调用

- [函数调用](func.md#函数调用)
  - [Go 调用 C 函数](func.md#go-调用-c-函数)
  - [C 函数的返回值](func.md#c-函数的返回值)
    - [数值型返回值](func.md#数值型返回值)
    - [获取错误状态码](func.md#获取错误状态码)
    - [void 型返回值](func.md#void-型返回值)
    - [字符串型返回值](func.md#字符串型返回值)
    - [字符串型返回值及其长度](func.md#字符串型返回值及其长度)
    - [结构体型返回值](func.md#结构体型返回值)
  - [C 调用 Go 导出函数](func.md#c-调用-go-导出函数)

## Go 调用 DLL 动态库

- [Go 调用 DLL 动态库](dll.md#go-调用-dll-动态库)
	- [简单 CGO 程序](dll.md#简单-cgo-程序)
	- [Go 程序编译成动态库](dll.md#go-程序编译成动态库)
		- [Windows](dll.md#windows)
		- [Linux/Unix/macOS](dll.md#linuxunixmacos)
	- [Go 调用 Windows DLL](dll.md#go-调用-windows-dll)
		- [syscall.Syscall 系列方法](dll.md#syscallsyscall-系列方法)
		- [调用 user32.dll 试验](dll.md#调用-user32dll-试验)
	- [调用自己编译的 DLL](dll.md#调用自己编译的-dll)
		- [准备文件](dll.md#准备文件)
		- [编译成 DLL 文件](dll.md#编译成-dll-文件)
		- [在 Go 中载入 DLL](dll.md#在-go-中载入-dll)

## Go 程序链接 C 库

- [Go 程序链接 C 库](link.md#go-程序链接-c-库)
  - [链接 C 静态库](link.md#链接-c-静态库)
    - [Linux](link.md#linux)
    - [Windows .lib 文件](link.md#windows-lib-文件)
  - [链接 C 动态库](link.md#链接-c-动态库)
    - [Linux](link.md#linux-1)
    - [Windows](link.md#windows)
  - [Go 导出 C 静态库](link.md#go-导出-c-静态库)
  - [Go 导出 C 动态库](link.md#go-导出-c-动态库)

## CGO 内部机制

- [CGO 内部机制](internal.md#cgo-内部机制)
  - [CGO 生成的中间文件](internal.md#cgo-生成的中间文件)
  - [Go 调用 C 函数](internal.md#go-调用-c-函数)
  - [C 调用 Go 函数](internal.md#c-调用-go-函数)
