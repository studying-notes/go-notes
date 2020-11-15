---
date: 2020-11-14T22:38:03+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "range 循环遍历"  # 文章标题
url:  "posts/go/abc/range"  # 设置网页永久链接
tags: [ "go"]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 热身

### 题目一

```go
func Process1(tasks []string) {
	for _, task := range tasks {
		// 启动协程并发处理任务
		go func() {
			fmt.Printf("Worker start process task: %s\n", task)
		}()
	}
}
```

函数 `Process1()` 用于处理任务，每个任务均启动一个协程进行处理。请问函数是否有问题？

### 题目二

```go
func Process2(tasks []string) {
	for _, task := range tasks {
		// 启动协程并发处理任务
		go func(t string) {
			fmt.Printf("Worker start process task: %s\n", t)
		}(task)
	}
}
```

函数 `Process2()` 用于处理任务，每个任务均启动一个协程进行处理。协程匿名函数接收一个任务作为参数，并进行处理。请问函数是否有问题？

### 题目三

项目中经常需要编写单元测试，而单元测试最常见的是 `table-driven` 风格的测试，如下所示：待测函数很简单，只是计算输入数值的 2 倍值。

```go
func Double(a int) int {
	return a * 2
}
```

测试函数如下：

```go
func TestDouble(t *testing.T) {
	var tests = []struct {
		name         string
		input        int
		expectOutput int
	}{
		{
			name:         "double 1 should got 2",
			input:        1,
			expectOutput: 2,
		},
		{
			name:         "double 2 should got 4",
			input:        2,
			expectOutput: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectOutput != Double(test.input) {
				t.Fatalf("expect: %d, but got: %d", test.input, test.expectOutput)
			}
		})
	}
}

```
上述测试函数也很简单，通过设计多个测试用例，标记输入输出，使用子测试进行验证。请问，上述测试有没有问题？

## 原理剖析

上述三个问题，有个共同点就是都引用了循环变量。即在 `for index, value := range xxx` 语句中，`index` 和 `value` 便是循环变量。不同点是循环变量的使用方式，有的是直接在协程中引用（题目一），有的作为参数传递（题目二），而题目三则是兼而有之。

回答以上问题，记住以下两点即可。

### 循环变量是易变的

首先，循环变量实际上只是一个普通的变量。

语句 `for index, value := range xxx` 中，每次循环 index 和 value 都会被重新赋值（并非生成新的变量）。

如果循环体中会启动协程（并且协程会使用循环变量），就需要格外注意了，因为很可能循环结束后协程才开始执行，此时，所有协程使用的循环变量有可能已被改写。（是否会改写取决于引用循环变量的方式）

### 循环变量需要绑定

在题目一中，协程函数体中引用了循环变量 `task`，协程从被创建到被调度执行期间循环变量极有可能被改写，这种情况下，我们称之为变量没有绑定。所以，题目一打印结果是混乱的。很有可能（随机）所有协程执行的 `task` 都是列表中的最后一个 task。

在题目二中，协程函数体中并没有直接引用循环变量 `task`，而是使用的参数。而在创建协程时，循环变量 `task` 作为函数参数传递给了协程。参数传递的过程实际上也生成了新的变量，也即间接完成了绑定。所以，题目二实际上是没有问题的。

在题目三中，测试用例名字 `test.name` 通过函数参数完成了绑定，而 `test.input 和 test.expectOutput` 则没有绑定。然而题目三实际执行却不会有问题，因为 `t.Run(...)` 并不会启动新的协程，也就是循环体并没有并发。此时，即便循环变量没有绑定也没有问题。但是风险在于，如果 `t.Run(...)` 执行的测试体有可能并发（比如通过 `t.Parallel()` ），此时就极有可能引入问题。

对于题目三，建议显式地绑定，例如：

```go
	for _, test := range tests {
		tc := test // 显式绑定，每次循环都会生成一个新的tc变量
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectOutput != Double(tc.input) {
				t.Fatalf("expect: %d, but got: %d", tc.input, tc.expectOutput)
			}
		})
	}
```

通过 `tc := test` 显式地绑定，每次循环会生成一个新的变量。

## 总结

- 如果循环体没有并发出现，则引用循环变量一般不会出现问题；
- 如果循环体有并发，则根据引用循环变量的位置不同而有所区别：
	-  通过参数完成绑定，则一般没有问题；
	-  函数体中引用，则需要显式地绑定。
