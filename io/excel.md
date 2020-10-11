---
date: 2020-07-20T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 读写 Excel 表格"  # 文章标题
url:  "posts/go/io/excel"  # 设置网页链接，默认使用文件名
tags: [ "go", "excel"]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```
go get github.com/tealeg/xlsx/v3
```

## 打开文件

```go
package main

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
)

func main() {
	wb, err := xlsx.OpenFile("io/xlsx/dev.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	for idx, sheet := range wb.Sheets {
		fmt.Println(idx, sheet.Name)
	}
}
```

## 读取单元格数据

```go
package main

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
)

func main() {
	wb, err := xlsx.OpenFile("io/xlsx/dev.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	// 遍历表
	for _, sheet := range wb.Sheets {
		fmt.Println(sheet.Name)
		// 遍历行读取
		for i := 0; i < sheet.MaxRow; i++ {
			// 遍历每行的列读取
			for j := 0; j < sheet.MaxCol; j++ {
				cell, _ := sheet.Cell(i, j)
				fmt.Print(cell.String(), "\t")
			}
			fmt.Println()
		}
	}
}
```

## 创建文件

```go
package main

import (
	"github.com/tealeg/xlsx/v3"
	"log"
)

func main() {
	var (
		file            *xlsx.File
		sheet           *xlsx.Sheet
		row, row1, row2 *xlsx.Row
		cell            *xlsx.Cell
		err             error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	row = sheet.AddRow()
	row.SetHeightCM(1)
	cell = row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "年龄"

	row1 = sheet.AddRow()
	row1.SetHeightCM(1)
	cell = row1.AddCell()
	cell.Value = "数学"
	cell = row1.AddCell()
	cell.Value = "18"

	row2 = sheet.AddRow()
	row2.SetHeightCM(1)
	cell = row2.AddCell()
	cell.Value = "语文"
	cell = row2.AddCell()
	cell.Value = "28"

	err = file.Save("io/xlsx/dev.xlsx")
	if err != nil {
		log.Fatal(err)
	}
}
```

## 修改文件

```go
package main

import (
	"github.com/tealeg/xlsx/v3"
	"log"
)

func main() {
	path := "io/xlsx/dev.xlsx"
	file, err := xlsx.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}
	first := file.Sheets[0]
	row := first.AddRow()
	row.SetHeightCM(1)
	
	// 在原有之上添加
	cell := row.AddCell()
	cell.Value = "松岛枫"
	cell = row.AddCell()
	cell.Value = "99"

	err = file.Save(path)
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

