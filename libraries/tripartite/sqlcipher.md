---
date: 2020-07-20T21:39:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "用 SQLCipher 加密 SQLite"  # 文章标题
description: "Windows 下加密 SQLite 的准备工作相当繁琐"
url:  "posts/go/libraries/tripartite/sqlcipher"  # 设置网页链接，默认使用文件名
tags: [ "go", "sqlcipher", "sqlite"]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
toc: true  # 是否自动生成目录
---

## 预准备

### 安装 Perl 64bit

有不同的发行版本，随便选一个即可，比如：

```
http://strawberryperl.com/download/5.30.2.1/strawberry-perl-5.30.2.1-64bit.msi
```

安装到默认或者指定路径，然后添加 `bin` 目录到 `PATH` 中。

### 安装 GCC 64bit

也有很多选择，比如 TDM-GCC-64：

```
https://jmeubank.github.io/tdm-gcc/
```

### 安装 MSYS

```
http://downloads.sourceforge.net/mingw/MSYS-1.0.11.exe
```

不可缺少，缺了，一些命令无法在 Windows 下无法运行。

> 建议改用 MSYS2，MSYS2 是 MSYS 的一个升级版, 准确的说是集成了 pacman 和 MinGW-w64 的 Cygwin 升级版, 提供了 bash shell 等 Linux 环境、版本控制软件和 MinGW-w64 工具链等。

### 安装 Make

官网

```
https://cmake.org
```

无论如何搞定，最后可以提供 `make` 命令即可，`Powershell` 下可以以管理员身份执行以下命令安装：

```powershell
choco install make
```

另一种：

```powershell
copy mingw32-make.exe make.exe

copy c:/developer/tdm64-gcc/bin/mingw32-make.exe c:/developer/tdm64-gcc/bin/make.exe
```

启动一个终端，运行以下命令，都可以正确显示就表明准备工作就绪了：

```powershell
perl -v
gcc -v
make -v
```

## 编译 OpenSSL 64bit && 32bit

1. 下载 OpenSSL

```
https://www.openssl.org/source/
```

新版本编译起来更复杂，这里用的旧版本：

```
https://www.openssl.org/source/openssl-1.0.2k.tar.gz
```

2. 解压

```powershell
tar -xvzf openssl-1.0.2k.tar.gz
```

用 360 解压缩这一类工具也没问题。

3. 编译

```powershell
cd openssl-1.0.2k

# x64
perl configure mingw64 no-shared no-asm

# x86 
perl configure mingw no-shared no-asm

# 在 Makefile 中添加参数
# CFLAGS = -m32
# LDFLAGS = -m32
make
```

编译过程大概十几分钟。

编译完成后，将 `openssl-1.0.2k` 目录下的两个文件 `libcrypto.a`、`libcrypto.pc` 复制到 `TDM-GCC-64\lib` 目录下，然后将 `openssl-1.0.2k\include\openssl` 这个文件夹复制到 `TDM-GCC-64\x86_64-w64-mingw32\include` 下。

## 安装 Go-SQLCipher 库

1. 首先下载源码

```powershell
go get github.com/xeodou/go-sqlcipher
```

不出意外，安装报错，但源码已经下载下来了，找到所在目录，一般在：

```
GOPATH\pkg\mod\github.com\xeodou
```

这里我遇到了一个问题，全局 GOPATH 和 Goland 的 GOPATH 设置的不一样，导致 Goland 没有问题而命令行运行一直出错。

2. 修改 `sqlite3_windows.go` 文件

打开一级目录下的  `sqlite3_windows.go` 文件，改成以下内容：

```go
// +build windows

package sqlite3

/*
#cgo CFLAGS: -I. -fno-stack-check -fno-stack-protector -mno-stack-arg-probe
#cgo windows,386 CFLAGS: -D_USE_32BIT_TIME_T
#cgo LDFLAGS: -lmingwex -lmingw32 -lgdi32
*/
import "C"
```

3. 最后再次安装

```powershell
go install -v .
```

完成！

### 附个人 fork 修改版

上面的方法在 Windows 下过于麻烦，因此我把改好的 fork 仓库上传了，这样就不用每次都修改文件了。

```powershell
go get github.com/fujiawei-dev/go-sqlcipher 
```

## 启用 CGO

```shell
GOARCH=386;CGO_ENABLED=1
```

## 创建加密数据库

```go
package main

import (
	"database/sql"
	_ "github.com/fujiawei-dev/go-sqlcipher"
)

func func main() {
    sql.Open("sqlite3", databasefile +"?_key=password")
}
```

经测试，现在只支持自定义密钥 `key`，其他设置了都无效，最终都会变成默认设置。

## 打开加密数据库

折腾半天，走了很多错误的方向，一是 Windows 平台找不到编译好的最新版本，找了个旧版本（当时还没发现这个问题），怎么也打不开生成的文件；二是根据 StackOverflow 的回答找了一个专门打开 SQLite 的可视化工具（SQLiteStudio），然而必须自己输除密钥外的其他设置，怎么输都错。于是我搜索寻找默认设置，然而从官方文档到源代码，只找到一丝眉目。

几近绝望，最终在一篇讲微信加密数据库的博客里发现了一款工具：

```shell
# 官网
https://download.sqlitebrowser.org/DB.Browser.for.SQLite-3.12.0-win64.msi

# 蓝奏云
https://wwa.lanzous.com/itiH8gqi0pa
```

提供了新旧版本的默认设置，问题解决了。
