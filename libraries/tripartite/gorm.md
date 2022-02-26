---
date: 2020-09-19T13:39:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "gorm - 数据库操作"  # 文章标题
description: "用 GORM 实现创建、查询、更新和删除操作"
url:  "posts/go/libraries/tripartite/gorm"  # 设置网页链接，默认使用文件名
tags: [ "go", "gorm", "orm", "mysql", "sql", "sqlite" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
toc: true  # 是否自动生成目录
---

## 目录

- [目录](#目录)
- [Distinct](#distinct)
- [Update/Updates](#updateupdates)
- [First/Find/Scan 区别](#firstfindscan-区别)
- [创建一张表](#创建一张表)
- [功能预览](#功能预览)
- [连接到数据库](#连接到数据库)
	- [引入驱动](#引入驱动)
	- [支持的数据库示例](#支持的数据库示例)
	- [MySQL](#mysql)
		- [自定义驱动](#自定义驱动)
		- [已存在的连接](#已存在的连接)
		- [MySQL 操作示例](#mysql-操作示例)
	- [PostgreSQL](#postgresql)
		- [自定义驱动](#自定义驱动-1)
		- [已存在的连接](#已存在的连接-1)
	- [SQLite3](#sqlite3)
	- [SQL Server](#sql-server)
	- [Clickhouse](#clickhouse)
	- [连接池](#连接池)
- [声明数据表模型](#声明数据表模型)
	- [约定](#约定)
	- [字段级权限控制](#字段级权限控制)
	- [支持的结构体标签](#支持的结构体标签)
	- [定义字段字符集](#定义字段字符集)
- [gorm 中的默认设置](#gorm-中的默认设置)
	- [开启日志输出打印 SQL 语句](#开启日志输出打印-sql-语句)
	- [gorm.Model](#gormmodel)
	- [默认 ID 作为 Primary Key](#默认-id-作为-primary-key)
	- [自定义表名而非结构体名](#自定义表名而非结构体名)
	- [根据条件切换表名](#根据条件切换表名)
	- [启用单数表名](#启用单数表名)
	- [在执行语句时指定表名](#在执行语句时指定表名)
	- [更改 grom 默认表名设置方式](#更改-grom-默认表名设置方式)
	- [默认列名为字段的下划线式](#默认列名为字段的下划线式)
- [常用功能](#常用功能)
	- [自动迁移数据表模型](#自动迁移数据表模型)
	- [设置数据表字符集、引擎等](#设置数据表字符集引擎等)
	- [检查表是否存在](#检查表是否存在)
	- [增删改数据表的结构](#增删改数据表的结构)
	- [索引和约束](#索引和约束)
- [创建](#创建)
	- [常用方式](#常用方式)
	- [指定字段](#指定字段)
	- [忽略字段](#忽略字段)
	- [批量插入](#批量插入)
	- [创建钩子](#创建钩子)
- [查询](#查询)
	- [基本查询](#基本查询)
		- [根据主键获取第一条记录](#根据主键获取第一条记录)
		- [根据主键获取指定的某条记录，仅当主键为整型时可用](#根据主键获取指定的某条记录仅当主键为整型时可用)
		- [根据主键获取最后一条记录](#根据主键获取最后一条记录)
		- [随机获取一条记录](#随机获取一条记录)
	- [查询结果分析](#查询结果分析)
		- [获取所有的记录](#获取所有的记录)
	- [条件查询](#条件查询)
		- [通过 结构体 / Map 查询](#通过-结构体--map-查询)
		- [通过 Where 条件语句查询](#通过-where-条件语句查询)
		- [通过 In 条件语句查询](#通过-in-条件语句查询)
		- [通过 LIKE 条件语句查询](#通过-like-条件语句查询)
		- [通过 Not 条件语句查询](#通过-not-条件语句查询)
		- [通过 Or 条件语句查询](#通过-or-条件语句查询)
		- [FirstOrCreate 不存在就插入记录](#firstorcreate-不存在就插入记录)
	- [子查询](#子查询)
	- [Select 部分字段查询](#select-部分字段查询)
		- [映射结构体](#映射结构体)
		- [一行一行赋值](#一行一行赋值)
	- [排序 Order](#排序-order)
		- [多字段排序](#多字段排序)
		- [覆盖排序](#覆盖排序)
	- [限制输出数量 LIMIT](#限制输出数量-limit)
	- [统计数量 COUNT](#统计数量-count)
	- [分组 Group & Having](#分组-group--having)
		- [单字段与多字段](#单字段与多字段)
		- [一行一行获取](#一行一行获取)
		- [一次性获取](#一次性获取)
	- [JOIN 连接查询](#join-连接查询)
	- [Pluck 查询：获取一个列作为切片](#pluck-查询获取一个列作为切片)
		- [Pluck + Where 查询](#pluck--where-查询)
	- [Scan 扫描：获取多个列的值](#scan-扫描获取多个列的值)
	- [原生 SQL Scan](#原生-sql-scan)
		- [Exec](#exec)
- [更新](#更新)
	- [Save 更新所有字段](#save-更新所有字段)
	- [Update 更新指定字段](#update-更新指定字段)
		- [根据主键更新单个属性](#根据主键更新单个属性)
		- [根据条件更新单个属性](#根据条件更新单个属性)
		- [用 map 更新多个属性](#用-map-更新多个属性)
		- [用 struct 更新多个属性](#用-struct-更新多个属性)
	- [Select 更新部分字段](#select-更新部分字段)
	- [Omit 忽略更新部分字段](#omit-忽略更新部分字段)
	- [只更新指定字段，不更新自动更新字段](#只更新指定字段不更新自动更新字段)
		- [更新单个属性](#更新单个属性)
		- [更新多个属性](#更新多个属性)
	- [批量更新](#批量更新)
	- [获取更新记录总数](#获取更新记录总数)
	- [使用 SQL 计算表达式](#使用-sql-计算表达式)
- [删除](#删除)
	- [删除记录](#删除记录)
	- [批量删除](#批量删除)
	- [软删除](#软删除)
	- [物理删除](#物理删除)
- [事务处理](#事务处理)
	- [禁用默认事务](#禁用默认事务)
	- [一般流程](#一般流程)
	- [嵌套事务](#嵌套事务)
	- [手动事务](#手动事务)
	- [手动事务示例](#手动事务示例)
	- [SavePoint、RollbackTo](#savepointrollbackto)
- [实体关联](#实体关联)
	- [自动创建、更新](#自动创建更新)
	- [Belongs To](#belongs-to)
		- [重写外键](#重写外键)
		- [重写引用](#重写引用)
		- [外键约束](#外键约束)
	- [Has One](#has-one)
		- [重写外键](#重写外键-1)
		- [重写引用](#重写引用-1)
		- [外键约束](#外键约束-1)
		- [预加载](#预加载)
		- [Joins 预加载](#joins-预加载)
	- [Has Many](#has-many)
		- [重写外键](#重写外键-2)
		- [重写引用](#重写引用-2)
		- [外键约束](#外键约束-2)
- [GORM 时区配置](#gorm-时区配置)
	- [系统默认时区](#系统默认时区)
	- [设置时区](#设置时区)
- [官网资料](#官网资料)

用 GORM 实现创建、查询、更新和删除操作。v2 版本与 v1 版本在增删查改方面基本没有区别，只在初始化时略有区别。

```shell
# v1
go get -u github.com/jinzhu/gorm

# v2
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

MySQL 的 8.0 以上版本不支持零日期格式，导致 gorm 插入默认数据出错。

## Distinct

只能多行全字段去重，不能多行单字段去重

## Update/Updates

更新时的数据类型必须与模型结构体保持一致，比如字符型数值不能用于整数，否则乱码。Where 则不存在这个问题。

## First/Find/Scan 区别

First / Find 的结构体的 TableName 必须是存在的表，否则报错，即使指定了 db.Table() or db.Modlel() ；

Scan 可以是任意结构体，但必须指定 db.Table() or db.Modlel()。

## 创建一张表

```sql
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL DEFAULT '',
  `age` int(3) NOT NULL DEFAULT '0',
  `sex` tinyint(3) NOT NULL DEFAULT '0',
  `phone` varchar(40) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4
```

当然最简单的还是用 GROM 的自动迁移功能。

## 功能预览

```go
package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

// gorm.Model 官方定义的通用模型
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Product struct {
	gorm.Model // 官方定义的通用模型
	Code       string
	Price      uint
}

func main() {
	db, err := gorm.Open("sqlite3", "gorm_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// AutoMigrate run auto migration for given models,
	// will only add missing fields, won't delete/change current data
	// 自动迁移，根据给定的模型自动创建或修改同名数据表
	db.AutoMigrate(&Product{})

	// 插入一条数据
	db.Create(&Product{Code: "SN2020001", Price: 122})
	db.Create(&Product{Code: "SN2020002", Price: 123})

	// 查询一条数据
	var product Product

	// find product with id 1
	db.First(&product, 1)
	//fmt.Printf("%+v\n", product)

	// find product with code SN2020001
	db.First(&product, "code = ?", "SN2020001")
	//fmt.Printf("%+v\n", product)

	// 查询多条数据，标记已删除的默认无法获取
	var products []Product
	db.Find(&products, "code = ?", "SN2020")
	//fmt.Printf("%+v\n", products)

	// Update - update product's price to 2000
	// 结构体的字段名，而非表的字段名
	db.Model(&product).Update("Price", 2000)
	// product price 字段同步更新
	fmt.Printf("%+v\n", product)

	// Delete - delete product
	// 只是标记了 DeletedAt 字段
	db.Delete(&product)
	// product DeletedAt 字段不会更新
	fmt.Printf("%+v\n", product)
}
```

## 连接到数据库

### 引入驱动

**官方**

```go
import _ "github.com/jinzhu/gorm/dialects/mysql"
// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"
```

**三方**

```go
import _ "github.com/go-sql-driver/mysql"
```

### 支持的数据库示例

### MySQL

```go
gorm.Open("mysql", "user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local")
```

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
  DefaultStringSize: 256, // default size for string fields
  DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
  DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
  DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
  SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
}), &gorm.Config{})
```

- `parseTime` - 需要解析 `time.Time` 时设置为 `True`
- `charset` - 为了完整支持 `UTF-8` 编码，需要设置为 `charset=utf8mb4`

#### 自定义驱动

```go
import (
  _ "example.com/my_mysql_driver"
  "gorm.io/gorm"
)

db, err := gorm.Open(mysql.New(mysql.Config{
  DriverName: "my_mysql_driver",
  DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
}), &gorm.Config{})
```

#### 已存在的连接

```go
import (
  "database/sql"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

sqlDB, err := sql.Open("mysql", "mydb_dsn")
gormDB, err := gorm.Open(mysql.New(mysql.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```

#### MySQL 操作示例

```go
package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

// 定义一个数据模型
// 列名是字段名的蛇形小写，即下划线形式
type User struct {
	Id       uint   `gorm:"AUTO_INCREMENT"`
	Name     string `gorm:"size:50"`
	Age      int    `gorm:"size:3"`
	Birthday *time.Time
	Email    string `gorm:"type:varchar(50);unique_index"`
	PassWord string `gorm:"type:varchar(25)"`
}

func main() {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/grom_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 默认的表名都是结构体名称的复数形式，User 结构体默认创建的表为 users
	// db.SingularTable(true) 可以取消表名的复数形式，使表名和结构体名称一致
	db.AutoMigrate(&User{})

	// 添加唯一索引
	db.Model(&User{}).AddUniqueIndex("name_email", "id", "name", "email")

	// 插入记录
	//db.Create(&User{Name: "xlj", Age: 18, Email: "xlj@xlj.org"})

	gm := db.Create(&User{Name: "xll", Age: 18, Email: "xlj@xlj.org"})
	if err := gm.Error; err != nil {// 处理错误
		fmt.Println(err)
	}

	// 查看插入后的全部元素
	var users []User
	db.Find(&users)
	fmt.Printf("%+v\n", users)

	var user User
	// 查询一条记录
	db.First(&user, "name = ?", "xlj")
	fmt.Printf("%+v\n", user)

	// 更新记录
	db.Model(&user).Update("name", "xxj")
	fmt.Printf("%+v\n", user)

	// 删除记录
	db.Delete(&user)
}
```

### PostgreSQL

```go
gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
```

```go
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

```go
// https://github.com/go-gorm/postgres
db, err := gorm.Open(postgres.New(postgres.Config{
  DSN: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
  PreferSimpleProtocol: true, // disables implicit prepared statement usage
}), &gorm.Config{})
```

#### 自定义驱动

```go
import (
  _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
  "gorm.io/gorm"
)

db, err := gorm.Open(postgres.New(postgres.Config{
  DriverName: "cloudsqlpostgres",
  DSN: "host=project:region:instance user=postgres dbname=postgres password=password sslmode=disable",
})
```

#### 已存在的连接

```go
import (
  "database/sql"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

sqlDB, err := sql.Open("postgres", "mydb_dsn")
gormDB, err := gorm.Open(postgres.New(postgres.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```

### SQLite3

```go
gorm.Open("sqlite3", "gorm.db")
```

```go
import (
  "gorm.io/driver/sqlite" // Sqlite driver based on GGO
  // "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
  "gorm.io/gorm"
)

// github.com/mattn/go-sqlite3
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
```

### SQL Server

```go
gorm.Open("mssql", "sqlserver://username:password@localhost:1433?database=dbname")
```

```go
import (
  "gorm.io/driver/sqlserver"
  "gorm.io/gorm"
)

// github.com/denisenkom/go-mssqldb
dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
```

### Clickhouse

```go
import (
  "gorm.io/driver/clickhouse"
  "gorm.io/gorm"
)

func main() {
  dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
  db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

  // Auto Migrate
  db.AutoMigrate(&User{})
  // Set table options
  db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})

  // Insert
  db.Create(&user)

  // Select
  db.Find(&user, "id = ?", 10)

  // Batch Insert
  var users = []User{user1, user2, user3}
  db.Create(&users)
  // ...
}
```

### 连接池

```go
sqlDB, err := db.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)
```

## 声明数据表模型

```go
type User struct {
  gorm.Model
  Name         string
  Age          sql.NullInt64
  Birthday     *time.Time
  Email        string  `gorm:"type:varchar(100);unique_index"`
  // 设置字段大小为 255
  Role         string  `gorm:"size:255"`
  // 设置唯一且不为空
  MemberNumber *string `gorm:"unique;not null"`
  // 设置为自增序列
  Num          int     `gorm:"AUTO_INCREMENT"`
  // 创建 `addr` 索引
  Address      string  `gorm:"index:addr"`
  // 忽略字段
  IgnoreMe     int     `gorm:"-"`
}
```

### 约定

GORM 倾向于约定，而不是配置。默认情况下，GORM 使用 ID 作为主键，使用结构体名的 蛇形复数 作为表名，字段名的 蛇形 作为列名，并使用 CreatedAt、UpdatedAt 字段追踪创建、更新时间。

https://gorm.io/zh_CN/docs/conventions.html

### 字段级权限控制

可导出的字段在使用 GORM 进行 CRUD 时拥有全部的权限，此外，GORM 允许您用标签控制字段级别的权限。这样您就可以让一个字段的权限是只读、只写、只创建、只更新或者被忽略。

```go
type User struct {
  Name string `gorm:"<-:create"` // allow read and create
  Name string `gorm:"<-:update"` // allow read and update
  Name string `gorm:"<-"`        // allow read and write (create and update)
  Name string `gorm:"<-:false"`  // allow read, disable write permission
  Name string `gorm:"->"`        // readonly (disable write permission unless it configured )
  Name string `gorm:"->;<-:create"` // allow read and create
  Name string `gorm:"->:false;<-:create"` // createonly (disabled read from db)
  Name string `gorm:"-"`  // ignore this field when write and read with struct
  Name string `gorm:"migration"` // // ignore this field when migration
}
```

### 支持的结构体标签

声明 model 时，tag 是可选的，GORM 支持以下 tag。 tag 名大小写不敏感，但建议使用 camelCase 风格

| 标签名                 | 说明                                                         |
| :--------------------- | :----------------------------------------------------------- |
| column                 | 指定 db 列名                                                 |
| type                   | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：`not null`、`size`, `autoIncrement`… 像 `varbinary(8)` 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：`MEDIUMINT UNSIGNED not NULL AUTO_INCREMENT` |
| size                   | 指定列大小，例如：`size:256`                                 |
| primaryKey             | 指定列为主键                                                 |
| unique                 | 指定列为唯一                                                 |
| default                | 指定列的默认值                                               |
| precision              | 指定列的精度                                                 |
| scale                  | 指定列大小                                                   |
| not null               | 指定列为 NOT NULL                                            |
| autoIncrement          | 指定列为自动增长                                             |
| autoIncrementIncrement | 自动步长，控制连续记录之间的间隔                             |
| embedded               | 嵌套字段                                                     |
| embeddedPrefix         | 嵌入字段的列名前缀                                           |
| autoCreateTime         | 创建时追踪当前时间，对于 `int` 字段，它会追踪秒级时间戳，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoCreateTime:nano` |
| autoUpdateTime         | 创建/更新时追踪当前时间，对于 `int` 字段，它会追踪秒级时间戳，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoUpdateTime:milli` |
| index                  | 根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 [索引](https://gorm.io/zh_CN/docs/indexes.html) 获取详情 |
| uniqueIndex            | 与 `index` 相同，但创建的是唯一索引                          |
| check                  | 创建检查约束，例如 `check:age > 13`，查看 [约束](https://gorm.io/zh_CN/docs/constraints.html) 获取详情 |
| <-                     | 设置字段写入的权限， `<-:create` 只创建、`<-:update` 只更新、`<-:false` 无写入权限、`<-` 创建和更新权限 |
| ->                     | 设置字段读的权限，`->:false` 无读权限                        |
| -                      | 忽略该字段，`-` 无读写权限                                   |
| comment                | 迁移时为字段添加注释                                         |

### 定义字段字符集

```go
type User struct {
    gorm.Model
    Name `sql:"type:VARCHAR(5) CHARACTER SET utf8 COLLATE utf8_general_ci"`
}
```

## gorm 中的默认设置

### 开启日志输出打印 SQL 语句

```go
db.LogMode(true)
```

### gorm.Model

预先定义的一个基础模型，可以嵌入自定义的模型中：

```go
// gorm.Model definition
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}

// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` into model `User`
type User struct {
  gorm.Model
  Name string
}

// Declaring model w/o gorm.Model
type User struct {
  ID   int
  Name string
}
```

### 默认 ID 作为 Primary Key

```go
type User struct {
  ID   string // field named `ID` will be used as primary field by default
  Name string
}

// Set field `AnimalID` as primary field
type Animal struct {
  AnimalID int64 `gorm:"primary_key"`
  Name     string
  Age      int64
}
```

### 自定义表名而非结构体名

```go
// 默认表名为 `users`
type User struct {}

// 设置结构体 User 的表名为 `profiles`
func (User) TableName() string {
  return "profiles"
}
```

### 根据条件切换表名

```go
func (u User) TableName() string {
  if u.Role == "admin" {
    return "admin_users"
  } else {
    return "users"
  }
}
```

### 启用单数表名

```go
// 启用单数表名 user 而非默认的复数 users
db.SingularTable(true)
```

### 在执行语句时指定表名

```go
// Create `deleted_users` table with struct User's definition
db.Table("deleted_users").CreateTable(&User{})

// SELECT * FROM deleted_users;
var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)

// DELETE FROM deleted_users WHERE name = 'bill';
db.Table("deleted_users").Delete("", "name = bill")
// 物理删除
```

### 更改 grom 默认表名设置方式

```go
gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
  return "prefix_" + defaultTableName;
}
```

### 默认列名为字段的下划线式

```go
type User struct {
  ID        uint      // column name is `id`
  Name      string    // column name is `name`
  Birthday  time.Time // column name is `birthday`
  CreatedAt time.Time // column name is `created_at`
}
```

未指定情况下，grom 自动转换结构体字段名称为下划线式作为表的列名。

```go
// Overriding Column Name
type Animal struct {
  AnimalId    int64     `gorm:"column:beast_id"`         // set column name to `beast_id`
  Birthday    time.Time `gorm:"column:day_of_the_beast"` // set column name to `day_of_the_beast`
  Age         int64     `gorm:"column:age_of_the_beast"` // set column name to `age_of_the_beast`
}
```

## 常用功能

### 自动迁移数据表模型

自动迁移模式将保持表的更新，但是**不会更新索引以及现有列的类型**或删除未使用的列。

```go
// 同时迁移多个模型
db.AutoMigrate(&User{}, &Product{}, &Order{})
```

### 设置数据表字符集、引擎等

```go
// 比如修改表的字符类型 CHARSET=utf8mb4
db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").AutoMigrate(&User{})
```

### 检查表是否存在

```go
// 检查模型是否存在
db.HasTable(&User{})

// 检查表是否存在
db.HasTable("users")
```

### 增删改数据表的结构

```go
// 使用模型创建数据表
db.CreateTable(&User{})

// 增加参数创建
db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&User{})

// 删除表
db.DropTable(&User{})
db.DropTable("users")

// 模型和表名的混搭：两张表都将被删除
db.DropTableIfExists(&User{}, "products")

// 修改列，修改字段类型：不一定有效
db.Model(&User{}).ModifyColumn("description", "text")

// 删除列
db.Model(&User{}).DropColumn("description")

// 指定表名创建表
db.Table("deleted_users").CreateTable(&User{})

// 指定表名查询
var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
```

### 索引和约束

```go
// 添加外键
// 1st param : 外键字段
// 2nd param : 外键表(字段)
// 3rd param : ONDELETE
// 4th param : ONUPDATE
db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")

// 单个索引
db.Model(&User{}).AddIndex("idx_user_name", "name")

// 多字段索引
db.Model(&User{}).AddIndex("idx_user_name_age", "name", "age")

// 添加唯一索引，通常使用多个字段来唯一标识一条记录
db.Model(&User{}).AddUniqueIndex("idx_user_name", "name")
db.Model(&User{}).AddUniqueIndex("idx_user_name_age", "name", "id","email")

// 删除索引
db.Model(&User{}).RemoveIndex("idx_user_name")
```

## 创建

### 常用方式

```go
user := User{Name: "Who", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // pass pointer of data to Create

user.ID             // returns inserted data's primary key
result.Error        // returns error
result.RowsAffected // returns inserted records count
```

### 指定字段

```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("Who", 18, "2020-07-04 11:05:21.775")
```

### 忽略字段

```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")
```

### 批量插入

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

for _, user := range users {
  user.ID // 1,2,3
}
```

可以指定一次批量插入的数量，上面的默认是组合成一条SQL插入，下面的是按数量分割多次插入

```go
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)
```

```go

```

```go

```

### 创建钩子

BeforeSave, BeforeCreate, AfterSave, AfterCreate

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New()

  if u.Role == "admin" {
    return errors.New("invalid role")
  }
  return
}
```

跳过钩子

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)

DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)
```

```go

```

```go

```

```go

```



## 查询

> 暂时没有找到只获取一个值的方法。

### 基本查询

#### 根据主键获取第一条记录

```go
// SELECT * FROM users ORDER BY id LIMIT 1;
db.First(&user)
```

#### 根据主键获取指定的某条记录，仅当主键为整型时可用

```go
// SELECT * FROM users WHERE id = 10;
db.First(&user, 10)
```

#### 根据主键获取最后一条记录

```go
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
db.Last(&user)
```

#### 随机获取一条记录

```go
// SELECT * FROM users LIMIT 1;
db.Take(&user)
```

### 查询结果分析

```go
var user Model
gm := db.First(&user)

fmt.Println(gm.RowsAffected)
fmt.Println(gm.Error)
```

```
[2020-10-10 11:06:59]  [22.01ms]  SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1
[1 rows affected or returned ]
1
<nil>
```

#### 获取所有的记录

```go
// SELECT * FROM users;
users := []user
db.Find(&users)
```

### 条件查询

#### 通过 结构体 / Map 查询

```go
// 结构体方式
// select * from users where name = 'grom'
db.Where(&User{Name: "grom"}).First(&user)

// Map 方式
// select * from users where name = 'grom' and age = 20;
db.Where(map[string]interface{}{"name": "grom", "age": 20}).Find(&users)

// 主键的切片
// select * from users where id in (20,21,22);
db.Where([]int64{20, 21, 22}).Find(&users)
```

#### 通过 Where 条件语句查询

```go
// 使用条件获取一条记录 First() 方法
db.Where("name = ?", "grom").First(&user)

users := []user
// 获取全部记录 Find() 方法
db.Where("name = ?", "bill").Find(&users)

// 不等于 !=
db.Where("name <> ?", "bill").Find(&users)
db.Where("name != ?", "bill").Find(&users)

// AND
db.Where("name = ? AND age >= ?", "bill", "22").Find(&users)

// 时间比较
// select * from users where updated_at > '2020-03-06 00:00:00'
db.Where("updated_at > ?", "2020-03-06 00:00:00").Find(&users)

// BETWEEN
// select * from users where created_at between '2020-03-06 00:00:00' and '2020-03-14 00:00:00'
db.Where("created_at BETWEEN ? AND ?", "2020-03-06 00:00:00", "2020-03-14 00:00:00").Find(&users)
```

#### 通过 In 条件语句查询

```go
// IN
db.Where("name IN (?)", []string{"bill", "grom"}).Find(&users)
```

#### 通过 LIKE 条件语句查询

```go
// LIKE
db.Where("name LIKE ?", "%bill%").Find(&users)
```

#### 通过 Not 条件语句查询

```go
// select * from users where name != 'grom';
db.Not("name", "bill").First(&user)

// select * from users where name not in ("bill","grom");
db.Not("name", []string{"bill", "grom"}).Find(&users)

// select * from users where id not in (1,2,3)
db.Not([]int64{1,2,3}).First(&user)

// select * from users;
db.Not([]int64{}).First(&user)

// 原生 SQL
// select * from users where not(name = 'grom');
db.Not("name = ?", "grom").First(&user)

// struct 方式查询
// select * from users where name != 'grom'
db.Not(User{Name: "grom"}).First(&user)
```

#### 通过 Or 条件语句查询

```go
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)

// struct 方式
// SELECT * FROM users WHERE name = 'bill' OR name = 'grom';
db.Where("name = 'bill'").Or(User{Name: "grom"}).Find(&users)

// Map 方式
// SELECT * FROM users WHERE name = 'bill' OR name = 'grom';
db.Where("name = 'bill'").Or(map[string]interface{}{"name": "grom"}).Find(&users)
```

#### FirstOrCreate 不存在就插入记录

获取匹配的第一条记录，否则根据给定的条件创建一个新的记录（仅支持 struct 和 map 条件）。

```go
// 不存在就插入记录
db.FirstOrCreate(&user, User{Name: "non_existing"})

// select * from users where name = 'grom'
db.Where(User{Name: "grom"}).FirstOrCreate(&user)

// Attrs 参数：如果记录未找到，将使用参数创建 struct 和记录
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)

db.Where(User{Name: "grom"}).Attrs(User{Age: 30}).FirstOrCreate(&user)

// Assign 参数：不管记录是否找到，都将参数赋值给 struct 并保存至数据库
db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrCreate(&user)
```

### 子查询

```go
/*
select *
from users
where deleted_at is null
  and year > 2019
  and (amount > (select avg(amount) from users where year = 1970))
*/
db.Where("amount > (?)", db.Table("users").Select("avg(amount)").Where("year = ?", 1970).QueryExpr()).Where("year > ?", 2019).Find(&models)
```

> 其中 (?) 的括号不可以少。

### Select 部分字段查询

通常情况下，我们只想选择几个字段进行查询，指定你想从数据库中检索出的字段，默认会选择全部字段。

#### 映射结构体

```go
// SELECT name, age FROM users;
db.Select("name, age").Find(&users)

// SELECT name, age FROM users;
db.Select([]string{"name", "age"}).Find(&users)

```

#### 一行一行赋值

```go
// select coalesce(year, 1997) from users;
// 两者等价
/*
select case
           when year is not null then year
           else 1997
           end
from users;
*/
// 当数据表中 year 字段为 null 就默认为 1997
rows, _ := db.Table("users").Select("coalesce(year,?)", 1997).Rows()
for rows.Next() {
	var r int
	if err = rows.Scan(&r); err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}
```

```
2002
1979
1997
2014
1973
```

### 排序 Order

> 默认升序

#### 多字段排序

```go
// SELECT * FROM users ORDER BY age desc, name;
db.Order("age desc, name").Find(&users)

// SELECT * FROM users ORDER BY age desc, name;
db.Order("age desc").Order("name").Find(&users)
```

#### 覆盖排序

> 执行了两条 SQL 语句

```go
// SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL ORDER BY year desc;
// SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL ORDER BY `uuid`;
db.Order("desc desc").Find(&users1).Order("uuid", true).Find(&users2)
```

### 限制输出数量 LIMIT

```go
// SELECT * FROM users LIMIT 3;
db.Limit(3).Find(&users)

// 设置 -1 取消 LIMIT 条件
// SELECT * FROM users LIMIT 10;
// SELECT * FROM users;
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
// 执行了两条 SQL 语句
```

### 统计数量 COUNT

```go
// 通过指针赋值给 count
var count int

// 这个方法实际上执行了两条 SQL 语句
// SELECT count(*) FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((year = 1997) OR (first_name = 'skye'))
db.Where("year = ?", 1997).Or("first_name = ?", "skye").Find(&models).Count(&count)

// SELECT count(*) FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((year != 1997))
db.Model(&Model{}).Where("year != ?", 1997).Count(&count)

// SELECT count(*) FROM `users`
db.Table(Model{}.TableName()).Count(&count)

// SELECT count(distinct(year)) FROM `users`
db.Table("users").Select("count(distinct(year))").Count(&count)
```

### 分组 Group & Having

#### 单字段与多字段

> 不能加括号

```go
Group("one")
Group("one, two")
```

#### 一行一行获取

```go
type Result struct {
  Date  time.Time
  Total float64
}
var results []Result

// SELECT date(created_at) as date, sum(amount) as total FROM `users`   GROUP BY date(created_at)
rows, _ := db.Table("users").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
for rows.Next() {
	r := Result{}
	_ = rows.Scan(&r.Date, &r.Total)
	results = append(results, r)
}

// SELECT date(created_at) as date, sum(amount) as total FROM `users`   GROUP BY date(created_at) HAVING (sum(amount) > 100)
rows, _ := db.Table("users").Select("date(created_at) as date, sum(amount) as total").
	Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
for rows.Next() {
	r := Result{}
	_ = rows.Scan(&r.Date, &r.Total)
	results = append(results, r)
}
```

#### 一次性获取

```go
type Result struct {
  Date  time.Time
  Total int64
}

// SELECT date(created_at) as date, sum(amount) as total FROM `users`   GROUP BY date(created_at)
db.Table("users").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Scan(&results)
```

### JOIN 连接查询

```go
type Result struct {
	Name  string
	Code string
}
var results []Result

// 当 books 表中不存在该 user_id 时 code 字段会置为 null，而不会忽略
// SELECT users.last_name, books.code FROM `users` left join books on books.user_id = users.id
rows, err := db.Table("users").Select("users.last_name, books.code").Joins("left join books on books.user_id = users.id").Rows()
for rows.Next() {
	r := Result{}
	_ = rows.Scan(&r.Name, &r.Code)
	results = append(results, r)
}

// 当一个 user_id 存在多个 books.code 时会各自成为一行

// 一次性获取
db.Table("users").Select("users.last_name, books.code").Joins("left join books on books.user_id = users.id").Scan(results)

// 多连接及参数
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "bill@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

### Pluck 查询：获取一个列作为切片

Pluck，查询 model 中的一个列作为切片，可以存在重复值。

```go
var years []int
db.Find(&users).Pluck("year", &years)
// SELECT year FROM `users`  WHERE `users`.`deleted_at` IS NULL

db.Find(&users).Pluck("distinct year", &years)
// SELECT distinct year FROM `users`  WHERE `users`.`deleted_at` IS NULL

var names []string
db.Model(&User{}).Pluck("name", &names)

db.Table("users").Pluck("name", &names)
```

#### Pluck + Where 查询

```go
var year []int
// SELECT year FROM `users`  WHERE (email = 'w')
db.Table("users").Where("email = ?", "w").Pluck("year", &year)
```

### Scan 扫描：获取多个列的值

```go
type Result struct {
	UUID string
	Year int
}
var results []Result

// SELECT uuid, year FROM `users`  WHERE (year = 1997)
db.Table("users").Select("uuid, year").Where("year = ?", 1997).Scan(&results)
```

### 原生 SQL Scan

```go
// 原生 SQL
db.Raw("SELECT uuid, year FROM users WHERE year = ?", 1997).Scan(&result)
```

#### Exec

```go
db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at=? WHERE id IN ?", time.Now(), []int64{1,2,3})

// Exec with SQL Expression
db.Exec("update users set money=? where name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

> Raw 用于查询， Exec 用于其他命令

## 更新

### Save 更新所有字段

不适用于部分字段的更新。

```go
db.First(&user)

user.Year = "2020"
user.UserName = "admin"

// 执行 SQL 时全部字段都更新了
db.Save(&user)
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `created_at`   = '2016-01-21 12:41:03',
`updated_at`   = '2020-10-02 13:34:17',
`deleted_at`   = NULL,
`email`        = 'qsZrdnD@CCDQE.ru',
`password`     = 'iIqxlUrmDOsCFcgKdyUqzarzMgseaDtrpmgtmgEjXMlrmghOIZ',
`phone_number` = '865-429-3107',
`user_name`    = 'admin',
`first_name`   = 'Michale',
`last_name`    = 'Fahey',
`uuid`         = 'b18f56bb19d14ae793564fd1dad426ab',
`year`         = '2020',
`amount`       = 727371.6
WHERE `users`.`deleted_at` IS NULL
AND `users`.`id` = 1
```

### Update 更新指定字段

#### 根据主键更新单个属性

```go
user.ID = 3
db.Model(&user).Update("password", "12345")
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `password`   = '12345',
    `updated_at` = '2020-10-02 13:54:02'
WHERE `users`.`deleted_at` IS NULL
  AND `users`.`id` = 3
```

#### 根据条件更新单个属性

```go
db.Model(&Model{}).Where("year = ?", 1997).Update("password", "12345")
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `password`   = '12345',
    `updated_at` = '2020-10-02 13:55:16'
WHERE `users`.`deleted_at` IS NULL
  AND ((year = 1997))
```

#### 用 map 更新多个属性

> 只会更新其中有变化的属性

```go
user.ID = 20
db.Model(&user).Updates(map[string]interface{}{"password": "12345"})
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `password`   = '12345',
    `updated_at` = '2020-10-02 13:57:36'
WHERE `users`.`deleted_at` IS NULL
  AND `users`.`id` = 20
```

#### 用 struct 更新多个属性

> 只会更新其中有变化且为非零值的字段

```go
user.ID = 30
db.Model(&user).Updates(Model{Password: "passwd", Year: "1998"})
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `password`   = 'passwd',
    `updated_at` = '2020-10-02 13:59:08',
    `year`       = '1998'
WHERE `users`.`deleted_at` IS NULL
  AND `users`.`id` = 30
```

当使用 struct 更新时，GORM 只会更新那些非零值的字段，对于下面的操作，不会发生任何更新，"", 0, false 都是其类型的零值。

```go
db.Model(&user).Updates(User{Name: "", Age: 0, Actived: false})
```

### Select 更新部分字段

```go
db.Model(&user).Select("year").Updates(map[string]interface{}{"password": "123", "year": "2008"})
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `updated_at` = '2020-10-02 14:03:07',
    `year`       = '2008'
WHERE `users`.`deleted_at` IS NULL
```

### Omit 忽略更新部分字段

```go
db.Model(&user).Omit("year").Updates(map[string]interface{}{"password": "123", "year": "2008"})
```

执行的 SQL语句：

```sql
UPDATE `users`
SET `password`   = '123',
    `updated_at` = '2020-10-02 14:04:17'
WHERE `users`.`deleted_at` IS NULL
```

### 只更新指定字段，不更新自动更新字段

上面的更新操作会自动运行 model 的 `BeforeUpdate`，`AfterUpdate` 方法，来更新一些类似 `UpdatedAt` 的字段在更新时保存其 `Associations`，如果不想调用这些方法，可以使用 `UpdateColumn`，`UpdateColumns`。

#### 更新单个属性

> 类似于 `Update`

```go
db.Model(&user).UpdateColumn("name", "hello")
```

执行的 SQL语句：

```sql
update users set name = 'hello' where id = user.id;
```

#### 更新多个属性

> 类似于 `Updates`

```go
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
```

执行的 SQL语句：

```sql
update users set name = 'hello', age=18 where id = user.id;
```

### 批量更新

```go
db.Table("users").Where("id IN (?)", []int{1, 2, 3, 4}).Updates(map[string]interface{}{"password": "admin"})
```

执行的 SQL语句：

```sql
UPDATE `users` SET `password` = 'admin'  WHERE (id IN (1,2,3,4))
```

### 获取更新记录总数

```go
// 使用 `RowsAffected` 获取更新记录总数
db.Model(User{}).Updates(User{Name: "rustle", Age: 18}).RowsAffected
```

### 使用 SQL 计算表达式

```go
// update products set price = price*2+100 where id = product.id
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))

// update products set price = price*2+100 where id = product.id;
db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})

// update products set quantity = quantity-1 where id = product.id;
db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))

// update products set quantity = quantity -1 where id = product.id and quantity > 1
db.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
```

## 删除

### 删除记录

删除记录时，请确保主键字段有值，GORM 会通过主键去删除记录，如果主键为空，GORM 会删除该 model 的所有记录。

```go
// 删除现有记录
// UPDATE `users` SET `deleted_at`='2020-10-02 14:15:14'  WHERE `users`.`deleted_at` IS NULL AND `users`.`id` = 20
user.ID = 20
db.Delete(&user)

// 为删除 SQL 添加额外的 SQL 操作
// delete from emails where id = email.id OPTION (OPTIMIZE FOR UNKNOWN)
db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&email)
```

### 批量删除

```go
// UPDATE `users` SET `deleted_at`='2020-10-02 14:19:27'  WHERE `users`.`deleted_at` IS NULL AND ((last_name LIKE '%o%'))
db.Where("last_name LIKE ?", "%o%").Delete(Model{})

// UPDATE `users` SET `deleted_at`='2020-10-02 14:20:24'  WHERE `users`.`deleted_at` IS NULL AND ((last_name LIKE))'%d%'
db.Delete(Model{}, "last_name LIKE", "%d%")
```

### 软删除

如果一个 model 有 DeletedAt 字段，将自动获得软删除的功能。当调用 Delete 方法时， 记录不会真正的从数据库中被删除，只会将 DeletedAt 字段的值会被设置为当前时间。

在之前，可能会使用 isDelete 之类的字段来标记记录删除，不过在 gorm 中内置了 DeletedAt 字段，并且有相关 HOOK 来保证软删除。

```go
// UPDATE users SET deleted_at="2020-03-13 10:23" WHERE id = user.id;
db.Delete(&user)

// 批量删除
// 软删除的批量删除其实就是把 deleted_at 改成当前时间
// 并且在查询时无法查到，所以底层用的是 update 的 sql
db.Where("age = ?", 20).Delete(&User{})

// 查询记录时会忽略被软删除的记录
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
db.Where("age = 20").Find(&user)

// Unscoped 方法可以查询被软删除的记录
// SELECT * FROM users WHERE age = 20;
db.Unscoped().Where("age = 20").Find(&users)
```

### 物理删除

使用 `Unscoped().Delete()` 方法才是真正执行 SQL 中的 `delete` 语句.

```go
// Unscoped 方法可以物理删除记录
// DELETE FROM orders WHERE id=10;
db.Unscoped().Delete(&order)
```

## 事务处理

### 禁用默认事务

为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，可以在初始化时禁用它，这将获得大约 30%+ 性能提升。

```go
// 全局禁用
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})

// 持续会话模式
tx := db.Session(&Session{SkipDefaultTransaction: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

### 一般流程

```go
db.Transaction(func(tx *gorm.DB) error {
  // 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    // 返回任何错误都会回滚事务
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    return err
  }

  // 返回 nil 提交事务
  return nil
})
```

### 嵌套事务

```go
db.Transaction(func(tx *gorm.DB) error {
  tx.Create(&user1)

  tx.Transaction(func(tx2 *gorm.DB) error {
    tx2.Create(&user2)
    return errors.New("rollback user2") // Rollback user2
  })

  tx.Transaction(func(tx2 *gorm.DB) error {
    tx2.Create(&user3)
    return nil
  })

  return nil
})

// Commit user1, user3
```

### 手动事务

```go
// 开始事务
tx := db.Begin()

// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
tx.Create(...)

// ...

// 遇到错误时回滚事务
tx.Rollback()

// 否则，提交事务
tx.Commit()
```

### 手动事务示例

```go
func (a ArticleTags) Create(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, articleTag := range a {
		if err := tx.Create(&articleTag).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
```

### SavePoint、RollbackTo

GORM 提供了 SavePoint、Rollbackto 来提供保存点以及回滚至保存点，例如：

```go
tx := db.Begin()
tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commit user1
```

## 实体关联

### 自动创建、更新

```go
user := User{
  Name:            "jinzhu",
  BillingAddress:  Address{Address1: "Billing Address - Address 1"},
  ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
  Emails:          []Email{
    {Email: "jinzhu@example.com"},
    {Email: "jinzhu-2@example.com"},
  },
  Languages:       []Language{
    {Name: "ZH"},
    {Name: "EN"},
  },
}

db.Create(&user)
// BEGIN TRANSACTION;
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "languages" ("name") VALUES ('ZH'), ('EN') ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "user_languages" ("user_id","language_id") VALUES (111, 1), (111, 2) ON DUPLICATE KEY DO NOTHING;
// COMMIT;

db.Save(&user)
```

### Belongs To

belongs to 会与另一个模型建立了一对一的连接。 这种模型的每一个实例都“属于”另一个模型的一个实例。

创建时不存在会一起创建记录。

belongs to 会与另一个模型建立了一对一的连接。 这种模型的每一个实例都“属于”另一个模型的一个实例。

#### 重写外键

```go
type User struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"foreignKey:id"`  // 自定义外键
}

type Company struct {
	ID   int
	Name string
}
```

要定义一个 belongs to 关系，必须存在外键，默认的外键使用拥有者的类型名加上主字段名。

#### 重写引用

对于 belongs to 关系，GORM 通常使用拥有者的主字段作为外键的值。 对于上面的例子，它是 Company 的 ID 字段，**当将 user 分配给某个 company 时，GORM 会将 company 的 ID 保存到用户的 CompanyID 字段**

此外，也可以使用标签 references 手动更改它：

```go
type User struct {
	gorm.Model
	Name      string
	CompanyID string

	// v2 使用 Code 作为引用； v1 为 association_foreignkey
	Company   Company `gorm:"references:Code"`
}

type Company struct {
	ID   int
	Code string
	Name string
}
```

#### 外键约束

可以通过为标签 constraint 配置 OnUpdate、OnDelete 实现外键约束，在使用 GORM 进行**迁移时它会被创建**。

```go
type User struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
	ID   int
	Name string
}
```

### Has One

has one 与另一个模型建立一对一的关联，但它和一对一关系有些许不同。 这种关联表明一个模型的每个实例都包含或拥有另一个模型的一个实例。

例如，您的应用包含 user 和 credit card 模型，且每个 user 只能有一张 credit card。

```go
// User 有一张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCard CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

#### 重写外键

对于 has one 关系，同样必须存在外键字段。拥有者将把属于它的模型的主键保存到这个字段。

这个字段的名称通常由 has one 模型的类型加上其 主键 生成，对于上面的例子，它是 UserID。

为 user 添加 credit card 时，它会将 user 的 ID 保存到自己的 UserID 字段。

如果你想要使用另一个字段来保存该关系，你同样可以使用标签 foreignKey 来更改它，例如：

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"foreignKey:UserName"`
  // 使用 UserName 作为外键
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

#### 重写引用

默认情况下，拥有者实体会将 has one 对应模型的主键保存为外键，您也可以修改它，用另一个字段来保存，例如下个这个使用 Name 来保存的例子。

您可以使用标签 references 来更改它，例如：

```go
type User struct {
  gorm.Model
  Name       string     `gorm:"index"`
  CreditCard CreditCard `gorm:"foreignkey:UserName;references:name"`
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

#### 外键约束

你可以通过为标签 constraint 配置 OnUpdate、OnDelete 实现外键约束，在使用 GORM 进行迁移时它会被创建，例如：

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

#### 预加载

GORM 可以通过 Preload、Joins 预加载 belongs to 关联的记录。

```go
db.Preload("Orders").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

db.Preload("Orders").Preload("Profiles").Preload("Role").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
```

#### Joins 预加载

Preload 在一个单独查询中加载关联数据。而 Join Preload 会使用 inner join 加载关联数据。

```go
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, 1)
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, "users.name = ?", "jinzhu")
db.Joins("Company").Joins("Manager").Joins("Account").Find(&users, "users.id IN ?", []int{1,2,3,4,5})
```

Join Preload 适用于一对一的关系。

### Has Many

has many 与另一个模型建立了一对多的连接。 不同于 has one，拥有者可以有零或多个关联模型。

例如，您的应用包含 user 和 credit card 模型，且每个 user 可以有多张 credit card。

```go
// User 有多张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCards []CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

#### 重写外键

要定义 has many 关系，同样必须存在外键。 默认的外键名是拥有者的类型名加上其主键字段名

例如，要定义一个属于 User 的模型，则其外键应该是 UserID。

此外，想要使用另一个字段作为外键，您可以使用 foreignKey 标签自定义它：

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
  gorm.Model
  Number    string
  UserRefer uint
}
```

#### 重写引用

GORM 通常使用拥有者的主键作为外键的值。 对于上面的例子，它是 User 的 ID 字段。

为 user 添加 credit card 时，GORM 会将 user 的 ID 字段保存到 credit card 的 UserID 字段。

同样的，您也可以使用标签 references 来更改它，例如：

```go
type User struct {
  gorm.Model
  MemberNumber string
  CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard struct {
  gorm.Model
  Number     string
  UserNumber string
}
```

将 MemberNumber 的值赋值给 UserNumber。

#### 外键约束

你可以通过为标签 constraint 配置 OnUpdate、OnDelete 实现外键约束，在使用 GORM 进行迁移时它会被创建，例如：

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

## GORM 时区配置

### 系统默认时区

```go
conStr := "root:123456@tcp(192.168.3.93:33061)/zxd?charset=utf8mb4&parseTime=true&loc=Local"

db, err := gorm.Open("mysql", conStr)
if err != nil {
    log.Fatalf("%v", err)
}
```

### 设置时区

```go
conStr := "root:123456@tcp(192.168.3.93:33061)/zxd?charset=utf8mb4&parseTime=true&loc=Asia%2fShanghai"

db, err := gorm.Open("mysql", conStr)
if err != nil {
    log.Fatalf("%v", err)
}
```

`loc=Asia%2fShanghai`，gorm 配置链接字符串要求对 Loc 做 UrlEncode 处理

/ -> `%2f`

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

```go

```






## 官网资料

```
https://gorm.io/docs/query.html
```
