---
date: 2020-11-15T12:18:19+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 结构体标签"  # 文章标题
url:  "posts/go/docs/grammar/struct/tag"  # 设置网页永久链接
tags: [ "Go", "struct" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 前言

Go 的 struct 声明允许字段附带 `Tag` 来对字段做一些标记。

该 `Tag` 不仅仅是一个字符串那么简单，因为其主要用于反射场景，`reflect` 包中提供了操作 `Tag` 的方法，所以 `Tag` 写法也要遵循一定的规则。

## Tag 的本质

### Tag 规则

`Tag` 本身是一个字符串，但字符串中却是：` 以空格分隔的 key:value 对 `。

- `key` : 必须是非空字符串，字符串不能包含控制字符、空格、引号、冒号。
- `value` : 以双引号标记的字符串
- 注意：冒号前后不能有空格

如下代码所示

```go
type Server struct {
    ServerName string `key1:"value1" key11:"value11"`
    ServerIP   string `key2:"value2"`
}
```

上述代码 `ServerName` 字段的 `Tag` 包含两个 key-value 对。`ServerIP` 字段的 `Tag` 只包含一个 key-value 对。

### Tag 是 struct 的一部分

前面说过，`Tag` 只有在反射场景中才有用，而反射包中提供了操作 `Tag` 的方法。在说方法前，有必要先了解一下 Go 是如何管理 struct 字段的。

以下是 `reflect` 包中的类型声明，省略了部分与本文无关的字段。

```go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string
	...
	Type      Type      // field type
	Tag       StructTag // field tag string
	...
}

type StructTag string
```

可见，描述一个结构体成员的结构中包含了 `StructTag`，而其本身是一个 `string`。也就是说 `Tag` 其实是结构体字段的一个组成部分。

### 获取 Tag

`StructTag` 提供了 `Get(key string) string` 方法来获取 `Tag`，示例如下：

```go
package main

import (
    "reflect"
    "fmt"
)

type Server struct {
    ServerName string `key1:"value1" key11:"value11"`
    ServerIP   string `key2:"value2"`
}

func main() {
    s := Server{}
    st := reflect.TypeOf(s)

    field1 := st.Field(0)
    fmt.Printf("key1:%v\n", field1.Tag.Get("key1"))
    fmt.Printf("key11:%v\n", field1.Tag.Get("key11"))

    filed2 := st.Field(1)
    fmt.Printf("key2:%v\n", filed2.Tag.Get("key2"))
}
```

程序输出如下：

```
key1:value1
key11:value11
key2:value2
```

## Tag 存在的意义

使用反射可以动态的给结构体成员赋值，正是因为有 tag，在赋值前可以使用 tag 来决定赋值的动作。
比如，官方的 `encoding/json` 包，可以将一个 JSON 数据 `Unmarshal` 进一个结构体，此过程中就使用了 Tag。该包定义一些规则，只要参考该规则设置 tag 就可以将不同的 JSON 数据转换成结构体。

正是基于 struct 的 tag 特性，才有了诸如 json、orm 等等的应用。

## Tag 常见用法

常见的 tag 用法，主要是 JSON 数据解析、ORM 映射等。

```go

```
