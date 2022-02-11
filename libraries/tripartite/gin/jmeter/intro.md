---
date: 2020-11-21T20:58:30+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Jmeter 压力测试简介"  # 文章标题
url:  "posts/gin/project/jmeter/intro"  # 设置网页链接，默认使用文件名
tags: [ "gin", "jmeter", "测试" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 安装 JDK

**JDK 8**

https://code.aliyun.com/kar/ojdk8-8u271/raw/master/jdk-8u271-windows-x64.exe

**JDK 11**

> 阿里网盘

> 以下命令仅在 Cmd 管理员权限下有效。

### 添加 JAVA_HOME

```shell
wmic ENVIRONMENT create name="JAVA_HOME", username="<system>", VariableValue="C:\Program Files\Java\jdk-11"
```

### 添加 Path

```shell
wmic ENVIRONMENT where "name='Path' and username='<system>'" set VariableValue="%Path%;%JAVA_HOME%\bin;%JAVA_HOME%\jre\bin"
```

### 添加 CLASSPATH

```shell
wmic ENVIRONMENT create name="CLASSPATH", username="<system>", VariableValue=".;%JAVA_HOME%\lib;%JAVA_HOME%\lib\dt.jar;%JAVA_HOME%\lib\tools.jar"
```

### 验证

```shell
java -version
```

```
java version "11" 2018-09-25
Java(TM) SE Runtime Environment 18.9 (build 11+28)
Java HotSpot(TM) 64-Bit Server VM 18.9 (build 11+28, mixed mode)
```

## 安装 Jmeter

进入 jmeter 官网下载

https://mirror.bit.edu.cn/apache//jmeter/binaries/apache-jmeter-5.2.1.zip

## 解压执行

在 Windows 下运行 jmeter.bat，或者在 UNIX 下运行文件 jmeter，这两个文件都可以在 bin 目录下找到。在一个很短的等待之后，JMeter 的图形用户界面就会出现。

在 bin 目录中，还有其他几个测试人员可能会用到的脚本。

```shell
wmic ENVIRONMENT where "name='Path' and username='<system>'" set VariableValue="%Path%;C:\Developer\apache-jmeter-5.2.1\bin"
```

### Windows 脚本

- jmeter.bat ——运行 JMeter （默认 GUI 模式）。
- jmeter-n.cmd ——加载一个 JMX 文件，并在非 GUI 模式下运行。
- jmeter-n-r.cmd ——加载一个 JMX 文件，并在远程非 GUI 模式下运行。
- jmeter-t.cmd ——加载一个 JMX 文件，并在 GUI 模式下运行。
- jmeter-server.bat ——以服务器模式启动 JMeter。
- mirror-server.cmd ——在非 GUI 模式下启动 JMeter 镜像服务器。
- shutdown.cmd ——关闭一个非 GUI 实例（优雅的）。
- stoptest.cmd ——停止一个非 GUI 实例（中断式的）。

注意：关键字 LAST 可以与 jmeter-n.cmd、jmeter-t.cmd 和 jmeter-n-r.cmd 一起使用，这就意味着最近一个测试计划是交互式运行的。

### UNIX 脚本

- jmeter ——运行 JMeter （默认 GUI 模式）。定义了一些 JVM 设置，但并不是对所有 JVM 都生效。
- jmeter-server ——以服务器模式启动 JMeter （通过合适的参数来调用 JMeter 脚本）。
- jmeter.sh ——没有指定 JVM 选项的基础 JMeter 脚本。
- mirror-server.sh ——在非 GUI 模式下启动 JMeter 镜像服务器。
- shutdown.sh ——关闭一个非 GUI 实例（优雅的）。
- stoptest.sh ——停止一个非 GUI 实例（中断式的）。

## 修改运行属性

可以修改 bin 目录下的 jmeter.properties 文件，或者根据 jmeter.properties 文件创建用户自己的属性文件，并在命令行中指定属性文件名。

```shell

```

```shell

```
