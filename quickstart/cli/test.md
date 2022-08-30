---
date: 2022-07-19T10:10:25+08:00
author: "Rustle Karl"

title: "go test 基础用法"
url:  "posts/go/quickstart/cli/test"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 运行指定测试文件

```shell
go test xxx_test.go
```

## 运行测试文件中的某个函数

```shell
go test -run TestXXX
```
