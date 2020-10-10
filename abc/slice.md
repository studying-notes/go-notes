---
date: 2020-07-26T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 语言切片"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/abc/slice"  # 设置网页永久链接
tags: [ "go", "slice" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 切片的结构定义

切片（slice）就是一种简化版的动态数组。因为动态数组的长度不固定，所以切片的长度自然也就不能是类型的组成部分了。数组虽然有适用的地方，但是数组的类型和操作都不够灵活，因此在 Go 代码中数组使用得并不多。而切片则使用得相当广泛。

切片的结构定义，即 `reflect.SliceHeader`：

```go
// SliceHeader is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

由此可以看出切片的开头部分和 Go 字符串是一样的，但是切片多了一个 Cap 成员表示切片指向的内存空间的最大容量（对应元素的个数，而不是字节数）。

`x := []int{2,3,5, 7,11}` 和 `y := x[1:3]` 两个切片对应的内存结构：

![06NAIJ.png](https://s1.ax1x.com/2020/10/10/06NAIJ.png)

切片的定义方式：

```go
a []int               // nil切片，和nil相等，一般用来表示一个不存在的切片
b = []int{}           // 空切片，和nil不相等，一般用来表示一个空的集合
c = []int{1, 2, 3}    // 有3个元素的切片，len和cap都为3
d = c[:2]             // 有2个元素的切片，len为2，cap为3
e = c[0:2:cap(c)]     // 有2个元素的切片，len为2，cap为3
f = c[:0]             // 有0个元素的切片，len为0，cap为3
g = make([]int, 3)    // 有3个元素的切片，len和cap都为3
h = make([]int, 2, 3) // 有2个元素的切片，len为2，cap为3
i = make([]int, 0, 3) // 有0个元素的切片，len为0，cap为3
```

和数组一样，内置的 len() 函数返回切片中有效元素的长度，内置的 cap() 函数返回切片容量大小，容量必须大于或等于切片的长度。也可以通过 reflect.SliceHeader 结构访问切片的信息（只是为了说明切片的结构，并不是推荐的做法）。切片可以和 nil 进行比较，只有当切片底层数据指针为空时切片本身才为 nil，这时候切片的长度和容量信息将是无效的。如果有切片的底层数据指针为空，但是长度和容量不为 0 的情况，那么说明切片本身已经被损坏了（例如，直接通过reflect.SliceHeader或unsafe包对切片作了不正确的修改）。

遍历切片的方式和遍历数组的方式类似：

```go
for i := range a {
    fmt.Printf("a[%d]: %d\n", i, a[i])
}
for i, v := range b {
    fmt.Printf("b[%d]: %d\n", i, v)
}
for i := 0; i < len(c); i++ {
    fmt.Printf("c[%d]: %d\n", i, c[i])
}
```

其实除了遍历之外，只要是切片的底层数据指针、长度和容量没有发生变化，对切片的遍历、元素的读取和修改就和数组一样。在对切片本身进行赋值或参数传递时，和数组指针的操作方式类似，但是**只复制切片头信息**（reflect.SliceHeader），而**不会复制底层的数据**。对于类型，和数组的最大不同是，切片的类型和长度信息无关，只要是相同类型元素构成的切片均对应相同的切片类型。

### 切片函数传参理解

通过函数可以修改切片的底层数据：

```go
func change(array []int) {
	array[0] = 0
}

func main() {
	array := []int{9}
	change(array)
	fmt.Println(array)
}
```

## 添加切片元素

### append

内置的泛型函数 `append()` 可以在切片的尾部追加 N 个元素：

```go
var a []int
a = append(a, 1)                 // 追加一个元素
a = append(a, 1, 2, 3)           // 追加多个元素，手写解包方式
a = append(a, []int{1, 2, 3}...) // 追加一个切片，切片需要解包
```

不过要注意的是，在容量不足的情况下，`append()` 操作会导致重新分配内存，可能导致巨大的内存分配和复制数据的代价。即使容量足够，依然需要用 `append()` 函数的返回值来更新切片本身，因为新切片的长度已经发生了变化。

除了在切片的尾部追加，还可以在切片的开头添加元素：

```go
var a = []int{1, 2, 3}
a = append([]int{0}, a...)          // 在开头添加一个元素
a = append([]int{-3, -2, -1}, a...) // 在开头添加一个切片
```

**在开头一般都会导致内存的重新分配**，而且会导致已有的元素全部复制一次。因此，从切片的开头添加元素的性能一般要比从尾部追加元素的性能差很多。

由于 append() 函数返回新的切片，也就是它支持链式操作，因此我们可以将多个 append () 操作组合起来，实现在切片中间插入元素：

```go
var a []int
a = append(a[:i], append([]int{x}, a[i:]...)...)       // 在第i个位置插入x
a = append(a[:i], append([]int{1, 2, 3}, a[i:]...)...) // 在第i个位置插入切片
```

每个添加操作中的第二个 append () 调用都会创建一个临时切片，并将 `a[i:]` 的内容复制到新创建的切片中，然后将临时创建的切片再追加到 `a[:i]`。

### copy

用 copy() 和 append() 组合可以避免创建中间的临时切片，同样是完成添加元素的操作：

```go
a = append(a, 0)     // 切片扩展一个空间
copy(a[i+1:], a[i:]) // a[i:]向后移动一个位置
a[i] = x             // 设置新添加的元素
```

第一句中的 append() 用于扩展切片的长度，为要插入的元素留出空间。第二句中的 copy() 操作将要插入位置开始之后的元素向后挪动一个位置。第三句真实地将新添加的元素赋值到对应的位置。操作语句虽然冗长了一点，但是相比前面的方法，可以减少中间创建的临时切片。

用 copy() 和 append() 组合也可以实现在中间位置插入多个元素（也就是插入一个切片）：

```go
a = append(a, x...)         // 为x切片扩展足够的空间
copy(a[i+len(x):], a[i:])   // a[i:]向后移动len(x)个位置
copy(a[i:], x)              // 复制新添加的切片
```

稍显不足的是，在第一句扩展切片容量的时候，扩展空间部分的元素复制是没有必要的。没有专门的内置函数用于扩展切片的容量，append() 本质是用于追加元素而不是扩展容量，扩展切片容量只是 append() 的一个副作用。

## 删除切片元素

根据要删除元素的位置，有从开头位置删除、从中间位置删除和从尾部删除3种情况，其中删除切片尾部的元素最快：

```go
a = []int{1, 2, 3}
a = a[:len(a)-1]   // 删除尾部1个元素
a = a[:len(a)-N]   // 删除尾部N个元素
```

删除开头的元素可以直接移动数据指针：

```go
a = []int{1, 2, 3}
a = a[1:] // 删除开头1个元素
a = a[N:] // 删除开头N个元素
```

删除开头的元素也可以不移动数据指针，而将后面的数据向开头移动。可以用 append() 原地完成（所谓原地完成是指在原有的切片数据对应的内存区间内完成，不会导致内存空间结构的变化）：

```go
a = []int{1, 2, 3}
a = append(a[:0], a[1:]...) // 删除开头1个元素
a = append(a[:0], a[N:]...) // 删除开头N个元素
```

也可以用 copy() 完成删除开头的元素：

```go
a = []int{1, 2, 3}
a = a[:copy(a, a[1:])] // 删除开头1个元素
a = a[:copy(a, a[N:])] // 删除开头N个元素
```

对于删除中间的元素，需要对剩余的元素进行一次整体挪动，同样可以用 append() 或 copy() 原地完成：

```go
a = []int{1, 2, 3, ...}

a = append(a[:i], a[i+1:]...) // 删除中间1个元素
a = append(a[:i], a[i+N:]...) // 删除中间N个元素

a = a[:i+copy(a[i:], a[i+1:])]  // 删除中间1个元素
a = a[:i+copy(a[i:], a[i+N:])]  // 删除中间N个元素
```

## 切片内存技巧

对于切片来说，len 为 0 但是 cap 容量不为 0 的切片则是非常有用的特性。当然，如果 len 和 cap 都为 0 的话，则变成一个真正的空切片，虽然它并不是一个 nil 的切片。在判断一个切片是否为空时，一般通过 len 获取切片的长度来判断，一般很少将切片和 nil 做直接的比较。

例如下面的 `TrimSpace()` 函数用于删除 `[]byte` 中的空格。函数实现利用了长度为 0 的切片的特性，实现高效而且简洁。

```go
func TrimSpace(s []byte) []byte {
    b := s[:0]
    for _, x := range s {
        if x != ' ' {
            b = append(b, x)
        }
    }
    return b
}
```

其实类似的根据过滤条件原地删除切片元素的算法都可以采用类似的方式处理（因为是删除操作，所以不会出现内存不足的情形）：

```go
func Filter(s []byte, fn func(x byte) bool) []byte {
    b := s[:0]
    for _, x := range s {
        if !fn(x) {
            b = append(b, x)
        }
    }
    return b
}
```

切片高效操作的要点是要降低内存分配的次数，尽量保证 append() 操作不会超出 cap 的容量，降低触发内存分配的次数和每次分配内存的大小。

## 避免切片内存泄漏

如前所述，切片操作并不会复制底层的数据。底层的数组会被保存在内存中，直到它不再被引用。但是有时候可能会因为一个小的内存引用而导致底层整个数组处于被使用的状态，这会延迟垃圾回收器对底层数组的回收。

例如，FindPhoneNumber() 函数加载整个文件到内存，然后搜索第一个出现的电话号码，最后结果以切片方式返回。

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return regexp.MustCompile("[0-9]+").Find(b)
}
```

这段代码返回的 []byte 指向保存整个文件的数组。由于切片引用了整个原始数组，导致垃圾回收器不能及时释放底层数组的空间。一个小的需求可能导致需要长时间保存整个文件数据。这虽然不是传统意义上的内存泄漏，但是可能会降低系统的整体性能。

要解决这个问题，可以将感兴趣的数据复制到一个新的切片中（数据的传值是 Go 语言编程的一个哲学，虽然传值有一定的代价，但是换取的好处是切断了对原始数据的依赖）：

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = regexp.MustCompile("[0-9]+").Find(b)
    return append([]byte{}, b...)
}
```

类似的问题在删除切片元素时可能会遇到。假设切片里存放的是指针对象，那么下面删除末尾的元素后，被删除的元素依然被切片底层数组引用，从而导致不能及时被垃圾回收器回收（这要依赖回收器的实现方式）：

```go
var a []*int{ ... }
a = a[:len(a)-1]    // 被删除的最后一个元素依然被引用，可能导致垃圾回收器操作被阻碍
```

保险的方式是先将指向需要提前回收内存的指针设置为 nil，保证垃圾回收器可以发现需要回收的对象，然后再进行切片的删除操作：

```go
var a []*int{ ... }
a[len(a)-1] = nil // 垃圾回收器回收最后一个元素内存
a = a[:len(a)-1]  // 从切片删除最后一个元素
```

当然，如果切片存在的周期很短的话，可以不用刻意处理这个问题。因为如果切片本身已经可以被垃圾回收器回收的话，切片对应的每个元素自然也就可以被回收了。

## 切片类型强制转换

为了安全，当两个切片类型 []T 和 []Y 的底层原始切片类型不同时，Go 语言是无法直接转换类型的。不过安全都是有一定代价的，有时候这种转换是有它的价值的——可以简化编码或者是提升代码的性能。

例如在64位系统上，需要对一个 []float64 切片进行高速排序，我们可以将它强制转换为 []int 整数切片，然后以整数的方式进行排序（因为 float64 遵循 IEEE 754 浮点数标准特性，所以当浮点数有序时对应的整数也必然是有序的）。

下面的代码通过两种方法将 []float64 类型的切片转换为 []int 类型的切片：

```go
var a = []float64{4, 2, 43, 6, 7, 43, 23}

func SortFloat64FastV1(a []float64) {
	// 强制类型转换
	var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]
	
	// 以int方式给float64排序
	sort.Ints(b)
}
}

func SortFloat64FastV2(a []float64) {
	// 通过reflect.SliceHeader更新切片头部信息实现转换
	var c []int
	aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*cHdr = *aHdr

	// 以int方式给float64排序
	sort.Ints(c)
}
```

第一种强制转换是先将切片数据的开始地址转换为一个较大的数组的指针，然后对数组指针对应的数组重新做切片操作。中间需要 `unsafe.Pointer` 来连接两个不同类型的指针传递。需要注意的是，Go 语言实现中非 0 大小数组的长度不得超过 2 GB，因此需要针对数组元素的类型大小计算数组的最大长度范围（`[]uint8` 最大 2 GB，`[]uint16` 最大 1 GB，依此类推，但是 `[]struct{}` 数组的长度可以超过 2 GB）。

第二种转换操作是分别取两个不同类型的切片头信息指针，任何类型的切片头部信息底层都对应 `reflect.SliceHeader` 结构，然后通过更新结构体方式来更新切片信息，从而实现 a 对应的 `[]float64` 切片到 c 对应的 `[]int` 切片的转换。

通过基准测试，可以发现用 sort.Ints 对转换后的 []int 排序的性能要比用 sort.Float64s 排序的性能高一点。不过需要注意的是，这个方法可行的前提是要保证 []float64 中没有 NaN 和 Inf 等非规范的浮点数（因为浮点数中 NaN 不可排序，正 0 和负 0 相等，但是整数中没有这类情形）。
