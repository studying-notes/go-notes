---
date: 2020-11-15T13:41:14+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "range 循环遍历实现机制"  # 文章标题
url:  "posts/go/abc/range2"  # 设置网页永久链接
tags: [ "go" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 前言

range 是 Golang 提供的一种迭代遍历手段，可操作的类型有数组、切片、Map、channel 等，实际使用频率非常高。

## 热身

按照惯例，我们看几个有意思的题目，用于检测对 range 的了解程度。

### 题目一：切片遍历

下面函数通过遍历切片，打印切片的下标和元素值，请问性能上有没有可优化的空间？

```go
func RangeSlice(slice []int) {
    for index, value := range slice {
        _, _ = index, value
    }
}
```

程序解释：

函数中使用 for-range 对切片进行遍历，获取切片的下标和元素素值，这里忽略函数的实际意义。

参考答案：

遍历过程中每次迭代会对 index 和 value 进行赋值，如果数据量大或者 value 类型为 string 时，对 value 的赋值操作可能是多余的，可以在 for-range 中忽略 value 值，使用 slice[index] 引用 value 值。

### 题目二：Map 遍历

下面函数通过遍历 Map，打印 Map 的 key 和 value，请问性能上有没有可优化的空间？

```go
func RangeMap(myMap map[int]string) {
    for key, _ := range myMap {
        _, _ = key, myMap[key]
    }
}
```

程序解释：

函数中使用 for-range 对 map 进行遍历，获取 map 的 key 值，并根据 key 值获取获取 value 值，这里忽略函数的实际意义。

参考答案：

函数中 for-range 语句中只获取 key 值，然后根据 key 值获取 value 值，虽然看似减少了一次赋值，但通过 key 值查找 value 值的性能消耗可能高于赋值消耗。能否优化取决于 map 所存储数据结构特征、结合实际情况进行。

### 题目三：动态遍历

请问如下程序是否能正常结束？

```go
func main() {
    v := []int{1, 2, 3}
    for i:= range v {
        v = append(v, i)
    }
}
```

程序解释：

main() 函数中定义一个切片 v，通过 range 遍历 v，遍历过程中不断向 v 中添加新的元素。

参考答案：

能够正常结束。循环内改变切片的长度，不影响循环次数，**循环次数在循环开始前就已经确定了**。

## 实现原理

对于for-range语句的实现，可以从编译器源码中找到答案。

编译器源码 `gofrontend/go/statements.cc/For_range_statement::do_lower()` 方法中有如下注释。

```go
// Arrange to do a loop appropriate for the type.  We will produce
//   for INIT ; COND ; POST {
//           ITER_INIT
//           INDEX = INDEX_TEMP
//           VALUE = VALUE_TEMP // If there is a value
//           original statements
//   }
```

可见 range 实际上是一个 C 风格的循环结构。range 支持数组、数组指针、切片、map 和 channel 类型，对于不同类型有些细节上的差异。

### range for slice

下面的注释解释了遍历 slice 的过程：

```go
// The loop we generate:
//   for_temp := range
//   len_temp := len(for_temp)
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = for_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

遍历 slice 前会先获取 slice 的长度 len_temp 作为循环次数，循环体中，每次循环会先获取元素值，如果 for-range 中接收 index 和 value 的话，则会对 index 和 value 进行一次赋值。

由于循环开始前循环次数就已经确定了，所以循环过程中新添加的元素是没办法遍历到的。

另外，数组与数组指针的遍历过程与 slice 基本一致，不再赘述。

### range for map

下面的注释解释了遍历 map 的过程：

```go
// The loop we generate:
//   var hiter map_iteration_struct
//   for mapiterinit(type, range, &hiter); hiter.key != nil; mapiternext(&hiter) {
//           index_temp = *hiter.key
//           value_temp = *hiter.val
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

遍历 map 时没有指定循环次数，循环体与遍历 slice 类似。由于 map 底层实现与 slice 不同，map 底层使用 hash 表实现，插入数据位置是随机的，所以遍历过程中新插入的数据不能保证遍历到。

### range for channel

遍历 channel 是最特殊的，这是由 channel 的实现机制决定的：

```go
// The loop we generate:
//   for {
//           index_temp, ok_temp = <-range
//           if !ok_temp {
//                   break
//           }
//           index = index_temp
//           original body
//   }
```

channel 遍历是依次从 channel 中读取数据,读取前是不知道里面有多少个元素的。如果 channel 中没有元素，则会阻塞等待，如果 channel 已被关闭，则会解除阻塞并退出循环。

注：
- 上述注释中 index_temp 实际上描述是有误的，应该为 value_temp，因为 index 对于 channel 是没有意义的。
- 使用 for-range 遍历 channel 时只能获取一个返回值。

- 遍历过程中可以视情况放弃接收 index 或 value，可以一定程度上提升性能
- 遍历 channel 时，如果 channel 中没有数据，可能会阻塞
- 尽量避免遍历过程中修改原数据

## 小结

- for-range 的实现实际上是 C 风格的 for 循环
- 使用 index,value 接收 range 返回值会发生一次数据拷贝
