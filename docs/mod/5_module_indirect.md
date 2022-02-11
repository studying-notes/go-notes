---
date: 2020-11-15T20:29:40+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gomod 深入讲解 5"  # 文章标题
url:  "posts/go/docs/mod/5_module_indirect"  # 设置网页永久链接
tags: [ "go", "gomod" ]  # 标签
series: [ "Go 学习笔记" ]  # 系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

在使用 Go module 过程中，随着引入的依赖增多，也许你会发现 `go.mod` 文件中部分依赖包后面会出现一个 `// indirect` 的标识。这个标识总是出现在 `require` 指令中，其中 `// ` 与代码的行注释一样表示注释的开始，`indirect` 表示间接的依赖。

比如开源软件 Kubernetes（v1.17.0 版本）的 go.mod 文件中就有数十个依赖包被标记为 `indirect`：

```
require (
	github.com/Rican7/retry v0.1.0 // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20170425171159-5493cabe49f7 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/checkpoint-restore/go-criu v0.0.0-20190109184317-bdb7599cd87b // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	...
)
```

在执行命令 `go mod tidy` 时，Go module 会自动整理 `go.mod 文件 `，如果有必要会在部分依赖包的后面增加 `// indirect` 注释。一般而言，被添加注释的包肯定是间接依赖的包，而没有添加 `// indirect` 注释的包则是直接依赖的包，即明确的出现在某个 `import` 语句中。

然而，这里需要着重强调的是：并不是所有的间接依赖都会出现在 `go.mod` 文件中。

间接依赖出现在 `go.mod` 文件的情况，可能符合下面所列场景的一种或多种：

- 直接依赖未启用 Go module
- 直接依赖 go.mod 文件中缺失部分依赖

## 直接依赖未启用 Go module

如下图所示，Module A 依赖 B，但是 B 还未切换成 Module，也即没有 `go.mod` 文件，此时，当使用 `go mod tidy` 命令更新 A 的 `go.mod` 文件时，B 的两个依赖 B1 和 B2 将会被添加到 A 的 `go.mod` 文件中（前提是 A 之前没有依赖 B1 和 B2），并且 B1 和 B2 还会被添加 `// indirect` 的注释。

![](images/gomodule_indirect_01.png)

此时 Module A 的 `go.mod` 文件中 require 部分将会变成：

```
require (
	B vx.x.x
	B1 vx.x.x // indirect
	B2 vx.x.x // indirect
)
```

依赖 B 及 B 的依赖 B1 和 B2 都会出现在 `go.mod` 文件中。

## 直接依赖 go.mod 文件不完整

如上面所述，如果依赖 B 没有 `go.mod` 文件，则 Module A 将会把 B 的所有依赖记录到 A 的 `go.mod` 文件中。即便 B 拥有 `go.mod`，如果 `go.mod` 文件不完整的话，Module A 依然会记录部分 B 的依赖到 `go.mod` 文件中。

如下图所示，Module B 虽然提供了 `go.mod` 文件中，但 `go.mod` 文件中只添加了依赖 B1，那么此时 A 在引用 B 时，则会在 A 的 `go.mod` 文件中添加 B2 作为间接依赖，B1 则不会出现在 A 的 `go.mod` 文件中。

![](images/gomodule_indirect_02.png)

此时 Module A 的 `go.mod` 文件中 require 部分将会变成：

```
require (
	B vx.x.x
	B2 vx.x.x // indirect
)
```

由于 B1 已经包含进 B 的 `go.mod` 文件中，A 的 `go.mod` 文件则不必再记录，只会记录缺失的 B2。

## 总结

#### 为什么要记录间接依赖

在上面的例子中，如果某个依赖 B 没有 `go.mod` 文件，在 A 的 `go.mod` 文件中已经记录了依赖 B 及其版本号，为什么还要增加间接依赖呢？

我们知道 Go module 需要精确地记录软件的依赖情况，虽然此处记录了依赖 B 的版本号，但 B 的依赖情况没有记录下来，所以如果 B 的 `go.mod` 文件缺失了（或没有）这个信息，则需要在 A 的 `go.mod` 文件中记录下来。此时间接依赖的版本号将会根据 Go module 的版本选择机制确定一个最优版本。

#### 如何处理间接依赖

综上所述间接依赖出现在 `go.mod` 中，可以一定程度上说明依赖有瑕疵，要么是其不支持 Go module，要么是其 `go.mod` 文件不完整。

由于 Go 语言从 v1.11 版本才推出 module 的特性，众多开源软件迁移到 go module 还需要一段时间，在过渡期必然会出现间接依赖，但随着时间的推进，在 `go.mod` 中出现 `// indirect` 的机率会越来越低。

出现间接依赖可能意味着你在使用过时的软件，如果有精力的话还是推荐尽快消除间接依赖。可以通过使用依赖的新版本或者替换依赖的方式消除间接依赖。

#### 如何查找间接依赖来源

Go module 提供了 `go mod why` 命令来解释为什么会依赖某个软件包，若要查看 `go.mod` 中某个间接依赖是被哪个依赖引入的，可以使用命令 `go mod why -m <pkg>` 来查看。

比如，我们有如下的 `go.mod` 文件片断：

```
require (
	github.com/Rican7/retry v0.1.0 // indirect
	github.com/google/uuid v1.0.0
	github.com/renhongcai/indirect v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/text v0.3.2
)
```

我们希望确定间接依赖 `github.com/Rican7/retry v0.1.0 // indirect` 是被哪个依赖引入的，则可以使用命令 `go mod why` 来查看：

```
[root@ecs-d8b6 gomodule]# go mod why -m github.com/Rican7/retry
# github.com/Rican7/retry
github.com/renhongcai/gomodule
github.com/renhongcai/indirect
github.com/Rican7/retry
```

上面的打印信息中 `# github.com/Rican7/retry` 表示当前正在分析的依赖，后面几行则表示依赖链。`github.com/renhongcai/gomodule` 依赖 `github.com/renhongcai/indirect`，而 `github.com/renhongcai/indirect`依赖 `github.com/Rican7/retry`。由此我们就可以判断出间接依赖 `github.com/Rican7/retry` 是被 `github.com/renhongcai/indirect` 引入的。

另外，命令 `go mod why -m all` 则可以分析所有依赖的依赖链。
