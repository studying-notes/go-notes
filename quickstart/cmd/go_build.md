---
date: 2022-05-22T18:39:22+08:00
author: "Rustle Karl"

title: "go build 命令详解"
url:  "posts/go/docs/cmd/go_build"  # 永久链接
tags: [ "Go", "README" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

toc: true  # 目录
draft: false  # 草稿
---

## 语法

```shell
go build [-o 输出名] [-i] [编译标记] [包名]
```

- 如果参数为 xxx.go 文件或文件列表，则编译为一个个单独的包。
- 当编译单个 main 包（文件），则生成可执行文件。
- 当编译单个或多个非主包时，只构建编译包，但丢弃生成的对象（.a 文件），仅用作检查包可以构建。
- 当编译包时，会自动忽略 `_test.go` 的测试文件。

## 选项

```shell
-a 强制重新编译
-n 打印命令但是不执行
-p n 编译时并行数，默认 CPU 数量（GOMAXPROCS）
-race 启用数据竞争状态检测
-msan 启用与内存清理程序的互操作
-v 打印出被编译的包名
-work 打印临时工作目录的名称，并在退出时不删除它
-x 打印命令
-asmflags 'flag list' 传递每个 go 工具 asm 调用的参数
-buildmode mode 编译模式 `go help buildmode` 可以看到所有模式
-buildvcs 是否使用版本控制信息标记二进制文件。 
-compiler name 编译器 (gccgo or gc)
-gccgoflags 'arg list' gccgo 编译/链接器参数
-gcflags 'arg list' 垃圾回收参数
-installsuffix suffix 安装目录
-ldflags 'flag list'
    '-s -w': 压缩编译后的体积
    -s: 去掉符号表
    -w: 去掉调试信息，不能gdb调试了
-linkshared 链接到以前使用 -buildmode=shared 创建的共享库
-pkgdir dir 从指定位置，而不是通常的位置安装和加载所有软件包
-tags 'tag list' 构建出带 tag 的版本，逗号分隔的列表，旧版本是空格分隔
# https://stackoverflow.com/questions/45279385/remove-file-paths-from-text-directives-in-go-binaries
-trimpath 从生成的可执行文件中删除所有文件系统路径，否则包含编译时候项目所在的绝对路径，报错崩溃的时候就会显示出来
-toolexec 'cmd args'
```
