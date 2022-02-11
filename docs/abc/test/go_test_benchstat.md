---
date: 2020-11-16T21:01:29+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "深入测试标准库之性能测试分析工具"  # 文章标题
url:  "posts/go/abc/memory/test/go_test_benchstat"  # 设置网页永久链接
tags: [ "go", "test" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

`benchmark`测试是实际项目中经常使用的性能测试方法，我们可以针对某个函数或者某个功能点增加`benchmark`测试，
以便在CI测试中监测其性能变化，当该函数或功能性能下降时能够及时发现。

此外，在日常开发活动中或者参与开源贡献时也有可能针对某个函数或功能点做一些性能优化，此时，如何把`benchmark`测试数据呈现出来便非常重要了，因为你很可能在优化前后执行多次`benchmark`测试，手工分析这些测试结果无疑是低效的。

本节结合笔者在`Golang`社区参与开源贡献时的经历，介绍一下由官方推荐的性能测试分析工具`benchstat`，权当抛砖引玉之用。

## 认识数据
我们先看一个`benchmark`测试样本：
```
BenchmarkReadGoSum-4   	    2223	    521556 ns/op
```
该样本包含一个测试名字`BenchmarkReadGoSum-4`（其中`-4`表示测试环境为4个cpu）、测试迭代次数（2223）和每次迭代的花费的时间（521556ns）。

尽管每个样本中的时间已经是多次迭代后的平均值，但为了更好的分析性能，往往需要多个样本。

使用`go test`的`-count=N`参数可以指定执行`benchmark`N次，从而产生N个样本，比如产生15个样本：
```
BenchmarkReadGoSum-4   	    2223	    521556 ns/op
BenchmarkReadGoSum-4   	    2347	    516675 ns/op
BenchmarkReadGoSum-4   	    2340	    538406 ns/op
BenchmarkReadGoSum-4   	    2130	    548440 ns/op
BenchmarkReadGoSum-4   	    2391	    514602 ns/op
BenchmarkReadGoSum-4   	    2394	    527955 ns/op
BenchmarkReadGoSum-4   	    2313	    536693 ns/op
BenchmarkReadGoSum-4   	    2330	    538244 ns/op
BenchmarkReadGoSum-4   	    2360	    516426 ns/op
BenchmarkReadGoSum-4   	    2407	    541435 ns/op
BenchmarkReadGoSum-4   	    2154	    544386 ns/op
BenchmarkReadGoSum-4   	    2362	    540411 ns/op
BenchmarkReadGoSum-4   	    2305	    581713 ns/op
BenchmarkReadGoSum-4   	    2204	    519633 ns/op
BenchmarkReadGoSum-4   	    1867	    602543 ns/op
```
手工分析多个样本将会是一项非常有挑战的工作，因为你可能需要根据统计学规则抛弃一些异常的样本，剩下的样本再取平均值。

## benchstat
`benchstat`为Golang官方推荐的一款命令行工具，可以针对一组或多组样本进行分析，如果同时分析两组样本（比如优化前和优化后），还可以给出性能变化结果。

使用命令`go get golang.org/x/perf/cmd/benchstat`即可快捷安装，它将被安装到`$GOPATH/bin`目录中。通常我们会将该目录添加到`PATH`环境变量中。

使用时我们需要把`benchmark`测试样子输出到文件中，`benchstat`会读取这些文件，命令格式如下：
```
benchstat [options] old.txt [new.txt] [more.txt ...]
```

#### 分析一组样本
我们把上面的15个样本输出到名为`BenchmarkReadGoSum.before`的文件，然后使用`benchstat`分析：
```
# benchstat BenchmarkReadGoSum.before
name         time/op
ReadGoSum-4  531µs ± 3%
```
输出结果包括一个耗时平均值（531µs）和样本离散值（3%）。

#### 分析两组样本
同上，我们把性能优化后的结果输出到名为`BenchmarkReadGoSum.after`的文件，然后使用`benchstat`分析优化的效果：
```
# benchstat BenchmarkReadGoSum.before BenchmarkReadGoSum.after
name         old time/op  new time/op  delta
ReadGoSum-4   531µs ± 3%   518µs ± 7%  -2.41%  (p=0.033 n=13+15)
```
当只有两组样本时，`benchstat`还会额外计算出差值，比如本例中，平均花费时间下降了`2.41%`。

另外，`p=0.033`表示结果的可信程度，p 值越大可信程度越低，统计学中通常把`p=0.05`作为临界值，超过此值说明结果不可信，可能是样本过少等原因。

`n=13+15`表示采用的样本数量，出于某些原因(比如数据值反常，过大或过小)，`benchstat`会舍弃某些样本，本例中优化前的数据中舍弃了两个样本，优化后的数据没有舍弃，所以`13+15`，表示两组样本分别采用了13和15个样本。

## 小结
在`Golang`贡献者指导文档中，特别提到如果提交的代码涉及性能变化，需要将`benchstat`结果上传，以便代码审核者查看。

当然，我们也可以在闭源项目中使用，比起手工分析样本，`benchstat`明显可以大大提升效率。
