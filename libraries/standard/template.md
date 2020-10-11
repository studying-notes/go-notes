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

## 入门示例

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

```go

```

```go

```

```go

```

```go

```

```go

```

```go

```

