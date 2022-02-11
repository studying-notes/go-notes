---
date: 2020-11-16T20:56:19+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "深入测试标准库之 Main 测试"  # 文章标题
url:  "posts/go/abc/memory/test/main"  # 设置网页永久链接
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

每一种测试（单元测试、性能测试或示例测试），都有一个数据类型与其对应。

* 单元测试：InternalTest
* 性能测试：InternalBenchmark
* 示例测试：InternalExample

测试编译阶段，每个测试都会被放到指定类型的切片中，测试执行时，这些测试将会被放到 testing.M 数据结构中进行调度。

而 testing.M 即是 MainTest 对应的数据结构。

## 数据结构

源码 `src/testing/testing.go:M` 定义了 testing.M 的数据结构：

```go
// M is a type passed to a TestMain function to run the actual tests.
type M struct {
	tests      []InternalTest       // 单元测试
	benchmarks []InternalBenchmark  // 性能测试
	examples   []InternalExample    // 示例测试
	timer     *time.Timer           // 测试超时时间
}
```

单元测试、性能测试和示例测试在经过编译后都会被存放到一个 testing.M 数据结构中，在测试执行时该数据结构将传递给 TestMain()，真正执行测试的是 testing.M 的 Run() 方法，这个后面我们会继续分析。

timer 用于指定测试的超时时间，可以通过参数 `timeout <n>` 指定，当测试执行超时后将会立即结束并判定为失败。

## 执行测试

TestMain() 函数通常会有一个 m.Run() 方法，该方法会执行单元测试、性能测试和示例测试，如果用户实现了 TestMain() 但没有调用 m.Run() 的话，那么什么测试都不会被执行。

m.Run() 不仅会执行测试，还会做一些初始化工作，比如解析参数、启动定时器、根据参数指示创建一系列的文件等。

m.Run() 使用三个独立的方法来执行三种测试：

* 单元测试：runTests(m.deps.MatchString, m.tests)
* 性能测试：runExamples(m.deps.MatchString, m.examples)
* 示例测试：runBenchmarks(m.deps.ImportPath(), m.deps.MatchString, m.benchmarks)

其中 m.deps 里存放了测试匹配相关的内容，暂时先不用关注。
