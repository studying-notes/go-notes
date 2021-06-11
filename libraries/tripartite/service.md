---
date: 2020-11-14T20:03:29+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "用 Go 语言 编写 Windows 服务"  # 文章标题
url:  "posts/go/libraries/tripartite/service"  # 设置网页链接，默认使用文件名
tags: [ "go", "service"]  # 标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```go
package main

import (
	"fmt"
	"github.com/kardianos/service"
	"os"
	"time"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// 运行逻辑
	for {
		fmt.Println("running...")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "AAAAAAAAA", //服务显示名称
		DisplayName: "AAAAAAAAA", //服务名称
		Description: "AAAAAAAAA", //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err = s.Install(); err != nil {
				panic(err)
			}
			fmt.Println("服务安装成功!")
			if err = s.Start(); err != nil {
				panic(err)
			}
			fmt.Println("服务启动成功!")
		case "start":
			if err = s.Start(); err != nil {
				panic(err)
			}
			fmt.Println("服务启动成功!")
		case "stop":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			fmt.Println("服务关闭成功!")
		case "restart":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			fmt.Println("服务关闭成功!")
			if err = s.Start(); err != nil {
				panic(err)
			}
			fmt.Println("服务启动成功!")
		case "remove", "uninstall":
			if err = s.Stop(); err != nil {
				panic(err)
			}
			fmt.Println("服务关闭成功!")
			if err = s.Uninstall(); err != nil {
				panic(err)
			}
			fmt.Println("服务卸载成功!")
		}
		return
	}
	fmt.Println("run")

	err = s.Run()
	fmt.Println("run over")

	if err != nil {
		panic(err)
	}
	fmt.Println("over")
}
```

- 安装服务 `xx.exe install`
- 启动服务 `xx.exe start`
- 停止服务 `xx.exe stop`
- 重启服务 `xx.exe restart`
- 删除服务 `xx.exe remove/uninstall`
