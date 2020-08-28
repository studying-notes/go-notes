# GORM CRUD

用 GORM 实现创建、查询、更新和删除操作。现在 v2 版本已经在测试了，不过这里暂时只写一下 v1 版本。

```shell
# v1
go get -u github.com/bill/gorm

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
	- [声明数据表模型](#声明数据表模型)
		- [支持的结构体标签](#支持的结构体标签)
	- [gorm 中的默认设置](#gorm-中的默认设置)
		- [gorm.Model](#gormmodel)
		- [默认 ID 作为 Primary Key](#默认-id-作为-primary-key)
		- [自定义表名而非结构体名](#自定义表名而非结构体名)
		- [在执行语句时指定表名](#在执行语句时指定表名)
		- [更改 grom 默认表名设置方式](#更改-grom-默认表名设置方式)
		- [默认列名为字段的下划线式](#默认列名为字段的下划线式)
	- [常用功能](#常用功能)
		- [自动迁移](#自动迁移)
		- [检查表](#检查表)
		- [增删改表的结构](#增删改表的结构)
		- [索引和约束](#索引和约束)
	- [查询](#查询)
		- [基本查询](#基本查询)
		- [结构体方式查询](#结构体方式查询)
		- [Where 条件查询](#where-条件查询)
		- [Not 条件查询](#not-条件查询)
		- [Or 条件查询](#or-条件查询)
		- [FirstOrCreate](#firstorcreate)
		- [子查询](#子查询)
		- [字段查询 Select](#字段查询-select)
		- [排序 Order](#排序-order)
		- [限制输出数量 LIMIT](#限制输出数量-limit)
		- [统计数量 COUNT](#统计数量-count)
		- [分组 Group & Having](#分组-group--having)
		- [连接查询](#连接查询)
		- [Pluck 查询](#pluck-查询)
		- [Scan 扫描](#scan-扫描)
	- [更新](#更新)
		- [更新所有字段 Save](#更新所有字段-save)
		- [更新修改字段 Update](#更新修改字段-update)
		- [更新或者忽略某些字段](#更新或者忽略某些字段)
		- [无 HOOK 更新](#无-hook-更新)
		- [批量更新](#批量更新)
		- [使用 SQL 计算表达式](#使用-sql-计算表达式)
	- [删除](#删除)
		- [删除记录](#删除记录)
		- [批量删除](#批量删除)
		- [软删除](#软删除)
		- [物理删除](#物理删除)

## 功能预览

```go
package main

import (
	"fmt"
	"github.com/bill/gorm"
	_ "github.com/bill/gorm/dialects/sqlite"
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
import _ "github.com/bill/gorm/dialects/mysql"
// import _ "github.com/bill/gorm/dialects/postgres"
// import _ "github.com/bill/gorm/dialects/sqlite"
// import _ "github.com/bill/gorm/dialects/mssql"
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
	"github.com/bill/gorm"
	_ "github.com/bill/gorm/dialects/mysql"
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

## gorm 中的默认设置

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

// 根据条件切换表名
func (u User) TableName() string {
  if u.Role == "admin" {
    return "admin_users"
  } else {
    return "users"
  }
}

// 启用单数表名 user 而非默认的复数 users
db.SingularTable(true)
```

### 在执行语句时指定表名

```go
// Create `deleted_users` table with struct User's definition
db.Table("deleted_users").CreateTable(&User{})

var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
//// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "bill").Delete()
//// DELETE FROM deleted_users WHERE name = 'bill';
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
// select * from users where name = 'grom'
db.Where(&User{Name: "grom"}).First(&user)

// Map 方式
// select * from users where name = 'grom' and age = 20;
db.Where(map[string]interface{}{"name": "grom", "age": 20}).Find(&users)

// 主键的切片
// select * from users where id in (20,21,22);
db.Where([]int64{20, 21, 22}).Find(&users)
```

### Where 条件查询

```go
// 使用条件获取一条记录 First() 方法
db.Where("name = ?", "grom").First(&user)

// 获取全部记录 Find() 方法
db.Where("name = ?", "bill").Find(&users)

// 不等于 !=
db.Where("name <> ?", "bill").Find(&users)

// IN
db.Where("name IN (?)", []string{"bill", "grom"}).Find(&users)

// LIKE
db.Where("name LIKE ?", "%bill%").Find(&users)

// AND
db.Where("name = ? AND age >= ?", "bill", "22").Find(&users)

// Time
// select * from users where updated_at > '2020-03-06 00:00:00'
db.Where("updated_at > ?", lastWeek).Find(&users)

// BETWEEN
// select * from users where created_at between '2020-03-06 00:00:00' and '2020-03-14 00:00:00'
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
```

### Not 条件查询

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

### Or 条件查询

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

### FirstOrCreate

获取匹配的第一条记录，否则根据给定的条件创建一个新的记录（仅支持 struct 和 map 条件）。

```go
// 未找到,就插入记录
// if select * from users where name = 'non_existing') is null; insert into users(name) values("non_existing")
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
// SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (amount > (SELECT AVG(amount) FROM "orders"  WHERE (state = 'paid')));
db.Where("amount > ?", DB.Table("orders").Select("AVG(amount)").Where("state = ?", "paid").QueryExpr()).Find(&orders)
```

### 字段查询 Select

通常情况下，我们只想选择几个字段进行查询，指定你想从数据库中检索出的字段，默认会选择全部字段。

```go
// SELECT name, age FROM users;
db.Select("name, age").Find(&users)

// SELECT name, age FROM users;
db.Select([]string{"name", "age"}).Find(&users)

// SELECT COALESCE(age,'42') FROM users;
db.Table("users").Select("COALESCE(age,?)", 42).Rows()
```

### 排序 Order

```go
// SELECT * FROM users ORDER BY age desc, name;
db.Order("age desc, name").Find(&users)

// 多字段排序
// SELECT * FROM users ORDER BY age desc, name;
db.Order("age desc").Order("name").Find(&users)

// 覆盖排序
db.Order("age desc").Find(&users1).Order("age", true).Find(&users2)
```

### 限制输出数量 LIMIT

```go
// SELECT * FROM users LIMIT 3;
db.Limit(3).Find(&users)

// 设置 -1 取消 Limit 条件
// SELECT * FROM users LIMIT 10;
// SELECT * FROM users;
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
```

### 统计数量 COUNT

```go
// SELECT count(*) from USERS WHERE name = 'bill' OR name = 'grom';
db.Where("name = ?", "bill").Or("name = ?", "grom").Find(&users).Count(&count)

// select count(*) from users where name = 'grom'
db.Model(&User{}).Where("name = ?", "grom").Count(&count)

// SELECT count(*) FROM deleted_users;
db.Table("deleted_users").Count(&count)

// SELECT count( distinct(name) ) FROM deleted_users;
db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
```

### 分组 Group & Having

```go
rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
for rows.Next() {
  ...
}

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
for rows.Next() {
  ...
}

type Result struct {
  Date  time.Time
  Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

### 连接查询

```go
rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
  ...
}

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

// 多连接及参数
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "bill@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

### Pluck 查询

Pluck，查询 model 中的一个列作为切片，如果您想要查询多个列，您应该使用 Scan。

```go
var ages []int64
db.Find(&users).Pluck("age", &ages)

var names []string
db.Model(&User{}).Pluck("name", &names)

db.Table("deleted_users").Pluck("name", &names)
```

### Scan 扫描

```go
type Result struct {
  Name string
  Age  int
}

var result Result
db.Table("users").Select("name, age").Where("name = ?", "Antonio").Scan(&result)

// 原生 SQL
db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
```

## 更新

### 更新所有字段 Save

```go
db.First(&user)

user.Name = "grom"
user.Age = 100

// update users set name = 'grom',age=100 where id = user.id
db.Save(&user)
```

### 更新修改字段 Update

```go
// 更新单个属性，如果它有变化
// update users set name = 'hello' where id = user.id
db.Model(&user).Update("name", "hello")

// 根据给定的条件更新单个属性
// update users set name = 'hello' where active = true
db.Model(&user).Where("active = ?", true).Update("name", "hello")

// 使用 map 更新多个属性，只会更新其中有变化的属性
// update users set name = 'hello',age=18,actived=false where id = user.id
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})

// 使用 struct 更新多个属性，只会更新其中有变化且为非零值的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18})

// 当使用 struct 更新时，GORM 只会更新那些非零值的字段
// 对于下面的操作，不会发生任何更新，"", 0, false 都是其类型的零值
db.Model(&user).Updates(User{Name: "", Age: 0, Actived: false})
```

### 更新或者忽略某些字段

```go
// update users set name = 'hello' where id = user.id;
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})

// Omit() 方法用来忽略字段
// update users set age=18,actived=false where id = user.id
db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
```

### 无 HOOK 更新

上面的更新操作会自动运行 model 的 `BeforeUpdate`，`AfterUpdate` 方法，来更新一些类似 `UpdatedAt` 的字段在更新时保存其 `Associations`，如果不想调用这些方法，可以使用 `UpdateColumn`，`UpdateColumns`。

```go
// 更新单个属性，类似于 `Update`
// update users set name = 'hello' where id = user.id;
db.Model(&user).UpdateColumn("name", "hello")

// 更新多个属性，类似于 `Updates`
// update users set name = 'hello',age=18 where id = user.id;
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
```

### 批量更新

```go
// update users set name = 'hello',age=18 where id in (10,11)
db.Table("users").Where("id IN (?)", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})

// 使用 struct 更新时，只会更新非零值字段，若想更新所有字段，请用 map[string]interface{}
db.Model(User{}).Updates(User{Name: "hello", Age: 18})

// 使用 `RowsAffected` 获取更新记录总数
db.Model(User{}).Updates(User{Name: "hello", Age: 18}).RowsAffected
```

### 使用 SQL 计算表达式

```go
// update products set price = price*2+100 where id = product.id
DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))

// update products set price = price*2+100 where id = product.id;
DB.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})

// update products set quantity = quantity-1 where id = product.id;
DB.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))

// update products set quantity = quantity -1 where id = product.id and quantity > 1
DB.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
```

## 删除

### 删除记录

删除记录时，请确保主键字段有值，GORM 会通过主键去删除记录，如果主键为空，GORM 会删除该 model 的所有记录。

```go
// 删除现有记录
// delete from emails where id = email.id;
db.Delete(&email)

// 为删除 SQL 添加额外的 SQL 操作
// delete from emails where id = email.id OPTION (OPTIMIZE FOR UNKNOWN)
db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&email)
```

### 批量删除

```go
// delete from emails where email like '%bill%'
db.Where("email LIKE ?", "%bill%").Delete(Email{})

db.Delete(Email{}, "email LIKE ?", "%bill%")
```

### 软删除

如果一个 model 有 DeletedAt 字段，将自动获得软删除的功能。当调用 Delete 方法时， 记录不会真正的从数据库中被删除，只会将 DeletedAt 字段的值会被设置为当前时间。

在之前，我们可能会使用 isDelete 之类的字段来标记记录删除，不过在 gorm 中内置了 DeletedAt 字段，并且有相关 HOOK 来保证软删除。

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

```
## 官网资料

```
https://gorm.io/docs/query.html
```
