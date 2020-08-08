# CGO 编程

CGO 的原理其实就是由编译器识别出 `import "C"` 的位置，然后在其上的注释中提取 C 代码，最后调用 C 编译器进行分开编译。

- [快速开始](quick/README.md)
	- [基于 C 标准库函数输出字符串](quick/README.md#基于-c-标准库函数输出字符串)
		- [手动释放资源](quick/README.md#手动释放资源)
	- [用自己定义的 C 函数](quick/README.md#用自己定义的-c-函数)
		- [额外的 C 语言源文件](quick/README.md#额外的-c-语言源文件)
	- [C 代码的模块化](quick/README.md#c-代码的模块化)
	- [用 Go 重新实现 C 函数](quick/README.md#用-go-重新实现-c-函数)
	- [面向 C 接口的 Go 编程](quick/README.md#面向-c-接口的-go-编程)

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
