---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 文本模板引擎"  # 文章标题
url:  "posts/go/libraries/standard/template"  # 设置网页链接，默认使用文件名
tags: [ "go", "template", "render" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

template 是 Go 语言的文本模板引擎，它提供了两个标准库。这两个标准库使用了同样的接口，但功能略有不同，具体如下：

- text/template：基于模板输出文本内容。
- html/template：基于模板输出安全的 HTML 格式的内容，可以理解为其进行了转义，以避免受某些注入攻击。

`html/template` 实现了数据驱动的模板，用于生成可防止代码注入的安全的 HTML 内容。它提供了和 `text/template` 相同的接口，Go 语言中输出 HTML 的场景都应使用 `html/template`。

Go 语言内置了文本模板引擎 `text/template` 和用于 HTML 文档的 `html/template`。它们的作用机制可以简单归纳如下：

1. 模板文件通常定义为 `.tmpl` 和 `.tpl` 为后缀，必须用 `UTF8` 编码；
2. 模板文件中使用 `{{ }}` 包裹和标识需要传入的数据；
3. 传给模板的数据可以通过点号 `.` 来访问，如果数据是复杂类型的数据，可以通过 `{{ .FieldName }}` 来访问它的字段；
4. 除 `{{ }}` 包裹的内容外，其他内容均不做修改原样输出。


## 模板引擎的用法

Go 语言模板引擎的使用可以分为三部分：定义模板文件、解析模板文件和模板渲染。

### 定义模板文件

其中，定义模板文件时需要我们按照相关语法规则去编写。

### 解析模板文件

上面定义好了模板文件之后，可以使用下面的常用方法去解析模板文件，得到模板对象：

```go
func (t *Template) Parse(src string) (*Template, error)
func ParseFiles(filenames ...string) (*Template, error)
func ParseGlob(pattern string) (*Template, error)
```

当然，也可以使用 `func New(name string) *Template` 函数创建一个名为 `name` 的模板，然后对其调用上面的方法去解析模板字符串或模板文件。

### 模板渲染

渲染模板简单来说就是使用数据去填充模板，当然实际上可能会复杂很多。

```go
func (t *Template) Execute(wr io.Writer, data interface{}) error
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error
```

## text/template 入门示例

```go
func main() {
	text := "1. {{ title .user }}\n2. {{ .tag | title}}"
	funcMap := template.FuncMap{"title": strings.Title}
	tmpl, _ := template.New("start").Funcs(funcMap).Parse(text)
	data := map[string]string{
		"user": "rustle",
		"tag":  "admin",
	}
	_ = tmpl.Execute(os.Stdout, data)
}
```

首先调用标准库 `text/template` 中的 `New` 方法，其根据我们给定的名称标识创建了一个全新的模板对象。接下来调用 `Parse` 方法，将 `text` 解析为当前文本模板的主体内容。最后调用 `Execute` 方法，进行模板渲染。简单来说，就是将传入的 `data` 动态参数渲染到对应的模板标识位上。因为 `io.Writer` 指定到了 `os.Stdout` 中，所以其最终输出到标准控制台中。

## 模板语法

- 双层大括号：也就是 `{{ }}` 标识符，在 `template` 中，所有的动作（Actions）、数据评估（Data Evaluations）、控制流转都需要用标识符双层大括号包裹，其余的模板内容均全部原样输出。
- 点（DOT）：会根据点（DOT）标识符进行模板变量的渲染，其参数可以为任何值，但特殊的复杂类型需进行特殊处理。例如，当为指针时，内部会在必要时自动表示为指针所指向的值。如果执行结果生成了一个函数类型的值，如结构体的函数类型字段，那么该函数不会自动调用。
- 函数调用：在前面的代码中，通过 `FuncMap` 方法注册了名 `title` 的自定义函数。在模板渲染中一共用了两类处理方法，即使用 `{{title .user}}` 和管道符 `|` 对 `.tag` 进行处理。在 `template` 中，会把管道符前面的运算结果作为参数传递给管道符后面的函数，最终，命令的输出结果就是这个管道的运算结果。

### 示例

1. 定义模板文件

```html
<html lang="zh-CN">
<head>
    <title>Hello</title>
</head>
<body>
    <p>Hello {{.}}</p>
</body>
</html>
```

2. 解析和渲染模板文件

```go
func sayHello(w http.ResponseWriter, r *http.Request) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println("failed:", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入
	tmpl.Execute(w, "World")
}
func main() {
	http.HandleFunc("/", sayHello)
	_ := http.ListenAndServe(":9090", nil)
}
```

## 模板语法

### {{.}}

模板语法都包含在 `{{` `}}` 中间，其中 `{{.}}` 中的点表示当前对象。

当我们传入一个结构体对象时，我们可以根据 `.` 来访问结构体的对应字段。

同理，当我们传入的变量是 Map 时，也可以在模板文件中通过 `.` 根据 `key` 来取值。

### 注释

注释不能嵌套，并且必须紧贴分界符始止。

```html
{{/* a comment */}}
```

### pipeline

`pipeline` 是指产生数据的操作，比如`{{.}}`、`{{.Name}}`等。Go 的模板语法中支持使用管道符号 `|` 链接多个命令，用法和 `unix` 下的管道类似：`|` 前面的命令会将运算结果传递给后一个命令的最后一个位置。

Go 的模板语法中，`pipeline` 的概念是传递数据，只要能产生数据的，都是 `pipeline`。

### 变量

我们还可以在模板中声明变量，用来保存传入模板的数据或其他语句生成的结果。具体语法如下：

```html
$obj := {{.}}
```

其中 `$obj` 是变量的名字，在后续的代码中就可以使用该变量了。

### 移除空格

有时候我们在使用模板语法的时候会不可避免的引入一下空格或者换行符，这样模板最终渲染出来的内容可能就和我们想的不一样，这个时候可以使用 `{{-` 语法去除模板内容左侧的所有空白符号， 使用 `-}}` 去除模板内容右侧的所有空白符号。

```html
{{- .Name -}}
```

`-` 紧挨 `{{`、`}}`，同时与模板值之间需要使用空格分隔。

### 条件判断

Go 模板语法中的条件判断有以下几种:

```html
{{if pipeline}} T1 {{end}}

{{if pipeline}} T1 {{else}} T0 {{end}}

{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
```

### range

Go 的模板语法中使用 `range` 关键字进行遍历，有以下两种写法，其中 `pipeline` 的值必须是数组、切片、字典或者通道。

1. pipeline 的值其长度为 0时，不会有任何输出

```html
{{range pipeline}} T1 {{end}}
```

2. pipeline 的值其长度为 0，则会执行 T0。

```html
{{range pipeline}} T1 {{else}} T0 {{end}}
```

### with

1. pipeline 的值其长度为 0时，不会有任何输出

```html
{{with pipeline}} T1 {{end}}
```

2. pipeline 的值其长度为 0，则会执行 T0。

```html
{{range pipeline}} T1 {{else}} T0 {{end}}
```

### 预定义函数

执行模板时，函数从两个函数字典中查找：首先是模板函数字典，然后是全局函数字典。一般不在模板内定义函数，而是使用 `Funcs` 方法添加函数到模板里。

预定义的全局函数如下：

- and
      函数返回它的第一个empty参数或者最后一个参数；
      就是说"and x y"等价于"if x then y else x"；所有参数都会执行；
- or
      返回第一个非empty参数或者最后一个参数；
      亦即"or x y"等价于"if x then x else y"；所有参数都会执行；
- not
      返回它的单个参数的布尔值的否定
- len
      返回它的参数的整数类型长度
- index
      执行结果为第一个参数以剩下的参数为索引/键指向的值；
      如"index x 1 2 3"返回x[1][2][3]的值；每个被索引的主体必须是数组、切片或者字典。
- print
      即fmt.Sprint
- printf
      即fmt.Sprintf
- println
      即fmt.Sprintln
- html
      返回与其参数的文本表示形式等效的转义HTML。
      这个函数在html/template中不可用。
- urlquery
      以适合嵌入到网址查询中的形式返回其参数的文本表示的转义值。
      这个函数在html/template中不可用。
- js
      返回与其参数的文本表示形式等效的转义JavaScript。
- call
      执行结果是调用第一个参数的返回值，该参数必须是函数类型，其余参数作为调用该函数的参数；
      如"call .X.Y 1 2"等价于go语言里的dot.X.Y(1, 2)；
      其中Y是函数类型的字段或者字典的值，或者其他类似情况；
      call的第一个参数的执行结果必须是函数类型的值（和预定义函数如print明显不同）；
      该函数类型值必须有1到2个返回值，如果有2个则后一个必须是error接口类型；
      如果有2个返回值的方法返回的error非nil，模板执行会中断并返回给调用模板执行者该错误；

### 比较函数

布尔函数会将任何类型的零值视为假，其余视为真。下面是定义为函数的二元比较运算的集合：

- eq - 如果 `arg1 == arg2` 则返回真
- ne - 如果 `arg1 != arg2` 则返回真
- lt - 如果 `arg1 < arg2` 则返回真
- le - 如果 `arg1 <= arg2` 则返回真
- gt - 如果 `arg1 > arg2` 则返回真
- ge - 如果 `arg1 >= arg2` 则返回真

为了简化多参数相等检测，`eq` 可以接受 2 个或更多个参数，它会将第一个参数和其余参数依次比较。

```template
{{eq arg1 arg2 arg3}}
```

比较函数只适用于基本类型。但是，整数和浮点数不能互相比较。

太多了，暂时用不到，放上链接：

```
https://www.liwenzhou.com/posts/Go/go_template/
```
