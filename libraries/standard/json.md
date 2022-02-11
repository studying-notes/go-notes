---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "json - JSON 序列化和反序列化"  # 文章标题
url:  "posts/go/libraries/standard/json"  # 设置网页链接，默认使用文件名
tags: [ "go", "json" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

- [处理 PUT 请求的反序列化问题](#处理-put-请求的反序列化问题)
- [序列化列表](#序列化列表)
- [Go 和 JSON 转换关系](#go-和-json-转换关系)
- [基本的序列化](#基本的序列化)
- [基本的反序列化](#基本的反序列化)
- [嵌套结构体的序列化](#嵌套结构体的序列化)
- [解析不知道格式的数据](#解析不知道格式的数据)
- [结构体标签](#结构体标签)
	- [指定字段名](#指定字段名)
	- [忽略某个字段](#忽略某个字段)
	- [忽略空值字段](#忽略空值字段)
	- [忽略嵌套结构体空值字段](#忽略嵌套结构体空值字段)
	- [不修改原结构体忽略空值字段](#不修改原结构体忽略空值字段)
	- [优雅处理字符串格式的数字](#优雅处理字符串格式的数字)
- [整数变浮点数](#整数变浮点数)
- [自定义解析时间字段](#自定义解析时间字段)

## 处理 PUT 请求的反序列化问题

首先，反对用指针。

在 Go 标准库的反序列化中，结构体和 JSON 中单方面存在的字段都自动忽略。

所以推荐做法是，先从数据库中取出该条数据，然后在该条数据的基础上反序列化。

对应的接口协议设计就要求将对象的 ID 放在 path/query 参数里，修改数据放 body 里，而不是将修改对象的 ID 也放 body 里。

```go
package main

import (
	"encoding/json"
)

type Folder struct {
	Path  string   `json:"path"`
	Count int      `json:"count"`
	UID   string   `json:"uid"`
	Items []string `json:"items"`
}

func main() {
	f := &Folder{Path: "path", Count: 1, UID: "uid", Items: []string{"1", "2"}}

	data := []byte(`{"path": "", "count": 0}`)

	_ = json.Unmarshal(data, &f)
}
```

## 序列化列表

偶尔会迷糊。。。

```go
import (
	"encoding/json"
	"fmt"
)

func main() {
	v := []string{"1", "2", "3"}
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))

	var val []string
	err = json.Unmarshal(buf, &val)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
```

## Go 和 JSON 转换关系

| Go 类型 | JSON 类型 |
| ------- | --------- |
| bool | booleans |
| float64 | numbers |
| string | strings |
| nil | null |

## 基本的序列化

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Box struct {
	Length string
	Width  int64
	Height float64
}

func main() {
	b1 := Box{
		Length: "120",
		Width:  75,
		Height: 4,
	}
	// b1 可以是指针
	buf, _ := json.Marshal(b1)
	fmt.Printf("%s\n", buf)

	var b2 Box
	_ = json.Unmarshal(buf, &b2)
	fmt.Printf("%+v\n", b2)
}
```

```
{"Length":"120","Width":75,"Height":4}
{Length:120 Width:75 Height:4}
```

## 基本的反序列化

结构体和 JSON 中单方面存在的字段都自动忽略。

## 嵌套结构体的序列化

```go
func main() {
	type Thing struct {
		Length int `json:"length"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type Person struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Parents []string `json:"parents"`
		Thing   `json:"thing"`
	}

	person := Person{Name: "Wetness", Age: 18,
		Parents: []string{"Gomez", "Morita"},
		// 类型字段也可以用于赋值，不用定义变量
		Thing: Thing{2, 2, 2}}

	fmt.Printf("%#v\n\n", person)

	buf, _ := json.Marshal(person)
	fmt.Printf("%s\n", buf)
}
```

```
main.Person{Name:"Wetness", Age:18, Parents:[]string{"Gomez", "Morita"}, Thing:main.Thing{Length:2, Width:2, Height:2}}

{"name":"Wetness","age":18,"parents":["Gomez","Morita"],"thing":{"length":2,"width":2,"height":2}}
```

## 解析不知道格式的数据

场景是 MQTT 服务，订阅 Topic 相同，但是消息的格式却是有多种，通过 `type` 字段区分。

```go
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Parents []string `json:"parents"`
	Thing   `json:"thing"`
}

type Thing struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func main() {
	b := []byte(`{"type":"1","content":{"name":"Wetness","age":18,"parents":["Gomez","Morita"],"thing":{"length":1,"width":2,"height":3}}}`)

	model := make(map[string]interface{})
	_ = json.Unmarshal(b, &model)
	t := model["type"].(string) // 根据类型判断解析
	switch t {
	case "1":
		//p := Message{Content: Person{}}
		//_ = json.Unmarshal(b, &p)
		//fmt.Printf("%#v\n\n", p.Content)
		// 不是期望得到的
		// map[string]interface {}{"age":18, "name":"Wetness", "parents":[]interface {}{"Gomez", "Morita"}}

		p2 := struct {
			Type    string `json:"type"`
			Content Person `json:"content"`
		}{}
		_ = json.Unmarshal(b, &p2)
		fmt.Printf("%#v\n\n", p2.Content)
		// main.Person{Name:"Wetness", Age:18, Parents:[]string{"Gomez", "Morita"}, Thing:main.Thing{Length:1, Width:2, Height:3}}
	}
}
```

> 由此可见，interface{} 类型的无法反序列化回来。

## 结构体标签

标签 Tag 是结构体的元信息，可以在运行的时候通过反射的机制读取出来。 Tag 在结构体字段的后方定义，由一对反引号包裹起来，格式如下：

```
`key1:"value1" key2:"value2"`
```

结构体 tag 由一个或多个键值对组成。**键与值使用冒号分隔**，**值用双引号括起来**。同一个结构体字段可以设置多个键值对 tag，**不同的键值对之间使用空格分隔**。

### 指定字段名

序列化与反序列化默认情况下使用结构体的字段名，我们可以通过给结构体字段添加 tag 来指定 JSON 序列化生成的字段名。

```go
type Box struct {
  Length string `json:"length" db:"length"`
	Width  int64  `json:"width"`
	Height float64
}
```

### 忽略某个字段

```go
type Box struct {
	Length string `json:"length"`
	Width  int64
	Height float64 `json:"-"`
}
```

### 忽略空值字段

当 struct 中的**字段没有值时**，` json.Marshal()` 序列化的时候不会忽略这些字段，而是**默认输出字段的类型零值**。如果想要在序列序列化时忽略这些没有值的字段时，可以在对应字段添加 `omitempty` 标签。

```go
type Box struct {
	Length string `json:"length"`
	Width  int64 `json:"width,omitempty"`
	Height float64 `json:"-"`
}
```

### 忽略嵌套结构体空值字段

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Box struct {
	Length string  `json:"length"`
	Width  int64   `json:"width,omitempty"`
	Height float64 `json:"-"`
	Things
}

type Things struct {
	Weight []float64 `json:"weight"`
	Logo   string    `json:"logo"`
}

func main() {
	b1 := Box{
		Length: "120",
		Width:  75,
		Height: 4,
	}
	buf, _ := json.Marshal(b1)
	fmt.Printf("%s\n", buf)

	var b2 Box
	_ = json.Unmarshal(buf, &b2)
	fmt.Printf("%+v\n", b2)
}
```

**匿名嵌套结构体**序列化后的 JSON 串为**单层**：

```
{"length":"120","width":75,"weight":null,"logo":""}
```

为了还原嵌套形式，必须改为具名嵌套或者定义字段 Tag：

```go
type Box struct {
	Length string  `json:"length"`
	Width  int64   `json:"width,omitempty"`
	Height float64 `json:"-"`
	Things `json:"things"`
}
```

```
{"length":"120","width":75,"things":{"weight":null,"logo":""}}
```

在嵌套的结构体为空值时，忽略该字段，仅添加 `omitempty` 是不够的，必须用嵌套结构体的指针：

```go
type Box struct {
	Length  string  `json:"length"`
	Width   int64   `json:"width,omitempty"`
	Height  float64 `json:"-"`
	*Things `json:"things,omitempty"`
}
```

### 不修改原结构体忽略空值字段

可以使用创建另外一个结构体匿名嵌套原结构体，同时指定字段为匿名结构体指针类型，并添加omitemptytag，示例代码如下：

```go
type Box struct {
	*Things           // 匿名嵌套
	Logo    *struct{} `json:"logo,omitempty"`
}

type Things struct {
	Weight []float64 `json:"weight"`
	Logo   string    `json:"logo"`
}
```

### 优雅处理字符串格式的数字

在结构体 tag 中添加 string 来告诉 JSON 包从字符串中解析相应字段的数据：

```go
type Card struct {
	ID    int64   `json:"id,string"`    // 添加 string tag
	Score float64 `json:"score,string"` // 添加 string tag
}
```

## 整数变浮点数

在 JSON 协议中是没有整型和浮点型之分的，它们统称为 `number`。JSON 字符串中的数字经过 Go 语言中的 JSON 包反序列化之后都会成为 `float64` 类型。

```go
func main() {
	var m1 = make(map[string]interface{}, 1)
	m1["count"] = 1
	buf, _ := json.Marshal(m1)
	fmt.Printf("%s\n", buf)  // {"count":1}

	var m2 map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(buf))
	// as a interface{} instead of as a float64
	decoder.UseNumber()

	_ = decoder.Decode(&m2)
	fmt.Printf("%T\n", m2["count"])  // json.Number

	// 类型转换
	count, _ := m2["count"].(json.Number).Int64()
	fmt.Printf("%T\n", int(count))  // int
}
```

## 自定义解析时间字段

```go
type Post struct {
	PostTime time.Time `json:"post_time"`
}

func main() {
	p1 := Post{
		PostTime: time.Now(),
	}
	buf, _ := json.Marshal(p1)
	fmt.Printf("%s\n", buf)

	s := `{"post_time":"2020-07-18 14:16:32"}`
	var p2 Post
	_ = json.Unmarshal([]byte(s), &p2)
	fmt.Printf("%+v\n", p2)
}
```

```
{"post_time":"2020-07-18T14:17:53.3245572+08:00"}
{PostTime:0001-01-01 00:00:00 +0000 UTC}
```

无法识别常用的字符串时间格式。

自定义事件格式解析，见：

```
https://zhuanlan.zhihu.com/p/158873918
```
