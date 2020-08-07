# Go 语言 JSON 技巧

- [Go 语言 JSON 技巧](#go-语言-json-技巧)
	- [基本的序列化](#基本的序列化)
	- [结构体 Tag](#结构体-tag)
	- [指定字段名](#指定字段名)
	- [忽略某个字段](#忽略某个字段)
	- [忽略空值字段](#忽略空值字段)
	- [忽略嵌套结构体空值字段](#忽略嵌套结构体空值字段)
	- [不修改原结构体忽略空值字段](#不修改原结构体忽略空值字段)
	- [优雅处理字符串格式的数字](#优雅处理字符串格式的数字)
	- [整数变浮点数](#整数变浮点数)
	- [自定义解析时间字段](#自定义解析时间字段)

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

## 结构体 Tag

Tag 是结构体的元信息，可以在运行的时候通过反射的机制读取出来。 Tag 在结构体字段的后方定义，由一对反引号包裹起来，格式如下：

```go
`key1:"value1" key2:"value2"`
```

结构体 tag 由一个或多个键值对组成。键与值使用冒号分隔，值用双引号括起来。同一个结构体字段可以设置多个键值对 tag，不同的键值对之间使用空格分隔。

## 指定字段名

序列化与反序列化默认情况下使用结构体的字段名，我们可以通过给结构体字段添加 tag 来指定 JSON 序列化生成的字段名。

```go
type Box struct {
	Length string `json:"length"`
	Width  int64
	Height float64
}
```

## 忽略某个字段

```go
type Box struct {
	Length string `json:"length"`
	Width  int64
	Height float64 `json:"-"`
}
```

## 忽略空值字段

当 struct 中的字段没有值时，` json.Marshal()` 序列化的时候不会忽略这些字段，而是默认输出字段的类型零值。如果想要在序列序列化时忽略这些没有值的字段时，可以在对应字段添加 `omitempty tag`。

```go
type Box struct {
	Length string `json:"length"`
	Width  int64 `json:"width,omitempty"`
	Height float64 `json:"-"`
}
```

## 忽略嵌套结构体空值字段

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

匿名嵌套结构体序列化后的 JSON 串为单层的：

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

## 不修改原结构体忽略空值字段

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

## 优雅处理字符串格式的数字

在结构体 tag 中添加 string 来告诉 JSON 包从字符串中解析相应字段的数据：

```go
type Card struct {
	ID    int64   `json:"id,string"`    // 添加 string tag
	Score float64 `json:"score,string"` // 添加 string tag
}
```

## 整数变浮点数

在 JSON 协议中是没有整型和浮点型之分的，它们统称为 number。JSON 字符串中的数字经过 Go 语言中的 JSON 包反序列化之后都会成为 float64 类型。

```go
func main() {
	var m1 = make(map[string]interface{}, 1)
	m1["count"] = 1
	buf, _ := json.Marshal(m1)
	fmt.Printf("%s\n", buf)  // {"count":1}

	var m2 map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(buf))
	decoder.UseNumber()
	_ = decoder.Decode(&m2)
	fmt.Printf("%T\n", m2["count"])  // json.Number
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
