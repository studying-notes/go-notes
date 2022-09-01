---
date: 2021-01-02T16:12:26+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Go 格式化占位符"  # 文章标题
url:  "posts/go/libraries/standard/fmt/placeholder"  # 设置网页永久链接
tags: [ "Go", "placeholder" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## 全局

```
%T : 变量的类型信息
%v : 变量的地址
```

```
%w : 包装错误
```

```go
fmt.Errorf("warp: %w", err)
```

## 指针类型

```
%p : 带 0x 的指针
%#p: 不带 0x 的指针
```

## 布尔类型

```
%t : bool,布尔型
```

## 整型

```
%d : 整数
%0nd : 规定输出长度为n的整数，其中开头的数字 0 是必须的，如果整数长度小于n，则用0补齐
%b : 2进制形式
%o : 8进制形式
%x : 16进制形式，小写
%X : 16进制形式，大写
```

```
\ : 后面紧跟长度为3的8进制数
\x : 后面紧跟长度为2的16进制数
\u : 后面紧跟长度为4的16进制数
\U : 后面紧跟长度为8的16进制数
```

## 浮点型

```
%f : 浮点型，默认保留6位小数
%.nf : 浮点型，保留n位小数
%e : 科学计数表示法
%.ne : 科学计数表示法，保留n位小数
%g : 浮点型,用最少的数字表示这个值
%.ng : 最多用n位数字表示这个值,默认浮点数形式，当整数部分长度大于n时，采用科学计数法形式
```

## 字符串

```
%s : 字符串
%q : 字符串带双引号
%#q : 字符串带反引号,字符串本身还有反引号时，则改为字符串带双引号
%x : 将字符串转换为小写的16进制格式
%X : 将字符串转换为大写的16进制格式
% x : 带空格的小写的16进制格式
% X : 带空格的小写的16进制格式
%c : 字符
```