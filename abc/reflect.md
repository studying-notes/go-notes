---
date: 2020-11-08T19:47:48+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 的反射机制"  # 文章标题
url:  "posts/go/abc/reflect"  # 设置网页永久链接
tags: [ "go", "reflect" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

Go 官方博客地址：https://blog.golang.org/laws-of-reflection

## 反射概念

reflect 实现了运行时的反射能力，反射提供一种让程序检查自身结构的能力，反射是一种检查 interface 变量的底层类型和值的机制，能够让程序操作不同类型的对象。

想深入了解反射，必须深入理解类型和接口概念。

## 静态类型

Go是静态类型语言，比如"int"、"float32"、"[]byte"等等。每个变量都有一个静态类型，且在编译时就确定了。考虑一下如下一种类型声明:

```go
type Myint int

var i int
var j Myint
```

i 和 j 类型相同吗？i 和 j 类型是不同的。二者拥有不同的静态类型，没有类型转换的话是不可以互相赋值的，尽管二者底层类型是一样的。

## 特殊的静态类型

interface 类型是一种特殊的类型，它代表方法集合。 它可以存放任何实现了其方法的值。

经常被拿来举例的是 io 包里的这两个接口类型：

```go
// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

任何类型，比如某 struct，只要实现了其中的 Read() 方法就被认为是实现了 Reader 接口，只要实现了 Write() 方法，就被认为是实现了 Writer 接口，不过方法参数和返回值要跟接口声明的一致。

接口类型的变量可以存储任何实现该接口的值。

## 特殊的 interface 类型

最特殊的 interface 类型为空 interface 类型，即 `interface {}`，前面说了，interface 用来表示一组方法集合，所有实现该方法集合的类型都被认为是实现了该接口。那么空 interface 类型的方法集合为空，也就是说所有类型都可以认为是实现了该接口。

一个类型实现空 interface 并不重要，重要的是一个空 interface 类型变量可以存放所有值，记住是所有值，这才是最最重要的。这也是有些人认为 Go 是动态类型的原因，这是个错觉。

## interface 类型是如何表示的

前面讲了，interface 类型的变量可以存放任何实现了该接口的值。还是以上面的 `io.Reader` 为例进行说明，`io.Reader` 是一个接口类型，`os.OpenFile()` 方法返回一个 `File` 结构体类型变量，该结构体类型实现了 `io.Reader` 的方法，那么 `io.Reader` 类型变量就可以用来接收该返回值。如下所示：

```go
var r io.Reader
tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil {
    return nil, err
}
r = tty
```

那么问题来了。

r 的类型是什么？ r 的类型始终是 `io.Reader`interface 类型，无论其存储什么值。

那 `File` 类型体现在哪里？ r 保存了一个 (value, type) 对来表示其所存储值的信息。value 即为 r 所持有元素的值，type 即为所持有元素的底层类型

如何将 r 转换成另一个类型结构体变量？比如转换成 `io.Writer` ？使用类型断言，如 `w = r.(io.Writer)`. 意思是如果 r 所持有的元素如果同样实现了 io.Writer 接口, 那么就把值传递给 w。

## 反射三定律

前面之所以讲类型，是为了引出 interface，之所以讲 interface 是想说 interface 类型有个 (value，type) 对，而反射就是检查 interface 的这个 (value, type) 对的。具体一点说就是 Go 提供一组方法提取 interface 的 value，提供另一组方法提取 interface 的 type. 

官方提供了三条定律来说明反射，比较清晰，下面也按照这三定律来总结。

反射包里有两个接口类型要先了解一下。

- `reflect.Type` 提供一组接口处理 interface 的类型，即（ value, type ）中的 type 
- `reflect.Value` 提供一组接口处理 interface 的值, 即 (value, type) 中的 value 

下面会提到反射对象，所谓反射对象即反射包里提供的两种类型的对象。

- `reflect.Type` 类型对象
- `reflect.Value` 类型对象

### 反射第一定律：反射可以将 interface 类型变量转换成反射对象

下面示例，看看是如何通过反射获取一个变量的值和类型的：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4
	t := reflect.TypeOf(x)  //t is reflect.Type
	fmt.Println("type:", t)

	v := reflect.ValueOf(x) //v is reflect.Value
	fmt.Println("value:", v)
}
```

程序输出如下：

```
type: float64
value: 3.4
```

反射是针对 interface 类型变量的，其中 `TypeOf()` 和 `ValueOf()` 接受的参数都是 `interface{}` 类型的，也即 x 值是被转成了 interface 传入的。

### 反射第二定律：反射可以将反射对象还原成 interface 对象

之所以叫'反射'，反射对象与 interface 对象是可以互相转化的。看以下例子：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4

	v := reflect.ValueOf(x) //v is reflect.Value

	var y float64 = v.Interface().(float64)
	fmt.Println("value:", y)
}
```

对象 x 转换成反射对象 v，v 又通过 Interface() 接口转换成 interface 对象，interface 对象通过.(float64) 类型断言获取 float64 类型的值。

### 反射第三定律：反射对象可修改，value 值必须是可设置的

通过反射可以将 interface 类型变量转换成反射对象，可以使用该反射对象设置其持有的值。在介绍何谓反射对象可修改前，先看一下失败的例子：

```go
package main

import (
	"reflect"
)

func main() {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	v.SetFloat(7.1) // Error: will panic.
}
```

如下代码，通过反射对象 v 设置新值，会出现 panic。报错如下：

```
panic: reflect: reflect.Value.SetFloat using unaddressable value
```

错误原因即是 v 是不可修改的。

反射对象是否可修改取决于其所存储的值，回想一下函数传参时是传值还是传址就不难理解上例中为何失败了。

上例中，传入 reflect.ValueOf() 函数的其实是 x 的值，而非 x 本身。即通过 v 修改其值是无法影响 x 的，也即是无效的修改，所以 golang 会报错。

想到此处，即可明白，如果构建 v 时使用 x 的地址就可实现修改了，但此时 v 代表的是指针地址，我们要设置的是指针所指向的内容，也即我们想要修改的是 ` * v`。那怎么通过 v 修改 x 的值呢？

`reflect.Value` 提供了 `Elem()` 方法，可以获得指针向指向的 `value`。看如下代码：

```go
package main

import (
"reflect"
	"fmt"
)

func main() {
	var x float64 = 3.4
	v := reflect.ValueOf(&x)
	v.Elem().SetFloat(7.1)
	fmt.Println("x :", v.Elem().Interface())
}
```

输出为：

```
x : 7.1
```

## 反射函数和类型

- `reflect.TypeOf` 获取类型信息，返回 `Type` 类型；
- `reflect.ValueOf` 获取数据的运行时表示，返回 `Value` 类型。

### Type

类型 `Type` 是反射包定义的一个接口，我们可以使用 reflect.TypeOf 函数获取任意变量的的类型，`Type` 接口中定义了多种方法，比如 `MethodByName` 可以获取当前类型对应方法的引用、`Implements` 可以判断当前类型是否实现了某个接口：

```go
type Type interface {
        Align() int
        FieldAlign() int
        Method(int) Method
        MethodByName(string) (Method, bool)
        NumMethod() int
        // ...
        Implements(u Type) bool
        // ...
}
```

### Value

`Value` 的类型与 `Type` 不同，它被声明成了结构体。这个结构体没有对外暴露的字段，但是提供了获取或者写入数据的方法：


```go
type Value struct {
        // ...
}

func (v Value) Addr() Value
func (v Value) Bool() bool
func (v Value) Bytes() []byte
// ...
```

反射中的所有方法基本都是围绕着 `Type` 和 `Value` 这两个类型设计的。通过 `reflect.TypeOf`、`reflect.ValueOf` 可以将一个普通的变量转换成 `Type` 和 `Value`，随后就可以使用相关方法对它们进行复杂的操作。

## 遍历结构体

### 获取标签值

提取非嵌套结构体指定标签的值

```go
func ExtractTagValue(i interface{}, tag string) (tagValues []string) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && i != nil {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		types := v.Type()
		for i := 0; i < v.NumField(); i++ {
			tagValues = append(tagValues, types.Field(i).Tag.Get(tag))
		}
	}
	return tagValues
}
```

```go
type Fruit struct {
	ID    string   `json:"id"`
	Name  []string `json:"name"`
	Price string   `json:"price"`
	Area  `json:"area"`
}

func main() {
	fmt.Println(ExtractTagValue(Fruit{}, "json"))
}
```

```
[id name price area]
```

### 获取字段值

获取结构体字段的值

```go
func ExtractFieldValue(i interface{}) (fieldValues []interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && i != nil {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fieldValues = append(fieldValues, f.Interface())
		}
	}
	return fieldValues
}
```

```go
func main() {
	fruit := Fruit{ID: "1", Name: []string{"apple", "nut"}, Price: "12",
		Area: Area{Length: "20", Width: "30"}}
	fmt.Println(ExtractFieldValue(fruit))
}
```

```
[1 [apple nut] 12 {20 30}]
```

### 修改字段值

修改结构体字段值

```go
func ModifyFieldValue(ptr interface{}, handle func(string) string) {
	types := reflect.TypeOf(ptr)
	values := reflect.ValueOf(ptr)
	if types.Kind() != reflect.Ptr {
		return // 必须传入指针才能修改原结构体
	}
	types = types.Elem()
	values = values.Elem()
	if values.Kind() == reflect.Struct {
		for i := 0; i < types.NumField(); i++ {
			f := values.Field(i)
			switch f.Kind() {
			case reflect.String: // 字符串类型
				f.Set(reflect.ValueOf(handle(f.String()))) // 设置新字段值
			case reflect.Slice:
				obj := reflect.ValueOf(f.Interface())
				for j := 0; j < obj.Len(); j++ {
					_ = obj.Index(j).String() // 提取切片数据
				}
			}
		}
	}
}
```

```go
func main() {
	ModifyFieldValue(&fruit, func(s string) string {
		return "modify-" + s
	})
	fmt.Printf("%+v\n", fruit)
}
```

```
{ID:modify-1 Name:[apple nut] Price:modify-12 Area:{Length:20 Width:30}}
```

```go

```

```go

```

```go

```
