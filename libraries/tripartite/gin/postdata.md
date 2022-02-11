---
date: 2020-11-29T19:44:04+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "form-data、x-www-form-urlencoded 的区别"  # 文章标题
# description: "文章描述"
url:  "posts/gin/abc/postdata"  # 设置网页永久链接
tags: [ "gin", "http"]  # 标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## form-data

即 http 请求中的 `multipart/form-data`, 它会将表单的数据处理为一条消息，以标签为单元，用分隔符分开。既可以上传键值对，也可以上传文件。

当上传的字段是文件时，会有 `Content-Type` 来说明文件类型； `content-disposition`，用来说明字段的一些信息；由于有 boundary 隔离，所以 `multipart/form-data` 既可以上传文件，也可以上传键值对，它采用了键值对的方式，所以可以上传多个文件。

## x-www-form-urlencoded

即 `application/x-www-form-urlencoded`，会将表单内的数据转换为键值对。


- `multipart/form-data`：既可以上传文件等二进制数据，也可以上传表单键值对，只是最后会转化为一条信息
- `x-www-form-urlencoded`：只能上传键值对，并且键值对都是间隔分开的
