---
date: 2020-11-15T12:31:46+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "Iota 自增表达式"  # 文章标题
url:  "posts/go/docs/grammar/iota"  # 设置网页永久链接
tags: [ "Go", "iota" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

iota 常用于 const 表达式中，其值是从零开始，const 声明块中每增加一行 iota 值自增 1。

## 规则

iota 代表了 const 声明块的行索引（下标从 0 开始）。

const 声明还有个特点，即**第一个常量必须指定一个表达式，后续的常量如果没有表达式，则继承上面的表达式**。

```go
const (
    bit0, mask0 = 1 << iota, 1<<iota - 1   //const声明第0行，即iota==0
    bit1, mask1                            //const声明第1行，即iota==1, 表达式继承上面的语句
    _, _                                   //const声明第2行，即iota==2
    bit3, mask3                            //const声明第3行，即iota==3
)
```

- 第 0 行的表达式展开即 `bit0, mask0 = 1 << 0, 1<<0 - 1`，所以 bit0 == 1，mask0 == 0 ；
- 第 1 行没有指定表达式继承第一行，即 `bit1, mask1 = 1 << 1, 1<<1 - 1`，所以 bit1 == 2，mask1 == 1 ；
- 第 2 行没有定义常量
- 第 3 行没有指定表达式继承第一行，即 `bit3, mask3 = 1 << 3, 1<<3 - 1`，所以 bit0 == 8，mask0 == 7 ；

## 编译原理

const 块中每一行在 go 中使用 spec 数据结构描述，spec 声明如下：

```go
    // A ValueSpec node represents a constant or variable declaration
    // (ConstSpec or VarSpec production).
    //
    ValueSpec struct {
        Doc     *CommentGroup // associated documentation; or nil
        Names   []*Ident      // value names (len(Names) > 0)
        Type    Expr          // value type; or nil
        Values  []Expr        // initial values; or nil
        Comment *CommentGroup // line comments; or nil
    }
```

这里我们只关注 ValueSpec.Names，这个切片中保存了一行中定义的常量，如果一行定义 N 个常量，那么 ValueSpec.Names 切片长度即为 N。

const 块实际上是 spec 类型的切片，用于表示 const 中的多行。

所以编译期间构造常量时的伪算法如下：

```go
    for iota, spec := range ValueSpecs {
        for i, name := range spec.Names {
            obj := NewConst(name, iota...) //此处将iota传入，用于构造常量
			...
        }
    }
```

从上面可以更清晰的看出 **iota 实际上是遍历 const 块的索引**，每行中即便多次使用 iota，其值也不会递增。

```go

```
