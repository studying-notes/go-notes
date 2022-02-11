---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 单元测试"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/docs/test/unittest"  # 设置网页永久链接
tags: [ "go", "unittest" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 源代码目录结构

在 gotest 包中创建两个文件，目录结构如下所示：

```
|--[src]
   |--[gotest]
      |--unit.go
      |--unit_test.go
```

其中 `unit.go` 为源代码文件，`unit_test.go` 为测试文件。要保证测试文件以"_test.go"结尾。

## 源代码文件

源代码文件 `unit.go` 中包含一个 `Add()` 方法，如下所示：

```go
package gotest

// Add 方法用于演示go test使用
func Add(a int, b int) int {
    return a + b
}
```

`Add()` 方法仅提供两数加法，实际项目中不可能出现类似的方法，此处仅供单元测试示例。

## 测试文件

测试文件 `unit_test.go` 中包含一个测试方法 `TestAdd()`，如下所示：

```go
package gotest_test

import (
    "testing"
    "gotest"
)

func TestAdd(t *testing.T) {
    var a = 1
    var b = 2
    var expected = 3

    actual := gotest.Add(a, b)
    if actual != expected {
        t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
    }
}
```

通过 package 语句可以看到，测试文件属于 "gotest_test" 包，测试文件也可以跟源文件在同一个包，但常见的做法是创建一个包专用于测试，这样可以使测试文件和源文件隔离。Go 源代码以及其他知名的开源框架通常会创建测试包，而且规则是在原包名上加上 "_test"。

测试函数命名规则为 "TestXxx"，其中 "Test" 为单元测试的固定开头，go test 只会执行以此为开头的方法。紧跟 "Test" 是以首字母大写的单词，用于识别待测试函数。

测试函数参数并不是必须要使用的，但 "testing.T" 提供了丰富的方法帮助控制测试流程。

t.Errorf() 用于标记测试失败，标记失败还有几个方法，在介绍 testing.T 结构时再详细介绍。

## 执行测试

命令行下，使用 `go test` 命令即可启动单元测试，如下所示：

```
src\gotest> go test
PASS
ok      gotest  0.378s
```

通过打印可知，测试通过，花费时间为 0.378s。

## 运行整个测试文件

```shell
go test xxx_test.go
```

## 运行测试文件中的某个函数

```shell
go test -run TestXXX
```
