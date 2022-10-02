---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 单元测试"  # 文章标题
url:  "posts/go/libraries/standard/testing/unittest"  # 设置网页永久链接
tags: [ "Go", "unittest" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

创建两个文件，其中 `unit.go` 为源代码文件，`unit_test.go` 为测试文件。要保证测试文件以"_test.go"结尾。
## 源代码文件

源代码文件 `unit.go` 中包含一个 `Add()` 方法，如下所示：

```go
package main

// Add 方法用于演示go test使用
func Add(a int, b int) int {
    return a + b
}
```

`Add()` 方法仅提供两数加法，实际项目中不可能出现类似的方法，此处仅供单元测试示例。

## 测试文件

测试文件 `unit_test.go` 中包含一个测试方法 `TestAdd()`，如下所示：

```go
package main

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

测试函数命名规则为 "TestXxx"，其中 "Test" 为单元测试的固定开头，go test 只会执行以此为开头的方法。紧跟 "Test" 是以首字母大写的单词，用于识别待测试函数。

测试函数参数并不是必须要使用的，但 "testing.T" 提供了丰富的方法帮助控制测试流程。比如 t.Errorf() 用于标记测试失败。

## 执行测试

命令行下，使用 `go test` 命令即可启动单元测试，如下所示：

```bash
go test
```

```
PASS
ok      gotest  0.378s
```

```go

```
