---
date: 2020-11-29T22:05:23+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Jmeter 命令行运行"  # 文章标题
url:  "posts/gin/project/jmeter/cmd"  # 设置网页链接，默认使用文件名
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

## 原因

使用 GUI 方式启动 jmeter，运行线程较多的测试时，会造成内存和 CPU 的大量消耗，导致客户机卡死。所以正确的打开方式是在 GUI 模式下调整测试脚本，再用命令行模式执行。

## 命令

```shell
jmeter -n -t <testplan filename> -l <listener filename>
```

[![DgttzQ.png](https://s3.ax1x.com/2020/11/29/DgttzQ.png)](https://imgchr.com/i/DgttzQ)


## 示例

### 测试计划与结果，都在当前目录

```shell
jmeter -n -t test1.jmx -l result.jtl 
```

### 指定日志路径

```shell
jmeter -n -t test1.jmx -l report\01-result.csv -j report\01-log.log      
```

### 默认分布式执行

```shell
jmeter -n -t test1.jmx -r -l report\01-result.csv -j report\01-log.log
```

### 指定 IP 分布式执行

```shell
jmeter -n -t test1.jmx -R 192.168.10.25:1036 -l report\01-result.csv -j report\01-log.log
```

### 生成测试报表

```shell
jmeter -n -t 【Jmx脚本位置】-l 【中间文件result.jtl位置】-e -o 【报告指定文件夹】
```

```shell
jmeter -n -t test1.jmx  -l  report\01-result.jtl  -e -o tableresult
```

### 根据结果文件生成报告

```shell
jmeter -g【已经存在的 .jtl 文件的路径】-o 【用于存放 html 报告的目录】
```

```shell
jmeter -g result.jtl -o ResultReport 
```

## 报告释义

- APDEX(Application Performance Index)

应用程序性能满意度的标准，范围在 0-1 之间，1 表示达到所有用户均满意。

- Requests Summary

请求的通过率 (OK) 与失败率 (KO)，百分比显示。

-  Statistics

数据分析，将 Summary Report 和 Aggrerate Report 的结果合并。

- Errors

错误情况，依据不同的错误类型，将所有错误结果展示。

- Top 5 Errors by sampler

Top 5 错误

- Over Time

Response Times Over Time: 响应时间

Bytes Throughput Over Time: 字节接收/发送的数量

Latencies Over Time:延迟时间

-  Throughput

Hits Per Second: 每秒点击率

Codes Per Second: 每秒状态码数量

Transactions Per Second: 每秒事务量

Response Time Vs Request: 响应时间点请求的成功/失败数

Latency Vs Request: 延迟时间点请求的成功/失败数

- Response Times

Response Time Percentiles: 响应时间百分比

Active Threads Over Time: 激活线程数

Time Vs Threads: 测试过程中的线程数时续图

Response Time Distribution: 响应时间分布
