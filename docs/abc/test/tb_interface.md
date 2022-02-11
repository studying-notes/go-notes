---
date: 2020-11-16T18:18:28+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "深入测试标准库之 TB 接口"  # 文章标题
url:  "posts/go/abc/memory/test/tb_interface"  # 设置网页永久链接
tags: [ "go", "test" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 简介

TB 接口，顾名思义，是 testing.T( 单元测试 ) 和 testing.B( 性能测试 ) 共用的接口。

TB 接口通过在接口中定义一个名为 private( ）的私有方法，保证了即使用户实现了类似的接口，也不会跟 testing.TB 接口冲突。

其实，这些接口在 testing.T 和 testing.B 公共成员 testing.common 中已经实现。

## 接口定义

在 `src/testing/testing.go` 中定义了 testing.TB 接口：

```go
// TB is the interface common to T and B.
type TB interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	Helper()

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate Go 1 compatibility.
	private()
}
```

其中对外接口需要 testing.T 和 testing.B 实现，但由于 testing.T 和 testing.B 都继承了 testing.common，而 testing.common 已经实现了这些接口，所以 testing.T 和 testing.B 天然实现了 TB 接口。

其中私有接口 `private()` 用于控制该接口的唯一性，即便用户代码中某个类型实现了这些方法，由于无法实现这个私有接口，也不能被认为是实现了 TB 接口，所以不会跟用户代码产生冲突。

## 接口分类

我们在 testing.common 部分介绍过每个接口的实现，我们接下来就从函数功能上对接口进行分类。

以单元测试为例，每个测试函数都需要接收一个 testing.T 类型的指针作为函数参数，该参数主要用于控制测试流程（如结束和跳过）和记录日志。

### 记录日志

* Log(args ...interface{})
* Logf(format string, args ...interface{})

Log() 和 Logf() 负责记录日志，其区别在于是否支持格式化参数。

### 标记失败+记录日志

* Error(args ...interface{})
* Errorf(format string, args ...interface{})

Error() 和 Errorf() 负责标记当前测试失败并记录日志。只标记测试状态为失败，并不影响测试函数流程，不会结束当前测试，也不会退出当前测试。

### 标记失败+记录日志+结束测试

* Fatal(args ...interface{})
* Fatalf(format string, args ...interface{})

Fatal() 和 Fatalf() 负责标记当前测试失败、记录日志，并退出当前测试。

### 标记失败

* Fail()

Fail() 仅标记录前测试状态为失败。

### 标记失败并退出

* FailNow()

FailNow() 标记当前测试状态为失败并退出当前测试。

### 跳过测试+记录日志并退出

* Skip(args ...interface{})
* Skipf(format string, args ...interface{})

Skip() 和 Skipf() 标记当前测试状态为跳过并记录日志，最后退出当前测试。

### 跳过测试并退出

* SkipNow()

SkipNow() 标记测试状态为跳过，并退出当前测试。

## 私有接口避免冲突

接口定义中的 private() 方法是一个值得学习的用法。其目的是限定 testing.TB 接口的全局唯一性，即便用户的某个类型实现了除 private() 方法以外的其他方法，也不能说明实现了 testing.TB 接口，因为无法实现 private() 方法，private() 方法属于 testing 包内部可见，外部不可见。
