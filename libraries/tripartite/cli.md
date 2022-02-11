---
date: 2021-01-01T20:12:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "urfave/cli - 构建 CLI 程序"  # 文章标题
url:  "posts/go/libraries/tripartite/cli"  # 设置网页链接，默认使用文件名
tags: [ "go", "cli" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## v1 & v2

### v1

```shell
go get github.com/urfave/cli
```

```go
import (
  "github.com/urfave/cli"
)
```

### v2

```shell
go get github.com/urfave/cli/v2
```

```go
import (
  "github.com/urfave/cli/v2" // imports as package "cli"
)
```

## 基本结构

通过 cli.NewApp() 创建一个实例，然后调用 Run() 方法就实现了一个最基本的命令行程序。为了让程序干点事情，可以指定一下入口函数 app.Action：
`

```go
package main

import (
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
```

## 公共配置

就是帮助里显示的提示信息：

```go
func main() {
	app := cli.NewApp()
	app.Name = "NewApp"
	app.Usage = "Browse Your Life in Pictures"
	app.Version = "1.0.0"
	app.Copyright = "(c) 2018-2020 Rustle Karl"
	app.EnableBashCompletion = true
}
```

## Flag 配置

显示以下帮助信息：

```ouput
   --lang FILE, -l FILE    read from FILE (default: "english")
   --port value, -p value  listening port (default: 8000)
```

实现代码：

```go
	var language string

	app.Flags = []cli.Flag {
		cli.IntFlag {
			Name: "port, p",
			Value: 8000,
			Usage: "listening port",
		},
		cli.StringFlag {
			Name: "lang, l",
			Value: "english",
			Usage: "read from `FILE`",
			Destination: &language,
		},
	}
```

- Name 字段中逗号后面的字符表示 flag 的简写，也就是说 "--port" 和 "-p" 是等价的。
- Value 字段可以指定 flag 的默认值。
- Usage 字段是 flag 的描述信息。
- Destination 字段可以为该 flag 指定一个接收者，比如上面的 language 变量。解析完 "--lang" 这个 flag 后会自动存储到这个变量里，后面的代码就可以直接使用这个变量的值了。
- 给用户增加一些属性值类型的提示，可以通过占位符（ placeholder ）来实现，比如上面的 "--lang FILE"。占位符通过 `` 符号来标识。

正常来说帮助信息里的 flag 是按照代码里的声明顺序排列的，如果想让它们按照字典序排列的话，可以借助于 sort：

```go
import "sort"
sort.Sort(cli.FlagsByName(app.Flags))
```

## Command 配置

命令行程序除了有 flag，还有 command（比如 git log, git commit 等等）。另外每个 command 可能还有 subcommand，也就必须要通过添加两个命令行参数才能完成相应的操作。

```shell
NAME:
   GoTest db - database operations

USAGE:
   GoTest db command [command options] [arguments...]

COMMANDS:
     insert  insert data
     delete  delete data

OPTIONS:
   --help, -h  Help!Help!
```

每个 command 都对应于一个 cli.Command 接口的实例，入口函数通过 Action 指定。如果想在帮助信息里实现分组显示，可以为每个 command 指定一个 Category。具体代码如下：

```go
	app.Commands = []cli.Command {
		{
			Name: "add",
			Aliases: []string{"a"},
			Usage: "calc 1+1",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("1 + 1 = ", 1 + 1)
				return nil
			},
		},
		{
			Name: "sub",
			Aliases: []string{"s"},
			Usage: "calc 5-3",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("5 - 3 = ", 5 - 3)
				return nil
			},
		},
		{
			Name: "db",
			Usage: "database operations",
			Category: "database",
			Subcommands: []cli.Command {
				{
					Name: "insert",
					Usage: "insert data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert subcommand")
						return nil
					},
				},
				{
					Name: "delete",
					Usage: "delete data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete subcommand")
						return nil
					},
				},
			},
		},
	}
```

如果你想在 command 执行前后执行后完成一些操作，可以指定 app.Before/app.After 这两个字段：

```go
	app.Before = func(c *cli.Context) error {
		fmt.Println("app Before")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Println("app After")
		return nil
	}
```

## 示例 Demo

```go
package cli

import (
	"fmt"
	"os"
	"log"
	"sort"
	"gopkg.in/urfave/cli.v1"
)

func Run() {
	var language string

	app := cli.NewApp()
	app.Name = "GoTest"
	app.Usage = "hello world"
	app.Version = "1.2.3"
	app.Flags = []cli.Flag {
		cli.IntFlag {
			Name: "port, p",
			Value: 8000,
			Usage: "listening port",
		},
		cli.StringFlag {
			Name: "lang, l",
			Value: "english",
			Usage: "read from `FILE`",
			Destination: &language,
		},
	}
	app.Commands = []cli.Command {
		{
			Name: "add",
			Aliases: []string{"a"},
			Usage: "calc 1+1",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("1 + 1 = ", 1 + 1)
				return nil
			},
		},
		{
			Name: "sub",
			Aliases: []string{"s"},
			Usage: "calc 5-3",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("5 - 3 = ", 5 - 3)
				return nil
			},
		},
		{
			Name: "db",
			Usage: "database operations",
			Category: "database",
			Subcommands: []cli.Command {
				{
					Name: "insert",
					Usage: "insert data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert subcommand")
						return nil
					},
				},
				{
					Name: "delete",
					Usage: "delete data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete subcommand")
						return nil
					},
				},
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		fmt.Println(c.String("lang"), c.Int("port"))
		fmt.Println(language)

		// if c.Int("port") == 8000 {
		// 	return cli.NewExitError("invalid port", 88)
		// }

		return nil
	}
	app.Before = func(c *cli.Context) error {
		fmt.Println("app Before")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Println("app After")
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	cli.HelpFlag = cli.BoolFlag {
		Name: "help, h",
		Usage: "Help!Help!",
	}

	cli.VersionFlag = cli.BoolFlag {
		Name: "print-version, v",
		Usage: "print version",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
```

```go

```



```go

```

```go

```



```go

```

```go

```
