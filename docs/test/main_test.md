---
date: 2020-11-15T20:07:06+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go Main 测试"  # 文章标题
url:  "posts/go/docs/test/sub_test"  # 设置网页永久链接
tags: [ "go", "test" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类


# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 简介

我们知道子测试的一个方便之处在于可以让多个测试共享 Setup 和 Tear-down。但这种程度的共享有时并不满足需求，有时希望在整个测试程序做一些全局的 setup 和 Tear-down，这时就需要 Main 测试了。

所谓 Main 测试，即声明一个 `func TestMain(m * testing.M)`，它是名字比较特殊的测试，参数类型为 `testing.M` 指针。如果声明了这样一个函数，当前测试程序将不是直接执行各项测试，而是将测试交给 TestMain 调度。

## 示例

下面通过一个例子来展示 Main 测试用法：

```go
// TestMain 用于主动执行各种测试，可以测试前后做setup和tear-down操作
func TestMain(m *testing.M) {
    println("TestMain setup.")

    retCode := m.Run() // 执行测试，包括单元测试、性能测试和示例测试

    println("TestMain tear-down.")

    os.Exit(retCode)
}
```

上述例子中，日志打印的两行分别对应 Setup 和 Tear-down 代码，m.Run() 即为执行所有的测试，m.Run() 的返回结果通过 os.Exit() 返回。

如果所有测试均通过测试，m.Run() 返回 0，否同 m.Run() 返回 1，代表测试失败。

有一点需要注意的是，TestMain 执行时，命令行参数还未解析，如果测试程序需要依赖参数，可以使用 `flag.Parse()` 解析参数，m.Run() 方法内部还会再次解析参数，此处解析不会影响原测试过程。
