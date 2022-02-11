---
date: 2020-12-11T08:43:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Data URI scheme"  # 文章标题
url:  "posts/gin/abc/scheme"  # 设置网页永久链接
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

网页上有些图片的 src 或 css 背景图片的 url 后面跟了一大串字符，比如：

```base64
data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAkAAAAJAQMAAADaX5RTAAAAA3NCSVQICAjb4U/gAAAABlBMVEX///+ZmZmOUEqyAAAAAnRSTlMA/1uRIrUAAAAJcEhZcwAACusAAArrAYKLDVoAAAAWdEVYdENyZWF0aW9uIFRpbWUAMDkvMjAvMTIGkKG+AAAAHHRFWHRTb2Z0d2FyZQBBZG9iZSBGaXJld29ya3MgQ1M26LyyjAAAAB1JREFUCJljONjA8LiBoZyBwY6BQQZMAtlAkYMNAF1fBs/zPvcnAAAAAElFTkSuQmCC
```

那么这是什么呢？这是 Data URI scheme。

Data URI scheme 是在 RFC2397 中定义的，目的是将一些小的数据，直接嵌入到网页中，从而不用再从外部文件载入。比如上面那串字符，其实是一张小图片，将这些字符复制黏贴到火狐的地址栏中并转到，就能看到它了，一张 1X36 的白灰 png 图片。

在上面的 Data URI 中，data 表示取得数据的协定名称，image/png 是数据类型名称，base64 是数据的编码方法，逗号后面就是这个 image/png 文件 base64 编码后的数据。

 目前，Data URI scheme 支持的类型有：

- data:,文本数据
- data:text/plain,文本数据
- data:text/html,HTML 代码
- data:text/html;base64,base64 编码的 HTML 代码
- data:text/css,CSS 代码
- data:text/css;base64,base64 编码的 CSS 代码
- data:text/javascript,javascript 代码
- data:text/javascript;base64,base64 编码的 javascript 代码
- data:image/gif;base64,base64 编码的 gif 图片数据
- data:image/png;base64,base64 编码的 png 图片数据
- data:image/jpeg;base64,base64 编码的 jpeg 图片数据
- data:image/x-icon;base64,base64 编码的 icon 图片数据

把图像文件的内容直接写在了 HTML 文件中，这样做的好处是，节省了一个 HTTP 请求。坏处是浏览器不会缓存这种图像。
