---
date: 2022-07-19T09:56:46+08:00
author: "Rustle Karl"

title: "air 监控文件实时重载"
url:  "posts/go/tools/air"  # 永久链接
tags: [ "Go", "README" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

> https://github.com/cosmtrek/air

## 安装

```shell
go install github.com/cosmtrek/air@latest
```

## 初始化

```shell
air init
```

将创建一个 .air.toml 文件。

## 运行

默认当前目录下的 .air.toml 文件。

```shell
air
```

或者指定文件：

```shell
# firstly find `.air.toml` in current directory, if not found, use defaults
air -c .air.toml
```

或者更直接：

```shell
go run github.com/cosmtrek/air@latest -c .air.toml
```

## 配置文件示例

> from gitea

```ini
root = "."
tmp_dir = ".air"

[build]
cmd = "make backend"
bin = "gitea"
include_ext = ["go", "tmpl"]
exclude_dir = ["modules/git/tests", "services/gitdiff/testdata", "modules/avatar/testdata"]
include_dir = ["cmd", "models", "modules", "options", "routers", "services", "templates"]
exclude_regex = ["_test.go$"]
```
