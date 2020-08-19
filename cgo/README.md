# CGO 编程

CGO 的原理其实就是由编译器识别出 `import "C"` 的位置，然后在其上的注释中提取 C 代码，最后调用 C 编译器进行分开编译。

## 快速开始

- [快速开始](quick/README.md)
	- [基于 C 标准库函数输出字符串](quick/README.md#基于-c-标准库函数输出字符串)
		- [手动释放资源](quick/README.md#手动释放资源)
	- [用自己定义的 C 函数](quick/README.md#用自己定义的-c-函数)
		- [额外的 C 语言源文件](quick/README.md#额外的-c-语言源文件)
	- [C 代码的模块化](quick/README.md#c-代码的模块化)
	- [用 Go 重新实现 C 函数](quick/README.md#用-go-重新实现-c-函数)
	- [面向 C 接口的 Go 编程](quick/README.md#面向-c-接口的-go-编程)

## CGO 入门

- [CGO 入门](base/README.md)
  - [`import "C"` 语句](base/README.md#import-c-语句)
  - [`#cgo` 语句](base/README.md#cgo-语句)
  - [build tag 条件编译](base/README.md#build-tag-条件编译)

## 类型转换

- [类型转换](type/README.md)
  - [C 与 Go 的数据类型](type/README.md#c-与-go-的数据类型)
    - [数值类型](type/README.md#数值类型)
      - [C 数值类型的 Go 表示](type/README.md#c-数值类型的-go-表示)
      - [Go 数值类型的 C 表示](type/README.md#go-数值类型的-c-表示)
    - [指针类型](type/README.md#指针类型)
      - [Go 字符串和切片的 C 语言表示](type/README.md#go-字符串和切片的-c-语言表示)
  - [结构体、联合和枚举类型](type/README.md#结构体联合和枚举类型)
  - [字符串、数组和切片](type/README.md#字符串数组和切片)
    - [字符串双向转换](type/README.md#字符串双向转换)
    - [数组双向转换](type/README.md#数组双向转换)
  - [指针间的转换](type/README.md#指针间的转换)
  - [数值和指针的转换](type/README.md#数值和指针的转换)
  - [切片间的转换](type/README.md#切片间的转换)

# 函数调用

- [函数调用](func/README.md)
  - [Go 调用 C 函数](func/README.md#go-调用-c-函数)
  - [C 函数的返回值](func/README.md#c-函数的返回值)
    - [数值型返回值](func/README.md#数值型返回值)
    - [错误返回值](func/README.md#错误返回值)
    - [void 返回值](func/README.md#void-返回值)
    - [字符串型返回值](func/README.md#字符串型返回值)
    - [字符串型返回值及其长度](func/README.md#字符串型返回值及其长度)
    - [结构体型返回值](func/README.md#结构体型返回值)
  - [C 调用 Go 导出函数](func/README.md#c-调用-go-导出函数)

# Go 调用 DLL

- [Go 调用 DLL](dll/README.md)
	- [简单 CGO 程序](dll/README.md#简单-cgo-程序)
	- [Go 程序编译成动态库](dll/README.md#go-程序编译成动态库)
		- [Windows](dll/README.md#windows)
		- [Linux/Unix/macOS](dll/README.md#linuxunixmacos)
	- [Go 调用 Windows DLL](dll/README.md#go-调用-windows-dll)
		- [syscall.Syscall 系列方法](dll/README.md#syscallsyscall-系列方法)
		- [调用 user32.dll 试验](dll/README.md#调用-user32dll-试验)
	- [调用自己编译的 DLL](dll/README.md#调用自己编译的-dll)
		- [准备文件](dll/README.md#准备文件)
		- [编译成 DLL 文件](dll/README.md#编译成-dll-文件)
		- [在 Go 中载入 DLL](dll/README.md#在-go-中载入-dll)

# Go 程序链接 C 库

- [Go 程序链接 C 库](link/README.md)
  - [链接 C 静态库](link/README.md#链接-c-静态库)
  - [链接 C 动态库](link/README.md#链接-c-动态库)
    - [Linux](link/README.md#linux)
    - [Windows](link/README.md#windows)
  - [Go 导出 C 静态库](link/README.md#go-导出-c-静态库)
  - [Go 导出 C 动态库](link/README.md#go-导出-c-动态库)
