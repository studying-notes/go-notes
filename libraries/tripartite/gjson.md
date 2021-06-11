---
date: 2021-01-20T15:22:22+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "GJSON 快速提取 JSON 值"  # 文章标题
url:  "posts/go/libraries/tripartite/gjson"  # 设置网页链接，默认使用文件名
tags: [ "go", "json" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```go
package main

import "github.com/tidwall/gjson"

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {
	value := gjson.Get(json, "name.last")
	println(value.String())
}
```

```go

```


```go

```


