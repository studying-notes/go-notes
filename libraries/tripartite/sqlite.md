---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "go-sqlite3 - SQLite / SQLCipher 操作示例"  # 文章标题
url:  "posts/go/libraries/tripartite/sqlite"  # 设置网页链接，默认使用文件名
tags: [ "go", "sqlcipher", "sqlite" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

> SQLCipher 相关语法完全相同，只是引入的驱动有差异。

- [引入驱动](#引入驱动)
- [连接数据库](#连接数据库)
- [创建数据表](#创建数据表)
	- [从 SQL 文件创建](#从-sql-文件创建)
- [插入](#插入)
- [查询](#查询)
	- [无条件查询](#无条件查询)
	- [条件查询](#条件查询)
- [删除](#删除)
- [事务插入](#事务插入)

## 引入驱动

```go
import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/fujiawei-dev/go-sqlcipher"
	"log"
	"os"
)
```

## 连接数据库

```go
func main() {
	_ = os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
```

自动创建一个空的不带任何表结构的数据库文件。

## 创建数据表

```go
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS userinfo
		(
			uid        INTEGER PRIMARY KEY AUTOINCREMENT,
			username   VARCHAR(64) NULL,
			department VARCHAR(64) NULL,
			created    DATE        NULL
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
```

### 从 SQL 文件创建

> SQLCipher 相关语法完全相同。

```go
	sqlStmt, _ := ioutil.ReadFile("./main.sql")
	_, err = db.Exec(string(sqlStmt))
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
```

## 插入

```go
	_, err = db.Exec("INSERT INTO userinfo(username, department, created) values('foo', 'bar', 'baz'), ('foo2', 'bar2', 'baz2')")
	if err != nil {
		log.Fatal(err)
	}
```

## 查询

### 无条件查询

```go
	rows, err := db.Query("select username, department from userinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var department string
		err = rows.Scan(&username, &department)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(username, department)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
```

### 条件查询

```go
	stmt, err = db.Prepare("select created from userinfo where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var created string
	err = stmt.QueryRow("3").Scan(&created)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created)
```

## 删除

```go
	_, err = db.Exec("delete from userinfo")
	if err != nil {
		log.Fatal(err)
	}
```

## 事务插入

```go
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec("Zhu", "Yue", fmt.Sprintf("VOL.%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
```
