---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 随机数标准库"  # 文章标题
url:  "posts/go/libraries/standard/rand"  # 设置网页链接，默认使用文件名
tags: [ "go", "random" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 基本函数

```go
// 该函数设置随机种子
// 若不调用此函数设置随机种子，则默认的种子值为 1
// 由于随机算法是固定的，如果每次都以 1 作为随机种
// 子开始产生随机数，则结果都是一样的，因此一般
// 都需要调用此函数来设置随机种子，通常的做法是
// 以当前时间作为随机种子以保证每次随机种子都不
// 同，从而产生的随机数也不同
// 该函数协程安全
func Seed(seed int64)

// 以下函数用来生成相应数据类型的随机数
// 带 n 的版本则生成 [0,n) 的随机数
// 生成的随机数都是非负数
func Float32() float32
func Float64() float64

func Int() int
func Intn(n int) int

// 该函数只返回 int32 表示范围内的非负数，
// 位数为 31，因此该函数叫做 Int31
func Int31() int32  
func Int31n(n int32) int32

func Int63() int64
func Int63n(n int64) int64

func Uint32() uint32
func Uint64() uint64

// 创建一个以 seed 为种子的源，不是协程安全的
func NewSource(seed int64) Source

// 以 src 为源创建随机对象
func New(src Source) *Rand

// 设置或重置种子，不是协程安全的
func (r *Rand) Seed(seed int64)

// 下面的函数和全局版本的函数功能一样
func (r *Rand) Float32() float32
func (r *Rand) Float64() float64
func (r *Rand) Int() int
func (r *Rand) Int31() int32
func (r *Rand) Int31n(n int32) int32
func (r *Rand) Int63() int64
func (r *Rand) Int63n(n int64) int64
func (r *Rand) Intn(n int) int
func (r *Rand) Uint32() uint32
func (r *Rand) Uint64() uint64
```

## 随机种子

```go
// 返回当前时间
func Now() Time

// 将 Time 类型转换为 int64 类型以作为随机种子

// 返回从 1970 年 1 月 1 日到 t 的秒数
func (t Time) Unix() int64
// 返回从 1970 年 1 月 1 日到 t 的纳秒数
func (t Time) UnixNano() int64
```

## 示例

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 全局函数
	rand.Seed(time.Now().Unix())

	fmt.Println(rand.Int())       // int随机值，返回值为int
	fmt.Println(rand.Intn(100))   // [0,100)的随机值，返回值为int
	fmt.Println(rand.Int31())     // 31位int随机值，返回值为int32
	fmt.Println(rand.Int31n(100)) // [0,100)的随机值，返回值为int32
	fmt.Println(rand.Float32())   // 32位float随机值，返回值为float32
	fmt.Println(rand.Float64())   // 64位float随机值，返回值为float64

	// 如果要产生负数到正数的随机值，只需要将生成的随机数减去相应数值即可
	fmt.Println(rand.Intn(100) - 50) // [-50, 50)的随机值

	// Rand对象
	r := rand.New(rand.NewSource(time.Now().Unix()))

	fmt.Println(r.Int())       // int随机值，返回值为int
	fmt.Println(r.Intn(100))   // [0,100)的随机值，返回值为int
	fmt.Println(r.Int31())     // 31位int随机值，返回值为int32
	fmt.Println(r.Int31n(100)) // [0,100)的随机值，返回值为int32
	fmt.Println(r.Float32())   // 32位float随机值，返回值为float32
	fmt.Println(r.Float64())   // 64位float随机值，返回值为float64

	// 如果要产生负数到正数的随机值，只需要将生成的随机数减去相应数值即可
	fmt.Println(r.Intn(100) - 50) // [-50, 50)的随机值
}
```
