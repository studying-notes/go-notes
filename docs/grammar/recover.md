---
date: 2020-11-15T13:28:06+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 捕获异常"  # 文章标题
url:  "posts/go/docs/grammar/recover"  # 设置网页永久链接
tags: [ "Go", "recover" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

项目中，有时为了让程序更健壮，也即不 `panic`，我们或许会使用 `recover()` 来接收异常并处理。

`panic` 会停掉当前正在执行的程序，但是与 `os.Exit(-1)` 这种退出不同，`panic` 的撤退比较有秩序，先处理完**当前** `goroutine` 已经 `defer` 挂上去的任务，即在 `panic` 语句之前注册的 `defer` 语句，执行完毕后再退出整个程序。

而 `defer` 的存在，让我们可以在 `defer` 中通过 `recover` 获取 `panic`，从而达到捕获异常的效果。

`panic` 允许传递一个参数，参数通常是将出错的信息以字符串的形式来表示，`panic` 会打印这个字符串，以及触发 `panic` 的调用栈。

`panic` 的原则是：执行且只执行当前 `goroutine` 的 `defer`，且在 `panic` 之前注册。

`panic` 仅保证当前 goroutine 下的 `defer` 都会被执行到，但不保证其他协程的 `defer` 也会执行到。如果是在同一 `goroutine` 下的调用者的 `defer`，那么可以一路回溯回去执行；但如果是不同 `goroutine`，那就不做保证了。

`recover` 只在 `defer` 的函数中有效，如果不是在 `defer` 上下文中调用，`recover` 会直接返回 `nil`。

比如以下代码：

```go
func NoPanic() {
	if err := recover(); err != nil {
		fmt.Println("Recover success...")
	}
}

func Dived(n int) {
	defer NoPanic()

	fmt.Println(1/n)
}
```

`func NoPanic()` 会自动接收异常，并打印相关日志，算是一个通用的异常处理函数。

业务处理函数中只要使用了 `defer NoPanic()`，那么就不会再有 `panic` 发生。

## 误区

在项目中，有众多的数据库更新操作，正常的更新操作需要提交，而失败的就需要回滚，如果异常分支比较多，就会有很多重复的回滚代码，所以有人尝试了一个做法：即在 defer 中判断是否出现异常，有异常则回滚，否则提交。

简化代码如下所示：

```go
func IsPanic() bool {
	if err := recover(); err != nil {
		fmt.Println("Recover success...")
		return true
	}

	return false
}

func UpdateTable() {
    // defer中决定提交还是回滚
	defer func() {
		if IsPanic() {
			// Rollback transaction
		} else {
			// Commit transaction
		}
	}()

	// Database update operation...
}
```

`func IsPanic() bool` 用来接收异常，返回值用来说明是否发生了异常。`func UpdateTable()` 函数中，使用 defer 来判断最终应该提交还是回滚。

上面代码初步看起来还算合理，但是此处的 `IsPanic()` 再也不会返回 `true`，不是 `IsPanic()` 函数的问题，而是其调用的位置不对。

## 失效的条件

上面代码 `IsPanic()` 失效了，其原因是违反了 recover 的一个限制，导致 recover() 失效（永远返回 `nil` ）。

以下三个条件会让 recover() 返回 `nil` :

1. panic 时指定的参数为 `nil` ；(一般 panic 语句如 `panic( " xxx failed... " )`)
2. 当前协程没有发生 panic ；
3. recover 没有被 defer 方法**直接调用**；

前两条都比较容易理解，上述例子正是匹配第 3 个条件。

本例中，recover() 调用栈为“ defer （匿名）函数” --> IsPanic() --> recover()。也就是说，recover 并没有被 defer 方法直接调用。符合第 3 个条件，所以 recover() 永远返回 nil。

```go

```
