---
date: 2022-10-08T09:44:47+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "反射的基本使用方法"  # 文章标题
url:  "posts/go/docs/internal/reflect/operation"  # 设置网页永久链接
tags: [ "Go", "operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [反射的两种基本类型](#反射的两种基本类型)
- [反射转换为接口](#反射转换为接口)
- [Elem() 间接访问](#elem-间接访问)
- [修改反射的值](#修改反射的值)
- [结构体与反射](#结构体与反射)
- [遍历结构体字段](#遍历结构体字段)
- [修改结构体字段](#修改结构体字段)
- [嵌套结构体的赋值](#嵌套结构体的赋值)
- [结构体方法与动态调用](#结构体方法与动态调用)
- [反射在运行时创建结构体](#反射在运行时创建结构体)
- [函数与反射](#函数与反射)
- [反射与其他类型](#反射与其他类型)

## 反射的两种基本类型

Go语言中提供了两种基本方法可以让我们构建反射的两个基本类型：`reflect.Type` 和 `reflect.Value`。

```go
func ValueOf(i interface{}) Value
func TypeOf(i interface{}) Type
```

这两个函数的参数都是空接口 interface{}，内部存储了即将被反射的变量。因此，反射与接口之间存在很强的联系。可以说，不理解接口就无法深入理解反射。

可以将 reflect.Value 看作反射的值，reflect.Type 看作反射的实际类型。其中，reflect.Type 是一个接口，包含和类型有关的许多方法签名，例如 Align 方法、String 方法等。

reflect.Value 是一个结构体，其内部包含了很多方法。可以简单地用 fmt 打印 reflect.TypeOf 与 reflect.ValueOf 函数生成的结果。reflect.ValueOf 将打印出反射内部的值，reflect.TypeOf 会打印出反射的类型。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var n = 1.23
	fmt.Println("reflect.TypeOf(n) =", reflect.TypeOf(n))
	fmt.Println("reflect.ValueOf(n) =", reflect.ValueOf(n))
}
```

```go
reflect.TypeOf(n) = float64
reflect.ValueOf(n) = 1.23
```

reflect.Value 类型中的 Type 方法可以获取当前反射的类型。

因此，reflect.Value 可以转换为 reflect.Type。reflect.Value 与 reflect.Type 都具有 Kind 方法，可以获取标识类型的 Kind，其底层是 unit。

Go 语言中的内置类型都可以用唯一的整数进行标识。

```go
// A Kind represents the specific kind of type that a Type represents.
// The zero Kind is not a valid kind.
type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Pointer
	Slice
	String
	Struct
	UnsafePointer
)
```

如下所示，通过 Kind 类型可以方便地验证反射的类型是否相同。

```go
package main

import (
	"reflect"
)

func main() {
	var n = 1.23
	println(reflect.TypeOf(n).Kind() == reflect.Float64)
}
```

## 反射转换为接口

reflect.Value 中的 Interface 方法以空接口的形式返回 reflect.Value 中的值。如果要进一步获取空接口的真实值，可以通过接口的断言语法对接口进行转换。

下例实现了从值到反射，再从反射到值的过程。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var n = 1.23
	pointer := reflect.ValueOf(&n)
	value := reflect.ValueOf(n)
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(*convertPointer)
	fmt.Println(convertValue)
}
```

除了使用接口进行转换，reflect.Value 还提供了一些转换到具体类型的方法，这些特殊的方法可以加快转换的速度。

```go
func (v Value) Int() int64
func (v Value) Uint() uint64
func (v Value) Float() float64
func (v Value) String() string
```

另外，这些方法经过了特殊的处理，因此不管反射内部类型是 int8、int16，还是 int32，通过 Int 方法后都将转换为 int64。

```go
// Int returns v's underlying value, as an int64.
// It panics if v's Kind is not Int, Int8, Int16, Int32, or Int64.
func (v Value) Int() int64 {
	k := v.kind()
	p := v.ptr
	switch k {
	case Int:
		return int64(*(*int)(p))
	case Int8:
		return int64(*(*int8)(p))
	case Int16:
		return int64(*(*int16)(p))
	case Int32:
		return int64(*(*int32)(p))
	case Int64:
		return *(*int64)(p)
	}
	panic(&ValueError{"reflect.Value.Int", v.kind()})
}
```

## Elem() 间接访问

如果反射中存储的是指针或接口，那么如何访问指针指向的数据呢？reflect.Value 提供了 Elem 方法返回指针或接口指向的数据。如果 Value 存储的不是指针或接口，则使用 Elem 方法时会出错。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var n = 1.23
	pointer := reflect.ValueOf(&n)
	value := pointer.Elem()
	fmt.Println(value)
}
```

当涉及修改反射的值时，Elem 方法是非常必要的。

我们已经知道，接口中存储的是指针，那么我们要修改的究竟是指针本身还是指针指向的数据呢？这个时候 Elem 方法就起到了关键作用。为了更好地理解 Elem 方法的功能，下面举一个特殊的例子——反射类型是一个空接口，而空接口中包含了 int 类型的指针。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var n = 123
	var y = &n
	var x interface{} = y

	v := reflect.ValueOf(&x) // v is a reflect.Value of type *interface{}
	fmt.Println(v.Type())

	vx := v.Elem() // vx is a reflect.Value of type interface{}
	fmt.Println(vx.Kind())

	vy := vx.Elem() // vy is a reflect.Value of type *int
	fmt.Println(vy.Kind())

	vz := vy.Elem() // vz is a reflect.Value of type int
	fmt.Println(vz.Kind())
}
```

```
*interface {}
interface
ptr
int
```

后面还会看到，在修改反射值时也需要使用到 Elem 方法。

reflect.Type 类型仍然有 Elem 方法，但是该方法只用于获取类型。该方法不仅仅可以返回指针和接口指向的类型，还可以返回数组、通道、切片、指针、哈希表存储的类型。下面用一个复杂的例子来说明该方法的功能，如果反射的类型在这些类型之外，那么仍然会报错。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	type A [16]int16
	var c <-chan map[A][]byte

	tc := reflect.TypeOf(c)
	fmt.Println(tc.Kind())    // chan
	fmt.Println(tc.ChanDir()) // <-chan

	tm := tc.Elem()
	ta, tb := tm.Key(), tm.Elem()
	fmt.Println(tm.Kind(), ta.Kind(), tb.Kind()) // map array slice

	tx, ty := ta.Elem(), tb.Elem()
	// byte is an alias of uint8
	fmt.Println(tx.Kind(), ty.Kind()) // int16 uint8
}
```

## 修改反射的值

有多种方式可以修改反射中存储的值，例如 reflect.Value 的 Set 方法：

```go
func (v Value) Set(x Value)
```

Set 方法接收一个 reflect.Value 类型的参数，该参数的类型必须和 v 的类型一致，否则会报错。如果 v 是一个指针，那么 x 的类型必须和 v 指向的类型一致。如果 v 是一个接口，那么 x 的类型必须和 v 中存储的类型一致。

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var n = 123
    var x interface{} = &n

    v := reflect.ValueOf(&x)
    vx := v.Elem()
    vy := vx.Elem()
    vz := vy.Elem()

    vz.SetInt(456)
    fmt.Println(n)
}
```

只有当反射中存储的实际值是指针时才能赋值，否则是没有意义的，因为在反射之前，实际值被转换为了空接口，如果空接口中存储的值是一个副本，那么修改它会引起混淆，因此Go语言禁止这样做。

必须使用 Elem 方法才能够让值可以被赋值。可以通过 Elem 方法区分要修改的是指针还是指针指向的数据。

## 结构体与反射

假设现在有 User 结构体及 ReflectCallFunc 方法：

```go
type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFunc() {
	fmt.Println("ReflectCallFunc")
}
```

下例通过反射的两种基本方法将结构体转换为反射类型，用 fmt 简单打印出类型与值：

```go
func main() {
	user := User{1, "Tom", 12}
	rType := reflect.TypeOf(user)
	fmt.Println("rType:", rType)
	rValue := reflect.ValueOf(user)
	fmt.Println("rValue:", rValue)
}
```

## 遍历结构体字段

```go
func main() {
	user := User{1, "Tom", 12}
	rType := reflect.TypeOf(user)
	rValue := reflect.ValueOf(user)
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		value := rValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}
}
```

通过 reflect.Type 类型的 NumField 函数获取结构体中字段的个数。relect.Type 与 reflect.Value 都有 Field 方法，relect.Type 的 Field 方法主要用于获取结构体的元信息，其返回 StructField 结构，该结构包含字段名、所在包名、Tag 名等基础信息。

```go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string

	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}
```

reflect.Value 的 Field 方法主要返回结构体字段的值类型，后续可以使用它修改结构体字段的值。

## 修改结构体字段

以下例中最简单的结构体 s 为例，X 为大写，表示可以导出的字段，而 y 为小写，表示未导出的字段。

```go
type User struct {
	X int
	y float64
}
```

要修改结构体字段，可以使用 reflect.Value 提供的 Set 方法。初学者可能选择使用如下方式进行赋值操作，但这种方式是错误的。

```go
func main() {
	var s = User{X: 1, y: 2.0}
	rValue := reflect.ValueOf(s)
	rValueX := rValue.Field(0)
	rValueX.SetInt(100)
}
```

错误的原因正如我们在介绍 Elem 方法时提到的，由于 reflect.ValueOf 函数的参数是空接口，如果我们将值类型复制到空接口会产生一次复制，那么值就不是原来的值了，因此 Go 语言禁止了这种容易带来混淆的写法。要想修改原始值，需要在构造反射时传递结构体指针。

但是只修改为指针还不够，因为在 Field 方法中调用的方法必须为结构体。

因此，需要先通过 Elem 方法获取指针指向的结构体值类型，才能调用 field 方法。正确的使用方式如下所示。同时要注意，未导出的字段 y 是不能被赋值的。

```go
func main() {
	var s = User{X: 1, y: 2.0}
	rValue := reflect.ValueOf(&s).Elem()
	rValueX := rValue.Field(0)
	rValueX.SetInt(100)
}
```

```go
// SetInt sets v's underlying value to x.
// It panics if v's Kind is not Int, Int8, Int16, Int32, or Int64, or if CanSet() is false.
func (v Value) SetInt(x int64) {
	v.mustBeAssignable()
	switch k := v.kind(); k {
	default:
		panic(&ValueError{"reflect.Value.SetInt", v.kind()})
	case Int:
		*(*int)(v.ptr) = int(x)
	case Int8:
		*(*int8)(v.ptr) = int8(x)
	case Int16:
		*(*int16)(v.ptr) = int16(x)
	case Int32:
		*(*int32)(v.ptr) = int32(x)
	case Int64:
		*(*int64)(v.ptr) = x
	}
}
```

## 嵌套结构体的赋值

下例中，Nested 结构体中包含了 Child 字段，Child 也是一个结构体，那么 Child 字段的值能被修改吗？能够被修改的前提是 Child 字段对应的 children 结构体的所有字段都是可导出的。

```go
package main

import "reflect"

type children struct {
	Age int
}

type Nested struct {
	X     int
	Child children
}

func main() {
	vs := reflect.ValueOf(&Nested{X: 1, Child: children{Age: 2}}).Elem()
	vz := vs.Field(1)
	vz.Set(reflect.ValueOf(children{Age: 3}))
}
```

## 结构体方法与动态调用

要获取任意类型对应的方法，可以使用 reflect.Type 提供的 Method 方法，Method 方法需要传递方法的 index 序号。

```go
Method(i int) Method
```

如果 index 序号超出了范围，则会在运行时报错。该方法在大部分时候如下例所示，用于遍历反射结构体的方法。

```go
func main() {
    var s = User{X: 1, y: 2.0}
    rType := reflect.TypeOf(s)
    for i := 0; i < rType.NumMethod(); i++ {
        method := rType.Method(i)
        fmt.Println(method.Name)
    }
}
```

在实践中，更多时候我们使用 reflect.Value 的 MethodByName 方法，参数为方法名并返回代表该方法的 reflect.Value 对象。如果该方法不存在，则会返回空。

如下所示，通过 Type 方法将 reflect.Value 转换为 reflect.Type，reflect.Type 接口中有一系列方法可以获取函数的参数个数、返回值个数、方法个数等属性。

```go
func main() {
    var s = User{X: 1, y: 2.0}
    rValue := reflect.ValueOf(s)
    rType := rValue.Type()
    fmt.Println(rType.NumMethod())
    fmt.Println(rType.NumField())
    fmt.Println(rType.NumIn())
    fmt.Println(rType.NumOut())
}
```

获取代表方法的 reflectv.Value 对象后，可以通过 call 方法在运行时调用方法。

```go
func (v Value) Call(in []Value) []Value
```

Call 方法的参数为实际方法中传入参数的 reflect.Value 切片。因此，对于无参数的调用，参数需要构造一个长度为 0 的切片，如下所示。

```go
func main() {
    var s = User{X: 1, y: 2.0}
    rValue := reflect.ValueOf(s)
    method := rValue.MethodByName("Get")
    method.Call([]reflect.Value{})
}
```

对于有参数的调用，需要先构造出 reflect.Value 类型的参数切片，如下所示。

```go
func main() {
    var s = User{X: 1, y: 2.0}
    rValue := reflect.ValueOf(s)
    method := rValue.MethodByName("Set")
    method.Call([]reflect.Value{reflect.ValueOf(100)})
}
```

如果参数是一个指针类型，那么只需要构造指针类型的 reflect.Value 即可。

## 反射在运行时创建结构体

除了使用 reflect.TypeOf 函数生成已知类型的反射类型，还可以使用 reflect 标准库中的 ArrayOf、SliceOf 等函数生成一些在编译时完全不存在的类型或对象。对于结构体，需要使用 reflect.StructOf 函数在运行时生成特定的结构体对象。

```go
func StructOf(fields []StructField) Type
```

reflect.StructOf 函数参数是 StructField 的切片，StructField 代表结构体中的字段。其中，Name 代表该字段名，Type 代表该字段的类型。

```go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string

	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}
```

下面看一个生成结构体反射对象的例子。该函数可变参数中的类型依次构建为结构体的字段，并返回结构体变量。

```go
package main

import (
	"fmt"
	"reflect"
)

func MakeStruct(vals ...interface{}) reflect.Value {
	var structFields []reflect.StructField
	for k, v := range vals {
		t := reflect.TypeOf(v)
		structField := reflect.StructField{
			Name: fmt.Sprintf("F%d", k),
			Type: t,
		}
		structFields = append(structFields, structField)
	}
	structType := reflect.StructOf(structFields)
	structValue := reflect.New(structType)
	return structValue
}

func main() {
	s := MakeStruct(0, "", []int{})
	s.Elem().Field(0).SetInt(1)
	s.Elem().Field(1).SetString("hello")
	s.Elem().Field(2).Set(reflect.ValueOf([]int{1, 2, 3}))
	fmt.Println(s)
}
```

## 函数与反射

下例实现函数的动态调用，这和方法的调用是相同的，同样使用了 reflect.Call。如果函数中的参数为指针，那么可以借助 reflect.New 生成指定类型的反射指针对象。

```go
package main

import "reflect"

func fn(x int, y *int) {
	*y = x
}

func main() {
	rType := reflect.TypeOf(fn)
	rValue := reflect.ValueOf(fn)
	args := make([]reflect.Value, rType.NumIn())
	args[0] = reflect.ValueOf(0)
	args[1] = reflect.New(reflect.TypeOf(3))
	rValue.Call(args)
}
```

## 反射与其他类型

对于其他的一些类型，可以通过 XXXof 方法构造特定的 reflect.Type 类型，下例中介绍了一些复杂类型的反射实现过程。

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	ta := reflect.ArrayOf(10, reflect.TypeOf(0))
	tc := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0))
	tp := reflect.PtrTo(reflect.TypeOf(0))
	ts := reflect.SliceOf(reflect.TypeOf(0))
	tm := reflect.MapOf(reflect.TypeOf(0), reflect.TypeOf(0))
	tf := reflect.FuncOf([]reflect.Type{reflect.TypeOf(0)}, []reflect.Type{reflect.TypeOf(0)}, false)
	tt := reflect.StructOf([]reflect.StructField{
		{Name: "A", Type: reflect.TypeOf(0)},
		{Name: "B", Type: reflect.TypeOf(0)},
	})

	fmt.Println(ta, tc, tp, ts, tm, tf, tt)
}
```

根据 reflect.Type 生成对应的 reflect.Value，reflect 包中提供了对应类型的 makeXXX 方法。

```go
func MakeChan(typ Type, size int) Value
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value
func MakeMap(typ Type) Value
func MakeMapWithSize(typ Type, cap int) Value
func MakeSlice(typ Type, len, cap int) Value
```

除此之外，还可以使用 reflect.New 方法根据反射的类型分配相应大小的内存。

```go

```
