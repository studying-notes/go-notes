# GORM CRUD

用 GORM 实现创建、查询、更新和删除操作。现在 v2 版本已经在测试了，不过这里暂时只写一下 v1 版本。


```shell
# v1
go get -u github.com/jinzhu/gorm

# v2
go get gorm.io/gorm
```

- [GORM CRUD](#gorm-crud)
	- [功能预览](#功能预览)
	- [连接到数据库](#连接到数据库)
		- [引入驱动](#引入驱动)
		- [支持的数据库示例](#支持的数据库示例)
		- [MySQL](#mysql)
			- [MySQL 操作示例](#mysql-操作示例)
		- [PostgreSQL](#postgresql)
		- [SQLite3](#sqlite3)
		- [SQL Server](#sql-server)
	- [声明模型](#声明模型)
		- [支持的结构体标签](#支持的结构体标签)
	- [约定](#约定)
		- [gorm.Model](#gormmodel)
		- [默认 ID 作为 Primary Key](#默认-id-作为-primary-key)
		- [多元化表名](#多元化表名)
		- [指定表名](#指定表名)
		- [更改默认表名](#更改默认表名)
		- [下划线式列名](#下划线式列名)
	- [常用功能](#常用功能)
		- [自动迁移](#自动迁移)
		- [检查表](#检查表)
		- [增删改表的结构](#增删改表的结构)
		- [索引和约束](#索引和约束)
	- [查询](#查询)
		- [基本查询](#基本查询)
		- [结构体方式查询](#结构体方式查询)
		- [Where 条件查询](#where-条件查询)
	- [更多资料](#更多资料)

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

- `parseTime` - 需要解析 `time.Time` 时设置为 `True`
- `charset` - 为了完整支持 `UTF-8` 编码，需要设置为 `charset=utf8mb4`

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

### SQLite3

```go
gorm.Open("sqlite3", "gorm.db")
```

### SQL Server

```go
gorm.Open("mssql", "sqlserver://username:password@localhost:1433?database=dbname")
```

## 声明模型

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

### 支持的结构体标签

| Tag | Description|
| :-------------- | :-------------- |
| column | 指定列名 |
| type | 指定列的数据类型 `gorm:"type:varchar(25)"` |
| size | 指定列的大小，默认 255 `gorm:"size:50"` |
| PRIMARY_KEY | 指定列为 primary key |
| unique | 指定列为 unique |
| DEFAULT | 指定列的默认值 |
| PRECISION | 指定列精度 |
| not null | 指定列不为空 |
| AUTO_INCREMENT | 指定列为自增序列 |
| index | 创建带有或不带有名称的索引，相同名称将创建复合索引 |
| unique_index | 类似 `index`，但是唯一 |
| EMBEDDED | 将结构设置为嵌入式 |
| EMBEDDED_PREFIX | 设置嵌入式结构的前缀名称 |
| - | 忽略字段 |

## 约定

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

### 多元化表名

```go
// default table name is `users`
type User struct {} 

// Set User's table name to be `profiles`
func (User) TableName() string {
  return "profiles"
}

func (u User) TableName() string {
  if u.Role == "admin" {
    return "admin_users"
  } else {
    return "users"
  }
}

// Disable table name's pluralization, if set to
// true, `User`'s table name will be `user`
db.SingularTable(true)
```

### 指定表名

```go
// Create `deleted_users` table with struct User's definition
db.Table("deleted_users").CreateTable(&User{})

var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
//// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "jinzhu").Delete()
//// DELETE FROM deleted_users WHERE name = 'jinzhu';
```

### 更改默认表名

```go
gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
  return "prefix_" + defaultTableName;
}
```

### 下划线式列名

```go
type User struct {
  ID        uint      // column name is `id`
  Name      string    // column name is `name`
  Birthday  time.Time // column name is `birthday`
  CreatedAt time.Time // column name is `created_at`
}

// Overriding Column Name
type Animal struct {
  AnimalId    int64     `gorm:"column:beast_id"`         // set column name to `beast_id`
  Birthday    time.Time `gorm:"column:day_of_the_beast"` // set column name to `day_of_the_beast`
  Age         int64     `gorm:"column:age_of_the_beast"` // set column name to `age_of_the_beast`
}
```

## 常用功能

### 自动迁移

自动迁移模式将保持表的更新，但是不会更新索引以及现有列的类型或删除未使用的列。

```go
// 同时迁移多个模型
db.AutoMigrate(&User{}, &Product{}, &Order{})

// 创建表时增加相关参数
// 比如修改表的字符类型 CHARSET=utf8
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```

### 检查表

```go
// 检查模型是否存在
db.HasTable(&User{})

// 检查表是否存在
db.HasTable("users")
```

### 增删改表的结构

```go
// 使用模型创建
db.CreateTable(&User{})

// 增加参数创建
db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&User{})

// 删除表
db.DropTable(&User{})
db.DropTable("users")

// 模型和表名的混搭
db.DropTableIfExists(&User{}, "products")

// 修改列，修改字段类型
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

## 查询

### 基本查询

```go
// 根据主键查询第一条记录
// SELECT * FROM users ORDER BY id LIMIT 1;
db.First(&user)


// 随机获取一条记录
// SELECT * FROM users LIMIT 1;
db.Take(&user)


// 根据主键查询最后一条记录
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
db.Last(&user)


// 查询所有的记录
// SELECT * FROM users;
users := []user
db.Find(&users)


// 查询指定的某条记录，仅当主键为整型时可用
// SELECT * FROM users WHERE id = 10;
db.First(&user, 10)
```

### 结构体方式查询

```go
// 结构体方式
// select * from users where name = 'bgbiao.top'
db.Where(&User{Name: "bgbiao.top"}).First(&user)

// Map 方式
// select * from users where name = 'bgbiao.top' and age = 20;
db.Where(map[string]interface{}{"name": "bgbiao.top", "age": 20}).Find(&users)

// 主键的切片
// select * from users where id in (20,21,22);
db.Where([]int64{20, 21, 22}).Find(&users)
```

### Where 条件查询

```go
// 使用条件获取一条记录 First() 方法
db.Where("name = ?", "bgbiao.top").First(&user)

// 获取全部记录 Find() 方法
db.Where("name = ?", "jinzhu").Find(&users)

// 不等于
db.Where("name <> ?", "jinzhu").Find(&users)

// IN
db.Where("name IN (?)", []string{"jinzhu", "bgbiao.top"}).Find(&users)

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)

// Time
// select * from users where updated_at > '2020-03-06 00:00:00'
db.Where("updated_at > ?", lastWeek).Find(&users)

// BETWEEN
// select * from users where created_at between '2020-03-06 00:00:00' and '2020-03-14 00:00:00'
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
```

## 更多资料

```
https://zhuanlan.zhihu.com/p/113251066

https://gorm.io/docs/query.html
```
