---
date: 2020-10-10T14:33:53+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "sqlx - 扩展标准库 database/sql"  # 文章标题
description: "扩展标准库 database/sql"
url:  "posts/go/libraries/tripartite/sqlx"  # 设置网页链接，默认使用文件名
tags: [ "go", "sqlx", "sql", "mysql"]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

```
go get github.com/jmoiron/sqlx
```

- [连接数据库](#连接数据库)
- [查询](#查询)
	- [查询单行](#查询单行)
	- [查询多行](#查询多行)
- [插入、更新和删除](#插入更新和删除)
- [NamedExec](#namedexec)
- [NamedQuery](#namedquery)
- [事务操作](#事务操作)
- [sqlx.In 批量插入](#sqlxin-批量插入)
	- [bindvars（绑定变量）](#bindvars绑定变量)
	- [拼接语句实现](#拼接语句实现)
	- [sqlx.In 实现](#sqlxin-实现)
- [sqlx.In 查询](#sqlxin-查询)
	- [in](#in)
	- [FIND_IN_SET](#find_in_set)

## 连接数据库

```go
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db, _ = sqlx.Connect("mysql", dsn)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10) // 空闲连接
	return
}
```

也可以用 MustConnect：

```go
// MustConnect connects to a database and panics on error.
func MustConnect(driverName, dataSourceName string) *DB {
	db, err := Connect(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}
```

```go
// Connect to a database and verify with a ping.
func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
```

## 查询

### 查询单行

```go
type User struct {
	Id   int
	Name string
	Age  int
}

func main() {
	sqlStr := "select * from user where id=?"
	var u User
	_ := db.Get(&u, sqlStr, 1)
	fmt.Printf("%+v\n", u)
}
```

这里遇到一个巨大的坑，结构体的成员必须以大写字母开头，这点与标准库不一样。

### 查询多行

```go
func main() {
	sqlStr := "select * from user where id>?"
	var u []User
	err := db.Select(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("%+v\n", u)
}
```

## 插入、更新和删除

与标准库操作一致，不再重复。

## NamedExec

`DB.NamedExec` 方法用来绑定 SQL 语句与结构体或 Map 中的同名字段。

```go
func main() {
	sqlStr := "INSERT INTO user (name,age) VALUES (:name,:age)"
	_, _ = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "sqlite",
			"age":  28,
		})
	return
}
```

## NamedQuery

同上，但这里是查询。

```go
func main() {
	sqlStr := "SELECT * FROM user WHERE name=:name"

	// Map 映射
	rows, _ := db.NamedQuery(sqlStr, map[string]interface{}{"name": "go"})
	defer rows.Close()
	for rows.Next() {
		var u User
		_ := rows.StructScan(&u)
		fmt.Printf("user:%#v\n", u)
	}

	// 根据结构体字段进行映射
	u := User{Name: "redis"}
	rows, _ = db.NamedQuery(sqlStr, u)
	defer rows.Close()
	for rows.Next() {
		var u User
		_ := rows.StructScan(&u)
		fmt.Printf("user:%#v\n", u)
	}
}
```

## 事务操作

对于事务操作，我们可以使用 `sqlx` 中提供的 `db.Beginx()` 和 `tx.Exec()` 方法。

```go
func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)

	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	sqlStr1 := "update user set age=20 where id=?"

	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return
	}
	if n != 1 {
		return
	}

	sqlStr2 := "update user set age=50 where id=?"
	rs, err = tx.Exec(sqlStr2, 1234)
	if err != nil {
		return
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return
	}
	if n != 1 {
		return
	}
}
```

## sqlx.In 批量插入

`sqlx.In` 是 `sqlx` 提供的一个非常方便的函数。

```
https://github.com/jmoiron/sqlx/issues/123
```

### bindvars（绑定变量）

查询占位符 `?` 在内部称为 **bindvars（查询占位符）**，它非常重要，应该始终使用它向数据库发送值，因为它可以防止 SQL 注入攻击。`database/sql` 不尝试对查询文本进行任何验证；它与编码的参数一起按原样发送到服务器，除非驱动程序实现一个特殊的接口，否则在执行之前，查询是在服务器上准备的。因此 `bindvars` 是特定于数据库的:

- MySQL 中使用 `?`
- PostgreSQL 使用枚举的 `$1`、`$2` 等 bindvar 语法
- SQLite 中 `?` 和 `$1` 的语法都支持
- Oracle 中使用 `:name` 的语法

`bindvars` 的一个常见误解是，它用来在 SQL 语句中插入值，其实仅用于参数化，不允许更改 SQL 语句的结构。例如，使用 `bindvars` 尝试参数化列或表名将不起作用：

```go
// 不能用来插入表名（做SQL语句中表名的占位符）
db.Query("SELECT * FROM ?", "mytable")

// 也不能用来插入列名（做SQL语句中列名的占位符）
db.Query("SELECT ?, ? FROM people", "name", "location")
```

### 拼接语句实现

```go
type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)

	var users []User
	users = append(users, User{Name: "wx", Age: 1})
	users = append(users, User{Name: "qw", Age: 31})
	users = append(users, User{Name: "zq", Age: 16})

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*2)
	for _, u := range users {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}

	stmt := fmt.Sprintf("INSERT INTO user (name, age) VALUES %s",
		strings.Join(valueStrings, ","))
	_, _ = db.Exec(stmt, valueArgs...)
}
```

### sqlx.In 实现

前提是需要结构体实现 `driver.Valuer` 接口：

```go
func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}
```

```go
type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)

	var users []interface{}
	users = append(users, User{Name: "wre", Age: 1})
	users = append(users, User{Name: "qwrew", Age: 31})
	users = append(users, User{Name: "zsfdq", Age: 16})

	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users...,
	)
	fmt.Println(query)
	fmt.Println(args)
	_, _ = db.Exec(query, args...)
}
```

## sqlx.In 查询

在 `sqlx` 查询语句中实现 `in` 查询和 `FIND_IN_SET` 函数，即实现：

```sql
SELECT * FROM user WHERE id in (3, 2, 1);
SELECT * FROM user WHERE id in (3, 2, 1) ORDER BY FIND_IN_SET(id, '3,2,1');
```

### in

```go
type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)

	var users []User
	ids := []int{1, 1234, 1235}
	query, args, _ := sqlx.In("SELECT name, age FROM user WHERE id IN (?)", ids)
	query = db.Rebind(query)
	_ = db.Select(&users, query, args...)

	fmt.Printf("%+v\n", users)
}
```

### FIND_IN_SET

```go
type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

func main() {
	var db *sqlx.DB
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)

	var users []User
	ids := []int{1, 1234, 1235}

	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, _ := sqlx.In("SELECT name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))

	query = db.Rebind(query)
	_ = db.Select(&users, query, args...)

	fmt.Printf("%+v\n", users)
}
```
