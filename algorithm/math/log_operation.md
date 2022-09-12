---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "对数运算"  # 文章标题
url:  "posts/go/algorithm/math/log_operation"  # 设置网页永久链接
tags: [ "Go", "log-operation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

```go
func main() {
    // 以自然对数 e 为底
    fmt.Println(math.Log(10))
    // 以 2 为底
    fmt.Println(math.Log2(16))  // 4
    // 以 10 为底
	fmt.Println(math.Log10(100))  // 2
    // 对数运算法则
	fmt.Println(math.Log(16)/math.Log(2))  // 2
}
```

```go

```
