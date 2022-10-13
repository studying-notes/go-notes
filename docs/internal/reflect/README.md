---
date: 2022-10-07T16:57:39+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "反射"  # 文章标题
url:  "posts/go/docs/internal/reflect/README"  # 设置网页永久链接
tags: [ "Go", "readme" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 为什么需要反射

在计算机科学中，**反射是程序在运行时检查、修改自身结构和行为的能力**。比如，汇编语言具有固有的反射性，因为它可以通过**将指令定义为数据**并修改这些指令数据对原始体系结构进行修改。

假如现在有一个 Student 结构体：

```go
type Student struct {
	Name string
	Age  int
}
```

如果希望写一个可以将该结构体转换为 sql 语句的函数，按照过去的实现方式，可以为此结构体添加一个如下的生成方法：

```go
func (s *Student) CreateSQL() string {
	return fmt.Sprintf("INSERT INTO student VALUES (%s, %d)", s.Name, s.Age)
}
```

但是，假如我们有其他结构体也有相同的需求呢？很显然，可以为每个类型都添加一个 CreateSQL 方法，并定义一个接口：

```go
type SQL interface {
	CreateSQL() string
}
```

这种方法在项目初期，以及结构体类型简单的时候是比较方便的。但是有时候项目中定义的类型会非常多，而且可能当前类型还没有被创建出来（需要运行时创建或者通过远程过程调用触发），这时我们会写很多逻辑相同的重复代码。

那么是否有一种更加简单通用的办法来解决这一类问题呢？如果有办法在运行时探测到结构体变量中的成员名和方法名就好了，这就是反射为我们提供的功能。

如下所示，可以将上面的场景改造成反射的形式。在 createQuery 函数中，我们可以传递任何的结构体类型，该函数会遍历结构体中所有的字段，并构造 query 字符串。

```go
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func (s *Student) CreateSQL() string {
	return fmt.Sprintf("INSERT INTO student VALUES (%s, %d)", s.Name, s.Age)
}

type SQL interface {
	CreateSQL() string
}

func createQuery(q interface{}) string {
	if reflect.TypeOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		v := reflect.ValueOf(q)
		query := fmt.Sprintf("INSERT INTO %s VALUES (", t)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query += fmt.Sprintf("%d", v.Field(i).Int())
				} else {
					query += fmt.Sprintf(", %d", v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query += fmt.Sprintf("'%s'", v.Field(i).String())
				} else {
					query += fmt.Sprintf(",'%s'", v.Field(i).String())
				}
			}
		}
		query += ")"
		return query
	}
	return ""
}

type Trade struct {
	id    int
	Price int
}

func main() {
	s := Student{Name: "John", Age: 20}
	fmt.Println(s.CreateSQL())

	t := Trade{id: 1, Price: 100}
	fmt.Println(createQuery(t))
}
```

可以看到，上面的案例通过反射大大简化了代码的编写。

## 反射的基本使用方法

[反射的基本使用方法](operation.md)

## 反射底层原理

[反射底层原理](underlying_principle.md)

## 反射性能损耗

我们通过一个简单的测试来看看反射的性能损耗。

```go
package main

import (
	"fmt"
	"reflect"
)

type People struct {
	Age  int
	Name string
}

func NewPeople() *People {
	return &People{
		Age:  18,
		Name: "Karl",
	}
}

func NewPeopleByReflect() interface{} {
	var people People
	t := reflect.TypeOf(people)
	v := reflect.New(t)
	v.Elem().Field(0).SetInt(18)
	v.Elem().Field(1).SetString("Karl")
	return v.Interface()
}

func main() {
	p1 := NewPeople()
	p2 := NewPeopleByReflect().(*People)

	fmt.Println(p1)
	fmt.Println(p2)
}
```

```go
package main

import "testing"

func BenchmarkNewPeople(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPeople()
	}
}

func BenchmarkNewPeopleByReflect(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPeopleByReflect()
	}
}
```

```
BenchmarkNewPeople
BenchmarkNewPeople-12                   1000000000               0.2703 ns/op
      0 B/op           0 allocs/op
BenchmarkNewPeopleByReflect
BenchmarkNewPeopleByReflect-12          12622676                92.02 ns/op
     48 B/op           2 allocs/op
```

## 性能损耗猜测

通过上面的测试，我们可以看到，使用反射创建对象的性能损耗是非常大的，这是为什么呢？

### 结构体成员数量

增加结构体成员数量，反射的性能损耗会更大。

4 个成员的情况：

```
BenchmarkNewPeople
BenchmarkNewPeople-12                   1000000000               0.2574 ns/op
      0 B/op           0 allocs/op
BenchmarkNewPeopleByReflect
BenchmarkNewPeopleByReflect-12           9228043               129.4 ns/op
    128 B/op           2 allocs/op
```

0 个成员的情况：

```
BenchmarkNewPeople
BenchmarkNewPeople-12                   1000000000               0.2628 ns/op
      0 B/op           0 allocs/op
BenchmarkNewPeopleByReflect
BenchmarkNewPeopleByReflect-12          74903094                17.66 ns/op
      0 B/op           0 allocs/op
```

### reflect.TypeOf 和 reflect.New

```go
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
```

```go
// toType converts from a *rtype to a Type that can be returned
// to the client of package reflect. In gc, the only concern is that
// a nil *rtype must be replaced by a nil Type, but in gccgo this
// function takes care of ensuring that multiple *rtype for the same
// type are coalesced into a single Type.
func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}
```

该方法涉及接口转换，接口转换的前提是被转换的接口能够包含转换接口中的方法，这里需要在运行时判断，在接口转换时，使用了运行时 `runtime.convI2I` 函数。所以，这里的性能损耗是比较大的。

下面的测试可能是因为 itab 的缓存，循环调用导致接口转换的速度变快了。

```go
func NewPeopleByReflect() interface{} {
	var people People
	t := reflect.TypeOf(people)
	return t
}
```

```
BenchmarkNewPeople
BenchmarkNewPeople-12                   1000000000               0.2786 ns/op
      0 B/op           0 allocs/op
BenchmarkNewPeopleByReflect
BenchmarkNewPeopleByReflect-12          1000000000               1.077 ns/op
      0 B/op           0 allocs/op
```

```go
func NewPeopleByReflect() interface{} {
	var people People
	t := reflect.TypeOf(people)
	v := reflect.New(t)
	return v
}
```

```
BenchmarkNewPeople
BenchmarkNewPeople-12                   1000000000               0.2589 ns/op
      0 B/op           0 allocs/op
BenchmarkNewPeopleByReflect
BenchmarkNewPeopleByReflect-12           8251340               134.8 ns/op
    152 B/op           3 allocs/op
```

看起来 reflect.TypeOf 的性能损耗是比较小的。那么主要的性能损耗是在 reflect.New 上。

用 Go 原生自带的性能分析工具 pprof 来分析一下它们的主要耗时，来验证猜测。

```
go test -bench=".*" -benchmem -memprofile memprofile.out -cpuprofile profile.out
```

```
go tool pprof ./profile.out
```

```
list NewPeopleByReflect
```

```
      20ms      770ms (flat, cum) 68.14% of Total
         .          .     24:func NewPeopleByReflect() interface{} {
         .          .     25:   var people People
      10ms      280ms     26:   t := reflect.TypeOf(people)
      10ms      300ms     27:   v := reflect.New(t)
         .       30ms     28:   v.Elem().Field(0).SetInt(18)
         .       50ms     29:   v.Elem().Field(1).SetString("Karl")
         .       30ms     30:   v.Elem().Field(2).SetString("Test1")
         .       30ms     31:   v.Elem().Field(3).SetString("Test2")
         .       50ms     32:   return v.Interface()
         .          .     33:}
```

两者的耗时差不多，但上面单独 TypeOf 耗时少，可能是 itab 缓存的原因。

> itab 的构造比较麻烦，因此在 Go 语言中，相同的转换的 itab 会被存储到全局的 hash 表中。当全局的 hash 表中没有相同的 itab 时，会将此 itab 存储到全局的 hash 表中，第二次转换时可以直接到全局的 hash 表中查找此 itab，实现逻辑在 getitab 函数中。
> 详细的逻辑可以参考[接口底层原理](../interface/underlying_principle.md)中的接口转换部分。

reflect.New 的耗时主要是内存的分配，mallocgc 函数涉及 GC 等，比较耗时。

```go
//go:linkname reflect_unsafe_New reflect.unsafe_New
func reflect_unsafe_New(typ *_type) unsafe.Pointer {
	return mallocgc(typ.size, typ, true)
}
```

反射的使用增加了代码的复杂性并在一定程度上降低了效率，一般在实践中使用得不是很多。都在尝试用其他的方式来解决问题，比如通过代码生成的方式来解决反射的问题。

```go

```
