---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"

title: "Go 安装与卸载"  # 文章标题
url:  "posts/go/quickstart/install"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

## Go 模块代理

```shell
go env -w GOPROXY=https://goproxy.cn,direct
```

## 安装

### Windows

```shell
choco install golang
```

### Ubuntu

官方源现在更新很及时：

```shell
apt install -y golang-go
```

```shell
go version
```

或者通过 snap 安装：

```shell
snap install --classic go
```

不推荐手动下载安装，环境变量等配置管理起来不方便。
