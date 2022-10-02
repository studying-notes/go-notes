---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "位运算"  # 文章标题
url:  "posts/go/algorithm/math/bit_operation"  # 设置网页永久链接
tags: [ "Go", "bit-operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [位运算符号](#位运算符号)
- [数的二进制表示](#数的二进制表示)
- [问题](#问题)
  - [判断二进制数某位是否为 1](#判断二进制数某位是否为-1)
  - [求绝对值](#求绝对值)
- [二级](#二级)
  - [三级](#三级)

## 位运算符号

| 符号 | 含义 |
| ---- | -------- |
| & | AND 与 |
| \ | | OR 或 |
| ^ | XOR / 取反 |
| &^ | 位清空 |
| << | 左移 |
| >> | 右移 |

## 数的二进制表示

```go
package main

import "fmt"

func main() {
	fmt.Printf("%08b\n", 6)
	fmt.Printf("%016b\n", 24)
	fmt.Printf("%032b\n", 36)
}
```

## 问题

### 判断二进制数某位是否为 1

```go
// 判断二进制数从右往左数第 i 位是否为 1，位数从 0 开始
func IsOneInt(n, i int) bool {
	return (1<<i)&n == 1<<i
}

// 从结果上看两者没有区别
func IsOneUint(n, i int) bool {
	return (uint(n) & (uint(1) << uint(i))) == uint(1)<<uint(i)
}
```

### 求绝对值

> 正数求负数是取反加一，负数求正数就是减一取反。

仅用于学习，实际性能与取负数是差不多的，因为编译器优化了。

```go
func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return ^(n - 1)
}
```

## 二级

### 三级

```go

```

```go

```
