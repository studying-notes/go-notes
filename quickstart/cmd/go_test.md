---
date: 2022-03-06T13:12:36+08:00
author: "Rustle Karl"

title: "go test 基础用法"
url:  "posts/go/quickstart/cmd/go_test"  # 永久链接
tags: [ "go", "README" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

toc: true  # 目录
draft: false  # 草稿
---

## 运行指定测试文件

```shell
go test xxx_test.go
```

## 运行测试文件中的某个函数

```shell
go test -run TestXXX
```
