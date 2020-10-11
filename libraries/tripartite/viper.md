---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Viper 中文教程"  # 文章标题
url:  "posts/go/libraries/tripartite/viper"  # 设置网页链接，默认使用文件名
tags: [ "go", "viper" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 简介

Viper 是适用于 Go 应用程序的完整配置解决方案。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。它支持以下特性：

- 设置默认值
- 从 JSON、TOML、YAML、HCL、envfile 和 Java properties 格式的配置文件读取配置信息
- 实时监控和重新读取配置文件
- 从环境变量中读取配置
- 从远程配置系统读取并监控配置变化
- 从命令行参数读取配置
- 从 buffer 读取配置
- 显式配置值

**安装**

```shell
go get github.com/spf13/viper
```

**优先级**

1. 显示调用Set设置值
2. 命令行参数
3. 环境变量
4. 配置文件
5. key/value 存储
6. 默认值

> 目前 Viper 配置的键（Key）是大小写不敏感的。

## 建立默认值

```go
viper.SetDefault("ContentDir", "content")
viper.SetDefault("LayoutDir", "layouts")
viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
```

## 读取配置文件

Viper 支持 JSON、TOML、YAML、HCL、envfile 和 Java properties 格式的配置文件。Viper 可以搜索多个路径，但目前单个 Viper 实例只支持单个配置文件。Viper 不默认任何配置搜索路径，将默认决策留给应用程序。

```go
viper.SetConfigFile("config.yaml") // 配置文件
viper.SetConfigName("config") // 配置文件名称(无扩展名)
viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项
viper.AddConfigPath("/etc/appname/")   // 查找配置文件所在的路径
viper.AddConfigPath("$HOME/.appname")  // 多次调用以添加多个搜索路径
viper.AddConfigPath(".")               // 还可以在工作目录中查找配置

viper.ReadInConfig() // 查找并读取配置文件
```

## 写入配置文件

- WriteConfig - 将当前的 viper 配置写入预定义的路径并覆盖。如果没有预定义的路径，则报错。
- SafeWriteConfig - 将当前的 viper 配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
- WriteConfigAs - 将当前的 viper 配置写入给定的文件路径。将覆盖给定的文件。
- SafeWriteConfigAs - 将当前的 viper 配置写入给定的文件路径。不会覆盖给定的文件。

```go
viper.WriteConfig() // 将当前配置写入 “viper.AddConfigPath()”和“viper.SetConfigName” 设置的预定义路径
viper.SafeWriteConfig()
viper.WriteConfigAs("/path/to/my/.config")
viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
viper.SafeWriteConfigAs("/path/to/my/.other_config")
```

## 监控并重新读取配置文件

Viper 支持在运行时实时读取配置文件的功能。

确保在调用 WatchConfig() 之前添加了所有的配置路径。

```go
viper.WatchConfig()
// 配置文件发生变更之后会调用的回调函数
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
})
```

如果某个对象只在初始化时读取一次设置，那就起不到热更新的作用。

## 从 io.Reader 读取配置

Viper 预先定义了许多配置源，如文件、环境变量、标志和远程 K/V 存储，还可以实现自己所需的配置源并将其提供给 viper。

```go
viper.SetConfigType("yaml") // 或者 viper.SetConfigType("YAML")

// 任何需要将此配置添加到程序中的方法
var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

viper.ReadConfig(bytes.NewBuffer(yamlExample))
viper.Get("name") // "steve"
```

## 覆盖设置

```go
viper.Set("Verbose", true)
viper.Set("LogFile", LogFile)
```

## 注册和使用别名

别名允许多个键引用单个值。

```go
viper.RegisterAlias("loud", "Verbose")  // 注册别名

viper.Set("verbose", true) // 结果与下一行相同
viper.Set("loud", true)   // 结果与前一行相同

viper.GetBool("loud") // true
viper.GetBool("verbose") // true
```

## 使用环境变量

用到再说

```
https://zhuanlan.zhihu.com/p/138691244
```
