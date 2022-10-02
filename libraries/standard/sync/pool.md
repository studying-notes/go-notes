---
date: 2022-09-08T11:08:14+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 对象重用的机制"  # 文章标题
url:  "posts/go/libraries/standard/sync/pool"  # 设置网页永久链接
tags: [ "Go", "pool" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 简介

用于保存和复用临时对象，减少内存分配，降低 GC 压力。

sync.Pool 是可伸缩的，同时也是并发安全的，其大小仅受限于内存的大小。sync.Pool 用于存储那些被分配了但是没有被使用，而未来可能会使用的值。这样就可以不用再次经过内存分配，可直接复用已有对象，减轻 GC 的压力，从而提升系统的性能。

sync.Pool 的大小是可伸缩的，高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理。

## 用法

### 场景

```go
type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var buf, _ = json.Marshal(Student{Name: "Geektutu", Age: 25})

func unmarsh() {
	stu := &Student{}
	json.Unmarshal(buf, stu)
}
```

json 的反序列化在文本解析和网络通信过程中非常常见，当程序并发度非常高的情况下，短时间内需要创建大量的临时对象。而这些对象是都是分配在堆上的，会给 GC 造成很大压力，严重影响程序的性能。

### 声明对象池

只需要实现 New 函数即可。对象池中没有对象时，将会调用 New 函数创建。

```go
var studentPool = sync.Pool{
    New: func() interface{} {
        return new(Student)
    },
}
```

### 获取和归还

```go
// 获取
stu := studentPool.Get().(*Student)

json.Unmarshal(buf, stu)

// 归还
studentPool.Put(stu)
```

- Get() 用于从对象池中获取对象，因为返回值是 interface{}，因此需要类型转换。
- Put() 则是在对象使用完毕后，返回对象池。

### 性能测试

```go
func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := &Student{}
		json.Unmarshal(buf, stu)
	}
}

func BenchmarkUnmarshalWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stu := studentPool.Get().(*Student)
		json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}
```

```bash
go test -bench . -benchmem
```

在这个例子中，因为 Student 结构体内存占用较小，内存分配几乎不耗时间。而标准库 json 反序列化时利用了反射，效率是比较低的，占据了大部分时间，因此两种方式最终的执行时间几乎没什么变化。但是内存占用差了一个数量级，使用了 sync.Pool 后，内存占用仅为未使用的 234/5096 = 1/22，对 GC 的影响就很大了

## 标准库中的案例

Go 语言标准库也大量使用了 sync.Pool，例如 fmt 和 encoding/json。

fmt.Printf 的调用是非常频繁的，利用 sync.Pool 复用 pp 对象能够极大地提升性能，减少内存占用，同时降低 GC 压力。

```go
// go 1.13.6

// pp is used to store a printer's state and is reused with sync.Pool to avoid allocations.
type pp struct {
    buf buffer
    ...
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}

// free saves used pp structs in ppFree; avoids an allocation per invocation.
func (p *pp) free() {
	if cap(p.buf) > 64<<10 {
		return
	}

	p.buf = p.buf[:0]
	p.arg = nil
	p.value = reflect.Value{}
	p.wrappedErr = nil
	ppFree.Put(p)
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()
	return
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}
```
