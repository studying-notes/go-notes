---
date: 2020-07-20T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "用 fsnotify 监听文件系统事件"  # 文章标题
url:  "posts/go/libraries/tripartite/fsnotify"  # 设置网页链接，默认使用文件名
tags: [ "go", "fsnotify" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

# fsnotify 监听文件系统事件

开源库 fsnotify 是用 Go 语言编写的跨平台文件系统监听事件库，常用于文件监听，因此我们可以借助该库来实现这个功能。

```
go get -u github.com/fsnotify/fsnotify
```
