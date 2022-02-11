---
date: 2020-07-20T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Cobra - 构建 CLI 程序"  # 文章标题
url:  "posts/go/libraries/tripartite/cobra"  # 设置网页链接，默认使用文件名
tags: [ "go", "cobra", "cli" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```
go get -u github.com/spf13/cobra
```

## 基本概念

```
httpd start --port=1313
```

其中，`start` 是命令 `command`，`port` 是标志 `flag`，`1312` 是参数 `args`。


## 初始化新项目

```
cobra init --pkg-name example

cobra add server
cobra add config
cobra add create -p 'configCmd'
```


```go
var cmd = &cobra.Command{
	Use:   "子命令的命令标识",
	Short: "简短说明",
	Long:  "完整说明",
}
```

- Use：子命令的命令标识。
- Short：简短说明，在 help 命令输出的帮助信息中展示。
- Long：完整说明，在 help 命令输出的帮助信息中展示。


```
var mode int8
cmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词转换的模式")
```

第一个参数为需要绑定的变量，第二个参数为接收该参数的完整的命令标志，第三个参数为对应的短标识，第四个参数为默认值，第五个参数为使用说明。
