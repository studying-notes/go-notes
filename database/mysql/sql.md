# Go 操作 MySQL

Go 语言中的 `database/sql` 标准库提供了保证 SQL 或类 SQL 数据库的泛用接口，并不提供具体的数据库驱动。使用 `database/sql` 包时必须注入一个数据库驱动。

```shell
go get -u github.com/go-sql-driver/mysql
```

```go
func Open(driverName, dataSourceName string) (*DB, error)
```

## 初始化连接

```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "user:password@tcp(localhost:3306)/dev"
    db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
    // 发生错误时，db 有可能为 nil，所以不能
    // 在 panic 之前注册关闭 db 的 defer 函数
	defer db.Close()
}
```

`Open` 函数只是验证其参数格式是否正确，实际上并不创建与数据库的连接，即使参数格式正确但是密码之类错误，也不会报错。如果要检查数据源是否真实有效，可以调用 `Ping` 方法测试。

返回的 DB 对象可以安全地被多个 goroutine 并发使用，并且维护其自己的空闲连接池。因此，Open 函数应该仅被调用一次，很少需要关闭这个 DB 对象。

- SetMaxOpenConns - 设置与数据库建立连接的最大数目

```go
func (db *DB) SetMaxOpenConns(n int)
```

如果 n 大于 0 且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制；如果 n<=0，不会限制最大开启连接数；默认为 0，即无限制。

- SetMaxIdleConns - 设置连接池中的最大闲置连接数

```go
func (db *DB) SetMaxIdleConns(n int)
```

如果 n 大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制；如果 n<=0，不会保留闲置连接。

## CRUD 增删查改

创建测试数据库：

```sql
create database sql_test;
use sql_test;

create table `user`
(
    `id`   bigint(20) not null auto_increment,
    `name` varchar(20) default '',
    `age`  int(11)     default '0',
    primary key (`id`)
) engine = innodb
  auto_increment = 1
  default charset = utf8mb4;

insert into user
values (1, 'mysql', 12);
```

### 单行查询

```go
func (db *DB) QueryRow(query string, args ...interface{}) *Row
```

单行查询 `db.QueryRow()` 执行一次查询，并期望返回最多一行结果。`QueryRow` 总是返回非 `nil` 的值，直到返回值的 `Scan` 方法被调用时，才会返回被延迟的错误。

```go
type user struct {
	id   int
	name string
	age  int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/sql_test"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	sqlStr := "select id, name, age from user where id=?"
	var u user
	_ = db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}
```

### 多行查询

```go
func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
```

多行查询 `db.Query()` 执行一次查询，返回多行结果，一般用于执行 `select` 命令。参数 `args` 表示 `query` 中的占位参数。

```go
func main() {
	sqlStr := "select id, name, age from user"
    rows, _ := db.Query(sqlStr)
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		_ = rows.Scan(&u.id, &u.name, &u.age)
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}
```

### 插入数据

插入、更新和删除操作都使用 `Exec` 方法。

```go
func (db *DB) Exec(query string, args ...interface{}) (Result, error)
```

`Exec` 执行一次命令（查询、删除、更新、插入等），返回的 `Result` 是对已执行的 SQL 命令的总结。参数 `args` 表示 `query` 中的占位参数。

```go
func main() {
	sqlStr := "insert into user(name, age) values (?,?)"
	_, _ = db.Exec(sqlStr, "redis", 12)
}
```

### 更新数据

```go
func main() {
	sqlStr := "update user set age=? where id = ?"
	_, _ = db.Exec(sqlStr, 12, 1)
}
```

### 删除数据

```go
func main() {
	sqlStr := "delete from user where id = ?"
	_, _ = db.Exec(sqlStr, 1)
}
```

## MySQL 预处理

### 什么是预处理？

普通 SQL 语句执行过程：

1. 客户端对 SQL 语句进行占位符替换得到完整的 SQL 语句。
2. 客户端发送完整 SQL 语句到 MySQL 服务端
3. MySQL 服务端执行完整的 SQL 语句并将结果返回给客户端。

预处理执行过程：

1. 把 SQL 语句分成两部分，命令部分与数据部分。
2. 先把命令部分发送给 MySQL 服务端，MySQL 服务端进行 SQL 预处理。
3. 然后把数据部分发送给 MySQL 服务端，MySQL 服务端对 SQL 语句进行占位符替换。
4. MySQL 服务端执行完整的 SQL 语句并将结果返回给客户端。

### 为什么要预处理？

1. 优化 MySQL 服务器重复执行 SQL 的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
2. 避免 SQL 注入问题。

### Go 实现 MySQL 预处理

`database/sql` 中使用下面的 `Prepare` 方法来实现预处理操作。

```go
func (db *DB) Prepare(query string) (*Stmt, error)
```

`Prepare` 方法会先将 SQL 语句发送给 MySQL 服务端，返回一个准备好的状态用于之后的查询和命令。返回值可以同时执行多个查询和命令。

```go
func main() {
	dsn := "root:root@tcp(localhost:3306)/sql_test"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	sqlStr := "select id, name, age from user where id > ?"
	stmt, _ := db.Prepare(sqlStr)
	defer stmt.Close()
	rows, _ := stmt.Query(0)
	defer rows.Close()

	for rows.Next() {
		var u user
		_ = rows.Scan(&u.id, &u.name, &u.age)
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}
```

### SQL 注入问题

任何时候都不应该自己拼接 SQL 语句！

```go
func main() {
	dsn := "root:root@tcp(localhost:3306)/sql_test"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	//inject := "redis' or 1=1#"
	inject := "redis' union select * from user #"
	//inject := "redis' and (select count(*) from user) <10 #"
	sqlStr := fmt.Sprintf("select * from user where name='%s'", inject)
	fmt.Println(sqlStr)
	var u user
	_ = db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	fmt.Printf("%+v\n", u)
}
```

## Go 实现 MySQL 事务

### 什么是事务？

事务：一个最小的不可再分的工作单元；通常一个事务对应一个完整的业务，同时这个完整的业务需要执行多次的 DML（insert、update、delete）语句共同联合完成。

在 MySQL 中只有使用了 `Innodb` 数据库引擎的数据库或表才支持事务。事务处理可以用来维护数据库的完整性，保证成批的 SQL 语句要么全部执行，要么全部不执行。

### 事务的 ACID

通常事务必须满足 4 个条件（ACID）：

- 原子性（Atomicity，或称不可分割性）- 一个事务中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节。事务在执行过程中发生错误，会被回滚到事务开始前的状态，就像这个事务从来没有执行过一样。

- 一致性（Consistency）- 在事务开始之前和事务结束以后，数据库的完整性没有被破坏。这表示写入的资料必须完全符合所有的预设规则，这包含资料的精确度、串联性以及后续数据库可以自发性地完成预定的工作。

- 隔离性（Isolation，又称独立性）- 数据库允许多个并发事务同时对其数据进行读写和修改的能力，隔离性可以防止多个事务并发执行时由于交叉执行而导致数据的不一致。事务隔离分为不同级别，包括读未提交（Read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（Serializable）。

- 持久性（Durability）- 事务处理结束后，对数据的修改就是永久的，即便系统故障也不会丢失。

### 事务相关方法

Go 语言中使用以下三个方法实现 MySQL 中的事务操作。 

1. 开始事务

```go
func (db *DB) Begin() (*Tx, error)
```

2. 提交事务

```go
func (tx *Tx) Commit() error
```

3. 回滚事务

```go
func (tx *Tx) Rollback() error
```

### 事务示例

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id   int
	name string
	age  int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/sql_test"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin failed, err:%v\n", err)
		return
	}
	sqlStr1 := "update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 1232)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 1233)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret2.RowsAffected() failed, err:%v\n", err)
		return
	}

	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("提交事务")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚")
	}
}
```
