---
date: 2022-07-19T10:26:22+08:00
author: "Rustle Karl"

title: "Go 命令行简介"
url:  "posts/go/quickstart/cli/README"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

显示帮助：

```shell
go help command
```

## packages

```shell
go help packages
```

大部分命令都作用于 packages 集合：

```shell
go action [packages]
```

通常，[packages] 是一个导入路径列表。导入路径是作为根路径或以 . 或 .. 符号可被解释为文件系统路径并表示该目录中的包。

如果没有给出导入路径，则表示当前目录。

有 4 个保留名称：

- main 表示独立可执行文件的顶级包
- all GOPATH 中的全部包
- std 全部标准库
- cmd 内置命令和库

`cmd/` 起始的导入路径只匹配 Go 仓库源码。

如果导入路径包含一个或多个“...”通配符，则它是一种模式，每个都可以匹配任何字符串，包括空字符串和包含斜杠的字符串。 这样的模式扩展到所有包在 GOPATH 树中找到名称匹配的目录模式。

但存在两个特殊情况：

-  /... 可以匹配空字符串，所以 net/... 匹配 net 及其子目录，比如 net/http。
-  任何包含通配符的以斜线分隔的模式元素都不会参与到 vendor 包路径中的 vendor 元素的匹配中，因此 ./... 不匹配 ./vendor 或 ./mycode/ 子目录中的包 vendor，但是 ./vendor/... 和 ./mycode/vendor/... 可以。 但是请注意，一个名为 vendor 且本身包含代码的目录不是一个 vendored 包： cmd/vendor 将是一个名为 vendor 的命令，并且模式 cmd/... 与之匹配。

以 . 或者 _ 开始的目录或文件会被忽略。

名为 testdata 的目录也被忽略。

```shell

```

```shell

```
