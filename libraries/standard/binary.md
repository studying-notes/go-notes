---
date: 2020-10-27T09:25:05+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "二进制数据的序列化与反序列化"  # 文章标题
url:  "posts/go/libraries/standard/binary"  # 设置网页链接，默认使用文件名
tags: [ "go", "binary", "io" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

> TCP / IP 是大端字节序

大端字节序：高位字节在前，低位字节在后，这是人类读写数值的方法。
小端字节序：低位字节在前，高位字节在后，即以 0x1122 形式储存。

```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message struct {
	Head [4]byte
	Body [7]byte
	Tail [4]byte
}

func main() {
	msg := Message{
		Head: [4]byte{'h', 'e', 'a', 'd'},
		Body: [7]byte{',', 'h', 'e', 'l', 'l', 'o', ','},
		Tail: [4]byte{'t', 'a', 'i', 'l'},
	}
	writeBuf := new(bytes.Buffer)
	_ = binary.Write(writeBuf, binary.LittleEndian, msg)
	fmt.Println(writeBuf.String())

	buf := bytes.NewReader(writeBuf.Bytes())
	var message Message
	_ = binary.Read(buf, binary.LittleEndian, &message)
}
```

```go
package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	buf := []byte{0x9e, 0x8d}
	fmt.Println(buf)
	
	fmt.Printf("%d\n", binary.BigEndian.Uint16(buf))

	fmt.Printf("%d\n", int(buf[0])<<8|int(buf[1]))
	fmt.Printf("%d\n", int(buf[0])<<8|int(buf[1]))

	port := 40589
	fmt.Printf("%x\n", port)

	fmt.Println(byte(port >> 8))
	fmt.Println(byte(port))
}
```