---
date: 2022-07-19T10:33:41+08:00
author: "Rustle Karl"

title: "go install 命令详解"
url:  "posts/go/quickstart/cli/install"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

```shell
go help install
```

```shell
go install [build flags] [packages]
```

编译并安装源码包。

实际上 `go install` 与 `go build`、`go run` 命令在功能上相差不多，最大的区别是 `go install` 命令会**将编译后的相关文件（如可执行的二进制文件、归档文件等）安装到指定的目录中**。

对于可执行文件，会被安装到 GOBIN （默认是 GOPATH/bin 或 $HOME/go/bin） 目录下。如果未设置 GOBIN 环境变量，则会被安装到 $GOROOT/bin 或者 $GOTOOLDIR。

无法直接指定安装路径。

## run

编译并马上运行程序，它可接收一个或多个文件参数。但只接收 main 包下的文件作为参数，如果不是 main 包下的文件，则会出错。

在执行 go run 命令后，所编译的二进制文件最终存放在一个临时目录里。可以通过 -n 或 -x 参数进行查看。

这两个参数的作用是打印编译过程中的所有执行命令，-n 参数不会继续执行编译后的二进制文件，而 -x 参数会继续执行编译后的二进制文件。

编译器执行了绝大部分编译相关的工作，过程如下：

![06tX5j.png](https://dd-static.jd.com/ddimg/jfs/t1/49489/6/20961/11562/630dc1a3E73ebaa32/dd2b381e06063808.png)

- 创建编译依赖所需的临时目录。Go 编译器会设置一个临时环境变量 WORK，用于在此工作区编译应用程序，执行编译后的二进制文件，其默认值为系统的临时文件目录路径。可以通过设置 GOTMPDIR 来调整其执行目录。
- 编译和生成编译所需要的依赖。该阶段将会编译和生成标准库中的依赖（如 flag.a、log.a、net/http等）、应用程序中的外部依赖（如 gin-gonic/gin 等），以及应用程序自身的代码，然后生成、链接对应归档文件（.a 文件）和编译配置文件。
- 创建并进入编译二进制文件所需的临时目录。即创建 exe 目录。
- 生成可执行文件。这里主要用到的是 link 工具，该工具读取依赖文件的 Go 归档文件或对象及其依赖项，最终将它们组合为可执行的二进制文件。

![06tOaQ.png](https://dd-static.jd.com/ddimg/jfs/t1/155194/35/24206/28953/630dc1a3Ed2680a74/65b045817a21d7a0.png)

- 执行可执行文件。到先前指定的目录 $WORK/b001/exe/main 下执行生成的二进制文件。

在执行 go run 命令后，除非设置了-work 参数，否则会在应用程序结束时自动删除该目录下的相关临时文件。


```shell

```

```shell

```
