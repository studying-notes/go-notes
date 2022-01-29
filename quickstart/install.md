---
date: 2020-08-30T21:06:02+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 安装与配置指南"  # 文章标题
description: "纸上得来终觉浅，学到过知识点分分钟忘得一干二净，今后无论学什么，都做好笔记吧。"
url:  "posts/go/quickstart/install"  # 设置网页永久链接
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

与大多数开源软件一样，Go 安装包也分为二进制包、源码包。二进制包为基于源码编译出各个组件，并把这些组件打包在一起供人下载和安装，源码包为 Golang 语言源码，供人下载、编译后自行安装。

Go 语言安装比较简单，大体上分为三个步骤：

* 安装可执行命令
* 设置 PATH 环境变量
* 设置 GOPATH 环境变量

## Ubuntu

### 添加官网源获取最新版

> 这种方法最方便，不必手动设置环境变量。

```shell
sudo add-apt-repository ppa:longsleep/golang-backports
# 或者省略，但安装的不会是最新版本
```

```shell
apt update

apt update && apt upgrade -y

apt install golang-go -y
```

```shell
go version
```

### 通过 snap 安装

```shell
snap install --classic go
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

二进制安装包中包含二进制、文档、标准库等内容，我们需要将该二进制完整的解压出来。一般使用 `/usr/local/go` 来存放解压出来的文件，这个目录也就是 `GOROOT`，即 Go 的根目录。

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

#### 设置 PATH 和 GOPATH

Go 的二进制可执行文件存在于 ` $ GOROOT/bin` 目录，需要将该目录加入到 `PATH` 环境变量中。

比如，把下面语句放入 `/etc/profile` 文件中。

```shell
export PATH=$PATH:/usr/local/go/bin
```

Linux 下，自 Go 1.8 版本起，默认把 ` $ HOME/go` 作为 `GOPATH` 目录，可以根据需要设置自已的 `GOPATH` 目录。

`GOPATH` 值不可以与 `GOROOT` 相同，因为如果用户的项目与标准库重名会导致编译时产生歧义。

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

## CentOS

wget https://dl.google.com/go/go1.17.6.linux-amd64.tar.gz


```shell
tar -xzvf go1.17.6.linux-amd64.tar.gz
```

```shell
echo "
export GOROOT=/usr/local/go
export GOPATH=/home/fujiawei/gopath
export GOBIN=$GOROOT/bin
export PATH=$PATH:$GOROOT/bin
export PATH=$PATH:$GOPATH/bin
" >> /etc/profile
```

```shell
source /etc/profile
```

```shell
echo "
export GOROOT=/home/fujiawei/go
export GOPATH=/home/fujiawei/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
" > ~/.config/fish/config.fish && source ~/.config/fish/config.fish
```

```shell
echo "
export GOROOT=/home/fujiawei/go
export GOPATH=/home/fujiawei/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
" >> /home/fujiawei/.bashrc && source /home/fujiawei/.bashrc
```

```shell

```