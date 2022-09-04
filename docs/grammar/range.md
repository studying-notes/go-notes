---
date: 2022-09-04T15:15:47+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 循环遍历"  # 文章标题
url:  "posts/go/docs/grammar/range"  # 设置网页永久链接
tags: [ "Go", "range" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

- [下标访问与拷贝](#下标访问与拷贝)
- [循环次数在开始前就已经确定](#循环次数在开始前就已经确定)
- [循环变量是易变的](#循环变量是易变的)
- [循环变量需要绑定](#循环变量需要绑定)
  - [未绑定](#未绑定)
  - [参数绑定](#参数绑定)
  - [单元测试案例](#单元测试案例)

## 下标访问与拷贝

用 index,value 接收 range 返回值会发生一次数据拷贝。

```go
func RangeSlice(slice []int) {
    for index, value := range slice {
        _, _ = index, value
    }
}
```

函数中使用 for-range 对切片进行遍历，获取切片的下标和元素素值，这里忽略函数的实际意义。

遍历过程中每次迭代会对 index 和 value 进行赋值，如果数据量大或者 value 类型为 string 时，对 value 的赋值操作可能是多余的，可以在 for-range 中忽略 value 值，使用 slice[index] 引用 value 值。

```go
func RangeMap(myMap map[int]string) {
    for key, _ := range myMap {
        _, _ = key, myMap[key]
    }
}
```

函数中使用 for-range 对 map 进行遍历，获取 map 的 key 值，并根据 key 值获取获取 value 值，这里忽略函数的实际意义。

函数中 for-range 语句中只获取 key 值，然后根据 key 值获取 value 值，虽然看似减少了一次赋值，但通过 key 值查找 value 值的性能消耗可能高于赋值消耗。能否优化取决于 map 所存储数据结构特征、结合实际情况进行。

## 循环次数在开始前就已经确定

```go
func main() {
    v := []int{1, 2, 3}
    for i:= range v {
        v = append(v, i)
    }
}
```

main() 函数中定义一个切片 v，通过 range 遍历 v，遍历过程中不断向 v 中添加新的元素。

循环内改变切片的长度，不影响循环次数，**循环次数在循环开始前就已经确定了**。

## 循环变量是易变的

首先，循环变量实际上只是一个普通的变量。

语句 `for index, value := range xxx` 中，每次循环 index 和 value 都会被重新赋值（并非生成新的变量）。

如果循环体中会启动协程（并且协程会使用循环变量），就需要格外注意了，因为很可能**循环结束后协程才开始执行**，此时，所有协程使用的循环变量有可能已被改写。（是否会改写取决于引用循环变量的方式）

## 循环变量需要绑定

### 未绑定

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

函数 `Process1()` 用于处理任务，每个任务均启动一个协程进行处理。

协程函数体中引用了循环变量 `task`，协程从被创建到被调度执行期间循环变量极有可能被改写，这种情况下，我们称之为变量没有绑定。

所以，打印结果是混乱的。很有可能（随机）所有协程执行的 `task` 都是列表中的最后一个 task。

### 参数绑定

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

函数 `Process2()` 用于处理任务，每个任务均启动一个协程进行处理。协程匿名函数接收一个任务作为参数，并进行处理。

协程函数体中并没有直接引用循环变量 `task`，而是使用的参数。而在创建协程时，循环变量 `task` 作为函数参数传递给了协程。参数传递的过程实际上也生成了新的变量，也即间接完成了绑定。所以，实际上是没有问题的。

### 单元测试案例

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

测试用例名字 `test.name` 通过函数参数完成了绑定，而 `test.input 和 test.expectOutput` 则没有绑定。

然而实际执行却不会有问题，因为 `t.Run(...)` 并不会启动新的协程，也就是循环体并没有并发。此时，即便循环变量没有绑定也没有问题。

但是风险在于，如果 `t.Run(...)` 执行的测试体有可能并发（比如通过 `t.Parallel()` ），此时就极有可能引入问题。

建议显式地绑定，例如：

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

```go

```
