---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 安装与配置指南"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/doc/install"  # 设置网页永久链接
tags: [ "go", "config" ]  # 标签
series: [ "Go 学习笔记"]  # 系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## Ubuntu

### 添加官网源获取最新版

```shell
sudo add-apt-repository ppa:longsleep/golang-backports
# 或者省略，但安装的不会是最新版本
```

```shell
sudo apt-get update
sudo apt-get install golang-go
```

```shell
go version
```

### 手动卸载旧版本

```shell
rm -rf /usr/bin/go
rm -rf /usr/local/go
apt purge golang-go
```

### 手动安装最新版

这种方法对于 WSL 不太好。

[Download Page](https://studygolang.com/dl)

```shell
mkdir /usr/local/go
chmod -R 0777 /usr/local/go
```

```shell
wget https://studygolang.com/dl/golang/go1.13.7.linux-amd64.tar.gz
```

```shell
wget https://studygolang.com/dl/golang/go1.15.linux-amd64.tar.gz
```

```shell
tar -C /usr/local -xvf go1.15.linux-amd64.tar.gz
```

- Bash

```shell
export PATH=/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:/root/bin
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile && source ~/.profile
```

> WSL 十分特殊，PATH 可能带上 Windows 的 PATH，因此添加之前先重置 PATH。

- Fish

```shell
export PATH=/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin:/root/bin
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.config/fish/config.fish && source ~/.config/fish/config.fish
```

```shell
go env -w GOBIN=/usr/local/go/bin
```

## Go 模块代理

```shell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off
go env -w GO111MODULE=on
```
